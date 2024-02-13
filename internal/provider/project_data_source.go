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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSourceWithConfigure = &projectDataSource{}

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

type projectDataSource struct {
	client *nuodbaas.APIClient
}

// Schema implements datasource.DataSource.
func (d *projectDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "The state of a given project.",
		MarkdownDescription: "The state of a given project.",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Description:         "Name of the organization for which project is created",
				MarkdownDescription: "Name of the organization for which project is created",
				Required:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Name of the project",
				MarkdownDescription: "Name of the project",
				Required:            true,
			},
			"sla": schema.StringAttribute{
				Description:         "The SLA for the project. Cannot be updated once the project is created.",
				MarkdownDescription: "The SLA for the project. Cannot be updated once the project is created.",
				Computed:            true,
			},
			"tier": schema.StringAttribute{
				Description:         "The service tier for the project",
				MarkdownDescription: "The service tier for the project",
				Computed:            true,
			},
			"maintenance": schema.SingleNestedAttribute{
				Description: "Maintenance shutdown status of the project. " +
					"Shutting down a project also shuts down all databases belonging to it.",
				MarkdownDescription: "Maintenance shutdown status of the project. " +
					"Shutting down a project also shuts down all databases belonging to it.",
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"is_disabled": schema.BoolAttribute{
						Description:         "Whether the project or database should be shutdown",
						MarkdownDescription: "Whether the project or database should be shutdown",
						Computed:            true,
					},
				},
			},
			"resource_version": schema.StringAttribute{
				Description:         "The version of the resource. When specified in a `PUT` request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates.",
				MarkdownDescription: "The version of the resource. When specified in a `PUT` request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates.",
				Computed:            true,
			},
			"properties": schema.SingleNestedAttribute{
				Description:         "Project configuration properties.",
				MarkdownDescription: "Project configuration properties.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"tier_parameters": schema.MapAttribute{
						Description:         "Opaque parameters supplied to project service tier.",
						MarkdownDescription: "Opaque parameters supplied to project service tier.",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
		},
	}
}

// Metadata implements datasource.DataSource.
func (d *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Read implements datasource.DataSource.
func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state model.ProjectDataSourceModel
	if !helper.ReadResource(ctx, resp.Diagnostics, req.Config.Get, &state) {
		return
	}

	project, err := helper.GetProject(ctx, d.client, state.Organization, state.Name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting project",
			helper.GetApiErrorMessage(err, "Could not update project:"),
		)
		return
	}

	if !helper.ConvertResource(resp.Diagnostics, &project, &state) {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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
			fmt.Sprintf("Expected *nuodbaas.APIClient, got: %T. Please report this issue to NuoDB.Support@3ds.com", req.ProviderData),
		)
		return
	}

	d.client = client
}
