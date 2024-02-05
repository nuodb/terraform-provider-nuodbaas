/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	nuodbaas "github.com/nuodb/terraform-provider-nuodbaas/client"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/model"
)

// Return a generic error message if the provided provider attribute is missing
func GetProviderValidatorErrorMessage(valueType string, envVariable string) string {
	return fmt.Sprintf("The provider cannot create the NuoDbaas API client as there is a missing or empty value for the NuoDbaas API %v. "+
		"Set the %v value in the configuration or use the %v environment variable. "+
		"If either is already set, ensure the value is not empty.", valueType, valueType, envVariable)
}

// Extracts the error message from error object
func GetApiErrorMessage(err error, message string) string {
	errorObj := GetErrorContentObj(err)
	extendedErrorMessage := err.Error()
	if errorObj != nil {
		extendedErrorMessage = errorObj.GetDetail()
	}
	return fmt.Sprintf("%s %s", message, extendedErrorMessage)
}

// Return a Error Content object
func GetErrorContentObj(err error) *nuodbaas.ErrorContent {
	if serverErr, ok := err.(*nuodbaas.GenericOpenAPIError); ok {
		if errModel, ok := serverErr.Model().(nuodbaas.ErrorContent); ok {
			return &errModel
		}
	}
	return nil
}

func IsTimeoutError(err error) bool {
	return os.IsTimeout(err)
}

// Removes any extra double quotes that are added in the string
func RemoveDoubleQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

// Converts map[string]string to basetypes.MapValue
func ConvertMapToTfMap(mapObj *map[string]string) (basetypes.MapValue, diag.Diagnostics) {
	mapValue := map[string]attr.Value{}
	for k, v := range *mapObj {
		mapValue[k] = types.StringValue(v)
	}
	tfMapValue, diags := types.MapValue(types.StringType, mapValue)
	return tfMapValue, diags
}

func ReadResource(ctx context.Context, diags diag.Diagnostics, fn func(context.Context, any) diag.Diagnostics, target any) bool {
	var obj types.Object
	diags.Append(fn(ctx, &obj)...)
	if diags.HasError() {
		return false
	}
	diags.Append(obj.As(ctx, target, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})...)
	return !diags.HasError()
}

func ConvertResource(diags diag.Diagnostics, from any, to any) bool {
	err := Convert(from, to)
	if err != nil {
		diags.AddError(
			"Conversion Failed",
			fmt.Sprintf("Unable to convert from %+T to %+T resource: %s", from, to, err.Error()))
		return false
	}
	return true
}

func Convert(from any, to any) error {
	var bytes bytes.Buffer
	if err := json.NewEncoder(&bytes).Encode(from); err != nil {
		return err
	}
	if err := json.NewDecoder(&bytes).Decode(to); err != nil {
		return err
	}
	return nil
}

func CreateDatabase(ctx context.Context, client *nuodbaas.APIClient, resource model.DatabaseResourceModel) error {
	createModel := nuodbaas.DatabaseCreateUpdateModel{}
	if err := Convert(&resource, &createModel); err != nil {
		return err
	}
	_, err := client.DatabasesAPI.CreateDatabase(ctx, resource.Organization, resource.Project, resource.Name).DatabaseCreateUpdateModel(createModel).Execute()
	return err
}

func UpdateDatabase(ctx context.Context, client *nuodbaas.APIClient, resource model.DatabaseResourceModel) error {
	// Fetch database and get resourceVersion
	latest, err := GetDatabase(ctx, client, resource.Organization, resource.Project, resource.Name)
	if err != nil {
		return err
	}
	updateModel := nuodbaas.DatabaseCreateUpdateModel{}
	if err := Convert(&resource, &updateModel); err != nil {
		return err
	}
	updateModel.DbaPassword = nil
	for {
		updateModel.ResourceVersion = latest.ResourceVersion
		_, err := client.DatabasesAPI.CreateDatabase(ctx, resource.Organization, resource.Project, resource.Name).DatabaseCreateUpdateModel(updateModel).Execute()
		if err == nil {
			return nil
		}
		errorContent := GetErrorContentObj(err)
		if errorContent == nil || errorContent.Code == nil || *errorContent.Code != "CONCURRENT_UPDATE" {
			return err
		}
		// Re-fetch database and get resourceVersion
		latest, err = GetDatabase(ctx, client, resource.Organization, resource.Project, resource.Name)
	}
}

