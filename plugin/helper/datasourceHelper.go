package helper

import (
	"strings"
	"terraform-provider-nuodbaas/internal/model"

	nuodbaas "github.com/GIT_USER_ID/GIT_REPO_ID"
	"github.com/hashicorp/terraform-plugin-framework/types"
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