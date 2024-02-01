# DbaasAccessRuleModel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Allow** | Pointer to **[]string** | List of access rule entries in the form &#x60;&lt;verb&gt;:&lt;resource specifier&gt;[:&lt;SLA&gt;]&#x60; that specify requests to allow | [optional] 
**Deny** | Pointer to **[]string** | List of access rule entries in the form &#x60;&lt;verb&gt;:&lt;resource specifier&gt;&#x60; that specify requests to deny | [optional] 

## Methods

### NewDbaasAccessRuleModel

`func NewDbaasAccessRuleModel() *DbaasAccessRuleModel`

NewDbaasAccessRuleModel instantiates a new DbaasAccessRuleModel object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDbaasAccessRuleModelWithDefaults

`func NewDbaasAccessRuleModelWithDefaults() *DbaasAccessRuleModel`

NewDbaasAccessRuleModelWithDefaults instantiates a new DbaasAccessRuleModel object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAllow

`func (o *DbaasAccessRuleModel) GetAllow() []string`

GetAllow returns the Allow field if non-nil, zero value otherwise.

### GetAllowOk

`func (o *DbaasAccessRuleModel) GetAllowOk() (*[]string, bool)`

GetAllowOk returns a tuple with the Allow field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllow

`func (o *DbaasAccessRuleModel) SetAllow(v []string)`

SetAllow sets Allow field to given value.

### HasAllow

`func (o *DbaasAccessRuleModel) HasAllow() bool`

HasAllow returns a boolean if a field has been set.

### GetDeny

`func (o *DbaasAccessRuleModel) GetDeny() []string`

GetDeny returns the Deny field if non-nil, zero value otherwise.

### GetDenyOk

`func (o *DbaasAccessRuleModel) GetDenyOk() (*[]string, bool)`

GetDenyOk returns a tuple with the Deny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeny

`func (o *DbaasAccessRuleModel) SetDeny(v []string)`

SetDeny sets Deny field to given value.

### HasDeny

`func (o *DbaasAccessRuleModel) HasDeny() bool`

HasDeny returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


