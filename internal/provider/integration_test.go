package provider

import (
	"context"
	"os"
	"testing"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/database"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/project"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/stretchr/testify/require"
)

type testVars struct {
	providerCfg NuoDbaasProviderModel
	project     ProjectResourceModel
	database    DatabaseResourceModel
	builder     *TfConfigBuilder
}

func (vars *testVars) resetVars() {
	dbaPassword := "dba"
	vars.providerCfg = NuoDbaasProviderModel{}
	vars.project = ProjectResourceModel{
		Organization: "org",
		Name:         "proj",
		Sla:          "dev",
		Tier:         "n0.nano",
	}
	vars.database = DatabaseResourceModel{
		Organization: "org",
		Project:      "proj",
		Name:         "db",
		DbaPassword:  &dbaPassword,
	}
	vars.builder = NewTfConfigBuilder().WithProviderConfig("nuodbaas", &vars.providerCfg).
		WithProjectResource("proj", &vars.project).
		WithDatabaseResource("db", &vars.database, "nuodbaas_project.proj").
		WithProjectDataSource("proj", &ProjectNameModel{
			Organization: "org",
			Name:         "proj",
		}, "nuodbaas_project.proj").
		WithDatabaseDataSource("db", &DatabaseNameModel{
			Organization: "org",
			Project:      "proj",
			Name:         "db",
		}, "nuodbaas_database.db").
		WithProjectsDataSource("proj_list", &ProjectsDataSourceModel{}).
		WithDatabasesDataSource("db_list", &DatabasesDataSourceModel{})
}

func newTestVars() *testVars {
	var ret testVars
	ret.resetVars()
	return &ret
}

