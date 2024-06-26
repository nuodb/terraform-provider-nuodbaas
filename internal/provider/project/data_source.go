// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package project

import (
	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func GetProjectDataSourceAttributes() (map[string]schema.Attribute, error) {
	return framework.GetDataSourceAttributes("ProjectModel")
}

func NewProjectDataSourceState() framework.DataSourceState {
	return &ProjectResourceModel{}
}

func NewProjectDataSource() datasource.DataSource {
	return &framework.GenericDataSource{
		TypeName:                "project",
		Description:             "Data source for exposing information about NuoDB projects created using the DBaaS Control Plane",
		GetDataSourceAttributes: GetProjectDataSourceAttributes,
		Build:                   NewProjectDataSourceState,
	}
}
