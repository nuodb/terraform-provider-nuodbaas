/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewProjectDataSourceState() DataSourceState {
	return &ProjectResourceModel{}
}

func NewProjectDataSource() datasource.DataSource {
	return &GenericDataSource{
		resourceTypeName: "project",
		description:      "Data source for exposing information about NuoDB projects provisioned using the DBaaS Control Plane",
		getOpenApiSchema: framework.GetProjectSchema,
		build:            NewProjectDataSourceState,
	}
}
