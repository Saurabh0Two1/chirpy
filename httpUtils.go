package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorResp struct {
	Error string `json:"error"`
}

func returnError(w http.ResponseWriter, statusCode int, message string, e error) {

	respBody := errorResp{
		Error: fmt.Sprintf("%s: %s", message, e),
	}

	jsonData, _ := json.Marshal(respBody)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}
