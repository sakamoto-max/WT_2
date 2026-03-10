package models

import "time"

type Plan struct {
	PlanName string `json:"plan_name"`
	UserId   int    `json:"user_id"`
}

type Plan2 struct {
	PlanName  string   `json:"plan_name"`
	Exercises []string `json:"exercises"`
}

type Plan2Resp struct{
	Message string `json:"message"`
	PlanName  string   `json:"plan_name"`
	Exercises []string `json:"exercises"`
}

type Plan3 struct {
	Id       int
	PlanName string
}

type AllPlansResp struct {
	NumberOfPlans int     `json:"number_of_plan"`
	Plans         []Plan2 `json:"plans"`
}

type PlanResp struct {
	Message string `json:"message"`
}

type StartPlanWorkoutResp struct {
	Message   string   `json:"message"`
	Exercises []string `json:"exercises"`
}

type ExerciseArry struct {
	Exercises []string `json:"exercises"`
}

type Exercise struct {
	Id   int    `json:"id"`
	Name string `json:"exercise_name"`
}

type AddExerciseResp struct {
	PlanName  string   `json:"plan_name"`
	Exercises []string `json:"exercises"`
}

type AddExercise struct {
	ExerciseName string `json:"exercise_name"`
}

type RepsWeights struct {
	Set    int `json:"set"`
	Reps   int `json:"reps"`
	Weight int `json:"weight"`
}

type AddRepsWeightResp struct {
	ExerciseId int           `json:"exercise_id"`
	Sets       []RepsWeights `json:"sets"`
}

type MqMsg struct {
	UserId int
	Action string
	Time time.Time
	Email string
}
