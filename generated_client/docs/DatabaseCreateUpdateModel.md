# DatabaseCreateUpdateModel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Organization** | Pointer to **string** |  | [optional] 
**Project** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**DbaPassword** | Pointer to **string** | The password for the DBA user. Can only be specified when creating a database. | [optional] 
**Tier** | Pointer to **string** | The service tier for the database. If omitted, the project service tier is inherited. | [optional] 
**Maintenance** | Pointer to [**MaintenanceModel**](MaintenanceModel.md) |  | [optional] 
**Properties** | Pointer to [**DatabasePropertiesModel**](DatabasePropertiesModel.md) |  | [optional] 
**ResourceVersion** | Pointer to **string** | The version of the resource. When specified in a &#x60;PUT&#x60; request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates. | [optional] 
**Status** | Pointer to [**DatabaseStatusModel**](DatabaseStatusModel.md) |  | [optional] 

## Methods

### NewDatabaseCreateUpdateModel

`func NewDatabaseCreateUpdateModel() *DatabaseCreateUpdateModel`

NewDatabaseCreateUpdateModel instantiates a new DatabaseCreateUpdateModel object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDatabaseCreateUpdateModelWithDefaults

`func NewDatabaseCreateUpdateModelWithDefaults() *DatabaseCreateUpdateModel`

NewDatabaseCreateUpdateModelWithDefaults instantiates a new DatabaseCreateUpdateModel object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOrganization

`func (o *DatabaseCreateUpdateModel) GetOrganization() string`

GetOrganization returns the Organization field if non-nil, zero value otherwise.

### GetOrganizationOk

`func (o *DatabaseCreateUpdateModel) GetOrganizationOk() (*string, bool)`

GetOrganizationOk returns a tuple with the Organization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganization

`func (o *DatabaseCreateUpdateModel) SetOrganization(v string)`

SetOrganization sets Organization field to given value.

### HasOrganization

`func (o *DatabaseCreateUpdateModel) HasOrganization() bool`

HasOrganization returns a boolean if a field has been set.

### GetProject

`func (o *DatabaseCreateUpdateModel) GetProject() string`

GetProject returns the Project field if non-nil, zero value otherwise.

### GetProjectOk

`func (o *DatabaseCreateUpdateModel) GetProjectOk() (*string, bool)`

GetProjectOk returns a tuple with the Project field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProject

`func (o *DatabaseCreateUpdateModel) SetProject(v string)`

SetProject sets Project field to given value.

### HasProject

`func (o *DatabaseCreateUpdateModel) HasProject() bool`

HasProject returns a boolean if a field has been set.

### GetName

`func (o *DatabaseCreateUpdateModel) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DatabaseCreateUpdateModel) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DatabaseCreateUpdateModel) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DatabaseCreateUpdateModel) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDbaPassword

`func (o *DatabaseCreateUpdateModel) GetDbaPassword() string`

GetDbaPassword returns the DbaPassword field if non-nil, zero value otherwise.

### GetDbaPasswordOk

`func (o *DatabaseCreateUpdateModel) GetDbaPasswordOk() (*string, bool)`

GetDbaPasswordOk returns a tuple with the DbaPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDbaPassword

`func (o *DatabaseCreateUpdateModel) SetDbaPassword(v string)`

SetDbaPassword sets DbaPassword field to given value.

### HasDbaPassword

`func (o *DatabaseCreateUpdateModel) HasDbaPassword() bool`

HasDbaPassword returns a boolean if a field has been set.

### GetTier

`func (o *DatabaseCreateUpdateModel) GetTier() string`

GetTier returns the Tier field if non-nil, zero value otherwise.

### GetTierOk

`func (o *DatabaseCreateUpdateModel) GetTierOk() (*string, bool)`

GetTierOk returns a tuple with the Tier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTier

`func (o *DatabaseCreateUpdateModel) SetTier(v string)`

SetTier sets Tier field to given value.

### HasTier

`func (o *DatabaseCreateUpdateModel) HasTier() bool`

HasTier returns a boolean if a field has been set.

### GetMaintenance

`func (o *DatabaseCreateUpdateModel) GetMaintenance() MaintenanceModel`

GetMaintenance returns the Maintenance field if non-nil, zero value otherwise.

### GetMaintenanceOk

`func (o *DatabaseCreateUpdateModel) GetMaintenanceOk() (*MaintenanceModel, bool)`

GetMaintenanceOk returns a tuple with the Maintenance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaintenance

`func (o *DatabaseCreateUpdateModel) SetMaintenance(v MaintenanceModel)`

SetMaintenance sets Maintenance field to given value.

### HasMaintenance

`func (o *DatabaseCreateUpdateModel) HasMaintenance() bool`

HasMaintenance returns a boolean if a field has been set.

### GetProperties

`func (o *DatabaseCreateUpdateModel) GetProperties() DatabasePropertiesModel`

GetProperties returns the Properties field if non-nil, zero value otherwise.

### GetPropertiesOk

`func (o *DatabaseCreateUpdateModel) GetPropertiesOk() (*DatabasePropertiesModel, bool)`

GetPropertiesOk returns a tuple with the Properties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProperties

`func (o *DatabaseCreateUpdateModel) SetProperties(v DatabasePropertiesModel)`

SetProperties sets Properties field to given value.

### HasProperties

`func (o *DatabaseCreateUpdateModel) HasProperties() bool`

HasProperties returns a boolean if a field has been set.

### GetResourceVersion

`func (o *DatabaseCreateUpdateModel) GetResourceVersion() string`

GetResourceVersion returns the ResourceVersion field if non-nil, zero value otherwise.

### GetResourceVersionOk

`func (o *DatabaseCreateUpdateModel) GetResourceVersionOk() (*string, bool)`

GetResourceVersionOk returns a tuple with the ResourceVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceVersion

`func (o *DatabaseCreateUpdateModel) SetResourceVersion(v string)`

SetResourceVersion sets ResourceVersion field to given value.

### HasResourceVersion

`func (o *DatabaseCreateUpdateModel) HasResourceVersion() bool`

HasResourceVersion returns a boolean if a field has been set.

### GetStatus

`func (o *DatabaseCreateUpdateModel) GetStatus() DatabaseStatusModel`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *DatabaseCreateUpdateModel) GetStatusOk() (*DatabaseStatusModel, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *DatabaseCreateUpdateModel) SetStatus(v DatabaseStatusModel)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *DatabaseCreateUpdateModel) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


