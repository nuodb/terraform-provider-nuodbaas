// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package backup

import (
	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ framework.DataSourceState = &BackupResourceModel{}
)

func GetBackupDataSourceAttributes() (map[string]schema.Attribute, error) {
	return framework.GetDataSourceAttributes("BackupModel")
}

func NewBackupDataSourceState() framework.DataSourceState {
	return &BackupResourceModel{}
}

func NewBackupDataSource() datasource.DataSource {
	return &framework.GenericDataSource{
		TypeName:                "backup",
		Description:             "Data source for exposing information about NuoDB backups created using the DBaaS Control Plane",
		GetDataSourceAttributes: GetBackupDataSourceAttributes,
		Build:                   NewBackupDataSourceState,
	}
}
