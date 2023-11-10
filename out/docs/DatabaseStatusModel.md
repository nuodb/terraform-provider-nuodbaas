# DatabaseStatusModel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SqlEndpoint** | Pointer to **string** | The endpoint for SQL clients to connect to | [optional] 
**CaPem** | Pointer to **string** | The PEM-encoded certificate for SQL clients to verify database servers | [optional] 
**Ready** | Pointer to **bool** | Whether the database is ready | [optional] 
**Shutdown** | Pointer to **bool** | Whether the database has shutdown | [optional] 
**Message** | Pointer to **string** | Message summarizing the state of the database | [optional] 

## Methods

### NewDatabaseStatusModel

`func NewDatabaseStatusModel() *DatabaseStatusModel`

NewDatabaseStatusModel instantiates a new DatabaseStatusModel object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDatabaseStatusModelWithDefaults

`func NewDatabaseStatusModelWithDefaults() *DatabaseStatusModel`

NewDatabaseStatusModelWithDefaults instantiates a new DatabaseStatusModel object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSqlEndpoint

`func (o *DatabaseStatusModel) GetSqlEndpoint() string`

GetSqlEndpoint returns the SqlEndpoint field if non-nil, zero value otherwise.

### GetSqlEndpointOk

`func (o *DatabaseStatusModel) GetSqlEndpointOk() (*string, bool)`

GetSqlEndpointOk returns a tuple with the SqlEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSqlEndpoint

`func (o *DatabaseStatusModel) SetSqlEndpoint(v string)`

SetSqlEndpoint sets SqlEndpoint field to given value.

### HasSqlEndpoint

`func (o *DatabaseStatusModel) HasSqlEndpoint() bool`

HasSqlEndpoint returns a boolean if a field has been set.

### GetCaPem

`func (o *DatabaseStatusModel) GetCaPem() string`

GetCaPem returns the CaPem field if non-nil, zero value otherwise.

### GetCaPemOk

`func (o *DatabaseStatusModel) GetCaPemOk() (*string, bool)`

GetCaPemOk returns a tuple with the CaPem field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCaPem

`func (o *DatabaseStatusModel) SetCaPem(v string)`

SetCaPem sets CaPem field to given value.

### HasCaPem

`func (o *DatabaseStatusModel) HasCaPem() bool`

HasCaPem returns a boolean if a field has been set.

### GetReady

`func (o *DatabaseStatusModel) GetReady() bool`

GetReady returns the Ready field if non-nil, zero value otherwise.

### GetReadyOk

`func (o *DatabaseStatusModel) GetReadyOk() (*bool, bool)`

GetReadyOk returns a tuple with the Ready field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReady

`func (o *DatabaseStatusModel) SetReady(v bool)`

SetReady sets Ready field to given value.

### HasReady

`func (o *DatabaseStatusModel) HasReady() bool`

HasReady returns a boolean if a field has been set.

### GetShutdown

`func (o *DatabaseStatusModel) GetShutdown() bool`

GetShutdown returns the Shutdown field if non-nil, zero value otherwise.

### GetShutdownOk

`func (o *DatabaseStatusModel) GetShutdownOk() (*bool, bool)`

GetShutdownOk returns a tuple with the Shutdown field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShutdown

`func (o *DatabaseStatusModel) SetShutdown(v bool)`

SetShutdown sets Shutdown field to given value.

### HasShutdown

`func (o *DatabaseStatusModel) HasShutdown() bool`

HasShutdown returns a boolean if a field has been set.

### GetMessage

`func (o *DatabaseStatusModel) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *DatabaseStatusModel) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *DatabaseStatusModel) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *DatabaseStatusModel) HasMessage() bool`

HasMessage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


