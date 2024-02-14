//go:build tools

/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package tools

// List of build tools to fetch using `make install-tools`
import (
	_ "github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
	_ "gotest.tools/gotestsum"
)
