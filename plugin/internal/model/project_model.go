package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type ProjectResourceModel struct {
	Organization    types.String `tfsdk:"organization"`
	Name            types.String `tfsdk:"name"`
	Sla             types.String `tfsdk:"sla"`
	Tier         types.String `tfsdk:"tier"`
	Maintenance     types.Object `tfsdk:"maintenance"`
	ResourceVersion types.String `tfsdk:"resource_version"`
}

type MaintenanceModel struct {
	IsDisabled types.Bool   `tfsdk:"is_disabled"`
}