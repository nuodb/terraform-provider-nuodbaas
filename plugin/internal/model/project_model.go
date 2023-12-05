package model

import "github.com/hashicorp/terraform-plugin-framework/types"

type ProjectResourceModel struct {
	Organization    types.String `tfsdk:"organization"`
	Name            types.String `tfsdk:"name"`
	Sla             types.String `tfsdk:"sla"`
	Tier         	types.String `tfsdk:"tier"`
	Maintenance     *MaintenanceModel `tfsdk:"maintenance"`
	ResourceVersion types.String `tfsdk:"resource_version"`
	Properties      *ProjectProperties `tfsdk:"properties"`
}

type ProjectProperties struct {
	TierParameters types.Map  `tfsdk:"tier_parameters"`
}
type MaintenanceModel struct {
	IsDisabled types.Bool   `tfsdk:"is_disabled"`
}


type ProjectDataSourceResponseModel struct {
	Organization types.String `tfsdk:"organization"`
	Name		 types.String `tfsdk:"name"`	
}