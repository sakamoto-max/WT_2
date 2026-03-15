package user


type ValidationErr struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}

var (
	ErrPlanNameReq  = ValidationErr{Path: "plan_name", Message: "required"}
	ErrExercisesReq = ValidationErr{Path: "exercises", Message: "required", Type: "slice"}
	ErrWorkoutReq = ValidationErr{Path: "workout", Message: "required"}
	ErrExerciseIdReq = ValidationErr{Path: "exercise_id", Message: "required"}
	ErrTrackerReq = ValidationErr{Path: "tracker", Message: "required", Type: "slice"}
	ErrRepsReq = ValidationErr{Path: "reps", Message: "required"}
	ErrWeightReq = ValidationErr{Path: "reps", Message: "required"}
)	

type PlanName struct {
	PlanName string `json:"plan_name"`
}

func (p *PlanName) Validate() (*[]ValidationErr, bool) {
	var validationErrs []ValidationErr

	if p.PlanName == "" {
		validationErrs = append(validationErrs, ErrPlanNameReq)
	}

	if len(validationErrs) > 0 {
		return &validationErrs, true
	}

	return &validationErrs, false
}

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

func (t *Tracker) Validate() (*[]ValidationErr, bool) {

	var validationErrs []ValidationErr

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