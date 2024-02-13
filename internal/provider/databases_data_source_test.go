/* (C) Copyright 2016-2024 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	nuodbaas_client_test "github.com/nuodb/terraform-provider-nuodbaas/internal/client/testclient"
	"github.com/stretchr/testify/require"
)

func TestAccDatabasesDataSourceEmpty(t *testing.T) {
	resourceName := "database_list"
	resourcePath := fmt.Sprintf("data.%s.%s", getDatabasesDatasourceTypeName(), resourceName)

	organization := "org"
	project := "proj"

	// REST service on k8s does not like filtering by non-existent projects. Envtest does not care.
	ctx := context.TODO()
	client, err := nuodbaas_client_test.DefaultApiClient()
	require.NoError(t, err)
	require.NoError(t, nuodbaas_client_test.CreateProject(t, ctx, client, organization, project, "dev", "n0.nano"))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
					data "%s" "%s" {}
				`, getDatabasesDatasourceTypeName(), resourceName),
				Check: resource.TestCheckResourceAttr(resourcePath, "databases.#", "0"),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
					data "%s" "%s" {
						filter {
							organization = "%s"
						}
					}
				`, getDatabasesDatasourceTypeName(), resourceName, organization),
				Check: resource.TestCheckResourceAttr(resourcePath, "databases.#", "0"),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
					data "%s" "%s" {
						filter {
							organization = "%s"
							project = "%s"
						}
					}
				`, getDatabasesDatasourceTypeName(), resourceName, organization, project),
				Check: resource.TestCheckResourceAttr(resourcePath, "databases.#", "0"),
			},
		},
	})
}

func TestAccDatabasesDataSourceFilterOrganization(t *testing.T) {
	db1Org := "org1"
	db1Proj := "proj1"
	db1Name := "db1"
	db2Org := "org2"
	db2Proj := "proj2"
	db2Name := "db2"

	filter := fmt.Sprintf(`filter {
		organization = "%s"
	}`, db1Org)

	testFilter(db1Name, db1Proj, db1Org, db2Name, db2Proj, db2Org, t, filter)
}

func TestAccDatabasesDataSourceFilterProject(t *testing.T) {
	db1Org := "org"
	db1Proj := "proj1"
	db1Name := "db1"
	db2Org := db1Org
	db2Proj := "proj2"
	db2Name := "db2"

	filter := fmt.Sprintf(`filter {
		organization = "%s"
		project = "%s"
	}`, db1Org, db1Proj)

	testFilter(db1Name, db1Proj, db1Org, db2Name, db2Proj, db2Org, t, filter)
}

func testFilter(db1Name string, db1Proj string, db1Org string, db2Name string, db2Proj string, db2Org string, t *testing.T, filter string) {
	resourceName := "database_list"
	resourcePath := fmt.Sprintf("data.%s.%s", getDatabasesDatasourceTypeName(), resourceName)

	// TestCheckFunc that verifies the database list regardless of returned order
	outOfOrderCheck := func(s *terraform.State) error {
		projList, ok := s.RootModule().Resources[resourcePath]
		if !ok {
			return errors.New("Could not get the database list")
		}

		firstDb, ok := projList.Primary.Attributes["databases.0.name"]
		if !ok {
			return errors.New("Could not get the first name in the database list")
		}

		var db1Path, db2Path string
		if firstDb == db1Name {
			db1Path = "databases.0"
			db2Path = "databases.1"
		} else {
			db1Path = "databases.1"
			db2Path = "databases.0"
		}

		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resourcePath, db1Path+".name", db1Name),
			resource.TestCheckResourceAttr(resourcePath, db1Path+".project", db1Proj),
			resource.TestCheckResourceAttr(resourcePath, db1Path+".organization", db1Org),
			resource.TestCheckResourceAttr(resourcePath, db2Path+".name", db2Name),
			resource.TestCheckResourceAttr(resourcePath, db2Path+".project", db2Proj),
			resource.TestCheckResourceAttr(resourcePath, db2Path+".organization", db2Org),
		)(s)
	}

	// Create test databases
	ctx := context.TODO()
	client, err := nuodbaas_client_test.DefaultApiClient()
	require.NoError(t, err)
	require.NoError(t, nuodbaas_client_test.CreateProject(t, ctx, client, db1Org, db1Proj, "dev", "n0.nano"))
	require.NoError(t, nuodbaas_client_test.CreateDatabase(t, ctx, client, db1Org, db1Proj, db1Name, "pass"))
	require.NoError(t, nuodbaas_client_test.CreateProject(t, ctx, client, db2Org, db2Proj, "dev", "n0.nano"))
	require.NoError(t, nuodbaas_client_test.CreateDatabase(t, ctx, client, db2Org, db2Proj, db2Name, "pass"))

	dataSourceDefinition := fmt.Sprintf(`
						data "%s" "%s" {
							%s
						}
		`, getDatabasesDatasourceTypeName(), resourceName, filter)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
					data "%s" "%s" {}
					`, getDatabasesDatasourceTypeName(), resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePath, "databases.#", "2"),
					outOfOrderCheck,
				),
			},
			{
				Config: providerConfig + dataSourceDefinition,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePath, "databases.#", "1"),
					resource.TestCheckResourceAttr(resourcePath, "databases.0.name", db1Name),
					resource.TestCheckResourceAttr(resourcePath, "databases.0.organization", db1Org),
					resource.TestCheckResourceAttr(resourcePath, "databases.0.project", db1Proj),
				),
			},
		},
	})
}

func getDatabasesDatasourceTypeName() string {
	source := NewDatabasesDataSource()

	ctx := context.TODO() // Not used

	request := datasource.MetadataRequest{ProviderTypeName: getProviderTypeName()}
	response := datasource.MetadataResponse{}
	source.Metadata(ctx, request, &response)
	return response.TypeName
}
