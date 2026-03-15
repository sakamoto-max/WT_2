package utils

import (
	"encoding/json"
	"net/http"
	"tracker_service/internal/user"
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

func ValidationErrWriter(w http.ResponseWriter, errs *[]user.ValidationErr) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errs)
}
