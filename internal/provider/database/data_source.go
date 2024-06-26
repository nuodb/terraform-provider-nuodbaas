// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package database

import (
	"context"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ framework.DataSourceState = &DatabaseResourceModel{}
)

type DatabaseDataSourceModel openapi.DatabaseModel

func (state *DatabaseDataSourceModel) Read(ctx context.Context, client openapi.ClientInterface) error {
	resp, err := client.GetDatabase(ctx, state.Organization, state.Project, state.Name)
	if err != nil {
		return err
	}
	return helper.ParseResponse(resp, state)
}

func GetDatabaseDataSourceAttributes() (map[string]schema.Attribute, error) {
	return framework.GetDataSourceAttributes("DatabaseModel")
}

func NewDatabaseDataSourceState() framework.DataSourceState {
	return &DatabaseDataSourceModel{}
}

func NewDatabaseDataSource() datasource.DataSource {
	return &framework.GenericDataSource{
		TypeName:                "database",
		Description:             "Data source for exposing information about NuoDB databases created using the DBaaS Control Plane",
		GetDataSourceAttributes: GetDatabaseDataSourceAttributes,
		Build:                   NewDatabaseDataSourceState,
	}
}
