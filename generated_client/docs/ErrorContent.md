# ErrorContent

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Code** | Pointer to **string** | Application-level error code that describes how the error should be handled and how the &#x60;detail&#x60; field should be interpreted:   * &#x60;HTTP_ERROR&#x60; - The error should be handled based on the HTTP status code (&#x60;status&#x60;) of the response according to RFC-9910, and &#x60;detail&#x60; should be interpreted as a human-readable string.   * &#x60;CONCURRENT_UPDATE&#x60; - A concurrent update caused the &#x60;PUT&#x60; or &#x60;PATCH&#x60; request to fail. A &#x60;PUT&#x60; request can be retried after using &#x60;GET&#x60; to obtain the latest resource version and applying the desired change to it. A &#x60;PATCH&#x60; request can be retried without any changes to the request content. | [optional] 
**Status** | Pointer to **string** | HTTP status code and reason | [optional] 
**Detail** | Pointer to **string** | Detail about the error | [optional] 

## Methods

### NewErrorContent

`func NewErrorContent() *ErrorContent`

NewErrorContent instantiates a new ErrorContent object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewErrorContentWithDefaults

`func NewErrorContentWithDefaults() *ErrorContent`

NewErrorContentWithDefaults instantiates a new ErrorContent object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCode

`func (o *ErrorContent) GetCode() string`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *ErrorContent) GetCodeOk() (*string, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *ErrorContent) SetCode(v string)`

SetCode sets Code field to given value.

### HasCode

`func (o *ErrorContent) HasCode() bool`

HasCode returns a boolean if a field has been set.

### GetStatus

`func (o *ErrorContent) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ErrorContent) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ErrorContent) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ErrorContent) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetDetail

`func (o *ErrorContent) GetDetail() string`

GetDetail returns the Detail field if non-nil, zero value otherwise.

### GetDetailOk

`func (o *ErrorContent) GetDetailOk() (*string, bool)`

GetDetailOk returns a tuple with the Detail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetail

`func (o *ErrorContent) SetDetail(v string)`

SetDetail sets Detail field to given value.

### HasDetail

`func (o *ErrorContent) HasDetail() bool`

HasDetail returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


