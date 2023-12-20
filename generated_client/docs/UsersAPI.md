# \UsersAPI

All URIs are relative to *https://example.nuodb.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateUser**](UsersAPI.md#CreateUser) | **Put** /users/{organization}/{user} | Create or update a user
[**DeleteUser**](UsersAPI.md#DeleteUser) | **Delete** /users/{organization}/{user} | Delete an existing user
[**GetAllUsers**](UsersAPI.md#GetAllUsers) | **Get** /users | List the users in the cluster
[**GetUser**](UsersAPI.md#GetUser) | **Get** /users/{organization}/{user} | Get an existing user
[**GetUsers**](UsersAPI.md#GetUsers) | **Get** /users/{organization} | List the users in an organization
[**PatchUser**](UsersAPI.md#PatchUser) | **Patch** /users/{organization}/{user} | Update an existing user



## CreateUser

> CreateUser(ctx, organization, user).DbaasUserCreateUpdateModel(dbaasUserCreateUpdateModel).AllowCrossOrganizationAccess(allowCrossOrganizationAccess).Execute()

Create or update a user

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
    user := "user_example" // string | 
    dbaasUserCreateUpdateModel := *openapiclient.NewDbaasUserCreateUpdateModel(*openapiclient.NewDbaasAccessRuleModel()) // DbaasUserCreateUpdateModel | 
    allowCrossOrganizationAccess := true // bool | Whether to allow the user to have access outside of its organization (optional) (default to false)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.UsersAPI.CreateUser(context.Background(), organization, user).DbaasUserCreateUpdateModel(dbaasUserCreateUpdateModel).AllowCrossOrganizationAccess(allowCrossOrganizationAccess).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UsersAPI.CreateUser``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organization** | **string** |  | 
**user** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiCreateUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **dbaasUserCreateUpdateModel** | [**DbaasUserCreateUpdateModel**](DbaasUserCreateUpdateModel.md) |  | 
 **allowCrossOrganizationAccess** | **bool** | Whether to allow the user to have access outside of its organization | [default to false]

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


## DeleteUser

> DeleteUser(ctx, organization, user).TimeoutSeconds(timeoutSeconds).Execute()

Delete an existing user

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
    user := "user_example" // string | 
    timeoutSeconds := int32(56) // int32 | The number of seconds to wait for the deletion to be finalized, unless 0 is specified which indicates not to wait (optional) (default to 0)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.UsersAPI.DeleteUser(context.Background(), organization, user).TimeoutSeconds(timeoutSeconds).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UsersAPI.DeleteUser``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organization** | **string** |  | 
**user** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteUserRequest struct via the builder pattern


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


## GetAllUsers

> ItemListString GetAllUsers(ctx).LabelFilter(labelFilter).ListAccessible(listAccessible).Execute()

List the users in the cluster

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
    labelFilter := "labelFilter_example" // string | Comma-separated list of filters to apply based on labels, which are composed using `AND`. Acceptable filter expressions are: * `key` - Only return resources that have label with specified key * `key=value` - Only return resources that have label with specified key set to value * `!key` - Only return resources that do _not_ have label with specified key * `key!=value` - Only return resources that do _not_ have label with specified key set to value (optional)
    listAccessible := true // bool | Whether to return any accessible sub-resources even if the current user does not have access privileges to list all resources at this level (optional) (default to false)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.UsersAPI.GetAllUsers(context.Background()).LabelFilter(labelFilter).ListAccessible(listAccessible).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UsersAPI.GetAllUsers``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetAllUsers`: ItemListString
    fmt.Fprintf(os.Stdout, "Response from `UsersAPI.GetAllUsers`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetAllUsersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **labelFilter** | **string** | Comma-separated list of filters to apply based on labels, which are composed using &#x60;AND&#x60;. Acceptable filter expressions are: * &#x60;key&#x60; - Only return resources that have label with specified key * &#x60;key&#x3D;value&#x60; - Only return resources that have label with specified key set to value * &#x60;!key&#x60; - Only return resources that do _not_ have label with specified key * &#x60;key!&#x3D;value&#x60; - Only return resources that do _not_ have label with specified key set to value | 
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


## GetUser

> DbaasUserModel GetUser(ctx, organization, user).Execute()

Get an existing user

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
    user := "user_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.UsersAPI.GetUser(context.Background(), organization, user).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UsersAPI.GetUser``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetUser`: DbaasUserModel
    fmt.Fprintf(os.Stdout, "Response from `UsersAPI.GetUser`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organization** | **string** |  | 
**user** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



### Return type

[**DbaasUserModel**](DbaasUserModel.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetUsers

> ItemListString GetUsers(ctx, organization).LabelFilter(labelFilter).ListAccessible(listAccessible).Execute()

List the users in an organization

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
    labelFilter := "labelFilter_example" // string | Comma-separated list of filters to apply based on labels, which are composed using `AND`. Acceptable filter expressions are: * `key` - Only return resources that have label with specified key * `key=value` - Only return resources that have label with specified key set to value * `!key` - Only return resources that do _not_ have label with specified key * `key!=value` - Only return resources that do _not_ have label with specified key set to value (optional)
    listAccessible := true // bool | Whether to return any accessible sub-resources even if the current user does not have access privileges to list all resources at this level (optional) (default to false)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.UsersAPI.GetUsers(context.Background(), organization).LabelFilter(labelFilter).ListAccessible(listAccessible).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UsersAPI.GetUsers``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetUsers`: ItemListString
    fmt.Fprintf(os.Stdout, "Response from `UsersAPI.GetUsers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organization** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetUsersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **labelFilter** | **string** | Comma-separated list of filters to apply based on labels, which are composed using &#x60;AND&#x60;. Acceptable filter expressions are: * &#x60;key&#x60; - Only return resources that have label with specified key * &#x60;key&#x3D;value&#x60; - Only return resources that have label with specified key set to value * &#x60;!key&#x60; - Only return resources that do _not_ have label with specified key * &#x60;key!&#x3D;value&#x60; - Only return resources that do _not_ have label with specified key set to value | 
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


## PatchUser

> PatchUser(ctx, organization, user).JsonPatchOperation(jsonPatchOperation).AllowCrossOrganizationAccess(allowCrossOrganizationAccess).Execute()

Update an existing user

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
    user := "user_example" // string | 
    jsonPatchOperation := []openapiclient.JsonPatchOperation{*openapiclient.NewJsonPatchOperation("Op_example", "Path_example")} // []JsonPatchOperation | 
    allowCrossOrganizationAccess := true // bool | Whether to allow the user to have access outside of its organization (optional) (default to false)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.UsersAPI.PatchUser(context.Background(), organization, user).JsonPatchOperation(jsonPatchOperation).AllowCrossOrganizationAccess(allowCrossOrganizationAccess).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UsersAPI.PatchUser``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**organization** | **string** |  | 
**user** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiPatchUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **jsonPatchOperation** | [**[]JsonPatchOperation**](JsonPatchOperation.md) |  | 
 **allowCrossOrganizationAccess** | **bool** | Whether to allow the user to have access outside of its organization | [default to false]

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

