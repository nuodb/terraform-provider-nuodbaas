# Scope

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FieldSelector** | Pointer to [**FieldSelector**](FieldSelector.md) |  | [optional] 
**GroupByLabels** | Pointer to **[]string** | The label keys on which the selected databases are divided into groups. | [optional] 
**LabelSelector** | Pointer to [**LabelSelector**](LabelSelector.md) |  | [optional] 

## Methods

### NewScope

`func NewScope() *Scope`

NewScope instantiates a new Scope object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewScopeWithDefaults

`func NewScopeWithDefaults() *Scope`

NewScopeWithDefaults instantiates a new Scope object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFieldSelector

`func (o *Scope) GetFieldSelector() FieldSelector`

GetFieldSelector returns the FieldSelector field if non-nil, zero value otherwise.

### GetFieldSelectorOk

`func (o *Scope) GetFieldSelectorOk() (*FieldSelector, bool)`

GetFieldSelectorOk returns a tuple with the FieldSelector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFieldSelector

`func (o *Scope) SetFieldSelector(v FieldSelector)`

SetFieldSelector sets FieldSelector field to given value.

### HasFieldSelector

`func (o *Scope) HasFieldSelector() bool`

HasFieldSelector returns a boolean if a field has been set.

### GetGroupByLabels

`func (o *Scope) GetGroupByLabels() []string`

GetGroupByLabels returns the GroupByLabels field if non-nil, zero value otherwise.

### GetGroupByLabelsOk

`func (o *Scope) GetGroupByLabelsOk() (*[]string, bool)`

GetGroupByLabelsOk returns a tuple with the GroupByLabels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroupByLabels

`func (o *Scope) SetGroupByLabels(v []string)`

SetGroupByLabels sets GroupByLabels field to given value.

### HasGroupByLabels

`func (o *Scope) HasGroupByLabels() bool`

HasGroupByLabels returns a boolean if a field has been set.

### GetLabelSelector

`func (o *Scope) GetLabelSelector() LabelSelector`

GetLabelSelector returns the LabelSelector field if non-nil, zero value otherwise.

### GetLabelSelectorOk

`func (o *Scope) GetLabelSelectorOk() (*LabelSelector, bool)`

GetLabelSelectorOk returns a tuple with the LabelSelector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabelSelector

`func (o *Scope) SetLabelSelector(v LabelSelector)`

SetLabelSelector sets LabelSelector field to given value.

### HasLabelSelector

`func (o *Scope) HasLabelSelector() bool`

HasLabelSelector returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


