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

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/backup"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/backuppolicy"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/stretchr/testify/require"
)

func newBackup(dbName, backupName string) *BackupResourceModel {
	// Get organization from user name
	orgName := getOrganization()
	// Generate a random project name to avoid collisions
	projectName := withRandomSuffix("proj")
	return &BackupResourceModel{
		Organization: orgName,
		Project:      projectName,
		Database:     dbName,
		Name:         backupName,
	}
}

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
	}
}

func skipIfNotSupported(t *testing.T, requestFn func() (*http.Response, error)) {
	resp, err := requestFn()
	// function should not return 404 if request is supported
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		t.Skipf("Server does not support '%s %s'", resp.Request.Method, resp.Request.URL.Path)
	}
	// Make sure some other error did not occur
	require.NoError(t, err)
}

func skipIfBackupsNotSupported(t *testing.T, ctx context.Context, client openapi.ClientInterface) {
	skipIfNotSupported(t, func() (*http.Response, error) {
		return client.GetAllBackups(ctx, &openapi.GetAllBackupsParams{
			ListAccessible: ptr(true),
		})
	})
}

func skipIfBackupPoliciesNotSupported(t *testing.T, ctx context.Context, client openapi.ClientInterface) {
	skipIfNotSupported(t, func() (*http.Response, error) {
		return client.GetAllBackupPolicies(ctx, &openapi.GetAllBackupPoliciesParams{
			ListAccessible: ptr(true),
		})
	})
}

