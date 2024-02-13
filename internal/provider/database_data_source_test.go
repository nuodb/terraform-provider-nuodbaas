/* (C) Copyright 2016-2024 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	nuodbaas_client_test "github.com/nuodb/terraform-provider-nuodbaas/internal/client/testclient"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"
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
	tierParamKey := "foo"
	tierParamValue := "bar"

	model := openapi.DatabaseCreateUpdateModel{
		DbaPassword: &password,
		Properties: &openapi.DatabasePropertiesModel{
			ArchiveDiskSize: &archiveDisk,
			JournalDiskSize: &journalDisk,
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
					resource.TestCheckResourceAttr(resourcePath, "maintenance.is_disabled", strconv.FormatBool(disabled)),
					resource.TestCheckResourceAttr(resourcePath, "properties.archive_disk_size", archiveDisk),
					resource.TestCheckResourceAttr(resourcePath, "properties.journal_disk_size", journalDisk),
					resource.TestCheckResourceAttr(resourcePath, "properties.tier_parameters.%", "1"),
					resource.TestCheckResourceAttr(resourcePath, "properties.tier_parameters."+tierParamKey, tierParamValue),
					resource.TestCheckResourceAttrSet(resourcePath, "status.sql_endpoint"),
					// Flaky under envtest
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
