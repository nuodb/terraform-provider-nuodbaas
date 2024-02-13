/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
package helper

import (
	"fmt"
	"strings"

	"github.com/nuodb/terraform-provider-nuodbaas/internal/model"
)

func GetProjectDataSourceResponse(projects []string) ([]model.ProjectDataSourceNameModel, error) {
	var ret []model.ProjectDataSourceNameModel
	for _, project := range projects {
		parts := strings.Split(project, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Unexpected format for project name: %s", project)
		}
		ret = append(ret, model.ProjectDataSourceNameModel{
			Organization: parts[0],
			Name:         parts[1],
		})
	}
	return ret, nil
}

func GetDatabaseDataSourceResponse(databases []string) ([]model.DatabaseDataSourceNameModel, error) {
	var ret []model.DatabaseDataSourceNameModel
	for _, db := range databases {
		parts := strings.Split(db, "/")
		if len(parts) != 3 {
			return nil, fmt.Errorf("Unexpected format for database name: %s", db)
		}
		ret = append(ret, model.DatabaseDataSourceNameModel{
			Organization: parts[0],
			Project:      parts[1],
			Name:         parts[2],
		})
	}
	return ret, nil
}
