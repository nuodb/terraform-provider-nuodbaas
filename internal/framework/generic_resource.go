/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package framework

import (
	"context"
	"fmt"
	"time"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ resource.ResourceWithConfigure   = &GenericResource{}
	_ resource.ResourceWithImportState = &GenericResource{}
)

// GenericResource is a Resource implementation that handles all interactions
// with the Terraform API and delegates interaction with the provider API to
// ResourceState.
type GenericResource struct {
	client           *openapi.Client
	TypeName         string
	Description      string
	GetOpenApiSchema func() (*openapi3.Schema, error)
	Build            func() ResourceState
}

// State is a marker interface for all structs that model Terraform resources
// and data sources.
type State interface{}

// ResourceState handles interactions with the provider API.
type ResourceState interface {
	State

	// Reset resets the local state of the resource to the zero value.
	Reset()

	// GetResourceVersion returns the resource version used to guard against
	// concurrent updates from the local state.
	GetResourceVersion() string

	// IsReady returns true if the resource is in the desired state
	// according to its spec, based on the local state.
	IsReady() bool

	// Create creates the resource in the backend based on the local state
	// and waits for it to satisfy IsReady().
	Create(context.Context, *openapi.Client) error

	// Read retrieves the state of the resource from the backend.
	Read(context.Context, *openapi.Client) error

	// Update updates the resource in the backend to match the local state
	// and waits for it to satisfy IsReady().
	Update(context.Context, *openapi.Client) error

	// Delete deletes the resource from the backend and waits for it to be
	// cleaned up.
	Delete(context.Context, *openapi.Client) error

	// Deserialize the resource ID
	SetId(string) error
}

func (r *GenericResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.TypeName
}

func (r *GenericResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Get OpenAPI spec and convert it to Terraform schema
	oas, err := r.GetOpenApiSchema()
	if err != nil {
		resp.Diagnostics.AddError("Schema Creation Error", err.Error())
		return
	}
	resp.Schema = schema.Schema{
		Description:         r.Description,
		MarkdownDescription: r.Description,
		Attributes:          ToResourceSchema(oas, false),
	}
}

func (r *GenericResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = client
}

func (r *GenericResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read desired resource state from Terraform
	state := r.Build()
	if !ReadResource(ctx, &resp.Diagnostics, req.Plan.Get, state) {
		return
	}
	// Create the resource
	err := state.Create(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Error creating "+r.TypeName, err.Error())
		return
	}
	// Get resource state after creation
	err = state.Read(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Error refreshing "+r.TypeName+" after create", err.Error())
		return
	}
	// Save resource into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Wait for resource to become ready
	err = r.AwaitReady(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for "+r.TypeName+" to become ready", err.Error())
		return
	}
}

func (r *GenericResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read resource from Terraform state
	state := r.Build()
	if !ReadResource(ctx, &resp.Diagnostics, req.State.Get, state) {
		return
	}
	// Get latest resource state
	err := state.Read(ctx, r.client)
	if err != nil {
		if helper.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading "+r.TypeName, err.Error())
		return
	}
	// Save resource into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *GenericResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read desired resource state from Terraform
	state := r.Build()
	if !ReadResource(ctx, &resp.Diagnostics, req.Plan.Get, state) {
		return
	}
	// Update the resource
	err := state.Update(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Error updating "+r.TypeName, err.Error())
		return
	}
	// Get resource state after update
	err = state.Read(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Error refreshing "+r.TypeName+" after update", err.Error())
		return
	}
	// Save resource into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Wait for resource to become ready
	err = r.AwaitReady(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for "+r.TypeName+" to become ready", err.Error())
		return
	}
}

func (r *GenericResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Read resource from Terraform state
	state := r.Build()
	if !ReadResource(ctx, &resp.Diagnostics, req.State.Get, state) {
		return
	}
	// Delete the resource
	err := state.Delete(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting "+r.TypeName, err.Error())
		return
	}
	// Wait for resource to disappear
	err = r.AwaitDeleted(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for "+r.TypeName+" to be deleted", err.Error())
		return
	}
}

func (r *GenericResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := r.Build()
	err := state.SetId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

const (
	READINESS_TIMEOUT = 5 * time.Minute
	DELETION_TIMEOUT  = 1 * time.Minute
	POLLING_INTERVAL  = 1 * time.Second
)

func (r *GenericResource) AwaitReady(ctx context.Context, state ResourceState) error {
	ctx, cancel := context.WithTimeout(ctx, READINESS_TIMEOUT)
	defer cancel()
	for !state.IsReady() {
		time.Sleep(POLLING_INTERVAL)

		// Re-read resource state
		err := state.Read(ctx, r.client)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *GenericResource) AwaitDeleted(ctx context.Context, state ResourceState) error {
	ctx, cancel := context.WithTimeout(ctx, DELETION_TIMEOUT)
	defer cancel()
	for {
		// Try to read resource state and check that "404 Not Found"
		// error is returned
		err := state.Read(ctx, r.client)
		if err != nil {
			if helper.IsNotFound(err) {
				return nil
			}
			return err
		}

		time.Sleep(POLLING_INTERVAL)
	}
}

// ReadResource decodes Terraform configuration, state, or plan to a model
// struct containing ordinary Golang field types (e.g. bool, *int, []string)
// that have the `tfsdk:"..."` tag.
func ReadResource[T State](ctx context.Context, diags *diag.Diagnostics, fn func(context.Context, any) diag.Diagnostics, dest T) bool {
	// Decode to opaque object type
	var obj types.Object
	diags.Append(fn(ctx, &obj)...)
	if diags.HasError() {
		return false
	}
	// Convert to target type, ignoring null and unknown values, which
	// should deserialize as nil
	diags.Append(obj.As(ctx, dest, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})...)
	return !diags.HasError()
}
