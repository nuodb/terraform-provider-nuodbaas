/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
package nuodbaas_client

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nuodb/terraform-provider-nuodbaas/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/model"

	nuodbaas "github.com/nuodb/terraform-provider-nuodbaas/client"
)

type NuodbaasProjectClient struct {
	client      *nuodbaas.APIClient
	org         string
	projectName string
	ctx         context.Context
}

func (client *NuodbaasProjectClient) createUpdateProject(projectModel *nuodbaas.ProjectModel, projectResourceModel model.ProjectResourceModel) error {
	apiRequestObject := client.client.ProjectsAPI.CreateProject(client.ctx, client.org, client.projectName)
	projectModel.SetSla(projectResourceModel.Sla.ValueString())
	projectModel.SetTier(projectResourceModel.Tier.ValueString())

	maintenanceModel := projectResourceModel.Maintenance

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

	return err
}

func (client *NuodbaasProjectClient) CreateProject(projectResourceModel model.ProjectResourceModel) error {
	projectModel := nuodbaas.NewProjectModelWithDefaults()
	return client.createUpdateProject(projectModel, projectResourceModel)
}

func (client *NuodbaasProjectClient) UpdateProject(projectResourceModel model.ProjectResourceModel) error {
	if len(projectResourceModel.ResourceVersion.ValueString()) == 0 {
		errorMessage := "cannot update the project. Resource version is missing"
		return errors.New(errorMessage)
	}
	projectModel := nuodbaas.NewProjectModelWithDefaults()
	projectModel.SetResourceVersion(projectResourceModel.ResourceVersion.ValueString())
	return client.createUpdateProject(projectModel, projectResourceModel)
}

func (client *NuodbaasProjectClient) GetProject() (*nuodbaas.ProjectModel, error) {
	apiGetRequestObject := client.client.ProjectsAPI.GetProject(client.ctx, client.org, client.projectName)
	projectMdoel, _, err := client.client.ProjectsAPI.GetProjectExecute(apiGetRequestObject)
	return projectMdoel, err
}

func (client *NuodbaasProjectClient) GetProjects() (*nuodbaas.ItemListString, error) {
	// TODO: Allow listing of projects across organizations
	itemList, _, err := client.client.ProjectsAPI.GetProjects(client.ctx, client.org).Execute()
	if err != nil {
		return nil, err
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

	return itemList, err

}

func (client *NuodbaasProjectClient) DeleteProject() error {
	_, err := client.client.ProjectsAPI.DeleteProject(client.ctx, client.org, client.projectName).Execute()
	return err
}

func NewProjectClient(client *nuodbaas.APIClient, ctx context.Context, org string, projectName string) *NuodbaasProjectClient {
	nuoClient := NuodbaasProjectClient{
		client:      client,
		org:         org,
		projectName: projectName,
		ctx:         ctx,
	}
	return &nuoClient
}
