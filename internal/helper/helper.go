// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/nuodb/terraform-provider-nuodbaas/openapi"
)

func IsNotFound(err error) bool {
	apiError, ok := err.(*ApiError)
	return ok && apiError.GetStatusCode() == http.StatusNotFound
}

var _ error = &ApiError{}

type ApiError struct {
	HttpResponse *http.Response
	ErrorContent openapi.ErrorContent
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("Error response from Control Plane service: status='%s', code=%s, detail=[%s]", e.GetStatus(), e.GetCode(), e.GetDetail())
}

func (e *ApiError) GetCode() openapi.ErrorContentCode {
	if e.ErrorContent.Code != nil {
		return *e.ErrorContent.Code
	}
	return ""
}

func (e *ApiError) GetDetail() string {
	if e.ErrorContent.Detail != nil {
		return *e.ErrorContent.Detail
	}
	return ""
}

func (e *ApiError) GetStatus() string {
	if e.ErrorContent.Status != nil {
		return *e.ErrorContent.Status
	}
	return ""
}

func (e *ApiError) GetStatusCode() int {
	if e.HttpResponse != nil {
		return e.HttpResponse.StatusCode
	}
	return 0
}

func ParseResponse(resp *http.Response, dest any) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		return err
	}
	// Decode JSON response
	if strings.Contains(resp.Header.Get("Content-Type"), "json") {
		// Decode ErrorContent response
		if resp.StatusCode >= http.StatusBadRequest {
			apiError := ApiError{HttpResponse: resp}
			if err := json.Unmarshal(bodyBytes, &apiError.ErrorContent); err != nil {
				return err
			}
			return &apiError
		}
		// Decode response to supplied target
		if dest != nil {
			if err := json.Unmarshal(bodyBytes, &dest); err != nil {
				return err
			}
		}
	}
	// If an error response with an unexpected Content-Type was returned, return an error
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("Unexpected response: status=[%s], content=[%s]", resp.Status, string(bodyBytes))
	}
	return nil
}

func processListResponse(prefix string, resp *http.Response, err error) ([]string, error) {
	// Make sure request was successful
	if err != nil {
		return nil, err
	}
	// Decode as ItemList
	var itemList openapi.ItemList
	err = ParseResponse(resp, &itemList)
	if err != nil {
		return nil, err
	}
	// Return resource names with prefix
	var names []string
	if itemList.Items != nil {
		for _, item := range *itemList.Items {
			name, err := item.AsItemListItems0()
			if err != nil {
				return nil, err
			}
			names = append(names, prefix+name)
		}
	}
	return names, nil
}

func GetBackups(ctx context.Context, client openapi.ClientInterface, organization, project, database string, labelFilter *string, listAccessible bool) ([]string, error) {
	var prefix string
	var resp *http.Response
	var err error
	if len(organization) == 0 {
		// Make sure project was not specified without organization
		if len(project) != 0 {
			return nil, fmt.Errorf("Cannot specify project filter (%s) without organization", project)
		}
		// List all backups
		params := openapi.GetAllBackupsParams{
			LabelFilter:    labelFilter,
			ListAccessible: &listAccessible,
		}
		resp, err = client.GetAllBackups(ctx, &params)
	} else if len(project) == 0 {
		// Make sure database was not specified without project
		if len(database) != 0 {
			return nil, fmt.Errorf("Cannot specify database filter (%s) without project", database)
		}
		// List all backups within organization
		prefix = organization + "/"
		params := openapi.GetOrganizationBackupsParams{
			LabelFilter:    labelFilter,
			ListAccessible: &listAccessible,
		}
		resp, err = client.GetOrganizationBackups(ctx, organization, &params)
	} else if len(database) == 0 {
		// List all backups within project
		prefix = organization + "/" + project + "/"
		params := openapi.GetProjectBackupsParams{
			LabelFilter:    labelFilter,
			ListAccessible: &listAccessible,
		}
		resp, err = client.GetProjectBackups(ctx, organization, project, &params)
	} else {
		// List all backups within database
		prefix = organization + "/" + project + "/" + database + "/"
		params := openapi.GetBackupsParams{
			LabelFilter:    labelFilter,
			ListAccessible: &listAccessible,
		}
		resp, err = client.GetBackups(ctx, organization, project, database, &params)
	}
	return processListResponse(prefix, resp, err)
}

func GetDatabases(ctx context.Context, client openapi.ClientInterface, organization, project string, labelFilter *string, listAccessible bool) ([]string, error) {
	var prefix string
	var resp *http.Response
	var err error
	if len(organization) == 0 {
		// Make sure project was not specified without organization
		if len(project) != 0 {
			return nil, fmt.Errorf("Cannot specify project filter (%s) without organization", project)
		}
		// List all databases
		params := openapi.GetAllDatabasesParams{
			LabelFilter:    labelFilter,
			ListAccessible: &listAccessible,
		}
		resp, err = client.GetAllDatabases(ctx, &params)
	} else if len(project) == 0 {
		// List all databases within organization
		prefix = organization + "/"
		params := openapi.GetOrganizationDatabasesParams{
			LabelFilter:    labelFilter,
			ListAccessible: &listAccessible,
		}
		resp, err = client.GetOrganizationDatabases(ctx, organization, &params)
	} else {
		// List all databases within project
		prefix = organization + "/" + project + "/"
		params := openapi.GetDatabasesParams{
			LabelFilter:    labelFilter,
			ListAccessible: &listAccessible,
		}
		resp, err = client.GetDatabases(ctx, organization, project, &params)
	}
	return processListResponse(prefix, resp, err)
}

func GetProjects(ctx context.Context, client openapi.ClientInterface, organization string, labelFilter *string, listAccessible bool) ([]string, error) {
	var prefix string
	var resp *http.Response
	var err error
	if len(organization) == 0 {
		// List all projects
		params := openapi.GetAllProjectsParams{
			LabelFilter:    labelFilter,
			ListAccessible: &listAccessible,
		}
		resp, err = client.GetAllProjects(ctx, &params)
	} else {
		// List all project within organization
		prefix = organization + "/"
		params := openapi.GetProjectsParams{
			LabelFilter:    labelFilter,
			ListAccessible: &listAccessible,
		}
		resp, err = client.GetProjects(ctx, organization, &params)
	}
	return processListResponse(prefix, resp, err)
}

func GetBackupPolicies(ctx context.Context, client openapi.ClientInterface, organization string, labelFilter *string, listAccessible bool) ([]string, error) {
	var prefix string
	var resp *http.Response
	var err error
	if len(organization) == 0 {
		// List all backup policies
		params := openapi.GetAllBackupPoliciesParams{
			LabelFilter:    labelFilter,
			ListAccessible: &listAccessible,
		}
		resp, err = client.GetAllBackupPolicies(ctx, &params)
	} else {
		// List all backup policies within organization
		prefix = organization + "/"
		params := openapi.GetBackupPoliciesParams{
			LabelFilter:    labelFilter,
			ListAccessible: &listAccessible,
		}
		resp, err = client.GetBackupPolicies(ctx, organization, &params)
	}
	return processListResponse(prefix, resp, err)
}
