/* (C) Copyright 2016-2024 Dassault Systemes SE.
All Rights Reserved.
*/

package examples

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/provider"

	nuodbaas_client_test "github.com/nuodb/terraform-provider-nuodbaas/internal/client/testclient"
)

func testChart(t *testing.T, plan string, configVariables config.Variables, hasResources bool, noApply bool, checkClean bool) {
	if checkClean {
		nuodbaas_client_test.CheckClean()
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"nuodbaas": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config:             plan,
				PlanOnly:           noApply,
				ExpectNonEmptyPlan: noApply && hasResources,
				ConfigVariables:    configVariables,
			},
		},
		CheckDestroy: func(s *terraform.State) error {
			if checkClean {
				return nuodbaas_client_test.CheckClean()
			}
			return nil
		},
	})
}

func testChartDir(t *testing.T, path string, setUp string, configVariables config.Variables, hasResources bool, noApply bool, checkClean bool) {
	require.NoError(t, filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		require.NoError(t, err)
		if d.IsDir() {
			return nil
		}

		t.Run(path, func(t *testing.T) {
			plan, err := os.ReadFile(path)
			require.NoError(t, err)
			testChart(t, setUp+string(plan), configVariables, hasResources, noApply, checkClean)
		})

		return nil
	}))
}

// Should tests skip testing if charts apply successfully, in favor of faster execution.
func noApplyDefault(t *testing.T) bool {
	if testing.Short() {
		t.Log("Not trying to apply examples because -short flag set.")
		return true
	}
	return false
}

func TestExamplesResourceProject(t *testing.T) {
	path := "resources/nuodbaas_project"

	configVariables := config.Variables{}

	// Other plan parts that the example assumes exist
	setUp := ``

	testChartDir(t, path, setUp, configVariables, true, noApplyDefault(t), true)
}

func TestExamplesResourceDatabase(t *testing.T) {
	path := "resources/nuodbaas_database"

	configVariables := config.Variables{}

	// Other plan parts that the example assumes exist
	setUp := `resource "nuodbaas_project" "nuodb" {
		organization = "org"
		name         = "nuodb"
		sla          = "prod"
		tier         = "n0.nano"
	  }`

	testChartDir(t, path, setUp, configVariables, true, noApplyDefault(t), true)
}

func TestExamplesDatasources(t *testing.T) {
	path := "data-sources"

	configVariables := config.Variables{}

	// Other plan parts that the example assumes exist
	setUp := ``

	ctx := context.TODO()
	client, err := nuodbaas_client_test.DefaultApiClient()
	require.NoError(t, err)
	require.NoError(t, nuodbaas_client_test.CreateProject(t, ctx, client, "system", "nuodb", "dev", "n0.nano"))
	require.NoError(t, nuodbaas_client_test.CreateDatabase(t, ctx, client, "system", "nuodb", "dbaas", "pass"))

	testChartDir(t, path, setUp, configVariables, false, noApplyDefault(t), false)

	require.NoError(t, nuodbaas_client_test.DeleteDatabase(ctx, client, "system", "nuodb", "dbaas", false))
	require.NoError(t, nuodbaas_client_test.DeleteProject(ctx, client, "system", "nuodb", false))

	require.NoError(t, nuodbaas_client_test.CheckClean())
}

func TestExamplesProvider(t *testing.T) {
	path := "provider"

	url := os.Getenv("NUODB_CP_URL_BASE")
	if url == "" {
		url = nuodbaas_client_test.DEFAULT_URL
	}

	configVariables := config.Variables{
		"dbaas_credentials": config.MapVariable(map[string]config.Variable{
			"url_base": config.StringVariable(url),
			"user":     config.StringVariable(os.Getenv("NUODB_CP_USER")),
			"password": config.StringVariable(os.Getenv("NUODB_CP_PASSWORD")),
		}),
	}

	// Other plan parts that the example assumes exist
	setUp := `variable "dbaas_credentials" {
		type        = map
	}
	`

	require.NoError(t, filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		require.NoError(t, err)
		if d.IsDir() {
			return nil
		}

		t.Run(d.Name(), func(t *testing.T) {
			plan, err := os.ReadFile(path)
			require.NoError(t, err)

			// Remove our plugin from the config. We should use the provider factory.
			file, diag := hclwrite.ParseConfig(plan, path, hcl.Pos{Line: 1, Column: 1})
			require.Empty(t, diag.Errs())

			body := file.Body()
			for _, block := range body.Blocks() {
				if block.Type() == "terraform" {
					for _, subBlock := range block.Body().Blocks() {
						if subBlock.Type() == "required_providers" {
							subBlock.Body().RemoveAttribute("nuodbaas")
						}
					}
				}
			}

			testChart(t, setUp+string(file.Bytes()), configVariables, true, noApplyDefault(t), true)
		})

		return nil
	}))
}
