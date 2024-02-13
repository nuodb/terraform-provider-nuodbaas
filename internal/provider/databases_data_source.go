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
	_ DataSourceState = &DatabasesDataSourceModel{}
)

type DatabaseFilterModel struct {
	Organization *string `tfsdk:"organization"`
	Project      *string `tfsdk:"project"`
}

type DatabaseNameModel struct {
	Organization string `tfsdk:"organization"`
	Project      string `tfsdk:"project"`
	Name         string `tfsdk:"name"`
}

type DatabasesDataSourceModel struct {
	Filter    *DatabaseFilterModel `tfsdk:"filter"`
	Databases []DatabaseNameModel  `tfsdk:"databases"`
}

// GetDatabasesDataSourceSchema returns the schema for the databases (plural)
// data source. This has to be provided explicitly because there is no schema in
// the OpenAPI spec for the REST API that corresponds to it.
func GetDatabasesDataSourceSchema() *schema.Schema {
	return &schema.Schema{
		Description:         "Data source for listing NuoDB databases provisioned using the DBaaS Control Plane",
		MarkdownDescription: "Data source for listing NuoDB databases provisioned using the DBaaS Control Plane",
		Attributes: map[string]schema.Attribute{
			"databases": schema.ListNestedAttribute{
				Description:         "The list of databases that satisfy the filter requirements",
				MarkdownDescription: "The list of databases that satisfy the filter requirements",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "The name of the database",
							MarkdownDescription: "The name of the database",
							Computed:            true,
						},
						"organization": schema.StringAttribute{
							Description:         "The organization that the database belongs to",
							MarkdownDescription: "The organization that the database belongs to",
							Computed:            true,
						},
						"project": schema.StringAttribute{
							Description:         "The project that the database belongs to",
							MarkdownDescription: "The project that the database belongs to",
							Computed:            true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				Description:         "Filters to apply to database list",
				MarkdownDescription: "Filters to apply to database list",
				Attributes: map[string]schema.Attribute{
					"organization": schema.StringAttribute{
						Description:         "The organization to return databases for",
						MarkdownDescription: "The organization to return databases for",
						Optional:            true,
					},
					"project": schema.StringAttribute{
						Description:         "The project to return databases for. If specified, the organization must also be specified.",
						MarkdownDescription: "The project to return databases for. If specified, the organization must also be specified.",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (state *DatabasesDataSourceModel) Read(ctx context.Context, client *openapi.Client) error {
	var organization, project string
	if state.Filter != nil {
		if state.Filter.Organization != nil {
			organization = *state.Filter.Organization
		}
		if state.Filter.Project != nil {
			project = *state.Filter.Project
		}
	}
	databases, err := helper.GetDatabases(ctx, client, organization, project, true)
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

func NewDatabasesDataSourceState() DataSourceState {
	return &DatabasesDataSourceModel{}
}

func NewDatabasesDataSource() datasource.DataSource {
	return &GenericDataSource{
		resourceTypeName: "databases",
		schema:           GetDatabasesDataSourceSchema(),
		build:            NewDatabasesDataSourceState,
	}
}
