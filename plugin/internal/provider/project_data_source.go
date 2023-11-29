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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSourceWithConfigure = &projectDataSource{}

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

type projectDataSource struct {
	client *nuodbaas.APIClient
}

type projectsModel struct {
	Filter		 types.Object     `tfsdk:"filter"`
	Projects     []types.String   `tfsdk:"projects"`
}

type projectFilterModel struct {
	Organization types.String     `tfsdk:"organization"`
}

// Schema implements datasource.DataSource.
func (d *projectDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"projects": schema.ListAttribute{
				Computed: true,
				ElementType: types.StringType,
			},
			"filter" : schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"organization" : schema.StringAttribute{
						Required: true,
					},
				},
			}, 
		},
	}
}

// Metadata implements datasource.DataSource.
func (d *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

// Read implements datasource.DataSource.
func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state projectsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var filter projectFilterModel

	resp.Diagnostics.Append(state.Filter.As(ctx, &filter, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})...)

	if resp.Diagnostics.HasError() {
		return
	}

	projects, httpResponse, err := d.client.ProjectsAPI.GetProjects(ctx, filter.Organization.ValueString()).Execute()
	
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting projects",
			"Could not get projects, unexpected error: "+ helper.GetHttpResponseErrorMessage(httpResponse, err),
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("TAGGER projects are %+v", projects))

	for _, project := range projects.Items {
		state.Projects = append(state.Projects, types.StringValue(project))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure implements datasource.DataSourceWithConfigure.
func (d *projectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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