func GetDatabase(ctx context.Context, client *nuodbaas.APIClient, organization, project, database string) (*nuodbaas.DatabaseModel, error) {
	model, _, err := client.DatabasesAPI.GetDatabase(ctx, organization, project, database).Execute()
	return model, err
}

func DeleteDatabase(ctx context.Context, client *nuodbaas.APIClient, organization, project, database string) error {
	_, err := client.DatabasesAPI.DeleteDatabase(ctx, organization, project, database).Execute()
	return err
}

func GetDatabases(ctx context.Context, client *nuodbaas.APIClient, organization, project string, listAccessible bool) ([]string, error) {
	var databases []string
	if len(organization) == 0 {
		if len(project) != 0 {
			return nil, fmt.Errorf("Cannot specify project filter (%s) without organization", project)
		}

		itemList, _, err := client.DatabasesAPI.GetAllDatabases(ctx).ListAccessible(listAccessible).Execute()
		if err != nil {
			return nil, err
		}
		for _, db := range itemList.Items {
			databases = append(databases, db)
		}
	} else if len(project) == 0 {
		itemList, _, err := client.DatabasesAPI.GetOrganizationDatabases(ctx, organization).ListAccessible(listAccessible).Execute()
		if err != nil {
			return nil, err
		}
		for _, db := range itemList.Items {
			databases = append(databases, organization+"/"+db)
		}
	} else {
		itemList, _, err := client.DatabasesAPI.GetDatabases(ctx, organization, project).ListAccessible(listAccessible).Execute()
		if err != nil {
			return nil, err
		}
		for _, db := range itemList.Items {
			databases = append(databases, organization+"/"+project+"/"+db)
		}
	}
	return databases, nil
}

func CreateProject(ctx context.Context, client *nuodbaas.APIClient, resource model.ProjectResourceModel) error {
	createModel := nuodbaas.ProjectModel{}
	if err := Convert(&resource, &createModel); err != nil {
		return err
	}
	_, err := client.ProjectsAPI.CreateProject(ctx, resource.Organization, resource.Name).ProjectModel(createModel).Execute()
	return err
}

func UpdateProject(ctx context.Context, client *nuodbaas.APIClient, resource model.ProjectResourceModel) error {
	// Fetch project and get resourceVersion
	latest, err := GetProject(ctx, client, resource.Organization, resource.Name)
	if err != nil {
		return err
	}
	updateModel := nuodbaas.ProjectModel{}
	if err := Convert(&resource, &updateModel); err != nil {
		return err
	}
	for {
		updateModel.ResourceVersion = latest.ResourceVersion
		_, err := client.ProjectsAPI.CreateProject(ctx, resource.Organization, resource.Name).ProjectModel(updateModel).Execute()
		if err == nil {
			return nil
		}
		errorContent := GetErrorContentObj(err)
		if errorContent == nil || errorContent.Code == nil || *errorContent.Code != "CONCURRENT_UPDATE" {
			return err
		}
		// Re-fetch project and get resourceVersion
		latest, err = GetProject(ctx, client, resource.Organization, resource.Name)
	}
}

func GetProject(ctx context.Context, client *nuodbaas.APIClient, organization, project string) (*nuodbaas.ProjectModel, error) {
	model, _, err := client.ProjectsAPI.GetProject(ctx, organization, project).Execute()
	return model, err
}

func DeleteProject(ctx context.Context, client *nuodbaas.APIClient, organization, project string) error {
	_, err := client.ProjectsAPI.DeleteProject(ctx, organization, project).Execute()
	return err
}

func GetProjects(ctx context.Context, client *nuodbaas.APIClient, organization string, listAccessible bool) ([]string, error) {
	var projects []string
	if len(organization) == 0 {
		itemList, _, err := client.ProjectsAPI.GetAllProjects(ctx).ListAccessible(listAccessible).Execute()
		if err != nil {
			return nil, err
		}
		for _, project := range itemList.Items {
			projects = append(projects, project)
		}
	} else {
		itemList, _, err := client.ProjectsAPI.GetProjects(ctx, organization).ListAccessible(listAccessible).Execute()
		if err != nil {
			return nil, err
		}
		for _, project := range itemList.Items {
			projects = append(projects, organization+"/"+project)
		}
	}
	return projects, nil
}
