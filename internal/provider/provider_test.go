/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	nuodbaas_client_test "github.com/nuodb/terraform-provider-nuodbaas/internal/client/testclient"
)

const (
	providerConfig = `
		provider "nuodbaas" { }
  `
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	getProviderTypeName(): providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// TODO: Verify that DBaaS is up
	// Also, clean up left over projects if there are any from previous runs
}

func checkClean(_ *terraform.State) error {
	return nuodbaas_client_test.NewTestClient(context.TODO()).CheckClean()
}

func getProviderTypeName() string {
	instance := NuoDbaasProvider{}

	ctx := context.TODO() // Never actually used
	response := provider.MetadataResponse{}
	instance.Metadata(ctx, provider.MetadataRequest{}, &response)

	return response.TypeName
}
