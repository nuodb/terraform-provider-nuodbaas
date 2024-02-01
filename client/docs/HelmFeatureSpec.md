# HelmFeatureSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ChartCompatibility** | Pointer to **string** | The Helm chart version compatibility constraint for the Helm feature. | [optional] 
**Optional** | Pointer to **bool** | Whether the Helm feature is optional and does not emit an error if the Helm chart or product version is incompatible. | [optional] 
**Parameters** | Pointer to [**map[string]Parameters**](Parameters.md) | The parameter definitions referenced in values. For example, parameter named &#x60;foo&#x60; is referenced using &#x60;&lt;&lt; .meta.params.foo &gt;&gt;&#x60; template. | [optional] 
**ProductCompatibility** | Pointer to **string** | The NuoDB product version compatibility constraint for the Helm feature. | [optional] 
**Values** | Pointer to [**map[string]AnyTypeValue**](AnyTypeValue.md) |  | [optional] 

## Methods

### NewHelmFeatureSpec

`func NewHelmFeatureSpec() *HelmFeatureSpec`

NewHelmFeatureSpec instantiates a new HelmFeatureSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewHelmFeatureSpecWithDefaults

`func NewHelmFeatureSpecWithDefaults() *HelmFeatureSpec`

NewHelmFeatureSpecWithDefaults instantiates a new HelmFeatureSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChartCompatibility

`func (o *HelmFeatureSpec) GetChartCompatibility() string`

GetChartCompatibility returns the ChartCompatibility field if non-nil, zero value otherwise.

### GetChartCompatibilityOk

`func (o *HelmFeatureSpec) GetChartCompatibilityOk() (*string, bool)`

GetChartCompatibilityOk returns a tuple with the ChartCompatibility field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChartCompatibility

`func (o *HelmFeatureSpec) SetChartCompatibility(v string)`

SetChartCompatibility sets ChartCompatibility field to given value.

### HasChartCompatibility

`func (o *HelmFeatureSpec) HasChartCompatibility() bool`

HasChartCompatibility returns a boolean if a field has been set.

### GetOptional

`func (o *HelmFeatureSpec) GetOptional() bool`

GetOptional returns the Optional field if non-nil, zero value otherwise.

### GetOptionalOk

`func (o *HelmFeatureSpec) GetOptionalOk() (*bool, bool)`

GetOptionalOk returns a tuple with the Optional field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOptional

`func (o *HelmFeatureSpec) SetOptional(v bool)`

SetOptional sets Optional field to given value.

### HasOptional

`func (o *HelmFeatureSpec) HasOptional() bool`

HasOptional returns a boolean if a field has been set.

### GetParameters

`func (o *HelmFeatureSpec) GetParameters() map[string]Parameters`

GetParameters returns the Parameters field if non-nil, zero value otherwise.

### GetParametersOk

`func (o *HelmFeatureSpec) GetParametersOk() (*map[string]Parameters, bool)`

GetParametersOk returns a tuple with the Parameters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetParameters

`func (o *HelmFeatureSpec) SetParameters(v map[string]Parameters)`

SetParameters sets Parameters field to given value.

### HasParameters

`func (o *HelmFeatureSpec) HasParameters() bool`

HasParameters returns a boolean if a field has been set.

### GetProductCompatibility

`func (o *HelmFeatureSpec) GetProductCompatibility() string`

GetProductCompatibility returns the ProductCompatibility field if non-nil, zero value otherwise.

### GetProductCompatibilityOk

`func (o *HelmFeatureSpec) GetProductCompatibilityOk() (*string, bool)`

GetProductCompatibilityOk returns a tuple with the ProductCompatibility field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProductCompatibility

`func (o *HelmFeatureSpec) SetProductCompatibility(v string)`

SetProductCompatibility sets ProductCompatibility field to given value.

### HasProductCompatibility

`func (o *HelmFeatureSpec) HasProductCompatibility() bool`

HasProductCompatibility returns a boolean if a field has been set.

### GetValues

`func (o *HelmFeatureSpec) GetValues() map[string]AnyTypeValue`

GetValues returns the Values field if non-nil, zero value otherwise.

### GetValuesOk

`func (o *HelmFeatureSpec) GetValuesOk() (*map[string]AnyTypeValue, bool)`

GetValuesOk returns a tuple with the Values field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValues

`func (o *HelmFeatureSpec) SetValues(v map[string]AnyTypeValue)`

SetValues sets Values field to given value.

### HasValues

`func (o *HelmFeatureSpec) HasValues() bool`

HasValues returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


