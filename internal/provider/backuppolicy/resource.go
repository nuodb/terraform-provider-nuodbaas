// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package backuppolicy

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
	_ framework.ResourceState   = &BackupPolicyResourceModel{}
	_ framework.DataSourceState = &BackupPolicyResourceModel{}
)

type BackupPolicyResourceModel openapi.BackupPolicyModel

func (state *BackupPolicyResourceModel) Reset() {
	*state = BackupPolicyResourceModel{}
}

func (state *BackupPolicyResourceModel) CheckReady(ctx context.Context, client openapi.ClientInterface) error {
	return nil
}

func (state *BackupPolicyResourceModel) Create(ctx context.Context, client openapi.ClientInterface) error {
	resp, err := client.CreateBackupPolicy(ctx, state.Organization, state.Name, openapi.BackupPolicyModel(*state))
	if err != nil {
		return err
	}
	return helper.ParseResponse(resp, nil)
}

func (state *BackupPolicyResourceModel) Read(ctx context.Context, client openapi.ClientInterface) error {
	resp, err := client.GetBackupPolicy(ctx, state.Organization, state.Name)
	if err != nil {
		return err
	}
	state.Reset()
	return helper.ParseResponse(resp, state)
}

func (state *BackupPolicyResourceModel) Update(ctx context.Context, client openapi.ClientInterface, currentState framework.ResourceState) error {
	// Fetch backup policy and get resourceVersion
	latest := &BackupPolicyResourceModel{
		Organization: state.Organization,
		Name:         state.Name,
	}
	err := latest.Read(ctx, client)
	if err != nil {
		return err
	}
	for {
		state.ResourceVersion = latest.ResourceVersion
		resp, err := client.CreateBackupPolicy(ctx, state.Organization, state.Name, openapi.BackupPolicyModel(*state))
		if err != nil {
			return err
		}
		// Decode the response and check that there is no error
		err = helper.ParseResponse(resp, nil)
		if err == nil {
			return nil
		}
		// If error is not retriable (code=CONCURRENT_UPDATE), fail fast
		if apiError, ok := err.(*helper.ApiError); !ok || apiError.GetCode() != openapi.ErrorContentCodeCONCURRENTUPDATE {
			return err
		}
		// Re-fetch policy and get resourceVersion
		err = latest.Read(ctx, client)
		if err != nil {
			return err
		}
	}
}

func (state *BackupPolicyResourceModel) Delete(ctx context.Context, client openapi.ClientInterface) error {
	resp, err := client.DeleteBackupPolicy(ctx, state.Organization, state.Name, nil)
	if err != nil {
		return err
	}
	return helper.ParseResponse(resp, nil)
}

func (state *BackupPolicyResourceModel) SetId(id string) error {
	pathParts := strings.Split(id, "/")
	if len(pathParts) != 2 || pathParts[0] == "" || pathParts[1] == "" {
		return fmt.Errorf("Expected an id with format \"organization/name\". Got: %s", id)
	}
	state.Organization = pathParts[0]
	state.Name = pathParts[1]
	return nil
}

func (state *BackupPolicyResourceModel) GetEventPath() string {
	return fmt.Sprintf("events/backuppolicies/%s/%s", state.Organization, state.Name)
}

func GetBackupPolicyResourceAttributes() (map[string]schema.Attribute, error) {
	return framework.GetResourceAttributes("BackupPolicyModel")
}

func NewBackupPolicyResourceModel() framework.ResourceState {
	return &BackupPolicyResourceModel{}
}

func NewBackupPolicyResource() resource.Resource {
	return &framework.GenericResource{
		TypeName:              "backuppolicy",
		Description:           "Resource for managing NuoDB backup policies created using the DBaaS Control Plane",
		GetResourceAttributes: GetBackupPolicyResourceAttributes,
		Build:                 NewBackupPolicyResourceModel,
	}
}
