package customerrors

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrWorkoutNotEnded       = errors.New("user already has a workout that is not ended")
	ErrExerciseDoesnotExits  = errors.New("exercise doesnot exist in the DB. please create it")
	ErrGettingClaims         = errors.New("error getting claims from the context")
	ErrPlanNameDoesNotExists = errors.New("this plan name does not exist")
	ErrPlanAlreadyExists = errors.New("plan already exists")
)


type MyErrors struct {
	ErorrCode int 
	Message string `json:"message"`
}


var (
	ErrPlanAlreadyExists2 = MyErrors{ErorrCode: http.StatusBadRequest, Message: "this plan already exits"}
	ErrInternalServerError = MyErrors{ErorrCode: http.StatusInternalServerError, Message: "something went wrong"}
)

func CustomErrorWriter(w http.ResponseWriter, err *MyErrors) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(err.ErorrCode)
	json.NewEncoder(w).Encode(err)
}

func NewMyError(code int, message string) *MyErrors {
	return &MyErrors{ErorrCode: code, Message: message}
}