// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package provider_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/database"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/project"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	semver "github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"
)

type TestOption string

const (
	PAUSE_OPERATOR_COMMAND       TestOption = "PAUSE_OPERATOR_COMMAND"
	RESUME_OPERATOR_COMMAND      TestOption = "RESUME_OPERATOR_COMMAND"
	CONTAINER_SCHEDULING_ENABLED TestOption = "CONTAINER_SCHEDULING_ENABLED"
	WEBHOOKS_ENABLED             TestOption = "WEBHOOKS_ENABLED"
	ORGANIZATION_BOUND_USER      TestOption = "ORGANIZATION_BOUND_USER"
	RESOURCE_CREATE_TIMEOUT      TestOption = "RESOURCE_CREATE_TIMEOUT"
	RESOURCE_UPDATE_TIMEOUT      TestOption = "RESOURCE_UPDATE_TIMEOUT"
	RESOURCE_DELETE_TIMEOUT      TestOption = "RESOURCE_DELETE_TIMEOUT"
	NUODB_CP_VERSION             TestOption = "NUODB_CP_VERSION"
)

func (option TestOption) Get() string {
	return os.Getenv(string(option))
}

func (option TestOption) IsTrue() bool {
	return option.Get() == "true"
}

func (option TestOption) IsFalse() bool {
	return option.Get() == "false"
}

func (option TestOption) GetSemver() (*semver.Version, error) {
	return semver.NewVersion(option.Get())
}

func (option TestOption) IsVersionLessThan(version string) bool {
	thisVersion, err := option.GetSemver()
	if err != nil {
		return false
	}
	suppliedVersion, err := semver.NewVersion(version)
	return err == nil && thisVersion.LessThan(suppliedVersion)
}

func PauseOperator(t *testing.T) {
	// Pausing Operator when webhooks are enabled prevents CRUD operations
	// from being performed
	if WEBHOOKS_ENABLED.IsTrue() {
		t.Skip("Cannot pause operator if webhooks are enabled")
	}
	pauseCmd := PAUSE_OPERATOR_COMMAND.Get()
	if pauseCmd == "" {
		t.Skip("Pausing operator is not supported")
	}
	resumeCmd := RESUME_OPERATOR_COMMAND.Get()
	if resumeCmd == "" {
		t.Skip("Resuming operator is not supported")
	}
	// Pause reconciliation by stopping Operator
	cmd := exec.Command(pauseCmd)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, out)
	// On test teardown, resume Operator
	t.Cleanup(func() {
		cmd := exec.Command(resumeCmd)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("Unexpected error [%s]: %s", err, out)
		}
	})
}

func getOrganization() string {
	user := os.Getenv("NUODB_CP_USER")
	if user == "" {
		return "org"
	}
	parts := strings.Split(user, "/")
	return parts[0]
}

func withRandomSuffix(name string) string {
	suffix := rand.Intn(1000) //nolint:gosec // This is not a security concern
	return fmt.Sprintf("%s%d", name, suffix)
}

func setDefaultResourceTimeouts(providerCfg *NuoDbaasProviderModel) {
	var timeouts framework.OperationTimeouts
	if timeout := RESOURCE_CREATE_TIMEOUT.Get(); timeout != "" {
		timeouts.Create = &timeout
	}
	if timeout := RESOURCE_UPDATE_TIMEOUT.Get(); timeout != "" {
		timeouts.Update = &timeout
	}
	if timeout := RESOURCE_DELETE_TIMEOUT.Get(); timeout != "" {
		timeouts.Delete = &timeout
	}
	providerCfg.Timeouts = map[string]framework.OperationTimeouts{
		framework.DEFAULT_RESOURCE: timeouts,
	}
}

type testVars struct {
	providerCfg NuoDbaasProviderModel
	project     ProjectResourceModel
	database    DatabaseResourceModel
	builder     *TfConfigBuilder
}

func (vars *testVars) resetVars() {
	dbaPassword := "dba"
	vars.providerCfg = NuoDbaasProviderModel{}
	setDefaultResourceTimeouts(&vars.providerCfg)
	// Get organization from user name
	orgName := getOrganization()
	// Generate a random project name to avoid collisions
	projectName := withRandomSuffix("proj")
	vars.project = ProjectResourceModel{
		Organization: orgName,
		Name:         projectName,
		Sla:          "dev",
		Tier:         "n0.nano",
	}
	vars.database = DatabaseResourceModel{
		Organization: orgName,
		Project:      projectName,
		Name:         "db",
		DbaPassword:  &dbaPassword,
	}
	vars.builder = NewTfConfigBuilder().WithProviderConfig("nuodbaas", &vars.providerCfg).
		WithProjectResource("proj", &vars.project).
		WithDatabaseResource("db", &vars.database, "nuodbaas_project.proj").
		WithProjectDataSource("proj", &ProjectNameModel{
			Organization: orgName,
			Name:         projectName,
		}, "nuodbaas_project.proj").
		WithDatabaseDataSource("db", &DatabaseNameModel{
			Organization: orgName,
			Project:      projectName,
			Name:         "db",
		}, "nuodbaas_database.db").
		WithProjectsDataSource("proj_list", &ProjectsDataSourceModel{}).
		WithDatabasesDataSource("db_list", &DatabasesDataSourceModel{})
}

func newTestVars(overrideTimeouts bool) *testVars {
	var ret testVars
	ret.resetVars()
	// If overrideTimeouts=true and this is an end-to-end test, accelerate
	// test by skipping readiness checks.
	if overrideTimeouts && CONTAINER_SCHEDULING_ENABLED.IsTrue() {
		ret.providerCfg.Timeouts = map[string]framework.OperationTimeouts{
			"default": {
				Create: ptr("0"),
				Update: ptr("0"),
			},
		}
	}
	return &ret
}

func ptr[T any](v T) *T {
	return &v
}

