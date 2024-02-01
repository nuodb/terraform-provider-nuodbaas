# DatabaseQuotaModel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the resource | 
**Description** | Pointer to **string** | Human-readable description of the resource | [optional] 
**ResourceVersion** | Pointer to **string** | The version of the resource. When specified in a &#x60;PUT&#x60; request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates. | [optional] 
**Spec** | [**DatabaseQuotaSpec**](DatabaseQuotaSpec.md) |  | 
**Status** | Pointer to [**DatabaseQuotaStatus**](DatabaseQuotaStatus.md) |  | [optional] 

## Methods

### NewDatabaseQuotaModel

`func NewDatabaseQuotaModel(name string, spec DatabaseQuotaSpec, ) *DatabaseQuotaModel`

NewDatabaseQuotaModel instantiates a new DatabaseQuotaModel object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDatabaseQuotaModelWithDefaults

`func NewDatabaseQuotaModelWithDefaults() *DatabaseQuotaModel`

NewDatabaseQuotaModelWithDefaults instantiates a new DatabaseQuotaModel object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *DatabaseQuotaModel) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DatabaseQuotaModel) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DatabaseQuotaModel) SetName(v string)`

SetName sets Name field to given value.


### GetDescription

`func (o *DatabaseQuotaModel) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *DatabaseQuotaModel) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *DatabaseQuotaModel) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *DatabaseQuotaModel) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetResourceVersion

`func (o *DatabaseQuotaModel) GetResourceVersion() string`

GetResourceVersion returns the ResourceVersion field if non-nil, zero value otherwise.

### GetResourceVersionOk

`func (o *DatabaseQuotaModel) GetResourceVersionOk() (*string, bool)`

GetResourceVersionOk returns a tuple with the ResourceVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceVersion

`func (o *DatabaseQuotaModel) SetResourceVersion(v string)`

SetResourceVersion sets ResourceVersion field to given value.

### HasResourceVersion

`func (o *DatabaseQuotaModel) HasResourceVersion() bool`

HasResourceVersion returns a boolean if a field has been set.

### GetSpec

`func (o *DatabaseQuotaModel) GetSpec() DatabaseQuotaSpec`

GetSpec returns the Spec field if non-nil, zero value otherwise.

### GetSpecOk

`func (o *DatabaseQuotaModel) GetSpecOk() (*DatabaseQuotaSpec, bool)`

GetSpecOk returns a tuple with the Spec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpec

`func (o *DatabaseQuotaModel) SetSpec(v DatabaseQuotaSpec)`

SetSpec sets Spec field to given value.


### GetStatus

`func (o *DatabaseQuotaModel) GetStatus() DatabaseQuotaStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *DatabaseQuotaModel) GetStatusOk() (*DatabaseQuotaStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *DatabaseQuotaModel) SetStatus(v DatabaseQuotaStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *DatabaseQuotaModel) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


