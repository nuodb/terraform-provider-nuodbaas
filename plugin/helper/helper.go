package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"terraform-provider-nuodbaas/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func GetProviderValidatorErrorMessage(valueType string, envVariable string) string {
	return fmt.Sprintf("The provider cannot create the NuoDbaas API client as there is a missing or empty value for the NuoDbaas API %v. "+
		"Set the %v value in the configuration or use the %v environment variable. "+
		"If either is already set, ensure the value is not empty.", valueType, valueType, envVariable)
}

func GetHttpResponseObj(httpResponse *http.Response, target interface{}) error {
	defer httpResponse.Body.Close()
	return json.NewDecoder(httpResponse.Body).Decode(target)
}

func GetHttpResponseModel(httpResponse *http.Response) *model.ErrorModel {
	errorModel := &model.ErrorModel{}
	errorObj := GetHttpResponseObj(httpResponse, errorModel)
	if errorObj != nil {
		return nil
	} 
	return errorModel
}


func GetHttpResponseErrorMessage(httpResponse *http.Response, err error) string {
	errorModel := &model.ErrorModel{}
	errorObj := GetHttpResponseObj(httpResponse, errorModel)

	if errorObj != nil {
		return err.Error()
	} 
	return errorModel.Detail
}


func IsTimeoutError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
	  return true
	}
   return false
}

func RemoveDoubleQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

func ConvertMapToTfMap(mapObj *map[string]string) (basetypes.MapValue, diag.Diagnostics){
	tierParameters := map[string]attr.Value{}
	for k,v := range *mapObj {
		tierParameters[k] = types.StringValue(v)
	}
	mapValue, diags := types.MapValue(types.StringType, tierParameters)
	return mapValue, diags
}