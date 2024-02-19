/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ framework.DataSourceState = &ProjectsDataSourceModel{}
)

type ProjectFilterModel struct {
	Organization *string  `tfsdk:"organization" hcl:"organization"`
	Labels       []string `tfsdk:"labels" hcl:"labels"`
}

type ProjectNameModel struct {
	Organization string `tfsdk:"organization" hcl:"organization"`
	Name         string `tfsdk:"name" hcl:"name"`
}

type ProjectsDataSourceModel struct {
	Filter   *ProjectFilterModel `tfsdk:"filter" hcl:"filter"`
	Projects []ProjectNameModel  `tfsdk:"projects" hcl:"projects"`
}

// GetProjectsDataSourceSchema returns the schema for the projects (plural) data
// source. This has to be provided explicitly because there is no schema in the
// OpenAPI spec for the REST API that corresponds to it.
func GetProjectsDataSourceSchema() *schema.Schema {
	sb := framework.NewSchemaBuilder().WithDescription("Data source for listing NuoDB projects created using the DBaaS Control Plane")
	sb.WithOrganizationScopeFilters("projects")
	sb.WithOrganizationScopeList("project", "projects").WithNameAttribute("project")
	return sb.Build()
}

// Read implements datasource.DataSource.
func (state *ProjectsDataSourceModel) Read(ctx context.Context, client *openapi.Client) error {
	var organization string
	var labelFilter *string
	if state.Filter != nil {
		if state.Filter.Organization != nil {
			organization = *state.Filter.Organization
		}
		if state.Filter.Labels != nil {
			labelFilterStr := strings.Join(state.Filter.Labels, ",")
			labelFilter = &labelFilterStr
		}
	}
	projects, err := helper.GetProjects(ctx, client, organization, labelFilter, true)
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

func NewProjectsDataSourceState() framework.DataSourceState {
	return &ProjectsDataSourceModel{}
}

func NewProjectsDataSource() datasource.DataSource {
	return &framework.GenericDataSource{
		TypeName:       "projects",
		SchemaOverride: GetProjectsDataSourceSchema(),
		Build:          NewProjectsDataSourceState,
	}
}
