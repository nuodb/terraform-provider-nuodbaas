/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
package nuodbaas_test_client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	backoff "github.com/cenkalti/backoff/v4"

	"github.com/nuodb/terraform-provider-nuodbaas/helper"

	nuodbaas "github.com/nuodb/terraform-provider-nuodbaas/client"
	nuodbaas_client "github.com/nuodb/terraform-provider-nuodbaas/internal/client"
)

const (
	DEFAULT_URL = "http://localhost:8080"
)

func DefaultApiClient() *nuodbaas.APIClient {
	user := os.Getenv("NUODB_CP_USER")
	password := os.Getenv("NUODB_CP_PASSWORD")
	urlBase := os.Getenv("NUODB_CP_URL_BASE")

	if urlBase == "" {
		urlBase = DEFAULT_URL
	}

	return nuodbaas_client.NewApiClient(true, urlBase, user, password)
}

func CreateProject(t *testing.T, ctx context.Context, client *nuodbaas.APIClient, organization string, name string, sla string, teir string) error {
	model := nuodbaas.ProjectModel{
		Sla:  sla,
		Tier: teir,
	}

	return CreateProjectWithModel(t, ctx, client, organization, name, model)
}

func CreateProjectWithModel(t *testing.T, ctx context.Context, client *nuodbaas.APIClient, organization string, name string, model nuodbaas.ProjectModel) error {
	request := client.ProjectsAPI.CreateProject(ctx, organization, name)
	request = request.ProjectModel(model)
	_, err := request.Execute()

	t.Cleanup(func() {
		err := DeleteProject(context.TODO(), client, organization, name, true)
		if err != nil {
			t.Fatal(err)
		}
	})

	if err != nil {
		return errors.New(helper.GetApiErrorMessage(err, "Could not create project"))
	}

	return nil
}

func GetProject(ctx context.Context, client *nuodbaas.APIClient, organization string, name string) (*nuodbaas.ProjectModel, *http.Response, error) {
	project, response, err := client.ProjectsAPI.GetProject(ctx, organization, name).Execute()
	if err != nil {
		err = errors.New(helper.GetApiErrorMessage(err, "Could not get project"))
	}

	return project, response, err
}

func DeleteProject(ctx context.Context, client *nuodbaas.APIClient, organization string, name string, ignoreMissing bool) error {
	request := client.ProjectsAPI.DeleteProject(ctx, organization, name)

	response, err := request.Execute()
	if err != nil && (!ignoreMissing || response.StatusCode != http.StatusNotFound) {
		return errors.New(helper.GetApiErrorMessage(err, "Could not delete project"))
	}

	return backoff.Retry(func() error {
		_, status, err := GetProject(ctx, client, organization, name)
		if status.StatusCode == http.StatusNotFound {
			return nil
		}
		if err != nil {
			return backoff.Permanent(err)
		}

		return fmt.Errorf("Timed out waiting for project %s/%s to be deleted.", organization, name)
	}, backoff.NewExponentialBackOff())
}

func CreateDatabase(t *testing.T, ctx context.Context, client *nuodbaas.APIClient, organization string, project string, name string, password string) error {
	model := nuodbaas.DatabaseCreateUpdateModel{
		DbaPassword: &password,
	}
	return CreateDatabaseWithModel(t, ctx, client, organization, project, name, model)
}

func CreateDatabaseWithModel(t *testing.T, ctx context.Context, client *nuodbaas.APIClient, organization string, project string, name string, model nuodbaas.DatabaseCreateUpdateModel) error {
	request := client.DatabasesAPI.CreateDatabase(ctx, organization, project, name)
	request = request.DatabaseCreateUpdateModel(model)
	_, err := request.Execute()

	t.Cleanup(func() {
		err := DeleteDatabase(context.TODO(), client, organization, project, name, true)
		if err != nil {
			t.Fatal(err)
		}
	})

	if err != nil {
		return errors.New(helper.GetApiErrorMessage(err, "Could not create database"))
	}

	return nil
}

func GetDatabase(ctx context.Context, client *nuodbaas.APIClient, org string, project string, name string) (*nuodbaas.DatabaseModel, *http.Response, error) {
	database, response, err := client.DatabasesAPI.GetDatabase(ctx, org, project, name).Execute()
	if err != nil {
		err = errors.New(helper.GetApiErrorMessage(err, "Could not get database"))
	}

	return database, response, err
}

func DeleteDatabase(ctx context.Context, client *nuodbaas.APIClient, org string, project string, name string, ignoreMissing bool) error {
	request := client.DatabasesAPI.DeleteDatabase(ctx, org, project, name)

	response, err := request.Execute()
	if err != nil && (!ignoreMissing || response.StatusCode != http.StatusNotFound) {
		return errors.New(helper.GetApiErrorMessage(err, "Could not delete database"))
	}

	return backoff.Retry(func() error {
		_, status, err := GetDatabase(ctx, client, org, project, name)
		if status.StatusCode == http.StatusNotFound {
			return nil
		}
		if err != nil {
			return backoff.Permanent(err)
		}

		return fmt.Errorf("Timed out waiting for database %s/%s/%s to be deleted.", org, project, name)
	}, backoff.NewExponentialBackOff())
}

func CheckClean() error {
	var errList error = nil

	ctx := context.Background()
	client := DefaultApiClient()
	projects, err := GetProjects(ctx, client)
	errList = errors.Join(errList, err)

	if len(projects) > 0 {
		errList = errors.Join(errList, fmt.Errorf("Projects left behind: %s", projects))
	}

	databases, err := GetDatabases(ctx, client)
	errList = errors.Join(errList, err)

	if len(databases) > 0 {
		errList = errors.Join(errList, fmt.Errorf("Databases left behind: %s", databases))
	}

	return errList
}

func GetProjects(ctx context.Context, client *nuodbaas.APIClient) ([]string, error) {
	request := client.ProjectsAPI.GetAllProjects(ctx)
	request = request.ListAccessible(true)
	projects, _, err := request.Execute()
	if err != nil {
		return nil, errors.New(helper.GetApiErrorMessage(err, "Could not get projects"))
	}

	return projects.Items, nil
}

func GetDatabases(ctx context.Context, client *nuodbaas.APIClient) ([]string, error) {
	request := client.DatabasesAPI.GetAllDatabases(ctx)
	request = request.ListAccessible(true)
	projects, _, err := request.Execute()
	if err != nil {
		return nil, errors.New(helper.GetApiErrorMessage(err, "Could not get databases"))
	}

	return projects.Items, nil
}
