package nuodbaas_client

import (
	"context"
	"errors"
	"net/http"
	"terraform-provider-nuodbaas/internal/model"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"
)

type nuodbaasDatabaseClient struct {
	client			*openapi.APIClient
	org    			string
	projectName 	string
	ctx         	context.Context
	databaseName 	string
}

func (client *nuodbaasDatabaseClient) CreateDatabase(maintenanceModel model.MaintenanceModel, body model.DatabaseCreateUpdateModel) (*http.Response, error) {
	databaseModel := openapi.NewDatabaseCreateUpdateModel()
	return client.createDatabase(databaseModel, maintenanceModel, body, false)
}


func (client *nuodbaasDatabaseClient) createDatabase(databaseModel *openapi.DatabaseCreateUpdateModel,  maintenanceModel model.MaintenanceModel, body model.DatabaseCreateUpdateModel, isUpdate bool) (*http.Response, error) {
	apiRequestObject := client.client.DatabasesAPI.CreateDatabase(client.ctx, client.org, client.projectName, client.databaseName)

	if isUpdate==false {
		databaseModel.SetDbaPassword(body.Password)
	}
	
	databaseModel.SetTier(body.Tier)
	var openApiMaintenanceModel = openapi.MaintenanceModel{}
	if !maintenanceModel.ExpiresIn.IsNull() {
		openApiMaintenanceModel.ExpiresIn = maintenanceModel.ExpiresIn.ValueStringPointer()
	}
	if !maintenanceModel.IsDisabled.IsNull() {
		openApiMaintenanceModel.IsDisabled = maintenanceModel.IsDisabled.ValueBoolPointer()
	}

	var openApiDatabasePropertiesModel = openapi.DatabasePropertiesModel{}
	if len(body.ArchiveDiskSize) > 0 {
		openApiDatabasePropertiesModel.ArchiveDiskSize = &body.ArchiveDiskSize
	}
	if len(body.JournalDiskSize) > 0 {
		openApiDatabasePropertiesModel.JournalDiskSize = &body.JournalDiskSize
	}

	databaseModel.SetMaintenance(openApiMaintenanceModel)
	databaseModel.SetProperties(openApiDatabasePropertiesModel)
	apiRequestObject= apiRequestObject.DatabaseCreateUpdateModel(*databaseModel)
	return client.client.DatabasesAPI.CreateDatabaseExecute(apiRequestObject)
}

func (client *nuodbaasDatabaseClient) UpdateDatabase(maintenanceModel model.MaintenanceModel, body model.DatabaseCreateUpdateModel, resourceVersion string) (*http.Response, error) {
	if len(resourceVersion) == 0 {
		return nil, errors.New("cannot update the project. Resource version is missing")
	}
	databaseModel := openapi.NewDatabaseCreateUpdateModel()
	databaseModel.SetResourceVersion(resourceVersion)
	return client.createDatabase(databaseModel, maintenanceModel, body, true)
}

func (client *nuodbaasDatabaseClient) GetDatabase() (*openapi.DatabaseModel, *http.Response, error) {
	apiRequestObject := client.client.DatabasesAPI.GetDatabase(client.ctx, client.org, client.projectName, client.databaseName)
	return client.client.DatabasesAPI.GetDatabaseExecute(apiRequestObject)
}

func (client *nuodbaasDatabaseClient) DeleteDatabase() (*http.Response, error) {
	return client.client.ProjectsAPI.DeleteProject(client.ctx, client.org, client.projectName).Execute()
}

func NewDatabaseClient(client *openapi.APIClient, ctx context.Context, org string, projectName string, databaseName string) *nuodbaasDatabaseClient {
	databaseClient := nuodbaasDatabaseClient{
		client: 		client,
		org: 			org,
		projectName:	projectName,
		ctx:			ctx,
		databaseName: 	databaseName,
	}
	return &databaseClient
}