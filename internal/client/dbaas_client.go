/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
package nuodbaas_client

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"

	nuodbaas "github.com/nuodb/terraform-provider-nuodbaas/client"
)

func NewApiClient(skipVerify bool, urlBase string, user string, password string) *nuodbaas.APIClient {
	configuration := nuodbaas.NewConfiguration()

	if skipVerify {
		configuration.HTTPClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}
	configuration.Servers = nuodbaas.ServerConfigurations{
		{URL: urlBase, Description: "The base URL to use for the Terraform provider"},
	}
	if user != "" {
		configuration.DefaultHeader["Authorization"] = fmt.Sprintf("Basic %s", basicAuth(user, password))
	}

	apiClient := nuodbaas.NewAPIClient(configuration)
	return apiClient
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
