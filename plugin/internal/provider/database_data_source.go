package provider

import (
	"context"
	"fmt"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/helper"

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

type databaseFilterModel struct {
	Organization 	types.String   	`tfsdk:"organization"`
	Project      	types.String	`tfsdk:"project"`
}

// Schema implements datasource.DataSource.
func (d *databaseDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "Name of the organization for which project is created",
				Required: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the database",
				Required: true,
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The name of the project for which database is created",
				Required: true,
			},
			"tier": schema.StringAttribute{
				MarkdownDescription: "The Tier for the project. Cannot be updated once the database is created.",
				Optional: true,
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
			},
			"properties": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"archive_disk_size": schema.StringAttribute{
						MarkdownDescription: "The size of the archive volumes for the database. Can be only updated to increase the volume size",
						Optional: true,
					},
					"journal_disk_size": schema.StringAttribute{
						MarkdownDescription: "The size of the journal volumes for the database. Can be only updated to increase the volume size.",
						Optional: true,
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

	databaseResponseModel, httpResponse, err :=  d.client.DatabasesAPI.GetDatabase(ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString()).Execute()

	if err!=nil {
		resp.Diagnostics.AddError(
			"Error updating database",
			fmt.Sprintf("Could not update database, unexpected error: %+v", helper.GetHttpResponseErrorMessage(httpResponse, err)),
		)
		return
	}
	
	propertiesModel := &model.DatabasePropertiesResourceModel{}
	maintenanceModel := &model.MaintenanceModel{}

	if databaseResponseModel.Properties.ArchiveDiskSize != nil {
		propertiesModel.ArchiveDiskSize = types.StringValue(*databaseResponseModel.Properties.ArchiveDiskSize)
	}

	if databaseResponseModel.Properties.JournalDiskSize != nil {
		propertiesModel.JournalDiskSize = types.StringValue(*databaseResponseModel.Properties.JournalDiskSize)
	}

	if databaseResponseModel.Maintenance.ExpiresIn != nil {
		maintenanceModel.ExpiresIn =  types.StringValue(*databaseResponseModel.Maintenance.ExpiresIn)
	}

	if databaseResponseModel.Maintenance.IsDisabled != nil {
		maintenanceModel.IsDisabled =  types.BoolValue(*databaseResponseModel.Maintenance.IsDisabled)
	}

	var databaseResp = model.DatabaseDataSourceModel{
		Organization: types.StringValue(*databaseResponseModel.Organization),
		Name: types.StringValue(*databaseResponseModel.Name),
		Tier: types.StringValue(*databaseResponseModel.Tier),
		ResourceVersion: types.StringValue(*databaseResponseModel.ResourceVersion),
		Project: types.StringValue(*databaseResponseModel.Project),
		Properties: propertiesModel,
		Maintenance: maintenanceModel,
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
			fmt.Sprintf("Expected *openapi.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

