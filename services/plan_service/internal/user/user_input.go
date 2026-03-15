package user

type ValidationErr struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}

var (
	ErrPlanNameReq  = ValidationErr{Path: "plan_name", Message: "required"}
	ErrExercisesReq = ValidationErr{Path: "exercises", Message: "required", Type: "slice"}
)

type Plan2 struct {
	PlanName  string   `json:"plan_name"`
	Exercises []string `json:"exercises"`
}

func (p *Plan2) Validate() (*[]ValidationErr, bool) {
	var validationErrs []ValidationErr

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
