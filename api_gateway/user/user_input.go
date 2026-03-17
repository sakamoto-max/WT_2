package user

import (
	"errors"
	"time"
)

type validatonErrs struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
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

var (
	ErrExerciseNameReq = validatonErrs{Path: "exericse_name", Message: "required"}
	ErrBodyPartReq = validatonErrs{Path: "body_part", Message: "required"}
	ErrEquipmentReq = validatonErrs{Path: "equipment", Message: "required"}
)

var (
	// ErrPlanNameReq  = validatonErrs{Path: "plan_name", Message: "required"}
	// ErrExercisesReq = validatonErrs{Path: "exercises", Message: "required", Type: "slice"}
	ErrWorkoutReq = validatonErrs{Path: "workout", Message: "required"}
	ErrExerciseIdReq = validatonErrs{Path: "exercise_id", Message: "required"}
	ErrTrackerReq = validatonErrs{Path: "tracker", Message: "required", Type: "slice"}
	ErrRepsReq = validatonErrs{Path: "reps", Message: "required"}
	ErrWeightReq = validatonErrs{Path: "reps", Message: "required"}
)	


type Reps struct {
	Weight int `json:"weight"`
	Reps int `json:"reps"`
}


type Workout struct {
	ExerciseId int `json:"exercise_id"`
	RepsWeight []Reps `json:"tracker"`
}


type Tracker struct {
	Workout []Workout `json:"workout"`
}

func (t *Tracker) Validate() (*[]validatonErrs, bool) {

	var validationErrs []validatonErrs

	if len(t.Workout) == 0 {
		validationErrs = append(validationErrs, ErrWorkoutReq)
	}

	for _, eachExer := range t.Workout {
		var emptyExerId int
		if eachExer.ExerciseId == emptyExerId {
			validationErrs = append(validationErrs, ErrExerciseIdReq)
		}

		if len(eachExer.RepsWeight) == 0 {
			validationErrs = append(validationErrs, ErrTrackerReq)
		}

		for _, repsAndWeigh := range eachExer.RepsWeight{
			if repsAndWeigh.Reps == 0{
				validationErrs = append(validationErrs, ErrRepsReq)
			}

			if repsAndWeigh.Weight == 0{
				validationErrs = append(validationErrs, ErrWeightReq)
			}
		}
	}

	if len(validationErrs) > 0{
		return &validationErrs, true
	}

	return &validationErrs, false

}

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

type ExerciseName struct {
	Name string `json:"exercise_name"`
}

func (e *ExerciseName) Validate() (*[]validatonErrs, bool) {

	var validationErrs []validatonErrs

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

func (e *Exercise) Validate() (*[]validatonErrs, bool) {

	var validationErrs []validatonErrs

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

var (
	ErrPlanNameReq  = validatonErrs{Path: "plan_name", Message: "required"}
	ErrExercisesReq = validatonErrs{Path: "exercises", Message: "required", Type: "slice"}
)

type Plan struct {
	PlanName  string   `json:"plan_name"`
	Exercises []string `json:"exercises"`
}

func (p *Plan) Validate() (*[]validatonErrs, bool) {
	var validationErrs []validatonErrs

	if p.PlanName == "" {
		validationErrs = append(validationErrs, ErrPlanNameReq)
	}

	if len(p.Exercises) == 0 {
		validationErrs = append(validationErrs, ErrExercisesReq)
	}

	if len(validationErrs) > 0 {
		return &validationErrs, true
	}

	return &validationErrs, false
}

type PlanName struct {
	PlanName string `json:"plan_name"`
}

func (p *PlanName) Validate() (*[]validatonErrs, bool) {
	var validationErrs []validatonErrs

	if p.PlanName == "" {
		validationErrs = append(validationErrs, ErrPlanNameReq)
	}

	if len(validationErrs) > 0 {
		return &validationErrs, true
	}

	return &validationErrs, false
}

