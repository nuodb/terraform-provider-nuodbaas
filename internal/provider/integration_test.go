package provider

import (
	"context"
	"testing"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"
	"github.com/stretchr/testify/require"
)

var (
	providerCfg NuoDbaasProviderModel
	project     ProjectResourceModel
	database    DatabaseResourceModel
	builder     *TfConfigBuilder
)

func ResetVars() {
	dbaPassword := "dba"
	providerCfg = NuoDbaasProviderModel{}
	project = ProjectResourceModel{
		Organization: "org",
		Name:         "proj",
		Sla:          "dev",
		Tier:         "n0.nano",
	}
	database = DatabaseResourceModel{
		Organization: "org",
		Project:      "proj",
		Name:         "db",
		DbaPassword:  &dbaPassword,
	}
	builder = NewTfConfigBuilder().WithProviderConfig("nuodbaas", &providerCfg).
		WithProjectResource("proj", &project).
		WithDatabaseResource("db", &database, "nuodbaas_project.proj").
		WithProjectDataSource("proj", &ProjectNameModel{
			Organization: "org",
			Name:         "proj",
		}, "nuodbaas_project.proj").
		WithDatabaseDataSource("db", &DatabaseNameModel{
			Organization: "org",
			Project:      "proj",
			Name:         "db",
		}, "nuodbaas_database.db")
}

func TestFullLifecycle(t *testing.T) {
	ResetVars()

	// Create provider server that runs within test
	ctx := context.Background()
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	tf.SetReattachConfig(reattachCfg)
	tf.WriteConfigT(t, builder.Build())
	_, err := tf.Init()
	require.NoError(t, err)

	// Run `terraform apply` to create project and database
	_, err = tf.Apply()
	defer tf.DestroySilently()
	require.NoError(t, err)

	// Use client created from provider config to populate structs with
	// current state. TODO(asz6): Parse `terraform show` output to obtain
	// state from Terraform.
	client, err := providerCfg.CreateClient()
	actualProject := project
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

	actualDatabase := database
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

	// Update database config attributes (tier, labels, and product_version)
	// and execute `terraform apply`
	tier := "n1.small"
	database.Tier = &tier
	database.Labels = &map[string]string{
		"priority": "high",
	}
	productVersion := "6.0"
	database.Properties = &openapi.DatabasePropertiesModel{
		ProductVersion: &productVersion,
	}
	tf.WriteConfigT(t, builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)
	err = actualDatabase.Read(ctx, client)
	require.Equal(t, "n1.small", *actualDatabase.Tier)
	require.Equal(t, database.Labels, actualDatabase.Labels)
	require.NotNil(t, actualDatabase.Properties)
	require.Equal(t, database.Properties.ProductVersion, actualDatabase.Properties.ProductVersion)

	// Update project config attributes (tier, labels, and product_version)
	// and execute `terraform apply`
	project.Tier = tier
	project.Labels = &map[string]string{
		"priority": "high",
	}
	project.Properties = &openapi.ProjectPropertiesModel{
		ProductVersion: &productVersion,
	}
	tf.WriteConfigT(t, builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)
	err = actualProject.Read(ctx, client)
	require.Equal(t, "n1.small", actualProject.Tier)
	require.Equal(t, project.Labels, actualProject.Labels)
	require.NotNil(t, actualProject.Properties)
	require.Equal(t, project.Properties.ProductVersion, actualProject.Properties.ProductVersion)

	// Run `terraform destroy` to delete database and project
	_, err = tf.Destroy()
	require.NoError(t, err)

	// Obtain actual project and database state and check that 404 is returned
	err = actualProject.Read(ctx, client)
	require.True(t, helper.IsNotFound(err), "Unexpected error: "+err.Error())
	err = actualDatabase.Read(ctx, client)
	require.True(t, helper.IsNotFound(err), "Unexpected error: "+err.Error())

}

func TestTimeouts(t *testing.T) {
	ResetVars()

	// Specify 5s timeout for all resources
	timeout5s := "5s"
	providerCfg.Timeouts = map[string]framework.OperationTimeouts{
		"default": {
			Create: &timeout5s,
			Update: &timeout5s,
			Delete: &timeout5s,
		},
	}

	// Create provider server that runs within test
	ctx := context.Background()
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	tf.SetReattachConfig(reattachCfg)
	tf.WriteConfigT(t, builder.Build())
	_, err := tf.Init()
	require.NoError(t, err)

	// Disable reconciliation
	reset := SetMockReconcilePolicy(t, MockReconcilePolicy{MarkAsReady: "false"})
	defer reset()

	// Run `terraform apply` to create project and database, which should
	// timeout at project creation
	out, err := tf.Apply()
	defer tf.DestroySilently()
	require.Error(t, err)

	// Check expected error message
	require.Contains(t, string(out), "Error waiting for project to become ready")
	require.Contains(t, string(out), "Timed out after "+timeout5s)

	// Disable readiness check for project by specifying timeout=0, and
	// specify 5s timeout for database explicitly
	noTimeout := "0"
	providerCfg.Timeouts = map[string]framework.OperationTimeouts{
		"project": {
			Create: &noTimeout,
			Update: &noTimeout,
			Delete: &noTimeout,
		},
		"database": {
			Create: &timeout5s,
			Update: &timeout5s,
			Delete: &timeout5s,
		},
	}
	tf.WriteConfigT(t, builder.Build())

	// Run `terraform apply` to re-create project and create database.
	// Project creation should succeed because timeout=0 was specified, but
	// the database readiness check should timeout.
	out, err = tf.Apply()
	require.Error(t, err)

	// Check expected error message
	require.Contains(t, string(out), "Error waiting for database to become ready")
	require.Contains(t, string(out), "Timed out after "+timeout5s)

	// Disable readiness check for all resources and run `terraform apply`
	providerCfg.Timeouts = map[string]framework.OperationTimeouts{
		"default": {
			Create: &noTimeout,
			Update: &noTimeout,
			Delete: &noTimeout,
		},
	}
	tf.WriteConfigT(t, builder.Build())
	_, err = tf.Apply()
	require.NoError(t, err)

	// Set project as disabled and set update timeout to 5s
	providerCfg.Timeouts = map[string]framework.OperationTimeouts{
		"default": {
			Update: &timeout5s,
		},
	}
	disabled := true
	project.Maintenance = &openapi.MaintenanceModel{IsDisabled: &disabled}
	tf.WriteConfigT(t, builder.Build())

	// Run `terraform apply` and check that update fails with timeout
	out, err = tf.Apply()
	require.Error(t, err)

	// Check expected error message
	require.Contains(t, string(out), "Error waiting for project to become ready")
	require.Contains(t, string(out), "Timed out after "+timeout5s)

}
