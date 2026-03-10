package utils

import (
	"encoding/json"
	"net/http"
)

func ResponseWriter(w http.ResponseWriter, repsonse any, statusCode int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(repsonse)
}


func ErrorWriter(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err.Error())
}