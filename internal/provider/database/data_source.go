/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package database

import (
	"context"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ framework.DataSourceState = &DatabaseResourceModel{}
)

type DatabaseDataSourceModel openapi.DatabaseModel

func (state *DatabaseDataSourceModel) Read(ctx context.Context, client *openapi.Client) error {
	resp, err := client.GetDatabase(ctx, state.Organization, state.Project, state.Name)
	if err != nil {
		return err
	}
	return helper.ParseResponse(resp, state)
}

func NewDatabaseDataSourceState() framework.DataSourceState {
	return &DatabaseDataSourceModel{}
}

func NewDatabaseDataSource() datasource.DataSource {
	return &framework.GenericDataSource{
		TypeName:         "database",
		Description:      "Data source for exposing information about NuoDB databases created using the DBaaS Control Plane",
		GetOpenApiSchema: framework.GetDatabaseDataSourceSchema,
		Build:            NewDatabaseDataSourceState,
	}
}
