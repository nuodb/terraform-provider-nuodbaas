/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/helper"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
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
}

func (p *NuoDbaasProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nuodbaas"
	resp.Version = p.version
}

func (p *NuoDbaasProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				Description: "Name of the organization",
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
				Description: "Username for Dbaas Client.",
			},
			"password": schema.StringAttribute{
				Optional: true,
				Sensitive: true,
				Description: "Password for Dbaas Client",
			},
			"url_base": schema.StringAttribute{
				Optional: true,
				Description: "The base URL for the server, including the protocol",
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

	if config.Organization.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("organization"),
			"Unknown Dbaas organization type",
			"The provider cannot create the NuoDB Dbaas API client as there is an unknown configuration value for the NuoDB Dbaas API organization. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the NUODB_CP_ORGANIZATION environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Dbaas username type",
			"The provider cannot create the NuoDB Dbaas API client as there is an unknown configuration value for the NuoDB Dbaas API username. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the NUODB_CP_USER environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Dbaas password type",
			"The provider cannot create the NuoDB Dbaas API client as there is an unknown configuration value for the NuoDB Dbaas API password. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the NUODB_CP_PASSWORD environment variable.",
		)
	}

	if config.BaseUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url_base"),
			"Unknown url type",
			"The provider cannot create the NuoDB Dbaas API client as there is an unknown configuration value for the Url. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the NUODB_CP_URL_BASE environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	organization := os.Getenv("NUODB_CP_ORGANIZATION")
    username := os.Getenv("NUODB_CP_USER")
    password := os.Getenv("NUODB_CP_PASSWORD")
	host := os.Getenv("NUODB_CP_URL_BASE")

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
		host = config.BaseUrl.ValueString()
	}

	if organization == "" {
		resp.Diagnostics.AddAttributeError(
            path.Root("organization"),
            "Missing NuoDB DBAAS API Organization",
			helper.GetProviderValidatorErrorMessage("organization", "NUODB_CP_ORGANIZATION"),
        )
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
            path.Root("username"),
            "Missing NuoDB DBAAS API Username",
			helper.GetProviderValidatorErrorMessage("username", "NUODB_CP_USER"),
        )
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
            path.Root("password"),
            "Missing NuoDB DBAAS API Password",
			helper.GetProviderValidatorErrorMessage("password", "NUODB_CP_PASSWORD"),
        )
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
            path.Root("host"),
            "Missing Url base",
			helper.GetProviderValidatorErrorMessage("host", "NUODB_CP_URL_BASE"),
        )
	}

	if resp.Diagnostics.HasError() {
		return
	}

	configuration := nuodbaas.NewConfiguration()
	serverConfig := nuodbaas.ServerConfigurations{{URL: host, Description: "The base URL for the server, including the protocol"}}
	basicUsername := fmt.Sprintf("%s/%s", organization, username)
	configuration.DefaultHeader["Authorization"] = fmt.Sprintf("Basic %s", basicAuth(basicUsername, password))
	configuration.Servers = serverConfig
	apiClient := nuodbaas.NewAPIClient(configuration)
	httpsRes, error := apiClient.HealthzAPI.GetHealth(context.Background()).Execute()
	if httpsRes.StatusCode != 403 && httpsRes.StatusCode >= 300 {
		resp.Diagnostics.AddError(
            "Something went wrong",
			fmt.Sprintf("Something went wrong %s", error.Error()),
        )
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *NuoDbaasProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource {
		NewProjectResource,
		NewDatabaseResource,
	}
}

func (p *NuoDbaasProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource {
        NewProjectDataSource,
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
