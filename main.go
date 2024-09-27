// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package main

import (
	"context"
	"flag"
	"log"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Generate model structs and client for REST API:
//go:generate bin/oapi-codegen -config oapi-codegen.yaml -generate types -o openapi/types.go openapi.yaml
//go:generate bin/oapi-codegen -config oapi-codegen.yaml -generate client -o openapi/client.go openapi.yaml
//go:generate bin/oapi-codegen -config oapi-codegen.yaml -generate spec -o openapi/spec.go openapi.yaml

// Inject current version into examples:
//go:generate ./inject-version.sh

// Format Terraform examples:
//go:generate bin/terraform fmt -recursive ./examples/

// Generate documentation:
//go:generate bin/tfplugindocs generate --provider-name nuodbaas

var (
	// The version from the Git tag is used by GoReleaser when publishing a
	// release, but the release job is conditionalized on the Git tag matching
	// the version in the code. `{{version}}` in the comment is a marker to
	// enable scraping of the version.
	//
	// For more information on GoReleaser, see https://goreleaser.com/cookbooks/using-main.version/
	version string = "1.3.0" // {{version}}
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
