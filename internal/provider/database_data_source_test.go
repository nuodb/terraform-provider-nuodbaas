/* (C) Copyright 2016-2024 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	nuodbaas "github.com/nuodb/terraform-provider-nuodbaas/client"
	nuodbaas_client_test "github.com/nuodb/terraform-provider-nuodbaas/internal/client/testclient"
	"github.com/stretchr/testify/require"
)

func TestAccDatabaseDataSource(t *testing.T) {
	organizationName := "org"
	projectName := "proj"
	sla := "dev"
	tier := "n0.nano"
	dbName := "db"

	password := "pass"
	archiveDisk := "1Gi"
	journalDisk := "2Gi"
	disabled := true
	expAt := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	tierParamKey := "foo"
	tierParamValue := "bar"

	model := nuodbaas.DatabaseCreateUpdateModel{
		DbaPassword: &password,
		Properties: &nuodbaas.DatabasePropertiesModel{
			ArchiveDiskSize: &archiveDisk,
			JournalDiskSize: &journalDisk,
			TierParameters:  &map[string]string{tierParamKey: tierParamValue},
		},
		Maintenance: &nuodbaas.MaintenanceModel{
			IsDisabled:    &disabled,
			ExpiresAtTime: &expAt,
		},
	}

	client := nuodbaas_client_test.NewTestClient(context.TODO())
	require.NoError(t, client.CreateProject(t, organizationName, projectName, sla, tier))
	require.NoError(t, client.CreateDatabaseWithModel(t, organizationName, projectName, dbName, model))

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
					resource.TestCheckResourceAttr(resourcePath, "maintenance.is_disabled", strconv.FormatBool(disabled)),
					resource.TestCheckResourceAttr(resourcePath, "maintenance.expires_at", expAt.String()),
					resource.TestCheckResourceAttrSet(resourcePath, "resource_version"),
					resource.TestCheckResourceAttr(resourcePath, "properties.archive_disk_size", archiveDisk),
					resource.TestCheckResourceAttr(resourcePath, "properties.journal_disk_size", journalDisk),
					resource.TestCheckResourceAttr(resourcePath, "properties.tier_parameters.%", "1"),
					resource.TestCheckResourceAttr(resourcePath, "properties.tier_parameters."+tierParamKey, tierParamValue),
					resource.TestCheckResourceAttrSet(resourcePath, "status.sql_end_point"),
					// Flaky under envtest
					// resource.TestCheckResourceAttrSet(resourcePath, "status.ca_pem"),
				),
			},
		},
	})
}

func getDatabaseDatasourceTypeName() string {
	source := databaseDataSource{}

	ctx := context.TODO() // Not used

	request := datasource.MetadataRequest{ProviderTypeName: getProviderTypeName()}
	response := datasource.MetadataResponse{}
	source.Metadata(ctx, request, &response)
	return response.TypeName
}
