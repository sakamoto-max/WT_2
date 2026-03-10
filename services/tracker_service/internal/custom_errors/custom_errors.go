package customerrors

import (
	"encoding/json"
	"net/http"
)

type MyErrors struct {
	Message string
}

// var (
// 	PlanDoesNotExist = "PLAN_DOES_NOT_EXIST"
// )

var (
	PlanDoesNotExistErr = MyErrors{Message: "plan does not exist, please create it"}
)

func CustomErrWriter(w http.ResponseWriter, err MyErrors) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(err)
}