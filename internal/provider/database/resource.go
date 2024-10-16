// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package database

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"
)

var (
	_ framework.ResourceState = &DatabaseResourceModel{}
)

type DatabaseResourceModel openapi.DatabaseCreateUpdateModel

func (state *DatabaseResourceModel) Reset() {
	*state = DatabaseResourceModel{}
}

func (state *DatabaseResourceModel) DbaPasswordMatches(other *DatabaseResourceModel) bool {
	if other != nil {
		if state.DbaPassword == nil || other.DbaPassword == nil {
			return state.DbaPassword == other.DbaPassword
		}
		return *state.DbaPassword == *other.DbaPassword
	}
	return true
}

func (state *DatabaseResourceModel) CheckReady(ctx context.Context, client openapi.ClientInterface) error {
	// Check that database is either Available or Stopped based on maintenance.isDisabled value
	if state.Status == nil || state.Status.State == nil {
		return fmt.Errorf("Database %s/%s/%s has no status information", state.Organization, state.Project, state.Name)
	}
	expectedState := openapi.DatabaseStatusModelStateAvailable
	if state.Maintenance != nil && state.Maintenance.IsDisabled != nil && *state.Maintenance.IsDisabled {
		expectedState = openapi.DatabaseStatusModelStateStopped
	}
	switch *state.Status.State {
	case expectedState:
		break
	case openapi.DatabaseStatusModelStateFailed:
		message := "unknown reason"
		if state.Status.Message != nil {
			message = *state.Status.Message
		}
		return framework.ResourceFailed("Database %s/%s/%s failed: %s",
			state.Organization, state.Project, state.Name, message)
	default:
		return fmt.Errorf("Database %s/%s/%s has an unexpected state: expected=%s, found=%s",
			state.Organization, state.Project, state.Name, expectedState, *state.Status.State)
	}
	// Check that DBA password is up-to-date if password update is supported
	resp, err := state.UpdateDbaPassword(ctx, client, nil)
	if IsDbaPasswordUpdateUnsupportedError(resp, err) {
		return nil
	}
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DBA password for database %s/%s/%s has not been updated",
			state.Organization, state.Project, state.Name)
	}
	return nil
}

func (state *DatabaseResourceModel) Create(ctx context.Context, client openapi.ClientInterface) error {
	resp, err := client.CreateDatabase(ctx, state.Organization, state.Project, state.Name, openapi.DatabaseCreateUpdateModel(*state))
	if err != nil {
		return err
	}
	return helper.ParseResponse(resp, nil)
}

func (state *DatabaseResourceModel) Read(ctx context.Context, client openapi.ClientInterface) error {
	// If DBA password is set, then this is invoked in the context of create
	// or update, to refresh the state. Make sure to save the DBA password,
	// since it is not returned by GET response.
	dbaPassword := state.DbaPassword
	resp, err := client.GetDatabase(ctx, state.Organization, state.Project, state.Name)
	if err != nil {
		return err
	}
	state.Reset()
	state.DbaPassword = dbaPassword
	return helper.ParseResponse(resp, state)
}

const (
	DBA_PASSWORD_CHANGE_UNSUPPORTED_MSG = "Configured DBA password was changed and the server does not support updating the DBA password. " +
		"Revert the configured DBA password to the value in the Terraform state and retry."
)

func IsDbaPasswordUpdateUnsupportedError(resp *http.Response, err error) bool {
	if err != nil {
		// "404 Not Found" is returned with no "detail" message if /dbaPassword	sub-resource is not supported
		if resp.StatusCode == http.StatusNotFound {
			if apiError, ok := err.(*helper.ApiError); !ok || apiError.GetDetail() == "" {
				return true
			}
		}
	}
	return false
}

func (state *DatabaseResourceModel) UpdateDbaPassword(ctx context.Context, client openapi.ClientInterface, target *string) (*http.Response, error) {
	request := openapi.UpdateDbaPasswordModel{
		Current: *state.DbaPassword,
		Target:  target,
	}
	resp, err := client.UpdateDbaPassword(ctx, state.Organization, state.Project, state.Name, nil, request)
	if err != nil {
		return resp, err
	}
	// Decode the response and check that there is no error
	return resp, helper.ParseResponse(resp, nil)
}

