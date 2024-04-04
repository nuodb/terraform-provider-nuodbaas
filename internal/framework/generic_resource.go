// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package framework

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.ResourceWithConfigure   = &GenericResource{}
	_ resource.ResourceWithImportState = &GenericResource{}
)

type ClientWithOptions struct {
	Client   openapi.ClientInterface
	timeouts map[string]map[string]time.Duration
}

func NewClientWithOptions(client openapi.ClientInterface, timeouts map[string]map[string]time.Duration) *ClientWithOptions {
	return &ClientWithOptions{Client: client, timeouts: timeouts}
}

// GenericResource is a Resource implementation that handles all interactions
// with the Terraform API and delegates interaction with the provider API to
// ResourceState.
type GenericResource struct {
	client                *ClientWithOptions
	TypeName              string
	Description           string
	GetResourceAttributes func() (map[string]schema.Attribute, error)
	Build                 func() ResourceState
}

// State is a marker interface for all structs that model Terraform resources
// and data sources.
type State interface{}

// ResourceState handles interactions with the provider API.
type ResourceState interface {
	State

	// Reset resets the local state of the resource to the zero value.
	Reset()

	// CheckReady returns an error if the resource is not in the desired
	// state according to its spec, based on the local state. The supplied
	// client can be used to request additional readiness information from
	// the server.
	CheckReady(ctx context.Context, client openapi.ClientInterface) error

	// Create creates the resource in the backend based on the local state
	// and waits for it to satisfy IsReady().
	Create(ctx context.Context, client openapi.ClientInterface) error

	// Read retrieves the state of the resource from the backend.
	Read(ctx context.Context, client openapi.ClientInterface) error

	// Update updates the resource in the backend to match the local state
	// and waits for it to satisfy IsReady().
	Update(ctx context.Context, client openapi.ClientInterface, currentState ResourceState) error

	// Delete deletes the resource from the backend and waits for it to be
	// cleaned up.
	Delete(ctx context.Context, client openapi.ClientInterface) error

	// Deserialize the resource ID
	SetId(id string) error
}

func (r *GenericResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.TypeName
}

func (r *GenericResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Get OpenAPI spec and convert it to Terraform schema
	attributes, err := r.GetResourceAttributes()
	if err != nil {
		resp.Diagnostics.AddError("Schema Creation Error", err.Error())
		return
	}
	resp.Schema = schema.Schema{
		Description:         r.Description,
		MarkdownDescription: r.Description,
		Attributes:          attributes,
	}
}

func getClient(diags *diag.Diagnostics, providerData any) *ClientWithOptions {
	if providerData == nil {
		return nil
	}
	client, ok := providerData.(*ClientWithOptions)
	if !ok {
		diags.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected %T, got: %T. Please report this issue to NuoDB.Support@3ds.com",
				&ClientWithOptions{}, providerData))
	}
	return client
}

func (r *GenericResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = getClient(&resp.Diagnostics, req.ProviderData)
}

func (r *GenericResource) finalizeCreateOrUpdate(ctx context.Context, state ResourceState, operation string, diags *diag.Diagnostics, tfstate *tfsdk.State) {
	// Get resource state after create or update
	err := state.Read(ctx, r.client.Client)
	if err != nil {
		diags.AddError("Unable to refresh "+r.TypeName+" after "+operation, err.Error())
		return
	}
	// Save resource into Terraform state before waiting for it to become
	// ready. This allows Terraform to manage the resource even if the
	// readiness check times out.
	diags.Append(tfstate.Set(ctx, state)...)
	if diags.HasError() {
		return
	}
	// Wait for resource to become ready
	err = r.AwaitReady(ctx, state, operation)
	if err != nil {
		diags.AddError("Unable to achieve desired state for "+r.TypeName, err.Error())
		return
	}
	// Save resource into Terraform state again now that it is ready
	diags.Append(tfstate.Set(ctx, state)...)
}

func (r *GenericResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read desired resource state from Terraform
	state := r.Build()
	if !ReadResource(ctx, &resp.Diagnostics, req.Plan.Get, state) {
		return
	}
	// Create the resource
	err := state.Create(ctx, r.client.Client)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create "+r.TypeName, err.Error())
		return
	}
	r.finalizeCreateOrUpdate(ctx, state, CREATE_OPERATION, &resp.Diagnostics, &resp.State)
}

func (r *GenericResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read resource from Terraform state
	state := r.Build()
	if !ReadResource(ctx, &resp.Diagnostics, req.State.Get, state) {
		return
	}
	// Get latest resource state
	err := state.Read(ctx, r.client.Client)
	if err != nil {
		if helper.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read "+r.TypeName, err.Error())
		return
	}
	// Save resource into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *GenericResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read desired resource state from Terraform
	plan := r.Build()
	if !ReadResource(ctx, &resp.Diagnostics, req.Plan.Get, plan) {
		return
	}
	// Read current resource state from Terraform
	state := r.Build()
	if !ReadResource(ctx, &resp.Diagnostics, req.State.Get, state) {
		return
	}
	// Update the resource
	err := plan.Update(ctx, r.client.Client, state)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update "+r.TypeName, err.Error())
		return
	}
	r.finalizeCreateOrUpdate(ctx, plan, UPDATE_OPERATION, &resp.Diagnostics, &resp.State)
}

