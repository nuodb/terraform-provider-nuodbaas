# \ClusterDatabasequotasAPI

All URIs are relative to *https://example.nuodb.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateDatabaseQuota**](ClusterDatabasequotasAPI.md#CreateDatabaseQuota) | **Put** /cluster/databasequotas/{name} | Create or update a database quota
[**DeleteDatabaseQuota**](ClusterDatabasequotasAPI.md#DeleteDatabaseQuota) | **Delete** /cluster/databasequotas/{name} | Delete an existing database quota
[**GetDatabaseQuota**](ClusterDatabasequotasAPI.md#GetDatabaseQuota) | **Get** /cluster/databasequotas/{name} | Get an existing database quota
[**GetDatabaseQuotas**](ClusterDatabasequotasAPI.md#GetDatabaseQuotas) | **Get** /cluster/databasequotas | List the database quotas
[**PatchDatabaseQuota**](ClusterDatabasequotasAPI.md#PatchDatabaseQuota) | **Patch** /cluster/databasequotas/{name} | Update an existing database quota



## CreateDatabaseQuota

> CreateDatabaseQuota(ctx, name).DatabaseQuotaModel(databaseQuotaModel).Execute()

Create or update a database quota

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

func main() {
    name := "name_example" // string | 
    databaseQuotaModel := *openapiclient.NewDatabaseQuotaModel("Name_example", *openapiclient.NewDatabaseQuotaSpec()) // DatabaseQuotaModel | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.ClusterDatabasequotasAPI.CreateDatabaseQuota(context.Background(), name).DatabaseQuotaModel(databaseQuotaModel).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterDatabasequotasAPI.CreateDatabaseQuota``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateDatabaseQuotaRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **databaseQuotaModel** | [**DatabaseQuotaModel**](DatabaseQuotaModel.md) |  | 

### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteDatabaseQuota

> DeleteDatabaseQuota(ctx, name).TimeoutSeconds(timeoutSeconds).Execute()

Delete an existing database quota

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

func main() {
    name := "name_example" // string | 
    timeoutSeconds := int32(56) // int32 | The number of seconds to wait for the deletion to be finalized, unless 0 is specified which indicates not to wait (optional) (default to 0)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.ClusterDatabasequotasAPI.DeleteDatabaseQuota(context.Background(), name).TimeoutSeconds(timeoutSeconds).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterDatabasequotasAPI.DeleteDatabaseQuota``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteDatabaseQuotaRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **timeoutSeconds** | **int32** | The number of seconds to wait for the deletion to be finalized, unless 0 is specified which indicates not to wait | [default to 0]

### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDatabaseQuota

> DatabaseQuotaModel GetDatabaseQuota(ctx, name).Execute()

Get an existing database quota

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

func main() {
    name := "name_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ClusterDatabasequotasAPI.GetDatabaseQuota(context.Background(), name).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterDatabasequotasAPI.GetDatabaseQuota``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetDatabaseQuota`: DatabaseQuotaModel
    fmt.Fprintf(os.Stdout, "Response from `ClusterDatabasequotasAPI.GetDatabaseQuota`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDatabaseQuotaRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DatabaseQuotaModel**](DatabaseQuotaModel.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDatabaseQuotas

> ItemListString GetDatabaseQuotas(ctx).Execute()

List the database quotas

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ClusterDatabasequotasAPI.GetDatabaseQuotas(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterDatabasequotasAPI.GetDatabaseQuotas``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetDatabaseQuotas`: ItemListString
    fmt.Fprintf(os.Stdout, "Response from `ClusterDatabasequotasAPI.GetDatabaseQuotas`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetDatabaseQuotasRequest struct via the builder pattern


### Return type

[**ItemListString**](ItemListString.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PatchDatabaseQuota

> PatchDatabaseQuota(ctx, name).JsonPatchOperation(jsonPatchOperation).Execute()

Update an existing database quota

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

func main() {
    name := "name_example" // string | 
    jsonPatchOperation := []openapiclient.JsonPatchOperation{*openapiclient.NewJsonPatchOperation("Op_example", "Path_example")} // []JsonPatchOperation | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.ClusterDatabasequotasAPI.PatchDatabaseQuota(context.Background(), name).JsonPatchOperation(jsonPatchOperation).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterDatabasequotasAPI.PatchDatabaseQuota``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchDatabaseQuotaRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **jsonPatchOperation** | [**[]JsonPatchOperation**](JsonPatchOperation.md) |  | 

### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: application/json-patch+json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

