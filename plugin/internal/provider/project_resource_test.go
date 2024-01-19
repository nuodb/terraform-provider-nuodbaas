/* (C) Copyright 2016-2024 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create a project
				Config: providerConfig + `
					resource "nuodbaas_project" "proj" {
						organization = var.org_name
						name         = "proj"
						sla          = "dev"
						tier         = "n0.nano"
					}
				`,
				ConfigVariables: config.Variables{"org_name": config.StringVariable(testOrgName)},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "organization", testOrgName),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "name", "proj"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "sla", "dev"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "tier", "n0.nano"),
					resource.TestCheckResourceAttrSet("nuodbaas_project.proj", "resource_version"),
				),
			},
			{
				// Test that we can read it back
				RefreshState: true,
			},
			{
				// Import it
				ConfigVariables:         config.Variables{"org_name": config.StringVariable(testOrgName)},
				ResourceName:            "nuodbaas_project.proj",
				ImportState:             true,
				ImportStateVerify:       true,
				SkipFunc:                func() (bool, error) { return true, nil }, //TODO: Import does not work
				ImportStateVerifyIgnore: []string{"resource_version"},
			},
			{
				// Update the project by setting it to be disabled
				Config: providerConfig + `
				resource "nuodbaas_project" "proj" {
					organization = var.org_name
					name         = "proj"
					sla          = "dev"
					tier         = "n0.nano"
					maintenance = {
						is_disabled = true
					}
				}
				`,
				ConfigVariables: config.Variables{"org_name": config.StringVariable(testOrgName)},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "maintenance.is_disabled", "true"),
				),
			},
		},
		CheckDestroy: func(s *terraform.State) error {
			//TODO: Verify that project was cleaned up
			return nil
		},
	})
}
