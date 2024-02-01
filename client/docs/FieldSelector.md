# FieldSelector

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MatchExpressions** | Pointer to [**[]DatabasequotaspecScopeFieldselectorMatchExpressions**](DatabasequotaspecScopeFieldselectorMatchExpressions.md) | The list of field selector requirements, which are composed with &#x60;AND&#x60;. | [optional] 
**MatchFields** | Pointer to **map[string]string** | The field selector requirements as a map where each key-value pair is equivalent to an element of &#x60;matchExpressions&#x60; with &#x60;operator&#x60; set to &#x60;&#x3D;&#x3D;&#x60;. The requirements are composed with &#x60;AND&#x60;. | [optional] 

## Methods

### NewFieldSelector

`func NewFieldSelector() *FieldSelector`

NewFieldSelector instantiates a new FieldSelector object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewFieldSelectorWithDefaults

`func NewFieldSelectorWithDefaults() *FieldSelector`

NewFieldSelectorWithDefaults instantiates a new FieldSelector object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMatchExpressions

`func (o *FieldSelector) GetMatchExpressions() []DatabasequotaspecScopeFieldselectorMatchExpressions`

GetMatchExpressions returns the MatchExpressions field if non-nil, zero value otherwise.

### GetMatchExpressionsOk

`func (o *FieldSelector) GetMatchExpressionsOk() (*[]DatabasequotaspecScopeFieldselectorMatchExpressions, bool)`

GetMatchExpressionsOk returns a tuple with the MatchExpressions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchExpressions

`func (o *FieldSelector) SetMatchExpressions(v []DatabasequotaspecScopeFieldselectorMatchExpressions)`

SetMatchExpressions sets MatchExpressions field to given value.

### HasMatchExpressions

`func (o *FieldSelector) HasMatchExpressions() bool`

HasMatchExpressions returns a boolean if a field has been set.

### GetMatchFields

`func (o *FieldSelector) GetMatchFields() map[string]string`

GetMatchFields returns the MatchFields field if non-nil, zero value otherwise.

### GetMatchFieldsOk

`func (o *FieldSelector) GetMatchFieldsOk() (*map[string]string, bool)`

GetMatchFieldsOk returns a tuple with the MatchFields field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchFields

`func (o *FieldSelector) SetMatchFields(v map[string]string)`

SetMatchFields sets MatchFields field to given value.

### HasMatchFields

`func (o *FieldSelector) HasMatchFields() bool`

HasMatchFields returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