func TestFullLifecycle(t *testing.T) {
	vars := newTestVars()

	// Create provider server that runs within test
	ctx := context.Background()
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	tf.SetReattachConfig(reattachCfg)
	tf.WriteConfigT(t, vars.builder.Build())
	_, err := tf.Init()
	require.NoError(t, err)

	// Run `terraform apply` to create project and database
	_, err = tf.Apply()
	defer tf.DestroySilently()
	require.NoError(t, err)

	// Use client created from provider config to populate structs with
	// current state
	client, err := vars.providerCfg.CreateClient()
	actualProject := vars.project
	err = actualProject.Read(ctx, client)
	require.NoError(t, err)
	require.Equal(t, "org", actualProject.Organization)
	require.Equal(t, "proj", actualProject.Name)
	require.Equal(t, "dev", actualProject.Sla)
	require.Equal(t, "n0.nano", actualProject.Tier)
	require.NotNil(t, actualProject.Status)
	require.NotNil(t, actualProject.Status.CaPem)
	require.NotNil(t, actualProject.Status.State)
	require.Equal(t, openapi.ProjectStatusModelStateAvailable, *actualProject.Status.State)

	actualDatabase := vars.database
	err = actualDatabase.Read(ctx, client)
	require.NoError(t, err)
	require.Equal(t, "org", actualDatabase.Organization)
	require.Equal(t, "proj", actualDatabase.Project)
	require.Equal(t, "db", actualDatabase.Name)
	require.NotNil(t, actualDatabase.Tier)
	require.Equal(t, "n0.nano", *actualDatabase.Tier)
	require.NotNil(t, actualDatabase.Status)
	require.NotNil(t, actualDatabase.Status.CaPem)
	require.NotNil(t, actualDatabase.Status.State)
	require.Equal(t, openapi.DatabaseStatusModelStateAvailable, *actualDatabase.Status.State)

	// Run `terraform apply` again and verify that it does nothing
	out, err := tf.Apply()
	require.NoError(t, err)
	require.Contains(t, string(out), "No changes.")
	require.Contains(t, string(out), "Your infrastructure matches the configuration.")

	// Check attributes in data sources
	tf.CheckStateResource(t, "data.nuodbaas_project.proj").
		HasAttributeValue("organization", "org").
		HasAttributeValue("name", "proj").
		HasAttributeValue("labels", map[string]any{}).
		HasAttribute("properties.product_version").
		HasAttributeValue("status.state", string(openapi.ProjectStatusModelStateAvailable)).
		HasAttributeValue("status.ready", true).
		HasAttributeValue("status.shutdown", false)
	tf.CheckStateResource(t, "data.nuodbaas_database.db").
		HasAttributeValue("organization", "org").
		HasAttributeValue("project", "proj").
		HasAttributeValue("name", "db").
		HasAttributeValue("labels", map[string]any{}).
		HasAttributeValue("tier", "n0.nano").
		HasAttribute("properties.product_version").
		HasAttribute("status.sql_endpoint").
		HasAttributeValue("status.state", string(openapi.DatabaseStatusModelStateAvailable)).
		HasAttributeValue("status.ready", true).
		HasAttributeValue("status.shutdown", false)

	// Refresh state and check list data sources. These do not have
	// dependencies, so they may not be populated after the initial
	// `terraform apply`.
	_, err = tf.Run("refresh")
	require.NoError(t, err)
	tf.CheckStateResource(t, "data.nuodbaas_projects.proj_list").
		HasAttributeValue("projects", []any{
			map[string]any{
				"organization": "org",
				"name":         "proj",
			}})
	tf.CheckStateResource(t, "data.nuodbaas_databases.db_list").
		HasAttributeValue("databases", []any{
			map[string]any{
				"organization": "org",
				"project":      "proj",
				"name":         "db",
			}})

	// Update database config attributes (tier, labels, and product_version)
	// and execute `terraform apply`
	tier := "n1.small"
	if val, ok := os.LookupEnv("E2E_TEST"); ok && val == "true" {
		// Avoid setting an unavailable service tier or increasing the
		// replica count in end-to-end tests, which could be time
		// consuming
		tier = "n0.nano"
	}
	vars.database.Tier = &tier
	vars.database.Labels = &map[string]string{
		"priority": "high",
	}
	productVersion := "6.0"
	vars.database.Properties = &openapi.DatabasePropertiesModel{
		ProductVersion: &productVersion,
	}
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)
	err = actualDatabase.Read(ctx, client)
	require.NotNil(t, actualDatabase.Tier)
	require.Equal(t, tier, *actualDatabase.Tier)
	require.Equal(t, vars.database.Labels, actualDatabase.Labels)
	require.NotNil(t, actualDatabase.Properties)
	require.Equal(t, vars.database.Properties.ProductVersion, actualDatabase.Properties.ProductVersion)

	// Update project config attributes (tier, labels, and product_version)
	// and execute `terraform apply`
	vars.project.Tier = tier
	vars.project.Labels = &map[string]string{
		"priority": "high",
	}
	vars.project.Properties = &openapi.ProjectPropertiesModel{
		ProductVersion: &productVersion,
	}
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)
	err = actualProject.Read(ctx, client)
	require.NoError(t, err)
	require.Equal(t, tier, actualProject.Tier)
	require.Equal(t, vars.project.Labels, actualProject.Labels)
	require.NotNil(t, actualProject.Properties)
	require.Equal(t, vars.project.Properties.ProductVersion, actualProject.Properties.ProductVersion)

	// Check attributes in data sources
	tf.CheckStateResource(t, "data.nuodbaas_project.proj").
		HasAttributeValue("organization", "org").
		HasAttributeValue("name", "proj").
		HasAttributeValue("tier", tier).
		HasAttributeValue("labels", map[string]any{"priority": "high"}).
		HasAttributeValue("properties.product_version", productVersion).
		HasAttributeValue("status.state", string(openapi.ProjectStatusModelStateAvailable)).
		HasAttributeValue("status.ready", true).
		HasAttributeValue("status.shutdown", false)
	tf.CheckStateResource(t, "data.nuodbaas_database.db").
		HasAttributeValue("organization", "org").
		HasAttributeValue("project", "proj").
		HasAttributeValue("name", "db").
		HasAttributeValue("tier", tier).
		HasAttributeValue("labels", map[string]any{"priority": "high"}).
		HasAttributeValue("properties.product_version", productVersion).
		HasAttributeValue("status.state", string(openapi.DatabaseStatusModelStateAvailable)).
		HasAttributeValue("status.ready", true).
		HasAttributeValue("status.shutdown", false)

	// Set project and database as disabled and also update labels
	disabled := true
	vars.project.Maintenance = &openapi.MaintenanceModel{
		IsDisabled: &disabled,
	}
	// Omit labels from project, which should not cause them to be removed
	// because unknown values are resolved from state
	vars.project.Labels = nil
	vars.database.Maintenance = &openapi.MaintenanceModel{
		IsDisabled: &disabled,
	}
	// Explicitly set labels to empty for database, which should cause
	// labels to be removed
	vars.database.Labels = &map[string]string{}
	// Write config and apply
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)

	// Re-fetch state from REST service and check values
	err = actualProject.Read(ctx, client)
	require.NoError(t, err)
	err = actualDatabase.Read(ctx, client)
	require.NoError(t, err)
	require.NotNil(t, actualProject.Maintenance)
	require.NotNil(t, actualProject.Maintenance.IsDisabled)
	require.True(t, *actualProject.Maintenance.IsDisabled)
	require.NotNil(t, actualDatabase.Maintenance)
	require.NotNil(t, actualDatabase.Maintenance.IsDisabled)
	require.True(t, *actualDatabase.Maintenance.IsDisabled)

	// Check attributes in data sources
	tf.CheckStateResource(t, "data.nuodbaas_project.proj").
		HasAttributeValue("organization", "org").
		HasAttributeValue("name", "proj").
		HasAttributeValue("tier", tier).
		HasAttributeValue("labels", map[string]any{"priority": "high"}).
		HasAttributeValue("properties.product_version", productVersion).
		HasAttributeValue("maintenance.is_disabled", true).
		HasAttributeValue("status.state", string(openapi.ProjectStatusModelStateStopped)).
		HasAttributeValue("status.shutdown", true)
	tf.CheckStateResource(t, "data.nuodbaas_database.db").
		HasAttributeValue("organization", "org").
		HasAttributeValue("project", "proj").
		HasAttributeValue("name", "db").
		HasAttributeValue("tier", tier).
		HasAttributeValue("labels", map[string]any{}).
		HasAttributeValue("properties.product_version", productVersion).
		HasAttributeValue("maintenance.is_disabled", true).
		HasAttributeValue("status.state", string(openapi.DatabaseStatusModelStateStopped)).
		HasAttributeValue("status.shutdown", true)

	// Run `terraform destroy` to delete database and project
	_, err = tf.Destroy()
	require.NoError(t, err)

	// Obtain actual project and database state and check that 404 is returned
	err = actualProject.Read(ctx, client)
	require.Error(t, err)
	require.True(t, helper.IsNotFound(err), "Unexpected error: "+err.Error())
	err = actualDatabase.Read(ctx, client)
	require.Error(t, err)
	require.True(t, helper.IsNotFound(err), "Unexpected error: "+err.Error())
}

