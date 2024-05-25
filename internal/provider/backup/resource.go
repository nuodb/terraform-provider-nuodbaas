// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package backup

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/framework"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"
)

var (
	_ framework.ResourceState = &BackupResourceModel{}
)

type BackupResourceModel openapi.BackupModel

func (state *BackupResourceModel) Reset() {
	*state = BackupResourceModel{}
}

func (state *BackupResourceModel) CheckReady(ctx context.Context, client openapi.ClientInterface) error {
	// Check that backup is Succeeded
	if state.Status == nil || state.Status.State == nil {
		return fmt.Errorf("Backup %s/%s/%s/%s has no status information",
			state.Organization, state.Project, state.Database, state.Name)
	}
	if *state.Status.State != openapi.BackupStatusModelStateSucceeded {
		return fmt.Errorf("Backup %s/%s/%s/%s has an unexpected state: expected=%s, found=%s",
			state.Organization, state.Project, state.Database, state.Name,
			openapi.BackupStatusModelStateSucceeded, *state.Status.State)
	}
	return nil
}

func (state *BackupResourceModel) Create(ctx context.Context, client openapi.ClientInterface) error {
	resp, err := client.CreateOrUpdateBackup(ctx, state.Organization, state.Project, state.Database, state.Name, openapi.BackupModel(*state))
	if err != nil {
		return err
	}
	return helper.ParseResponse(resp, nil)
}

func (state *BackupResourceModel) Read(ctx context.Context, client openapi.ClientInterface) error {
	resp, err := client.GetBackup(ctx, state.Organization, state.Project, state.Database, state.Name)
	if err != nil {
		return err
	}
	state.Reset()
	return helper.ParseResponse(resp, state)
}

func (state *BackupResourceModel) Update(ctx context.Context, client openapi.ClientInterface, currentState framework.ResourceState) error {
	// Fetch database and get resourceVersion
	latest := &BackupResourceModel{
		Organization: state.Organization,
		Project:      state.Project,
		Database:     state.Database,
		Name:         state.Name,
	}
	err := latest.Read(ctx, client)
	if err != nil {
		return err
	}
	for {
		state.ResourceVersion = latest.ResourceVersion
		resp, err := client.CreateOrUpdateBackup(ctx, state.Organization, state.Project, state.Database, state.Name, openapi.BackupModel(*state))
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
		// Re-fetch database and get resourceVersion
		err = latest.Read(ctx, client)
		if err != nil {
			return err
		}
	}
}

func (state *BackupResourceModel) Delete(ctx context.Context, client openapi.ClientInterface) error {
	resp, err := client.DeleteBackup(ctx, state.Organization, state.Project, state.Database, state.Name, nil)
	if err != nil {
		return err
	}
	return helper.ParseResponse(resp, nil)
}

func (state *BackupResourceModel) SetId(id string) error {
	pathParts := strings.Split(id, "/")
	if len(pathParts) != 4 || pathParts[0] == "" || pathParts[1] == "" || pathParts[2] == "" || pathParts[3] == "" {
		return fmt.Errorf("Expected an id with format \"organization/project/database/name\". Got: %s", id)
	}
	state.Organization = pathParts[0]
	state.Project = pathParts[1]
	state.Database = pathParts[2]
	state.Name = pathParts[3]
	return nil
}

func GetBackupResourceAttributes() (map[string]schema.Attribute, error) {
	return framework.GetResourceAttributes("BackupModel")
}

func NewBackupResourceState() framework.ResourceState {
	return &BackupResourceModel{}
}

func NewBackupResource() resource.Resource {
	return &framework.GenericResource{
		TypeName:              "backup",
		Description:           "Resource for managing NuoDB backups created using the DBaaS Control Plane",
		GetResourceAttributes: GetBackupResourceAttributes,
		Build:                 NewBackupResourceState,
	}
}
