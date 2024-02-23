/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"os"
	"strings"

	nuodbaas_client "github.com/nuodb/terraform-provider-nuodbaas/internal/client"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
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
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// NuoDbaasProviderModel describes the provider data model.
type NuoDbaasProviderModel struct {
	User       *string                                `tfsdk:"user" hcl:"user"`
	Password   *string                                `tfsdk:"password" hcl:"password"`
	UrlBase    *string                                `tfsdk:"url_base" hcl:"url_base"`
	SkipVerify *bool                                  `tfsdk:"skip_verify" hcl:"skip_verify"`
	Timeouts   map[string]framework.OperationTimeouts `tfsdk:"timeouts" hcl:"timeouts"`
}

func (pm *NuoDbaasProviderModel) GetUser() string {
	if pm.User != nil {
		return *pm.User
	}
	return os.Getenv("NUODB_CP_USER")
}

func (pm *NuoDbaasProviderModel) GetPassword() string {
	if pm.Password != nil {
		return *pm.Password
	}
	return os.Getenv("NUODB_CP_PASSWORD")
}

func (pm *NuoDbaasProviderModel) GetUrlBase() string {
	if pm.UrlBase != nil {
		return *pm.UrlBase
	}
	if ret := os.Getenv("NUODB_CP_URL_BASE"); ret != "" {
		return ret
	}
	// TODO(asz6): Consider removing this hardcoded default, which is
	// unlikely to be useful in a non-test environment
	return "http://localhost:8080"
}

func (pm *NuoDbaasProviderModel) GetSkipVerify() bool {
	if pm.SkipVerify != nil {
		return *pm.SkipVerify
	}
	skipVerifyValue := os.Getenv("NUODB_CP_SKIP_VERIFY")
	return skipVerifyValue == "1" || strings.ToLower(skipVerifyValue) == "true"
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
		Description: "The NuoDB DBaaS provider provides the ability to manage the projects and databases running under the NuoDB Control Plane.",
		Attributes: map[string]schema.Attribute{
			"user": schema.StringAttribute{
				Description: "The name of the user in the format `<organization>/<user>`. " +
					"If not specified, defaults to the `NUODB_CP_USER` environment variable.",
				Optional: true,
			},
			"password": schema.StringAttribute{
				Description: "The password for the user. " +
					"If not specified, defaults to the `NUODB_CP_PASSWORD` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"url_base": schema.StringAttribute{
				Description: "The base URL for the server, including the protocol. " +
					"If not specified, defaults to the `NUODB_CP_URL_BASE` environment variable.",
				Optional: true,
			},
			"skip_verify": schema.BoolAttribute{
				Description: "Whether to skip server certificate verification",
				Optional:    true,
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

	// Validate timeout configuration
	timeouts, err := framework.ParseTimeouts(config.Timeouts, resourceTypes())
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Empty().AtName("timeouts"), "Invalid provider configuration", err.Error())
		return
	}

	// Create client
	client, err := config.CreateClient()
	if err != nil {
		resp.Diagnostics.AddError("Client error", err.Error())
		return
	}

	// Pass client as opaque data
	resp.DataSourceData = client
	resp.ResourceData = framework.NewClientWithTimeouts(client, timeouts)
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
