/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package database

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
	_ framework.DataSourceState = &DatabasesDataSourceModel{}
)

type DatabaseFilterModel struct {
	Organization *string  `tfsdk:"organization" hcl:"organization" cty:"organization"`
	Project      *string  `tfsdk:"project" hcl:"project" cty:"project"`
	Labels       []string `tfsdk:"labels" hcl:"labels" cty:"labels"`
}

type DatabaseNameModel struct {
	Organization string `tfsdk:"organization" hcl:"organization" cty:"organization"`
	Project      string `tfsdk:"project" hcl:"project" cty:"project"`
	Name         string `tfsdk:"name" hcl:"name" cty:"name"`
}

type DatabasesDataSourceModel struct {
	Filter    *DatabaseFilterModel `tfsdk:"filter" hcl:"filter" cty:"filter"`
	Databases []DatabaseNameModel  `tfsdk:"databases" hcl:"databases" cty:"databases"`
}

// GetDatabasesDataSourceSchema returns the schema for the databases (plural)
// data source. This has to be provided explicitly because there is no schema in
// the OpenAPI spec for the REST API that corresponds to it.
func GetDatabasesDataSourceSchema() *schema.Schema {
	sb := framework.NewSchemaBuilder().WithDescription("Data source for listing NuoDB databases created using the DBaaS Control Plane")
	sb.WithProjectScopeFilters("databases")
	sb.WithProjectScopeList("database", "databases").WithNameAttribute("database")
	return sb.Build()
}

func (state *DatabasesDataSourceModel) Read(ctx context.Context, client *openapi.Client) error {
	var organization, project string
	var labelFilter *string
	if state.Filter != nil {
		if state.Filter.Organization != nil {
			organization = *state.Filter.Organization
		}
		if state.Filter.Project != nil {
			project = *state.Filter.Project
		}
		if state.Filter.Labels != nil {
			labelFilterStr := strings.Join(state.Filter.Labels, ",")
			labelFilter = &labelFilterStr
		}
	}
	databases, err := helper.GetDatabases(ctx, client, organization, project, labelFilter, true)
	if err != nil {
		return err
	}
	state.Databases, err = GetDatabaseDataSourceResponse(databases)
	return err
}

func GetDatabaseDataSourceResponse(databases []string) ([]DatabaseNameModel, error) {
	var ret []DatabaseNameModel
	for _, db := range databases {
		parts := strings.Split(db, "/")
		if len(parts) != 3 {
			return nil, fmt.Errorf("Unexpected format for database name: %s", db)
		}
		ret = append(ret, DatabaseNameModel{
			Organization: parts[0],
			Project:      parts[1],
			Name:         parts[2],
		})
	}
	return ret, nil
}

func NewDatabasesDataSourceState() framework.DataSourceState {
	return &DatabasesDataSourceModel{}
}

func NewDatabasesDataSource() datasource.DataSource {
	return &framework.GenericDataSource{
		TypeName:       "databases",
		SchemaOverride: GetDatabasesDataSourceSchema(),
		Build:          NewDatabasesDataSourceState,
	}
}
