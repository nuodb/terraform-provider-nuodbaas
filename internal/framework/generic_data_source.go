/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package framework

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSourceWithConfigure = &GenericDataSource{}
)

// GenericDataSource is a DataSource implementation that handles all
// interactions with the Terraform API and delegates interaction with the
// provider API to DataSourceState.
type GenericDataSource struct {
	client           *openapi.Client
	TypeName         string
	Description      string
	GetOpenApiSchema func() (*openapi3.Schema, error)
	SchemaOverride   *schema.Schema
	Build            func() DataSourceState
}

// DataSourceState handles interactions with the provider API.
type DataSourceState interface {
	State

	// Read retrieves the state of the resource from the backend.
	Read(context.Context, *openapi.Client) error
}

func (d *GenericDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.TypeName
}

func (d *GenericDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// If explicit schema is supplied, return that
	if d.SchemaOverride != nil {
		resp.Schema = *d.SchemaOverride
		return
	}
	// Otherwise, build schema from OpenAPI specification
	oas, err := d.GetOpenApiSchema()
	if err != nil {
		resp.Diagnostics.AddError("Schema Creation Error", err.Error())
		return
	}
	resp.Schema = schema.Schema{
		Description:         d.Description,
		MarkdownDescription: d.Description,
		Attributes:          ToDataSourceSchema(oas),
	}
}

func (d *GenericDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	// TODO(asz6): Not sure if we should report an error in this case.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*openapi.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *openapi.Client, got: %T. Please report this issue to NuoDB.Support@3ds.com", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *GenericDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Read data source attributes from config
	state := d.Build()
	if !ReadResource(ctx, &resp.Diagnostics, req.Config.Get, state) {
		return
	}
	// Read actual data from provider
	err := state.Read(ctx, d.client)
	if err != nil {
		if helper.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading "+d.TypeName, err.Error())
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
