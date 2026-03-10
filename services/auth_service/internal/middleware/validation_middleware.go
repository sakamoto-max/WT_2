package middleware

import (
	"auth_service/internal/models"
	"context"
	"encoding/json"
	"net/http"
)

type inputKey string

var (
	UserSignUp       inputKey = "USER_SIGNUP"
	UserLogin        inputKey = "USER_LOGIN"
	RefreshTokenUUID inputKey = "REFRESH_TOKEN"
)

type validatonErrs struct {
	Path    string
	Message string
}

var (
	errUserNameReq     = validatonErrs{Path: "name", Message: "required"}
	errUserEmailReq    = validatonErrs{Path: "email", Message: "required"}
	errUserPassWordReq = validatonErrs{Path: "password", Message: "required"}
	errUUIDReq         = validatonErrs{Path: "UUID", Message: "required"}
	// errMinUserName = validatonErrs{Path: "name", Message: "the min required length of the name is 2"}
	// errMaxUserName = validatonErrs{Path: "name", Message: "the max length of the name is 20"}
)

func UserLoginValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var userInput models.Login
		var validationErrs []validatonErrs

		json.NewDecoder(r.Body).Decode(&userInput)

		if userInput.Email == "" {
			validationErrs = append(validationErrs, errUserEmailReq)
		}

		if userInput.Password == "" {
			validationErrs = append(validationErrs, errUserPassWordReq)
		}

		if len(validationErrs) > 0 {
			validationErrWriter(w, validationErrs)
			return
		}

		ctx := context.WithValue(r.Context(), UserLogin, &userInput)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func UserSignUpValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// decode the body
		var userSentInput models.Signup
		var validationErrs []validatonErrs

		json.NewDecoder(r.Body).Decode(&userSentInput)

		if userSentInput.Name == "" {
			validationErrs = append(validationErrs, errUserNameReq)
		}

		if userSentInput.Email == "" {
			validationErrs = append(validationErrs, errUserEmailReq)
		}

		if userSentInput.Password == "" {
			validationErrs = append(validationErrs, errUserPassWordReq)
		}

		if len(validationErrs) > 0 {
			validationErrWriter(w, validationErrs)
			return
		}

		ctx := context.WithValue(r.Context(), UserSignUp, &userSentInput)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func NewRefreshTokenValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var userInput models.UUIDReader
		var validationErrs []validatonErrs
		var uuid string

		json.NewDecoder(r.Body).Decode(&userInput)

		uuid = userInput.UUID.String()

		if uuid == "" {
			validationErrs = append(validationErrs, errUUIDReq)
			validationErrWriter(w, validationErrs)
			return
		}

		ctx := context.WithValue(r.Context(), RefreshTokenUUID, &userInput)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func validationErrWriter(w http.ResponseWriter, err any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(err)
}
