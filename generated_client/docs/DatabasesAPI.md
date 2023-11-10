# \DatabasesAPI

All URIs are relative to *https://example.nuodb.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateDatabase**](DatabasesAPI.md#CreateDatabase) | **Put** /databases/{organization}/{project}/{database} | Create or update a database
[**DeleteDatabase**](DatabasesAPI.md#DeleteDatabase) | **Delete** /databases/{organization}/{project}/{database} | Delete an existing database
[**GetDatabase**](DatabasesAPI.md#GetDatabase) | **Get** /databases/{organization}/{project}/{database} | Get an existing database
[**GetDatabases**](DatabasesAPI.md#GetDatabases) | **Get** /databases/{organization}/{project} | List the databases in a project
[**PatchDatabase**](DatabasesAPI.md#PatchDatabase) | **Patch** /databases/{organization}/{project}/{database} | Update an existing database



## CreateDatabase

> CreateDatabase(ctx, organization, project, database).DatabaseCreateUpdateModel(databaseCreateUpdateModel).Execute()

Create or update a database

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
    organization := "organization_example" // string | 
    project := "project_example" // string | 
    database := "database_example" // string | 
    databaseCreateUpdateModel := *openapiclient.NewDatabaseCreateUpdateModel() // DatabaseCreateUpdateModel | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DatabasesAPI.CreateDatabase(context.Background(), organization, project, database).DatabaseCreateUpdateModel(databaseCreateUpdateModel).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DatabasesAPI.CreateDatabase``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organization** | **string** |  | 
**project** | **string** |  | 
**database** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateDatabaseRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **databaseCreateUpdateModel** | [**DatabaseCreateUpdateModel**](DatabaseCreateUpdateModel.md) |  | 

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


## DeleteDatabase

> DeleteDatabase(ctx, organization, project, database).Execute()

Delete an existing database

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
    organization := "organization_example" // string | 
    project := "project_example" // string | 
    database := "database_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DatabasesAPI.DeleteDatabase(context.Background(), organization, project, database).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DatabasesAPI.DeleteDatabase``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organization** | **string** |  | 
**project** | **string** |  | 
**database** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteDatabaseRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




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


## GetDatabase

> DatabaseModel GetDatabase(ctx, organization, project, database).Execute()

Get an existing database

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
    organization := "organization_example" // string | 
    project := "project_example" // string | 
    database := "database_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DatabasesAPI.GetDatabase(context.Background(), organization, project, database).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DatabasesAPI.GetDatabase``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetDatabase`: DatabaseModel
    fmt.Fprintf(os.Stdout, "Response from `DatabasesAPI.GetDatabase`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organization** | **string** |  | 
**project** | **string** |  | 
**database** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDatabaseRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




### Return type

[**DatabaseModel**](DatabaseModel.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDatabases

> ItemListString GetDatabases(ctx, organization, project).Execute()

List the databases in a project

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
    organization := "organization_example" // string | 
    project := "project_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DatabasesAPI.GetDatabases(context.Background(), organization, project).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DatabasesAPI.GetDatabases``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetDatabases`: ItemListString
    fmt.Fprintf(os.Stdout, "Response from `DatabasesAPI.GetDatabases`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organization** | **string** |  | 
**project** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetDatabasesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



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


## PatchDatabase

> PatchDatabase(ctx, organization, project, database).JsonPatchOperation(jsonPatchOperation).Execute()

Update an existing database

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
    organization := "organization_example" // string | 
    project := "project_example" // string | 
    database := "database_example" // string | 
    jsonPatchOperation := []openapiclient.JsonPatchOperation{*openapiclient.NewJsonPatchOperation("Op_example", "Path_example")} // []JsonPatchOperation | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.DatabasesAPI.PatchDatabase(context.Background(), organization, project, database).JsonPatchOperation(jsonPatchOperation).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DatabasesAPI.PatchDatabase``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organization** | **string** |  | 
**project** | **string** |  | 
**database** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchDatabaseRequest struct via the builder pattern


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

