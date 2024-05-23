// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package backuppolicy

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
	_ framework.DataSourceState = &BackupPoliciesDataSourceModel{}
)

type BackupPolicyFilterModel struct {
	Organization *string  `tfsdk:"organization" hcl:"organization" cty:"organization"`
	Labels       []string `tfsdk:"labels" hcl:"labels" cty:"labels"`
}

type BackupPolicyNameModel struct {
	Organization string `tfsdk:"organization" hcl:"organization" cty:"organization"`
	Name         string `tfsdk:"name" hcl:"name" cty:"name"`
}

type BackupPoliciesDataSourceModel struct {
	Filter   *BackupPolicyFilterModel `tfsdk:"filter" hcl:"filter" cty:"filter"`
	Policies []BackupPolicyNameModel  `tfsdk:"policies" hcl:"policies" cty:"policies"`
}

// GetBackupPoliciesDataSourceSchema returns the schema for the backuppolicies
// (plural) data source. This has to be provided explicitly because there is no
// schema in the OpenAPI spec for the REST API that corresponds to it.
func GetBackupPoliciesDataSourceSchema() *schema.Schema {
	sb := framework.NewSchemaBuilder().WithDescription("Data source for listing NuoDB backup policies created using the DBaaS Control Plane")
	sb.WithOrganizationScopeFilters("policies")
	sb.WithOrganizationScopeList("policy", "policies").WithNameAttribute("policy")
	return sb.Build()
}

// Read implements datasource.DataSource.
func (state *BackupPoliciesDataSourceModel) Read(ctx context.Context, client openapi.ClientInterface) error {
	var organization string
	var labelFilter *string
	if state.Filter != nil {
		if state.Filter.Organization != nil {
			organization = *state.Filter.Organization
		}
		if state.Filter.Labels != nil {
			labelFilterStr := strings.Join(state.Filter.Labels, ",")
			labelFilter = &labelFilterStr
		}
	}
	policies, err := helper.GetBackupPolicies(ctx, client, organization, labelFilter, true)
	if err != nil {
		return err
	}
	state.Policies, err = GetBackupPoliciesDataSourceResponse(policies)
	return err
}

func GetBackupPoliciesDataSourceResponse(policies []string) ([]BackupPolicyNameModel, error) {
	var ret []BackupPolicyNameModel
	for _, policy := range policies {
		parts := strings.Split(policy, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Unexpected format for backup policy name: %s", policy)
		}
		ret = append(ret, BackupPolicyNameModel{
			Organization: parts[0],
			Name:         parts[1],
		})
	}
	return ret, nil
}

func NewBackupPoliciesDataSourceState() framework.DataSourceState {
	return &BackupPoliciesDataSourceModel{}
}

func NewBackupPoliciesDataSource() datasource.DataSource {
	return &framework.GenericDataSource{
		TypeName:       "backuppolicies",
		SchemaOverride: GetBackupPoliciesDataSourceSchema(),
		Build:          NewBackupPoliciesDataSourceState,
	}
}
