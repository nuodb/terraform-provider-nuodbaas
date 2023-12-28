/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"fmt"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/helper"

	nuodbaas_client "github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/client"
	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

var _ datasource.DataSourceWithConfigure = &databaseDataSource{}

func NewDatabaseDataSource() datasource.DataSource {
	return &databaseDataSource{}
}

type databaseDataSource struct {
	client *nuodbaas.APIClient
}

type databaseModel = model.DatabaseDataSourceModel

// Schema implements datasource.DataSource.
func (d *databaseDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "The organization that the database belongs to",
				Required: true,
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The name of the project for which database is created",
				Required: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the database",
				Required: true,
			},
			"tier": schema.StringAttribute{
				MarkdownDescription: "The Tier for the database. Cannot be updated once the database is created.",
				Computed: true,
			},
			"maintenance": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"is_disabled": schema.BoolAttribute{
						MarkdownDescription: "Whether the project or database should be shutdown",
						Computed: true,
					},
					"expires_in": schema.StringAttribute{
						MarkdownDescription: "The time until the project or database is disabled, e.g. 1d",
						Computed: true,
					},
					"expires_at": schema.StringAttribute{
						MarkdownDescription: "The time at which the project or database will be disabled",
						Computed: true,
					},
				},
			},
			"resource_version": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The version of the resource. When specified in a PUT request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates.",
			},
			"properties": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"archive_disk_size": schema.StringAttribute{
						MarkdownDescription: "The size of the archive volumes for the database. Can be only updated to increase the volume size",
						Computed: true,
					},
					"journal_disk_size": schema.StringAttribute{
						MarkdownDescription: "The size of the journal volumes for the database. Can be only updated to increase the volume size.",
						Computed: true,
					},
					"tier_parameters": schema.MapAttribute{
						Computed: true,
						ElementType: types.StringType,
					},
				},
			},
			"status": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"sql_end_point": schema.StringAttribute{
						MarkdownDescription: "The endpoint for SQL clients to connect to",
						Computed: true,
					},
					"ca_pem": schema.StringAttribute{
						MarkdownDescription: "The PEM-encoded certificate for SQL clients to verify database servers",
						Computed: true,
					},
				},
			},
		},
	}
}

// Metadata implements datasource.DataSource.
func (d *databaseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

// Read implements datasource.DataSource.
func (d *databaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state databaseModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	databaseResponseModel, err := nuodbaas_client.NewDatabaseClient(d.client, ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString()).GetDatabase()

	if err!=nil {
		resp.Diagnostics.AddError(
			"Error updating database",
			helper.GetApiErrorMessage(err, "Could not update database, unexpected error:"),
		)
		return
	}
	
	var databaseResp = model.DatabaseDataSourceModel{
		Organization: types.StringValue(*databaseResponseModel.Organization),
		Name: types.StringValue(*databaseResponseModel.Name),
		Tier: types.StringValue(*databaseResponseModel.Tier),
		ResourceVersion: types.StringValue(*databaseResponseModel.ResourceVersion),
		Project: types.StringValue(*databaseResponseModel.Project),
	}

	if databaseResponseModel.Properties != nil {
		propertiesModel := &model.DatabasePropertiesResourceModel{
			TierParameters: types.MapNull(types.StringType),
			ArchiveDiskSize: types.StringPointerValue(databaseResponseModel.Properties.ArchiveDiskSize),
			JournalDiskSize:  types.StringPointerValue(databaseResponseModel.Properties.JournalDiskSize),
		}

		if databaseResponseModel.Properties.TierParameters != nil {
			mapValue, diags := helper.ConvertMapToTfMap(databaseResponseModel.Properties.TierParameters)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			propertiesModel.TierParameters = mapValue
		}
		databaseResp.Properties = propertiesModel
	}

	if databaseResponseModel.Maintenance != nil {
		maintenanceModel := &model.MaintenanceDataSourceModel{
			IsDisabled: types.BoolPointerValue(databaseResponseModel.Maintenance.IsDisabled),
			ExpiresIn: types.StringPointerValue(databaseResponseModel.Maintenance.ExpiresIn),
		}

		if databaseResponseModel.Maintenance.ExpiresAtTime != nil {
			maintenanceModel.ExpiresAt = types.StringValue(databaseResponseModel.Maintenance.ExpiresAtTime.String())
		}
		databaseResp.Maintenance = maintenanceModel
	}

	if databaseResponseModel.Status != nil {
		statusModel := &model.StatusModel{
			SqlEndPoint:  types.StringPointerValue(databaseResponseModel.Status.SqlEndpoint),
			CaPem: types.StringPointerValue(databaseResponseModel.Status.CaPem),
		}
		databaseResp.Status = statusModel
	}

	state = databaseResp

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure implements datasource.DataSourceWithConfigure.
func (d *databaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

