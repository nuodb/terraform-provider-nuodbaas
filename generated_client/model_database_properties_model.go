/*
NuoDB Control Plane REST API

NuoDB Control Plane (CP) allows users to create and manage NuoDB databases remotely using a Database as a Service (DBaaS) model.

API version: 2.2.0
Contact: NuoDB.Support@3ds.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package nuodbaas

import (
	"encoding/json"
)

// checks if the DatabasePropertiesModel type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DatabasePropertiesModel{}

// DatabasePropertiesModel struct for DatabasePropertiesModel
type DatabasePropertiesModel struct {
	// The size of the archive volumes for the database. Can be only updated to increase the volume size.
	ArchiveDiskSize *string `json:"archiveDiskSize,omitempty"`
	// The size of the journal volumes for the database. Can be only updated to increase the volume size.
	JournalDiskSize *string `json:"journalDiskSize,omitempty"`
}

// NewDatabasePropertiesModel instantiates a new DatabasePropertiesModel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDatabasePropertiesModel() *DatabasePropertiesModel {
	this := DatabasePropertiesModel{}
	return &this
}

// NewDatabasePropertiesModelWithDefaults instantiates a new DatabasePropertiesModel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDatabasePropertiesModelWithDefaults() *DatabasePropertiesModel {
	this := DatabasePropertiesModel{}
	return &this
}

// GetArchiveDiskSize returns the ArchiveDiskSize field value if set, zero value otherwise.
func (o *DatabasePropertiesModel) GetArchiveDiskSize() string {
	if o == nil || IsNil(o.ArchiveDiskSize) {
		var ret string
		return ret
	}
	return *o.ArchiveDiskSize
}

// GetArchiveDiskSizeOk returns a tuple with the ArchiveDiskSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DatabasePropertiesModel) GetArchiveDiskSizeOk() (*string, bool) {
	if o == nil || IsNil(o.ArchiveDiskSize) {
		return nil, false
	}
	return o.ArchiveDiskSize, true
}

// HasArchiveDiskSize returns a boolean if a field has been set.
func (o *DatabasePropertiesModel) HasArchiveDiskSize() bool {
	if o != nil && !IsNil(o.ArchiveDiskSize) {
		return true
	}

	return false
}

// SetArchiveDiskSize gets a reference to the given string and assigns it to the ArchiveDiskSize field.
func (o *DatabasePropertiesModel) SetArchiveDiskSize(v string) {
	o.ArchiveDiskSize = &v
}

// GetJournalDiskSize returns the JournalDiskSize field value if set, zero value otherwise.
func (o *DatabasePropertiesModel) GetJournalDiskSize() string {
	if o == nil || IsNil(o.JournalDiskSize) {
		var ret string
		return ret
	}
	return *o.JournalDiskSize
}

// GetJournalDiskSizeOk returns a tuple with the JournalDiskSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DatabasePropertiesModel) GetJournalDiskSizeOk() (*string, bool) {
	if o == nil || IsNil(o.JournalDiskSize) {
		return nil, false
	}
	return o.JournalDiskSize, true
}

// HasJournalDiskSize returns a boolean if a field has been set.
func (o *DatabasePropertiesModel) HasJournalDiskSize() bool {
	if o != nil && !IsNil(o.JournalDiskSize) {
		return true
	}

	return false
}

// SetJournalDiskSize gets a reference to the given string and assigns it to the JournalDiskSize field.
func (o *DatabasePropertiesModel) SetJournalDiskSize(v string) {
	o.JournalDiskSize = &v
}

func (o DatabasePropertiesModel) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DatabasePropertiesModel) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ArchiveDiskSize) {
		toSerialize["archiveDiskSize"] = o.ArchiveDiskSize
	}
	if !IsNil(o.JournalDiskSize) {
		toSerialize["journalDiskSize"] = o.JournalDiskSize
	}
	return toSerialize, nil
}

type NullableDatabasePropertiesModel struct {
	value *DatabasePropertiesModel
	isSet bool
}

func (v NullableDatabasePropertiesModel) Get() *DatabasePropertiesModel {
	return v.value
}

func (v *NullableDatabasePropertiesModel) Set(val *DatabasePropertiesModel) {
	v.value = val
	v.isSet = true
}

func (v NullableDatabasePropertiesModel) IsSet() bool {
	return v.isSet
}

func (v *NullableDatabasePropertiesModel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDatabasePropertiesModel(val *DatabasePropertiesModel) *NullableDatabasePropertiesModel {
	return &NullableDatabasePropertiesModel{value: val, isSet: true}
}

func (v NullableDatabasePropertiesModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDatabasePropertiesModel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


