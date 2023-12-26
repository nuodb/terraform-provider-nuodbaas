/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
package nuodbaas_client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/helper"
	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/model"

	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

type nuodbaasProjectClient struct {
	client		*nuodbaas.APIClient
	org    		string
	projectName string
	ctx         context.Context
}


func (client *nuodbaasProjectClient) createUpdateProject(projectModel *nuodbaas.ProjectModel, projectResourceModel model.ProjectResourceModel) (*nuodbaas.ErrorContent) {
	apiRequestObject := client.client.ProjectsAPI.CreateProject(client.ctx,client.org, client.projectName)
	projectModel.SetSla(projectResourceModel.Sla.ValueString())
	projectModel.SetTier(projectResourceModel.Tier.ValueString())

	maintenanceModel:= projectResourceModel.Maintenance

	if maintenanceModel != nil {
		var openApiMaintenanceModel = nuodbaas.MaintenanceModel{}

		if !maintenanceModel.IsDisabled.IsNull() {
			openApiMaintenanceModel.IsDisabled = maintenanceModel.IsDisabled.ValueBoolPointer()
		}
		projectModel.SetMaintenance(openApiMaintenanceModel)
	}

	if projectResourceModel.Properties != nil {
		var openApiProjectPropertiesModel = nuodbaas.NewProjectPropertiesModel()
		var projectProperties = projectResourceModel.Properties

		if !projectProperties.TierParameters.IsNull() {
			elements := projectProperties.TierParameters.Elements()
			var tierParamters = map[string]string{}
			for key, element := range elements {
				tierParamters[key] = strings.ReplaceAll(helper.RemoveDoubleQuotes(element.String()), "\\\"", "\"")
			}
			openApiProjectPropertiesModel.TierParameters = &tierParamters
		}

		projectModel.SetProperties(*openApiProjectPropertiesModel)
	}

	apiRequestObject = apiRequestObject.ProjectModel(*projectModel)
	_, err := client.client.ProjectsAPI.CreateProjectExecute(apiRequestObject)

	if serverErr, ok := err.(*nuodbaas.GenericOpenAPIError); ok {
		tflog.Debug(client.ctx, fmt.Sprintf("TAGGER err is %+v %+v %+v", serverErr.Error(), string(serverErr.Body()), serverErr.Model()))
	}	else {
		tflog.Debug(client.ctx, fmt.Sprintf("TAGGER not ok and is %+v", err))
	}
	errorModel := helper.GetErrorModelFromError(client.ctx, err)
	tflog.Debug(client.ctx, fmt.Sprintf("TAGGER converted errorMessage is %+v %+v %+v", errorModel.GetDetail(), errorModel.GetCode(), errorModel.GetStatus()))

	return errorModel
}

func (client *nuodbaasProjectClient) CreateProject(projectResourceModel model.ProjectResourceModel) *nuodbaas.ErrorContent {
	projectModel := nuodbaas.NewProjectModelWithDefaults()
	return client.createUpdateProject(projectModel, projectResourceModel)
}

func (client *nuodbaasProjectClient) UpdateProject(projectResourceModel model.ProjectResourceModel)  *nuodbaas.ErrorContent {
	if len(projectResourceModel.ResourceVersion.ValueString()) == 0 {
		errorMessage := "cannot update the project. Resource version is missing"
		return &nuodbaas.ErrorContent{ 
			Detail: &errorMessage,
		}
	}
	projectModel := nuodbaas.NewProjectModelWithDefaults()
	projectModel.SetResourceVersion(projectResourceModel.ResourceVersion.ValueString())
	return client.createUpdateProject(projectModel, projectResourceModel)
}

func (client *nuodbaasProjectClient) GetProject() (*nuodbaas.ProjectModel, *nuodbaas.ErrorContent) {
	apiGetRequestObject := client.client.ProjectsAPI.GetProject(client.ctx, client.org, client.projectName)
	projectMdoel, _, err := client.client.ProjectsAPI.GetProjectExecute(apiGetRequestObject)
	errModel := helper.GetErrorModelFromError(client.ctx, err)
	return projectMdoel, errModel
}

func (client *nuodbaasProjectClient) GetProjects() (*nuodbaas.ItemListString, *http.Response, error) {
	itemList, response, err := client.client.ProjectsAPI.GetProjects(client.ctx, client.org).Execute()
	if err != nil {
		return nil, response, err
	}
	newListItems := make([]string, 0)
	if len(client.org) > 0 {
		for _, item := range itemList.GetItems() {
			newListItems = append(newListItems, fmt.Sprintf("%s/%s", client.org, item))
		}
	}
	if len(newListItems) > 0 {
		itemList.SetItems(newListItems)
	}

	return itemList, response, err

}

func (client *nuodbaasProjectClient) DeleteProject() *nuodbaas.ErrorContent {
	_, err := client.client.ProjectsAPI.DeleteProject(client.ctx, client.org, client.projectName).Execute()
	return helper.GetErrorModelFromError(client.ctx, err)
}

func NewProjectClient(client *nuodbaas.APIClient, ctx context.Context, org string, projectName string) *nuodbaasProjectClient {
	nuoClient := nuodbaasProjectClient{
		client: 		client,
		org: 			org,
		projectName:	projectName,
		ctx:			ctx,
	}
	return &nuoClient
}