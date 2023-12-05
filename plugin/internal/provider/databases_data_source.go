package provider

import (
	"context"
	"fmt"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/helper"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/model"

	nuodbaas_client "github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

var _ datasource.DataSourceWithConfigure = &databasesDataSource{}

func NewDatabasesDataSource() datasource.DataSource {
	return &databasesDataSource{}
}

type databasesDataSource struct {
	client *nuodbaas.APIClient
}

type databasesModel struct {
	Filter		*databaseFilterModel   `tfsdk:"filter"`
	Databases   []model.DatabasesDataSourceResponseModel  `tfsdk:"databases"`
}

// Schema implements datasource.DataSource.
func (d *databasesDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// "filter" : schema.SingleNestedAttribute{
			// 	Optional: true,
			// 	Attributes: map[string]schema.Attribute{
			// 		"organization" : schema.StringAttribute{
			// 			Optional: true,
			// 		},
			// 		"project" : schema.StringAttribute{
			// 			Optional: true,
			// 		},
			// 	},
			// },
			"databases": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes : map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"organization": schema.StringAttribute{
							Computed: true,
						},
						"project": schema.StringAttribute{
							Computed: true,
						},
					},

				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter" : schema.SingleNestedBlock{
				Attributes:  map[string]schema.Attribute{
					"organization" : schema.StringAttribute{
						Optional: true,
					},
					"project" : schema.StringAttribute{
						Optional: true,
					},
				},
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

	var (
		organization = ""
		project = ""
	)

	if state.Filter != nil && state.Filter.Organization.IsNull() && !state.Filter.Project.IsNull() {
		resp.Diagnostics.AddError(
			"Organization Missing",
			"Organization is required with project name to get databases",
		)
		return
	}

	if state.Filter != nil && !state.Filter.Organization.IsNull() {
		organization = state.Filter.Organization.ValueString()
	}

	if state.Filter != nil && !state.Filter.Project.IsNull() {
		project = state.Filter.Project.ValueString()
	}

	databaseClient := nuodbaas_client.NewDatabaseClient(d.client,ctx, organization, project, "")

	databases, httpResponse, err := databaseClient.GetDatabases()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting databases",
			"Could not get databases, unexpected error: "+ helper.GetHttpResponseErrorMessage(httpResponse, err),
		)
		return
	}
	// tflog.Debug(ctx, fmt.Sprintf("TAGGER projects are %+v", databases))

	state.Databases = helper.GetDatabaseDataSourceResponse(databases)

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

