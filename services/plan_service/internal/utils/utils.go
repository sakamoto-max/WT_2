package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"plan_service/internal/models"
)

func OkResponseWriter(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
func CreatedResponseWriter(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func BadReqErrorWriter(w http.ResponseWriter, err string) {

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(err)

}
func BadReqErrorWriterResp(w http.ResponseWriter, err string) {

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(err)

}

func InternalServerErrorWriter(w http.ResponseWriter, err string) {

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(err)
}

func MakeJSON(msg models.MqMsg) ([]byte, error) {
	dataInBytes, err := json.Marshal(msg)
	if err != nil {
		return dataInBytes, fmt.Errorf("error occured while making json : %w", err)
	}

	return dataInBytes, nil

}
