// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package provider_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/database"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/project"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/rogpeppe/go-internal/diff"
	"github.com/stretchr/testify/require"
)

type TfConfigBuilder struct {
	providers            map[string]any
	resources            map[string]any
	dataSources          map[string]any
	resourcesDependsOn   map[string][]string
	dataSourcesDependsOn map[string][]string
}

func NewTfConfigBuilder() *TfConfigBuilder {
	return &TfConfigBuilder{
		providers:            make(map[string]any),
		resources:            make(map[string]any),
		dataSources:          make(map[string]any),
		resourcesDependsOn:   make(map[string][]string),
		dataSourcesDependsOn: make(map[string][]string),
	}
}

func (b *TfConfigBuilder) WithProviderConfig(name string, provider *NuoDbaasProviderModel) *TfConfigBuilder {
	b.providers["nuodbaas"] = provider
	return b
}

func (b *TfConfigBuilder) WithResource(key string, resource any, dependsOn ...string) *TfConfigBuilder {
	b.resources[key] = resource
	if len(dependsOn) != 0 {
		b.resourcesDependsOn[key] = []string(dependsOn)
	}
	return b
}

func (b *TfConfigBuilder) WithDataSource(key string, dataSource any, dependsOn ...string) *TfConfigBuilder {
	b.dataSources[key] = dataSource
	if len(dependsOn) != 0 {
		b.dataSourcesDependsOn[key] = []string(dependsOn)
	}
	return b
}

func (b *TfConfigBuilder) WithoutResource(key string) *TfConfigBuilder {
	delete(b.resources, key)
	delete(b.resourcesDependsOn, key)
	return b
}

func (b *TfConfigBuilder) WithoutDataSource(key string) *TfConfigBuilder {
	delete(b.dataSources, key)
	delete(b.dataSourcesDependsOn, key)
	return b
}

func (b *TfConfigBuilder) WithDatabaseResource(name string, database *DatabaseResourceModel, dependsOn ...string) *TfConfigBuilder {
	return b.WithResource("nuodbaas_database."+name, database, dependsOn...)
}

func (b *TfConfigBuilder) WithProjectResource(name string, project *ProjectResourceModel, dependsOn ...string) *TfConfigBuilder {
	return b.WithResource("nuodbaas_project."+name, project, dependsOn...)
}

func (b *TfConfigBuilder) WithDatabaseDataSource(name string, database *DatabaseNameModel, dependsOn ...string) *TfConfigBuilder {
	return b.WithDataSource("nuodbaas_database."+name, database, dependsOn...)
}

func (b *TfConfigBuilder) WithProjectDataSource(name string, project *ProjectNameModel, dependsOn ...string) *TfConfigBuilder {
	return b.WithDataSource("nuodbaas_project."+name, project, dependsOn...)
}

func (b *TfConfigBuilder) WithDatabasesDataSource(name string, databases *DatabasesDataSourceModel, dependsOn ...string) *TfConfigBuilder {
	return b.WithDataSource("nuodbaas_databases."+name, databases, dependsOn...)
}

func (b *TfConfigBuilder) WithProjectsDataSource(name string, projects *ProjectsDataSourceModel, dependsOn ...string) *TfConfigBuilder {
	return b.WithDataSource("nuodbaas_projects."+name, projects, dependsOn...)
}

func (b *TfConfigBuilder) WithoutDatabaseResource(name string) *TfConfigBuilder {
	return b.WithoutResource("nuodbaas_database." + name)
}

func (b *TfConfigBuilder) WithoutProjectResource(name string) *TfConfigBuilder {
	return b.WithoutResource("nuodbaas_project." + name)
}

func (b *TfConfigBuilder) WithoutDatabaseDataSource(name string) *TfConfigBuilder {
	return b.WithoutDataSource("nuodbaas_database." + name)
}

func (b *TfConfigBuilder) WithoutProjectDataSource(name string) *TfConfigBuilder {
	return b.WithoutDataSource("nuodbaas_project." + name)
}

func (b *TfConfigBuilder) WithoutDatabasesDataSource(name string) *TfConfigBuilder {
	return b.WithoutDataSource("nuodbaas_databases." + name)
}

func (b *TfConfigBuilder) WithoutProjectsDataSource(name string) *TfConfigBuilder {
	return b.WithoutDataSource("nuodbaas_projects." + name)
}

