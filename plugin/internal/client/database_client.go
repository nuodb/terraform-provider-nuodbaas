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

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/helper"
	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/model"

	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

type NuodbaasDatabaseClient struct {
	client			*nuodbaas.APIClient
	org    			string
	projectName 	string
	ctx         	context.Context
	databaseName 	string
}

func (client *NuodbaasDatabaseClient) CreateDatabase(databaseResourceModel model.DatabaseResourceModel, propertiesResourceModel *model.DatabasePropertiesResourceModel) error {
	databaseModel := nuodbaas.NewDatabaseCreateUpdateModel()
	return client.createDatabase(databaseModel, databaseResourceModel, propertiesResourceModel)
}

func (client *NuodbaasDatabaseClient) createDatabase(databaseModel *nuodbaas.DatabaseCreateUpdateModel, databaseResourceModel model.DatabaseResourceModel,
	 propertiesResourceModel *model.DatabasePropertiesResourceModel) error {
	apiRequestObject := client.client.DatabasesAPI.CreateDatabase(client.ctx, client.org, client.projectName, client.databaseName)
	if databaseModel.ResourceVersion == nil {
		databaseModel.SetDbaPassword(databaseResourceModel.Password.ValueString())
	}
	databaseModel.SetTier(databaseResourceModel.Tier.ValueString())
	maintenanceModel := databaseResourceModel.Maintenance
	if maintenanceModel != nil {
		var openApiMaintenanceModel = nuodbaas.MaintenanceModel{}
		if !maintenanceModel.IsDisabled.IsNull() {
			openApiMaintenanceModel.IsDisabled = maintenanceModel.IsDisabled.ValueBoolPointer()
		}
		databaseModel.SetMaintenance(openApiMaintenanceModel)
	}
	

	var openApiDatabasePropertiesModel = nuodbaas.DatabasePropertiesModel{}
	if propertiesResourceModel != nil {
		if len(propertiesResourceModel.ArchiveDiskSize.ValueString()) > 0 {
			openApiDatabasePropertiesModel.ArchiveDiskSize = propertiesResourceModel.ArchiveDiskSize.ValueStringPointer()
		}
		if len(propertiesResourceModel.JournalDiskSize.ValueString()) > 0 {
			openApiDatabasePropertiesModel.JournalDiskSize = propertiesResourceModel.JournalDiskSize.ValueStringPointer()
		}
		if !propertiesResourceModel.TierParameters.IsNull() {
			elements := propertiesResourceModel.TierParameters.Elements()
			var tierParamters = map[string]string{}
			for key, element := range elements {
				tierParamters[key] = strings.ReplaceAll(helper.RemoveDoubleQuotes(element.String()), "\\\"", "\"")
			}
			openApiDatabasePropertiesModel.TierParameters = &tierParamters
		}
	}
	
	databaseModel.SetProperties(openApiDatabasePropertiesModel)
	apiRequestObject= apiRequestObject.DatabaseCreateUpdateModel(*databaseModel)
	_, err := client.client.DatabasesAPI.CreateDatabaseExecute(apiRequestObject)
	return err
}

func (client *NuodbaasDatabaseClient) UpdateDatabase(databaseResourceModel model.DatabaseResourceModel, propertiesResourceModel *model.DatabasePropertiesResourceModel) error {
	if len(databaseResourceModel.ResourceVersion.ValueString()) == 0 {
		return errors.New("cannot update the project. Resource version is missing")
	}
	databaseModel := nuodbaas.NewDatabaseCreateUpdateModel()
	databaseModel.SetResourceVersion(databaseResourceModel.ResourceVersion.ValueString())
	return client.createDatabase(databaseModel, databaseResourceModel, propertiesResourceModel)
}

func (client *NuodbaasDatabaseClient) GetDatabase() (*nuodbaas.DatabaseModel, error) {
	apiRequestObject := client.client.DatabasesAPI.GetDatabase(client.ctx, client.org, client.projectName, client.databaseName)
	model, _, err:= client.client.DatabasesAPI.GetDatabaseExecute(apiRequestObject)
	return model, err
}

func (client *NuodbaasDatabaseClient) DeleteDatabase() error {
	
	_, err:= client.client.DatabasesAPI.DeleteDatabase(client.ctx, client.org, client.projectName, client.databaseName).Execute()
	return err
}

func (client *NuodbaasDatabaseClient) GetDatabases() (*nuodbaas.ItemListString, error) {
	var (
		itemList *nuodbaas.ItemListString
		err error
	)
	if len(client.org) == 0 && len(client.projectName) == 0 {
		itemList, _, err = client.client.DatabasesAPI.GetAllDatabases(client.ctx).ListAccessible(true).Execute()
	} else {
		itemList, _, err = client.client.DatabasesAPI.GetDatabases(client.ctx, client.org, client.projectName).Execute()
	}

	if err != nil {
		return nil, err
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
	return itemList, err
	
}

func NewDatabaseClient(client *nuodbaas.APIClient, ctx context.Context, org string, projectName string, databaseName string) *NuodbaasDatabaseClient {
	databaseClient := NuodbaasDatabaseClient{
		client: 		client,
		org: 			org,
		projectName:	projectName,
		ctx:			ctx,
		databaseName: 	databaseName,
	}
	return &databaseClient
}