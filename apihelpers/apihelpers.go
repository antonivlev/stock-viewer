/*
Functions for writing data into the response
TODO: add tests
*/
package apihelpers

import (
	"encoding/json"
	"net/http"
)

// Writes message, error text and Bad Request code. If err is nil, just writes message.
func WriteError(w http.ResponseWriter, message string, err error) {
	errString := ""
	if err != nil {
		errString = err.Error()
	}
	http.Error(w, message+"\n\n"+errString, http.StatusBadRequest)
}

// Writes apiResponse as json to w, sets Bad Request code
func WriteErrorResponse(w http.ResponseWriter, apiResponse map[string]interface{}) {
	// parsed response, but it constains error; return it
	errResponseBytes, errMarshal := json.Marshal(apiResponse)
	// TODO: tedious error handling, any way to guarantee correctness?
	if errMarshal != nil {
		WriteError(w, "Error encoding api error response", errMarshal)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(errResponseBytes)
}
