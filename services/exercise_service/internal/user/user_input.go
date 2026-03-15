package user

import "time"

type ValidatonErrs struct {
	Path    string
	Message string
}

var (
	ErrExerciseNameReq = ValidatonErrs{Path: "exericse_name", Message: "required"}
	ErrBodyPartReq = ValidatonErrs{Path: "body_part", Message: "required"}
	ErrEquipmentReq = ValidatonErrs{Path: "equipment", Message: "required"}
)

type ExerciseName struct {
	Name string `json:"exercise_name"`
}

func (e *ExerciseName) Validate() (*[]ValidatonErrs, bool) {

	var validationErrs []ValidatonErrs

	if e.Name == "" {
		validationErrs = append(validationErrs, ErrExerciseNameReq)
	}

	if len(validationErrs) > 0 {
		return &validationErrs, true
	}

	return &validationErrs, false
}

type Exercise struct {
	Id        int       `json:"id,omitempty"`
	Name      string    `json:"name"`
	RestTime  int       `json:"rest_time_in_seconds"`
	BodyPart  string    `json:"body_part"`
	Equipment string    `json:"equipment"`
	CreatedAt time.Time `json:"created_at"`
}

func (e *Exercise) Validate() (*[]ValidatonErrs, bool) {

	var validationErrs []ValidatonErrs
	
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
		return &validationErrs, true
	}

	return &validationErrs, false
}
