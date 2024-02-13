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
				// Create a project
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
				ResourceName:      "nuodbaas_database.db",
				ImportState:       true,
				ImportStateVerify: true,
				SkipFunc:          func() (bool, error) { return true, nil }, //TODO: Import does not work
			},
			{
				// Update the database by setting it to be disabled
				Config: projConfig + `
				resource "nuodbaas_database" "db" {
					organization = nuodbaas_project.proj.organization
					project      = nuodbaas_project.proj.name
					name         = "db"
					dba_password = "changeMe"
					maintenance = {
						is_disabled = true
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO: Test that the resources match what is in the REST service
					resource.TestCheckResourceAttr("nuodbaas_database.db", "maintenance.is_disabled", "true"),
				),
			},
		},
		CheckDestroy: checkClean,
	})
}
