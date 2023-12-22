/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/helper"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/model"

	nuodbaas_client "github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.ResourceWithImportState = &DatabaseResource{}
)

func NewDatabaseResource() resource.Resource {
	return &DatabaseResource{}
}

// DatabaseResource defines the resource implementation.
type DatabaseResource struct {
	client *nuodbaas.APIClient
}


func (r *DatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *DatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Database Resource",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "Name of the organization for which project is created",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the database",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The name of the project for which database is created",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				
			},
			"dba_password": schema.StringAttribute{
				MarkdownDescription: "Database password. Cannot be updated once database is created",
				Required:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tier": schema.StringAttribute{
				MarkdownDescription: "The service tier for the database. If omitted, the project service tier is inherited.",
				Optional: true,
				Computed: true,
			},
			"maintenance": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"is_disabled": schema.BoolAttribute{
						MarkdownDescription: "Whether the project or database should be shutdown",
						Optional: true,
					},
				},
			},
			"resource_version": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "The version of the resource. When specified in a PUT request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates.",
				// This plan modifier is necessary since it is used in updating the database. Without it the value of resource_version would be unknown
				PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
			},
			"properties": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"archive_disk_size": schema.StringAttribute{
						MarkdownDescription: "The size of the archive volumes for the database. Can be only updated to increase the volume size",
						Optional: true,
						Computed: true,
					},
					"journal_disk_size": schema.StringAttribute{
						MarkdownDescription: "The size of the journal volumes for the database. Can be only updated to increase the volume size.",
						Optional: true,
					},
					"tier_parameters": schema.MapAttribute{
						MarkdownDescription: "Opaque parameters supplied to database service tier.",
						ElementType: types.StringType,
						Optional: true,
					},
				},
			},
			"status": schema.SingleNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"sql_end_point": schema.StringAttribute{
						MarkdownDescription: "The endpoint for SQL clients to connect to",
						Computed: true,
					},
					"ca_pem": schema.StringAttribute{
						MarkdownDescription: "The PEM-encoded certificate for SQL clients to verify database servers",
						Computed: true,
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts {
				Create: true,
				Update: true,
			}),
		},
	}
}

func (r *DatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

	r.client = client
}

func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state model.DatabaseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	var propertiesModel model.DatabasePropertiesResourceModel

	resp.Diagnostics.Append(state.Properties.As(ctx, &propertiesModel, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true, UnhandledNullAsEmpty: true})...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, diags, cancel := r.updateContextWithTimeout(ctx, state, "create")

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	defer cancel()
	
	databaseClient := nuodbaas_client.NewDatabaseClient(r.client, ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString())
	httpResponse, err := databaseClient.CreateDatabase(state, state.Maintenance, &propertiesModel)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database",
			"Could not create database, unexpected error: "+ helper.GetHttpResponseErrorMessage(httpResponse, err),
		)
		return
	}

	getDatabaseModel, httpResponse, err := r.waitForDatabase(ctx, databaseClient)

	if err!= nil && helper.IsTimeoutError(err) {
		resp.Diagnostics.AddError("Timeout error", fmt.Sprintf("Unable to get database %+v in ready. You can go ahead and retry creating it", state.Name.ValueString()))
		databaseClient = nuodbaas_client.NewDatabaseClient(r.client, context.Background(), state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString())
		httpResponse, err = databaseClient.DeleteDatabase()
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Database",
			"Could not get NuoDbaas database " + state.Name.ValueString()+" : " + helper.GetHttpResponseErrorMessage(httpResponse, err))
		return
	}

	updateState, diags := r.updateStateWithComputedValues(ctx, &state, getDatabaseModel, "create")

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, updateState)...)

}

func (r *DatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.DatabaseResourceModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	databaseClient := nuodbaas_client.NewDatabaseClient(r.client, ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString())
	getDatabaseModel, httpResponse, err := databaseClient.GetDatabase()

	if err != nil {
		errorModel := helper.GetHttpResponseModel(httpResponse)
		if errorModel != nil && errorModel.GetStatus() == "HTTP 404 Not Found"{
			resp.State.RemoveResource(ctx)
			return
		}
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Database",
			"Could not get NuoDbaas database " + state.Name.ValueString()+" : " + helper.GetHttpResponseErrorMessage(httpResponse, err))
		return
	}

	readState, diags := r.updateStateWithComputedValues(ctx, &state, getDatabaseModel, "read")
	
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, readState)...)

}

func (r *DatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state model.DatabaseResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var propertiesModel model.DatabasePropertiesResourceModel
	resp.Diagnostics.Append(state.Properties.As(ctx, &propertiesModel, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true, UnhandledNullAsEmpty: true})...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, diags, cancel := r.updateContextWithTimeout(ctx, state, "update")

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	defer cancel()

	databaseClient := nuodbaas_client.NewDatabaseClient(r.client, ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString())
	httpResponse, err := databaseClient.UpdateDatabase(state, state.Maintenance, &propertiesModel)
	httpResponseContent := helper.GetHttpResponseModel(httpResponse)

	if httpResponseContent != nil && httpResponseContent.GetCode() == "CONCURRENT_UPDATE" {
		err = r.retryUpdate(ctx, state, state.Maintenance, &propertiesModel, databaseClient)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating database",
			fmt.Sprintf("Could not update database, unexpected error: %+v", helper.GetHttpResponseErrorMessage(httpResponse, err)),
		)
		return
	}

	getDatabaseModel, httpResponse, err := r.waitForDatabase(ctx, databaseClient)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Database",
			"Could not get NuoDbaas database " + state.Name.ValueString()+" : " + helper.GetHttpResponseErrorMessage(httpResponse, err))
		return
	}

	updatedState, diags := r.updateStateWithComputedValues(ctx, &state, getDatabaseModel, "update")

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, updatedState)...)
}

