/* (C) Copyright 2016-2024 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseResource(t *testing.T) {
	projConfig := providerConfig + `
	resource "nuodbaas_project" "proj" {
		organization = "org"
		name         = "proj"
		sla          = "dev"
		tier         = "n0.nano"
	}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create project and database
				Config: projConfig + `
				resource "nuodbaas_database" "db" {
					organization = nuodbaas_project.proj.organization
					project      = nuodbaas_project.proj.name
					name         = "db"
					dba_password = "changeMe"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO: Test that the resources match what is in the REST service
					resource.TestCheckResourceAttr("nuodbaas_database.db", "organization", "org"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "name", "db"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "project", "proj"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "dba_password", "changeMe"),
					// Tier is inherited from project
					resource.TestCheckResourceAttr("nuodbaas_database.db", "tier", "n0.nano"),
				),
			},
			{
				// Test that we can read it back
				RefreshState: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO: Test that the resources match what is in the REST service
					resource.TestCheckResourceAttr("nuodbaas_database.db", "organization", "org"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "name", "db"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "project", "proj"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "dba_password", "changeMe"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "tier", "n0.nano"),
				),
			},
			{
				// Import it
				ResourceName:                         "nuodbaas_database.db",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "org/proj/db",
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"dba_password"},
			},
			{
				// Update the database by adding labels, setting it to be disabled, and setting product version
				Config: projConfig + `
				resource "nuodbaas_database" "db" {
					organization = nuodbaas_project.proj.organization
					project      = nuodbaas_project.proj.name
					name         = "db"
					dba_password = "changeMe"
					labels       = {
						x   = "y"
						foo = "bar"
					}
					maintenance  = {
						is_disabled = true
					}
					properties = {
						product_version = "9.9.9"
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO: Test that the resources match what is in the REST service
					resource.TestCheckResourceAttr("nuodbaas_database.db", "organization", "org"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "name", "db"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "project", "proj"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "dba_password", "changeMe"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "tier", "n0.nano"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "labels.%", "2"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "labels.x", "y"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "labels.foo", "bar"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "maintenance.is_disabled", "true"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "properties.product_version", "9.9.9"),
				),
			},
		},
		CheckDestroy: checkClean,
	})
}
