// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	_ framework.ResourceState   = &ProjectResourceModel{}
	_ framework.DataSourceState = &ProjectResourceModel{}
)

type ProjectResourceModel openapi.ProjectModel

func (state *ProjectResourceModel) Reset() {
	*state = ProjectResourceModel{}
}

func (state *ProjectResourceModel) CheckReady(ctx context.Context, client openapi.ClientInterface) error {
	if state.Status == nil || state.Status.State == nil {
		return fmt.Errorf("Project %s/%s has no status information", state.Organization, state.Name)
	}
	expectedState := openapi.ProjectStatusModelStateAvailable
	if state.Maintenance != nil && state.Maintenance.IsDisabled != nil && *state.Maintenance.IsDisabled {
		expectedState = openapi.ProjectStatusModelStateStopped
	}
	if *state.Status.State != expectedState {
		return fmt.Errorf("Project %s/%s has an unexpected state: expected=%s, found=%s",
			state.Organization, state.Name, expectedState, *state.Status.State)
	}
	return nil
}

func (state *ProjectResourceModel) Create(ctx context.Context, client openapi.ClientInterface) error {
	resp, err := client.CreateProject(ctx, state.Organization, state.Name, openapi.ProjectModel(*state))
	if err != nil {
		return err
	}
	return helper.ParseResponse(resp, nil)
}

func (state *ProjectResourceModel) Read(ctx context.Context, client openapi.ClientInterface) error {
	resp, err := client.GetProject(ctx, state.Organization, state.Name)
	if err != nil {
		return err
	}
	state.Reset()
	return helper.ParseResponse(resp, state)
}

func (state *ProjectResourceModel) Update(ctx context.Context, client openapi.ClientInterface, currentState framework.ResourceState) error {
	// Fetch project and get resourceVersion
	latest := &ProjectResourceModel{
		Organization: state.Organization,
		Name:         state.Name,
	}
	err := latest.Read(ctx, client)
	if err != nil {
		return err
	}
	for {
		state.ResourceVersion = latest.ResourceVersion
		resp, err := client.CreateProject(ctx, state.Organization, state.Name, openapi.ProjectModel(*state))
		if err != nil {
			return err
		}
		// Decode the response and check that there is no error
		err = helper.ParseResponse(resp, nil)
		if err == nil {
			return nil
		}
		// If error is not retriable (code=CONCURRENT_UPDATE), fail fast
		if apiError, ok := err.(*helper.ApiError); !ok || apiError.GetCode() != openapi.CONCURRENTUPDATE {
			return err
		}
		// Re-fetch project and get resourceVersion
		err = latest.Read(ctx, client)
		if err != nil {
			return err
		}
	}
}

func (state *ProjectResourceModel) Delete(ctx context.Context, client openapi.ClientInterface) error {
	resp, err := client.DeleteProject(ctx, state.Organization, state.Name, nil)
	if err != nil {
		return err
	}
	return helper.ParseResponse(resp, nil)
}

func (state *ProjectResourceModel) SetId(id string) error {
	pathParts := strings.Split(id, "/")
	if len(pathParts) != 2 || pathParts[0] == "" || pathParts[1] == "" {
		return fmt.Errorf("Expected an id with format \"organization/name\". Got: %s", id)
	}
	state.Organization = pathParts[0]
	state.Name = pathParts[1]
	return nil
}

func GetProjectResourceAttributes() (map[string]schema.Attribute, error) {
	return framework.GetResourceAttributes("ProjectModel")
}

func NewProjectResourceModel() framework.ResourceState {
	return &ProjectResourceModel{}
}

func NewProjectResource() resource.Resource {
	return &framework.GenericResource{
		TypeName:              "project",
		Description:           "Resource for managing NuoDB projects created using the DBaaS Control Plane",
		GetResourceAttributes: GetProjectResourceAttributes,
		Build:                 NewProjectResourceModel,
	}
}
