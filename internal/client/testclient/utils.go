/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
package nuodbaas_test_client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	backoff "github.com/cenkalti/backoff/v4"

	nuodbaas_client "github.com/nuodb/terraform-provider-nuodbaas/internal/client"
	"github.com/nuodb/terraform-provider-nuodbaas/internal/helper"
	"github.com/nuodb/terraform-provider-nuodbaas/openapi"
)

const (
	DEFAULT_URL = "http://localhost:8080"
)

func DefaultApiClient() (*openapi.Client, error) {
	user := os.Getenv("NUODB_CP_USER")
	password := os.Getenv("NUODB_CP_PASSWORD")
	urlBase := os.Getenv("NUODB_CP_URL_BASE")

	if urlBase == "" {
		urlBase = DEFAULT_URL
	}

	return nuodbaas_client.NewApiClient(urlBase, user, password, true)
}

func CreateProject(t *testing.T, ctx context.Context, client *openapi.Client, organization string, name string, sla string, tier string) error {
	model := openapi.ProjectModel{
		Sla:  sla,
		Tier: tier,
	}
	return CreateProjectWithModel(t, ctx, client, organization, name, model)
}

func CreateProjectWithModel(t *testing.T, ctx context.Context, client *openapi.Client, organization string, name string, model openapi.ProjectModel) error {
	// Set name fields explicitly, since JSON serialization does not omit them
	if model.Organization == "" {
		model.Organization = organization
	}
	if model.Name == "" {
		model.Name = name
	}
	resp, err := client.CreateProject(ctx, organization, name, model)
	t.Cleanup(func() {
		err := DeleteProject(context.TODO(), client, organization, name, true)
		if err != nil {
			t.Error(err)
		}
	})

	if err != nil {
		return err
	}
	err = helper.ParseResponse(resp, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetProject(ctx context.Context, client *openapi.Client, organization string, name string) (*openapi.ProjectModel, error) {
	var dest openapi.ProjectModel
	return &dest, helper.GetProjectByName(ctx, client, organization, name, &dest)
}

func DeleteProject(ctx context.Context, client *openapi.Client, organization string, name string, ignoreMissing bool) error {
	if err := helper.DeleteProjectByName(ctx, client, organization, name); err != nil {
		if ignoreMissing && helper.IsNotFound(err) {
			return nil
		}
		return err
	}

	return backoff.Retry(func() error {
		err := helper.GetProjectByName(ctx, client, organization, name, nil)
		if helper.IsNotFound(err) {
			return nil
		}
		if err != nil {
			return backoff.Permanent(err)
		}

		return fmt.Errorf("Timed out waiting for project %s/%s to be deleted.", organization, name)
	}, backoff.NewExponentialBackOff())
}

func CreateDatabase(t *testing.T, ctx context.Context, client *openapi.Client, organization string, project string, name string, password string) error {
	model := openapi.DatabaseCreateUpdateModel{
		DbaPassword: &password,
	}
	return CreateDatabaseWithModel(t, ctx, client, organization, project, name, model)
}

func CreateDatabaseWithModel(t *testing.T, ctx context.Context, client *openapi.Client, organization string, project string, name string, model openapi.DatabaseCreateUpdateModel) error {
	// Set name fields explicitly, since JSON serialization does not omit them
	if model.Organization == "" {
		model.Organization = organization
	}
	if model.Project == "" {
		model.Project = project
	}
	if model.Name == "" {
		model.Name = name
	}
	resp, err := client.CreateDatabase(ctx, organization, project, name, model)
	t.Cleanup(func() {
		err := DeleteDatabase(context.TODO(), client, organization, project, name, true)
		if err != nil {
			t.Error(err)
		}
	})

	if err != nil {
		return err
	}
	err = helper.ParseResponse(resp, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetDatabase(ctx context.Context, client *openapi.Client, organization, project, name string) (*openapi.DatabaseModel, error) {
	var dest openapi.DatabaseModel
	return &dest, helper.GetDatabaseByName(ctx, client, organization, project, name, &dest)
}

func DeleteDatabase(ctx context.Context, client *openapi.Client, organization, project, name string, ignoreMissing bool) error {
	if err := helper.DeleteDatabaseByName(ctx, client, organization, project, name); err != nil {
		if ignoreMissing && helper.IsNotFound(err) {
			return nil
		}
		return err
	}

	return backoff.Retry(func() error {
		_, err := GetDatabase(ctx, client, organization, project, name)
		if helper.IsNotFound(err) {
			return nil
		}
		if err != nil {
			return backoff.Permanent(err)
		}

		return fmt.Errorf("Timed out waiting for database %s/%s/%s to be deleted.", organization, project, name)
	}, backoff.NewExponentialBackOff())
}

func CheckClean() error {
	var errList error = nil

	ctx := context.Background()
	client, err := DefaultApiClient()
	if err != nil {
		return err
	}
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

func GetProjects(ctx context.Context, client *openapi.Client) ([]string, error) {
	return helper.GetProjects(ctx, client, "", true)
}

func GetDatabases(ctx context.Context, client *openapi.Client) ([]string, error) {
	return helper.GetDatabases(ctx, client, "", "", true)
}
