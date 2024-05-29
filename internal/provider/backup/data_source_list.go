// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package backup

import (
	"context"
	"fmt"
	"strings"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ framework.DataSourceState = &BackupsDataSourceModel{}
)

type BackupFilterModel struct {
	Organization *string  `tfsdk:"organization" hcl:"organization" cty:"organization"`
	Project      *string  `tfsdk:"project" hcl:"project" cty:"project"`
	Database     *string  `tfsdk:"database" hcl:"database" cty:"database"`
	Labels       []string `tfsdk:"labels" hcl:"labels" cty:"labels"`
}

type BackupNameModel struct {
	Organization string `tfsdk:"organization" hcl:"organization" cty:"organization"`
	Project      string `tfsdk:"project" hcl:"project" cty:"project"`
	Database     string `tfsdk:"database" hcl:"database" cty:"database"`
	Name         string `tfsdk:"name" hcl:"name" cty:"name"`
}

type BackupsDataSourceModel struct {
	Filter  *BackupFilterModel `tfsdk:"filter" hcl:"filter" cty:"filter"`
	Backups []BackupNameModel  `tfsdk:"backups" hcl:"backups" cty:"backups"`
}

// GetBackupsDataSourceSchema returns the schema for the backups (plural) data
// source. This has to be provided explicitly because there is no schema in the
// OpenAPI spec for the REST API that corresponds to it.
func GetBackupsDataSourceSchema() *schema.Schema {
	sb := framework.NewSchemaBuilder().WithDescription("Data source for listing NuoDB backups created using the DBaaS Control Plane")
	sb.WithDatabaseScopeFilters("backups")
	sb.WithDatabaseScopeList("backup", "backups").WithNameAttribute("backup")
	return sb.Build()
}

func (state *BackupsDataSourceModel) Read(ctx context.Context, client openapi.ClientInterface) error {
	var organization, project, database string
	var labelFilter *string
	if state.Filter != nil {
		if state.Filter.Organization != nil {
			organization = *state.Filter.Organization
		}
		if state.Filter.Project != nil {
			project = *state.Filter.Project
		}
		if state.Filter.Database != nil {
			database = *state.Filter.Database
		}
		if state.Filter.Labels != nil {
			labelFilterStr := strings.Join(state.Filter.Labels, ",")
			labelFilter = &labelFilterStr
		}
	}
	backups, err := helper.GetBackups(ctx, client, organization, project, database, labelFilter, true)
	if err != nil {
		return err
	}
	state.Backups, err = GetBackupDataSourceResponse(backups)
	return err
}

func GetBackupDataSourceResponse(backups []string) ([]BackupNameModel, error) {
	var ret []BackupNameModel
	for _, backup := range backups {
		parts := strings.Split(backup, "/")
		if len(parts) != 4 {
			return nil, fmt.Errorf("Unexpected format for backup name: %s", backup)
		}
		ret = append(ret, BackupNameModel{
			Organization: parts[0],
			Project:      parts[1],
			Database:     parts[2],
			Name:         parts[3],
		})
	}
	return ret, nil
}

func NewBackupsDataSourceState() framework.DataSourceState {
	return &BackupsDataSourceModel{}
}

func NewBackupsDataSource() datasource.DataSource {
	return &framework.GenericDataSource{
		TypeName:       "backups",
		SchemaOverride: GetBackupsDataSourceSchema(),
		Build:          NewBackupsDataSourceState,
	}
}
