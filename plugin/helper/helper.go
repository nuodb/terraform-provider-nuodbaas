/*
(C) Copyright 2016-2023 Dassault Systemes SE.
All Rights Reserved.
*/
package helper

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	nuodbaas "github.com/nuodb/nuodbaas-tf-plugin/generated_client"
)

// Return a generic error message if the provided provider attribute is missing
func GetProviderValidatorErrorMessage(valueType string, envVariable string) string {
	return fmt.Sprintf("The provider cannot create the NuoDbaas API client as there is a missing or empty value for the NuoDbaas API %v. "+
		"Set the %v value in the configuration or use the %v environment variable. "+
		"If either is already set, ensure the value is not empty.", valueType, valueType, envVariable)
}


//Return a Error Content object
func GetErrorContentObj(err error) *nuodbaas.ErrorContent {
	if serverErr, ok := err.(*nuodbaas.GenericOpenAPIError); ok {
		errorModel := nuodbaas.ErrorContent{}
		json.Unmarshal(serverErr.Body(), &errorModel)
		return &errorModel
	} else {
		errorMessage := err.Error()
		code := ""
		status := ""
		return &nuodbaas.ErrorContent{
			Detail: &errorMessage,
			Code: &code,
			Status: &status,
		}
	}
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