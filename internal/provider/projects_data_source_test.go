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

func TestAccProjectsDataSourceEmpty(t *testing.T) {
	resourceName := "project_list"
	resourcePath := fmt.Sprintf("data.%s.%s", getProjectsDatasourceTypeName(), resourceName)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
					data "%s" "%s" {}
				`, getProjectsDatasourceTypeName(), resourceName),
				Check: resource.TestCheckResourceAttr(resourcePath, "projects.#", "0"),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
					data "%s" "%s" {
						filter {
							organization = "someorg"
						}
					}
				`, getProjectsDatasourceTypeName(), resourceName),
				Check: resource.TestCheckResourceAttr(resourcePath, "projects.#", "0"),
			},
		},
	})
}

func TestAccProjectsDataSourceNotEmpty(t *testing.T) {
	resourceName := "project_list"
	resourcePath := fmt.Sprintf("data.%s.%s", getProjectsDatasourceTypeName(), resourceName)

	proj1Org := "org1"
	proj1Name := "proj1"
	proj2Org := "org2"
	proj2Name := "proj2"

	// TestCheckFunc that verifies the project list regardless of returned order
	outOfOrderCheck := func(s *terraform.State) error {
		projList, ok := s.RootModule().Resources[resourcePath]
		if !ok {
			return errors.New("Could not get projectsList")
		}

		firstOrg, ok := projList.Primary.Attributes["projects.0.organization"]
		if !ok {
			return errors.New("Could not get the first organization in the projectsList")
		}

		var proj1path, proj2path string
		if firstOrg == proj1Org {
			proj1path = "projects.0"
			proj2path = "projects.1"
		} else {
			proj1path = "projects.1"
			proj2path = "projects.0"
		}

		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resourcePath, proj1path+".name", proj1Name),
			resource.TestCheckResourceAttr(resourcePath, proj1path+".organization", proj1Org),
			resource.TestCheckResourceAttr(resourcePath, proj2path+".name", proj2Name),
			resource.TestCheckResourceAttr(resourcePath, proj2path+".organization", proj2Org),
		)(s)
	}

	// Create a couple of test projects
	client := nuodbaas_client_test.NewTestClient(context.TODO())
	require.NoError(t, client.CreateProject(t, proj1Org, proj1Name, "dev", "n0.nano"))
	require.NoError(t, client.CreateProject(t, proj2Org, proj2Name, "dev", "n0.nano"))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
					data "%s" "%s" {}
				`, getProjectsDatasourceTypeName(), resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePath, "projects.#", "2"), // Fail fast if there is not 2 projects
					outOfOrderCheck,
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
					data "%s" "%s" {
						filter {
							organization = "%s"
						}
					}
				`, getProjectsDatasourceTypeName(), resourceName, proj1Org),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePath, "projects.#", "1"),
					resource.TestCheckResourceAttr(resourcePath, "projects.0.name", proj1Name),
					resource.TestCheckResourceAttr(resourcePath, "projects.0.organization", proj1Org),
				),
			},
		},
	})
}

func getProjectsDatasourceTypeName() string {
	source := projectsDataSource{}

	ctx := context.TODO() // Not used

	request := datasource.MetadataRequest{ProviderTypeName: getProviderTypeName()}
	response := datasource.MetadataResponse{}
	source.Metadata(ctx, request, &response)
	return response.TypeName
}
