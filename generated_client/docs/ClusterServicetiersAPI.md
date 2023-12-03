# \ClusterServicetiersAPI

All URIs are relative to *http://}*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateServiceTier**](ClusterServicetiersAPI.md#CreateServiceTier) | **Put** /cluster/servicetiers/{name} | Create or update a service tier
[**DeleteServiceTier**](ClusterServicetiersAPI.md#DeleteServiceTier) | **Delete** /cluster/servicetiers/{name} | Delete an existing service tier
[**GetServiceTier**](ClusterServicetiersAPI.md#GetServiceTier) | **Get** /cluster/servicetiers/{name} | Get an existing service tier
[**GetServiceTiers**](ClusterServicetiersAPI.md#GetServiceTiers) | **Get** /cluster/servicetiers | List the service tiers
[**PatchServiceTier**](ClusterServicetiersAPI.md#PatchServiceTier) | **Patch** /cluster/servicetiers/{name} | Update an existing service tier



## CreateServiceTier

> CreateServiceTier(ctx, name).ServiceTierModel(serviceTierModel).Execute()

Create or update a service tier

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    name := "name_example" // string | 
    serviceTierModel := *openapiclient.NewServiceTierModel("Name_example", *openapiclient.NewServiceTierSpec()) // ServiceTierModel | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.ClusterServicetiersAPI.CreateServiceTier(context.Background(), name).ServiceTierModel(serviceTierModel).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterServicetiersAPI.CreateServiceTier``: %v\n", err)
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

Other parameters are passed through a pointer to a apiCreateServiceTierRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **serviceTierModel** | [**ServiceTierModel**](ServiceTierModel.md) |  | 

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


## DeleteServiceTier

> DeleteServiceTier(ctx, name).TimeoutSeconds(timeoutSeconds).Execute()

Delete an existing service tier

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    name := "name_example" // string | 
    timeoutSeconds := int32(56) // int32 | The number of seconds to wait for the deletion to be finalized, unless 0 is specified which indicates not to wait (optional) (default to 0)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.ClusterServicetiersAPI.DeleteServiceTier(context.Background(), name).TimeoutSeconds(timeoutSeconds).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterServicetiersAPI.DeleteServiceTier``: %v\n", err)
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

Other parameters are passed through a pointer to a apiDeleteServiceTierRequest struct via the builder pattern


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


## GetServiceTier

> ServiceTierModel GetServiceTier(ctx, name).Execute()

Get an existing service tier

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    name := "name_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ClusterServicetiersAPI.GetServiceTier(context.Background(), name).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterServicetiersAPI.GetServiceTier``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetServiceTier`: ServiceTierModel
    fmt.Fprintf(os.Stdout, "Response from `ClusterServicetiersAPI.GetServiceTier`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetServiceTierRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ServiceTierModel**](ServiceTierModel.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetServiceTiers

> ItemListString GetServiceTiers(ctx).Execute()

List the service tiers

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ClusterServicetiersAPI.GetServiceTiers(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterServicetiersAPI.GetServiceTiers``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetServiceTiers`: ItemListString
    fmt.Fprintf(os.Stdout, "Response from `ClusterServicetiersAPI.GetServiceTiers`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetServiceTiersRequest struct via the builder pattern


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


## PatchServiceTier

> PatchServiceTier(ctx, name).JsonPatchOperation(jsonPatchOperation).Execute()

Update an existing service tier

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    name := "name_example" // string | 
    jsonPatchOperation := []openapiclient.JsonPatchOperation{*openapiclient.NewJsonPatchOperation("Op_example", "Path_example")} // []JsonPatchOperation | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.ClusterServicetiersAPI.PatchServiceTier(context.Background(), name).JsonPatchOperation(jsonPatchOperation).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterServicetiersAPI.PatchServiceTier``: %v\n", err)
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

Other parameters are passed through a pointer to a apiPatchServiceTierRequest struct via the builder pattern


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

