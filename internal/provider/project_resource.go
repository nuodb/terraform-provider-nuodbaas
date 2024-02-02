/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/nuodb/terraform-provider-nuodbaas/helper"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/model"

	nuodbaas_client "github.com/nuodb/terraform-provider-nuodbaas/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	nuodbaas "github.com/nuodb/terraform-provider-nuodbaas/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.ResourceWithImportState = &ProjectResource{}
)

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

// ProjectResource defines the resource implementation.
type ProjectResource struct {
	client *nuodbaas.APIClient
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A resource to create new DBaaS projects. " +
			"Projects allow you to group databases. " +
			"Every databases must belong to a project.",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "Name of the organization for which project is created",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the project",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sla": schema.StringAttribute{
				Description:         "The SLA for the project. Cannot be updated once the project is created.",
				MarkdownDescription: "The SLA for the project. Cannot be updated once the project is created.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tier": schema.StringAttribute{
				Description:         "The service tier for the project",
				MarkdownDescription: "The service tier for the project",
				Required:            true,
			},
			"maintenance": schema.SingleNestedAttribute{
				MarkdownDescription: "Maintenance shutdown status of the project. " +
					"Shutting down a project also shuts down all databases belonging to it.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"is_disabled": schema.BoolAttribute{
						Description:         "Whether the project or database should be shutdown",
						MarkdownDescription: "Whether the project or database should be shutdown",
						Optional:            true,
					},
				},
			},
			"resource_version": schema.StringAttribute{
				Description:         "The version of the resource. When specified in a `PUT` request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates.",
				MarkdownDescription: "The version of the resource. When specified in a `PUT` request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"properties": schema.SingleNestedAttribute{
				MarkdownDescription: "Project configuration properties.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"tier_parameters": schema.MapAttribute{
						Description:         "Opaque parameters supplied to project service tier.",
						MarkdownDescription: "Opaque parameters supplied to project service tier.",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Update: true,
			}),
		},
	}
}

func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*nuodbaas.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *nuodbaas.APIClient, got: %T. Please report this issue to NuoDB.Support@3ds.com", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state model.ProjectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Set timeout
	timeout, diags := state.Timeouts.Create(ctx, 20*time.Minute)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	projectClient := nuodbaas_client.NewProjectClient(r.client, ctx, state.Organization.ValueString(), state.Name.ValueString())
	err := projectClient.CreateProject(state)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			helper.GetApiErrorMessage(err, "Could not create project, unexpected error:"),
		)
		return
	}

	getProjectModel, err := waitForProject(projectClient)

	if err != nil && helper.IsTimeoutError(err) {
		resp.Diagnostics.AddError("Timeout error", fmt.Sprintf("Unable to get project %+v in ready. You can go ahead and retry creating it", state.Name.ValueString()))
		projectClient = nuodbaas_client.NewProjectClient(r.client, context.Background(), state.Organization.ValueString(), state.Name.ValueString())
		err = projectClient.DeleteProject()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting failed project",
				helper.GetApiErrorMessage(err, "Could not clean up timed out project deploy:"),
			)
		}
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Project",
			helper.GetApiErrorMessage(err, fmt.Sprintf("Could not get NuoDbaas project %+v", state.Name.ValueString())),
		)
		return
	}

	state.ResourceVersion = types.StringValue(*getProjectModel.ResourceVersion)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.ProjectResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectClient := nuodbaas_client.NewProjectClient(r.client, ctx, state.Organization.ValueString(), state.Name.ValueString())
	projectModel, err := projectClient.GetProject()

	if err != nil {
		if errObj := helper.GetErrorContentObj(err); errObj != nil {
			if errObj.GetStatus() == "HTTP 404 Not Found" {
				resp.State.RemoveResource(ctx)
				return
			}
		}
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Project",
			helper.GetApiErrorMessage(err, fmt.Sprintf("Could not get NuoDbaas project %+v", state.Name.ValueString())),
		)
		return
	}

	if projectModel.Maintenance != nil {
		var maintenance = state.Maintenance

		if projectModel.Maintenance.IsDisabled != nil {
			if (state.Maintenance != nil && state.Maintenance.IsDisabled.IsNull() && *projectModel.Maintenance.IsDisabled) || (state.Maintenance != nil && !state.Maintenance.IsDisabled.IsNull()) {
				maintenance.IsDisabled = types.BoolValue(*projectModel.Maintenance.IsDisabled)
			}
		}
		state.Maintenance = maintenance
	}

	if projectModel.Properties != nil {
		propertiesModel := model.ProjectProperties{
			TierParameters: types.MapNull(types.StringType),
		}

		if len(*projectModel.Properties.TierParameters) != 0 {
			mapValue, diags := helper.ConvertMapToTfMap(projectModel.Properties.TierParameters)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			propertiesModel.TierParameters = mapValue
		}
		if len(*projectModel.Properties.TierParameters) != 0 {
			state.Properties = &propertiesModel
		}

	}

	state.Tier = types.StringValue(projectModel.Tier)
	state.ResourceVersion = types.StringValue(*projectModel.ResourceVersion)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data model.ProjectResourceModel

	// TODO: Refresh project version

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Set timeout
	timeout, diags := data.Timeouts.Create(ctx, 20*time.Minute)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	projectClient := nuodbaas_client.NewProjectClient(r.client, ctx, data.Organization.ValueString(), data.Name.ValueString())
	err := projectClient.UpdateProject(data)

	//TODO: Retry on CONCURRENT_UPDATE

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating project",
			helper.GetApiErrorMessage(err, "Could not update project, unexpected error:"),
		)
		return
	}

	_, err = waitForProject(projectClient)

	if err != nil && helper.IsTimeoutError(err) {
		resp.Diagnostics.AddError("Timeout error", fmt.Sprintf("Project %+v failed to become ready.", data.Name.ValueString()))
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Project",
			helper.GetApiErrorMessage(err, fmt.Sprintf("Could not get project status %+v", data.Name.ValueString())),
		)
		return
	}

	// TODO: Save other values?

	// TODO: Uncomment once we remove UseStateForUnknown from ResourceVersion in schema.
	// This requires fetching resource version at the start of the Update
	// data.ResourceVersion = types.StringValue(*projectModel.ResourceVersion)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.ProjectResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := nuodbaas_client.NewProjectClient(r.client, ctx, state.Organization.ValueString(), state.Name.ValueString()).DeleteProject()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting project",
			helper.GetApiErrorMessage(err, fmt.Sprintf("Unable to delete project %s, unexpected error:", state.Name.ValueString())),
		)
		return
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//TODO: Does not work
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func waitForProject(projectClient *nuodbaas_client.NuodbaasProjectClient) (*nuodbaas.ProjectModel, error) {
	var projectModel *nuodbaas.ProjectModel
	var err error
	var waitTime = 1
	for {
		projectModel, err = projectClient.GetProject()
		if err != nil {
			return nil, err
		}

		if projectModel.Status != nil {
			isDisabled := projectModel.Maintenance != nil && projectModel.Maintenance.IsDisabled != nil && *projectModel.Maintenance.IsDisabled

			if (!isDisabled && projectModel.Status.GetState() == "Available") || (isDisabled && projectModel.Status.GetState() == "Stopped") {
				break
			}
		}

		time.Sleep(time.Duration(waitTime) * time.Second)
		waitTime = helper.ComputeWaitTime(waitTime, 10)
	}
	return projectModel, nil
}
