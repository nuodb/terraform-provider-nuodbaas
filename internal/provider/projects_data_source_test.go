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
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"
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
						filter = {
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
	ctx := context.TODO()

	client, err := nuodbaas_client_test.DefaultApiClient()
	require.NoError(t, err)
	require.NoError(t, nuodbaas_client_test.CreateProject(t, ctx, client, proj1Org, proj1Name, "dev", "n0.nano"))

	// Create project with label
	model := openapi.ProjectModel{
		Sla:  "dev",
		Tier: "n0.nano",
		Labels: &map[string]string{
			"key": "value",
		},
	}
	require.NoError(t, nuodbaas_client_test.CreateProjectWithModel(t, ctx, client, proj2Org, proj2Name, model))

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
						filter = {
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
			// Specify label filter key=value, which should return proj2
			{
				Config: providerConfig + fmt.Sprintf(`
					data "%s" "%s" {
						filter = {
							labels = ["key=value"]
						}
					}
				`, getProjectsDatasourceTypeName(), resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePath, "projects.#", "1"),
					resource.TestCheckResourceAttr(resourcePath, "projects.0.name", proj2Name),
					resource.TestCheckResourceAttr(resourcePath, "projects.0.organization", proj2Org),
				),
				// Skip if running end-to-end test because path rewrite used with Nginx is broken and does not preserve query parameters
				// TODO: Remove once path rewrite is fixed
				SkipFunc: func() (bool, error) { return IsE2eTest(), nil },
			},
			// Specify label filter key!=value at org1 scope, which should return proj1
			{
				Config: providerConfig + fmt.Sprintf(`
					data "%s" "%s" {
						filter = {
							organization = "%s"
							labels = ["key!=value"]
						}
					}
				`, getProjectsDatasourceTypeName(), resourceName, proj1Org),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePath, "projects.#", "1"),
					resource.TestCheckResourceAttr(resourcePath, "projects.0.name", proj1Name),
					resource.TestCheckResourceAttr(resourcePath, "projects.0.organization", proj1Org),
				),
				// Skip if running end-to-end test because path rewrite used with Nginx is broken and does not preserve query parameters
				// TODO: Remove once path rewrite is fixed
				SkipFunc: func() (bool, error) { return IsE2eTest(), nil },
			},
			// Specify label filter key=value at org1 scope, which should return nothing
			{
				Config: providerConfig + fmt.Sprintf(`
					data "%s" "%s" {
						filter = {
							organization = "%s"
							labels = ["key=value"]
						}
					}
				`, getProjectsDatasourceTypeName(), resourceName, proj1Org),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePath, "projects.#", "0"),
				),
				// Skip if running end-to-end test because path rewrite used with Nginx is broken and does not preserve query parameters
				// TODO: Remove once path rewrite is fixed
				SkipFunc: func() (bool, error) { return IsE2eTest(), nil },
			},
		},
	})
}

func getProjectsDatasourceTypeName() string {
	source := NewProjectsDataSource()

	ctx := context.TODO() // Not used

	request := datasource.MetadataRequest{ProviderTypeName: getProviderTypeName()}
	response := datasource.MetadataResponse{}
	source.Metadata(ctx, request, &response)
	return response.TypeName
}
