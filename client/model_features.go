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
	"fmt"
)

// checks if the Features type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Features{}

// Features struct for Features
type Features struct {
	// The name of the resource.
	Name string `json:"name"`
	// The namespace of the resource. When not specified, the current namespace is assumed.
	Namespace *string `json:"namespace,omitempty"`
}

type _Features Features

// NewFeatures instantiates a new Features object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFeatures(name string) *Features {
	this := Features{}
	this.Name = name
	return &this
}

// NewFeaturesWithDefaults instantiates a new Features object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFeaturesWithDefaults() *Features {
	this := Features{}
	return &this
}

// GetName returns the Name field value
func (o *Features) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *Features) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *Features) SetName(v string) {
	o.Name = v
}

// GetNamespace returns the Namespace field value if set, zero value otherwise.
func (o *Features) GetNamespace() string {
	if o == nil || IsNil(o.Namespace) {
		var ret string
		return ret
	}
	return *o.Namespace
}

// GetNamespaceOk returns a tuple with the Namespace field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Features) GetNamespaceOk() (*string, bool) {
	if o == nil || IsNil(o.Namespace) {
		return nil, false
	}
	return o.Namespace, true
}

// HasNamespace returns a boolean if a field has been set.
func (o *Features) HasNamespace() bool {
	if o != nil && !IsNil(o.Namespace) {
		return true
	}

	return false
}

// SetNamespace gets a reference to the given string and assigns it to the Namespace field.
func (o *Features) SetNamespace(v string) {
	o.Namespace = &v
}

func (o Features) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Features) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["name"] = o.Name
	if !IsNil(o.Namespace) {
		toSerialize["namespace"] = o.Namespace
	}
	return toSerialize, nil
}

func (o *Features) UnmarshalJSON(bytes []byte) (err error) {
    // This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"name",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(bytes, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varFeatures := _Features{}

	err = json.Unmarshal(bytes, &varFeatures)

	if err != nil {
		return err
	}

	*o = Features(varFeatures)

	return err
}

type NullableFeatures struct {
	value *Features
	isSet bool
}

func (v NullableFeatures) Get() *Features {
	return v.value
}

func (v *NullableFeatures) Set(val *Features) {
	v.value = val
	v.isSet = true
}

func (v NullableFeatures) IsSet() bool {
	return v.isSet
}

func (v *NullableFeatures) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFeatures(val *Features) *NullableFeatures {
	return &NullableFeatures{value: val, isSet: true}
}

func (v NullableFeatures) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFeatures) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