func TestTimeouts(t *testing.T) {
	vars := newTestVars()

	// Disable reconciliation
	reset := SetMockReconcilePolicy(t, MockReconcilePolicy{MarkAsReady: "false"})
	defer reset()

	// Specify timeout for all resources
	timeout := "1s"
	vars.providerCfg.Timeouts = map[string]framework.OperationTimeouts{
		"default": {
			Create: &timeout,
			Update: &timeout,
			Delete: &timeout,
		},
	}

	// Create provider server that runs within test
	ctx := context.Background()
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	tf.SetReattachConfig(reattachCfg)
	tf.WriteConfigT(t, vars.builder.Build())
	_, err := tf.Init()
	require.NoError(t, err)

	// Run `terraform apply` to create project and database, which should
	// timeout at project creation
	out, err := tf.Apply()
	defer tf.DestroySilently()
	require.Error(t, err)

	// Check expected error message
	require.Contains(t, string(out), "Unable to achieve desired state for project")
	require.Contains(t, string(out), "Timed out after "+timeout)

	// Check that state contains project resource despite readiness
	// failure, but not project data source
	tf.CheckStateResource(t, "nuodbaas_project.proj").
		HasAttributeValue("status.state", string(openapi.ProjectStatusModelStateCreating)).
		HasAttributeValue("status.ready", false)
	data, err := tf.GetStateResource("data.nuodbaas_project.proj")
	require.NoError(t, err)
	require.Nil(t, data)

	// Disable readiness check for project by specifying timeout=0, and
	// specify timeout for database explicitly
	noTimeout := "0"
	vars.providerCfg.Timeouts = map[string]framework.OperationTimeouts{
		"project": {
			Create: &noTimeout,
			Update: &noTimeout,
			Delete: &noTimeout,
		},
		"database": {
			Create: &timeout,
			Update: &timeout,
			Delete: &timeout,
		},
	}
	tf.WriteConfigT(t, vars.builder.Build())

	// Run `terraform apply` to re-create project and create database.
	// Project creation should succeed because timeout=0 was specified, but
	// the database readiness check should timeout.
	out, err = tf.Apply()
	require.Error(t, err)

	// Check expected error message
	require.Contains(t, string(out), "Unable to achieve desired state for database")
	require.Contains(t, string(out), "Timed out after "+timeout)

	// Check attributes in data sources
	tf.CheckStateResource(t, "data.nuodbaas_project.proj").
		HasAttributeValue("status.state", string(openapi.ProjectStatusModelStateCreating)).
		HasAttributeValue("status.ready", false)
	// Check that state contains database resource despite readiness
	// failure, but not database data source
	tf.CheckStateResource(t, "nuodbaas_database.db").
		HasAttributeValue("status.state", string(openapi.DatabaseStatusModelStateCreating)).
		HasAttributeValue("status.ready", false)
	data, err = tf.GetStateResource("data.nuodbaas_database.db")
	require.NoError(t, err)
	require.Nil(t, data)

	// Disable readiness check for all resources and run `terraform apply`
	vars.providerCfg.Timeouts = map[string]framework.OperationTimeouts{
		"default": {
			Create: &noTimeout,
			Update: &noTimeout,
			Delete: &noTimeout,
		},
	}
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)

	// Check attributes in data sources
	tf.CheckStateResource(t, "data.nuodbaas_project.proj").
		HasAttributeValue("status.state", string(openapi.ProjectStatusModelStateCreating)).
		HasAttributeValue("status.ready", false)
	tf.CheckStateResource(t, "data.nuodbaas_database.db").
		HasAttributeValue("status.state", string(openapi.DatabaseStatusModelStateCreating)).
		HasAttributeValue("status.ready", false)

	// Set project as disabled and set update timeout
	vars.providerCfg.Timeouts = map[string]framework.OperationTimeouts{
		"default": {
			Update: &timeout,
		},
	}
	disabled := true
	vars.project.Maintenance = &openapi.MaintenanceModel{IsDisabled: &disabled}
	tf.WriteConfigT(t, vars.builder.Build())

	// Run `terraform apply` and check that update fails with timeout
	out, err = tf.Apply()
	require.Error(t, err)

	// Check expected error message
	require.Contains(t, string(out), "Unable to achieve desired state for project")
	require.Contains(t, string(out), "Timed out after "+timeout)
}

