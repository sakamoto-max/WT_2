package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func InternalServerErr(w http.ResponseWriter, err error) {
	log.Printf("error occured : %v", err.Error())
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "server encountered a problem",
	})
}

func BadReq(w http.ResponseWriter, err error) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error" : err.Error(),
	})
}

func NotFoundErr(w http.ResponseWriter, err error) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error" : err.Error(),
	})
}


