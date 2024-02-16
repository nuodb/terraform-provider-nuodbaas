/* (C) Copyright 2016-2024 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"regexp"
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
						labels       = { "key": "value" }
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO: Test that the resources match what is in the REST service
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "organization", "org"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "name", "proj"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "sla", "dev"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "tier", "n0.nano"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "labels.%", "1"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "labels.key", "value"),
					// product_version should be computed by the REST service
					resource.TestCheckResourceAttrSet("nuodbaas_project.proj", "properties.product_version"),
					// The create operation should block until the state is Available
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "status.state", "Available"),
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
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "labels.%", "1"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "labels.key", "value"),
					resource.TestCheckResourceAttrSet("nuodbaas_project.proj", "properties.product_version"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "status.state", "Available"),
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
				// Update the project by setting is_disabled and product_version
				Config: providerConfig + `
				resource "nuodbaas_project" "proj" {
					organization = "org"
					name         = "proj"
					sla          = "dev"
					tier         = "n0.nano"
					maintenance  = {
						is_disabled = true
					}
					properties   = {
						product_version = "9.9.9"
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO: Test that the resources match what is in the REST service
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "organization", "org"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "name", "proj"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "sla", "dev"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "tier", "n0.nano"),
					// Labels should remain due to them being computed / UseStateForUnknown
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "labels.%", "1"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "labels.key", "value"),
					// product_version should be computed by the REST service
					resource.TestCheckResourceAttrSet("nuodbaas_project.proj", "properties.product_version"),
					// Check updated attribute values
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "maintenance.is_disabled", "true"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "properties.product_version", "9.9.9"),
					// The update operation should block until the state is Stopped when is_disabled=true
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "status.state", "Stopped"),
				),
			},
			{
				// Update the project by removing labels
				Config: providerConfig + `
				resource "nuodbaas_project" "proj" {
					organization = "org"
					name         = "proj"
					sla          = "dev"
					tier         = "n0.nano"
					labels       = {}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO: Test that the resources match what is in the REST service
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "organization", "org"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "name", "proj"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "sla", "dev"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "tier", "n0.nano"),
					// Labels should be removed due to explicit setting of labels={}
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "labels.%", "0"),
					// Maintenance setting and properties should remain due to them being computed / UseStateForUnknown
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "maintenance.is_disabled", "true"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "properties.product_version", "9.9.9"),
					resource.TestCheckResourceAttr("nuodbaas_project.proj", "status.state", "Stopped"),
				),
			},
			{
				// Negative test: specify an invalid product_version
				Config: providerConfig + `
				resource "nuodbaas_project" "proj" {
					organization = "org"
					name         = "proj"
					sla          = "dev"
					tier         = "n0.nano"
					properties   = {
						product_version = "x.y.z"
					}
				}
				`,
				// The testing library annoyingly does not allow checking of the detail message, only the summary
				ExpectError: regexp.MustCompile("Error updating project"),
			},
		},
		CheckDestroy: checkClean,
	})
}
