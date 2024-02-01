# DatabasePropertiesModel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ArchiveDiskSize** | Pointer to **string** | The size of the archive volumes for the database. Can be only updated to increase the volume size. | [optional] 
**JournalDiskSize** | Pointer to **string** | The size of the journal volumes for the database. Can be only updated to increase the volume size. | [optional] 
**TierParameters** | Pointer to **map[string]string** | Opaque parameters supplied to database service tier. | [optional] 
**InheritTierParameters** | Pointer to **bool** | Whether to inherit tier parameters from the project if the database service tier matches the project. | [optional] 
**ProductVersion** | Pointer to **string** | The version/tag of the NuoDB image to use. For available tags, see https://hub.docker.com/r/nuodb/nuodb-ce/tags. If omitted, the database version will be inherited from the project. | [optional] 

## Methods

### NewDatabasePropertiesModel

`func NewDatabasePropertiesModel() *DatabasePropertiesModel`

NewDatabasePropertiesModel instantiates a new DatabasePropertiesModel object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDatabasePropertiesModelWithDefaults

`func NewDatabasePropertiesModelWithDefaults() *DatabasePropertiesModel`

NewDatabasePropertiesModelWithDefaults instantiates a new DatabasePropertiesModel object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetArchiveDiskSize

`func (o *DatabasePropertiesModel) GetArchiveDiskSize() string`

GetArchiveDiskSize returns the ArchiveDiskSize field if non-nil, zero value otherwise.

### GetArchiveDiskSizeOk

`func (o *DatabasePropertiesModel) GetArchiveDiskSizeOk() (*string, bool)`

GetArchiveDiskSizeOk returns a tuple with the ArchiveDiskSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetArchiveDiskSize

`func (o *DatabasePropertiesModel) SetArchiveDiskSize(v string)`

SetArchiveDiskSize sets ArchiveDiskSize field to given value.

### HasArchiveDiskSize

`func (o *DatabasePropertiesModel) HasArchiveDiskSize() bool`

HasArchiveDiskSize returns a boolean if a field has been set.

### GetJournalDiskSize

`func (o *DatabasePropertiesModel) GetJournalDiskSize() string`

GetJournalDiskSize returns the JournalDiskSize field if non-nil, zero value otherwise.

### GetJournalDiskSizeOk

`func (o *DatabasePropertiesModel) GetJournalDiskSizeOk() (*string, bool)`

GetJournalDiskSizeOk returns a tuple with the JournalDiskSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJournalDiskSize

`func (o *DatabasePropertiesModel) SetJournalDiskSize(v string)`

SetJournalDiskSize sets JournalDiskSize field to given value.

### HasJournalDiskSize

`func (o *DatabasePropertiesModel) HasJournalDiskSize() bool`

HasJournalDiskSize returns a boolean if a field has been set.

### GetTierParameters

`func (o *DatabasePropertiesModel) GetTierParameters() map[string]string`

GetTierParameters returns the TierParameters field if non-nil, zero value otherwise.

### GetTierParametersOk

`func (o *DatabasePropertiesModel) GetTierParametersOk() (*map[string]string, bool)`

GetTierParametersOk returns a tuple with the TierParameters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTierParameters

`func (o *DatabasePropertiesModel) SetTierParameters(v map[string]string)`

SetTierParameters sets TierParameters field to given value.

### HasTierParameters

`func (o *DatabasePropertiesModel) HasTierParameters() bool`

HasTierParameters returns a boolean if a field has been set.

### GetInheritTierParameters

`func (o *DatabasePropertiesModel) GetInheritTierParameters() bool`

GetInheritTierParameters returns the InheritTierParameters field if non-nil, zero value otherwise.

### GetInheritTierParametersOk

`func (o *DatabasePropertiesModel) GetInheritTierParametersOk() (*bool, bool)`

GetInheritTierParametersOk returns a tuple with the InheritTierParameters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInheritTierParameters

`func (o *DatabasePropertiesModel) SetInheritTierParameters(v bool)`

SetInheritTierParameters sets InheritTierParameters field to given value.

### HasInheritTierParameters

`func (o *DatabasePropertiesModel) HasInheritTierParameters() bool`

HasInheritTierParameters returns a boolean if a field has been set.

### GetProductVersion

`func (o *DatabasePropertiesModel) GetProductVersion() string`

GetProductVersion returns the ProductVersion field if non-nil, zero value otherwise.

### GetProductVersionOk

`func (o *DatabasePropertiesModel) GetProductVersionOk() (*string, bool)`

GetProductVersionOk returns a tuple with the ProductVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProductVersion

`func (o *DatabasePropertiesModel) SetProductVersion(v string)`

SetProductVersion sets ProductVersion field to given value.

### HasProductVersion

`func (o *DatabasePropertiesModel) HasProductVersion() bool`

HasProductVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


