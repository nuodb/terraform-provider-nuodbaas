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
		Description: "The NuoDB DBaaS provider provides the ability to manage the projects and databases running under the NuoDB Control Plane.",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Description: "The Control Plane organization that the user belongs to. " +
					"If not specified, defaults to the NUODB_CP_ORGANIZATION environment variable.",
				Optional: true,
			},
			"username": schema.StringAttribute{
				Description: "The name of the user. " +
					"If not specified, defaults to the NUODB_CP_USER environment variable.",
				Optional: true,
			},
			"password": schema.StringAttribute{
				Description: "The password for the user. " +
					"If not specified, defaults to the NUODB_CP_PASSWORD environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"url_base": schema.StringAttribute{
				Description: "The base URL for the server, including the protocol. " +
					"If not specified, defaults to the NUODB_CP_PASSWORD environment variable.",
				Optional: true,
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

	if urlBase == "" {
		urlBase = "http://localhost:8080"
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
	configuration.Servers = nuodbaas.ServerConfigurations{
		{URL: urlBase, Description: "The base URL for the server, including the protocol"},
	}
	if organization != "" && username != "" && password != "" {
		basicUsername := fmt.Sprintf("%s/%s", organization, username)
		configuration.DefaultHeader["Authorization"] = fmt.Sprintf("Basic %s", basicAuth(basicUsername, password))
	}
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
