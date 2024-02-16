/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ DataSourceState = &ProjectsDataSourceModel{}
)

type ProjectFilterModel struct {
	Organization *string `tfsdk:"organization"`
}

type ProjectNameModel struct {
	Organization string `tfsdk:"organization"`
	Name         string `tfsdk:"name"`
}

type ProjectsDataSourceModel struct {
	Filter   *ProjectFilterModel `tfsdk:"filter"`
	Projects []ProjectNameModel  `tfsdk:"projects"`
}

// GetProjectsDataSourceSchema returns the schema for the projects (plural) data
// source. This has to be provided explicitly because there is no schema in the
// OpenAPI spec for the REST API that corresponds to it.
func GetProjectsDataSourceSchema() *schema.Schema {
	return &schema.Schema{
		Description:         "Data source for listing NuoDB projects provisioned using the DBaaS Control Plane",
		MarkdownDescription: "Data source for listing NuoDB projects provisioned using the DBaaS Control Plane",
		Attributes: map[string]schema.Attribute{
			"projects": schema.ListNestedAttribute{
				Description:         "The list of projects that satisfy the filter requirements",
				MarkdownDescription: "The list of projects that satisfy the filter requirements",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"organization": schema.StringAttribute{
							Description:         "The name of the organization the project belongs to",
							MarkdownDescription: "The name of the organization the project belongs to",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description:         "The name of the project",
							MarkdownDescription: "The name of the project",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Description:         "Filters to apply to project list",
				MarkdownDescription: "Filters to apply to project list",
				Attributes: map[string]schema.Attribute{
					"organization": schema.StringAttribute{
						Description:         "The organization to return projects for",
						MarkdownDescription: "The organization to return projects for",
						Optional:            true,
					},
				},
			},
		},
	}
}

// Read implements datasource.DataSource.
func (state *ProjectsDataSourceModel) Read(ctx context.Context, client *openapi.Client) error {
	var organization string
	if state.Filter != nil {
		if state.Filter.Organization != nil {
			organization = *state.Filter.Organization
		}
	}
	projects, err := helper.GetProjects(ctx, client, organization, true)
	if err != nil {
		return err
	}
	state.Projects, err = GetProjectDataSourceResponse(projects)
	return err
}

func GetProjectDataSourceResponse(projects []string) ([]ProjectNameModel, error) {
	var ret []ProjectNameModel
	for _, project := range projects {
		parts := strings.Split(project, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Unexpected format for project name: %s", project)
		}
		ret = append(ret, ProjectNameModel{
			Organization: parts[0],
			Name:         parts[1],
		})
	}
	return ret, nil
}

func NewProjectsDataSourceState() DataSourceState {
	return &ProjectsDataSourceModel{}
}

func NewProjectsDataSource() datasource.DataSource {
	return &GenericDataSource{
		resourceTypeName: "projects",
		schema:           GetProjectsDataSourceSchema(),
		build:            NewProjectsDataSourceState,
	}
}
