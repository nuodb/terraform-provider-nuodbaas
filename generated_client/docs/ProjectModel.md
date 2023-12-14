# ProjectModel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Organization** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Sla** | **string** | The SLA for the project. Cannot be updated once the project is created. | 
**Tier** | **string** | The service tier for the project | 
**Maintenance** | Pointer to [**MaintenanceModel**](MaintenanceModel.md) |  | [optional] 
**ResourceVersion** | Pointer to **string** | The version of the resource. When specified in a &#x60;PUT&#x60; request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates. | [optional] 
**Status** | Pointer to [**ProjectStatusModel**](ProjectStatusModel.md) |  | [optional] 

## Methods

### NewProjectModel

`func NewProjectModel(sla string, tier string, ) *ProjectModel`

NewProjectModel instantiates a new ProjectModel object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectModelWithDefaults

`func NewProjectModelWithDefaults() *ProjectModel`

NewProjectModelWithDefaults instantiates a new ProjectModel object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOrganization

`func (o *ProjectModel) GetOrganization() string`

GetOrganization returns the Organization field if non-nil, zero value otherwise.

### GetOrganizationOk

`func (o *ProjectModel) GetOrganizationOk() (*string, bool)`

GetOrganizationOk returns a tuple with the Organization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganization

`func (o *ProjectModel) SetOrganization(v string)`

SetOrganization sets Organization field to given value.

### HasOrganization

`func (o *ProjectModel) HasOrganization() bool`

HasOrganization returns a boolean if a field has been set.

### GetName

`func (o *ProjectModel) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ProjectModel) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ProjectModel) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ProjectModel) HasName() bool`

HasName returns a boolean if a field has been set.

### GetSla

`func (o *ProjectModel) GetSla() string`

GetSla returns the Sla field if non-nil, zero value otherwise.

### GetSlaOk

`func (o *ProjectModel) GetSlaOk() (*string, bool)`

GetSlaOk returns a tuple with the Sla field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSla

`func (o *ProjectModel) SetSla(v string)`

SetSla sets Sla field to given value.


### GetTier

`func (o *ProjectModel) GetTier() string`

GetTier returns the Tier field if non-nil, zero value otherwise.

### GetTierOk

`func (o *ProjectModel) GetTierOk() (*string, bool)`

GetTierOk returns a tuple with the Tier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTier

`func (o *ProjectModel) SetTier(v string)`

SetTier sets Tier field to given value.


### GetMaintenance

`func (o *ProjectModel) GetMaintenance() MaintenanceModel`

GetMaintenance returns the Maintenance field if non-nil, zero value otherwise.

### GetMaintenanceOk

`func (o *ProjectModel) GetMaintenanceOk() (*MaintenanceModel, bool)`

GetMaintenanceOk returns a tuple with the Maintenance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaintenance

`func (o *ProjectModel) SetMaintenance(v MaintenanceModel)`

SetMaintenance sets Maintenance field to given value.

### HasMaintenance

`func (o *ProjectModel) HasMaintenance() bool`

HasMaintenance returns a boolean if a field has been set.

### GetResourceVersion

`func (o *ProjectModel) GetResourceVersion() string`

GetResourceVersion returns the ResourceVersion field if non-nil, zero value otherwise.

### GetResourceVersionOk

`func (o *ProjectModel) GetResourceVersionOk() (*string, bool)`

GetResourceVersionOk returns a tuple with the ResourceVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResourceVersion

`func (o *ProjectModel) SetResourceVersion(v string)`

SetResourceVersion sets ResourceVersion field to given value.

### HasResourceVersion

`func (o *ProjectModel) HasResourceVersion() bool`

HasResourceVersion returns a boolean if a field has been set.

### GetStatus

`func (o *ProjectModel) GetStatus() ProjectStatusModel`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ProjectModel) GetStatusOk() (*ProjectStatusModel, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ProjectModel) SetStatus(v ProjectStatusModel)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ProjectModel) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


