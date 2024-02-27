/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"

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

func (state *DatabaseResourceModel) GetResourceVersion() string {
	if state.ResourceVersion != nil {
		return *state.ResourceVersion
	}
	return ""
}

func (state *DatabaseResourceModel) CheckReady() error {
	if state.Status == nil || state.Status.State == nil {
		return fmt.Errorf("Database %s/%s/%s has no status information", state.Organization, state.Project, state.Name)
	}
	expectedState := openapi.DatabaseStatusModelStateAvailable
	if state.Maintenance != nil && state.Maintenance.IsDisabled != nil && *state.Maintenance.IsDisabled {
		expectedState = openapi.DatabaseStatusModelStateStopped
	}
	if *state.Status.State != expectedState {
		return fmt.Errorf("Database %s/%s/%s has an unexpected state: expected=%s, found=%s",
			state.Organization, state.Project, state.Name, expectedState, *state.Status.State)
	}
	return nil
}

func (state *DatabaseResourceModel) Create(ctx context.Context, client *openapi.Client) error {
	resp, err := client.CreateDatabase(ctx, state.Organization, state.Project, state.Name, openapi.DatabaseCreateUpdateModel(*state))
	if err != nil {
		return err
	}
	return helper.ParseResponse(resp, nil)
}

func (state *DatabaseResourceModel) Read(ctx context.Context, client *openapi.Client) error {
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

func (state *DatabaseResourceModel) Update(ctx context.Context, client *openapi.Client) error {
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
		if apiError, ok := err.(*helper.ApiError); !ok || apiError.GetCode() != openapi.CONCURRENTUPDATE {
			return err
		}
		// Re-fetch database and get resourceVersion
		err = latest.Read(ctx, client)
		if err != nil {
			return err
		}
	}
}

func (state *DatabaseResourceModel) Delete(ctx context.Context, client *openapi.Client) error {
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

func NewDatabaseResourceState() framework.ResourceState {
	return &DatabaseResourceModel{}
}

func NewDatabaseResource() resource.Resource {
	return &framework.GenericResource{
		TypeName:         "database",
		Description:      "Resource for managing NuoDB databases created using the DBaaS Control Plane",
		GetOpenApiSchema: framework.GetDatabaseResourceSchema,
		Build:            NewDatabaseResourceState,
	}
}
