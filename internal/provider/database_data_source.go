/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ DataSourceState = &DatabaseResourceModel{}
)

type DatabaseDataSourceModel openapi.DatabaseModel

func (state *DatabaseDataSourceModel) Read(ctx context.Context, client *openapi.Client) error {
	resp, err := client.GetDatabase(ctx, state.Organization, state.Project, state.Name)
	if err != nil {
		return err
	}
	return helper.ParseResponse(resp, state)
}

func NewDatabaseDataSourceState() DataSourceState {
	return &DatabaseDataSourceModel{}
}

func NewDatabaseDataSource() datasource.DataSource {
	return &GenericDataSource{
		resourceTypeName: "database",
		description:      "Data source for exposing information about NuoDB databases provisioned using the DBaaS Control Plane",
		getOpenApiSchema: framework.GetDatabaseDataSourceSchema,
		build:            NewDatabaseDataSourceState,
	}
}
