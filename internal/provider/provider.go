// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package provider

import (
	"context"
	"net/url"
	"os"
	"strings"
	"time"

	nuodbaas_client "github.com/nuodb/terraform-provider-nuodbaas/internal/client"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/backup"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/backuppolicy"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/database"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/project"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// Ensure NuoDbaasProvider satisfies various provider interfaces.
var (
	_ provider.Provider                   = &NuoDbaasProvider{}
	_ provider.ProviderWithValidateConfig = &NuoDbaasProvider{}
)

// NuoDbaasProvider defines the provider implementation.
type NuoDbaasProvider struct {
	version string
}

// NuoDbaasProviderModel describes the provider data model.
type NuoDbaasProviderModel struct {
	User       *string                                `tfsdk:"user" hcl:"user" cty:"user"`
	Password   *string                                `tfsdk:"password" hcl:"password" cty:"password"`
	Token      *string                                `tfsdk:"token" hcl:"token" cty:"token"`
	UrlBase    *string                                `tfsdk:"url_base" hcl:"url_base" cty:"url_base"`
	SkipVerify *bool                                  `tfsdk:"skip_verify" hcl:"skip_verify" cty:"skip_verify"`
	Timeouts   map[string]framework.OperationTimeouts `tfsdk:"timeouts" hcl:"timeouts" cty:"timeouts"`
}

var _ framework.ProviderConfig = &NuoDbaasProviderModel{}

const (
	NUODB_CP_USER        = "NUODB_CP_USER"
	NUODB_CP_PASSWORD    = "NUODB_CP_PASSWORD" //nolint:gosec // This is not a hardcoded password
	NUODB_CP_TOKEN       = "NUODB_CP_TOKEN"    //nolint:gosec // This is not a hardcoded authentication token
	NUODB_CP_URL_BASE    = "NUODB_CP_URL_BASE"
	NUODB_CP_SKIP_VERIFY = "NUODB_CP_SKIP_VERIFY"
)

func (pm *NuoDbaasProviderModel) GetUser() string {
	if pm.User != nil {
		return *pm.User
	}
	return os.Getenv(NUODB_CP_USER)
}

func (pm *NuoDbaasProviderModel) GetPassword() string {
	if pm.Password != nil {
		return *pm.Password
	}
	return os.Getenv(NUODB_CP_PASSWORD)
}

func (pm *NuoDbaasProviderModel) GetToken() string {
	if pm.Token != nil {
		return *pm.Token
	}
	return os.Getenv(NUODB_CP_TOKEN)
}

func (pm *NuoDbaasProviderModel) GetUrlBase() string {
	if pm.UrlBase != nil {
		return *pm.UrlBase
	}
	return os.Getenv(NUODB_CP_URL_BASE)
}

func (pm *NuoDbaasProviderModel) GetSkipVerify() bool {
	if pm.SkipVerify != nil {
		return *pm.SkipVerify
	}
	return os.Getenv(NUODB_CP_SKIP_VERIFY) == "true"
}

func (pm *NuoDbaasProviderModel) CreateClient() (openapi.ClientInterface, error) {
	return nuodbaas_client.NewApiClient(pm.GetUrlBase(), pm.GetUser(), pm.GetPassword(), pm.GetToken(), pm.GetSkipVerify())
}

func (p *NuoDbaasProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nuodbaas"
	resp.Version = p.version
}

func (p *NuoDbaasProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The NuoDB DBaaS Provider manages NuoDB databases using the NuoDB Control Plane.",
		Attributes: map[string]schema.Attribute{
			"user": schema.StringAttribute{
				Description: "The name of the user in the format `<organization>/<user>`. " +
					"If not specified, defaults to the value of the `" + NUODB_CP_USER + "` environment variable.",
				Optional: true,
			},
			"password": schema.StringAttribute{
				Description: "The password for the user. " +
					"If not specified, defaults to the value of the `" + NUODB_CP_PASSWORD + "` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"token": schema.StringAttribute{
				Description: "The token to use to authenticate the user. " +
					"If not specified, defaults to the value of the `" + NUODB_CP_TOKEN + "` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"url_base": schema.StringAttribute{
				Description: "The base URL for the server, including the protocol. " +
					"If not specified, defaults to the value of the `" + NUODB_CP_URL_BASE + "` environment variable.",
				Optional: true,
			},
			"skip_verify": schema.BoolAttribute{
				Description: "Whether to skip server certificate verification. " +
					"If not specified, defaults to the value of the `" + NUODB_CP_SKIP_VERIFY + "` environment variable.",
				Optional: true,
			},
			"timeouts": schema.MapNestedAttribute{
				Description: "Timeouts by resource type and operation. A resource type of `default` is used to supply timeouts for all resources that are not specified explicitly.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						framework.CREATE_OPERATION: schema.StringAttribute{
							Description: "The timeout for resource readiness after creation, specified as a duration with time unit suffix, e.g. `10m`. " +
								"A timeout of `0` indicates not to wait.",
							Optional: true,
						},
						framework.UPDATE_OPERATION: schema.StringAttribute{
							Description: "The timeout for resource readiness after update, specified as a duration with time unit suffix, e.g. `1m`. " +
								"A timeout of `0` indicates not to wait.",
							Optional: true,
						},
						framework.DELETE_OPERATION: schema.StringAttribute{
							Description: "The timeout for resource deletion, specified as a duration with time unit suffix, e.g. `30s`. " +
								"A timeout of `0` indicates not to wait.",
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func (p *NuoDbaasProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	config, timeouts := parseAndValidate(ctx, req.Config, &resp.Diagnostics)

	// Check that no errors occurred
	if resp.Diagnostics.HasError() {
		return
	}

	// Create client
	client, err := config.CreateClient()
	if err != nil {
		resp.Diagnostics.AddError("Unable to create client", err.Error())
		return
	}

	// Pass client as opaque data
	providerClient := framework.NewProviderClient(&config, client, timeouts)
	resp.DataSourceData = providerClient
	resp.ResourceData = providerClient
}

func parseAndValidate(ctx context.Context, rawConfig tfsdk.Config, diags *diag.Diagnostics) (NuoDbaasProviderModel, map[string]map[string]time.Duration) {
	var config NuoDbaasProviderModel
	if !framework.ReadResource(ctx, diags, rawConfig.Get, &config) {
		return config, nil
	}

	// Validate server URL
	if config.GetUrlBase() == "" {
		diags.AddError("Invalid provider configuration", "Must specify url_base or the environment variable "+NUODB_CP_URL_BASE)
	} else {
		url, err := url.Parse(config.GetUrlBase())
		// url.Parse() does not return error if scheme is missing, so
		// check that explicitly
		if err != nil {
			diags.AddAttributeError(path.Empty().AtName("url_base"), "Invalid provider configuration", err.Error())
		} else if url.Scheme == "" {
			diags.AddAttributeError(path.Empty().AtName("url_base"), "Invalid provider configuration", "No scheme found in URL")
		}
	}

	// Validate timeout configuration
	timeouts, err := framework.ParseTimeouts(config.Timeouts, resourceTypes())
	if err != nil {
		diags.AddAttributeError(path.Empty().AtName("timeouts"), "Invalid provider configuration", err.Error())
	}

	// Validate credentials
	hasUser := config.GetUser() != ""
	hasPassword := config.GetPassword() != ""

	if (hasUser && !hasPassword) || (hasPassword && !hasUser) { // user xnor password
		diags.AddError("Partial credentials", "To use basic authentication, both user name and password should be provided.")
	}

	if hasUser {
		userParts := strings.Split(config.GetUser(), "/")
		if len(userParts) != 2 || len(userParts[0]) < 1 || len(userParts[1]) < 1 {
			diags.AddAttributeError(path.Root("user"), "Malformed user name", "User name should be in the format \"<organization>/<user>\".")
		}
		// Make sure that token was not also supplied
		if config.GetToken() != "" {
			diags.AddError("Multiple credentials", "Both basic and token authentication credentials were supplied.")
		}
	}

	return config, timeouts
}

func (p *NuoDbaasProvider) Resources(ctx context.Context) []func() resource.Resource {
	return resources()
}

// resourceTypes returns the set of available resource types by name.
func resourceTypes() map[string]struct{} {
	set := make(map[string]struct{})
	for _, resourceFn := range resources() {
		resource, ok := resourceFn().(*framework.GenericResource)
		if ok {
			set[resource.TypeName] = struct{}{}
		}
	}
	return set
}

func resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewDatabaseResource,
		NewBackupPolicyResource,
		NewBackupResource,
	}
}

func (p *NuoDbaasProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
		NewProjectsDataSource,
		NewDatabaseDataSource,
		NewDatabasesDataSource,
		NewBackupPolicyDataSource,
		NewBackupPoliciesDataSource,
		NewBackupDataSource,
		NewBackupsDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NuoDbaasProvider{
			version: version,
		}
	}
}

func (p *NuoDbaasProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	parseAndValidate(ctx, req.Config, &resp.Diagnostics)
}