func TestFullLifecycle(t *testing.T) {
	vars := newTestVars(false)

	// Create provider server that runs within test
	ctx := context.Background()
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	err := tf.SetReattachConfig(reattachCfg)
	require.NoError(t, err)
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Init()
	require.NoError(t, err)

	// Run `terraform apply` to create project and database
	_, err = tf.Apply()
	defer tf.DestroySilently()
	require.NoError(t, err)

	// Use client created from provider config to populate structs with
	// current state
	client, err := vars.providerCfg.CreateClient()
	require.NoError(t, err)
	actualProject := vars.project
	err = actualProject.Read(ctx, client)
	require.NoError(t, err)
	require.Equal(t, vars.project.Organization, actualProject.Organization)
	require.Equal(t, vars.project.Name, actualProject.Name)
	require.Equal(t, "dev", actualProject.Sla)
	require.Equal(t, "n0.nano", actualProject.Tier)
	require.NotNil(t, actualProject.Status)
	require.NotNil(t, actualProject.Status.CaPem)
	require.NotNil(t, actualProject.Status.State)
	require.Equal(t, openapi.ProjectStatusModelStateAvailable, *actualProject.Status.State)

	actualDatabase := vars.database
	err = actualDatabase.Read(ctx, client)
	require.NoError(t, err)
	require.Equal(t, vars.project.Organization, actualDatabase.Organization)
	require.Equal(t, vars.project.Name, actualDatabase.Project)
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
		HasAttributeValue("organization", vars.project.Organization).
		HasAttributeValue("name", vars.project.Name).
		HasAttributeValue("labels", map[string]any{}).
		HasAttribute("properties.product_version").
		HasAttributeValue("status.state", string(openapi.ProjectStatusModelStateAvailable)).
		HasAttributeValue("status.ready", true).
		HasAttributeValue("status.shutdown", false)
	tf.CheckStateResource(t, "data.nuodbaas_database.db").
		HasAttributeValue("organization", vars.project.Organization).
		HasAttributeValue("project", vars.project.Name).
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
		HasListAttributeContaining("projects", map[string]any{
			"organization": vars.project.Organization,
			"name":         vars.project.Name,
		})
	tf.CheckStateResource(t, "data.nuodbaas_databases.db_list").
		HasListAttributeContaining("databases", map[string]any{
			"organization": vars.project.Organization,
			"project":      vars.project.Name,
			"name":         "db",
		})

	// Update database config attributes (tier, labels, and product_version)
	// and execute `terraform apply`
	tier := "n0.small"
	vars.database.Tier = &tier
	vars.database.Labels = &map[string]string{
		"priority": "high",
	}
	productVersion := "5.1"
	var expectedProductVersion string
	// Avoid triggering rolling upgrade if real processes are being started,
	// since this will result in the database transitioning back and forth
	// between Modifying and Available
	if !CONTAINER_SCHEDULING_ENABLED.IsTrue() {
		expectedProductVersion = productVersion
		vars.database.Properties = &openapi.DatabasePropertiesModel{
			ProductVersion: &productVersion,
		}
	} else {
		require.NotNil(t, actualDatabase.Properties)
		require.NotNil(t, actualDatabase.Properties.ProductVersion)
		expectedProductVersion = *actualDatabase.Properties.ProductVersion
	}
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)
	err = actualDatabase.Read(ctx, client)
	require.NoError(t, err)
	require.NotNil(t, actualDatabase.Tier)
	require.Equal(t, tier, *actualDatabase.Tier)
	require.Equal(t, vars.database.Labels, actualDatabase.Labels)
	require.NotNil(t, actualDatabase.Properties)
	if !CONTAINER_SCHEDULING_ENABLED.IsTrue() {
		require.Equal(t, vars.database.Properties.ProductVersion, actualDatabase.Properties.ProductVersion)
	}

	// Update project config attributes (tier, labels, and product_version)
	// and execute `terraform apply`
	vars.project.Tier = tier
	vars.project.Labels = &map[string]string{
		"priority": "high",
	}
	// Avoid triggering rolling upgrade if real processes are being started
	if !CONTAINER_SCHEDULING_ENABLED.IsTrue() {
		vars.project.Properties = &openapi.ProjectPropertiesModel{
			ProductVersion: &productVersion,
		}
	}
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)
	err = actualProject.Read(ctx, client)
	require.NoError(t, err)
	require.Equal(t, tier, actualProject.Tier)
	require.Equal(t, vars.project.Labels, actualProject.Labels)
	require.NotNil(t, actualProject.Properties)
	if !CONTAINER_SCHEDULING_ENABLED.IsTrue() {
		require.Equal(t, vars.project.Properties.ProductVersion, actualProject.Properties.ProductVersion)
	}

	// Check attributes in data sources
	tf.CheckStateResource(t, "data.nuodbaas_project.proj").
		HasAttributeValue("organization", vars.project.Organization).
		HasAttributeValue("name", vars.project.Name).
		HasAttributeValue("tier", tier).
		HasAttributeValue("labels", map[string]any{"priority": "high"}).
		HasAttributeValue("properties.product_version", expectedProductVersion).
		HasAttributeValue("status.state", string(openapi.ProjectStatusModelStateAvailable)).
		HasAttributeValue("status.ready", true).
		HasAttributeValue("status.shutdown", false)
	tf.CheckStateResource(t, "data.nuodbaas_database.db").
		HasAttributeValue("organization", vars.project.Organization).
		HasAttributeValue("project", vars.project.Name).
		HasAttributeValue("name", "db").
		HasAttributeValue("tier", tier).
		HasAttributeValue("labels", map[string]any{"priority": "high"}).
		HasAttributeValue("properties.product_version", expectedProductVersion).
		HasAttributeValue("status.state", string(openapi.DatabaseStatusModelStateAvailable)).
		HasAttributeValue("status.ready", true).
		HasAttributeValue("status.shutdown", false)

	// Set project and database as disabled and also update labels
	vars.project.Maintenance = &openapi.MaintenanceModel{
		IsDisabled: ptr(true),
	}
	// Omit labels from project, which should not cause them to be removed
	// because unknown values are resolved from state
	vars.project.Labels = nil
	vars.database.Maintenance = &openapi.MaintenanceModel{
		IsDisabled: ptr(true),
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
		HasAttributeValue("organization", vars.project.Organization).
		HasAttributeValue("name", vars.project.Name).
		HasAttributeValue("tier", tier).
		HasAttributeValue("labels", map[string]any{"priority": "high"}).
		HasAttributeValue("properties.product_version", expectedProductVersion).
		HasAttributeValue("maintenance.is_disabled", true).
		HasAttributeValue("status.state", string(openapi.ProjectStatusModelStateStopped)).
		HasAttributeValue("status.shutdown", true)
	tf.CheckStateResource(t, "data.nuodbaas_database.db").
		HasAttributeValue("organization", vars.project.Organization).
		HasAttributeValue("project", vars.project.Name).
		HasAttributeValue("name", "db").
		HasAttributeValue("tier", tier).
		HasAttributeValue("labels", map[string]any{}).
		HasAttributeValue("properties.product_version", expectedProductVersion).
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

func TestAttributeSerialization(t *testing.T) {
	if WEBHOOKS_ENABLED.IsTrue() {
		t.Skip("Do not test attributes exhaustively in end-to-end configuration, which may reject some settings")
	}

	vars := newTestVars(false)
	// Disable readiness checks since some configurations may not be valid
	// and result in resources never becoming ready
	vars.providerCfg.Timeouts = map[string]framework.OperationTimeouts{
		"default": {
			Create: ptr("0"),
			Update: ptr("0"),
		},
	}

	// Create provider server that runs within test
	ctx := context.Background()
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	err := tf.SetReattachConfig(reattachCfg)
	require.NoError(t, err)
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Init()
	require.NoError(t, err)
	defer tf.DestroySilently()

	// Run `terraform apply` to create resources
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)

	// Check attributes in project resource and data source
	checkResourceAndDataSource(t, "nuodbaas_project.proj", tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", vars.project.Organization).
			HasAttributeValue("name", vars.project.Name).
			HasAttributeValue("sla", "dev").
			HasAttributeValue("tier", "n0.nano").
			HasAttributeValue("labels", map[string]any{}).
			HasAttribute("properties.product_version").
			DoesNotHaveAttribute("maintenance").
			HasAttribute("status.state").
			HasAttribute("status.ready").
			HasAttribute("status.shutdown")
	})
	// Check attributes in database resource and data source
	checkResourceAndDataSource(t, "nuodbaas_database.db", tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", vars.project.Organization).
			HasAttributeValue("project", vars.project.Name).
			HasAttributeValue("name", "db").
			HasAttributeValue("tier", "n0.nano").
			HasAttributeValue("labels", map[string]any{}).
			HasAttribute("properties.product_version").
			DoesNotHaveAttribute("maintenance").
			HasAttribute("status.state").
			HasAttribute("status.ready").
			HasAttribute("status.shutdown")
	})

	// Save original database and project structs
	var (
		originalDatabase = vars.database
		originalProject  = vars.project
	)

	// Populate optional fields of database
	vars.database.Labels = &map[string]string{
		"one": "1",
		"two": "2",
	}
	vars.database.Maintenance = &openapi.MaintenanceModel{
		IsDisabled: ptr(true),
	}
	vars.database.Properties = &openapi.DatabasePropertiesModel{
		ArchiveDiskSize: ptr("100Gi"),
		JournalDiskSize: ptr("10Gi"),
		TierParameters: &map[string]string{
			"zones":      `["us-east-1", "us-east-2"]`,
			"smReplicas": "2",
			"teReplicas": "5",
		},
		ProductVersion: ptr("6.0"),
	}

	// Populate optional fields fo project
	vars.project.Labels = vars.database.Labels
	// Do not populate sub-attributes of maintenance to exercise values
	// being computed by the server
	vars.project.Maintenance = &openapi.MaintenanceModel{}
	vars.project.Properties = &openapi.ProjectPropertiesModel{
		TierParameters: &map[string]string{
			"zones":         `["us-east-1", "us-east-2"]`,
			"adminReplicas": "3",
		},
		ProductVersion: vars.database.Properties.ProductVersion,
	}

	// Run `terraform apply` to update resources
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)

	// Check attributes in project resource and data source
	checkResourceAndDataSource(t, "nuodbaas_project.proj", tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", vars.project.Organization).
			HasAttributeValue("name", vars.project.Name).
			HasAttributeValue("sla", "dev").
			HasAttributeValue("tier", "n0.nano").
			HasAttributeValue("labels", map[string]any{"one": "1", "two": "2"}).
			HasAttributeValue("properties.tier_parameters.zones", `["us-east-1", "us-east-2"]`).
			HasAttributeValue("properties.tier_parameters.adminReplicas", "3").
			HasAttributeValue("properties.product_version", "6.0").
			HasAttributeValue("maintenance.is_disabled", false).
			HasAttribute("status.state").
			HasAttribute("status.ready").
			HasAttribute("status.shutdown")
	})
	// Check attributes in database resource and data source
	checkResourceAndDataSource(t, "nuodbaas_database.db", tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", vars.project.Organization).
			HasAttributeValue("project", vars.project.Name).
			HasAttributeValue("name", "db").
			HasAttributeValue("tier", "n0.nano").
			HasAttributeValue("labels", map[string]any{"one": "1", "two": "2"}).
			HasAttributeValue("properties.archive_disk_size", "100Gi").
			HasAttributeValue("properties.journal_disk_size", "10Gi").
			HasAttributeValue("properties.tier_parameters.zones", `["us-east-1", "us-east-2"]`).
			HasAttributeValue("properties.tier_parameters.smReplicas", "2").
			HasAttributeValue("properties.tier_parameters.teReplicas", "5").
			HasAttributeValue("properties.product_version", "6.0").
			HasAttributeValue("maintenance.is_disabled", true).
			HasAttribute("status.state").
			HasAttribute("status.ready").
			HasAttribute("status.shutdown")
	})

	// Revert project and database to original, sparsely-populated settings
	// and `terraform apply` to update resources. This should not change the
	// resources because these are optional, computed attributes that are
	// resolved from Terraform state when unconfigured.
	vars.project = originalProject
	vars.database = originalDatabase
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)

	// Check attributes in project resource and data source
	checkResourceAndDataSource(t, "nuodbaas_project.proj", tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", vars.project.Organization).
			HasAttributeValue("name", vars.project.Name).
			HasAttributeValue("sla", "dev").
			HasAttributeValue("tier", "n0.nano").
			HasAttributeValue("labels", map[string]any{"one": "1", "two": "2"}).
			HasAttributeValue("properties.tier_parameters.zones", `["us-east-1", "us-east-2"]`).
			HasAttributeValue("properties.tier_parameters.adminReplicas", "3").
			HasAttributeValue("properties.product_version", "6.0").
			HasAttributeValue("maintenance.is_disabled", false).
			HasAttribute("status.state").
			HasAttribute("status.ready").
			HasAttribute("status.shutdown")
	})
	// Check attributes in database resource and data source
	checkResourceAndDataSource(t, "nuodbaas_database.db", tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", vars.project.Organization).
			HasAttributeValue("project", vars.project.Name).
			HasAttributeValue("name", "db").
			HasAttributeValue("tier", "n0.nano").
			HasAttributeValue("labels", map[string]any{"one": "1", "two": "2"}).
			HasAttributeValue("properties.archive_disk_size", "100Gi").
			HasAttributeValue("properties.journal_disk_size", "10Gi").
			HasAttributeValue("properties.tier_parameters.zones", `["us-east-1", "us-east-2"]`).
			HasAttributeValue("properties.tier_parameters.smReplicas", "2").
			HasAttributeValue("properties.tier_parameters.teReplicas", "5").
			HasAttributeValue("properties.product_version", "6.0").
			HasAttributeValue("maintenance.is_disabled", true).
			HasAttribute("status.state").
			HasAttribute("status.ready").
			HasAttribute("status.shutdown")
	})

	// Explicitly set labels and tier_parameters to empty maps so that they
	// get cleared
	vars.project.Labels = &map[string]string{}
	vars.database.Labels = &map[string]string{}
	vars.project.Properties = &openapi.ProjectPropertiesModel{TierParameters: &map[string]string{}}
	vars.database.Properties = &openapi.DatabasePropertiesModel{TierParameters: &map[string]string{}}
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)

	// Check attributes in project resource and data source
	checkResourceAndDataSource(t, "nuodbaas_project.proj", tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", vars.project.Organization).
			HasAttributeValue("name", vars.project.Name).
			HasAttributeValue("sla", "dev").
			HasAttributeValue("tier", "n0.nano").
			HasAttributeValue("labels", map[string]any{}).
			HasAttributeValue("properties.tier_parameters", map[string]any{}).
			HasAttributeValue("properties.product_version", "6.0").
			HasAttributeValue("maintenance.is_disabled", false).
			HasAttribute("status.state").
			HasAttribute("status.ready").
			HasAttribute("status.shutdown")
	})
	// Check attributes in database resource and data source
	checkResourceAndDataSource(t, "nuodbaas_database.db", tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", vars.project.Organization).
			HasAttributeValue("project", vars.project.Name).
			HasAttributeValue("name", "db").
			HasAttributeValue("tier", "n0.nano").
			HasAttributeValue("labels", map[string]any{}).
			HasAttributeValue("properties.archive_disk_size", "100Gi").
			HasAttributeValue("properties.journal_disk_size", "10Gi").
			HasAttributeValue("properties.tier_parameters", map[string]any{}).
			HasAttributeValue("properties.product_version", "6.0").
			HasAttributeValue("maintenance.is_disabled", true).
			HasAttribute("status.state").
			HasAttribute("status.ready").
			HasAttribute("status.shutdown")
	})
}

