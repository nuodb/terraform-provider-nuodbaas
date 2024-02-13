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

type TestClient struct {
	Client *nuodbaas.APIClient
	ctx    context.Context
}

func NewTestClient(ctx context.Context) TestClient {
	return TestClient{
		Client: DefaultApiClient(),
		ctx:    ctx,
	}
}

func DefaultApiClient() *nuodbaas.APIClient {
	user := os.Getenv("NUODB_CP_USER")
	password := os.Getenv("NUODB_CP_PASSWORD")
	urlBase := os.Getenv("NUODB_CP_URL_BASE")

	if urlBase == "" {
		urlBase = DEFAULT_URL
	}

	return nuodbaas_client.NewApiClient(true, urlBase, user, password)
}

func (client TestClient) CreateProject(t *testing.T, organization string, name string, sla string, teir string) error {
	model := nuodbaas.ProjectModel{
		Sla:  sla,
		Tier: teir,
	}

	return client.CreateProjectWithModel(t, organization, name, model)
}

func (client TestClient) CreateProjectWithModel(t *testing.T, organization string, name string, model nuodbaas.ProjectModel) error {
	request := client.Client.ProjectsAPI.CreateProject(client.ctx, organization, name)
	request = request.ProjectModel(model)
	_, err := request.Execute()

	t.Cleanup(func() {
		err := client.DeleteProject(organization, name, true)
		if err != nil {
			t.Error(err)
		}
	})

	if err != nil {
		return errors.New(helper.GetApiErrorMessage(err, "Could not create project"))
	}

	// Wait for project to exist.
	return backoff.Retry(func() error {
		_, status, err := client.GetProject(organization, name)
		if status.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Timed out waiting for project %s/%s to be created.", organization, name)
		}
		if err != nil {
			return backoff.Permanent(err)
		}

		return nil
	}, backoff.NewExponentialBackOff())
}

func (client TestClient) GetProject(organization string, name string) (*nuodbaas.ProjectModel, *http.Response, error) {
	project, response, err := client.Client.ProjectsAPI.GetProject(client.ctx, organization, name).Execute()
	if err != nil {
		err = errors.New(helper.GetApiErrorMessage(err, "Could not get project"))
	}

	return project, response, err
}

func (client TestClient) DeleteProject(organization string, name string, ignoreMissing bool) error {
	request := client.Client.ProjectsAPI.DeleteProject(client.ctx, organization, name)

	response, err := request.Execute()
	if err != nil && (!ignoreMissing || response.StatusCode != http.StatusNotFound) {
		return errors.New(helper.GetApiErrorMessage(err, "Could not delete project"))
	}

	return backoff.Retry(func() error {
		_, status, err := client.GetProject(organization, name)
		if status.StatusCode == http.StatusNotFound {
			return nil
		}
		if err != nil {
			return backoff.Permanent(err)
		}

		return fmt.Errorf("Timed out waiting for project %s/%s to be deleted.", organization, name)
	}, backoff.NewExponentialBackOff())
}

func (client TestClient) CreateDatabase(t *testing.T, organization string, project string, name string, password string) error {
	model := nuodbaas.DatabaseCreateUpdateModel{
		DbaPassword: &password,
	}
	return client.CreateDatabaseWithModel(t, organization, project, name, model)
}

func (client TestClient) CreateDatabaseWithModel(t *testing.T, organization string, project string, name string, model nuodbaas.DatabaseCreateUpdateModel) error {
	request := client.Client.DatabasesAPI.CreateDatabase(client.ctx, organization, project, name)
	request = request.DatabaseCreateUpdateModel(model)
	_, err := request.Execute()

	t.Cleanup(func() {
		err := client.DeleteDatabase(organization, project, name, true)
		if err != nil {
			t.Error(err)
		}
	})

	if err != nil {
		return errors.New(helper.GetApiErrorMessage(err, "Could not create database"))
	}

	// Wait for the database to exist.
	return backoff.Retry(func() error {
		_, status, err := client.GetDatabase(organization, project, name)
		if status.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Timed out waiting for database %s/%s/%s to be created.", organization, project, name)
		}
		if err != nil {
			return backoff.Permanent(err)
		}

		return nil
	}, backoff.NewExponentialBackOff())
}

func (client TestClient) GetDatabase(org string, project string, name string) (*nuodbaas.DatabaseModel, *http.Response, error) {
	database, response, err := client.Client.DatabasesAPI.GetDatabase(client.ctx, org, project, name).Execute()
	if err != nil {
		err = errors.New(helper.GetApiErrorMessage(err, "Could not get database"))
	}

	return database, response, err
}

func (client TestClient) DeleteDatabase(org string, project string, name string, ignoreMissing bool) error {
	request := client.Client.DatabasesAPI.DeleteDatabase(client.ctx, org, project, name)

	response, err := request.Execute()
	if err != nil && (!ignoreMissing || response.StatusCode != http.StatusNotFound) {
		return errors.New(helper.GetApiErrorMessage(err, "Could not delete database"))
	}

	return backoff.Retry(func() error {
		_, status, err := client.GetDatabase(org, project, name)
		if status.StatusCode == http.StatusNotFound {
			return nil
		}
		if err != nil {
			return backoff.Permanent(err)
		}

		return fmt.Errorf("Timed out waiting for database %s/%s/%s to be deleted.", org, project, name)
	}, backoff.NewExponentialBackOff())
}

func (client TestClient) CheckClean() error {
	var errList error = nil

	projects, err := client.GetProjects()
	errList = errors.Join(errList, err)

	if len(projects) > 0 {
		errList = errors.Join(errList, fmt.Errorf("Projects left behind: %s", projects))
	}

	databases, err := client.GetDatabases()
	errList = errors.Join(errList, err)

	if len(databases) > 0 {
		errList = errors.Join(errList, fmt.Errorf("Databases left behind: %s", databases))
	}

	return errList
}

func (client TestClient) GetProjects() ([]string, error) {
	request := client.Client.ProjectsAPI.GetAllProjects(client.ctx)
	request = request.ListAccessible(true)
	projects, _, err := request.Execute()
	if err != nil {
		return nil, errors.New(helper.GetApiErrorMessage(err, "Could not get projects"))
	}

	return projects.Items, nil
}

func (client TestClient) GetDatabases() ([]string, error) {
	request := client.Client.DatabasesAPI.GetAllDatabases(client.ctx)
	request = request.ListAccessible(true)
	projects, _, err := request.Execute()
	if err != nil {
		return nil, errors.New(helper.GetApiErrorMessage(err, "Could not get databases"))
	}

	return projects.Items, nil
}
