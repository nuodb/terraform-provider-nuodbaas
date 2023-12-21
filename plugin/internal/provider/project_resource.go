/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"fmt"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/helper"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/model"

	nuodbaas_client "github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
)

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

// ProjectResource defines the resource implementation.
type ProjectResource struct {
	client *nuodbaas.APIClient
}

type projectResourceModel = model.ProjectResourceModel
type maintenanceModel = model.MaintenanceModel

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "Name of the organization for which project is created",
				Required: true,
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
				MarkdownDescription: "The SLA for the project. Cannot be updated once the project is created.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				
			},
			"tier": schema.StringAttribute{
				MarkdownDescription: "The Tier for the project. Cannot be updated once the project is created.",
				Required:            true,
			},
			"maintenance": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"is_disabled": schema.BoolAttribute{
						Optional: true,
					},
				},
			},
			"resource_version": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
			},
			"properties": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"tier_parameters": schema.MapAttribute{
						Optional: true,
						ElementType: types.StringType,
					},
				},
			},
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
			fmt.Sprintf("Expected *nuodbaas.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state projectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectClient := nuodbaas_client.NewProjectClient(r.client, ctx, state.Organization.ValueString(), state.Name.ValueString())
	httpResponse, err := projectClient.CreateProject(state, state.Maintenance)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+ helper.GetHttpResponseErrorMessage(httpResponse, err),
		)
		return
	}

	getProjectModel, httpResponse, err := projectClient.GetProject()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Project",
			"Could not get NuoDbaas project " + state.Name.ValueString()+" : " + helper.GetHttpResponseErrorMessage(httpResponse, err),
		)
		return
	}

	state.ResourceVersion = types.StringValue(*getProjectModel.ResourceVersion)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state projectResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectClient := nuodbaas_client.NewProjectClient(r.client, ctx, state.Organization.ValueString(), state.Name.ValueString())
	projectModel, httpResp, err := projectClient.GetProject()

	if err != nil {
		errorModel := helper.GetHttpResponseModel(httpResp)
		if errorModel != nil && errorModel.GetStatus() == "HTTP 404 Not Found" {
			resp.State.RemoveResource(ctx)
			return
		}
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Project",
			"Could not get NuoDbaas project " + state.Name.ValueString()+" : " + helper.GetHttpResponseErrorMessage(httpResp, err),
		)
		return
	}

	if projectModel.Maintenance != nil {
		var maintenance = state.Maintenance
		
		if projectModel.Maintenance.IsDisabled != nil {
			if (state.Maintenance != nil && state.Maintenance.IsDisabled.IsNull() && *projectModel.Maintenance.IsDisabled) || (state.Maintenance!= nil && !state.Maintenance.IsDisabled.IsNull()) {
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
	var data projectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	
	projectClient := nuodbaas_client.NewProjectClient(r.client, ctx, data.Organization.ValueString(), data.Name.ValueString())
	httpResponse, err := projectClient.UpdateProject(data, data.Maintenance)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating project",
			"Could not update project, unexpected error: "+ helper.GetHttpResponseErrorMessage(httpResponse, err),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state projectResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResponse, err :=  nuodbaas_client.NewProjectClient(r.client, ctx, state.Organization.ValueString(), state.Name.ValueString()).DeleteProject()

	if err!=nil {
		resp.Diagnostics.AddError("Error deleting project", 
			fmt.Sprintf("Unable to delete project %s, unexpected error: %v", 
			state.Name.ValueString(), helper.GetHttpResponseErrorMessage(httpResponse, err)))
		return
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
