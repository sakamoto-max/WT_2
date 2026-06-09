package models

import (
	"encoding/json"
	"errors"
	// "fmt"
	"net/http"
)

type ValidationErrs struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}

type validationErrs []ValidationErrs

func NewValidationErr(Path string, Message string, w http.ResponseWriter) {
	err := ValidationErrs{
		Path:    Path,
		Message: Message,
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(err)
}

func (v validationErrs) Error() string {
	return "validation error occured"
}

var (
	ErrValidationErr = errors.New("validation error occured")
)

type Validatable interface {
	Validate() error
}

var (
	errUserNameReq               = ValidationErrs{Path: "name", Message: "required"}
	errUserEmailReq              = ValidationErrs{Path: "email", Message: "required"}
	errUserPassWordReq           = ValidationErrs{Path: "password", Message: "required"}
	errUUIDReq                   = ValidationErrs{Path: "UUID", Message: "required"}
	ErrExerciseNameReq           = ValidationErrs{Path: "name", Message: "required"}
	ErrBodyPartReq               = ValidationErrs{Path: "body_part", Message: "required"}
	ErrEquipmentReq              = ValidationErrs{Path: "equipment", Message: "required"}
	ErrWorkoutReq                = ValidationErrs{Path: "data", Message: "required"}
	ErrExerciseNameForTrackerReq = ValidationErrs{Path: "exercise_name", Message: "required"}
	ErrTrackerReq                = ValidationErrs{Path: "tracker", Message: "required", Type: "slice"}
	ErrRepsReq                   = ValidationErrs{Path: "reps", Message: "required"}
	ErrWeightReq                 = ValidationErrs{Path: "reps", Message: "required"}
	ErrPlanNameReq               = ValidationErrs{Path: "plan_name", Message: "required"}
	ErrExercisesReq              = ValidationErrs{Path: "exercises", Message: "required", Type: "slice"}
	ErrOldPassReq                = ValidationErrs{Path: "old_password", Message: "required"}
	ErrNewPassReq                = ValidationErrs{Path: "new_password", Message: "required"}
	ErrOldEmailReq               = ValidationErrs{Path: "old_email", Message: "required"}
	ErrNewEmailReq               = ValidationErrs{Path: "new_email", Message: "required"}
	ErrConflictResponse          = ValidationErrs{Path: "user_response", Message: `can only be "yes" or "no"`}
)

func (c *ChangeEmail) Validate() error {
	var validationErrs validationErrs

	if c.OldEmail == "" {
		validationErrs = append(validationErrs, ErrOldEmailReq)
	}

	if c.NewEmail == "" {
		validationErrs = append(validationErrs, ErrNewEmailReq)
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}

	return nil
}

func (t *Tracker) Validate() error {

	// fmt.Println("len of workout", len(t.Workout))

	// return nil

	// check if user response is there -> not there -> check tracker data -> atleast one exercise data should be there
	// there -> skip the tracker validation

	var validationErrs validationErrs

	if t.UserResponse == "" {

		if len(t.Workout) == 0 {
			validationErrs = append(validationErrs, ErrWorkoutReq)
		}

		for _, eachExer := range t.Workout {

			if eachExer.ExerciseName == "" {
				validationErrs = append(validationErrs, ErrExerciseNameForTrackerReq)
			}

			if len(eachExer.RepsWeight) == 0 {
				validationErrs = append(validationErrs, ErrTrackerReq)
			}

			for _, repsAndWeigh := range eachExer.RepsWeight {
				if repsAndWeigh.Reps == 0 {
					validationErrs = append(validationErrs, ErrRepsReq)
				}

				if repsAndWeigh.Weight == 0 {
					validationErrs = append(validationErrs, ErrWeightReq)
				}
			}
		}
	} else {

		if t.UserResponse != "yes" {
			if t.UserResponse != "no" {
				validationErrs = append(validationErrs, ErrConflictResponse)
			}
		}
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}
	return nil

}

func (s *Signup) Validate() error {

	var validationErrs validationErrs

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
		return validationErrs
	}

	return nil
}

func (l *Login) Validate() error {

	var validationErrs validationErrs
	if l.Email == "" {
		validationErrs = append(validationErrs, errUserEmailReq)
	}

	if l.Password == "" {
		validationErrs = append(validationErrs, errUserPassWordReq)
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}

	return nil
}

func (u *UUIDReader) Validate() error {

	var validationErrs validationErrs

	if u.UUID == "" {
		validationErrs = append(validationErrs, errUUIDReq)
		return validationErrs
	}

	return nil
}

func (e *ExerciseName) Validate() error {

	var validationErrs validationErrs

	if e.Name == "" {
		validationErrs = append(validationErrs, ErrExerciseNameReq)
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}

	return nil
}

func (e *Exercise) Validate() error {

	var validationErrs validationErrs

	if e.Name == "" {
		validationErrs = append(validationErrs, ErrExerciseNameReq)
	}

	if e.BodyPart == "" {
		validationErrs = append(validationErrs, ErrBodyPartReq)
	}

	if e.Equipment == "" {
		validationErrs = append(validationErrs, ErrEquipmentReq)
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}

	return nil
}

func (p *Plan) Validate() error {
	var validationErrs validationErrs

	if p.PlanName == "" {
		validationErrs = append(validationErrs, ErrPlanNameReq)
	}

	if len(p.Exercises) == 0 {
		validationErrs = append(validationErrs, ErrExercisesReq)
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}

	return nil
}

func (p *PlanName) Validate() error {
	var validationErrs validationErrs

	if p.PlanName == "" {
		validationErrs = append(validationErrs, ErrPlanNameReq)
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}

	return nil
}

func (c *ChangePass) Validate() error {

	var validationErrs validationErrs

	if c.OldPass == "" {
		validationErrs = append(validationErrs, ErrOldPassReq)
	}

	if c.NewPass == "" {
		validationErrs = append(validationErrs, ErrNewPassReq)
	}

	if len(validationErrs) > 0 {
		return validationErrs
	}

	return nil

}

func ValidationErrWriter(w http.ResponseWriter, err error) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(err)
}
