// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package provider

import (
	"context"
	"net/url"
	"os"

	nuodbaas_client "github.com/nuodb/terraform-provider-nuodbaas/internal/client"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/database"
	. "github.com/nuodb/terraform-provider-nuodbaas/internal/provider/project"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure NuoDbaasProvider satisfies various provider interfaces.
var (
	_ provider.Provider = &NuoDbaasProvider{}
)

// NuoDbaasProvider defines the provider implementation.
type NuoDbaasProvider struct {
	version string
}

// NuoDbaasProviderModel describes the provider data model.
type NuoDbaasProviderModel struct {
	User       *string                                `tfsdk:"user" hcl:"user" cty:"user"`
	Password   *string                                `tfsdk:"password" hcl:"password" cty:"password"`
	UrlBase    *string                                `tfsdk:"url_base" hcl:"url_base" cty:"url_base"`
	SkipVerify *bool                                  `tfsdk:"skip_verify" hcl:"skip_verify" cty:"skip_verify"`
	Timeouts   map[string]framework.OperationTimeouts `tfsdk:"timeouts" hcl:"timeouts" cty:"timeouts"`
}

const (
	NUODB_CP_USER        = "NUODB_CP_USER"
	NUODB_CP_PASSWORD    = "NUODB_CP_PASSWORD" //nolint:gosec // This is not a hardcoded password
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

func (pm *NuoDbaasProviderModel) CreateClient() (*openapi.Client, error) {
	return nuodbaas_client.NewApiClient(pm.GetUrlBase(), pm.GetUser(), pm.GetPassword(), pm.GetSkipVerify())
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
	var config NuoDbaasProviderModel
	if !framework.ReadResource(ctx, &resp.Diagnostics, req.Config.Get, &config) {
		return
	}

	// Validate server URL
	if config.GetUrlBase() == "" {
		resp.Diagnostics.AddError("Invalid provider configuration", "Must specify url_base or the environment variable "+NUODB_CP_URL_BASE)
	} else {
		url, err := url.Parse(config.GetUrlBase())
		// url.Parse() does not return error if scheme is missing, so
		// check that explicitly
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Empty().AtName("url_base"), "Invalid provider configuration", err.Error())
		} else if url.Scheme == "" {
			resp.Diagnostics.AddAttributeError(path.Empty().AtName("url_base"), "Invalid provider configuration", "No scheme found in URL")
		}
	}

	// Validate timeout configuration
	timeouts, err := framework.ParseTimeouts(config.Timeouts, resourceTypes())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Empty().AtName("timeouts"), "Invalid provider configuration", err.Error())
	}

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
	clientWithOptions := framework.NewClientWithOptions(client, timeouts)
	resp.DataSourceData = clientWithOptions
	resp.ResourceData = clientWithOptions
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
	}
}

func (p *NuoDbaasProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
		NewProjectsDataSource,
		NewDatabaseDataSource,
		NewDatabasesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NuoDbaasProvider{
			version: version,
		}
	}
}
