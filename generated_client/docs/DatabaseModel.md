# DatabaseModel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Organization** | Pointer to **string** |  | [optional] 
**Project** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **map[string]string** | User-defined labels attached to the resource that can be used for filtering | [optional] 
**Tier** | Pointer to **string** | The service tier for the database. If omitted, the project service tier is inherited. | [optional] 
**Maintenance** | Pointer to [**MaintenanceModel**](MaintenanceModel.md) |  | [optional] 
**Properties** | Pointer to [**DatabasePropertiesModel**](DatabasePropertiesModel.md) |  | [optional] 
**ResourceVersion** | Pointer to **string** | The version of the resource. When specified in a &#x60;PUT&#x60; request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates. | [optional] 
**Status** | Pointer to [**DatabaseStatusModel**](DatabaseStatusModel.md) |  | [optional] 

## Methods

### NewDatabaseModel

`func NewDatabaseModel() *DatabaseModel`

NewDatabaseModel instantiates a new DatabaseModel object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDatabaseModelWithDefaults

`func NewDatabaseModelWithDefaults() *DatabaseModel`

NewDatabaseModelWithDefaults instantiates a new DatabaseModel object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOrganization

`func (o *DatabaseModel) GetOrganization() string`

GetOrganization returns the Organization field if non-nil, zero value otherwise.

### GetOrganizationOk

`func (o *DatabaseModel) GetOrganizationOk() (*string, bool)`

GetOrganizationOk returns a tuple with the Organization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganization

`func (o *DatabaseModel) SetOrganization(v string)`

SetOrganization sets Organization field to given value.

### HasOrganization

`func (o *DatabaseModel) HasOrganization() bool`

HasOrganization returns a boolean if a field has been set.

### GetProject

`func (o *DatabaseModel) GetProject() string`

GetProject returns the Project field if non-nil, zero value otherwise.

### GetProjectOk

`func (o *DatabaseModel) GetProjectOk() (*string, bool)`

GetProjectOk returns a tuple with the Project field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProject

`func (o *DatabaseModel) SetProject(v string)`

SetProject sets Project field to given value.

### HasProject

`func (o *DatabaseModel) HasProject() bool`

HasProject returns a boolean if a field has been set.

### GetName

`func (o *DatabaseModel) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DatabaseModel) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DatabaseModel) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DatabaseModel) HasName() bool`

HasName returns a boolean if a field has been set.

### GetLabels

`func (o *DatabaseModel) GetLabels() map[string]string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *DatabaseModel) GetLabelsOk() (*map[string]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *DatabaseModel) SetLabels(v map[string]string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *DatabaseModel) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetTier

`func (o *DatabaseModel) GetTier() string`

GetTier returns the Tier field if non-nil, zero value otherwise.

### GetTierOk

`func (o *DatabaseModel) GetTierOk() (*string, bool)`

GetTierOk returns a tuple with the Tier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTier

`func (o *DatabaseModel) SetTier(v string)`

SetTier sets Tier field to given value.

### HasTier

`func (o *DatabaseModel) HasTier() bool`

HasTier returns a boolean if a field has been set.

### GetMaintenance

`func (o *DatabaseModel) GetMaintenance() MaintenanceModel`

GetMaintenance returns the Maintenance field if non-nil, zero value otherwise.

### GetMaintenanceOk

`func (o *DatabaseModel) GetMaintenanceOk() (*MaintenanceModel, bool)`

GetMaintenanceOk returns a tuple with the Maintenance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaintenance

`func (o *DatabaseModel) SetMaintenance(v MaintenanceModel)`

SetMaintenance sets Maintenance field to given value.

### HasMaintenance

`func (o *DatabaseModel) HasMaintenance() bool`

HasMaintenance returns a boolean if a field has been set.

### GetProperties

`func (o *DatabaseModel) GetProperties() DatabasePropertiesModel`

GetProperties returns the Properties field if non-nil, zero value otherwise.

### GetPropertiesOk

`func (o *DatabaseModel) GetPropertiesOk() (*DatabasePropertiesModel, bool)`

GetPropertiesOk returns a tuple with the Properties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProperties

`func (o *DatabaseModel) SetProperties(v DatabasePropertiesModel)`

SetProperties sets Properties field to given value.

### HasProperties

`func (o *DatabaseModel) HasProperties() bool`

HasProperties returns a boolean if a field has been set.

### GetResourceVersion

`func (o *DatabaseModel) GetResourceVersion() string`

GetResourceVersion returns the ResourceVersion field if non-nil, zero value otherwise.

### GetResourceVersionOk

`func (o *DatabaseModel) GetResourceVersionOk() (*string, bool)`

GetResourceVersionOk returns a tuple with the ResourceVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceVersion

`func (o *DatabaseModel) SetResourceVersion(v string)`

SetResourceVersion sets ResourceVersion field to given value.

### HasResourceVersion

`func (o *DatabaseModel) HasResourceVersion() bool`

HasResourceVersion returns a boolean if a field has been set.

### GetStatus

`func (o *DatabaseModel) GetStatus() DatabaseStatusModel`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *DatabaseModel) GetStatusOk() (*DatabaseStatusModel, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *DatabaseModel) SetStatus(v DatabaseStatusModel)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *DatabaseModel) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