func (b *TfConfigBuilder) Build() string {
	f := hclwrite.NewEmptyFile()
	ForEachInOrder(b.providers, func(key string, value any) {
		block := f.Body().AppendNewBlock("provider", []string{key}).Body()
		gohcl.EncodeIntoBody(value, block)
		f.Body().AppendNewline()
	})
	ForEachInOrder(b.resources, func(key string, value any) {
		block := f.Body().AppendNewBlock("resource", strings.Split(key, ".")).Body()
		gohcl.EncodeIntoBody(value, block)
		// Add depends_on attribute to block
		if dependsOn, ok := b.resourcesDependsOn[key]; ok {
			block.SetAttributeRaw("depends_on", tokensForIdentifierList(dependsOn))
		}
		f.Body().AppendNewline()
	})
	ForEachInOrder(b.dataSources, func(key string, value any) {
		block := f.Body().AppendNewBlock("data", strings.Split(key, ".")).Body()
		gohcl.EncodeIntoBody(value, block)
		// Add depends_on attribute to block
		if dependsOn, ok := b.dataSourcesDependsOn[key]; ok {
			block.SetAttributeRaw("depends_on", tokensForIdentifierList(dependsOn))
		}
		f.Body().AppendNewline()
	})
	return string(f.Bytes())
}

// ForEachInOrder iterates over entries in map in order and applies supplied function.
func ForEachInOrder(m map[string]any, fn func(string, any)) {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fn(key, m[key])
	}
}

func tokensForIdentifierList(args []string) hclwrite.Tokens {
	var tokens hclwrite.Tokens
	tokens = append(tokens, &hclwrite.Token{
		Type:  hclsyntax.TokenOBrack,
		Bytes: []byte("["),
	})
	for i, arg := range args {
		if i != 0 {
			tokens = append(tokens, &hclwrite.Token{
				Type:  hclsyntax.TokenComma,
				Bytes: []byte(","),
			})
		}
		tokens = append(tokens, hclwrite.TokensForIdentifier(arg)[0])
	}
	tokens = append(tokens, &hclwrite.Token{
		Type:  hclsyntax.TokenCBrack,
		Bytes: []byte("]"),
	})
	return tokens
}

