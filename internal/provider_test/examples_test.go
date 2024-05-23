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

	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider"
)

func CombineConfigs(t *testing.T, root string) string {
	var combinedConfigs string
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		require.NoError(t, err)
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".tf") {
			configContent, err := os.ReadFile(path)
			require.NoError(t, err)
			combinedConfigs += "\n" + string(configContent)
		}
		return nil
	})
	return combinedConfigs
}

func TestValidateExamples(t *testing.T) {
	// Create provider server that runs within test
	ctx := context.Background()
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Combine all example resource and data source configs
	projectRoot := GetProjectRoot(t)
	resourcesDir := filepath.Join(projectRoot, "examples", "resources")
	dataSourcesDir := filepath.Join(projectRoot, "examples", "data-sources")
	configContent := NewTfConfigBuilder().WithProviderConfig("nuodbaas", &NuoDbaasProviderModel{}).Build()
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

	// Validate config
	_, err = tf.Validate()
	require.NoError(t, err)
}
