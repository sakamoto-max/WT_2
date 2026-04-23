package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"wt/pkg/types"
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
		"error": err.Error(),
	})
}

func NotFoundErr(w http.ResponseWriter, err error) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
func DeletedNotFoundWriter(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(resp)
}



func CreatedWriter(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func OkRespWriter(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func ConflictWriter(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(resp)
}

func ConvertIntoBytes(payload any) (*[]byte, error) {

	dataInBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error in converting data into bytes : %w", err)
	}

	return &dataInBytes, nil
}

func ConvertIntoJosn(data *[]byte) *types.Data {

	var D types.Data

	_ = json.Unmarshal(*data, &D)

	return &D
}

func MakeJSONV2(msg any) (*string, error) {

	var data string

	dataInBytes, err := json.Marshal(msg)
	if err != nil{
		return nil, fmt.Errorf("error occured while making json : %w", err)
	}

	data = string(dataInBytes)

	return &data, nil
}
 