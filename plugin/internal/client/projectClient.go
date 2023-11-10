package nuodbaas_client

import (
	"context"
	"errors"
	"net/http"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/model"

	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

type nuodbaasProjectClient struct {
	client		*nuodbaas.APIClient
	org    		string
	projectName string
	ctx         context.Context
}


func (client *nuodbaasProjectClient) createUpdateProject(projectModel *nuodbaas.ProjectModel, projectResourceModel model.ProjectResourceModel, maintenanceModel model.MaintenanceModel) (*http.Response, error) {
	apiRequestObject := client.client.ProjectsAPI.CreateProject(client.ctx,client.org, client.projectName)
	projectModel.SetSla(projectResourceModel.Sla.ValueString())
	projectModel.SetTier(projectResourceModel.Tier.ValueString())
	var openApiMaintenanceModel = nuodbaas.MaintenanceModel{}
	if !maintenanceModel.IsDisabled.IsNull() {
		openApiMaintenanceModel.IsDisabled = maintenanceModel.IsDisabled.ValueBoolPointer()
	}
	projectModel.SetMaintenance(openApiMaintenanceModel)
	apiRequestObject = apiRequestObject.ProjectModel(*projectModel)
	return client.client.ProjectsAPI.CreateProjectExecute(apiRequestObject)
}

func (client *nuodbaasProjectClient) CreateProject(projectResourceModel model.ProjectResourceModel,  maintenanceModel model.MaintenanceModel) (*http.Response, error) {
	projectModel := nuodbaas.NewProjectModelWithDefaults()
	return client.createUpdateProject(projectModel, projectResourceModel, maintenanceModel)
}

func (client *nuodbaasProjectClient) UpdateProject(projectResourceModel model.ProjectResourceModel,  maintenanceModel model.MaintenanceModel) (*http.Response, error) {
	if len(projectResourceModel.ResourceVersion.ValueString()) == 0 {
		return nil, errors.New("Cannot update the project. Resource version is missing")
	}
	projectModel := nuodbaas.NewProjectModelWithDefaults()
	projectModel.SetResourceVersion(projectResourceModel.ResourceVersion.ValueString())
	return client.createUpdateProject(projectModel, projectResourceModel, maintenanceModel)
}

func (client *nuodbaasProjectClient) GetProject() (*nuodbaas.ProjectModel, *http.Response, error) {
	apiGetRequestObject := client.client.ProjectsAPI.GetProject(client.ctx, client.org, client.projectName)
	return client.client.ProjectsAPI.GetProjectExecute(apiGetRequestObject)
}

func (client *nuodbaasProjectClient) DeleteProject() (*http.Response, error) {
	return client.client.ProjectsAPI.DeleteProject(client.ctx, client.org, client.projectName).Execute()
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