//go:build tools

// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package tools

// List of build tools to fetch using `make install-tools`
import (
	_ "github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
	_ "gotest.tools/gotestsum"
)
