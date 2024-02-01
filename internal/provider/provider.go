/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/nuodb/terraform-provider-nuodbaas/helper"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	nuodbaas "github.com/nuodb/terraform-provider-nuodbaas/client"
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
	Organization types.String `tfsdk:"organization"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	BaseUrl      types.String `tfsdk:"url_base"`
	SkipVerify   types.Bool   `tfsdk:"skip_verify"`
}

func (p *NuoDbaasProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nuodbaas"
	resp.Version = p.version
}

func (p *NuoDbaasProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Description: "The name of the organization for the user",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "The name of the user",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for the user",
				Optional:    true,
				Sensitive:   true,
			},
			"url_base": schema.StringAttribute{
				Description: "The base URL for the server, including the protocol",
				Optional:    true,
			},
			"skip_verify": schema.BoolAttribute{
				Description: "Whether to skip server certificate verification",
				Optional:    true,
			},
		},
	}
}

func (p *NuoDbaasProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config NuoDbaasProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Why would any of these be unknown? Unknown is different from
	// unspecified. All of these fields are marked as optional, so is the
	// IsUnknown() check relevant for any of them?
	if config.Organization.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("organization"),
			"Unknown organization",
			"The provider cannot create the NuoDB DBaaS client as there is an unknown configuration value for the organization of the user. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NUODB_CP_ORGANIZATION environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown username",
			"The provider cannot create the NuoDB DBaaS client as there is an unknown configuration value for the user. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NUODB_CP_USER environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown password",
			"The provider cannot create the NuoDB DBaaS client as there is an unknown configuration value for the password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NUODB_CP_PASSWORD environment variable.",
		)
	}

	if config.BaseUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url_base"),
			"Unknown server URL",
			"The provider cannot create the NuoDB DBaaS client as there is an unknown configuration value for the server URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NUODB_CP_URL_BASE environment variable.",
		)
	}

	if config.SkipVerify.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("skip_verify"),
			"Unknown skip verify value",
			"Unknown value for skip_verify",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	organization := os.Getenv("NUODB_CP_ORGANIZATION")
	username := os.Getenv("NUODB_CP_USER")
	password := os.Getenv("NUODB_CP_PASSWORD")
	urlBase := os.Getenv("NUODB_CP_URL_BASE")
	skipVerify := false
	if skipVerifyValue := os.Getenv("NUODB_CP_SKIP_VERIFY"); skipVerifyValue == "1" || strings.ToLower(skipVerifyValue) == "true" {
		skipVerify = true
	}

	if !config.Organization.IsNull() {
		organization = config.Organization.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if !config.BaseUrl.IsNull() {
		urlBase = config.BaseUrl.ValueString()
	}

	if !config.SkipVerify.IsNull() {
		skipVerify = config.SkipVerify.ValueBool()
	}

	if organization == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("organization"),
			"Missing organization for user",
			helper.GetProviderValidatorErrorMessage("organization", "NUODB_CP_ORGANIZATION"),
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing username",
			helper.GetProviderValidatorErrorMessage("username", "NUODB_CP_USER"),
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing password",
			helper.GetProviderValidatorErrorMessage("password", "NUODB_CP_PASSWORD"),
		)
	}

	if urlBase == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url_base"),
			"Missing server URL",
			helper.GetProviderValidatorErrorMessage("url_base", "NUODB_CP_URL_BASE"),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	configuration := nuodbaas.NewConfiguration()
	// Disable server certificate verification if skip_verify=true
	if skipVerify {
		configuration.HTTPClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}
	serverConfig := nuodbaas.ServerConfigurations{{URL: urlBase, Description: "The base URL for the server, including the protocol"}}
	basicUsername := fmt.Sprintf("%s/%s", organization, username)
	configuration.DefaultHeader["Authorization"] = fmt.Sprintf("Basic %s", basicAuth(basicUsername, password))
	configuration.Servers = serverConfig
	apiClient := nuodbaas.NewAPIClient(configuration)
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	// TODO: This is issuing a health check and then checking if the request
	// was successful or 403 Forbidden, which means that the user was
	// authenticated but does not have access to 'GET /healthz'. Do we
	// actually need to check this eagerly?
	httpsRes, err := apiClient.HealthzAPI.GetHealth(ctx).Execute()

	// Check for error on client side
	serverErr := helper.GetErrorContentObj(err)
	if err != nil && serverErr == nil {
		resp.Diagnostics.AddError("Unable to connect to server", err.Error())
		return
	}
	// Check for error other than 403 Forbidden
	if serverErr != nil && httpsRes.StatusCode != http.StatusForbidden {
		msg := fmt.Sprintf("code=%s, status=%s, detail=%s", serverErr.GetCode(), serverErr.GetStatus(), serverErr.GetDetail())
		resp.Diagnostics.AddError("Error response from server", msg)
		return
	}

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *NuoDbaasProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewDatabaseResource,
	}
}

func (p *NuoDbaasProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
		NewProjectsDataSource,
		NewDatabasesDataSource,
		NewDatabaseDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NuoDbaasProvider{
			version: version,
		}
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
