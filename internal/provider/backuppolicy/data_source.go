// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package backuppolicy

import (
	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func GetBackupPolicyDataSourceAttributes() (map[string]schema.Attribute, error) {
	return framework.GetDataSourceAttributes("BackupPolicyModel")
}

func NewBackupPolicyDataSourceState() framework.DataSourceState {
	return &BackupPolicyResourceModel{}
}

func NewBackupPolicyDataSource() datasource.DataSource {
	return &framework.GenericDataSource{
		TypeName:                "backuppolicy",
		Description:             "Data source for exposing information about NuoDB backup policies created using the DBaaS Control Plane",
		GetDataSourceAttributes: GetBackupPolicyDataSourceAttributes,
		Build:                   NewBackupPolicyDataSourceState,
	}
}
