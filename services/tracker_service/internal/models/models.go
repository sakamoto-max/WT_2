package models

type PlanName struct {
	Name string `json:"plan_name"`
}

type Reps struct {
	Weight int `json:"weight"`
	Reps int `json:"reps"`
}

type Workout struct {
	ExerciseId int `json:"exercise_id"`
	Tracker []Reps `json:"tracker"`
}

type Tracker struct {
	PlanId int `json:"plan_id"`
	Workout []Workout `json:"workout"`
}

type GeneralResp struct {
	Message string `json:"message"`
}