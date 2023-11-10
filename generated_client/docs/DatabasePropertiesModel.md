# DatabasePropertiesModel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ArchiveDiskSize** | Pointer to **string** | The size of the archive volumes for the database. Can be only updated to increase the volume size. | [optional] 
**JournalDiskSize** | Pointer to **string** | The size of the journal volumes for the database. Can be only updated to increase the volume size. | [optional] 

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


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