func (r *GenericResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Read resource from Terraform state
	state := r.Build()
	if !ReadResource(ctx, &resp.Diagnostics, req.State.Get, state) {
		return
	}
	// Delete the resource
	err := state.Delete(ctx, r.client.Client)
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete "+r.TypeName, err.Error())
		return
	}
	// Wait for resource to disappear
	err = r.AwaitDeleted(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Unable to finalize deletion of "+r.TypeName, err.Error())
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
	DEFAULT_RESOURCE  = "default"
	CREATE_OPERATION  = "create"
	UPDATE_OPERATION  = "update"
	DELETE_OPERATION  = "delete"
)

type OperationTimeouts struct {
	Create *string `tfsdk:"create" hcl:"create" cty:"create"`
	Update *string `tfsdk:"update" hcl:"update" cty:"update"`
	Delete *string `tfsdk:"delete" hcl:"delete" cty:"delete"`
}

func ParseTimeouts(timeouts map[string]OperationTimeouts, resourceTypes map[string]struct{}) (map[string]map[string]time.Duration, error) {
	var errList []error
	to := make(map[string]map[string]time.Duration, len(timeouts))
	for resource, resourceTimeouts := range timeouts {
		// Validate resource name
		if resource != DEFAULT_RESOURCE {
			if _, ok := resourceTypes[resource]; !ok {
				errList = append(errList, fmt.Errorf("Invalid resource type: %s", resource))
			}
		}
		rto := make(map[string]time.Duration)
		for operation, operationTimeout := range map[string]*string{
			CREATE_OPERATION: resourceTimeouts.Create,
			UPDATE_OPERATION: resourceTimeouts.Update,
			DELETE_OPERATION: resourceTimeouts.Delete,
		} {
			if operationTimeout == nil {
				continue
			}
			// Parse time duration
			parsed, err := time.ParseDuration(*operationTimeout)
			if err != nil {
				// Trim superfluous "time: " prefix from error message
				errMsg := strings.TrimPrefix(err.Error(), "time: ")
				errList = append(errList, fmt.Errorf("Invalid timeout for %s %s: %s", resource, operation, errMsg))
				continue
			}
			if parsed < 0 {
				errList = append(errList, fmt.Errorf("Timeout for %s %s is negative: %s", resource, operation, parsed))
				continue
			}
			rto[operation] = parsed
		}
		to[resource] = rto
	}
	return to, errors.Join(errList...)
}

func (r *GenericResource) getTimeout(resource, operation string) *time.Duration {
	// Get timeouts for resource
	timeouts, ok := r.client.timeouts[resource]
	if ok {
		// Get timeout for operation
		timeout, ok := timeouts[operation]
		if ok {
			return &timeout
		}
	}
	return nil
}

func (r *GenericResource) GetTimeout(operation string, defaultTimeout time.Duration) time.Duration {
	// Get timeout for resource and operation
	if timeout := r.getTimeout(r.TypeName, operation); timeout != nil {
		return *timeout
	}
	// Get configured default timeout for operation across all resources
	if timeout := r.getTimeout(DEFAULT_RESOURCE, operation); timeout != nil {
		return *timeout
	}
	return defaultTimeout
}

func (r *GenericResource) AwaitReady(ctx context.Context, state ResourceState, operation string) error {
	timeout := r.GetTimeout(operation, READINESS_TIMEOUT)
	if timeout == 0 {
		tflog.Info(ctx, "Not waiting for resource to achieve desired state because timeout is 0",
			map[string]any{"resourceType": r.TypeName, "operation": operation})
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		// Check if resource is ready
		readyErr := state.CheckReady(ctx, r.client.Client)
		if readyErr == nil {
			return nil
		}
		time.Sleep(POLLING_INTERVAL)

		// Re-read resource state and check for timeout error
		err := state.Read(ctx, r.client.Client)
		if err != nil {
			if os.IsTimeout(err) && ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("Timed out after %s: %s", timeout, readyErr.Error())
			}
			// Return verbatim error if not timeout
			return err
		}
	}
}

func (r *GenericResource) AwaitDeleted(ctx context.Context, state ResourceState) error {
	timeout := r.GetTimeout(DELETE_OPERATION, DELETION_TIMEOUT)
	if timeout == 0 {
		tflog.Info(ctx, "Not waiting for deletion to be finalized because timeout is 0",
			map[string]any{"resourceType": r.TypeName})
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		// Try to read resource state and check that "404 Not Found"
		// error is returned
		err := state.Read(ctx, r.client.Client)
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
func ReadResource(ctx context.Context, diags *diag.Diagnostics, fn func(context.Context, any) diag.Diagnostics, dest any) bool {
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
