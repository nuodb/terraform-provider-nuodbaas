// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package provider_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/backuppolicy"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/stretchr/testify/require"
)

func newBackupPolicy() *BackupPolicyResourceModel {
	// Get organization from user name
	orgName := getOrganization()
	// Generate a random policy name to avoid collisions
	policyName := withRandomSuffix("policy")
	return &BackupPolicyResourceModel{
		Organization: orgName,
		Name:         policyName,
		Frequency:    "@daily",
		Selector: openapi.SelectorModel{
			Scope: orgName,
		},
		Retention: &openapi.RetentionModel{
			Yearly: ptr(int32(3)),
		},
	}
}

func skipIfBackupPoliciesNotSupported(t *testing.T, ctx context.Context, client openapi.ClientInterface) {
	resp, err := client.GetAllBackupPolicies(ctx, &openapi.GetAllBackupPoliciesParams{
		ListAccessible: ptr(true),
	})
	// GET /backuppolicies?listAccessible=true should not return 404 if supported
	if resp.StatusCode == http.StatusNotFound {
		t.Skip("Server does not have /backuppolicies resource")
	}
	// Make sure some other error did not occur
	require.NoError(t, err)
}

func TestBackupPolicy(t *testing.T) {
	// Skip test if /backuppolicies resource is not implemented by REST server
	var providerCfg NuoDbaasProviderModel
	client, err := providerCfg.CreateClient()
	require.NoError(t, err)
	ctx := context.Background()
	skipIfBackupPoliciesNotSupported(t, ctx, client)

	// Create provider server that runs within test
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	err = tf.SetReattachConfig(reattachCfg)
	require.NoError(t, err)

	policy := newBackupPolicy()
	builder := NewTfConfigBuilder().WithProviderConfig("nuodbaas", &providerCfg).
		WithBackupPolicyResource("pol", policy).
		WithBackupPolicyDataSource("pol", &BackupPolicyNameModel{
			Organization: policy.Organization,
			Name:         policy.Name,
		}, "nuodbaas_backuppolicy.pol").
		WithBackupPoliciesDataSource("all_policies", &BackupPoliciesDataSourceModel{}, "nuodbaas_backuppolicy.pol").
		WithBackupPoliciesDataSource("org_policies", &BackupPoliciesDataSourceModel{
			Filter: &BackupPolicyFilterModel{
				Organization: ptr(policy.Organization),
			},
		}, "nuodbaas_backuppolicy.pol").
		WithBackupPoliciesDataSource("otherorg_policies", &BackupPoliciesDataSourceModel{
			Filter: &BackupPolicyFilterModel{
				Organization: ptr("otherorg"),
			},
		}, "nuodbaas_backuppolicy.pol").
		WithBackupPoliciesDataSource("labelled_policies", &BackupPoliciesDataSourceModel{
			Filter: &BackupPolicyFilterModel{
				Organization: ptr(policy.Organization),
				Labels:       []string{"rpo"},
			},
		}, "nuodbaas_backuppolicy.pol")
	tf.WriteConfigT(t, builder.Build())
	_, err = tf.Init()
	require.NoError(t, err)

	// Run `terraform apply` to create backup policy
	_, err = tf.Apply()
	defer tf.DestroySilently()
	require.NoError(t, err)

	// Check attributes in resource
	tf.CheckStateResource(t, "data.nuodbaas_backuppolicy.pol").
		HasAttributeValue("organization", policy.Organization).
		HasAttributeValue("name", policy.Name).
		HasAttributeValue("labels", map[string]any{}).
		HasAttributeValue("frequency", policy.Frequency).
		HasAttributeValue("selector.scope", policy.Selector.Scope).
		HasAttributeValue("selector.slas", []any{}).
		HasAttributeValue("selector.tiers", []any{}).
		HasAttributeValue("selector.labels", map[string]any{}).
		DoesNotHaveAttribute("retention.hourly").
		DoesNotHaveAttribute("retention.daily").
		DoesNotHaveAttribute("retention.weekly").
		DoesNotHaveAttribute("retention.monthly").
		HasAttributeValue("retention.yearly", float64(3))

	// Check attributes in data source
	tf.CheckStateResource(t, "data.nuodbaas_backuppolicy.pol").
		HasAttributeValue("organization", policy.Organization).
		HasAttributeValue("name", policy.Name).
		HasAttributeValue("labels", map[string]any{}).
		HasAttributeValue("frequency", policy.Frequency).
		HasAttributeValue("selector.scope", policy.Selector.Scope).
		HasAttributeValue("selector.slas", []any{}).
		HasAttributeValue("selector.tiers", []any{}).
		HasAttributeValue("selector.labels", map[string]any{}).
		DoesNotHaveAttribute("retention.hourly").
		DoesNotHaveAttribute("retention.daily").
		DoesNotHaveAttribute("retention.weekly").
		DoesNotHaveAttribute("retention.monthly").
		HasAttributeValue("retention.yearly", float64(3))
	tf.CheckStateResource(t, "data.nuodbaas_backuppolicies.all_policies").
		HasAttributeValue("policies", []any{
			map[string]any{
				"organization": policy.Organization,
				"name":         policy.Name,
			},
		})
	tf.CheckStateResource(t, "data.nuodbaas_backuppolicies.org_policies").
		HasAttributeValue("policies", []any{
			map[string]any{
				"organization": policy.Organization,
				"name":         policy.Name,
			},
		})
	tf.CheckStateResource(t, "data.nuodbaas_backuppolicies.otherorg_policies").
		HasAttributeValue("policies", nil)
	tf.CheckStateResource(t, "data.nuodbaas_backuppolicies.labelled_policies").
		HasAttributeValue("policies", nil)

	// Update backup policy resource
	policy.Labels = &map[string]string{
		"rpo":     "1d",
		"purpose": "test",
	}
	policy.Selector.Slas = &[]string{"dev", "qa", "prod"}
	policy.Selector.Tiers = &[]string{"n0.nano", "n0.small", "n1.small"}
	policy.Selector.Labels = &map[string]string{
		"backup-schedule": "daily",
	}
	policy.Retention = &openapi.RetentionModel{
		Hourly:  ptr(int32(24)),
		Daily:   ptr(int32(7)),
		Weekly:  ptr(int32(4)),
		Monthly: ptr(int32(12)),
		Yearly:  ptr(int32(0)),
	}
	tf.WriteConfigT(t, builder.Build())

	// Run `terraform apply` again to update resource
	out, err := tf.Apply()
	require.NoError(t, err)
	require.Contains(t, string(out), "nuodbaas_backuppolicy.pol: Modifying...")

	// Check attributes in resource
	tf.CheckStateResource(t, "data.nuodbaas_backuppolicy.pol").
		HasAttributeValue("organization", policy.Organization).
		HasAttributeValue("name", policy.Name).
		HasAttributeValue("labels", map[string]interface{}{"purpose": "test", "rpo": "1d"}).
		HasAttributeValue("frequency", policy.Frequency).
		HasAttributeValue("selector.scope", policy.Selector.Scope).
		HasAttributeValue("selector.slas", []any{"dev", "qa", "prod"}).
		HasAttributeValue("selector.tiers", []any{"n0.nano", "n0.small", "n1.small"}).
		HasAttributeValue("selector.labels", map[string]any{"backup-schedule": "daily"}).
		HasAttributeValue("retention.hourly", float64(24)).
		HasAttributeValue("retention.daily", float64(7)).
		HasAttributeValue("retention.weekly", float64(4)).
		HasAttributeValue("retention.monthly", float64(12)).
		HasAttributeValue("retention.yearly", float64(0))

	// Check attributes in data source
	tf.CheckStateResource(t, "data.nuodbaas_backuppolicy.pol").
		HasAttributeValue("organization", policy.Organization).
		HasAttributeValue("name", policy.Name).
		HasAttributeValue("labels", map[string]interface{}{"purpose": "test", "rpo": "1d"}).
		HasAttributeValue("frequency", policy.Frequency).
		HasAttributeValue("selector.scope", policy.Selector.Scope).
		HasAttributeValue("selector.slas", []any{"dev", "qa", "prod"}).
		HasAttributeValue("selector.tiers", []any{"n0.nano", "n0.small", "n1.small"}).
		HasAttributeValue("selector.labels", map[string]any{"backup-schedule": "daily"}).
		HasAttributeValue("retention.hourly", float64(24)).
		HasAttributeValue("retention.daily", float64(7)).
		HasAttributeValue("retention.weekly", float64(4)).
		HasAttributeValue("retention.monthly", float64(12)).
		HasAttributeValue("retention.yearly", float64(0))
	tf.CheckStateResource(t, "data.nuodbaas_backuppolicies.all_policies").
		HasAttributeValue("policies", []any{
			map[string]any{
				"organization": policy.Organization,
				"name":         policy.Name,
			},
		})
	tf.CheckStateResource(t, "data.nuodbaas_backuppolicies.org_policies").
		HasAttributeValue("policies", []any{
			map[string]any{
				"organization": policy.Organization,
				"name":         policy.Name,
			},
		})
	tf.CheckStateResource(t, "data.nuodbaas_backuppolicies.otherorg_policies").
		HasAttributeValue("policies", nil)
	// Check that data source with label filter returns policy with label
	tf.CheckStateResource(t, "data.nuodbaas_backuppolicies.labelled_policies").
		HasAttributeValue("policies", []any{
			map[string]any{
				"organization": policy.Organization,
				"name":         policy.Name,
			},
		})

	// Run `terraform apply` again and verify that it does nothing
	out, err = tf.Apply()
	require.NoError(t, err)
	require.Contains(t, string(out), "No changes.")
	require.Contains(t, string(out), "Your infrastructure matches the configuration.")

	// Run `terraform destroy` to delete policy
	_, err = tf.Destroy()
	require.NoError(t, err)

	// Obtain actual policy state and check that 404 is returned
	actualPolicy := *policy
	err = actualPolicy.Read(ctx, client)
	require.Error(t, err)
	require.True(t, helper.IsNotFound(err), "Unexpected error: "+err.Error())
}

