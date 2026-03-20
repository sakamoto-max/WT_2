package utils

import (
	"auth_service/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	myErrs "wt/pkg/my_errors"
)

func BadReq(w http.ResponseWriter, err error) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error" : err.Error(),
	})
}

func InternalServerErr(w http.ResponseWriter, err error) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{
		"error" : "server encountered a problem",
	})	
}

func CreatedRespWriter(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func HashThePassword(password string) (string, error) {
	passInBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error encrypting the password : %w", err)
	}
	return string(passInBytes), nil
}

func MatchPasswords(password string, passFromDb string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passFromDb), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword){
			return myErrs.ErrIncorrectPassword
		}
		return fmt.Errorf("some error occured while authenticating the password : %w", err)
	}
	return nil
}
func ErrorWriter(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err.Error())
}

func ResponseWriter(w http.ResponseWriter, repsonse any, statusCode int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(repsonse)
}

func ValidationErrWriter(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resp)
}

func RequestTimeoutErrWriter(w http.ResponseWriter) {
	w.WriteHeader(http.StatusRequestTimeout)
}

func MakeJSON(msg models.MqMsg) ([]byte, error) {
	dataInBytes, err := json.Marshal(msg)
	if err != nil {
		return dataInBytes, fmt.Errorf("error occured while making json : %w", err)
	}

	return dataInBytes, nil

}

