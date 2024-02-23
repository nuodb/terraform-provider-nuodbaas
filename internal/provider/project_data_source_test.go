/* (C) Copyright 2016-2024 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"fmt"
	"testing"

	nuodbaas_client_test "github.com/nuodb/terraform-provider-nuodbaas/internal/client/testclient"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/project"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAccProjectDataSource(t *testing.T) {
	var (
		organizationName = "org"
		projectName      = "proj"
		sla              = "dev"
		tier             = "n0.nano"
		disabled         = true
		productVersion   = "6.0"
		tierParamKey     = "foo"
		tierParamValue   = "bar"
	)

	model := openapi.ProjectModel{
		Sla:  sla,
		Tier: tier,
		Labels: &map[string]string{
			"key": "value",
		},
		Maintenance: &openapi.MaintenanceModel{
			IsDisabled: &disabled,
		},
		Properties: &openapi.ProjectPropertiesModel{
			ProductVersion: &productVersion,
			TierParameters: &map[string]string{tierParamKey: tierParamValue},
		},
	}

	ctx := context.TODO()

	client, err := nuodbaas_client_test.DefaultApiClient()
	require.NoError(t, err)
	require.NoError(t, nuodbaas_client_test.CreateProjectWithModel(t, ctx, client, organizationName, projectName, model))

	resourceName := "project_details"
	resourcePath := fmt.Sprintf("data.%s.%s", getProjectDatasourceTypeName(), resourceName)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
					data "%s" "%s" {
						organization = "%s"
						name      = "%s"
					}
				`, getProjectDatasourceTypeName(), resourceName, organizationName, projectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourcePath, "organization", organizationName),
					resource.TestCheckResourceAttr(resourcePath, "name", projectName),
					resource.TestCheckResourceAttr(resourcePath, "sla", sla),
					resource.TestCheckResourceAttr(resourcePath, "tier", tier),
					resource.TestCheckResourceAttr(resourcePath, "labels.%", "1"),
					resource.TestCheckResourceAttr(resourcePath, "labels.key", "value"),
					resource.TestCheckResourceAttr(resourcePath, "maintenance.is_disabled", "true"),
					resource.TestCheckResourceAttr(resourcePath, "properties.product_version", "6.0"),
					resource.TestCheckResourceAttr(resourcePath, "properties.tier_parameters.%", "1"),
					resource.TestCheckResourceAttr(resourcePath, "properties.tier_parameters."+tierParamKey, tierParamValue),
				),
			},
		},
	})
}

func getProjectDatasourceTypeName() string {
	source := NewProjectDataSource()

	ctx := context.TODO() // Not used

	request := datasource.MetadataRequest{ProviderTypeName: getProviderTypeName()}
	response := datasource.MetadataResponse{}
	source.Metadata(ctx, request, &response)
	return response.TypeName
}
