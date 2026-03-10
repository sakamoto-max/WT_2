package validators

import (
	"encoding/json"
	"errors"
	"net/http"
	"plan_service/internal/models"
)

type validationKey string

var (
	CreatePlan validationKey = "create_plan"
)

var (
	ErrValidationErrOccured = errors.New("validation err occured")
)

type validationErr struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}

var (
	ErrPlanNameReq  = validationErr{Path: "plan_name", Message: "required"}
	ErrExercisesReq = validationErr{Path: "exercises", Message: "required", Type: "slice"}
)

func ValidationErrWriter(w http.ResponseWriter, errs *[]validationErr) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errs)
}

func CreatePlanValidator(userInput *models.Plan2) (*[]validationErr, error) {

	var validationErrs []validationErr

	if userInput.PlanName == "" {
		validationErrs = append(validationErrs, ErrPlanNameReq)
	}

	if len(userInput.Exercises) == 0 {
		validationErrs = append(validationErrs, ErrExercisesReq)
	}

	if len(validationErrs) > 0 {
		return &validationErrs, ErrValidationErrOccured
	}

	return &validationErrs, nil
}
