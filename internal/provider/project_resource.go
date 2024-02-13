/* (C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/

package provider

import (
	"context"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ ResourceState = &ProjectResourceModel{}
)

type ProjectResourceModel openapi.ProjectModel

func (state *ProjectResourceModel) Reset() {
	*state = ProjectResourceModel{}
}

func (state *ProjectResourceModel) GetResourceVersion() string {
	if state.ResourceVersion != nil {
		return *state.ResourceVersion
	}
	return ""
}

func (state *ProjectResourceModel) IsReady() bool {
	if state.Status.State == nil {
		return false
	}
	if state.Maintenance != nil && state.Maintenance.IsDisabled != nil && *state.Maintenance.IsDisabled {
		return *state.Status.State == openapi.ProjectStatusModelStateStopped
	}
	return *state.Status.State == openapi.ProjectStatusModelStateAvailable
}

func (state *ProjectResourceModel) Create(ctx context.Context, client *openapi.Client) error {
	resp, err := client.CreateProject(ctx, state.Organization, state.Name, openapi.ProjectModel(*state))
	if err != nil {
		return err
	}
	return helper.ParseResponse(resp, nil)
}

func (state *ProjectResourceModel) Read(ctx context.Context, client *openapi.Client) error {
	resp, err := client.GetProject(ctx, state.Organization, state.Name)
	if err != nil {
		return err
	}
	state.Reset()
	return helper.ParseResponse(resp, state)
}

func (state *ProjectResourceModel) Update(ctx context.Context, client *openapi.Client) error {
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
		if err == nil {
			break
		}
		err = helper.ParseResponse(resp, nil)
		if apiError, ok := err.(*helper.ApiError); !ok || apiError.GetCode() != openapi.CONCURRENTUPDATE {
			return err
		}
		// Re-fetch project and get resourceVersion
		err = latest.Read(ctx, client)
	}
	return nil
}

func (state *ProjectResourceModel) Delete(ctx context.Context, client *openapi.Client) error {
	resp, err := client.DeleteProject(ctx, state.Organization, state.Name, nil)
	if err != nil {
		return err
	}
	return helper.ParseResponse(resp, nil)
}

func NewProjectResourceModel() ResourceState {
	return &ProjectResourceModel{}
}

func NewProjectResource() resource.Resource {
	return &GenericResource{
		resourceTypeName: "project",
		description:      "Resource for managing NuoDB projects provisioned using the DBaaS Control Plane",
		getOpenApiSchema: framework.GetProjectSchema,
		build:            NewProjectResourceModel,
	}
}
