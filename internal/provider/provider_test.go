/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
		variable "org_name" {
			description = "The name of the organization for the user"
			type        = string
		}
		provider "nuodbaas" { }
  `
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"nuodbaas": providerserver.NewProtocol6WithError(New("test")()),
}

var testOrgName = getOrganization()

// TODO: The tests seem to rely on the organization name being resolved from the
// environment. Not sure why this is a requirement, but for now, just use the
// environment variable if it is set and otherwise return an arbitrary
// organization name.
func getOrganization() string {
	org := os.Getenv("NUODB_CP_ORGANIZATION")
	if org != "" {
		return org
	}
	return "org"
}

func testAccPreCheck(t *testing.T) {
	// TODO: Verify that DBaaS is up
	// Also, clean up left over projects if there are any from previous runs
}