func TestBackup(t *testing.T) {
	// Skip test if /backups resource is not implemented by REST server
	var providerCfg NuoDbaasProviderModel
	client, err := providerCfg.CreateClient()
	require.NoError(t, err)
	ctx := context.Background()
	skipIfBackupsNotSupported(t, ctx, client)

	// Create provider server that runs within test
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	err = tf.SetReattachConfig(reattachCfg)
	require.NoError(t, err)

	// Disable readiness check for backup resource so that backup controller does not have to be enabled
	providerCfg.Timeouts = map[string]framework.OperationTimeouts{
		"backup": {
			Create: ptr("0"),
			Update: ptr("0"),
		},
	}
	backup := newBackup("db", "backup")
	backup.ImportSource = &openapi.ImportSourceModel{
		BackupHandle: "backup",
		BackupPlugin: "plugin",
	}
	builder := NewTfConfigBuilder().WithProviderConfig("nuodbaas", &providerCfg).
		WithBackupResource("backup", backup).
		WithBackupDataSource("backup", &BackupNameModel{
			Organization: backup.Organization,
			Project:      backup.Project,
			Database:     backup.Database,
			Name:         backup.Name,
		}, "nuodbaas_backup.backup").
		WithBackupsDataSource("all_backups", &BackupsDataSourceModel{}, "nuodbaas_backup.backup").
		WithBackupsDataSource("org_backups", &BackupsDataSourceModel{
			Filter: &BackupFilterModel{
				Organization: ptr(backup.Organization),
			},
		}, "nuodbaas_backup.backup").
		WithBackupsDataSource("proj_backups", &BackupsDataSourceModel{
			Filter: &BackupFilterModel{
				Organization: ptr(backup.Organization),
				Project:      ptr(backup.Project),
			},
		}, "nuodbaas_backup.backup").
		WithBackupsDataSource("db_backups", &BackupsDataSourceModel{
			Filter: &BackupFilterModel{
				Organization: ptr(backup.Organization),
				Project:      ptr(backup.Project),
				Database:     ptr(backup.Database),
			},
		}, "nuodbaas_backup.backup").
		WithBackupsDataSource("otherorg_backups", &BackupsDataSourceModel{
			Filter: &BackupFilterModel{
				Organization: ptr("otherorg"),
			},
		}, "nuodbaas_backup.backup").
		WithBackupsDataSource("otherproj_backups", &BackupsDataSourceModel{
			Filter: &BackupFilterModel{
				Organization: ptr(backup.Organization),
				Project:      ptr("otherproj"),
			},
		}, "nuodbaas_backup.backup").
		WithBackupsDataSource("otherdb_backups", &BackupsDataSourceModel{
			Filter: &BackupFilterModel{
				Organization: ptr(backup.Organization),
				Project:      ptr(backup.Project),
				Database:     ptr("otherdb"),
			},
		}, "nuodbaas_backup.backup").
		WithBackupsDataSource("labelled_backups", &BackupsDataSourceModel{
			Filter: &BackupFilterModel{
				Organization: ptr(backup.Organization),
				Labels:       []string{"purpose=test"},
			},
		}, "nuodbaas_backup.backup")
	tf.WriteConfigT(t, builder.Build())
	_, err = tf.Init()
	require.NoError(t, err)

	// Run `terraform apply` to create backup
	_, err = tf.Apply()
	defer tf.DestroySilently()
	require.NoError(t, err)

	// Check attributes in resource
	tf.CheckStateResource(t, "data.nuodbaas_backup.backup").
		HasAttributeValue("organization", backup.Organization).
		HasAttributeValue("project", backup.Project).
		HasAttributeValue("database", backup.Database).
		HasAttributeValue("name", backup.Name).
		HasAttributeValue("labels", map[string]any{}).
		HasAttributeValue("import_source.backup_handle", "backup").
		HasAttributeValue("import_source.backup_plugin", "plugin").
		HasAttribute("status.ready_to_use").
		HasAttribute("status.state")

	// Check attributes in data source
	tf.CheckStateResource(t, "data.nuodbaas_backup.backup").
		HasAttributeValue("organization", backup.Organization).
		HasAttributeValue("project", backup.Project).
		HasAttributeValue("database", backup.Database).
		HasAttributeValue("name", backup.Name).
		HasAttributeValue("labels", map[string]any{}).
		HasAttributeValue("import_source.backup_handle", "backup").
		HasAttributeValue("import_source.backup_plugin", "plugin").
		HasAttribute("status.ready_to_use").
		HasAttribute("status.state")
	tf.CheckStateResource(t, "data.nuodbaas_backups.all_backups").
		HasAttributeValue("backups", []any{
			map[string]any{
				"organization": backup.Organization,
				"project":      backup.Project,
				"database":     backup.Database,
				"name":         backup.Name,
			},
		})
	tf.CheckStateResource(t, "data.nuodbaas_backups.org_backups").
		HasAttributeValue("backups", []any{
			map[string]any{
				"organization": backup.Organization,
				"project":      backup.Project,
				"database":     backup.Database,
				"name":         backup.Name,
			},
		})
	tf.CheckStateResource(t, "data.nuodbaas_backups.proj_backups").
		HasAttributeValue("backups", []any{
			map[string]any{
				"organization": backup.Organization,
				"project":      backup.Project,
				"database":     backup.Database,
				"name":         backup.Name,
			},
		})
	tf.CheckStateResource(t, "data.nuodbaas_backups.db_backups").
		HasAttributeValue("backups", []any{
			map[string]any{
				"organization": backup.Organization,
				"project":      backup.Project,
				"database":     backup.Database,
				"name":         backup.Name,
			},
		})
	tf.CheckStateResource(t, "data.nuodbaas_backups.otherorg_backups").
		HasAttributeValue("backups", nil)
	tf.CheckStateResource(t, "data.nuodbaas_backups.otherproj_backups").
		HasAttributeValue("backups", nil)
	tf.CheckStateResource(t, "data.nuodbaas_backups.otherdb_backups").
		HasAttributeValue("backups", nil)
	tf.CheckStateResource(t, "data.nuodbaas_backups.labelled_backups").
		HasAttributeValue("backups", nil)

	// Update backup labels
	backup.Labels = &map[string]string{
		"purpose": "test",
	}
	tf.WriteConfigT(t, builder.Build())

	// Run `terraform apply` again to update resource
	out, err := tf.Apply()
	require.NoError(t, err)
	require.Contains(t, string(out), "nuodbaas_backup.backup: Modifying...")

	// Check attributes in resource
	tf.CheckStateResource(t, "nuodbaas_backup.backup").
		HasAttributeValue("organization", backup.Organization).
		HasAttributeValue("project", backup.Project).
		HasAttributeValue("database", backup.Database).
		HasAttributeValue("name", backup.Name).
		HasAttributeValue("labels", map[string]any{"purpose": "test"}).
		HasAttributeValue("import_source.backup_handle", "backup").
		HasAttributeValue("import_source.backup_plugin", "plugin").
		HasAttribute("status.ready_to_use").
		HasAttribute("status.state")

	// Check attributes in data source
	tf.CheckStateResource(t, "data.nuodbaas_backup.backup").
		HasAttributeValue("organization", backup.Organization).
		HasAttributeValue("project", backup.Project).
		HasAttributeValue("database", backup.Database).
		HasAttributeValue("name", backup.Name).
		HasAttributeValue("labels", map[string]any{"purpose": "test"}).
		HasAttributeValue("import_source.backup_handle", "backup").
		HasAttributeValue("import_source.backup_plugin", "plugin").
		HasAttribute("status.ready_to_use").
		HasAttribute("status.state")
	tf.CheckStateResource(t, "data.nuodbaas_backups.all_backups").
		HasAttributeValue("backups", []any{
			map[string]any{
				"organization": backup.Organization,
				"project":      backup.Project,
				"database":     backup.Database,
				"name":         backup.Name,
			},
		})
	tf.CheckStateResource(t, "data.nuodbaas_backups.org_backups").
		HasAttributeValue("backups", []any{
			map[string]any{
				"organization": backup.Organization,
				"project":      backup.Project,
				"database":     backup.Database,
				"name":         backup.Name,
			},
		})
	tf.CheckStateResource(t, "data.nuodbaas_backups.proj_backups").
		HasAttributeValue("backups", []any{
			map[string]any{
				"organization": backup.Organization,
				"project":      backup.Project,
				"database":     backup.Database,
				"name":         backup.Name,
			},
		})
	tf.CheckStateResource(t, "data.nuodbaas_backups.db_backups").
		HasAttributeValue("backups", []any{
			map[string]any{
				"organization": backup.Organization,
				"project":      backup.Project,
				"database":     backup.Database,
				"name":         backup.Name,
			},
		})
	tf.CheckStateResource(t, "data.nuodbaas_backups.otherorg_backups").
		HasAttributeValue("backups", nil)
	tf.CheckStateResource(t, "data.nuodbaas_backups.otherproj_backups").
		HasAttributeValue("backups", nil)
	tf.CheckStateResource(t, "data.nuodbaas_backups.otherdb_backups").
		HasAttributeValue("backups", nil)
	tf.CheckStateResource(t, "data.nuodbaas_backups.labelled_backups").
		HasAttributeValue("backups", []any{
			map[string]any{
				"organization": backup.Organization,
				"project":      backup.Project,
				"database":     backup.Database,
				"name":         backup.Name,
			},
		})

	// Run `terraform apply` again and verify that it does nothing
	out, err = tf.Apply()
	require.NoError(t, err)
	require.Contains(t, string(out), "No changes.")
	require.Contains(t, string(out), "Your infrastructure matches the configuration.")

	t.Run("immutableAttributeChange", func(t *testing.T) {
		original := *backup
		defer func() {
			*backup = original
			tf.WriteConfigT(t, builder.Build())
		}()

		// Try to change an immutable attribute `name`
		backup.Name = "notbackup"
		tf.WriteConfigT(t, builder.Build())

		// Run `terraform apply` and check that "Immutable Attribute Change" warning is shown
		out, err := tf.Apply()
		require.Error(t, err)
		require.Contains(t, string(out), "Warning: Immutable Attribute Change")
		require.Contains(t, string(out), "A change has been made to the immutable attribute `name`, which may be rejected by the server or ignored. "+
			"In order for an immutable attribute change to take effect, it is necessary to delete and re-create the resource, which may result in data loss.")
		require.Contains(t, string(out), "Error: Unable to update backup")
		require.Contains(t, string(out), fmt.Sprintf(
			"Error response from Control Plane service: status='HTTP 404 Not Found', code=HTTP_ERROR, detail=[No backup %s for database %s/%s/%s]",
			backup.Name, backup.Organization, backup.Project, backup.Database))
		backup.Name = original.Name

		// Try to change immutable attribute `import_source.backup_handle`
		backup.ImportSource = &openapi.ImportSourceModel{
			BackupHandle: "notbackup",
			BackupPlugin: original.ImportSource.BackupPlugin,
		}
		tf.WriteConfigT(t, builder.Build())

		// Run `terraform apply` and check that "Immutable Attribute Change" warning is shown
		out, err = tf.Apply()
		require.Error(t, err)
		require.Contains(t, string(out), "Warning: Immutable Attribute Change")
		require.Contains(t, string(out), "A change has been made to the immutable attribute `import_source.backup_handle`, which may be rejected by the server or ignored. "+
			"In order for an immutable attribute change to take effect, it is necessary to delete and re-create the resource, which may result in data loss.")
		require.Contains(t, string(out), "Error response from Control Plane service: status='HTTP 409 Conflict', code=HTTP_ERROR, detail=[Backup import source cannot be updated: ")
	})

	t.Run("dataSourceFilterNegative", func(t *testing.T) {
		var backupsNeg BackupsDataSourceModel
		builder.WithBackupsDataSource("neg", &backupsNeg)
		defer func() {
			builder.WithoutBackupsDataSource("neg")
			tf.WriteConfigT(t, builder.Build())
		}()

		// Specify invalid filter: project but no organization
		backupsNeg.Filter = &BackupFilterModel{
			Project: ptr("proj"),
		}
		tf.WriteConfigT(t, builder.Build())
		// Run `terraform apply` and check that filter is rejected
		out, err := tf.Apply()
		require.Error(t, err)
		require.Contains(t, string(out), "Error: Unable to read backups")
		require.Contains(t, string(out), "Cannot specify project filter (proj) without organization")

		// Specify invalid filter: database but no project
		backupsNeg.Filter = &BackupFilterModel{
			Organization: ptr("org"),
			Database:     ptr("db"),
		}
		tf.WriteConfigT(t, builder.Build())
		// Run `terraform apply` and check that filter is rejected
		out, err = tf.Apply()
		require.Error(t, err)
		require.Contains(t, string(out), "Error: Unable to read backups")
		require.Contains(t, string(out), "Cannot specify database filter (db) without project")

		// Specify invalid filter: database but no organization/project
		backupsNeg.Filter = &BackupFilterModel{
			Database: ptr("db"),
		}
		tf.WriteConfigT(t, builder.Build())
		// Run `terraform apply` and check that filter is rejected
		out, err = tf.Apply()
		require.Error(t, err)
		require.Contains(t, string(out), "Error: Unable to read backups")
		require.Contains(t, string(out), "Cannot specify database filter (db) without project")
	})

	// Run `terraform destroy` to delete policy
	_, err = tf.Destroy()
	require.NoError(t, err)

	// Obtain actual policy state and check that 404 is returned
	actualBackup := *backup
	err = actualBackup.Read(ctx, client)
	require.Error(t, err)
	require.True(t, helper.IsNotFound(err), "Unexpected error: "+err.Error())
}

