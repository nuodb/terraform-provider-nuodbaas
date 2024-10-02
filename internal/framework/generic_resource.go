// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package framework

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"sync"
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
	"github.com/tmaxmax/go-sse"
)

var (
	_ resource.ResourceWithConfigure   = &GenericResource{}
	_ resource.ResourceWithImportState = &GenericResource{}
)

type ProviderConfig interface {
	// GetUser returns the user name in the provider configuration.
	GetUser() string

	// GetPassword returns the password for the user.
	GetPassword() string

	// GetToken returns the authentication token.
	GetToken() string

	// GetUrlBase returns the URL base.
	GetUrlBase() string

	// GetSkipVerify returns whether certificate verification should be skipped.
	GetSkipVerify() bool

	// CreateClient creates a REST API client.
	CreateClient() (openapi.ClientInterface, error)

	// ConsumeEvents creates an SSE connection and consumes events using the
	// supplied callback until the connection is closed.
	ConsumeEvents(ctx context.Context, path string, callback func(sse.Event)) error
}

type ProviderClient struct {
	ProviderConfig ProviderConfig
	Client         openapi.ClientInterface
	timeouts       map[string]map[string]time.Duration
}

func NewProviderClient(providerConfig ProviderConfig, client openapi.ClientInterface, timeouts map[string]map[string]time.Duration) *ProviderClient {
	return &ProviderClient{ProviderConfig: providerConfig, Client: client, timeouts: timeouts}
}

// GenericResource is a Resource implementation that handles all interactions
// with the Terraform API and delegates interaction with the provider API to
// ResourceState.
type GenericResource struct {
	client                *ProviderClient
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

	// Populate the local state with the resource ID.
	SetId(id string) error

	// Get the path to obtain an event stream for the resource.
	GetEventPath() string
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

func getClient(diags *diag.Diagnostics, providerData any) *ProviderClient {
	if providerData == nil {
		return nil
	}
	client, ok := providerData.(*ProviderClient)
	if !ok {
		diags.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected %T, got: %T. Please report this issue to NuoDB.Support@3ds.com",
				&ProviderClient{}, providerData))
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
	READINESS_TIMEOUT = 10 * time.Minute
	DELETION_TIMEOUT  = 1 * time.Minute
	POLLING_INTERVAL  = 10 * time.Second
	FAILURE_THRESHOLD = POLLING_INTERVAL + 1*time.Second
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

type resourceFailedError struct {
	message string
}

func ResourceFailed(format string, args ...any) error {
	return &resourceFailedError{message: fmt.Sprintf(format, args...)}
}

func (err *resourceFailedError) Error() string {
	return err.message
}

func isNetworkError(err error) bool {
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}
	urlErr, ok := err.(*url.Error)
	return ok && urlErr.Temporary()
}

const (
	SSE_EVENT_RESYNC     = "RESYNC"
	SSE_EVENT_CREATED    = "CREATED"
	SSE_EVENT_UPDATED    = "UPDATED"
	SSE_EVENT_DELETED    = "DELETED"
	SSE_EVENT_HEARTBEAT  = "HEARTBEAT"
	SSE_DATA_NO_RESOURCE = "null"
)

type eventStream struct {
	wg   sync.WaitGroup
	lock sync.Mutex
	// eventChannel is used to signal events on the resource. Each value
	// consumed from the channel indicates whether the resource existed when
	// the event was generated, in which case, its state can be observed in
	// the ResourceState that data is deserialized to.
	eventChannel chan bool
	// errChannel is used to notify the consumer about errors encountered
	// during event processing.
	errChannel chan error
}

func (stream *eventStream) withLock(fn func() error) error {
	stream.lock.Lock()
	defer stream.lock.Unlock()
	return fn()
}

func sendToChannel[T any](ctx context.Context, ch chan T, v T) {
	select {
	case ch <- v:
	case <-ctx.Done():
		tflog.Debug(ctx, "Cancelling send because context is done")
	}
}

