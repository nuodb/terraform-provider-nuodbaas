# DatabaseQuotaSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Hard** | Pointer to [**map[string]IntOrString**](IntOrString.md) | The set of desired hard limits for each named resource. | [optional] 
**Scope** | Pointer to [**Scope**](Scope.md) |  | [optional] 

## Methods

### NewDatabaseQuotaSpec

`func NewDatabaseQuotaSpec() *DatabaseQuotaSpec`

NewDatabaseQuotaSpec instantiates a new DatabaseQuotaSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDatabaseQuotaSpecWithDefaults

`func NewDatabaseQuotaSpecWithDefaults() *DatabaseQuotaSpec`

NewDatabaseQuotaSpecWithDefaults instantiates a new DatabaseQuotaSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHard

`func (o *DatabaseQuotaSpec) GetHard() map[string]IntOrString`

GetHard returns the Hard field if non-nil, zero value otherwise.

### GetHardOk

`func (o *DatabaseQuotaSpec) GetHardOk() (*map[string]IntOrString, bool)`

GetHardOk returns a tuple with the Hard field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHard

`func (o *DatabaseQuotaSpec) SetHard(v map[string]IntOrString)`

SetHard sets Hard field to given value.

### HasHard

`func (o *DatabaseQuotaSpec) HasHard() bool`

HasHard returns a boolean if a field has been set.

### GetScope

`func (o *DatabaseQuotaSpec) GetScope() Scope`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *DatabaseQuotaSpec) GetScopeOk() (*Scope, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *DatabaseQuotaSpec) SetScope(v Scope)`

SetScope sets Scope field to given value.

### HasScope

`func (o *DatabaseQuotaSpec) HasScope() bool`

HasScope returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


