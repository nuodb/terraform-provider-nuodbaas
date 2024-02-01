# DbaasUserCreateUpdateModel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Organization** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **map[string]string** | User-defined labels attached to the resource that can be used for filtering | [optional] 
**Password** | Pointer to **string** | The password for the user | [optional] 
**AccessRule** | [**DbaasAccessRuleModel**](DbaasAccessRuleModel.md) |  | 
**ResourceVersion** | Pointer to **string** | The version of the resource. When specified in a &#x60;PUT&#x60; request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates. | [optional] 

## Methods

### NewDbaasUserCreateUpdateModel

`func NewDbaasUserCreateUpdateModel(accessRule DbaasAccessRuleModel, ) *DbaasUserCreateUpdateModel`

NewDbaasUserCreateUpdateModel instantiates a new DbaasUserCreateUpdateModel object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDbaasUserCreateUpdateModelWithDefaults

`func NewDbaasUserCreateUpdateModelWithDefaults() *DbaasUserCreateUpdateModel`

NewDbaasUserCreateUpdateModelWithDefaults instantiates a new DbaasUserCreateUpdateModel object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOrganization

`func (o *DbaasUserCreateUpdateModel) GetOrganization() string`

GetOrganization returns the Organization field if non-nil, zero value otherwise.

### GetOrganizationOk

`func (o *DbaasUserCreateUpdateModel) GetOrganizationOk() (*string, bool)`

GetOrganizationOk returns a tuple with the Organization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganization

`func (o *DbaasUserCreateUpdateModel) SetOrganization(v string)`

SetOrganization sets Organization field to given value.

### HasOrganization

`func (o *DbaasUserCreateUpdateModel) HasOrganization() bool`

HasOrganization returns a boolean if a field has been set.

### GetName

`func (o *DbaasUserCreateUpdateModel) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DbaasUserCreateUpdateModel) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DbaasUserCreateUpdateModel) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DbaasUserCreateUpdateModel) HasName() bool`

HasName returns a boolean if a field has been set.

### GetLabels

`func (o *DbaasUserCreateUpdateModel) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *DbaasUserCreateUpdateModel) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *DbaasUserCreateUpdateModel) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *DbaasUserCreateUpdateModel) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetPassword

`func (o *DbaasUserCreateUpdateModel) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *DbaasUserCreateUpdateModel) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *DbaasUserCreateUpdateModel) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *DbaasUserCreateUpdateModel) HasPassword() bool`

HasPassword returns a boolean if a field has been set.

### GetAccessRule

`func (o *DbaasUserCreateUpdateModel) GetAccessRule() DbaasAccessRuleModel`

GetAccessRule returns the AccessRule field if non-nil, zero value otherwise.

### GetAccessRuleOk

`func (o *DbaasUserCreateUpdateModel) GetAccessRuleOk() (*DbaasAccessRuleModel, bool)`

GetAccessRuleOk returns a tuple with the AccessRule field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessRule

`func (o *DbaasUserCreateUpdateModel) SetAccessRule(v DbaasAccessRuleModel)`

SetAccessRule sets AccessRule field to given value.


### GetResourceVersion

`func (o *DbaasUserCreateUpdateModel) GetResourceVersion() string`

GetResourceVersion returns the ResourceVersion field if non-nil, zero value otherwise.

### GetResourceVersionOk

`func (o *DbaasUserCreateUpdateModel) GetResourceVersionOk() (*string, bool)`

GetResourceVersionOk returns a tuple with the ResourceVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceVersion

`func (o *DbaasUserCreateUpdateModel) SetResourceVersion(v string)`

SetResourceVersion sets ResourceVersion field to given value.

### HasResourceVersion

`func (o *DbaasUserCreateUpdateModel) HasResourceVersion() bool`

HasResourceVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


