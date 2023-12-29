# ProjectPropertiesModel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TierParameters** | Pointer to **map[string]string** | Opaque parameters supplied to project service tier. | [optional] 
**ProductVersion** | Pointer to **string** | The version/tag of the NuoDB image to use. For available tags, see https://hub.docker.com/r/nuodb/nuodb-ce/tags. If omitted, the project version will be resolved based on the SLA and cluster configuration. | [optional] 

## Methods

### NewProjectPropertiesModel

`func NewProjectPropertiesModel() *ProjectPropertiesModel`

NewProjectPropertiesModel instantiates a new ProjectPropertiesModel object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectPropertiesModelWithDefaults

`func NewProjectPropertiesModelWithDefaults() *ProjectPropertiesModel`

NewProjectPropertiesModelWithDefaults instantiates a new ProjectPropertiesModel object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTierParameters

`func (o *ProjectPropertiesModel) GetTierParameters() map[string]string`

GetTierParameters returns the TierParameters field if non-nil, zero value otherwise.

### GetTierParametersOk

`func (o *ProjectPropertiesModel) GetTierParametersOk() (*map[string]string, bool)`

GetTierParametersOk returns a tuple with the TierParameters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTierParameters

`func (o *ProjectPropertiesModel) SetTierParameters(v map[string]string)`

SetTierParameters sets TierParameters field to given value.

### HasTierParameters

`func (o *ProjectPropertiesModel) HasTierParameters() bool`

HasTierParameters returns a boolean if a field has been set.

### GetProductVersion

`func (o *ProjectPropertiesModel) GetProductVersion() string`

GetProductVersion returns the ProductVersion field if non-nil, zero value otherwise.

### GetProductVersionOk

`func (o *ProjectPropertiesModel) GetProductVersionOk() (*string, bool)`

GetProductVersionOk returns a tuple with the ProductVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProductVersion

`func (o *ProjectPropertiesModel) SetProductVersion(v string)`

SetProductVersion sets ProductVersion field to given value.

### HasProductVersion

`func (o *ProjectPropertiesModel) HasProductVersion() bool`

HasProductVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


