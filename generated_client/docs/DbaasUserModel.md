# DbaasUserModel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Organization** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **map[string]string** | User-defined labels attached to the resource that can be used for filtering | [optional] 
**AccessRule** | [**DbaasAccessRuleModel**](DbaasAccessRuleModel.md) |  | 
**ResourceVersion** | Pointer to **string** | The version of the resource. When specified in a &#x60;PUT&#x60; request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates. | [optional] 

## Methods

### NewDbaasUserModel

`func NewDbaasUserModel(accessRule DbaasAccessRuleModel, ) *DbaasUserModel`

NewDbaasUserModel instantiates a new DbaasUserModel object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDbaasUserModelWithDefaults

`func NewDbaasUserModelWithDefaults() *DbaasUserModel`

NewDbaasUserModelWithDefaults instantiates a new DbaasUserModel object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOrganization

`func (o *DbaasUserModel) GetOrganization() string`

GetOrganization returns the Organization field if non-nil, zero value otherwise.

### GetOrganizationOk

`func (o *DbaasUserModel) GetOrganizationOk() (*string, bool)`

GetOrganizationOk returns a tuple with the Organization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganization

`func (o *DbaasUserModel) SetOrganization(v string)`

SetOrganization sets Organization field to given value.

### HasOrganization

`func (o *DbaasUserModel) HasOrganization() bool`

HasOrganization returns a boolean if a field has been set.

### GetName

`func (o *DbaasUserModel) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DbaasUserModel) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DbaasUserModel) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DbaasUserModel) HasName() bool`

HasName returns a boolean if a field has been set.

### GetLabels

`func (o *DbaasUserModel) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *DbaasUserModel) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *DbaasUserModel) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *DbaasUserModel) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetAccessRule

`func (o *DbaasUserModel) GetAccessRule() DbaasAccessRuleModel`

GetAccessRule returns the AccessRule field if non-nil, zero value otherwise.

### GetAccessRuleOk

`func (o *DbaasUserModel) GetAccessRuleOk() (*DbaasAccessRuleModel, bool)`

GetAccessRuleOk returns a tuple with the AccessRule field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessRule

`func (o *DbaasUserModel) SetAccessRule(v DbaasAccessRuleModel)`

SetAccessRule sets AccessRule field to given value.


### GetResourceVersion

`func (o *DbaasUserModel) GetResourceVersion() string`

GetResourceVersion returns the ResourceVersion field if non-nil, zero value otherwise.

### GetResourceVersionOk

`func (o *DbaasUserModel) GetResourceVersionOk() (*string, bool)`

GetResourceVersionOk returns a tuple with the ResourceVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceVersion

`func (o *DbaasUserModel) SetResourceVersion(v string)`

SetResourceVersion sets ResourceVersion field to given value.

### HasResourceVersion

`func (o *DbaasUserModel) HasResourceVersion() bool`

HasResourceVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


