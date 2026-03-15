package customerrors

import (
	"errors"
)

type myErrors string

type MyErrors struct{
	ErrCode myErrors `json:"error code"` 
	Message string `json:"message"`
}


var (
	EmailDoesntExist = errors.New("email doesnot exist")
	ErrUserAlreadyExists = errors.New("user already exists")
	PleaseSignUp = errors.New("email not found. please signup first")
	ErrRefreshTokenExp = errors.New("the refresh token is expired. please login again")
)

// db errors
var (
	ErrUserNotfound      = errors.New("User not found")
	ErrNameAlreadyExits  = errors.New("user with this name already exits")
	ErrEmailAlreadyExits = errors.New("user with this email already exits")
)