func checkResourceAndDataSource(t *testing.T, address string, tf *TfHelper, assertFn func(*AttributeChecker)) {
	assertFn(tf.CheckStateResource(t, address))
	assertFn(tf.CheckStateResource(t, "data."+address))
}

func TestTimeouts(t *testing.T) {
	vars := newTestVars(false)

	// Disable reconciliation
	PauseOperator(t)

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
	err := tf.SetReattachConfig(reattachCfg)
	require.NoError(t, err)
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Init()
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
	vars.project.Maintenance = &openapi.MaintenanceModel{IsDisabled: ptr(true)}
	tf.WriteConfigT(t, vars.builder.Build())

	// Run `terraform apply` and check that update fails with timeout
	out, err = tf.Apply()
	require.Error(t, err)

	// Check expected error message
	require.Contains(t, string(out), "Unable to achieve desired state for project")
	require.Contains(t, string(out), "Timed out after "+timeout)
}

func TestNegative(t *testing.T) {
	vars := newTestVars(true)

	// Create provider server that runs within test
	ctx := context.Background()
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	err := tf.SetReattachConfig(reattachCfg)
	require.NoError(t, err)
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Init()
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
	out, err = tf.Run("import", "nuodbaas_project.nonexistent", vars.project.Organization+"/proj")
	require.Error(t, err)
	require.Contains(t, string(out), `resource address "nuodbaas_project.nonexistent" does not exist in the configuration.`)

	// Try to import resource not in remote state
	out, err = tf.Run("import", "nuodbaas_database.db", vars.project.Organization+"/proj/nonexistent")
	require.Error(t, err)
	require.Contains(t, string(out), "Cannot import non-existent remote object")

	// Specify an invalid attribute and run `terraform apply`
	vars.project.Sla = ""
	tf.WriteConfigT(t, vars.builder.Build())
	out, err = tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Unable to create project")
	require.Contains(t, string(out), "400 Bad Request")
	vars.project.Sla = "dev"

	// Specify a read-only attribute
	vars.project.Status = &openapi.ProjectStatusModel{CaPem: ptr("...")}
	tf.WriteConfigT(t, vars.builder.Build())
	out, err = tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Invalid Configuration for Read-Only Attribute")
	vars.project.Status = nil

	// Specify a data source without a dependency and run `terraform apply`.
	// This should fail with 404 Not Found.
	vars.builder.WithProjectDataSource("nodep", &ProjectNameModel{
		Organization: vars.project.Organization,
		Name:         vars.project.Name,
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
				Organization: vars.project.Organization,
				Name:         vars.project.Name,
			},
		},
	})
	tf.WriteConfigT(t, vars.builder.Build())
	out, err = tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Invalid Configuration for Read-Only Attribute")
	vars.builder.WithoutProjectsDataSource("project_list")

	// Specify an invalid database filter
	vars.builder.WithDatabasesDataSource("database_list", &DatabasesDataSourceModel{
		Filter: &DatabaseFilterModel{
			Project: ptr(vars.project.Name),
		},
	})
	tf.WriteConfigT(t, vars.builder.Build())
	out, err = tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Unable to read databases")
	require.Contains(t, string(out), "Cannot specify project filter ("+vars.project.Name+") without organization")
	vars.builder.WithoutDatabasesDataSource("database_list")

	// Specify a database resource without a project dependency and run
	// `terraform apply`. This should fail with 404 Not Found. Give project
	// non-existent name because there is small chance that org/proj gets
	// created in time
	vars.builder.WithDatabaseResource("nodep", &DatabaseResourceModel{
		Organization: vars.project.Organization,
		Project:      "nonexistent",
		Name:         "nodep",
		DbaPassword:  vars.database.DbaPassword,
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
		HasAttribute("status.ca_pem").
		HasAttribute("status.state").
		HasAttribute("status.ready")
	tf.CheckStateResource(t, "data.nuodbaas_database.db").
		HasAttributeValue("tier", "n0.nano").
		HasAttribute("status.ca_pem").
		HasAttribute("status.state").
		HasAttribute("status.ready")
	tf.CheckStateResource(t, "nuodbaas_project.proj").
		HasAttributeValue("sla", "dev").
		HasAttributeValue("tier", "n0.nano").
		HasAttribute("status.state").
		HasAttribute("status.ready")
	tf.CheckStateResource(t, "data.nuodbaas_project.proj").
		HasAttributeValue("sla", "dev").
		HasAttributeValue("tier", "n0.nano").
		HasAttribute("status.state").
		HasAttribute("status.ready")

	// Try to import resource already being managed
	out, err = tf.Run("import", "nuodbaas_database.db", vars.project.Organization+"/"+vars.project.Name+"/db")
	require.Error(t, err)
	tfName := "Terraform"
	if USE_TOFU.IsTrue() {
		tfName = "OpenTofu"
	}
	require.Contains(t, string(out), "Resource already managed by "+tfName)

	t.Run("invalidProviderConfiguration", func(t *testing.T) {
		// Clear any credentials set via environment variables
		t.Setenv("NUODB_CP_USER", "")
		t.Setenv("NUODB_CP_PASSWORD", "")
		t.Setenv("NUODB_CP_TOKEN", "")

		// Revert provider configuration when finished
		defer func() {
			vars.providerCfg = NuoDbaasProviderModel{}
		}()

		// Configure invalid credentials and verify that reads fail
		vars.providerCfg.User = ptr("org/user")
		vars.providerCfg.Password = ptr("badpassword")
		tf.WriteConfigT(t, vars.builder.Build())
		out, err = tf.Apply()
		require.Error(t, err)
		require.Contains(t, string(out), "Unable to read project")
		require.Contains(t, string(out), "Unable to read projects")
		require.Contains(t, string(out), "Unable to read databases")

		// Configure unreachable URL and verify that reads fail
		vars.providerCfg.UrlBase = ptr("http://unreachable/")
		tf.WriteConfigT(t, vars.builder.Build())
		out, err = tf.Apply()
		require.Error(t, err)
		require.Contains(t, string(out), "Unable to read project")
		require.Contains(t, string(out), "Unable to read projects")
		require.Contains(t, string(out), "Unable to read databases")
		vars.providerCfg.UrlBase = nil

		// Specify bad timeout values
		noSuffix := "999"
		negative := "-1s"
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
		vars.providerCfg.Timeouts = nil

		// Temporarily override environment variable NUODB_CP_URL_BASE
		// so that URL is not specified at all
		t.Setenv(NUODB_CP_URL_BASE, "")
		out, err = tf.Apply()
		require.Error(t, err)
		require.Contains(t, string(out), "Must specify url_base or the environment variable "+NUODB_CP_URL_BASE)

		// Specify an invalid URL (bad port)
		vars.providerCfg.UrlBase = ptr("http://host:-80")
		tf.WriteConfigT(t, vars.builder.Build())
		out, err = tf.Apply()
		require.Error(t, err)
		require.Contains(t, string(out), "invalid port")

		// Specify an invalid URL (no scheme)
		vars.providerCfg.UrlBase = ptr("badurl")
		tf.WriteConfigT(t, vars.builder.Build())
		out, err = tf.Apply()
		require.Error(t, err)
		require.Contains(t, string(out), "No scheme found in URL")
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
			Organization: vars.project.Organization,
			Project:      vars.project.Name,
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

	if WEBHOOKS_ENABLED.IsFalse() && !NUODB_CP_VERSION.IsVersionLessThan("2.5.0") {
		t.Run("failedProjectAndDatabase", func(t *testing.T) {
			// Specify non-existent project tier and check that it goes into failed state
			vars.project.Tier = "n5.huge"
			tf.WriteConfigT(t, vars.builder.Build())
			out, err = tf.Apply()
			require.Error(t, err)
			require.Contains(t, string(out), fmt.Sprintf("Project %s/%s failed: ", vars.project.Organization, vars.project.Name))

			// Specify non-existent database tier and check that it goes into failed state
			vars.project.Tier = "n0.nano"
			vars.database.Tier = ptr("n5.huge")
			tf.WriteConfigT(t, vars.builder.Build())
			out, err = tf.Apply()
			require.Error(t, err)
			require.Contains(t, string(out), fmt.Sprintf("Database %s/%s/%s failed: ", vars.database.Organization, vars.database.Project, vars.database.Name))
			// Revert tier change
			vars.database.Tier = nil
		})
	}

	// Run `terraform destroy` again, which should succeed now that
	// unmanaged database has been deleted
	tf.WriteConfigT(t, vars.builder.Build())
	out, err = tf.Destroy()
	require.NoError(t, err)
}

func TestImmutableAttributeChange(t *testing.T) {
	vars := newTestVars(true)

	// Create provider server that runs within test
	ctx := context.Background()
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	err := tf.SetReattachConfig(reattachCfg)
	require.NoError(t, err)
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Init()
	require.NoError(t, err)
	defer tf.DestroySilently()

	// Run `terraform apply` to create resources
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)
	// Sample some values in state file to validate
	tf.CheckStateResource(t, "nuodbaas_database.db").
		HasAttributeValue("tier", "n0.nano").
		HasAttribute("status.ca_pem").
		HasAttribute("status.state").
		HasAttribute("status.ready")
	tf.CheckStateResource(t, "data.nuodbaas_database.db").
		HasAttributeValue("tier", "n0.nano").
		HasAttribute("status.ca_pem").
		HasAttribute("status.state").
		HasAttribute("status.ready")
	tf.CheckStateResource(t, "nuodbaas_project.proj").
		HasAttributeValue("sla", "dev").
		HasAttributeValue("tier", "n0.nano").
		HasAttribute("status.state").
		HasAttribute("status.ready")
	tf.CheckStateResource(t, "data.nuodbaas_project.proj").
		HasAttributeValue("sla", "dev").
		HasAttributeValue("tier", "n0.nano").
		HasAttribute("status.state").
		HasAttribute("status.ready")

	// Check whether the REST server supports updating the DBA password
	client, err := vars.providerCfg.CreateClient()
	require.NoError(t, err)
	resp, err := client.UpdateDbaPassword(
		ctx, vars.database.Organization, vars.database.Project, vars.database.Name, nil,
		openapi.UpdateDbaPasswordModel{Current: *vars.database.DbaPassword})
	require.NoError(t, err)

	// "404 Not Found" with no "detail" message indicates that /dbaPassword
	// sub-resource does not exist. Run appropriate test case based on what
	// behavior is supported.
	err = helper.ParseResponse(resp, nil)
	if IsDbaPasswordUpdateUnsupportedError(resp, err) {
		t.Run("dbaPasswordChangeRejected", func(t *testing.T) {
			// Change DBA password and run `terraform apply`. Revert
			// DBA password change when finished.
			originalPassword := vars.database.DbaPassword
			vars.database.DbaPassword = ptr("updated")
			defer func() {
				vars.database.DbaPassword = originalPassword
			}()

			tf.WriteConfigT(t, vars.builder.Build())
			out, err := tf.Apply()
			// Password change should be rejected if the REST server
			// does not support it
			require.Error(t, err)
			require.Contains(t, string(out), "Configured DBA password was changed")
		})
	} else if CONTAINER_SCHEDULING_ENABLED.IsTrue() {
		require.NoError(t, err)
		t.Run("dbaPasswordChange", func(t *testing.T) {
			// Change DBA password and run `terraform apply`
			vars.database.DbaPassword = ptr("updated")
			tf.WriteConfigT(t, vars.builder.Build())
			out, err := tf.Apply()
			// Check that DBA password change succeeded
			require.NoError(t, err)
			require.Contains(t, string(out), "0 to add, 1 to change, 0 to destroy.")
		})
	} else {
		require.NoError(t, err)
		t.Run("dbaPasswordChangeTimeout", func(t *testing.T) {
			// Change DBA password and run `terraform apply`
			vars.database.DbaPassword = ptr("updated")
			// Expect readiness check to fail due to DBA password
			// not being updated and specify small timeout
			timeouts := vars.providerCfg.Timeouts
			defer func() {
				vars.providerCfg.Timeouts = timeouts
			}()
			vars.providerCfg.Timeouts = map[string]framework.OperationTimeouts{
				"database": {
					Update: ptr("2s"),
				},
			}
			tf.WriteConfigT(t, vars.builder.Build())
			out, err := tf.Apply()
			// Check that readiness check failed due to DBA password
			require.Error(t, err)
			require.Contains(t, string(out), "DBA password for database "+vars.project.Organization+"/"+vars.project.Name+"/db has not been updated")
			require.Contains(t, string(out), "0 to add, 1 to change, 0 to destroy.")
		})
	}

	// Change the project SLA and verify that a warning is displayed by
	// Terraform when running `terraform plan` unless the environment
	// variable ALLOW_DESTRUCTIVE_REPLACE=true is set
	t.Run("planSlaChange", func(t *testing.T) {
		vars.project.Sla = "qa"
		defer func() {
			vars.project.Sla = "dev"
		}()

		tf.WriteConfigT(t, vars.builder.Build())
		out, err := tf.Plan()
		require.NoError(t, err)
		require.Contains(t, string(out), "Immutable Attribute Change")
		require.Contains(t, string(out), "`sla`")
		require.Contains(t, string(out), "(`\"dev\"`)")
		require.NotContains(t, string(out), "Unable to update project")

		t.Setenv(framework.ALLOW_DESTRUCTIVE_REPLACE_VAR, "true")
		tf.WriteConfigT(t, vars.builder.Build())
		out, err = tf.Plan()
		require.NoError(t, err)
		require.Contains(t, string(out), "# forces replacement")
		require.Contains(t, string(out), "1 to add, 0 to change, 1 to destroy.")
	})

	// Run `terraform apply` and verify that a warning is displayed by
	// Terraform and that the request is also rejected by the server
	t.Run("applySlaChangeRejected", func(t *testing.T) {
		vars.project.Sla = "qa"
		defer func() {
			vars.project.Sla = "dev"
		}()

		tf.WriteConfigT(t, vars.builder.Build())
		out, err := tf.Apply()
		require.Error(t, err)
		require.Contains(t, string(out), "Immutable Attribute Change")
		require.Contains(t, string(out), "`sla`")
		require.Contains(t, string(out), "(`\"dev\"`)")
		require.Contains(t, string(out), "Unable to update project")
	})

	// Explicitly replace the project and database by running `terraform
	// destroy -target=...` followed by `terraform apply`
	t.Run("applySlaChangeExplicitly", func(t *testing.T) {
		// Change project SLA
		vars.project.Sla = "qa"
		tf.WriteConfigT(t, vars.builder.Build())

		// Destroy the project and database
		out, err := tf.Run("destroy", "-target=nuodbaas_project.proj", "-auto-approve")
		require.NoError(t, err)
		require.NotContains(t, string(out), "Immutable Attribute Change")
		require.Contains(t, string(out), "Destroy complete! Resources: 2 destroyed.")

		// Re-create the project and database
		out, err = tf.Apply()
		require.NoError(t, err)
		require.NotContains(t, string(out), "Immutable Attribute Change")
		require.Contains(t, string(out), "Apply complete! Resources: 2 added, 0 changed, 0 destroyed.")

		// Validate Terraform state and check that project has updated SLA value
		tf.CheckStateResource(t, "nuodbaas_database.db").
			HasAttributeValue("tier", "n0.nano").
			HasAttribute("status.state").
			HasAttribute("status.ready")
		tf.CheckStateResource(t, "data.nuodbaas_database.db").
			HasAttributeValue("tier", "n0.nano").
			HasAttribute("status.state").
			HasAttribute("status.ready")
		tf.CheckStateResource(t, "nuodbaas_project.proj").
			HasAttributeValue("sla", "qa").
			HasAttributeValue("tier", "n0.nano").
			HasAttribute("status.state").
			HasAttribute("status.ready")
		tf.CheckStateResource(t, "data.nuodbaas_project.proj").
			HasAttributeValue("sla", "qa").
			HasAttributeValue("tier", "n0.nano").
			HasAttribute("status.state").
			HasAttribute("status.ready")
	})
}

func TestImport(t *testing.T) {
	vars := newTestVars(true)

	// Create a project and database by directly invoking the REST service
	client, err := vars.providerCfg.CreateClient()
	require.NoError(t, err)

	// Create project
	project := ProjectResourceModel{
		Organization: vars.project.Organization,
		Name:         vars.project.Name,
		Sla:          "dev",
		Tier:         "n0.nano",
	}
	ctx := context.Background()
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
		Organization: vars.project.Organization,
		Project:      vars.project.Name,
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

	// Create provider server that runs within test
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	err = tf.SetReattachConfig(reattachCfg)
	require.NoError(t, err)
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Init()
	require.NoError(t, err)
	defer tf.DestroySilently()

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
	out, err = tf.Run("import", "nuodbaas_project.proj", vars.project.Organization+"/"+vars.project.Name)
	require.NoError(t, err)
	require.Contains(t, string(out), "Import successful!")
	out, err = tf.Run("import", "nuodbaas_database.db", vars.project.Organization+"/"+vars.project.Name+"/db")
	require.NoError(t, err)
	require.Contains(t, string(out), "Import successful!")

	// Verify that project and database are in state
	tf.CheckStateResource(t, "nuodbaas_project.proj").
		HasAttributeValue("organization", vars.project.Organization).
		HasAttributeValue("name", vars.project.Name).
		HasAttributeValue("tier", "n0.nano")
	tf.CheckStateResource(t, "nuodbaas_database.db").
		HasAttributeValue("organization", vars.project.Organization).
		HasAttributeValue("project", vars.project.Name).
		HasAttributeValue("name", "db").
		HasAttributeValue("tier", "n0.nano").
		HasAttributeValue("labels", map[string]any{"color": "blue"})

	// Run `terraform apply` and verify that it fails due to the DBA
	// password being in the configuration. The presence of the dba_password
	// attribute triggers an unnecessary update because it is not in state
	// after `terraform import`, which is populated by a `GET /databases`
	// response.
	out, err = tf.Apply()
	require.Error(t, err)
	require.Contains(t, string(out), "Cannot update DBA password on database "+vars.project.Organization+"/"+vars.project.Name+"/db since current password is unknown")

	// Remove DBA password from configuration, which is not needed because
	// the database is already created
	vars.database.DbaPassword = nil
	tf.WriteConfigT(t, vars.builder.Build())

	// Run `terraform apply` and verify that there is nothing to do
	out, err = tf.Apply()
	require.NoError(t, err)
	require.Contains(t, string(out), "No changes.")
	require.Contains(t, string(out), "Your infrastructure matches the configuration.")
}

func TestDataSourceFiltering(t *testing.T) {
	if ORGANIZATION_BOUND_USER.IsTrue() {
		t.Skipf("Current user is bound to organization")
	}

	var vars testVars
	vars.resetVars()
	// Disable readiness checks
	vars.providerCfg.Timeouts = map[string]framework.OperationTimeouts{
		"default": {
			Create: ptr("0"),
			Update: ptr("0"),
		},
	}

	// Create a projects and databases by directly invoking the REST service
	client, err := vars.providerCfg.CreateClient()
	require.NoError(t, err)
	ctx := context.Background()

	numProjects := 5
	numDatabases := 5
	if CONTAINER_SCHEDULING_ENABLED.IsTrue() {
		numProjects = 2
		numDatabases = 2
	}

	// Create projects and databases
	for i := 0; i != numProjects; i += 1 {
		projectName := fmt.Sprintf("proj%d", i)
		project := ProjectResourceModel{
			Organization: "unmanaged",
			Name:         projectName,
			Sla:          "dev",
			Tier:         "n0.nano",
			Labels: &map[string]string{
				"type": "unmanaged",
				"name": projectName,
			},
		}
		err = project.Create(ctx, client)
		require.NoError(t, err)
		defer func() {
			var timeoutSeconds int32 = 10
			_, _ = client.DeleteProject(
				ctx, project.Organization, project.Name,
				&openapi.DeleteProjectParams{TimeoutSeconds: &timeoutSeconds})
		}()

		// Create databases
		for j := 0; j != numDatabases; j += 1 {
			dbName := fmt.Sprintf("db%d", j)
			database := DatabaseResourceModel{
				Organization: "unmanaged",
				Project:      projectName,
				Name:         dbName,
				DbaPassword:  ptr("password"),
				Labels: &map[string]string{
					"type": "unmanaged",
					"name": dbName,
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
		}
	}

	// Add data sources that list projects and databases using various filters
	vars.builder = vars.builder.
		WithDatabasesDataSource("all", &DatabasesDataSourceModel{}).
		WithDatabasesDataSource("unmanaged", &DatabasesDataSourceModel{Filter: &DatabaseFilterModel{Organization: ptr("unmanaged")}}).
		WithDatabasesDataSource("proj0", &DatabasesDataSourceModel{Filter: &DatabaseFilterModel{Organization: ptr("unmanaged"), Project: ptr("proj0")}}).
		WithDatabasesDataSource("name_label", &DatabasesDataSourceModel{Filter: &DatabaseFilterModel{Labels: []string{"name"}}}).
		WithDatabasesDataSource("name_label_db0", &DatabasesDataSourceModel{Filter: &DatabaseFilterModel{Labels: []string{"name=db0"}}}).
		WithDatabasesDataSource("name_label_negative", &DatabasesDataSourceModel{Filter: &DatabaseFilterModel{Labels: []string{"!name"}}}).
		WithDatabasesDataSource("multiple_labels", &DatabasesDataSourceModel{Filter: &DatabaseFilterModel{Labels: []string{"name!=db0", "type"}}}).
		WithProjectsDataSource("all", &ProjectsDataSourceModel{}).
		WithProjectsDataSource("unmanaged", &ProjectsDataSourceModel{Filter: &ProjectFilterModel{Organization: ptr("unmanaged")}}).
		WithProjectsDataSource("name_label", &ProjectsDataSourceModel{Filter: &ProjectFilterModel{Labels: []string{"name"}}}).
		WithProjectsDataSource("name_label_proj0", &ProjectsDataSourceModel{Filter: &ProjectFilterModel{Labels: []string{"name=proj0"}}}).
		WithProjectsDataSource("name_label_negative", &ProjectsDataSourceModel{Filter: &ProjectFilterModel{Labels: []string{"!name"}}}).
		WithProjectsDataSource("multiple_labels", &ProjectsDataSourceModel{Filter: &ProjectFilterModel{Labels: []string{"name!=proj0", "type"}}})

	// Create provider server that runs within test
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	err = tf.SetReattachConfig(reattachCfg)
	require.NoError(t, err)
	tf.WriteConfigT(t, vars.builder.Build())
	_, err = tf.Init()
	require.NoError(t, err)
	defer tf.DestroySilently()

	// Apply Terraform config
	_, err = tf.Apply()
	require.NoError(t, err)
	// Refresh Terraform state so that the project and database managed by
	// Terraform is also populated in the data sources.
	_, err = tf.Run("refresh")
	require.NoError(t, err)

	// Check each database list data source
	managedDatabases := 1
	unmanagedDatabases := numProjects * numDatabases
	totalDatabases := unmanagedDatabases + managedDatabases
	checkDataSourceList(t, "databases", "all", totalDatabases, tf, func(ac *AttributeChecker) {})
	checkDataSourceList(t, "databases", "unmanaged", unmanagedDatabases, tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", "unmanaged")
	})
	checkDataSourceList(t, "databases", "proj0", numDatabases, tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", "unmanaged").HasAttributeValue("project", "proj0")
	})
	checkDataSourceList(t, "databases", "name_label", unmanagedDatabases, tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", "unmanaged")
	})
	checkDataSourceList(t, "databases", "name_label_db0", numProjects, tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("name", "db0")
	})
	checkDataSourceList(t, "databases", "name_label_negative", managedDatabases, tf, func(ac *AttributeChecker) {
		ac.DoesNotHaveAttributeValue("organization", "unmanaged")
	})
	checkDataSourceList(t, "databases", "multiple_labels", unmanagedDatabases-numProjects, tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", "unmanaged").
			DoesNotHaveAttributeValue("name", "db0")
	})

	// Check each project list data source
	managedProjects := 1
	unmanagedProjects := numProjects
	totalProjects := unmanagedProjects + managedProjects
	checkDataSourceList(t, "projects", "all", totalProjects, tf, func(ac *AttributeChecker) {})
	checkDataSourceList(t, "projects", "unmanaged", unmanagedProjects, tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", "unmanaged")
	})
	checkDataSourceList(t, "projects", "name_label", unmanagedProjects, tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", "unmanaged")
	})
	checkDataSourceList(t, "projects", "name_label_proj0", 1, tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("name", "proj0")
	})
	checkDataSourceList(t, "projects", "name_label_negative", managedProjects, tf, func(ac *AttributeChecker) {
		ac.DoesNotHaveAttributeValue("organization", "unmanaged")
	})
	checkDataSourceList(t, "projects", "multiple_labels", unmanagedProjects-1, tf, func(ac *AttributeChecker) {
		ac.HasAttributeValue("organization", "unmanaged").
			DoesNotHaveAttributeValue("name", "proj0")
	})
}

