// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/helper"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/model"

	nuodbaas_client "github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
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

type databaseResourceModel = model.DatabaseResourceModel

type propertiesResourceModel = model.DatabasePropertiesResourceModel

func (r *DatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *DatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Database Resource",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "Name of the organization for which project is created",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the database",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The name of the project for which database is created",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				
			},
			"dba_password": schema.StringAttribute{
				MarkdownDescription: "Database password. Cannot be updated once database is created",
				Required:            true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tier": schema.StringAttribute{
				MarkdownDescription: "The Tier for the project. Cannot be updated once the database is created.",
				Required:            true,
			},
			"maintenance": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"expires_in": schema.StringAttribute{
						MarkdownDescription: "The time until the project or database is disabled, e.g. 1d",
						Optional: true,
					},
					"is_disabled": schema.BoolAttribute{
						MarkdownDescription: "Whether the project or database should be shutdown",
						Optional: true,
					},
				},
			},
			"resource_version": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The version of the resource. When specified in a PUT request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates.",
				PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
			},
			"properties": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"archive_disk_size": schema.StringAttribute{
						MarkdownDescription: "The size of the archive volumes for the database. Can be only updated to increase the volume size",
						Optional: true,
						Computed: true,
					},
					"journal_disk_size": schema.StringAttribute{
						MarkdownDescription: "The size of the journal volumes for the database. Can be only updated to increase the volume size.",
						Optional: true,
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts {
				Create: true,
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
			fmt.Sprintf("Expected *openapi.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state databaseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var propertiesModel *propertiesResourceModel = state.Properties
	var maintenanceModel maintenanceModel

	resp.Diagnostics.Append(state.Maintenance.As(ctx, &maintenanceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("State value is"))
	
	createTimeout, diags:= state.Timeouts.Create(ctx, 30*time.Minute)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)

	defer cancel()
	
	databaseClient := nuodbaas_client.NewDatabaseClient(r.client, ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString())
	httpResponse, err := databaseClient.CreateDatabase(state, maintenanceModel, propertiesModel)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database",
			"Could not create database, unexpected error: "+ helper.GetHttpResponseErrorMessage(httpResponse, err),
		)
		return
	}

	var getDatabaseModel *nuodbaas.DatabaseModel
	databaseModel, httpResponse, err := databaseClient.GetDatabase()
	getDatabaseModel = databaseModel
	// for i := 0;i<15; i++ {
	// 	databaseModel, httpResponse, err := databaseClient.GetDatabase()
	// 	getDatabaseModel = databaseModel
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Error reading Database",
	// 			"Could not get NuoDbaas database " + state.Name.ValueString()+" : " + helper.GetHttpResponseErrorMessage(httpResponse, err),
	// 		)
	// 		return
	// 	}
	// 	if *getDatabaseModel.Status.Ready {
	// 		break
	// 	}
	// 	time.Sleep(10 * time.Second)
	// }

	state.ResourceVersion = types.StringValue(*getDatabaseModel.ResourceVersion)
	tflog.Debug(ctx, "TAGGER idhar tak aaya")


	propertiesValue := propertiesResourceModel{
		ArchiveDiskSize: types.StringValue(*getDatabaseModel.Properties.ArchiveDiskSize),
	}

	if getDatabaseModel.Properties.JournalDiskSize != nil {
		propertiesValue.JournalDiskSize = types.StringValue(*getDatabaseModel.Properties.JournalDiskSize)
	}

	// if getDatabaseModel.Properties.TierParameters != nil {
	// 	tierParameters := map[string]attr.Value{}
	// 	for k,v := range *getDatabaseModel.Properties.TierParameters {
	// 		tierParameters[k] = types.StringValue(v)
	// 	}
	// 	mapValue, diags := types.MapValue(types.StringType, tierParameters)
	// 	tflog.Debug(ctx, fmt.Sprintf("TAGGER idhar tak aaya again %+v", mapValue))
	// 	resp.Diagnostics.Append(diags...)
	// 	if resp.Diagnostics.HasError() {
	// 		return
	// 	}
	// 	propertiesValue.TierParameters = mapValue
	// }

	state.Properties = &propertiesValue

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state databaseResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	databaseModel, httpResponse, err := nuodbaas_client.NewDatabaseClient(r.client, ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString()).GetDatabase()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Project",
			"Could not get NuoDbaas project " + state.Name.ValueString()+" : " + helper.GetHttpResponseErrorMessage(httpResponse, err),
		)
		return
	}

	state.ResourceVersion = types.StringValue(*databaseModel.ResourceVersion)
	state.Tier = types.StringValue(*databaseModel.Tier)

	propertiesValue := propertiesResourceModel{}

	if databaseModel.Properties.ArchiveDiskSize != nil {
		propertiesValue.ArchiveDiskSize = types.StringValue(*databaseModel.Properties.ArchiveDiskSize)
	}

	if databaseModel.Properties.JournalDiskSize != nil {
		propertiesValue.JournalDiskSize = types.StringValue(*databaseModel.Properties.JournalDiskSize)
	}

	// if databaseModel.Properties.TierParameters != nil {
	// 	tierParameters := map[string]attr.Value{}
	// 	for k,v := range *databaseModel.Properties.TierParameters {
	// 		tierParameters[k] = types.StringValue(v)
	// 	}
	// 	mapValue, diags := types.MapValue(types.StringType, tierParameters)
	// 	resp.Diagnostics.Append(diags...)
	// 	if resp.Diagnostics.HasError() {
	// 		return
	// 	}
	// 	propertiesValue.TierParameters = mapValue
	// }

	state.Properties = &propertiesValue

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state databaseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var propertiesModel *propertiesResourceModel = state.Properties
	var maintenanceModel maintenanceModel
	resp.Diagnostics.Append(state.Maintenance.As(ctx, &maintenanceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})...)
	if resp.Diagnostics.HasError() {
		return
	}

	databaseClient := nuodbaas_client.NewDatabaseClient(r.client, ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString())
	httpResponse, err := databaseClient.UpdateDatabase(state, maintenanceModel, propertiesModel)

	if httpResponse.StatusCode == 409 {
		updateResponseObj, retryError, isUpdated := r.retryUpdate(ctx, state, maintenanceModel, propertiesModel)
		if !isUpdated {
			if retryError != nil {
				err = retryError
				httpResponse = updateResponseObj
			}
		}
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating database",
			fmt.Sprintf("Could not update database, unexpected error: %+v", helper.GetHttpResponseErrorMessage(httpResponse, err)),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state databaseResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResponse, err := r.client.DatabasesAPI.DeleteDatabase(ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString()).Execute()

	if err!=nil {
		resp.Diagnostics.AddError("Unable to delete database", 
			fmt.Sprintf("Unable to database project %s, unexpected error: %v", 
			state.Name.ValueString(), helper.GetHttpResponseErrorMessage(httpResponse, err)))
		return
	}
}

func (r *DatabaseResource) retryUpdate(ctx context.Context, state databaseResourceModel, maintenanceModel maintenanceModel, propertiesModel *model.DatabasePropertiesResourceModel) (*http.Response, error, bool) {
	databaseClient := nuodbaas_client.NewDatabaseClient(r.client, ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString())
	databaseModel, httpResponse, err := databaseClient.GetDatabase()
	if err != nil {
		return httpResponse, err, false
	}
	if *databaseModel.ResourceVersion != state.ResourceVersion.ValueString() {
		state.ResourceVersion = types.StringValue(*databaseModel.ResourceVersion)
		httpResponse, err = databaseClient.UpdateDatabase(state, maintenanceModel, propertiesModel)
		if err != nil {
			return httpResponse, err, false
		} else {
			return nil, nil, true
		}
	}
	return nil, nil, false
}

func (r *DatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
