/* (C) Copyright 2016-2024 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	nuodbaas_client_test "github.com/nuodb/terraform-provider-nuodbaas/internal/client/testclient"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/database"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAccDatabaseDataSource(t *testing.T) {
	var (
		organizationName = "org"
		projectName      = "proj"
		dbName           = "db"
		sla              = "dev"
		tier             = "n0.nano"
		dbaPassword      = "pass"
		disabled         = true
		archiveDisk      = "1Gi"
		journalDisk      = "2Gi"
		productVersion   = "5.1"
		tierParamKey     = "foo"
		tierParamValue   = "bar"
	)

	model := openapi.DatabaseCreateUpdateModel{
		DbaPassword: &dbaPassword,
		Labels: &map[string]string{
			"key0": "value0",
			"key1": "value1",
		},
		Properties: &openapi.DatabasePropertiesModel{
			ArchiveDiskSize: &archiveDisk,
			JournalDiskSize: &journalDisk,
			ProductVersion:  &productVersion,
			TierParameters:  &map[string]string{tierParamKey: tierParamValue},
		},
		Maintenance: &openapi.MaintenanceModel{
			IsDisabled: &disabled,
		},
	}

	ctx := context.TODO()

	client, err := nuodbaas_client_test.DefaultApiClient()
	require.NoError(t, err)
	require.NoError(t, nuodbaas_client_test.CreateProject(t, ctx, client, organizationName, projectName, sla, tier))
	require.NoError(t, nuodbaas_client_test.CreateDatabaseWithModel(t, ctx, client, organizationName, projectName, dbName, model))

	resourceName := "database_details"
	resourcePath := fmt.Sprintf("data.%s.%s", getDatabaseDatasourceTypeName(), resourceName)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
					data "%s" "%s" {
						name         = "%s"
						organization = "%s"
						project      = "%s"
					}
				`, getDatabaseDatasourceTypeName(), resourceName, dbName, organizationName, projectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePath, "organization", organizationName),
					resource.TestCheckResourceAttr(resourcePath, "project", projectName),
					resource.TestCheckResourceAttr(resourcePath, "name", dbName),
					resource.TestCheckResourceAttr(resourcePath, "tier", tier),
					resource.TestCheckResourceAttr(resourcePath, "labels.%", "2"),
					resource.TestCheckResourceAttr(resourcePath, "labels.key0", "value0"),
					resource.TestCheckResourceAttr(resourcePath, "labels.key1", "value1"),
					resource.TestCheckResourceAttr(resourcePath, "maintenance.is_disabled", strconv.FormatBool(disabled)),
					resource.TestCheckResourceAttr(resourcePath, "properties.archive_disk_size", archiveDisk),
					resource.TestCheckResourceAttr(resourcePath, "properties.journal_disk_size", journalDisk),
					resource.TestCheckResourceAttr(resourcePath, "properties.product_version", productVersion),
					resource.TestCheckResourceAttr(resourcePath, "properties.tier_parameters.%", "1"),
					resource.TestCheckResourceAttr(resourcePath, "properties.tier_parameters."+tierParamKey, tierParamValue),
					resource.TestCheckResourceAttrSet(resourcePath, "status.sql_endpoint"),
					resource.TestCheckNoResourceAttr(resourcePath, "dba_password"),
					resource.TestCheckNoResourceAttr(resourcePath, "resource_version"),
					// TODO(asz6): This is flaky because the mock controller that simulates the
					// DBaaS operator sets the database as ready asynchronously with generating the
					// CA certificate. This makes it is possible for the database to be marked as
					// ready before the CA PEM is available, which is unrealistic. We can uncomment
					// this once the mock controller is updated to have more realistic behavior.
					// resource.TestCheckResourceAttrSet(resourcePath, "status.ca_pem"),
				),
			},
		},
	})
}

func getDatabaseDatasourceTypeName() string {
	source := NewDatabaseDataSource()

	ctx := context.TODO() // Not used

	request := datasource.MetadataRequest{ProviderTypeName: getProviderTypeName()}
	response := datasource.MetadataResponse{}
	source.Metadata(ctx, request, &response)
	return response.TypeName
}
