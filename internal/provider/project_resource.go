/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/nuodb/terraform-provider-nuodbaas/helper"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/model"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
		Description: "A resource to create new DBaaS projects. " +
			"Projects allow you to group databases. " +
			"Every databases must belong to a project.",
		MarkdownDescription: "A resource to create new DBaaS projects. " +
			"Projects allow you to group databases. " +
			"Every databases must belong to a project.",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Description:         "Name of the organization for which project is created",
				MarkdownDescription: "Name of the organization for which project is created",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					framework.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the project",
				MarkdownDescription: "Name of the project",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					framework.RequiresReplace(),
				},
			},
			"sla": schema.StringAttribute{
				Description:         "The SLA for the project. Cannot be updated once the project is created.",
				MarkdownDescription: "The SLA for the project. Cannot be updated once the project is created.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					framework.RequiresReplace(),
				},
			},
			"tier": schema.StringAttribute{
				Description:         "The service tier for the project",
				MarkdownDescription: "The service tier for the project",
				Required:            true,
			},
			"maintenance": schema.SingleNestedAttribute{
				Description: "Maintenance shutdown status of the project. " +
					"Shutting down a project also shuts down all databases belonging to it.",
				MarkdownDescription: "Maintenance shutdown status of the project. " +
					"Shutting down a project also shuts down all databases belonging to it.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					framework.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"is_disabled": schema.BoolAttribute{
						Description:         "Whether the project or database should be shutdown",
						MarkdownDescription: "Whether the project or database should be shutdown",
						Optional:            true,
					},
				},
			},
			"properties": schema.SingleNestedAttribute{
				Description:         "Project configuration properties.",
				MarkdownDescription: "Project configuration properties.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					framework.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"tier_parameters": schema.MapAttribute{
						Description:         "Opaque parameters supplied to project service tier.",
						MarkdownDescription: "Opaque parameters supplied to project service tier.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
						PlanModifiers: []planmodifier.Map{
							framework.UseStateForUnknown(),
						},
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
	if !helper.ReadResource(ctx, resp.Diagnostics, req.Plan.Get, &state) {
		return
	}

	timeout, diags := state.Timeouts.Create(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := helper.CreateProject(ctx, r.client, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			helper.GetApiErrorMessage(err, "Could not create project, unexpected error:"),
		)
		return
	}

	latest, err := r.waitForProject(ctx, state.Organization, state.Name)
	if err != nil && helper.IsTimeoutError(err) {
		resp.Diagnostics.AddError(
			"Timeout error",
			fmt.Sprintf("Timed out waiting for project %s/%s", state.Organization, state.Name))
		return
	}

	if !helper.ConvertResource(resp.Diagnostics, &latest, &state) {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.ProjectResourceModel
	if !helper.ReadResource(ctx, resp.Diagnostics, req.State.Get, &state) {
		return
	}

	latest, err := helper.GetProject(ctx, r.client, state.Organization, state.Name)
	if err != nil {
		if errObj := helper.GetErrorContentObj(err); errObj != nil {
			if errObj.GetStatus() == "HTTP 404 Not Found" {
				resp.State.RemoveResource(ctx)
				return
			}
		}
		resp.Diagnostics.AddError(
			"Error reading project",
			helper.GetApiErrorMessage(err, "Could not get project "+state.Name),
		)
		return
	}

	if !helper.ConvertResource(resp.Diagnostics, &latest, &state) {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state model.ProjectResourceModel
	if !helper.ReadResource(ctx, resp.Diagnostics, req.Plan.Get, &state) {
		return
	}

	timeout, diags := state.Timeouts.Update(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := helper.UpdateProject(ctx, r.client, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating project",
			helper.GetApiErrorMessage(err, "Could not update project, unexpected error:"),
		)
		return
	}

	latest, err := r.waitForProject(ctx, state.Organization, state.Name)
	if err != nil && helper.IsTimeoutError(err) {
		resp.Diagnostics.AddError(
			"Timeout error",
			fmt.Sprintf("Timed out waiting for project %s/%s", state.Organization, state.Name))
		return
	}

	if !helper.ConvertResource(resp.Diagnostics, &latest, &state) {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.ProjectResourceModel
	if !helper.ReadResource(ctx, resp.Diagnostics, req.State.Get, &state) {
		return
	}

	err := helper.DeleteProject(ctx, r.client, state.Organization, state.Name)
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete project",
			helper.GetApiErrorMessage(err, fmt.Sprintf("Unable to delete project %s, unexpected error:", state.Name)))
		return
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//TODO: Does not work
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ProjectResource) waitForProject(ctx context.Context, organization, project string) (*nuodbaas.ProjectModel, error) {
	var waitTime = 1 * time.Second
	for {
		projectModel, err := helper.GetProject(ctx, r.client, organization, project)
		if err != nil {
			return nil, err
		}
		expectedState := "Available"
		if projectModel.Maintenance != nil && projectModel.Maintenance.GetIsDisabled() {
			expectedState = "Stopped"
		}
		if projectModel.Status != nil && projectModel.Status.GetState() == expectedState {
			return projectModel, nil
		}
		// TODO: Error out if the project is in a failed state?
		time.Sleep(waitTime)
	}
}
