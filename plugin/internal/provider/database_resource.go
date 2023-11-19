// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	nuodbaas_client "terraform-provider-nuodbaas/internal/client"
	"terraform-provider-nuodbaas/internal/model"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &DatabaseResource{}
	_ resource.ResourceWithImportState = &DatabaseResource{}
)

func NewDatabaseResource() resource.Resource {
	return &DatabaseResource{}
}

// DatabaseResource defines the resource implementation.
type DatabaseResource struct {
	client *openapi.APIClient
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
			"password": schema.StringAttribute{
				MarkdownDescription: "The name of the project for which database is created",
				Required:            true,
				Sensitive: true,
			},
			"tier": schema.StringAttribute{
				MarkdownDescription: "The Tier for the project. Cannot be updated once the project is created.",
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
			"archive_disk_size": schema.StringAttribute{
				MarkdownDescription: "The size of the archive volumes for the database. Can be only updated to increase the volume size",
				Computed: true,
				Optional: true,
			},
			"journal_disk_size": schema.StringAttribute{
				MarkdownDescription: "The size of the journal volumes for the database. Can be only updated to increase the volume size.",
				Optional: true,
			},

			// "properties": schema.SingleNestedAttribute{
			// 	Optional: true,
			// 	Computed: true,
			// 	Attributes: map[string]schema.Attribute{
			// 		"archive_disk_size": schema.StringAttribute{
			// 			MarkdownDescription: "The size of the archive volumes for the database. Can be only updated to increase the volume size",
			// 			Optional: true,
			// 		},
			// 		"journal_disk_size": schema.StringAttribute{
			// 			MarkdownDescription: "The size of the journal volumes for the database. Can be only updated to increase the volume size.",
			// 			Optional: true,
			// 		},
			// 	},
			// },
			
		},
	}
}

func (r *DatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openapi.APIClient)

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

	// var propertiesModel propertiesResourceModel
	var maintenanceModel maintenanceModel
	resp.Diagnostics.Append(state.Maintenance.As(ctx, &maintenanceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})...)
	if resp.Diagnostics.HasError() {
		return
	}
	// resp.Diagnostics.Append(state.Properties.As(ctx,&propertiesModel,  basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})...)
	databaseClient := nuodbaas_client.NewDatabaseClient(r.client, ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString())
	var archiveDiskSize, journalDiskSize string = "", ""
	if !state.ArchiveDiskSize.IsNull() && !state.ArchiveDiskSize.IsUnknown() {
		archiveDiskSize = state.ArchiveDiskSize.ValueString()
	}
	if !state.JournalDiskSize.IsNull() {
		journalDiskSize = state.JournalDiskSize.ValueString()
	}
	databaseBody := model.DatabaseCreateUpdateModel{
		Password: state.Password.ValueString(),
		Tier: state.Tier.ValueString(),
		ArchiveDiskSize: archiveDiskSize,
		JournalDiskSize: journalDiskSize,
	}

	_, err := databaseClient.CreateDatabase(maintenanceModel,databaseBody)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database",
			"Could not create database, unexpected error: "+ err.Error(),
		)
		return
	}

	getDatabaseModel, _, err := databaseClient.GetDatabase()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Project",
			"Could not get NuoDbaas project " + state.Name.ValueString()+" : " + err.Error(),
		)
		return
	}

	state.ResourceVersion = types.StringValue(*getDatabaseModel.ResourceVersion)

	if getDatabaseModel.Properties.ArchiveDiskSize != nil {
		state.ArchiveDiskSize = types.StringValue(*getDatabaseModel.Properties.ArchiveDiskSize)
	}

	if getDatabaseModel.Properties.JournalDiskSize != nil {
		state.JournalDiskSize = types.StringValue(*getDatabaseModel.Properties.JournalDiskSize)
	}
	// propertiesType := map[string]attr.Type{
	// 	"archive_disk_size" : types.StringType,
	// 	"journal_disk_size" : types.StringType,
	// }

	// propertiesValue := propertiesResourceModel{
	// 	ArchiveDiskSize : types.StringValue(archive_disk_size),
	// 	JournalDiskSize : types.StringValue(journal_disk_size),
	// }

	// objVal, diag := types.ObjectValueFrom(ctx, propertiesType, propertiesValue)
	// resp.Diagnostics.Append(diag...)
	// if resp.Diagnostics.HasError() {
	// 	tflog.Debug(ctx, "TAGGER error while converting to objVal")
	// 	return
	// }
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "TAGGER Last me error")
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
	databaseModel, _, err := nuodbaas_client.NewDatabaseClient(r.client, ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString()).GetDatabase()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Project",
			"Could not get NuoDbaas project " + state.Name.ValueString()+" : " + err.Error(),
		)
		return
	}

	state.ResourceVersion = types.StringValue(*databaseModel.ResourceVersion)

	journal_disk_size, archive_disk_size  := "" , ""

	if databaseModel.Properties.ArchiveDiskSize != nil {
		archive_disk_size = *databaseModel.Properties.ArchiveDiskSize
	}

	if databaseModel.Properties.JournalDiskSize != nil {
		journal_disk_size = *databaseModel.Properties.JournalDiskSize
	}

	// propertiesType := map[string]attr.Type{
	// 	"archive_disk_size" : types.StringType,
	// 	"journal_disk_size" : types.StringType,
	// }

	// propertiesValue := propertiesResourceModel{
	// 	ArchiveDiskSize : types.StringValue(archive_disk_size),
	// 	JournalDiskSize : types.StringValue(journal_disk_size),
	// }

	// objVal, diag := types.ObjectValueFrom(ctx, propertiesType,propertiesValue)
	// resp.Diagnostics.Append(diag...)
	// if resp.Diagnostics.HasError() {
	// 	tflog.Debug(ctx, "TAGGER error while converting to objVal read")
	// 	return
	// }
	state.ArchiveDiskSize = types.StringValue(archive_disk_size)
	state.JournalDiskSize = types.StringValue(journal_disk_size)

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

	// var propertiesModel propertiesResourceModel
	var maintenanceModel maintenanceModel
	resp.Diagnostics.Append(state.Maintenance.As(ctx, &maintenanceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})...)
	if resp.Diagnostics.HasError() {
		return
	}
	// resp.Diagnostics.Append(state.Properties.As(ctx, &propertiesModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})...)
	databaseClient := nuodbaas_client.NewDatabaseClient(r.client, ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString())
	var archiveDiskSize, journalDiskSize string = "", ""
	if !state.ArchiveDiskSize.IsNull() && !state.ArchiveDiskSize.IsUnknown() {
		archiveDiskSize = state.ArchiveDiskSize.ValueString()
	}
	if !state.JournalDiskSize.IsNull() {
		journalDiskSize = state.JournalDiskSize.ValueString()
	}
	databaseBody := model.DatabaseCreateUpdateModel{
		Password: state.Password.ValueString(),
		Tier: state.Tier.ValueString(),
		ArchiveDiskSize: archiveDiskSize,
		JournalDiskSize: journalDiskSize,
	}

	_, err := databaseClient.UpdateDatabase(maintenanceModel,databaseBody, state.ResourceVersion.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating database",
			fmt.Sprintf("Could not update database, unexpected error: %+v", err.Error()),
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

	_, err := r.client.DatabasesAPI.DeleteDatabase(ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString()).Execute()

	if err!=nil {
		resp.Diagnostics.AddError("Unable to delete project", 
			fmt.Sprintf("Unable to delete project %s, unexpected error: %v", 
			state.Name.ValueString(), err.Error()))
		return
	}
}

func (r *DatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
