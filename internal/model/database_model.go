/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
package model

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

type DatabaseResourceModel struct {
	Organization string                    `tfsdk:"organization" json:"organization"`
	Project      string                    `tfsdk:"project" json:"project"`
	Name         string                    `tfsdk:"name" json:"name"`
	DbaPassword  string                    `tfsdk:"dba_password" json:"dbaPassword"`
	Tier         *string                   `tfsdk:"tier" json:"tier"`
	Maintenance  *MaintenanceDatabaseModel `tfsdk:"maintenance" json:"maintenance"`
	Properties   *DatabasePropertiesModel  `tfsdk:"properties" json:"properties"`
	Timeouts     timeouts.Value            `tfsdk:"timeouts"`
}

type MaintenanceDatabaseModel struct {
	IsDisabled *bool `tfsdk:"is_disabled" json:"isDisabled"`
}

type DatabasePropertiesModel struct {
	ArchiveDiskSize *string            `tfsdk:"archive_disk_size" json:"archiveDiskSize"`
	JournalDiskSize *string            `tfsdk:"journal_disk_size" json:"journalDiskSize"`
	TierParameters  *map[string]string `tfsdk:"tier_parameters" json:"tierParameters"`
}

type DatabaseDataSourceNameModel struct {
	Organization string `tfsdk:"organization" json:"organization"`
	Project      string `tfsdk:"project" json:"project"`
	Name         string `tfsdk:"name" json:"name"`
}

type DatabaseDataSourceModel struct {
	Organization    string                    `tfsdk:"organization" json:"organization"`
	Project         string                    `tfsdk:"project" json:"project"`
	Name            string                    `tfsdk:"name" json:"name"`
	Tier            string                    `tfsdk:"tier" json:"tier"`
	Maintenance     *MaintenanceDatabaseModel `tfsdk:"maintenance" json:"maintenance"`
	Properties      *DatabasePropertiesModel  `tfsdk:"properties" json:"properties"`
	ResourceVersion string                    `tfsdk:"resource_version" json:"resourceVersion"`
	Status          *StatusModel              `tfsdk:"status" json:"status"`
}

type StatusModel struct {
	SqlEndpoint *string `tfsdk:"sql_endpoint" json:"sqlEndpoint"`
	CaPem       *string `tfsdk:"ca_pem" json:"caPem"`
	State       *string `tfsdk:"state" json:"state"`
}
