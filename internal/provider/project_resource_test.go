/* (C) Copyright 2016-2024 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
						organization = "org"
						name         = "proj"
						sla          = "dev"
						tier         = "n0.nano"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO: Test that the resources match what is in the REST service
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "organization", "org"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "name", "proj"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "sla", "dev"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "tier", "n0.nano"),
				),
			},
			{
				// Test that we can read it back
				RefreshState: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "organization", "org"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "name", "proj"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "sla", "dev"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "tier", "n0.nano"),
				),
			},
			{
				// Import it
				ResourceName:                         "nuodbaas_project.proj",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "org/proj",
				ImportStateVerifyIdentifierAttribute: "name",
			},
			{
				// Update the project by setting it to be disabled
				Config: providerConfig + `
				resource "nuodbaas_project" "proj" {
					organization = "org"
					name         = "proj"
					sla          = "dev"
					tier         = "n0.nano"
					maintenance = {
						is_disabled = true
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO: Test that the resources match what is in the REST service
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "maintenance.is_disabled", "true"),
				),
			},
		},
		CheckDestroy: checkClean,
	})
}
