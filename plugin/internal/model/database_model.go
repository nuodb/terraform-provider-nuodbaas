/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
package model

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DatabaseResourceModel struct {
	Organization    types.String 						`tfsdk:"organization"`
	Name            types.String 						`tfsdk:"name"`
	Project         types.String 						`tfsdk:"project"`
	Password        types.String 						`tfsdk:"dba_password"`
	Tier            types.String 						`tfsdk:"tier"`
	Properties      types.Object 						`tfsdk:"properties"`
	ResourceVersion types.String 						`tfsdk:"resource_version"`
	Maintenance     *MaintenanceModel 					`tfsdk:"maintenance"`
	Timeouts		timeouts.Value 						`tfsdk:"timeouts"`
	Status          types.Object						`tfsdk:"status"`
}

type DatabasePropertiesResourceModel struct {
	ArchiveDiskSize types.String `tfsdk:"archive_disk_size"`
	JournalDiskSize types.String `tfsdk:"journal_disk_size"`
	TierParameters  types.Map    `tfsdk:"tier_parameters"`
}

type MaintenanceDataSourceModel struct {
	IsDisabled types.Bool   `tfsdk:"is_disabled"`
	ExpiresIn types.String  `tfsdk:"expires_in"`
	ExpiresAt types.String  `tfsdk:"expires_at"` 
}

type DatabaseCreateUpdateModel struct {
	Password        string
	Tier            string
	ArchiveDiskSize string
	JournalDiskSize string
}


type DatabaseDataSourceModel struct {
	Organization    types.String 						`tfsdk:"organization"`
	Name            types.String 						`tfsdk:"name"`
	Project         types.String 						`tfsdk:"project"`
	Tier            types.String 						`tfsdk:"tier"`
	Properties      *DatabasePropertiesResourceModel	`tfsdk:"properties"`
	ResourceVersion types.String 						`tfsdk:"resource_version"`
	Maintenance     *MaintenanceDataSourceModel 		`tfsdk:"maintenance"`
	Status			*StatusModel                        `tfsdk:"status"`	
}

type DatabasesDataSourceResponseModel struct {
	Organization types.String `tfsdk:"organization"`
	Project      types.String `tfsdk:"project"`
	Name		 types.String `tfsdk:"name"`	
}

type StatusModel struct {
	SqlEndPoint 	types.String 	`tfsdk:"sql_end_point"`
	CaPem			types.String    `tfsdk:"ca_pem"`
}