func TestNegative(t *testing.T) {
	vars := newTestVars()

	// Create provider server that runs within test
	ctx := context.Background()
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	tf.SetReattachConfig(reattachCfg)
	tf.WriteConfigT(t, vars.builder.Build())
	_, err := tf.Init()
	require.NoError(t, err)
	defer tf.DestroySilently()

	// Try to import invalid path for project
	out, err := tf.Run("import", "nuodbaas_project.proj", "too/many/parts")
	require.Error(t, err)
	require.Contains(t, string(out), "Unexpected Import Identifier")
	require.Contains(t, string(out), `Expected an id with format "organization/name".`)

	// Try to import invalid path for database
	out, err = tf.Run("import", "nuodbaas_database.db", "too/few")
	require.Error(t, err)
	require.Contains(t, string(out), "Unexpected Import Identifier")
	require.Contains(t, string(out), `Expected an id with format "organization/project/name".`)

	// Try to import resource not in config
	out, err = tf.Run("import", "nuodbaas_project.nonexistent", "org/proj")
	require.Error(t, err)
	require.Contains(t, string(out), `resource address "nuodbaas_project.nonexistent" does not exist in the configuration.`)

	// Try to import resource not in remote state
	out, err = tf.Run("import", "nuodbaas_database.db", "org/proj/nonexistent")
	require.Error(t, err)
	require.Contains(t, string(out), "Cannot import non-existent remote object")

	// Omit a required attribute and run `terraform apply`
	vars.database.DbaPassword = nil
	tf.WriteConfigT(t, vars.builder.Build())
	out, err = tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Missing required argument")
	dbaPassword := "password"
	vars.database.DbaPassword = &dbaPassword

	// Specify an invalid attribute and run `terraform apply`
	vars.project.Sla = ""
	tf.WriteConfigT(t, vars.builder.Build())
	out, err = tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Unable to create project")
	require.Contains(t, string(out), "400 Bad Request")
	vars.project.Sla = "dev"

	// Specify a read-only attribute
	caPem := "..."
	vars.project.Status = &openapi.ProjectStatusModel{CaPem: &caPem}
	tf.WriteConfigT(t, vars.builder.Build())
	out, err = tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Invalid Configuration for Read-Only Attribute")
	vars.project.Status = nil

	// Specify a data source without a dependency and run `terraform apply`.
	// This should fail with 404 Not Found.
	vars.builder.WithProjectDataSource("nodep", &ProjectNameModel{
		Organization: "org",
		Name:         "proj",
	})
	tf.WriteConfigT(t, vars.builder.Build())
	out, err = tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Unable to read project")
	require.Contains(t, string(out), "404 Not Found")
	vars.builder.WithoutProjectDataSource("nodep")

	// Specify read-only data source attributes
	vars.builder.WithProjectsDataSource("project_list", &ProjectsDataSourceModel{
		Projects: []ProjectNameModel{
			{
				Organization: "org",
				Name:         "proj",
			},
		},
	})
	tf.WriteConfigT(t, vars.builder.Build())
	out, err = tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Invalid Configuration for Read-Only Attribute")
	vars.builder.WithoutProjectsDataSource("project_list")

	// Specify a database resource without a project dependency and run
	// `terraform apply`. This should fail with 404 Not Found.
	vars.builder.WithDatabaseResource("nodep", &DatabaseResourceModel{
		Organization: "org",
		Project:      "proj",
		Name:         "nodep",
		DbaPassword:  &dbaPassword,
	})
	tf.WriteConfigT(t, vars.builder.Build())
	out, err = tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Unable to create database")
	require.Contains(t, string(out), "404 Not Found")
	vars.builder.WithoutDatabaseResource("nodep")

	// Specify a duplicate project resource and run `terraform apply`. This
	// should fail with 409 Conflict.
	vars.builder.WithProjectResource("dup", &vars.project)
	tf.WriteConfigT(t, vars.builder.Build())
	out, err = tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Unable to create project")
	require.Contains(t, string(out), "409 Conflict")
	vars.builder.WithoutProjectResource("dup")

	// Run `terraform apply` on valid config to create resources
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)
	// Sample some values in state file to validate
	tf.CheckStateResource(t, "nuodbaas_database.db").
		HasAttributeValue("tier", "n0.nano").
		HasAttributeValue("status.state", string(openapi.DatabaseStatusModelStateAvailable)).
		HasAttributeValue("status.ready", true)
	tf.CheckStateResource(t, "data.nuodbaas_database.db").
		HasAttributeValue("tier", "n0.nano").
		HasAttributeValue("status.state", string(openapi.DatabaseStatusModelStateAvailable)).
		HasAttributeValue("status.ready", true)
	tf.CheckStateResource(t, "nuodbaas_project.proj").
		HasAttributeValue("tier", "n0.nano").
		HasAttributeValue("status.state", string(openapi.ProjectStatusModelStateAvailable)).
		HasAttributeValue("status.ready", true)
	tf.CheckStateResource(t, "data.nuodbaas_project.proj").
		HasAttributeValue("tier", "n0.nano").
		HasAttributeValue("status.state", string(openapi.ProjectStatusModelStateAvailable)).
		HasAttributeValue("status.ready", true)

	// Try to import resource already being managed
	out, err = tf.Run("import", "nuodbaas_database.db", "org/proj/db")
	require.Error(t, err)
	require.Contains(t, string(out), "Resource already managed by Terraform")

	t.Run("invalidProviderConfiguration", func(t *testing.T) {
		defer func() {
			vars.providerCfg = NuoDbaasProviderModel{}
		}()

		// Configure invalid URL and verify that reads fail
		badUrl := "http://badhost/"
		vars.providerCfg.UrlBase = &badUrl
		tf.WriteConfigT(t, vars.builder.Build())
		out, err = tf.Apply()
		require.Error(t, err)
		require.Contains(t, string(out), "Unable to read project")
		require.Contains(t, string(out), "Unable to read projects")
		require.Contains(t, string(out), "Unable to read databases")

		// Specify bad timeout values
		noSuffix := "999"
		negative := "-1s"
		vars.providerCfg.UrlBase = nil
		vars.providerCfg.Timeouts = map[string]framework.OperationTimeouts{
			"badresource": {},
			"database":    {Create: &noSuffix, Update: &negative},
		}
		tf.WriteConfigT(t, vars.builder.Build())
		out, err = tf.Apply()
		require.Error(t, err)
		require.Contains(t, string(out), "Invalid timeout for database create")
		require.Contains(t, string(out), `missing unit in duration "999"`)
		require.Contains(t, string(out), "Timeout for database update is negative: -1s")
		require.Contains(t, string(out), "Invalid resource type: badresource")
	})

	// Specify an invalid attribute and run `terraform apply`. This should
	// cause an update to be attempted that fails.
	vars.project.Tier = ""
	tf.WriteConfigT(t, vars.builder.Build())
	out, err = tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Unable to update project")
	require.Contains(t, string(out), "400 Bad Request")
	vars.project.Tier = "n0.nano"

	t.Run("deleteProjectWithDatabase", func(t *testing.T) {
		client, err := vars.providerCfg.CreateClient()
		require.NoError(t, err)

		// Create an unmanaged database that should cause `terraform
		// destroy` to fail
		dbaPassword := "db"
		database := DatabaseResourceModel{
			Organization: "org",
			Project:      "proj",
			Name:         "unmanaged",
			DbaPassword:  &dbaPassword,
		}
		err = database.Create(ctx, client)
		require.NoError(t, err)
		defer func() {
			var timeoutSeconds int32 = 10
			_, _ = client.DeleteDatabase(
				ctx, database.Organization, database.Project, database.Name,
				&openapi.DeleteDatabaseParams{
					TimeoutSeconds: &timeoutSeconds,
				})
		}()

		// Run `terraform destroy`, which should fail on project
		// deletion because it contains the unmanaged database
		tf.WriteConfigT(t, vars.builder.Build())
		out, err = tf.Destroy()
		require.Error(t, err)
		require.Contains(t, string(out), "Unable to delete project")
		require.Contains(t, string(out), "409 Conflict")
	})

	// Run `terraform destroy` again, which should succeed now that
	// unmanaged database has been deleted
	out, err = tf.Destroy()
	require.NoError(t, err)
}

