/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
package model

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

type ProjectResourceModel struct {
	Organization string                   `tfsdk:"organization" json:"organization"`
	Name         string                   `tfsdk:"name" json:"name"`
	Sla          string                   `tfsdk:"sla" json:"sla"`
	Tier         string                   `tfsdk:"tier" json:"tier"`
	Maintenance  *MaintenanceProjectModel `tfsdk:"maintenance" json:"maintenance"`
	Properties   *ProjectProperties       `tfsdk:"properties" json:"properties"`
	Timeouts     timeouts.Value           `tfsdk:"timeouts"`
}

type ProjectProperties struct {
	TierParameters *map[string]string `tfsdk:"tier_parameters" json:"tierParameters"`
}
type MaintenanceProjectModel struct {
	IsDisabled *bool `tfsdk:"is_disabled" json:"isDisabled"`
}

type ProjectDataSourceModel struct {
	Organization    string                   `tfsdk:"organization" json:"organization"`
	Name            string                   `tfsdk:"name" json:"name"`
	Sla             string                   `tfsdk:"sla" json:"sla"`
	Tier            string                   `tfsdk:"tier" json:"tier"`
	Maintenance     *MaintenanceProjectModel `tfsdk:"maintenance" json:"maintenance"`
	Properties      *ProjectProperties       `tfsdk:"properties" json:"properties"`
	ResourceVersion string                   `tfsdk:"resource_version" json:"resourceVersion"`
}

type ProjectDataSourceNameModel struct {
	Organization string `tfsdk:"organization" json:"organization"`
	Name         string `tfsdk:"name" json:"name"`
}
