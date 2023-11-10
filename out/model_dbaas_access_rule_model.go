/*
NuoDB Control Plane REST API

NuoDB Control Plane (CP) allows users to create and manage NuoDB databases remotely using a Database as a Service (DBaaS) model.

API version: 2.2.0
Contact: NuoDB.Support@3ds.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the DbaasAccessRuleModel type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DbaasAccessRuleModel{}

// DbaasAccessRuleModel The rule specifying access for the user
type DbaasAccessRuleModel struct {
	// List of access rule entries in the form `<verb>:<resource specifier>[:<SLA>]` that specify requests to allow
	Allow []string `json:"allow,omitempty"`
	// List of access rule entries in the form `<verb>:<resource specifier>` that specify requests to deny
	Deny []string `json:"deny,omitempty"`
}

// NewDbaasAccessRuleModel instantiates a new DbaasAccessRuleModel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDbaasAccessRuleModel() *DbaasAccessRuleModel {
	this := DbaasAccessRuleModel{}
	return &this
}

// NewDbaasAccessRuleModelWithDefaults instantiates a new DbaasAccessRuleModel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDbaasAccessRuleModelWithDefaults() *DbaasAccessRuleModel {
	this := DbaasAccessRuleModel{}
	return &this
}

// GetAllow returns the Allow field value if set, zero value otherwise.
func (o *DbaasAccessRuleModel) GetAllow() []string {
	if o == nil || IsNil(o.Allow) {
		var ret []string
		return ret
	}
	return o.Allow
}

// GetAllowOk returns a tuple with the Allow field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DbaasAccessRuleModel) GetAllowOk() ([]string, bool) {
	if o == nil || IsNil(o.Allow) {
		return nil, false
	}
	return o.Allow, true
}

// HasAllow returns a boolean if a field has been set.
func (o *DbaasAccessRuleModel) HasAllow() bool {
	if o != nil && !IsNil(o.Allow) {
		return true
	}

	return false
}

// SetAllow gets a reference to the given []string and assigns it to the Allow field.
func (o *DbaasAccessRuleModel) SetAllow(v []string) {
	o.Allow = v
}

// GetDeny returns the Deny field value if set, zero value otherwise.
func (o *DbaasAccessRuleModel) GetDeny() []string {
	if o == nil || IsNil(o.Deny) {
		var ret []string
		return ret
	}
	return o.Deny
}

// GetDenyOk returns a tuple with the Deny field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DbaasAccessRuleModel) GetDenyOk() ([]string, bool) {
	if o == nil || IsNil(o.Deny) {
		return nil, false
	}
	return o.Deny, true
}

// HasDeny returns a boolean if a field has been set.
func (o *DbaasAccessRuleModel) HasDeny() bool {
	if o != nil && !IsNil(o.Deny) {
		return true
	}

	return false
}

// SetDeny gets a reference to the given []string and assigns it to the Deny field.
func (o *DbaasAccessRuleModel) SetDeny(v []string) {
	o.Deny = v
}

func (o DbaasAccessRuleModel) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DbaasAccessRuleModel) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Allow) {
		toSerialize["allow"] = o.Allow
	}
	if !IsNil(o.Deny) {
		toSerialize["deny"] = o.Deny
	}
	return toSerialize, nil
}

type NullableDbaasAccessRuleModel struct {
	value *DbaasAccessRuleModel
	isSet bool
}

func (v NullableDbaasAccessRuleModel) Get() *DbaasAccessRuleModel {
	return v.value
}

func (v *NullableDbaasAccessRuleModel) Set(val *DbaasAccessRuleModel) {
	v.value = val
	v.isSet = true
}

func (v NullableDbaasAccessRuleModel) IsSet() bool {
	return v.isSet
}

func (v *NullableDbaasAccessRuleModel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDbaasAccessRuleModel(val *DbaasAccessRuleModel) *NullableDbaasAccessRuleModel {
	return &NullableDbaasAccessRuleModel{value: val, isSet: true}
}

func (v NullableDbaasAccessRuleModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDbaasAccessRuleModel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