func TestImport(t *testing.T) {
	vars := newTestVars()

	// Create provider server that runs within test
	ctx := context.Background()
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	tf.SetReattachConfig(reattachCfg)
	tf.WriteConfigT(t, vars.builder.Build())
	_, err := tf.Init()
	require.NoError(t, err)
	defer tf.DestroySilently()

	// Create a project and database by directly invoking the REST service
	client, err := vars.providerCfg.CreateClient()
	require.NoError(t, err)

	// Create project
	project := ProjectResourceModel{
		Organization: "org",
		Name:         "proj",
		Sla:          "dev",
		Tier:         "n0.nano",
	}
	err = project.Create(ctx, client)
	require.NoError(t, err)
	defer func() {
		var timeoutSeconds int32 = 10
		_, _ = client.DeleteProject(
			ctx, project.Organization, project.Name,
			&openapi.DeleteProjectParams{TimeoutSeconds: &timeoutSeconds})
	}()

	// Create database
	database := DatabaseResourceModel{
		Organization: "org",
		Project:      "proj",
		Name:         "db",
		DbaPassword:  vars.database.DbaPassword,
		// Include label that is not in config
		Labels: &map[string]string{
			"color": "blue",
		},
	}
	err = database.Create(ctx, client)
	require.NoError(t, err)
	defer func() {
		var timeoutSeconds int32 = 10
		_, _ = client.DeleteDatabase(
			ctx, database.Organization, database.Project, database.Name,
			&openapi.DeleteDatabaseParams{TimeoutSeconds: &timeoutSeconds})
	}()

	// Run `terraform apply` and verify that it fails due to the resources
	// already existing
	out, err := tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Unable to create project")
	require.Contains(t, string(out), "409 Conflict")

	// Run `terraform refresh` and verify that it does not do anything,
	// since it is only operating on the Terraform state
	out, err = tf.Run("refresh")
	require.NoError(t, err)
	require.Contains(t, string(out), "Empty or non-existent state")

	// Run `terraform import` for project and database
	out, err = tf.Run("import", "nuodbaas_project.proj", "org/proj")
	require.NoError(t, err)
	require.Contains(t, string(out), "Import successful!")
	out, err = tf.Run("import", "nuodbaas_database.db", "org/proj/db")
	require.NoError(t, err)
	require.Contains(t, string(out), "Import successful!")

	// Verify that project and database are in state
	tf.CheckStateResource(t, "nuodbaas_project.proj").
		HasAttributeValue("organization", "org").
		HasAttributeValue("name", "proj").
		HasAttributeValue("tier", "n0.nano")
	tf.CheckStateResource(t, "nuodbaas_database.db").
		HasAttributeValue("organization", "org").
		HasAttributeValue("project", "proj").
		HasAttributeValue("name", "db").
		HasAttributeValue("tier", "n0.nano").
		HasAttributeValue("labels", map[string]any{"color": "blue"})

	// Run `terraform apply` and verify that there is nothing to do
	out, err = tf.Apply()
	require.NoError(t, err)
	// TODO(asz6): On import, we do not have the DBA password in the state
	// since it is populated by the 'GET /databases' response, and the DBA
	// password is marked as required, so the configured value always
	// differs from the state.
	//require.Contains(t, string(out), "No changes.")
	//require.Contains(t, string(out), "Your infrastructure matches the configuration.")
}
