/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func NewProjectDataSourceState() framework.DataSourceState {
	return &ProjectResourceModel{}
}

func NewProjectDataSource() datasource.DataSource {
	return &framework.GenericDataSource{
		TypeName:         "project",
		Description:      "Data source for exposing information about NuoDB projects created using the DBaaS Control Plane",
		GetOpenApiSchema: framework.GetProjectSchema,
		Build:            NewProjectDataSourceState,
	}
}