func (r *GenericResource) stream(ctx context.Context, state ResourceState) *eventStream {
	stream := eventStream{
		eventChannel: make(chan bool),
		errChannel:   make(chan error),
	}
	stream.wg.Add(1)
	go func() {
		defer stream.wg.Done()
		// Try to use SSE to stream events on resource
		err := r.client.ProviderConfig.ConsumeEvents(ctx, state.GetEventPath(), func(event sse.Event) {
			// Do nothing on heartbeat messages
			if event.Type == SSE_EVENT_HEARTBEAT {
				return
			}
			// If event is DELETED or has no data, notify that the resource was deleted
			if event.Type == SSE_EVENT_DELETED || event.Data == SSE_DATA_NO_RESOURCE {
				sendToChannel(ctx, stream.eventChannel, false)
			} else {
				// Deserialize message and notify event
				err := stream.withLock(func() error {
					tflog.Debug(ctx, "Unmarshalling data from SSE message",
						map[string]any{"event": event.Type, "data": event.Data})
					return json.Unmarshal([]byte(event.Data), state)
				})
				if err == nil {
					sendToChannel(ctx, stream.eventChannel, true)
				} else {
					sendToChannel(ctx, stream.errChannel, err)
				}
			}
		})
		// If context is not done, there was an error consuming SSE messages
		select {
		case <-ctx.Done():
			return
		default:
			tflog.Info(ctx, "Downgrading from SSE to polling", map[string]any{"error": err})
		}
		// Use polling to generate resource events
		for {
			err := stream.withLock(func() error {
				return state.Read(ctx, r.client.Client)
			})
			if err == nil {
				sendToChannel(ctx, stream.eventChannel, true)
			} else if helper.IsNotFound(err) {
				sendToChannel(ctx, stream.eventChannel, false)
			} else if isNetworkError(err) {
				// Suppress network errors, which may be transient and retriable
				tflog.Info(ctx, "Suppressing network error while awaiting readiness",
					map[string]any{"resourceType": r.TypeName, "error": err.Error()})
			} else {
				sendToChannel(ctx, stream.errChannel, err)
			}
			// Wait for polling interval or until context is done
			select {
			case <-time.After(POLLING_INTERVAL):
				continue
			case <-ctx.Done():
				return
			}
		}
	}()
	return &stream
}

func (r *GenericResource) AwaitReady(ctx context.Context, state ResourceState, operation string) error {
	timeout := r.GetTimeout(operation, READINESS_TIMEOUT)
	if timeout == 0 {
		tflog.Info(ctx, "Not waiting for resource to achieve desired state because timeout is 0",
			map[string]any{"resourceType": r.TypeName, "operation": operation})
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	stream := r.stream(ctx, state)
	defer func() {
		cancel()
		stream.wg.Wait()
	}()
	var readyErr error
	var failedSince time.Time
	for {
		var err error
		done := false
		select {
		case exists := <-stream.eventChannel:
			// Check that resource still exists
			if !exists {
				return errors.New("Resource no longer exists")
			}
		case <-time.After(FAILURE_THRESHOLD):
			// Check readiness periodically even if not triggered by channel
		case err = <-stream.errChannel:
			// Check error encountered on channel
		case <-ctx.Done():
			// Check that timeout has not expired
			done = true
		}
		// If error was encountered on channel or timeout expired, return error
		if err != nil || done {
			if done || os.IsTimeout(err) && ctx.Err() == context.DeadlineExceeded {
				if readyErr != nil {
					return fmt.Errorf("Timed out after %s: %s", timeout, readyErr.Error())
				} else {
					return fmt.Errorf("Timed out after %s", timeout)
				}
			}
			return err
		}
		// Check if resource is ready
		readyErr = stream.withLock(func() error {
			return state.CheckReady(ctx, r.client.Client)
		})
		if readyErr == nil {
			return nil
		}
		// Return early if resource in failed state for some time
		if _, ok := readyErr.(*resourceFailedError); ok {
			if failedSince.IsZero() {
				failedSince = time.Now()
			} else if failedSince.Add(FAILURE_THRESHOLD).Before(time.Now()) {
				return readyErr
			}
		} else {
			failedSince = time.Time{}
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
	stream := r.stream(ctx, state)
	defer func() {
		cancel()
		stream.wg.Wait()
	}()
	for {
		var err error
		done := false
		select {
		case exists := <-stream.eventChannel:
			// Check if resource does not exist
			if !exists {
				return nil
			}
		case err = <-stream.errChannel:
			// Check error encountered on channel
		case <-ctx.Done():
			// Check that timeout has not expired
			done = true
		}
		// If error was encountered on channel or timeout expired, return error
		if err != nil || done {
			if done || os.IsTimeout(err) && ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("Timed out after %s", timeout)
			}
			return err
		}
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