// GetProjectRoot returns the root directory for the local clone of the
// terraform-provider-nuodbaas Git repository.
func GetProjectRoot(t *testing.T) string {
	_, filename, _, _ := runtime.Caller(0)
	t.Logf("Finding project root directory from %s", filename)
	dir := filepath.Dir(filename)
	for dir != "/" {
		fileInfo, err := os.Stat(filepath.Join(dir, ".git"))
		if err == nil && fileInfo.IsDir() {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	require.FailNow(t, "Unable to find root directory of project")
	return ""
}

const (
	TFRC_FORMAT = `provider_installation {
    filesystem_mirror {
        path    = "%s/dist/pkg_mirror"
        include = ["registry.terraform.io/nuodb/nuodbaas"]
    }
    direct {
        exclude = ["registry.terraform.io/nuodb/nuodbaas"]
    }
}`
	REQUIRED_PROVIDERS = `terraform {
  required_providers {
    nuodbaas = {
      source  = "registry.terraform.io/nuodb/nuodbaas"
    }
  }
}`
)

type TfHelper struct {
	WorkingDir       string
	TfrcFile         string
	Silent           bool
	ReattachProvider string
}

// SetReattachConfig configures the Terraform helper to use a local Terraform
// provider server instead of one packaged and installed into the filesystem
// mirror. This is useful because it allows code coverage data to be obtained
// using `go test` (instrumentation of the provider binary does not seem to work
// when Terraform is invoking it).
func (tf *TfHelper) SetReattachConfig(config plugin.ReattachConfig) error {
	encoded, err := json.Marshal(map[string]any{"registry.terraform.io/nuodb/nuodbaas": config})
	if err != nil {
		return err
	}
	tf.ReattachProvider = string(encoded)
	return nil
}

func (tf *TfHelper) Run(args ...string) ([]byte, error) {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = tf.WorkingDir
	if tf.ReattachProvider != "" {
		cmd.Env = append(os.Environ(), "TF_REATTACH_PROVIDERS="+tf.ReattachProvider)
	} else {
		cmd.Env = append(os.Environ(), "TF_CLI_CONFIG_FILE="+tf.TfrcFile)
	}
	out, err := cmd.CombinedOutput()
	if !tf.Silent {
		switch err.(type) {
		case nil, *exec.ExitError:
			fmt.Println()
			fmt.Printf("> terraform %s\n", strings.Join(args, " "))
			fmt.Printf("%s\n", out)
		}
	}
	return out, err
}

func (tf *TfHelper) Init() ([]byte, error) {
	return tf.Run("init")
}

func (tf *TfHelper) Plan() ([]byte, error) {
	return tf.Run("plan")
}

func (tf *TfHelper) Apply() ([]byte, error) {
	return tf.Run("apply", "-auto-approve")
}

func (tf *TfHelper) Destroy() ([]byte, error) {
	return tf.Run("destroy", "-auto-approve")
}

func (tf *TfHelper) Show() ([]byte, error) {
	return tf.Run("show")
}

func (tf *TfHelper) ShowJson() ([]byte, error) {
	return tf.Run("show", "-json")
}

func (tf *TfHelper) GetStateResources() ([]any, error) {
	copy := *tf
	copy.Silent = true
	out, err := copy.ShowJson()
	if err != nil {
		return nil, err
	}
	path := "values.root_module.resources"
	node, err := GetField(out, path)
	if err != nil {
		return nil, err
	}
	list, ok := node.([]any)
	if !ok {
		return nil, fmt.Errorf("List not found at field path %s", path)
	}
	return list, nil
}

func (tf *TfHelper) GetStateResource(address string) (any, error) {
	resources, err := tf.GetStateResources()
	if err != nil {
		return nil, err
	}
	for _, resource := range resources {
		v, err := FindChildNode(resource, "address")
		if err != nil {
			return nil, err
		}
		if addr, ok := v.(string); ok && addr == address {
			return FindChildNode(resource, "values")
		}
	}
	return nil, nil
}

func (tf *TfHelper) CheckStateResource(t *testing.T, address string) *AttributeChecker {
	resource, err := tf.GetStateResource(address)
	require.NoError(t, err)
	require.NotNil(t, resource)
	return &AttributeChecker{t, resource}
}

func (tf *TfHelper) DestroySilently() {
	copy := *tf
	copy.Silent = true
	_, _ = copy.Destroy()
}

func (tf *TfHelper) WriteConfig(tfConfig string) error {
	filename := "main.tf"
	tfConfigFile := filepath.Join(tf.WorkingDir, filename)
	tfConfig = REQUIRED_PROVIDERS + "\n\n" + tfConfig
	var orig []byte
	if _, err := os.Stat(tfConfigFile); err == nil {
		orig, err = os.ReadFile(tfConfigFile)
		if err != nil {
			orig = nil
		}
	}
	fmt.Println()
	// If there is an existing config, create patch from it to new config
	// and display it in output
	if orig != nil {
		patch := diff.Diff(filename, orig, filename, []byte(tfConfig))
		fmt.Printf("> patch -p0 <<EOF\n%sEOF\n", patch)
	} else {
		// Otherwise, just display the new config
		fmt.Printf("> cat <<EOF > %s\n%sEOF\n", filename, tfConfig)
	}
	return os.WriteFile(tfConfigFile, []byte(tfConfig), 0644)
}

func (tf *TfHelper) WriteConfigT(t *testing.T, tfConfig string) {
	err := tf.WriteConfig(tfConfig)
	require.NoError(t, err)
}

func GetField(jsonData []byte, path string) (any, error) {
	// Deserialize JSON to opaque map
	dest := make(map[string]any)
	err := json.Unmarshal(jsonData, &dest)
	if err != nil {
		return nil, err
	}
	return FindChildNode(dest, path)
}

func FindChildNode(node any, path string) (any, error) {
	// Traverse field path
	var ret any
	for _, field := range strings.Split(path, ".") {
		// Check that current node is an object
		object, ok := node.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("Invalid field path: %s", path)
		}
		// Get field within current object
		ret, ok = object[field]
		if !ok {
			return nil, nil
		}
		// Update node for next iteration
		node = ret
	}
	return ret, nil
}

type AttributeChecker struct {
	t        *testing.T
	resource any
}

func (ac *AttributeChecker) HasAttributeValue(attributePath string, expected any) *AttributeChecker {
	actual, err := FindChildNode(ac.resource, attributePath)
	require.NoError(ac.t, err)
	require.Equalf(ac.t, expected, actual, "Unexpected value for attribute %s", attributePath)
	return ac
}

func (ac *AttributeChecker) HasAttribute(attributePath string) *AttributeChecker {
	actual, err := FindChildNode(ac.resource, attributePath)
	require.NoError(ac.t, err)
	require.NotNil(ac.t, actual, "No attribute %s", attributePath)
	return ac
}

func (ac *AttributeChecker) DoesNotHaveAttributeValue(attributePath string, unexpected any) *AttributeChecker {
	actual, err := FindChildNode(ac.resource, attributePath)
	require.NoError(ac.t, err)
	require.NotEqualf(ac.t, unexpected, actual, "Unexpected value for attribute %s", attributePath)
	return ac
}

func (ac *AttributeChecker) DoesNotHaveAttribute(attributePath string) *AttributeChecker {
	actual, err := FindChildNode(ac.resource, attributePath)
	require.NoError(ac.t, err)
	require.Nil(ac.t, actual, "Unexpected attribute %s", attributePath)
	return ac
}

func (ac *AttributeChecker) ForEach(attributePath string, expectedCount int, assertFn func(*AttributeChecker)) *AttributeChecker {
	actual, err := FindChildNode(ac.resource, attributePath)
	require.NoError(ac.t, err)
	require.NotNil(ac.t, actual)

	// Check that value is a list or slice and has required number of elements
	v := reflect.ValueOf(actual)
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		require.Len(ac.t, actual, expectedCount)
	default:
		require.FailNow(ac.t, "Unexpected type: %T", actual)
	}

	// Iterate over elements and apply assert function
	for i := 0; i != v.Len(); i += 1 {
		elem := v.Index(i)
		assertFn(&AttributeChecker{ac.t, elem.Interface()})
	}
	return ac
}

