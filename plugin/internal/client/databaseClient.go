package nuodbaas_client

import (
	"context"
	"errors"
	"fmt"
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

func (client *nuodbaasDatabaseClient) CreateDatabase(databaseResourceModel model.DatabaseResourceModel, maintenanceModel model.MaintenanceModel, propertiesResourceModel model.DatabasePropertiesResourceModel)  (*http.Response, error) {
	databaseModel := openapi.NewDatabaseCreateUpdateModel()
	return client.createDatabase(databaseModel, databaseResourceModel, maintenanceModel, propertiesResourceModel, false)
}

func (client *nuodbaasDatabaseClient) createDatabase(databaseModel *openapi.DatabaseCreateUpdateModel, databaseResourceModel model.DatabaseResourceModel,
	 maintenanceModel model.MaintenanceModel, propertiesResourceModel model.DatabasePropertiesResourceModel, isUpdate bool)  (*http.Response, error) {
	apiRequestObject := client.client.DatabasesAPI.CreateDatabase(client.ctx, client.org, client.projectName, client.databaseName)
	if isUpdate==false {
		databaseModel.SetDbaPassword(databaseResourceModel.Password.ValueString())
	}
	databaseModel.SetTier(databaseResourceModel.Tier.ValueString())
	var openApiMaintenanceModel = openapi.MaintenanceModel{}
	if !maintenanceModel.ExpiresIn.IsNull() {
		openApiMaintenanceModel.ExpiresIn = maintenanceModel.ExpiresIn.ValueStringPointer()
	}
	if !maintenanceModel.IsDisabled.IsNull() {
		openApiMaintenanceModel.IsDisabled = maintenanceModel.IsDisabled.ValueBoolPointer()
	}

	var openApiDatabasePropertiesModel = openapi.DatabasePropertiesModel{}
	if len(propertiesResourceModel.ArchiveDiskSize.ValueString()) > 0 {
		openApiDatabasePropertiesModel.ArchiveDiskSize = propertiesResourceModel.ArchiveDiskSize.ValueStringPointer()
	}
	if len(propertiesResourceModel.JournalDiskSize.ValueString()) > 0 {
		openApiDatabasePropertiesModel.JournalDiskSize = propertiesResourceModel.JournalDiskSize.ValueStringPointer()
	}

	databaseModel.SetMaintenance(openApiMaintenanceModel)
	databaseModel.SetProperties(openApiDatabasePropertiesModel)
	apiRequestObject= apiRequestObject.DatabaseCreateUpdateModel(*databaseModel)
	return client.client.DatabasesAPI.CreateDatabaseExecute(apiRequestObject)
}

func (client *nuodbaasDatabaseClient) UpdateDatabase(databaseResourceModel model.DatabaseResourceModel, maintenanceModel model.MaintenanceModel, propertiesResourceModel model.DatabasePropertiesResourceModel) (*http.Response, error) {
	if len(databaseResourceModel.ResourceVersion.ValueString()) == 0 {
		return nil, errors.New("cannot update the project. Resource version is missing")
	}
	databaseModel := openapi.NewDatabaseCreateUpdateModel()
	databaseModel.SetResourceVersion(databaseResourceModel.ResourceVersion.ValueString())
	return client.createDatabase(databaseModel, databaseResourceModel, maintenanceModel, propertiesResourceModel, true)
}

func (client *nuodbaasDatabaseClient) GetDatabase() (*openapi.DatabaseModel, *http.Response, error) {
	apiRequestObject := client.client.DatabasesAPI.GetDatabase(client.ctx, client.org, client.projectName, client.databaseName)
	return client.client.DatabasesAPI.GetDatabaseExecute(apiRequestObject)
}

func (client *nuodbaasDatabaseClient) DeleteDatabase() (*http.Response, error) {
	return client.client.ProjectsAPI.DeleteProject(client.ctx, client.org, client.projectName).Execute()
}

func (client *nuodbaasDatabaseClient) GetDatabases() (*openapi.ItemListString, *http.Response, error) {
	var (
		itemList *openapi.ItemListString
		httpResponse *http.Response
		err error
	)
	if len(client.org) == 0 && len(client.projectName) == 0 {
		itemList, httpResponse, err = client.client.DatabasesAPI.GetAllDatabases(client.ctx).Execute()
	} else {
		itemList, httpResponse, err = client.client.DatabasesAPI.GetDatabases(client.ctx, client.org, client.projectName).Execute()
	}

	if err != nil {
		return nil, httpResponse, err
	}

	newListItems := itemList.Items
	if len(client.projectName) > 0 {
		for index, item := range itemList.GetItems() {
			newListItems[index] =  fmt.Sprintf("%s/%s", client.projectName, item)
		}
	}

	if len(client.org) > 0 {
		for index, item := range itemList.GetItems() {
			newListItems[index] =  fmt.Sprintf("%s/%s", client.org, item)
		}
	}

	itemList.SetItems(newListItems)
	return itemList, httpResponse, err
	
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