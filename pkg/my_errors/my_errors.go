package myerrors

import (
	"encoding/json"
	"errors"
	"log"

	// "fmt"
	// "log"
	"net/http"

	// "wt/pkg/utils"

	// "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	ErrTokenExpired     = errors.New("token is expired, get a new access token at /refresh")
	// ErrUserExits2       = status.Error(codes.AlreadyExists, "user already exits")
	ErrTokenMalformed   = errors.New("token is malformed. please check the token again")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrTokenIsMissing   = errors.New("token is missing, please provide the token")
	ErrRefreshExpired   = errors.New("referesh token is expired, please login again")
	ErrSignatureInvalid = errors.New("token's signature is invalid")
	ErrOldPassNewPassSame = errors.New("the old pass and new pass cannot be the same")
)

// claims error

var (
	ErrGettingClaims = errors.New("error getting claims from token")
)

var (
	ErrUserNotfound      = errors.New("User not found")
	ErrNameAlreadyExits  = errors.New("user with this name already exits")
	ErrEmailAlreadyExits = errors.New("user with this email already exits")
)

// user_login
var (
	ErrEmailNotFound     = errors.New("email not found")
	ErrIncorrectPassword = errors.New("password is incorrect")
	ErrOldEmailNewEmailSame = errors.New("old email and new email shouldn't be the same")
	ErrEmailDoesntMatch = errors.New("the email user sent is incorrect")
)

var (
	PlanServerNotResponding = errors.New("plan server is not responding")
)

var (
	ErrWorkoutOngoing = errors.New("user has ongoing workout which is not ended")
)

var (
	CodeInternalServer = 500
	CodeBadRequest = 400
	CodeNotFound = 404
)

func ErrMatcher(w http.ResponseWriter, err error) {
	st, _ := status.FromError(err)
	code := st.Code()
	switch code {
	case codes.AlreadyExists:
		appErr := &AppErrs{
			code: CodeBadRequest,
			Msg:  st.Message(),
		}
		appErr.AppErrWriter(w)
	case codes.Internal:
		log.Printf("error occured : %v", err)
		appErr := &AppErrs{
			code: CodeInternalServer,
			Msg:  st.Message(),
		}
		appErr.AppErrWriter(w)
	case codes.Canceled:
		log.Printf("error occured : %v", err.Error())
		appErr := &AppErrs{
			code: CodeInternalServer,
			Msg:  "server encountered a problem",
		}
		appErr.AppErrWriter(w)
	case codes.Unavailable:
		log.Printf("error occured : %v", err.Error())
		appErr := &AppErrs{
			code: CodeInternalServer,
			Msg:  "server encountered a problem",
		}
		appErr.AppErrWriter(w)
	case codes.PermissionDenied:
		appErr := &AppErrs{
			code: CodeBadRequest,
			Msg: st.Message(),
		}
		appErr.AppErrWriter(w)
	case codes.NotFound:
		appErr := &AppErrs{
			code: CodeNotFound,
			Msg: st.Message(),
		}
		appErr.AppErrWriter(w)
	}
}


func ErrMaker(err error) error {
	Err := &status.Status{}
	switch {
	case errors.Is(err, ErrEmailAlreadyExits):
		Err = status.Newf(codes.AlreadyExists, "user with the email already exists")
	case errors.Is(err, ErrNameAlreadyExits):
		Err = status.Newf(codes.AlreadyExists, "user with the name already exists")
	case errors.Is(err, PlanServerNotResponding):
		Err = status.New(codes.Canceled, err.Error())
	case errors.Is(err, ErrIncorrectPassword):
		Err = status.New(codes.PermissionDenied, "the password is incorrect")
	case errors.Is(err, ErrOldEmailNewEmailSame):
		Err = status.New(codes.InvalidArgument, err.Error())
	case errors.Is(err, ErrEmailDoesntMatch):
		Err = status.New(codes.InvalidArgument, err.Error())
	case errors.Is(err, ErrEmailNotFound):
		Err = status.New(codes.NotFound, err.Error())
	default:
		Err = status.Newf(codes.Internal, "some internal error occured : %v", err)
	}
	err = Err.Err()

	return err
}