// CreateTerraformWorkspace creates an empty directory to serve as a workspace for Terraform.
func CreateTerraformWorkspace(t *testing.T) *TfHelper {
	projectRoot := GetProjectRoot(t)
	workspaceDir := filepath.Join(projectRoot, "tmp", "tfworkspace")
	t.Logf("Removing workspace directory %s if it exists", workspaceDir)
	err := os.RemoveAll(workspaceDir)
	require.NoError(t, err)

	t.Logf("Creating workspace directory %s", workspaceDir)
	err = os.MkdirAll(workspaceDir, os.ModePerm)
	require.NoError(t, err)

	tfrcFile := filepath.Join(workspaceDir, "terraform.rc")
	config := fmt.Sprintf(TFRC_FORMAT, projectRoot)
	t.Logf("Writing registry configuration to %s:\n%s", tfrcFile, config)
	err = os.WriteFile(tfrcFile, []byte(config), 0644)
	require.NoError(t, err)

	return &TfHelper{
		WorkingDir: workspaceDir,
		TfrcFile:   tfrcFile,
	}
}

// CreateProviderServer creates a server that runs the nuodbaas Terraform provider.
func CreateProviderServer(t *testing.T, ctx context.Context) (plugin.ReattachConfig, func()) {
	ctx, cancel := context.WithCancel(ctx)
	config, closeCh, err := plugin.DebugServe(ctx, &plugin.ServeOpts{
		GRPCProviderV6Func: func() tfprotov6.ProviderServer {
			return providerserver.NewProtocol6(&NuoDbaasProvider{})()
		},
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:   "plugintest",
			Level:  hclog.Trace,
			Output: io.Discard,
		}),
		NoLogOutputOverride: true,
		UseTFLogSink:        t,
		ProviderAddr:        "registry.terraform.io/nuodb/nuodbaas",
	})
	require.NoError(t, err)
	// Cancel context and wait for channel to signal that the provider
	// server has been closed
	closeFn := func() {
		cancel()
		<-closeCh
	}
	return config, closeFn
}