func (r *DatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.DatabaseResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResponse, err := r.client.DatabasesAPI.DeleteDatabase(ctx, state.Organization.ValueString(), state.Project.ValueString(), state.Name.ValueString()).Execute()

	if err!=nil {
		resp.Diagnostics.AddError("Unable to delete database", 
			fmt.Sprintf("Unable to database project %s, unexpected error: %v", 
			state.Name.ValueString(), helper.GetHttpResponseErrorMessage(httpResponse, err)))
		return
	}
}

func (r *DatabaseResource) retryUpdate(ctx context.Context, state model.DatabaseResourceModel, maintenanceModel *maintenanceModel, propertiesModel *model.DatabasePropertiesResourceModel, databaseClient *nuodbaas_client.NuodbaasDatabaseClient) error {
	databaseModel, _, err := databaseClient.GetDatabase()
	if err != nil {
		return err
	}
	state.ResourceVersion = types.StringValue(*databaseModel.ResourceVersion)
	_, err = databaseClient.UpdateDatabase(state, maintenanceModel, propertiesModel)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (r *DatabaseResource) waitForDatabase(ctx context.Context, databaseClient *nuodbaas_client.NuodbaasDatabaseClient) (*nuodbaas.DatabaseModel, *http.Response, error){

	var getDatabaseModel *nuodbaas.DatabaseModel
	var waitTime = 1
	for {
		databaseModel, httpResponse, err := databaseClient.GetDatabase()
		getDatabaseModel = databaseModel
		if err != nil {
			return nil, httpResponse, err
		}
		if databaseModel.Maintenance != nil {
			if !*databaseModel.Maintenance.IsDisabled && databaseModel.Status != nil && databaseModel.Status.GetState() == "Available" {
				break
			} else if *databaseModel.Maintenance.IsDisabled && databaseModel.Status != nil && databaseModel.Status.GetState() == "Stopped" {
				break
			}
		} else if databaseModel.Status != nil && databaseModel.Status.GetState() == "Available" {
			break
		}
		time.Sleep(time.Duration(waitTime) * time.Second)
		waitTime = helper.ComputeWaitTime(waitTime, 10)
	}
	return getDatabaseModel, nil, nil
}

func (r *DatabaseResource) updateStateWithComputedValues(ctx context.Context, state *model.DatabaseResourceModel, databaseModel *nuodbaas.DatabaseModel, stateType string) (*model.DatabaseResourceModel, diag.Diagnostics ) {
	
	var diagnostics diag.Diagnostics
	
	if stateType != "update" {
		state.ResourceVersion = types.StringValue(*databaseModel.ResourceVersion)
	}
	state.Tier = types.StringValue(*databaseModel.Tier)

	propertiesValue := map[string]attr.Value{
		"tier_parameters": types.MapNull(types.StringType),
		"archive_disk_size": types.StringValue(*databaseModel.Properties.ArchiveDiskSize),
		"journal_disk_size": types.StringNull(),
	}

	if databaseModel.Properties.JournalDiskSize != nil {
		propertiesValue["journal_disk_size"] = types.StringValue(*databaseModel.Properties.JournalDiskSize)
	}

	if len(databaseModel.Properties.GetTierParameters()) != 0 {
		mapValue, diags := helper.ConvertMapToTfMap(databaseModel.Properties.TierParameters)
		diagnostics = append(diagnostics, diags...)
		if diagnostics.HasError() {
			return nil, diagnostics
		}
		propertiesValue["tier_parameters"] = mapValue
	}

	if databaseModel.Status != nil {
		elementTypes := map[string]attr.Type{
			"sql_end_point": types.StringType,
			"ca_pem": types.StringType,
		}
		elements := map[string]attr.Value{
			"sql_end_point": types.StringValue(*databaseModel.Status.SqlEndpoint),
			"ca_pem": types.StringValue(*databaseModel.Status.CaPem),
		}
		status, diags := types.ObjectValue(elementTypes, elements)
		diagnostics = append(diagnostics, diags...)
		if diagnostics.HasError() {
			return nil, diagnostics
		}
		state.Status = status
	}

	convertPropertiesType := map[string]attr.Type{
		"archive_disk_size" : types.StringType,
		"journal_disk_size" : types.StringType,
		"tier_parameters" : types.MapType{ElemType: types.StringType},
	}
	
	state.Properties = basetypes.NewObjectValueMust(convertPropertiesType, propertiesValue)

	return state, diagnostics
}

func (r *DatabaseResource) updateContextWithTimeout(ctx context.Context, state model.DatabaseResourceModel, timeoutType string) (context.Context, diag.Diagnostics, context.CancelFunc) {
	var timeout time.Duration
	var diags diag.Diagnostics

	if timeoutType == "create" {
		timeout, diags = state.Timeouts.Create(ctx, 30*time.Minute)
	} else if timeoutType == "update" {
		timeout, diags = state.Timeouts.Update(ctx, 30*time.Minute)
	}

	if diags.HasError() {
		return ctx, diags, nil
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)

	return ctx, diags, cancel
}

func (r *DatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}