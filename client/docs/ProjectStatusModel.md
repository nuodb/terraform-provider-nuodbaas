# ProjectStatusModel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CaPem** | Pointer to **string** | The PEM-encoded certificate for SQL clients to verify database servers within the project | [optional] 
**Ready** | Pointer to **bool** | Whether the project is ready | [optional] 
**Shutdown** | Pointer to **bool** | Whether the project and all of its databases have shutdown | [optional] 
**Message** | Pointer to **string** | Message summarizing the state of the project | [optional] 
**State** | Pointer to **string** | The state of the project:   * &#x60;Available&#x60; - The project is available   * &#x60;Creating&#x60; - The project is being created and not yet available   * &#x60;Modifying&#x60; - The project is being modified   * &#x60;Stopping&#x60; - Shutdown is in progress for this project   * &#x60;Stopped&#x60; - The project and its databases have been stopped   * &#x60;Expired&#x60; - The project and its databases have expired   * &#x60;Failed&#x60; - The project has failed to achieve a usable state   * &#x60;Deleting&#x60; - The project has been marked for deletion, which is in progress | [optional] 

## Methods

### NewProjectStatusModel

`func NewProjectStatusModel() *ProjectStatusModel`

NewProjectStatusModel instantiates a new ProjectStatusModel object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectStatusModelWithDefaults

`func NewProjectStatusModelWithDefaults() *ProjectStatusModel`

NewProjectStatusModelWithDefaults instantiates a new ProjectStatusModel object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCaPem

`func (o *ProjectStatusModel) GetCaPem() string`

GetCaPem returns the CaPem field if non-nil, zero value otherwise.

### GetCaPemOk

`func (o *ProjectStatusModel) GetCaPemOk() (*string, bool)`

GetCaPemOk returns a tuple with the CaPem field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCaPem

`func (o *ProjectStatusModel) SetCaPem(v string)`

SetCaPem sets CaPem field to given value.

### HasCaPem

`func (o *ProjectStatusModel) HasCaPem() bool`

HasCaPem returns a boolean if a field has been set.

### GetReady

`func (o *ProjectStatusModel) GetReady() bool`

GetReady returns the Ready field if non-nil, zero value otherwise.

### GetReadyOk

`func (o *ProjectStatusModel) GetReadyOk() (*bool, bool)`

GetReadyOk returns a tuple with the Ready field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReady

`func (o *ProjectStatusModel) SetReady(v bool)`

SetReady sets Ready field to given value.

### HasReady

`func (o *ProjectStatusModel) HasReady() bool`

HasReady returns a boolean if a field has been set.

### GetShutdown

`func (o *ProjectStatusModel) GetShutdown() bool`

GetShutdown returns the Shutdown field if non-nil, zero value otherwise.

### GetShutdownOk

`func (o *ProjectStatusModel) GetShutdownOk() (*bool, bool)`

GetShutdownOk returns a tuple with the Shutdown field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShutdown

`func (o *ProjectStatusModel) SetShutdown(v bool)`

SetShutdown sets Shutdown field to given value.

### HasShutdown

`func (o *ProjectStatusModel) HasShutdown() bool`

HasShutdown returns a boolean if a field has been set.

### GetMessage

`func (o *ProjectStatusModel) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ProjectStatusModel) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ProjectStatusModel) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *ProjectStatusModel) HasMessage() bool`

HasMessage returns a boolean if a field has been set.

### GetState

`func (o *ProjectStatusModel) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *ProjectStatusModel) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *ProjectStatusModel) SetState(v string)`

SetState sets State field to given value.

### HasState

`func (o *ProjectStatusModel) HasState() bool`

HasState returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

