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

func TestAccDatabaseResource(t *testing.T) {
	projConfig := providerConfig + `
	resource "nuodbaas_project" "proj" {
		organization = var.org_name
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
					organization = var.org_name
					project      = nuodbaas_project.proj.name
					name         = "db"
					dba_password = "changeMe"
				}
				`,
				ConfigVariables: config.Variables{"org_name": config.StringVariable(testOrgName)},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuodbaas_database.db", "organization", testOrgName),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "name", "db"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "project", "proj"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "dba_password", "changeMe"),
					resource.TestCheckResourceAttr("nuodbaas_database.db", "tier", "n0.nano"),
					resource.TestCheckResourceAttrSet("nuodbaas_database.db", "resource_version"),
				),
			},
			{
				// Test that we can read it back
				RefreshState: true,
			},
			{
				// Import it
				ConfigVariables:         config.Variables{"org_name": config.StringVariable(testOrgName)},
				ResourceName:            "nuodbaas_database.db",
				ImportState:             true,
				ImportStateVerify:       true,
				SkipFunc:                func() (bool, error) { return true, nil }, //TODO: Import does not work
				ImportStateVerifyIgnore: []string{"resource_version"},
			},
			{
				// Update the database by setting it to be disabled
				Config: projConfig + `
				resource "nuodbaas_database" "db" {
					organization = var.org_name
					project      = nuodbaas_project.proj.name
					name         = "db"
					dba_password = "changeMe"
					maintenance = {
						is_disabled = true
					}
				}
				`,
				ConfigVariables: config.Variables{"org_name": config.StringVariable(testOrgName)},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("nuodbaas_database.db", "maintenance.is_disabled", "true"),
				),
			},
		},
		CheckDestroy: func(s *terraform.State) error {
			//TODO: Verify that database was cleaned up
			return nil
		},
	})
}