func TestTfConfigBuilder(t *testing.T) {
	timeout := "5s"
	tfConfig := NewTfConfigBuilder().WithProviderConfig("nuodbaas", &NuoDbaasProviderModel{
		Timeouts: map[string]framework.OperationTimeouts{
			"default": {
				Create: &timeout,
				Update: &timeout,
				Delete: &timeout,
			},
		},
	}).WithProjectResource("proj", &ProjectResourceModel{
		Organization: "org",
		Name:         "proj",
	}).WithDatabaseResource("db", &DatabaseResourceModel{
		Organization: "org",
		Project:      "proj",
		Name:         "db",
	}, "nuodbaas_project.proj").WithProjectDataSource("proj", &ProjectNameModel{
		Organization: "org",
		Name:         "proj",
	}, "nuodbaas_project.proj").WithDatabaseDataSource("db", &DatabaseNameModel{
		Organization: "org",
		Project:      "proj",
		Name:         "db",
	}, "nuodbaas_database.db").Build()

	// Check that provider appears in config
	require.Contains(t, tfConfig, `provider "nuodbaas" {
  timeouts = {
    default = {
      create = "5s"
      delete = "5s"
      update = "5s"
    }
  }
}`)
	// Check that project resource appears in config
	require.Contains(t, tfConfig, `resource "nuodbaas_project" "proj" {
  organization = "org"
  name         = "proj"
  sla          = ""
  tier         = ""
}`)
	// Check that database resource appears in config
	require.Contains(t, tfConfig, `resource "nuodbaas_database" "db" {
  organization = "org"
  project      = "proj"
  name         = "db"
  depends_on   = [nuodbaas_project.proj]
`)
	// Check that project data source appears in config
	require.Contains(t, tfConfig, `data "nuodbaas_project" "proj" {
  organization = "org"
  name         = "proj"
  depends_on   = [nuodbaas_project.proj]
}`)
	// Check that database data source appears in config
	require.Contains(t, tfConfig, `data "nuodbaas_database" "db" {
  organization = "org"
  project      = "proj"
  name         = "db"
  depends_on   = [nuodbaas_database.db]
`)

	// Create provider server that runs within test
	ctx := context.Background()
	reattachCfg, closeFn := CreateProviderServer(t, ctx)
	defer closeFn()

	// Create Terraform workspace and initialize it with config
	tf := CreateTerraformWorkspace(t)
	tf.SetReattachConfig(reattachCfg)
	tf.WriteConfigT(t, tfConfig)
	_, err := tf.Init()
	require.NoError(t, err)
}

const (
	MOCK_OPERATOR_POLICY_CM = "mock-operator-policy"
	PATCH_FMT               = `[{"op": "add", "path": "/data", "value": {"markAsReady": "%s", "readinessDelaySeconds": "%s"}}]`
)

type MockReconcilePolicy struct {
	MarkAsReady           string `json:"markAsReady"`
	ReadinessDelaySeconds string `json:"readinessDelaySeconds"`
}

// GetMockReconcilePolicy returns the current policy being used by the mock reconcilier for domain and database resources.
func GetMockReconcilePolicy(t *testing.T) *MockReconcilePolicy {
	cmd := exec.Command("kubectl", "get", "configmap", MOCK_OPERATOR_POLICY_CM, "-o", "jsonpath={.data}", "--ignore-not-found")
	out, err := cmd.Output()
	require.NoError(t, err)

	// If configmap does not exist, return nil
	if strings.TrimSpace(string(out)) == "" {
		return nil
	}

	// Deserialize data fields
	var ret MockReconcilePolicy
	err = json.Unmarshal(out, &ret)
	require.NoError(t, err)

	return &ret
}

// SetMockReconcilePolicy updates the policy used by the mock reconcilier for
// domain and database custom resources and returns a function that reverts the
// configuration when invoked.
func SetMockReconcilePolicy(t *testing.T, newPolicy MockReconcilePolicy) func() {
	currentPolicy := GetMockReconcilePolicy(t)
	if currentPolicy == nil {
		t.Skipf("Configmap %s does not exist", MOCK_OPERATOR_POLICY_CM)
	}

	// Configmap exists, so patch it to have the supplied values
	patch := fmt.Sprintf(PATCH_FMT, newPolicy.MarkAsReady, newPolicy.ReadinessDelaySeconds)
	cmd := exec.Command("kubectl", "patch", "configmap", MOCK_OPERATOR_POLICY_CM, "--type=json", "-p", patch)
	_, err := cmd.Output()
	require.NoError(t, err)
	return func() { SetMockReconcilePolicy(t, *currentPolicy) }
}

func TestMockReconcilePolicy(t *testing.T) {
	// Save mock reconciliation policy
	policy := GetMockReconcilePolicy(t)
	if policy == nil {
		t.Skipf("Configmap %s does not exist", MOCK_OPERATOR_POLICY_CM)
	}

	t.Run("updateAndRevertPolicy", func(t *testing.T) {
		newPolicy := MockReconcilePolicy{MarkAsReady: "false", ReadinessDelaySeconds: "999"}
		reset := SetMockReconcilePolicy(t, newPolicy)
		defer reset()

		// Make sure reconcile policy is set to the supplied value
		require.Equal(t, &newPolicy, GetMockReconcilePolicy(t))
	})

	// Make sure reconcile policy is reverted to the original value
	require.Equal(t, policy, GetMockReconcilePolicy(t))
}
