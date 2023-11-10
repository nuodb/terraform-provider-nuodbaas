# MaintenanceModel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ExpiresAtTime** | Pointer to **time.Time** | The time at which the project or database will be disabled | [optional] 
**ExpiresIn** | Pointer to **string** | The time until the project or database is disabled, e.g. &#x60;1d&#x60; | [optional] 
**IsDisabled** | Pointer to **bool** | Whether the project or database should be shutdown | [optional] 

## Methods

### NewMaintenanceModel

`func NewMaintenanceModel() *MaintenanceModel`

NewMaintenanceModel instantiates a new MaintenanceModel object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMaintenanceModelWithDefaults

`func NewMaintenanceModelWithDefaults() *MaintenanceModel`

NewMaintenanceModelWithDefaults instantiates a new MaintenanceModel object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetExpiresAtTime

`func (o *MaintenanceModel) GetExpiresAtTime() time.Time`

GetExpiresAtTime returns the ExpiresAtTime field if non-nil, zero value otherwise.

### GetExpiresAtTimeOk

`func (o *MaintenanceModel) GetExpiresAtTimeOk() (*time.Time, bool)`

GetExpiresAtTimeOk returns a tuple with the ExpiresAtTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAtTime

`func (o *MaintenanceModel) SetExpiresAtTime(v time.Time)`

SetExpiresAtTime sets ExpiresAtTime field to given value.

### HasExpiresAtTime

`func (o *MaintenanceModel) HasExpiresAtTime() bool`

HasExpiresAtTime returns a boolean if a field has been set.

### GetExpiresIn

`func (o *MaintenanceModel) GetExpiresIn() string`

GetExpiresIn returns the ExpiresIn field if non-nil, zero value otherwise.

### GetExpiresInOk

`func (o *MaintenanceModel) GetExpiresInOk() (*string, bool)`

GetExpiresInOk returns a tuple with the ExpiresIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresIn

`func (o *MaintenanceModel) SetExpiresIn(v string)`

SetExpiresIn sets ExpiresIn field to given value.

### HasExpiresIn

`func (o *MaintenanceModel) HasExpiresIn() bool`

HasExpiresIn returns a boolean if a field has been set.

### GetIsDisabled

`func (o *MaintenanceModel) GetIsDisabled() bool`

GetIsDisabled returns the IsDisabled field if non-nil, zero value otherwise.

### GetIsDisabledOk

`func (o *MaintenanceModel) GetIsDisabledOk() (*bool, bool)`

GetIsDisabledOk returns a tuple with the IsDisabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsDisabled

`func (o *MaintenanceModel) SetIsDisabled(v bool)`

SetIsDisabled sets IsDisabled field to given value.

### HasIsDisabled

`func (o *MaintenanceModel) HasIsDisabled() bool`

HasIsDisabled returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


