/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package main

import (
	"context"
	"flag"
	"log"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Generate model structs and client for REST API:
//go:generate bin/oapi-codegen -generate types -include-tags databases,projects -package openapi -o openapi/types.go openapi.yaml
//go:generate bin/oapi-codegen -generate client -include-tags databases,projects -package openapi -o openapi/client.go openapi.yaml
//go:generate bin/oapi-codegen -generate spec -include-tags databases,projects -package openapi -o openapi/spec.go openapi.yaml

// Format Terraform examples:
//go:generate terraform fmt -recursive ./examples/

// Generate documentation:
//go:generate bin/tfplugindocs generate --provider-name nuodbaas

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "0.1.0"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/nuodb/nuodbaas",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
