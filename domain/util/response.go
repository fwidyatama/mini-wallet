package util

import (
	"encoding/json"
	"mini-wallet/domain/constant"
	"net/http"
)

type ResponseAPI struct {
	Status  string      `json:"status,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func SuccessResponseWriter(rw http.ResponseWriter, data interface{}, statusCode int) {
	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(ResponseAPI{
		Status: constant.SuccessMessage,
		Data:   data,
	})
}

func FailedResponseWriter(rw http.ResponseWriter, error string, statusCode int) {
	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(ResponseAPI{
		Status: constant.FailMessage,
		Data:   ErrorResponse{Error: error},
	})
}

func ResponseWriter(rw http.ResponseWriter, data interface{}, statusCode int) {
	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(data)
}
