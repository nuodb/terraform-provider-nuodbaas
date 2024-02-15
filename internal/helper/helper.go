/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
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
	if strings.Contains(resp.Header.Get("Content-Type"), "json") && resp.StatusCode >= http.StatusBadRequest {
		apiError := ApiError{HttpResponse: resp}
		if err := json.Unmarshal(bodyBytes, &apiError.ErrorContent); err != nil {
			return err
		}
		return &apiError
	}
	if strings.Contains(resp.Header.Get("Content-Type"), "json") && dest != nil {
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return err
		}
	}
	return nil
}

func GetDatabaseByName(ctx context.Context, client *openapi.Client, organization, project, name string, dest any) error {
	resp, err := client.GetDatabase(ctx, organization, project, name)
	if err != nil {
		return err
	}
	return ParseResponse(resp, dest)
}

func DeleteDatabaseByName(ctx context.Context, client *openapi.Client, organization, project, name string) error {
	resp, err := client.DeleteDatabase(ctx, organization, project, name, nil)
	if err != nil {
		return err
	}
	return ParseResponse(resp, nil)
}

func GetDatabases(ctx context.Context, client *openapi.Client, organization, project string, listAccessible bool) ([]string, error) {
	var databases []string
	if len(organization) == 0 {
		if len(project) != 0 {
			return nil, fmt.Errorf("Cannot specify project filter (%s) without organization", project)
		}

		params := openapi.GetAllDatabasesParams{ListAccessible: &listAccessible}
		resp, err := client.GetAllDatabases(ctx, &params)
		if err != nil {
			return nil, err
		}
		var itemList openapi.ItemListString
		err = ParseResponse(resp, &itemList)
		if err != nil {
			return nil, err
		}
		databases = *itemList.Items
	} else if len(project) == 0 {
		params := openapi.GetOrganizationDatabasesParams{ListAccessible: &listAccessible}
		resp, err := client.GetOrganizationDatabases(ctx, organization, &params)
		if err != nil {
			return nil, err
		}
		var itemList openapi.ItemListString
		err = ParseResponse(resp, &itemList)
		if err != nil {
			return nil, err
		}
		for _, db := range *itemList.Items {
			databases = append(databases, organization+"/"+db)
		}
	} else {
		params := openapi.GetDatabasesParams{ListAccessible: &listAccessible}
		resp, err := client.GetDatabases(ctx, organization, project, &params)
		if err != nil {
			return nil, err
		}
		var itemList openapi.ItemListString
		err = ParseResponse(resp, &itemList)
		if err != nil {
			return nil, err
		}
		for _, db := range *itemList.Items {
			databases = append(databases, organization+"/"+project+"/"+db)
		}
	}
	return databases, nil
}

func GetProjectByName(ctx context.Context, client *openapi.Client, organization, name string, dest any) error {
	resp, err := client.GetProject(ctx, organization, name)
	if err != nil {
		return err
	}
	return ParseResponse(resp, dest)
}

func DeleteProjectByName(ctx context.Context, client *openapi.Client, organization, name string) error {
	resp, err := client.DeleteProject(ctx, organization, name, nil)
	if err != nil {
		return err
	}
	return ParseResponse(resp, nil)
}

func GetProjects(ctx context.Context, client *openapi.Client, organization string, listAccessible bool) ([]string, error) {
	var projects []string
	if len(organization) == 0 {
		params := openapi.GetAllProjectsParams{ListAccessible: &listAccessible}
		resp, err := client.GetAllProjects(ctx, &params)
		if err != nil {
			return nil, err
		}
		var itemList openapi.ItemListString
		err = ParseResponse(resp, &itemList)
		if err != nil {
			return nil, err
		}
		projects = *itemList.Items
	} else {
		params := openapi.GetProjectsParams{ListAccessible: &listAccessible}
		resp, err := client.GetProjects(ctx, organization, &params)
		if err != nil {
			return nil, err
		}
		var itemList openapi.ItemListString
		err = ParseResponse(resp, &itemList)
		if err != nil {
			return nil, err
		}
		for _, project := range *itemList.Items {
			projects = append(projects, organization+"/"+project)
		}
	}
	return projects, nil
}
