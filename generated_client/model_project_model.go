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
	"fmt"
)

// checks if the ProjectModel type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ProjectModel{}

// ProjectModel struct for ProjectModel
type ProjectModel struct {
	Organization *string `json:"organization,omitempty"`
	Name *string `json:"name,omitempty"`
	// The SLA for the project. Cannot be updated once the project is created.
	Sla string `json:"sla"`
	// The service tier for the project
	Tier string `json:"tier"`
	Maintenance *MaintenanceModel `json:"maintenance,omitempty"`
	// The version of the resource. When specified in a `PUT` request payload, indicates that the resoure should be updated, and is used by the system to guard against concurrent updates.
	ResourceVersion *string `json:"resourceVersion,omitempty"`
	Status *ProjectStatusModel `json:"status,omitempty"`
}

type _ProjectModel ProjectModel

// NewProjectModel instantiates a new ProjectModel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewProjectModel(sla string, tier string) *ProjectModel {
	this := ProjectModel{}
	this.Sla = sla
	this.Tier = tier
	return &this
}

// NewProjectModelWithDefaults instantiates a new ProjectModel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewProjectModelWithDefaults() *ProjectModel {
	this := ProjectModel{}
	return &this
}

// GetOrganization returns the Organization field value if set, zero value otherwise.
func (o *ProjectModel) GetOrganization() string {
	if o == nil || IsNil(o.Organization) {
		var ret string
		return ret
	}
	return *o.Organization
}

// GetOrganizationOk returns a tuple with the Organization field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectModel) GetOrganizationOk() (*string, bool) {
	if o == nil || IsNil(o.Organization) {
		return nil, false
	}
	return o.Organization, true
}

// HasOrganization returns a boolean if a field has been set.
func (o *ProjectModel) HasOrganization() bool {
	if o != nil && !IsNil(o.Organization) {
		return true
	}

	return false
}

// SetOrganization gets a reference to the given string and assigns it to the Organization field.
func (o *ProjectModel) SetOrganization(v string) {
	o.Organization = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ProjectModel) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectModel) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ProjectModel) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ProjectModel) SetName(v string) {
	o.Name = &v
}

// GetSla returns the Sla field value
func (o *ProjectModel) GetSla() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Sla
}

// GetSlaOk returns a tuple with the Sla field value
// and a boolean to check if the value has been set.
func (o *ProjectModel) GetSlaOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Sla, true
}

// SetSla sets field value
func (o *ProjectModel) SetSla(v string) {
	o.Sla = v
}

// GetTier returns the Tier field value
func (o *ProjectModel) GetTier() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Tier
}

// GetTierOk returns a tuple with the Tier field value
// and a boolean to check if the value has been set.
func (o *ProjectModel) GetTierOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Tier, true
}

// SetTier sets field value
func (o *ProjectModel) SetTier(v string) {
	o.Tier = v
}

// GetMaintenance returns the Maintenance field value if set, zero value otherwise.
func (o *ProjectModel) GetMaintenance() MaintenanceModel {
	if o == nil || IsNil(o.Maintenance) {
		var ret MaintenanceModel
		return ret
	}
	return *o.Maintenance
}

// GetMaintenanceOk returns a tuple with the Maintenance field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectModel) GetMaintenanceOk() (*MaintenanceModel, bool) {
	if o == nil || IsNil(o.Maintenance) {
		return nil, false
	}
	return o.Maintenance, true
}

// HasMaintenance returns a boolean if a field has been set.
func (o *ProjectModel) HasMaintenance() bool {
	if o != nil && !IsNil(o.Maintenance) {
		return true
	}

	return false
}

// SetMaintenance gets a reference to the given MaintenanceModel and assigns it to the Maintenance field.
func (o *ProjectModel) SetMaintenance(v MaintenanceModel) {
	o.Maintenance = &v
}

// GetResourceVersion returns the ResourceVersion field value if set, zero value otherwise.
func (o *ProjectModel) GetResourceVersion() string {
	if o == nil || IsNil(o.ResourceVersion) {
		var ret string
		return ret
	}
	return *o.ResourceVersion
}

// GetResourceVersionOk returns a tuple with the ResourceVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectModel) GetResourceVersionOk() (*string, bool) {
	if o == nil || IsNil(o.ResourceVersion) {
		return nil, false
	}
	return o.ResourceVersion, true
}

// HasResourceVersion returns a boolean if a field has been set.
func (o *ProjectModel) HasResourceVersion() bool {
	if o != nil && !IsNil(o.ResourceVersion) {
		return true
	}

	return false
}

// SetResourceVersion gets a reference to the given string and assigns it to the ResourceVersion field.
func (o *ProjectModel) SetResourceVersion(v string) {
	o.ResourceVersion = &v
}

// GetStatus returns the Status field value if set, zero value otherwise.
func (o *ProjectModel) GetStatus() ProjectStatusModel {
	if o == nil || IsNil(o.Status) {
		var ret ProjectStatusModel
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectModel) GetStatusOk() (*ProjectStatusModel, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}
	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *ProjectModel) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given ProjectStatusModel and assigns it to the Status field.
func (o *ProjectModel) SetStatus(v ProjectStatusModel) {
	o.Status = &v
}

func (o ProjectModel) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ProjectModel) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Organization) {
		toSerialize["organization"] = o.Organization
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	toSerialize["sla"] = o.Sla
	toSerialize["tier"] = o.Tier
	if !IsNil(o.Maintenance) {
		toSerialize["maintenance"] = o.Maintenance
	}
	if !IsNil(o.ResourceVersion) {
		toSerialize["resourceVersion"] = o.ResourceVersion
	}
	if !IsNil(o.Status) {
		toSerialize["status"] = o.Status
	}
	return toSerialize, nil
}

func (o *ProjectModel) UnmarshalJSON(bytes []byte) (err error) {
    // This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"sla",
		"tier",
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

	varProjectModel := _ProjectModel{}

	err = json.Unmarshal(bytes, &varProjectModel)

	if err != nil {
		return err
	}

	*o = ProjectModel(varProjectModel)

	return err
}

type NullableProjectModel struct {
	value *ProjectModel
	isSet bool
}

func (v NullableProjectModel) Get() *ProjectModel {
	return v.value
}

func (v *NullableProjectModel) Set(val *ProjectModel) {
	v.value = val
	v.isSet = true
}

func (v NullableProjectModel) IsSet() bool {
	return v.isSet
}

func (v *NullableProjectModel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableProjectModel(val *ProjectModel) *NullableProjectModel {
	return &NullableProjectModel{value: val, isSet: true}
}

func (v NullableProjectModel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableProjectModel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


