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