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

var testOrgName = os.Getenv("NUODB_CP_ORGANIZATION")

func testAccPreCheck(t *testing.T) {
	// TODO: Verify that DBaaS is up
	// Also, clean up left over projects if there are any from previous runs
}
