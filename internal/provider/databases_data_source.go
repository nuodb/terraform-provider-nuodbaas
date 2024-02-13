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

var _ datasource.DataSourceWithConfigure = &databasesDataSource{}

func NewDatabasesDataSource() datasource.DataSource {
	return &databasesDataSource{}
}

type databasesDataSource struct {
	client *nuodbaas.APIClient
}

type databaseFilterModel struct {
	Organization *string `tfsdk:"organization"`
	Project      *string `tfsdk:"project"`
}

type databasesModel struct {
	Filter    *databaseFilterModel                `tfsdk:"filter"`
	Databases []model.DatabaseDataSourceNameModel `tfsdk:"databases"`
}

// Schema implements datasource.DataSource.
func (d *databasesDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "A listing of databases deployed in NuoDB DBaaS.",
		MarkdownDescription: "A listing of databases deployed in NuoDB DBaaS.",
		Attributes: map[string]schema.Attribute{
			"databases": schema.ListNestedAttribute{
				Description:         "The databases that exist.",
				MarkdownDescription: "The databases that exist.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "Name of the database.",
							MarkdownDescription: "Name of the database.",
							Computed:            true,
						},
						"organization": schema.StringAttribute{
							Description:         "The organization that the database belongs to.",
							MarkdownDescription: "The organization that the database belongs to.",
							Computed:            true,
						},
						"project": schema.StringAttribute{
							Description:         "The name of the project to which the database belongs.",
							MarkdownDescription: "The name of the project to which the database belongs.",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Description:         "Filters to narrow the list of fetched databases.",
				MarkdownDescription: "Filters to narrow the list of fetched databases.",
				Attributes: map[string]schema.Attribute{
					"organization": schema.StringAttribute{
						Description:         "Only return databases in a given organization.",
						MarkdownDescription: "Only return databases in a given organization.",
						Optional:            true,
					},
					"project": schema.StringAttribute{
						Description:         "Only return databases in a given project. If supplied, the `organization` must also be provided.",
						MarkdownDescription: "Only return databases in a given project. If supplied, the `organization` must also be provided.",
						Optional:            true,
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
	if !helper.ReadResource(ctx, resp.Diagnostics, req.Config.Get, &state) {
		return
	}

	var organization, project string
	if state.Filter != nil {
		if state.Filter.Organization != nil {
			organization = *state.Filter.Organization
		}
		if state.Filter.Project != nil {
			project = *state.Filter.Project
		}
	}
	databases, err := helper.GetDatabases(ctx, d.client, organization, project, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting databases",
			helper.GetApiErrorMessage(err, "Could not get databases, unexpected error:"),
		)
		return
	}

	state.Databases, err = helper.GetDatabaseDataSourceResponse(databases)
	if err != nil {
		resp.Diagnostics.AddError(
			"Conversion Failure",
			"Could not get convert database names: "+err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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
