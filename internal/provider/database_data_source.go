/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewDatabaseDataSourceState() DataSourceState {
	return &DatabaseResourceModel{}
}

func NewDatabaseDataSource() datasource.DataSource {
	return &GenericDataSource{
		resourceTypeName: "database",
		description:      "Data source for exposing information about NuoDB databases provisioned using the DBaaS Control Plane",
		getOpenApiSchema: framework.GetDatabaseSchema,
		build:            NewDatabaseDataSourceState,
	}
}
