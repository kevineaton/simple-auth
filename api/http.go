package api

import (
	"encoding/json"
	"net/http"
)

// APIReturn represents a standard API return object
type APIReturn struct {
	Data interface{} `json:"data,omitempty"`
}

// APIError represents the data key of an error for the API
type APIError struct {
	Code    string      `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// Send standardizes the return from the API
func Send(w http.ResponseWriter, code int, payload interface{}) {
	ret := APIReturn{}
	ret.Data = payload
	response, _ := json.Marshal(ret)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// SendError sends an error to the client
func SendError(w *http.ResponseWriter, r *http.Request, status int, systemCode string, message string, data *map[string]interface{}) {
	if data == nil {
		data = &map[string]interface{}{}
	}

	// TODO: we could integrate sentry and logging here

	Send(*w, status, APIError{
		Code:    systemCode,
		Message: message,
	})
}
