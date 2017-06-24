package utility

import (
	"encoding/json"
	"net/http"
)

// HTTPUtilityFunctions ... utility functions for doing basic HTTP operations
type HTTPUtilityFunctions struct {
}

//DecodeHTTPBody ...
func (huf HTTPUtilityFunctions) DecodeHTTPBody(req *http.Request) (map[string]interface{}, error) {
	decoder := json.NewDecoder(req.Body)
	var queryDict map[string]interface{}
	err := decoder.Decode(&queryDict)
	if err == nil {
		return queryDict, nil
	}
	return nil, err
}