func TestImportBackup(t *testing.T) {
	// Skip test if /backup resource is not implemented by REST server
	var providerCfg NuoDbaasProviderModel
	client, err := providerCfg.CreateClient()
	require.NoError(t, err)
	ctx := context.Background()
	skipIfBackupsNotSupported(t, ctx, client)

	// Disable readiness check for backup resource so that backup controller does not have to be enabled
	providerCfg.Timeouts = map[string]framework.OperationTimeouts{
		"backup": {
			Create: ptr("0"),
			Update: ptr("0"),
		},
	}
	// Create backup resource
	backup := newBackup("db", "backup")
	backup.ImportSource = &openapi.ImportSourceModel{
		BackupHandle: "backup",
		BackupPlugin: "plugin",
	}
	err = backup.Create(ctx, client)
	require.NoError(t, err)
	defer func() {
		var timeoutSeconds int32 = 10
		_, _ = client.DeleteBackup(
			ctx, backup.Organization, backup.Project, backup.Database, backup.Name,
			&openapi.DeleteBackupParams{TimeoutSeconds: &timeoutSeconds})
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
		WithBackupResource("backup", backup)
	tf.WriteConfigT(t, builder.Build())
	_, err = tf.Init()
	require.NoError(t, err)
	defer tf.DestroySilently()

	// Run `terraform apply` and verify that it fails due to the resources
	// already existing
	out, err := tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Unable to create backup")
	require.Contains(t, string(out), "409 Conflict")
	require.Contains(t, string(out), fmt.Sprintf(
		"Backup %s for database %s/%s/%s already exists",
		backup.Name, backup.Organization, backup.Project, backup.Database))

	// Run `terraform refresh` and verify that it does not do anything,
	// since it is only operating on the Terraform state
	out, err = tf.Run("refresh")
	require.NoError(t, err)
	require.Contains(t, string(out), "Empty or non-existent state")

	// Run `terraform import` for backup
	out, err = tf.Run("import", "nuodbaas_backup.backup",
		fmt.Sprintf("%s/%s/%s/%s", backup.Organization, backup.Project, backup.Database, backup.Name))
	require.NoError(t, err)
	require.Contains(t, string(out), "Import successful!")

	// Verify that backup is in state
	tf.CheckStateResource(t, "nuodbaas_backup.backup").
		HasAttributeValue("organization", backup.Organization).
		HasAttributeValue("project", backup.Project).
		HasAttributeValue("database", backup.Database).
		HasAttributeValue("name", backup.Name).
		HasAttributeValue("import_source.backup_handle", backup.ImportSource.BackupHandle).
		HasAttributeValue("import_source.backup_plugin", backup.ImportSource.BackupPlugin)

	// Run `terraform apply` and verify that there is nothing to do
	out, err = tf.Apply()
	require.NoError(t, err)
	require.Contains(t, string(out), "No changes.")
	require.Contains(t, string(out), "Your infrastructure matches the configuration.")
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
	policy.Retention = &openapi.RetentionModel{
		Yearly: ptr(int32(3)),
	}
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
	tf.CheckStateResource(t, "nuodbaas_backuppolicy.pol").
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
	tf.CheckStateResource(t, "nuodbaas_backuppolicy.pol").
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
