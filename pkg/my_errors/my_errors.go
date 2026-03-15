package myerrors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type AppErrs struct {
	Msg  string `json:"message"`
	err  error
	code int
}

func NewAppErr(err error, code int) *AppErrs {
	return &AppErrs{Msg: err.Error(), err: err, code: code}
}

func (appErr *AppErrs) Error() string {
	return appErr.Msg
}

func (a *AppErrs) AppErrWriter(w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(a.code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": a.Msg,
	})
}

// jwt errors
var (
	ErrTokenExpired   = errors.New("token is expired, get a new access token at /refresh")
	ErrTokenMalformed = errors.New("token is malformed. please check the token again")
	ErrTokenInvalid   = errors.New("token is invalid")
	ErrTokenIsMissing = errors.New("token is missing, please provide the token")
	ErrRefreshExpired = errors.New("referesh token is expired, please login again")
	ErrSignatureInvalid = errors.New("token's signature is invalid")
)

// user_service_errors
// db

var (
	ErrUserNotfound      = errors.New("User not found")
	ErrNameAlreadyExits  = errors.New("user with this name already exits")
	ErrEmailAlreadyExits = errors.New("user with this email already exits")
)

// user_login
var (
	ErrEmailNotFound = errors.New("email not found")
	ErrIncorrectPassword = errors.New("password is incorrect")
)




