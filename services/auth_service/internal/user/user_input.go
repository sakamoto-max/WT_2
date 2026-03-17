package user

import (
	"errors"
)

type validatonErrs struct {
	Path    string
	Message string
}

var (
	ErrValidationErr = errors.New("validation error occured")
)

var (
	errUserNameReq     = validatonErrs{Path: "name", Message: "required"}
	errUserEmailReq    = validatonErrs{Path: "email", Message: "required"}
	errUserPassWordReq = validatonErrs{Path: "password", Message: "required"}
	errUUIDReq         = validatonErrs{Path: "UUID", Message: "required"}
	// errMinUserName = validatonErrs{Path: "name", Message: "the min required length of the name is 2"}
	// errMaxUserName = validatonErrs{Path: "name", Message: "the max length of the name is 20"}
)

type Signup struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Signup) Validate() (*[]validatonErrs, bool) {

	var validationErrs []validatonErrs

	if s.Name == "" {
		validationErrs = append(validationErrs, errUserNameReq)
	}

	if s.Email == "" {
		validationErrs = append(validationErrs, errUserEmailReq)
	}

	if s.Password == "" {
		validationErrs = append(validationErrs, errUserPassWordReq)
	}

	if len(validationErrs) > 0 {
		return &validationErrs, true
	}

	return &validationErrs, false
}

type Login struct {
	Message  string `json:"message"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *Login) Validate() (*[]validatonErrs, bool) {

	var validationErrs []validatonErrs
	if l.Email == "" {
		validationErrs = append(validationErrs, errUserEmailReq)
	}

	if l.Password == "" {
		validationErrs = append(validationErrs, errUserPassWordReq)
	}

	if len(validationErrs) > 0 {
		return &validationErrs, true
	}

	return &validationErrs, false
}

type UUIDReader struct {
	UUID string `json:"uuid"`
}

func (u *UUIDReader) Validate() (*[]validatonErrs, bool) {

	var validationErrs []validatonErrs

	// uuid := u.UUID.String()

	if u.UUID == "" {
		validationErrs = append(validationErrs, errUUIDReq)
		return &validationErrs, true
	}

	return &validationErrs, false

}
