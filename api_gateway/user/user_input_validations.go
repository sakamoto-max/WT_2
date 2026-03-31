package user

import "errors"

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
	ErrExerciseNameReq = validatonErrs{Path: "name", Message: "required"}
	ErrBodyPartReq     = validatonErrs{Path: "body_part", Message: "required"}
	ErrEquipmentReq    = validatonErrs{Path: "equipment", Message: "required"}
	ErrWorkoutReq      = validatonErrs{Path: "workout", Message: "required"}
	ErrExerciseIdReq   = validatonErrs{Path: "exercise_id", Message: "required"}
	ErrTrackerReq      = validatonErrs{Path: "tracker", Message: "required", Type: "slice"}
	ErrRepsReq         = validatonErrs{Path: "reps", Message: "required"}
	ErrWeightReq       = validatonErrs{Path: "reps", Message: "required"}
	ErrPlanNameReq     = validatonErrs{Path: "plan_name", Message: "required"}
	ErrExercisesReq    = validatonErrs{Path: "exercises", Message: "required", Type: "slice"}
	ErrOldPassReq = validatonErrs{Path: "old_password", Message: "required"}
	ErrNewPassReq = validatonErrs{Path: "new_password", Message: "required"}
	ErrOldEmailReq = validatonErrs{Path: "old_email", Message: "required"}
	ErrNewEmailReq = validatonErrs{Path: "new_email", Message: "required"}
)

func (c *ChangeEmail) Validate() (*[]validatonErrs, bool) {	
	var validationErrs []validatonErrs

	if c.OldEmail == "" {
		validationErrs = append(validationErrs, ErrOldEmailReq)
	}

	if c.NewEmail == "" {
		validationErrs = append(validationErrs, ErrNewEmailReq)
	}

	if len(validationErrs) > 0 {
		return &validationErrs, true
	}

	return nil, false
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

		for _, repsAndWeigh := range eachExer.RepsWeight {
			if repsAndWeigh.Reps == 0 {
				validationErrs = append(validationErrs, ErrRepsReq)
			}

			if repsAndWeigh.Weight == 0 {
				validationErrs = append(validationErrs, ErrWeightReq)
			}
		}
	}

	if len(validationErrs) > 0 {
		return &validationErrs, true
	}

	return nil, false

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

	return nil, false
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

	return nil, false
}

func (u *UUIDReader) Validate() (*[]validatonErrs, bool) {

	var validationErrs []validatonErrs

	// uuid := u.UUID.String()

	if u.UUID == "" {
		validationErrs = append(validationErrs, errUUIDReq)
		return &validationErrs, true
	}

	return nil, false
}

func (e *ExerciseName) Validate() (*[]validatonErrs, bool) {

	var validationErrs []validatonErrs

	if e.Name == "" {
		validationErrs = append(validationErrs, ErrExerciseNameReq)
	}

	if len(validationErrs) > 0 {
		return &validationErrs, true
	}

	return nil, false
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

	return nil, false
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

	return nil, false
}

func (p *PlanName) Validate() (*[]validatonErrs, bool) {
	var validationErrs []validatonErrs

	if p.PlanName == "" {
		validationErrs = append(validationErrs, ErrPlanNameReq)
	}

	if len(validationErrs) > 0 {
		return &validationErrs, true
	}

	return nil, false
}

func (c *ChangePass) Validate() (*[]validatonErrs, bool) {

	var validationErrs []validatonErrs

	if c.OldPass == "" {
		validationErrs = append(validationErrs, ErrOldPassReq)
	}

	if c.NewPass == "" {
		validationErrs = append(validationErrs, ErrNewPassReq)
	}

	if len(validationErrs) > 0{
		return &validationErrs, true
	}

	return &validationErrs, false

}