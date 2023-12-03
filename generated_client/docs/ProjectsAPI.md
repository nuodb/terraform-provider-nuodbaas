# \ProjectsAPI

All URIs are relative to *http://}*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateProject**](ProjectsAPI.md#CreateProject) | **Put** /projects/{organization}/{project} | Create or update a project
[**DeleteProject**](ProjectsAPI.md#DeleteProject) | **Delete** /projects/{organization}/{project} | Delete an existing project
[**GetAllProjects**](ProjectsAPI.md#GetAllProjects) | **Get** /projects | List the projects in the cluster
[**GetProject**](ProjectsAPI.md#GetProject) | **Get** /projects/{organization}/{project} | Get an existing project
[**GetProjects**](ProjectsAPI.md#GetProjects) | **Get** /projects/{organization} | List the projects in an organization
[**PatchProject**](ProjectsAPI.md#PatchProject) | **Patch** /projects/{organization}/{project} | Update an existing project



## CreateProject

> CreateProject(ctx, organization, project).ProjectModel(projectModel).Execute()

Create or update a project

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
    projectModel := *openapiclient.NewProjectModel("dev", "n0.small") // ProjectModel | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.ProjectsAPI.CreateProject(context.Background(), organization, project).ProjectModel(projectModel).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.CreateProject``: %v\n", err)
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

### Other Parameters

Other parameters are passed through a pointer to a apiCreateProjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **projectModel** | [**ProjectModel**](ProjectModel.md) |  | 

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


## DeleteProject

> DeleteProject(ctx, organization, project).TimeoutSeconds(timeoutSeconds).Execute()

Delete an existing project

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
    timeoutSeconds := int32(56) // int32 | The number of seconds to wait for the deletion to be finalized, unless 0 is specified which indicates not to wait (optional) (default to 0)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.ProjectsAPI.DeleteProject(context.Background(), organization, project).TimeoutSeconds(timeoutSeconds).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.DeleteProject``: %v\n", err)
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

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteProjectRequest struct via the builder pattern


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


## GetAllProjects

> ItemListString GetAllProjects(ctx).ListAccessible(listAccessible).Execute()

List the projects in the cluster

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
    listAccessible := true // bool | Whether to return any accessible sub-resources even if the current user does not have access privileges to list all resources at this level (optional) (default to false)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ProjectsAPI.GetAllProjects(context.Background()).ListAccessible(listAccessible).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.GetAllProjects``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetAllProjects`: ItemListString
    fmt.Fprintf(os.Stdout, "Response from `ProjectsAPI.GetAllProjects`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetAllProjectsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **listAccessible** | **bool** | Whether to return any accessible sub-resources even if the current user does not have access privileges to list all resources at this level | [default to false]

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


## GetProject

> ProjectModel GetProject(ctx, organization, project).Execute()

Get an existing project

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
    resp, r, err := apiClient.ProjectsAPI.GetProject(context.Background(), organization, project).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.GetProject``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetProject`: ProjectModel
    fmt.Fprintf(os.Stdout, "Response from `ProjectsAPI.GetProject`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organization** | **string** |  | 
**project** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetProjectRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**ProjectModel**](ProjectModel.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetProjects

> ItemListString GetProjects(ctx, organization).ListAccessible(listAccessible).Execute()

List the projects in an organization

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
    listAccessible := true // bool | Whether to return any accessible sub-resources even if the current user does not have access privileges to list all resources at this level (optional) (default to false)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ProjectsAPI.GetProjects(context.Background(), organization).ListAccessible(listAccessible).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.GetProjects``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetProjects`: ItemListString
    fmt.Fprintf(os.Stdout, "Response from `ProjectsAPI.GetProjects`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organization** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetProjectsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **listAccessible** | **bool** | Whether to return any accessible sub-resources even if the current user does not have access privileges to list all resources at this level | [default to false]

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


## PatchProject

> PatchProject(ctx, organization, project).JsonPatchOperation(jsonPatchOperation).Execute()

Update an existing project

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
    jsonPatchOperation := []openapiclient.JsonPatchOperation{*openapiclient.NewJsonPatchOperation("Op_example", "Path_example")} // []JsonPatchOperation | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.ProjectsAPI.PatchProject(context.Background(), organization, project).JsonPatchOperation(jsonPatchOperation).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ProjectsAPI.PatchProject``: %v\n", err)
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

### Other Parameters

Other parameters are passed through a pointer to a apiPatchProjectRequest struct via the builder pattern


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

