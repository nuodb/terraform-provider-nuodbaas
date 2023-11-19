package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type DatabaseResourceModel struct {
	Organization    types.String `tfsdk:"organization"`
	Name            types.String `tfsdk:"name"`
	Project         types.String `tfsdk:"project"`
	Password        types.String `tfsdk:"password"`
	Tier            types.String `tfsdk:"tier"`
	ArchiveDiskSize types.String `tfsdk:"archive_disk_size"`
	JournalDiskSize types.String `tfsdk:"journal_disk_size"`
	ResourceVersion types.String `tfsdk:"resource_version"`
	Maintenance     types.Object `tfsdk:"maintenance"`
}

type DatabasePropertiesResourceModel struct {
	ArchiveDiskSize types.String `tfsdk:"archive_disk_size"`
	JournalDiskSize types.String `tfsdk:"journal_disk_size"`
}


type DatabaseCreateUpdateModel struct {
	Password string
	Tier     string
	ArchiveDiskSize string
	JournalDiskSize string
}