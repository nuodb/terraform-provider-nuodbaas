// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package provider_test

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider"
)

const (
	USING_LATEST_API TestOption = "USING_LATEST_API"
)

func CombineConfigs(t *testing.T, root string) string {
	var combinedConfigs string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		require.NoError(t, err)
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".tf") {
			configContent, err := os.ReadFile(path)
			require.NoError(t, err)
			combinedConfigs += "\n" + string(configContent)
		}
		return nil
	})
	require.NoError(t, err)
	return combinedConfigs
}

func TestExamples(t *testing.T) {
	// Create provider server that runs within test
	ctx := context.Background()
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Disable readiness check for backup resource so that backup controller does not have to be enabled
	var providerCfg NuoDbaasProviderModel
	providerCfg.Timeouts = map[string]OperationTimeouts{
		"backup": {
			Create: ptr("0"),
			Update: ptr("0"),
		},
	}

	// Combine all example resource and data source configs
	projectRoot := GetProjectRoot(t)
	resourcesDir := filepath.Join(projectRoot, "examples", "resources")
	dataSourcesDir := filepath.Join(projectRoot, "examples", "data-sources")
	configContent := NewTfConfigBuilder().WithProviderConfig("nuodbaas", &providerCfg).Build()
	configContent += CombineConfigs(t, resourcesDir)
	configContent += "\n" + CombineConfigs(t, dataSourcesDir)

	// Create Terraform workspace
	tf := CreateTerraformWorkspace(t)
	err := tf.SetReattachConfig(reattachCfg)
	require.NoError(t, err)

	// Initialize workspace with config
	tf.WriteConfigT(t, configContent)
	_, err = tf.Init()
	require.NoError(t, err)

	if USING_LATEST_API.IsTrue() {
		// Run `terraform apply` on config
		_, err = tf.Apply()
		defer tf.DestroySilently()
		require.NoError(t, err)

		// Run `terraform refresh` to update data sources
		_, err = tf.Run("refresh")
		require.NoError(t, err)
	} else {
		// The server may not support all resources and data sources, if
		// it is running at all. Run `terraform validate` only, which
		// does not require a server connection.
		_, err = tf.Validate()
		require.NoError(t, err)
	}
}
