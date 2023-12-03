# Parameters

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Default** | Pointer to **string** |  | [optional] 
**JsonSchema** | Pointer to **string** | A JSONSchema used to validate the parameter&#39;s value. | [optional] 

## Methods

### NewParameters

`func NewParameters() *Parameters`

NewParameters instantiates a new Parameters object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewParametersWithDefaults

`func NewParametersWithDefaults() *Parameters`

NewParametersWithDefaults instantiates a new Parameters object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDefault

`func (o *Parameters) GetDefault() string`

GetDefault returns the Default field if non-nil, zero value otherwise.

### GetDefaultOk

`func (o *Parameters) GetDefaultOk() (*string, bool)`

GetDefaultOk returns a tuple with the Default field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefault

`func (o *Parameters) SetDefault(v string)`

SetDefault sets Default field to given value.

### HasDefault

`func (o *Parameters) HasDefault() bool`

HasDefault returns a boolean if a field has been set.

### GetJsonSchema

`func (o *Parameters) GetJsonSchema() string`

GetJsonSchema returns the JsonSchema field if non-nil, zero value otherwise.

### GetJsonSchemaOk

`func (o *Parameters) GetJsonSchemaOk() (*string, bool)`

GetJsonSchemaOk returns a tuple with the JsonSchema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJsonSchema

`func (o *Parameters) SetJsonSchema(v string)`

SetJsonSchema sets JsonSchema field to given value.

### HasJsonSchema

`func (o *Parameters) HasJsonSchema() bool`

HasJsonSchema returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


