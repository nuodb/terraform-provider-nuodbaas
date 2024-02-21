package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
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

func (b *TfConfigBuilder) WithDatabaseResource(name string, database *DatabaseResourceModel, dependsOn ...string) *TfConfigBuilder {
	key := "nuodbaas_database " + name
	b.resources[key] = database
	if len(dependsOn) != 0 {
		b.resourcesDependsOn[key] = []string(dependsOn)
	}
	return b
}

func (b *TfConfigBuilder) WithProjectResource(name string, project *ProjectResourceModel, dependsOn ...string) *TfConfigBuilder {
	key := "nuodbaas_project " + name
	b.resources[key] = project
	if len(dependsOn) != 0 {
		b.resourcesDependsOn[key] = []string(dependsOn)
	}
	return b
}

func (b *TfConfigBuilder) WithDatabaseDataSource(name string, database *DatabaseNameModel, dependsOn ...string) *TfConfigBuilder {
	key := "nuodbaas_database " + name
	b.dataSources[key] = database
	if len(dependsOn) != 0 {
		b.dataSourcesDependsOn[key] = []string(dependsOn)
	}
	return b
}

func (b *TfConfigBuilder) WithProjectDataSource(name string, project *ProjectNameModel, dependsOn ...string) *TfConfigBuilder {
	key := "nuodbaas_project " + name
	b.dataSources[key] = project
	if len(dependsOn) != 0 {
		b.dataSourcesDependsOn[key] = []string(dependsOn)
	}
	return b
}

func (b *TfConfigBuilder) Build() string {
	f := hclwrite.NewEmptyFile()
	for key, value := range b.providers {
		block := f.Body().AppendNewBlock("provider", []string{key}).Body()
		gohcl.EncodeIntoBody(value, block)
		f.Body().AppendNewline()
	}
	for key, value := range b.resources {
		block := f.Body().AppendNewBlock("resource", strings.Split(key, " ")).Body()
		gohcl.EncodeIntoBody(value, block)
		// Add depends_on attribute to block
		if dependsOn, ok := b.resourcesDependsOn[key]; ok {
			block.SetAttributeRaw("depends_on", tokensForIdentifierList(dependsOn))
		}
		f.Body().AppendNewline()
	}
	for key, value := range b.dataSources {
		block := f.Body().AppendNewBlock("data", strings.Split(key, " ")).Body()
		gohcl.EncodeIntoBody(value, block)
		// Add depends_on attribute to block
		if dependsOn, ok := b.dataSourcesDependsOn[key]; ok {
			block.SetAttributeRaw("depends_on", tokensForIdentifierList(dependsOn))
		}
		f.Body().AppendNewline()
	}
	return string(f.Bytes())
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

func (tf *TfHelper) DestroySilently() {
	copy := *tf
	copy.Silent = true
	_, _ = copy.Destroy()
}

func (tf *TfHelper) WriteConfig(tfConfig string) error {
	tfConfigFile := filepath.Join(tf.WorkingDir, "main.tf")
	tfConfig = REQUIRED_PROVIDERS + "\n\n" + tfConfig
	fmt.Println()
	fmt.Printf("> cat <<EOF > %s\n%s\nEOF\n", tfConfigFile, tfConfig)
	return os.WriteFile(tfConfigFile, []byte(tfConfig), 0644)
}

func (tf *TfHelper) WriteConfigT(t *testing.T, tfConfig string) {
	err := tf.WriteConfig(tfConfig)
	require.NoError(t, err)
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