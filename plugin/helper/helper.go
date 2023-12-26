/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

// Return a generic error message if the provided provider attribute is missing
func GetProviderValidatorErrorMessage(valueType string, envVariable string) string {
	return fmt.Sprintf("The provider cannot create the NuoDbaas API client as there is a missing or empty value for the NuoDbaas API %v. "+
		"Set the %v value in the configuration or use the %v environment variable. "+
		"If either is already set, ensure the value is not empty.", valueType, valueType, envVariable)
}

func GetErrorModelFromError(ctx context.Context, err error) *nuodbaas.ErrorContent {
	if serverErr, ok := err.(*nuodbaas.GenericOpenAPIError); ok {
		errorModel := nuodbaas.ErrorContent{}
		json.Unmarshal(serverErr.Body(), &errorModel)
		return &errorModel
	}	else {
		tflog.Debug(ctx, fmt.Sprintf("TAGGER not ok and is %+v", err))
		return nil
	}
}

func getHttpResponseObj(httpResponse *http.Response, target interface{}) error {
	defer httpResponse.Body.Close()
	return json.NewDecoder(httpResponse.Body).Decode(target)
}

// Extracts the Error Model from the http.Response struct
func GetHttpResponseModel(httpResponse *http.Response) *nuodbaas.ErrorContent {
	errorModel := &nuodbaas.ErrorContent{}
	errorObj := getHttpResponseObj(httpResponse, errorModel)
	if errorObj != nil {
		return nil
	} 
	return errorModel
}


// Returns the readable error message string provided by the client
func GetHttpResponseErrorMessage(httpResponse *http.Response, err error) string {
	errorModel := &nuodbaas.ErrorContent{}
	errorObj := getHttpResponseObj(httpResponse, errorModel)

	if errorObj != nil {
		return err.Error()
	} 
	return errorModel.GetDetail()
}


func IsTimeoutError(err error) bool {
	return os.IsTimeout(err)
}

// Removes any extra double quotes that are added in the string
func RemoveDoubleQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

// Converts map[string]string to basetypes.MapValue
func ConvertMapToTfMap(mapObj *map[string]string) (basetypes.MapValue, diag.Diagnostics){
	mapValue := map[string]attr.Value{}
	for k,v := range *mapObj {
		mapValue[k] = types.StringValue(v)
	}
	tfMapValue, diags := types.MapValue(types.StringType, mapValue)
	return tfMapValue, diags
}

func ComputeWaitTime(i int, maxWait int) int {
	if(i * 2 < maxWait) {
		return i * 2
	} else {
		return maxWait
	}
}