# \ClusterHelmfeaturesAPI

All URIs are relative to *https://example.nuodb.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateHelmFeature**](ClusterHelmfeaturesAPI.md#CreateHelmFeature) | **Put** /cluster/helmfeatures/{name} | Create or update a Helm feature
[**DeleteHelmFeature**](ClusterHelmfeaturesAPI.md#DeleteHelmFeature) | **Delete** /cluster/helmfeatures/{name} | Delete an existing Helm feature
[**GetHelmFeature**](ClusterHelmfeaturesAPI.md#GetHelmFeature) | **Get** /cluster/helmfeatures/{name} | Get an existing Helm feature
[**GetHelmFeatures**](ClusterHelmfeaturesAPI.md#GetHelmFeatures) | **Get** /cluster/helmfeatures | List the Helm features
[**PatchHelmFeature**](ClusterHelmfeaturesAPI.md#PatchHelmFeature) | **Patch** /cluster/helmfeatures/{name} | Update an existing Helm feature



## CreateHelmFeature

> CreateHelmFeature(ctx, name).HelmFeatureModel(helmFeatureModel).Execute()

Create or update a Helm feature

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
    helmFeatureModel := *openapiclient.NewHelmFeatureModel("Name_example", *openapiclient.NewHelmFeatureSpec()) // HelmFeatureModel | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.ClusterHelmfeaturesAPI.CreateHelmFeature(context.Background(), name).HelmFeatureModel(helmFeatureModel).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterHelmfeaturesAPI.CreateHelmFeature``: %v\n", err)
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

Other parameters are passed through a pointer to a apiCreateHelmFeatureRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **helmFeatureModel** | [**HelmFeatureModel**](HelmFeatureModel.md) |  | 

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


## DeleteHelmFeature

> DeleteHelmFeature(ctx, name).TimeoutSeconds(timeoutSeconds).Execute()

Delete an existing Helm feature

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
    r, err := apiClient.ClusterHelmfeaturesAPI.DeleteHelmFeature(context.Background(), name).TimeoutSeconds(timeoutSeconds).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterHelmfeaturesAPI.DeleteHelmFeature``: %v\n", err)
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

Other parameters are passed through a pointer to a apiDeleteHelmFeatureRequest struct via the builder pattern


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


## GetHelmFeature

> HelmFeatureModel GetHelmFeature(ctx, name).Execute()

Get an existing Helm feature

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
    resp, r, err := apiClient.ClusterHelmfeaturesAPI.GetHelmFeature(context.Background(), name).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterHelmfeaturesAPI.GetHelmFeature``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetHelmFeature`: HelmFeatureModel
    fmt.Fprintf(os.Stdout, "Response from `ClusterHelmfeaturesAPI.GetHelmFeature`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetHelmFeatureRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**HelmFeatureModel**](HelmFeatureModel.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetHelmFeatures

> ItemListString GetHelmFeatures(ctx).Execute()

List the Helm features

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
    resp, r, err := apiClient.ClusterHelmfeaturesAPI.GetHelmFeatures(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterHelmfeaturesAPI.GetHelmFeatures``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetHelmFeatures`: ItemListString
    fmt.Fprintf(os.Stdout, "Response from `ClusterHelmfeaturesAPI.GetHelmFeatures`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetHelmFeaturesRequest struct via the builder pattern


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


## PatchHelmFeature

> PatchHelmFeature(ctx, name).JsonPatchOperation(jsonPatchOperation).Execute()

Update an existing Helm feature

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
    r, err := apiClient.ClusterHelmfeaturesAPI.PatchHelmFeature(context.Background(), name).JsonPatchOperation(jsonPatchOperation).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterHelmfeaturesAPI.PatchHelmFeature``: %v\n", err)
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

Other parameters are passed through a pointer to a apiPatchHelmFeatureRequest struct via the builder pattern


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

