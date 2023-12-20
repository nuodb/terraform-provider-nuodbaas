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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

var _ datasource.DataSourceWithConfigure = &projectsDataSource{}

func NewProjectsDataSource() datasource.DataSource {
	return &projectsDataSource{}
}

type projectsDataSource struct {
	client *nuodbaas.APIClient
}

type projectsModel struct {
	Filter		 *projectFilterModel     				 `tfsdk:"filter"`
	Projects     []model.ProjectDataSourceResponseModel  `tfsdk:"projects"`
}

type projectFilterModel struct {
	Organization types.String     `tfsdk:"organization"`
}

// Schema implements datasource.DataSource.
func (d *projectsDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"projects": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes : map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"organization": schema.StringAttribute{
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
				},
			},
		},
	}
}

// Metadata implements datasource.DataSource.
func (d *projectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

// Read implements datasource.DataSource.
func (d *projectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state projectsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	filter := state.Filter

	if resp.Diagnostics.HasError() {
		return
	}

	var organization = ""

	if filter != nil && !filter.Organization.IsNull() {
		organization = filter.Organization.ValueString()
	}

	projectClient := nuodbaas_client.NewProjectClient(d.client,ctx,organization,"")

	projects, httpResponse, err := projectClient.GetProjects()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting projects",
			"Could not get projects, unexpected error: "+ helper.GetHttpResponseErrorMessage(httpResponse, err),
		)
		return
	}

	projectDataSourceResponseList := helper.GetProjectDataSourceResponse(projects)
	

	state.Projects = projectDataSourceResponseList

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

// Configure implements datasource.DataSourceWithConfigure.
func (d *projectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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