func checkDataSourceList(t *testing.T, dataSourceType, dataSource string, expectedCount int, tf *TfHelper, assertFn func(*AttributeChecker)) {
	address := fmt.Sprintf("data.nuodbaas_%s.%s", dataSourceType, dataSource)
	t.Run(address, func(t *testing.T) {
		tf.CheckStateResource(t, address).ForEach(dataSourceType, expectedCount, assertFn)
	})
}

func TestValidation(t *testing.T) {
	// Create provider server that runs within test
	ctx := context.Background()
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	err := tf.SetReattachConfig(reattachCfg)
	require.NoError(t, err)

	t.Run("invalid project name", func(t *testing.T) {
		vars := newTestVars(false)
		projName := "this is not a valid project name"

		vars.project.Name = projName

		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		out, err := tf.Validate()
		require.Error(t, err)

		require.Contains(t, string(out), "must match pattern: ^[a-z][a-z0-9]*$")
		require.Contains(t, string(out), projName)
	})

	t.Run("invalid database product version", func(t *testing.T) {
		vars := newTestVars(false)
		productVersion := "six"

		vars.database.Properties = &openapi.DatabasePropertiesModel{
			ProductVersion: &productVersion,
		}

		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		out, err := tf.Validate()
		require.Error(t, err)

		require.Contains(t, string(out), "must match pattern:")
		require.Contains(t, string(out), "^([1-9][0-9]*|[1-9][0-9]*\\.[0-9]+|[1-9][0-9]*\\.[0-9]+\\.[0-9]+)([._-][a-z0-9._-]+)?$")
		require.Contains(t, string(out), productVersion)
	})

	t.Run("partial credentials", func(t *testing.T) {
		// Clear any credentials that might exist in the environment, for example when running as an e2e test
		t.Setenv(NUODB_CP_USER, "")
		t.Setenv(NUODB_CP_PASSWORD, "")
		t.Setenv(NUODB_CP_TOKEN, "")

		vars := newTestVars(false)

		errorString := "Partial credentials"
		errorDescription := "To use basic authentication, both user name and password should be provided"

		// Test user without a password
		vars.providerCfg.User = ptr("org/user")

		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		out, err := tf.Validate()
		require.Error(t, err)

		require.Contains(t, string(out), errorString)
		require.Contains(t, string(out), errorDescription)

		vars.providerCfg.User = nil

		// And passing user via the environment
		t.Setenv(NUODB_CP_USER, "org/user")

		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		out, err = tf.Validate()
		require.Error(t, err)

		require.Contains(t, string(out), errorString)
		require.Contains(t, string(out), errorDescription)

		t.Setenv(NUODB_CP_USER, "")

		// Test password without a user name
		vars.providerCfg.Password = ptr("password")

		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		out, err = tf.Validate()
		require.Error(t, err)

		require.Contains(t, string(out), errorString)
		require.Contains(t, string(out), errorDescription)

		vars.providerCfg.Password = nil

		// And passing password via the environment
		t.Setenv(NUODB_CP_PASSWORD, "password")

		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		out, err = tf.Validate()
		require.Error(t, err)

		require.Contains(t, string(out), errorString)
		require.Contains(t, string(out), errorDescription)
	})

	t.Run("multiple credentials", func(t *testing.T) {
		// Specify both basic and token credentials via environment variables
		t.Setenv(NUODB_CP_USER, "org/user")
		t.Setenv(NUODB_CP_PASSWORD, "pass")
		t.Setenv(NUODB_CP_TOKEN, "token")
		vars := newTestVars(false)
		tf.WriteConfigT(t, vars.builder.Build())

		errorString := "Multiple credentials"
		errorDescription := "Both basic and token authentication credentials were supplied."

		// Run `terraform validate`
		out, err := tf.Validate()
		require.Error(t, err)
		require.Contains(t, string(out), errorString)
		require.Contains(t, string(out), errorDescription)

		// Specify both basic and token credentials via config
		t.Setenv(NUODB_CP_USER, "")
		t.Setenv(NUODB_CP_PASSWORD, "")
		t.Setenv(NUODB_CP_TOKEN, "")
		vars.providerCfg.User = ptr("org/user")
		vars.providerCfg.Password = ptr("pass")
		vars.providerCfg.Token = ptr("token")
		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		out, err = tf.Validate()
		require.Error(t, err)
		require.Contains(t, string(out), errorString)
		require.Contains(t, string(out), errorDescription)
	})

	t.Run("malformed user", func(t *testing.T) {
		vars := newTestVars(false)
		// Clear any credentials that might exist in the environment, for example when running as an e2e test
		t.Setenv(NUODB_CP_USER, "")
		t.Setenv(NUODB_CP_PASSWORD, "somePassword")
		t.Setenv(NUODB_CP_TOKEN, "")

		errorString := "Malformed user name"
		errorDescription := "User name should be in the format \"<organization>/<user>\"."

		// Test user name without an org
		vars.providerCfg.User = ptr("org.user")
		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		out, err := tf.Validate()
		require.Error(t, err)
		require.Contains(t, string(out), errorString)
		require.Contains(t, string(out), errorDescription)

		// Test user with an empty org
		vars.providerCfg.User = ptr("/user")
		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		out, err = tf.Validate()
		require.Error(t, err)
		require.Contains(t, string(out), errorString)
		require.Contains(t, string(out), errorDescription)

		// Test user with only an org
		vars.providerCfg.User = ptr("org/")
		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		out, err = tf.Validate()
		require.Error(t, err)
		require.Contains(t, string(out), errorString)
		require.Contains(t, string(out), errorDescription)

		vars.providerCfg.User = nil

		// And passing user via the environment
		t.Setenv(NUODB_CP_USER, "orguser")

		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		out, err = tf.Validate()
		require.Error(t, err)

		require.Contains(t, string(out), errorString)
		require.Contains(t, string(out), errorDescription)
	})

	t.Run("validate url and timeout", func(t *testing.T) {
		// There is more extensive testing in TestNegative so only test that
		// they are checked by `terraform validate`
		vars := newTestVars(false)

		// Try an invalid timeout
		vars.providerCfg.Timeouts = map[string]framework.OperationTimeouts{
			"database": {Update: ptr("-1s")},
		}
		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		out, err := tf.Validate()
		require.Error(t, err)
		require.Contains(t, string(out), "Timeout for database update is negative: -1s")
		vars.providerCfg.Timeouts = nil

		// Try an invalid url
		vars.providerCfg.UrlBase = ptr("hostname.com")
		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		out, err = tf.Validate()
		require.Error(t, err)
		require.Contains(t, string(out), "No scheme found in URL")
	})

	t.Run("validate restore_from.backup", func(t *testing.T) {
		vars := newTestVars(false)

		// Specify a partial backup name
		vars.database.RestoreFrom = &openapi.RestoreFromModel{
			Backup: ptr("backup"),
		}
		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		out, err := tf.Validate()
		require.Error(t, err)
		require.Contains(t, string(out), "Attribute restore_from.backup must match pattern:")
		vars.providerCfg.Timeouts = nil

		// Specify fully-qualified backup name
		vars.database.RestoreFrom.Backup = ptr("org/proj/db/backup")
		tf.WriteConfigT(t, vars.builder.Build())

		// Run `terraform validate`
		_, err = tf.Validate()
		require.NoError(t, err)
	})
}