func TestImportBackupPolicy(t *testing.T) {
	// Skip test if /backuppolicies resource is not implemented by REST server
	var providerCfg NuoDbaasProviderModel
	client, err := providerCfg.CreateClient()
	require.NoError(t, err)
	ctx := context.Background()
	skipIfBackupPoliciesNotSupported(t, ctx, client)

	// Create backup policy
	policy := newBackupPolicy()
	err = policy.Create(ctx, client)
	require.NoError(t, err)
	defer func() {
		var timeoutSeconds int32 = 10
		_, _ = client.DeleteBackupPolicy(
			ctx, policy.Organization, policy.Name,
			&openapi.DeleteBackupPolicyParams{TimeoutSeconds: &timeoutSeconds})
	}()

	// Create provider server that runs within test
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace
	tf := CreateTerraformWorkspace(t)
	err = tf.SetReattachConfig(reattachCfg)
	require.NoError(t, err)

	// Initialize Terraform workspace with config
	builder := NewTfConfigBuilder().WithProviderConfig("nuodbaas", &providerCfg).
		WithBackupPolicyResource("pol", policy)
	tf.WriteConfigT(t, builder.Build())
	_, err = tf.Init()
	require.NoError(t, err)
	defer tf.DestroySilently()

	// Run `terraform apply` and verify that it fails due to the resources
	// already existing
	out, err := tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Unable to create backuppolicy")
	require.Contains(t, string(out), "409 Conflict")
	require.Contains(t, string(out), fmt.Sprintf("Backup policy %s for organization %s already exists", policy.Name, policy.Organization))

	// Run `terraform refresh` and verify that it does not do anything,
	// since it is only operating on the Terraform state
	out, err = tf.Run("refresh")
	require.NoError(t, err)
	require.Contains(t, string(out), "Empty or non-existent state")

	// Run `terraform import` for backup policy
	out, err = tf.Run("import", "nuodbaas_backuppolicy.pol", policy.Organization+"/"+policy.Name)
	require.NoError(t, err)
	require.Contains(t, string(out), "Import successful!")

	// Verify that backup policy is in state
	tf.CheckStateResource(t, "nuodbaas_backuppolicy.pol").
		HasAttributeValue("organization", policy.Organization).
		HasAttributeValue("name", policy.Name).
		HasAttributeValue("frequency", policy.Frequency).
		HasAttributeValue("selector.scope", policy.Selector.Scope)

	// Run `terraform apply` and verify that there is nothing to do
	out, err = tf.Apply()
	require.NoError(t, err)
	require.Contains(t, string(out), "No changes.")
	require.Contains(t, string(out), "Your infrastructure matches the configuration.")
}
