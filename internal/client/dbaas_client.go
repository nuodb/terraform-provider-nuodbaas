// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package nuodbaas_client

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/nuodb/terraform-provider-nuodbaas/openapi"
)

func WithBasicCredentials(user, password string) openapi.ClientOption {
	return openapi.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		if user != "" {
			req.SetBasicAuth(user, password)
		}
		return nil
	})
}

func NewApiClient(urlBase string, user string, password string, skipVerify bool) (*openapi.Client, error) {
	client, err := openapi.NewClient(urlBase, WithBasicCredentials(user, password))
	if err != nil {
		return nil, err
	}

	if skipVerify {
		client.Client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}
	return client, nil
}
