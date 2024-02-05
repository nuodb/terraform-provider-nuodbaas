/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"fmt"

	nuodbaas "github.com/nuodb/terraform-provider-nuodbaas/client"
	"github.com/nuodb/terraform-provider-nuodbaas/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSourceWithConfigure = &projectsDataSource{}

func NewProjectsDataSource() datasource.DataSource {
	return &projectsDataSource{}
}

type projectsDataSource struct {
	client *nuodbaas.APIClient
}

type projectsModel struct {
	Filter   *projectFilterModel                `tfsdk:"filter"`
	Projects []model.ProjectDataSourceNameModel `tfsdk:"projects"`
}

type projectFilterModel struct {
	Organization *string `tfsdk:"organization"`
}

// Schema implements datasource.DataSource.
func (d *projectsDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A listing of projects that exist in NuoDB DBaaS.",
		MarkdownDescription: "A listing of projects that exist in NuoDB DBaaS.",
		Attributes: map[string]schema.Attribute{
			"projects": schema.ListNestedAttribute{
				Description:         "The databases that exist.",
				MarkdownDescription: "The databases that exist.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "Name of the project",
							MarkdownDescription: "Name of the project",
							Computed:            true,
						},
						"organization": schema.StringAttribute{
							Description:         "Name of the organization for which project is created",
							MarkdownDescription: "Name of the organization for which project is created",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Description:         "Filters to narrow the list of fetched projects.",
				MarkdownDescription: "Filters to narrow the list of fetched projects.",
				Attributes: map[string]schema.Attribute{
					"organization": schema.StringAttribute{
						Description:         "Only return databases in a given organization.",
						MarkdownDescription: "Only return databases in a given organization.",
						Optional:            true,
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
	if !helper.ReadResource(ctx, resp.Diagnostics, req.Config.Get, &state) {
		return
	}

	var organization string
	if state.Filter != nil {
		if state.Filter.Organization != nil {
			organization = *state.Filter.Organization
		}
	}
	projects, err := helper.GetProjects(ctx, d.client, organization, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting projects",
			helper.GetApiErrorMessage(err, "Could not get projects, unexpected error:"),
		)
		return
	}

	state.Projects, err = helper.GetProjectDataSourceResponse(projects)
	if err != nil {
		resp.Diagnostics.AddError(
			"Conversion Failure",
			"Could not get convert project names: "+err.Error())
		return
	}
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
			fmt.Sprintf("Expected *nuodbaas.APIClient, got: %T. Please report this issue to NuoDB.Support@3ds.com", req.ProviderData),
		)
		return
	}

	d.client = client
}
