# DatabaseQuotaStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Used** | Pointer to [**map[string]map[string]IntOrString**](map.md) | The current observed total usage of the named resources per scoped group. | [optional] 

## Methods

### NewDatabaseQuotaStatus

`func NewDatabaseQuotaStatus() *DatabaseQuotaStatus`

NewDatabaseQuotaStatus instantiates a new DatabaseQuotaStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDatabaseQuotaStatusWithDefaults

`func NewDatabaseQuotaStatusWithDefaults() *DatabaseQuotaStatus`

NewDatabaseQuotaStatusWithDefaults instantiates a new DatabaseQuotaStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUsed

`func (o *DatabaseQuotaStatus) GetUsed() map[string]map[string]IntOrString`

GetUsed returns the Used field if non-nil, zero value otherwise.

### GetUsedOk

`func (o *DatabaseQuotaStatus) GetUsedOk() (*map[string]map[string]IntOrString, bool)`

GetUsedOk returns a tuple with the Used field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsed

`func (o *DatabaseQuotaStatus) SetUsed(v map[string]map[string]IntOrString)`

SetUsed sets Used field to given value.

### HasUsed

`func (o *DatabaseQuotaStatus) HasUsed() bool`

HasUsed returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