func (state *DatabaseResourceModel) Update(ctx context.Context, client openapi.ClientInterface, currentState framework.ResourceState) error {
	// Try to update DBA password if it was changed in config
	currentDatabase, _ := currentState.(*DatabaseResourceModel)
	if !state.DbaPasswordMatches(currentDatabase) {
		resp, err := currentDatabase.UpdateDbaPassword(ctx, client, state.DbaPassword)
		if IsDbaPasswordUpdateUnsupportedError(resp, err) {
			return errors.New(DBA_PASSWORD_CHANGE_UNSUPPORTED_MSG)
		}
		if err != nil {
			return err
		}
	}
	// Fetch database and get resourceVersion
	latest := &DatabaseResourceModel{
		Organization: state.Organization,
		Project:      state.Project,
		Name:         state.Name,
	}
	err := latest.Read(ctx, client)
	if err != nil {
		return err
	}
	// Stash DBA password and set it to null in request, since PUT requests
	// do not accept it
	dbaPassword := state.DbaPassword
	state.DbaPassword = nil
	for {
		state.ResourceVersion = latest.ResourceVersion
		resp, err := client.CreateDatabase(ctx, state.Organization, state.Project, state.Name, openapi.DatabaseCreateUpdateModel(*state))
		if err != nil {
			return err
		}
		// Decode the response and check that there is no error
		err = helper.ParseResponse(resp, nil)
		if err == nil {
			// Add back DBA password so that it is preserved in state
			state.DbaPassword = dbaPassword
			return nil
		}
		// If error is not retriable (code=CONCURRENT_UPDATE), fail fast
		if apiError, ok := err.(*helper.ApiError); !ok || apiError.GetCode() != openapi.ErrorContentCodeCONCURRENTUPDATE {
			return err
		}
		// Re-fetch database and get resourceVersion
		err = latest.Read(ctx, client)
		if err != nil {
			return err
		}
	}
}

func (state *DatabaseResourceModel) Delete(ctx context.Context, client openapi.ClientInterface) error {
	resp, err := client.DeleteDatabase(ctx, state.Organization, state.Project, state.Name, nil)
	if err != nil {
		return err
	}
	return helper.ParseResponse(resp, nil)
}

func (state *DatabaseResourceModel) SetId(id string) error {
	pathParts := strings.Split(id, "/")
	if len(pathParts) != 3 || pathParts[0] == "" || pathParts[1] == "" || pathParts[2] == "" {
		return fmt.Errorf("Expected an id with format \"organization/project/name\". Got: %s", id)
	}
	state.Organization = pathParts[0]
	state.Project = pathParts[1]
	state.Name = pathParts[2]
	return nil
}

func (state *DatabaseResourceModel) GetEventPath() string {
	return fmt.Sprintf("events/databases/%s/%s/%s", state.Organization, state.Project, state.Name)
}

func GetDatabaseResourceAttributes() (map[string]schema.Attribute, error) {
	return framework.GetResourceAttributes("DatabaseCreateUpdateModel",
		// DBA password can be updated from configuration, so remove note about it only being accepted on create
		framework.WithDescription("dbaPassword", "The password for the DBA user"),
		// Require fully-qualified backup name to prevent normalization from causing Terraform to fail due to change in attribute value
		framework.WithDescription("restoreFrom.backup", "The fully-qualified name of the backup to restore the database from"),
		framework.WithPattern("restoreFrom.backup", "([a-z][a-z0-9]*/){3}([0-9]+|[a-z][a-z0-9]*)"),
	)
}

func NewDatabaseResourceState() framework.ResourceState {
	return &DatabaseResourceModel{}
}

func NewDatabaseResource() resource.Resource {
	return &framework.GenericResource{
		TypeName:              "database",
		Description:           "Resource for managing NuoDB databases created using the DBaaS Control Plane",
		GetResourceAttributes: GetDatabaseResourceAttributes,
		Build:                 NewDatabaseResourceState,
	}
}
