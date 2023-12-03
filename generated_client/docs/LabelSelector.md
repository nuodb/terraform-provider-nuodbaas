# LabelSelector

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MatchExpressions** | Pointer to [**[]DatabasequotaspecScopeLabelselectorMatchExpressions**](DatabasequotaspecScopeLabelselectorMatchExpressions.md) | matchExpressions is a list of label selector requirements. The requirements are ANDed. | [optional] 
**MatchLabels** | Pointer to **map[string]string** | matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is \&quot;key\&quot;, the operator is \&quot;In\&quot;, and the values array contains only \&quot;value\&quot;. The requirements are ANDed. | [optional] 

## Methods

### NewLabelSelector

`func NewLabelSelector() *LabelSelector`

NewLabelSelector instantiates a new LabelSelector object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLabelSelectorWithDefaults

`func NewLabelSelectorWithDefaults() *LabelSelector`

NewLabelSelectorWithDefaults instantiates a new LabelSelector object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMatchExpressions

`func (o *LabelSelector) GetMatchExpressions() []DatabasequotaspecScopeLabelselectorMatchExpressions`

GetMatchExpressions returns the MatchExpressions field if non-nil, zero value otherwise.

### GetMatchExpressionsOk

`func (o *LabelSelector) GetMatchExpressionsOk() (*[]DatabasequotaspecScopeLabelselectorMatchExpressions, bool)`

GetMatchExpressionsOk returns a tuple with the MatchExpressions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchExpressions

`func (o *LabelSelector) SetMatchExpressions(v []DatabasequotaspecScopeLabelselectorMatchExpressions)`

SetMatchExpressions sets MatchExpressions field to given value.

### HasMatchExpressions

`func (o *LabelSelector) HasMatchExpressions() bool`

HasMatchExpressions returns a boolean if a field has been set.

### GetMatchLabels

`func (o *LabelSelector) GetMatchLabels() map[string]string`

GetMatchLabels returns the MatchLabels field if non-nil, zero value otherwise.

### GetMatchLabelsOk

`func (o *LabelSelector) GetMatchLabelsOk() (*map[string]string, bool)`

GetMatchLabelsOk returns a tuple with the MatchLabels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchLabels

`func (o *LabelSelector) SetMatchLabels(v map[string]string)`

SetMatchLabels sets MatchLabels field to given value.

### HasMatchLabels

`func (o *LabelSelector) HasMatchLabels() bool`

HasMatchLabels returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


