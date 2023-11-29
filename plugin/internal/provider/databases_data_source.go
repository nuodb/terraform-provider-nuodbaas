package provider

import (
	"context"
	"fmt"
	"terraform-provider-nuodbaas/helper"

	nuodbaas "github.com/GIT_USER_ID/GIT_REPO_ID"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ datasource.DataSourceWithConfigure = &databasesDataSource{}

func NewDatabasesDataSource() datasource.DataSource {
	return &databasesDataSource{}
}

type databasesDataSource struct {
	client *nuodbaas.APIClient
}

type databasesModel struct {
	Filter		types.Object   `tfsdk:"filter"`
	Databases   []types.String  `tfsdk:"databases"`
}

// Schema implements datasource.DataSource.
func (d *databasesDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"filter" : schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"organization" : schema.StringAttribute{
						Required: true,
					},
					"project" : schema.StringAttribute{
						Required: true,
					},
				},
			},
			"databases": schema.ListAttribute{
				ElementType: types.StringType,
				Optional: true,
			},
		},
	}
}

// Metadata implements datasource.DataSource.
func (d *databasesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_databases"
}

// Read implements datasource.DataSource.
func (d *databasesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state databasesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var filterModel databaseFilterModel

	resp.Diagnostics.Append(state.Filter.As(ctx, &filterModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})...)

	if resp.Diagnostics.HasError() {
		return
	}

	databases, httpResponse, err := d.client.DatabasesAPI.GetDatabases(ctx, filterModel.Organization.ValueString(), filterModel.Project.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting databases",
			"Could not get databases, unexpected error: "+ helper.GetHttpResponseErrorMessage(httpResponse, err),
		)
		return
	}
	// tflog.Debug(ctx, fmt.Sprintf("TAGGER projects are %+v", databases))

	for _, database := range databases.Items {
		state.Databases = append(state.Databases, types.StringValue(database))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure implements datasource.DataSourceWithConfigure.
func (d *databasesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

