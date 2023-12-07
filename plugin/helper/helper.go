package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"terraform-provider-nuodbaas/internal/model"
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