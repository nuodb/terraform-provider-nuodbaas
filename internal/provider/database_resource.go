/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"fmt"
	"time"

	nuodbaas "github.com/nuodb/terraform-provider-nuodbaas/client"
	"github.com/nuodb/terraform-provider-nuodbaas/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/model"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.ResourceWithImportState = &DatabaseResource{}
)

func NewDatabaseResource() resource.Resource {
	return &DatabaseResource{}
}

// DatabaseResource defines the resource implementation.
type DatabaseResource struct {
	client *nuodbaas.APIClient
}

func (r *DatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *DatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A resource to create a new database." +
			" When creating a project and database in the same chart, make sure that the database resource has an explicit dependency on the project resource (for example, by using values from the project in the database, like in the examples).",
		MarkdownDescription: "A resource to create a new database." +
			"\n\n ~> **Note** When creating a project and database in the same chart, make sure that the database resource has an explicit dependency on the project resource (for example, by using values from the project in the database, like in the examples).",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Description:         "Name of the organization which this database belongs to (should match the organization of the project).",
				MarkdownDescription: "Name of the organization which this database belongs to (should match the organization of the project).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					framework.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the database.",
				MarkdownDescription: "Name of the database.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					framework.RequiresReplace(),
				},
			},
			"project": schema.StringAttribute{
				Description:         "The name of the project for which database belongs to.",
				MarkdownDescription: "The name of the project for which database belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					framework.RequiresReplace(),
				},
			},
			"dba_password": schema.StringAttribute{
				Description:         "The password for the DBA user. Can only be specified when creating a database.",
				MarkdownDescription: "The password for the DBA user. Can only be specified when creating a database.",
				Required:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					framework.RequiresReplace(),
				},
			},
			"tier": schema.StringAttribute{
				Description:         "The service tier for the database. If omitted, the project service tier is inherited.",
				MarkdownDescription: "The service tier for the database. If omitted, the project service tier is inherited.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					framework.UseStateForUnknown(),
				},
			},
			"maintenance": schema.SingleNestedAttribute{
				Description:         "Maintenance shutdown status of the database.",
				MarkdownDescription: "Maintenance shutdown status of the database.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					framework.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"is_disabled": schema.BoolAttribute{
						Description:         "Whether the project or database should be shutdown",
						MarkdownDescription: "Whether the project or database should be shutdown",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							framework.UseStateForUnknown(),
						},
					},
				},
			},
			"properties": schema.SingleNestedAttribute{
				Description:         "Database configuration properties.",
				MarkdownDescription: "Database configuration properties.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					framework.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"archive_disk_size": schema.StringAttribute{
						Description:         "The size of the archive volumes for the database. Can be only updated to increase the volume size.",
						MarkdownDescription: "The size of the archive volumes for the database. Can be only updated to increase the volume size.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							framework.UseStateForUnknown(),
						},
					},
					"journal_disk_size": schema.StringAttribute{
						Description:         "The size of the journal volumes for the database. Can be only updated to increase the volume size.",
						MarkdownDescription: "The size of the journal volumes for the database. Can be only updated to increase the volume size.",
						Optional:            true,
					},
					"tier_parameters": schema.MapAttribute{
						Description:         "Opaque parameters supplied to database service tier.",
						MarkdownDescription: "Opaque parameters supplied to database service tier.",
						ElementType:         types.StringType,
						Optional:            true,
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

func (r *DatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state model.DatabaseResourceModel
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

	err := helper.CreateDatabase(ctx, r.client, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database",
			helper.GetApiErrorMessage(err, "Could not create database, unexpected error:"),
		)
		return
	}

	latest, err := r.waitForDatabase(ctx, state.Organization, state.Project, state.Name)
	if err != nil && helper.IsTimeoutError(err) {
		resp.Diagnostics.AddError(
			"Timeout error",
			fmt.Sprintf("Timed out waiting for database %s/%s/%s", state.Organization, state.Project, state.Name))
		return
	}

	if !helper.ConvertResource(resp.Diagnostics, &latest, &state) {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.DatabaseResourceModel
	if !helper.ReadResource(ctx, resp.Diagnostics, req.State.Get, &state) {
		return
	}

	latest, err := helper.GetDatabase(ctx, r.client, state.Organization, state.Project, state.Name)
	if err != nil {
		if errObj := helper.GetErrorContentObj(err); errObj != nil {
			if errObj.GetStatus() == "HTTP 404 Not Found" {
				resp.State.RemoveResource(ctx)
				return
			}
		}
		resp.Diagnostics.AddError(
			"Error reading database",
			helper.GetApiErrorMessage(err, "Could not get database "+state.Name),
		)
		return
	}

	if !helper.ConvertResource(resp.Diagnostics, &latest, &state) {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

}

func (r *DatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state model.DatabaseResourceModel
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

	err := helper.UpdateDatabase(ctx, r.client, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating database",
			helper.GetApiErrorMessage(err, "Could not update database, unexpected error:"),
		)
		return
	}

	latest, err := r.waitForDatabase(ctx, state.Organization, state.Project, state.Name)
	if err != nil && helper.IsTimeoutError(err) {
		resp.Diagnostics.AddError(
			"Timeout error",
			fmt.Sprintf("Timed out waiting for database %s/%s/%s", state.Organization, state.Project, state.Name))
		return
	}

	if !helper.ConvertResource(resp.Diagnostics, &latest, &state) {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.DatabaseResourceModel
	if !helper.ReadResource(ctx, resp.Diagnostics, req.State.Get, &state) {
		return
	}

	err := helper.DeleteDatabase(ctx, r.client, state.Organization, state.Project, state.Name)
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete database",
			helper.GetApiErrorMessage(err, fmt.Sprintf("Unable to delete database %s, unexpected error:", state.Name)))
		return
	}
}

func (r *DatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//TODO: Does not work
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DatabaseResource) waitForDatabase(ctx context.Context, organization, project, database string) (*nuodbaas.DatabaseModel, error) {
	var waitTime = 1 * time.Second
	for {
		databaseModel, err := helper.GetDatabase(ctx, r.client, organization, project, database)
		if err != nil {
			return nil, err
		}
		expectedState := "Available"
		if databaseModel.Maintenance != nil && databaseModel.Maintenance.GetIsDisabled() {
			expectedState = "Stopped"
		}
		if databaseModel.Status != nil && databaseModel.Status.GetState() == expectedState {
			return databaseModel, nil
		}
		// TODO: Error out if the database is in a failed state?
		time.Sleep(waitTime)
	}
}
