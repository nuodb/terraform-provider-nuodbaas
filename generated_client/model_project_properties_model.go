/*
NuoDB Control Plane REST API

NuoDB Control Plane (CP) allows users to create and manage NuoDB databases remotely using a Database as a Service (DBaaS) model.

API version: 2.3.0
Contact: NuoDB.Support@3ds.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package nuodbaas

import (
	"encoding/json"
)

// checks if the ProjectPropertiesModel type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ProjectPropertiesModel{}

// ProjectPropertiesModel struct for ProjectPropertiesModel
type ProjectPropertiesModel struct {
	// Opaque parameters supplied to project service tier.
	TierParameters *map[string]string `json:"tierParameters,omitempty"`
}

// NewProjectPropertiesModel instantiates a new ProjectPropertiesModel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewProjectPropertiesModel() *ProjectPropertiesModel {
	this := ProjectPropertiesModel{}
	return &this
}

// NewProjectPropertiesModelWithDefaults instantiates a new ProjectPropertiesModel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewProjectPropertiesModelWithDefaults() *ProjectPropertiesModel {
	this := ProjectPropertiesModel{}
	return &this
}

// GetTierParameters returns the TierParameters field value if set, zero value otherwise.
func (o *ProjectPropertiesModel) GetTierParameters() map[string]string {
	if o == nil || IsNil(o.TierParameters) {
		var ret map[string]string
		return ret
	}
	return *o.TierParameters
}

// GetTierParametersOk returns a tuple with the TierParameters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectPropertiesModel) GetTierParametersOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.TierParameters) {
		return nil, false
	}
	return o.TierParameters, true
}

// HasTierParameters returns a boolean if a field has been set.
func (o *ProjectPropertiesModel) HasTierParameters() bool {
	if o != nil && !IsNil(o.TierParameters) {
		return true
	}

	return false
}

// SetTierParameters gets a reference to the given map[string]string and assigns it to the TierParameters field.
func (o *ProjectPropertiesModel) SetTierParameters(v map[string]string) {
	o.TierParameters = &v
}

func (o ProjectPropertiesModel) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ProjectPropertiesModel) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.TierParameters) {
		toSerialize["tierParameters"] = o.TierParameters
	}
	return toSerialize, nil
}

type NullableProjectPropertiesModel struct {
	value *ProjectPropertiesModel
	isSet bool
}

func (v NullableProjectPropertiesModel) Get() *ProjectPropertiesModel {
	return v.value
}

func (v *NullableProjectPropertiesModel) Set(val *ProjectPropertiesModel) {
	v.value = val
	v.isSet = true
}

func (v NullableProjectPropertiesModel) IsSet() bool {
	return v.isSet
}

func (v *NullableProjectPropertiesModel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableProjectPropertiesModel(val *ProjectPropertiesModel) *NullableProjectPropertiesModel {
	return &NullableProjectPropertiesModel{value: val, isSet: true}
}

func (v NullableProjectPropertiesModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableProjectPropertiesModel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


