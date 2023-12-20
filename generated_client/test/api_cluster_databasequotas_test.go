/*
NuoDB Control Plane REST API

Testing ClusterDatabasequotasAPIService

*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech);

package nuodbaas

import (
	"context"
	"testing"

	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_nuodbaas_ClusterDatabasequotasAPIService(t *testing.T) {

	configuration := nuodbaas.NewConfiguration()
	apiClient := nuodbaas.NewAPIClient(configuration)

	t.Run("Test ClusterDatabasequotasAPIService CreateDatabaseQuota", func(t *testing.T) {

		t.Skip("skip test")  // remove to run test

		var name string

		httpRes, err := apiClient.ClusterDatabasequotasAPI.CreateDatabaseQuota(context.Background(), name).Execute()

		require.Nil(t, err)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test ClusterDatabasequotasAPIService DeleteDatabaseQuota", func(t *testing.T) {

		t.Skip("skip test")  // remove to run test

		var name string

		httpRes, err := apiClient.ClusterDatabasequotasAPI.DeleteDatabaseQuota(context.Background(), name).Execute()

		require.Nil(t, err)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test ClusterDatabasequotasAPIService GetDatabaseQuota", func(t *testing.T) {

		t.Skip("skip test")  // remove to run test

		var name string

		resp, httpRes, err := apiClient.ClusterDatabasequotasAPI.GetDatabaseQuota(context.Background(), name).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test ClusterDatabasequotasAPIService GetDatabaseQuotas", func(t *testing.T) {

		t.Skip("skip test")  // remove to run test

		resp, httpRes, err := apiClient.ClusterDatabasequotasAPI.GetDatabaseQuotas(context.Background()).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test ClusterDatabasequotasAPIService PatchDatabaseQuota", func(t *testing.T) {

		t.Skip("skip test")  // remove to run test

		var name string

		httpRes, err := apiClient.ClusterDatabasequotasAPI.PatchDatabaseQuota(context.Background(), name).Execute()

		require.Nil(t, err)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

}
