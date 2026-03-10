package customerrors

import (
	"errors"
)

// type MyErrors struct{
// 	ErrType string
// 	ErrName error
// 	StatusCode int
// 	Resp any
// }

type myErrors string

// var (
// 	errNoClaims myErrors = "ERROR_GETTING_CLAIMS"
// 	errUserNotFound myErrors = "ERROR_USER_NOT_FOUND"
// 	errUserAlreadyExists2 myErrors = "ERROR_USER_ALREADY_EXISTS"
// )


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


// var (
// 	ErrGettingClaims = MyErrors{ErrCode: ErrNoClaims, Message: "error getting claims from context"}
// 	UserAlrExists = MyErrors{ErrCode: ErrUserAlreadyExists2, Message: "user already exists."}
// )
