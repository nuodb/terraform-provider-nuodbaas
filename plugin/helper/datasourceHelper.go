/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
package helper

import (
	"strings"

	"github.com/nuodb/nuodbaas-tf-plugin/plugin/terraform-provider-nuodbaas/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

func GetProjectDataSourceResponse(list *nuodbaas.ItemListString) []model.ProjectDataSourceResponseModel {
	var projectDataSourceList []model.ProjectDataSourceResponseModel

	for _, item := range list.GetItems() {
		splitArr := strings.Split(item, "/")
		projectDataSourceList = append(projectDataSourceList, model.ProjectDataSourceResponseModel{
			Organization: types.StringValue(splitArr[0]),
			Name:         types.StringValue(splitArr[1]),
		})
	}
	return projectDataSourceList
}

func GetDatabaseDataSourceResponse(list *nuodbaas.ItemListString) []model.DatabasesDataSourceResponseModel {
	var databaseDataSourceList []model.DatabasesDataSourceResponseModel

	for _, item := range list.GetItems() {
		splitArr := strings.Split(item, "/")
		databaseDataSourceList = append(databaseDataSourceList, model.DatabasesDataSourceResponseModel{
			Organization: types.StringValue(splitArr[0]),
			Project:         types.StringValue(splitArr[1]),
			Name:		types.StringValue(splitArr[2]),
		})
	}
	return databaseDataSourceList
}