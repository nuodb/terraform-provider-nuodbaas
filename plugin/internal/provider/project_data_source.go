package provider

import (
	"context"
	"fmt"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/helper"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/model"

	nuodbaas_client "github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

var _ datasource.DataSourceWithConfigure = &projectDataSource{}

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

type projectDataSource struct {
	client *nuodbaas.APIClient
}

// Schema implements datasource.DataSource.
func (d *projectDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "Name of the organization for which project is created",
				Required: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the project",
				Required: true,
			},
			"sla": schema.StringAttribute{
				MarkdownDescription: "The SLA for the project. Cannot be updated once the project is created.",
				Optional: true,
			},
			"tier": schema.StringAttribute{
				MarkdownDescription: "The Tier for the project. Cannot be updated once the project is created.",
				Optional: true,
			},
			"maintenance": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"expires_in": schema.StringAttribute{
						MarkdownDescription: "The time until the project or database is disabled, e.g. 1d",
						Optional: true,
					},
					"is_disabled": schema.BoolAttribute{
						Optional: true,
					},
					"expires_at_time": schema.StringAttribute{
						Optional: true,
					},
				},
			},
			"resource_version": schema.StringAttribute{
				Computed: true,
			},
			"properties": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"tier_parameters": schema.MapAttribute{
						Optional: true,
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

// Metadata implements datasource.DataSource.
func (d *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Read implements datasource.DataSource.
func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state model.ProjectResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectClient := nuodbaas_client.NewProjectClient(d.client,ctx,state.Organization.ValueString(),state.Name.ValueString())

	project, httpResponse, err := projectClient.GetProject()
	
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting projects",
			"Could not get projects, unexpected error: "+ helper.GetHttpResponseErrorMessage(httpResponse, err),
		)
		return
	}

	projectStateModel := model.ProjectResourceModel{
		Organization: types.StringValue(*project.Organization),
		Name: types.StringValue(*project.Name),
		Sla: types.StringValue(project.Sla),
		Tier: types.StringValue(project.Tier),
		ResourceVersion: types.StringValue(*project.ResourceVersion),
	}

	if project.Maintenance != nil {
		maintenanceModel := model.MaintenanceModel{}
		if project.Maintenance.ExpiresIn != nil {
			maintenanceModel.ExpiresIn = types.StringValue(*project.Maintenance.ExpiresIn)
		}
	
		if project.Maintenance.IsDisabled != nil {
			maintenanceModel.IsDisabled = types.BoolValue(*project.Maintenance.IsDisabled)
		}

		if project.Maintenance.ExpiresAtTime != nil {
			maintenanceModel.ExpiresAtTime = types.StringValue(project.Maintenance.ExpiresAtTime.String())
		}
		projectStateModel.Maintenance = &maintenanceModel
	}

	if project.Properties != nil {
		properties := model.ProjectProperties{}
		if project.Properties.TierParameters != nil {
			mapValue, diag := helper.ConvertMapToTfMap(project.Properties.TierParameters)
			resp.Diagnostics.Append(diag...)
			if resp.Diagnostics.HasError() {
				return
			}
			properties.TierParameters = mapValue
		}
		projectStateModel.Properties = &properties
	}

	state = projectStateModel

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure implements datasource.DataSourceWithConfigure.
func (d *projectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*nuodbaas.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *openapi.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}