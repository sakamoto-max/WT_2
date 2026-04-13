package models

type PlanName struct {
	Name string `json:"plan_name"`
}

type Reps struct {
	Weight int `json:"weight"`
	Reps int `json:"reps"`
}

type Workout struct {
	ExerciseId string `json:"exercise_id"`
	ExerciseName string `json:"exercise_name"`
	RepsWeight []Reps `json:"tracker"`
}

type Tracker struct {
	// PlanId int `json:"plan_id"`
	Workout []Workout `json:"workout"`
}

type GeneralResp struct {
	Message string `json:"message"`
}

type Plan struct{
	Message string `json:"message"`
	PlanName string `json:"plan_name"`
	Exercises []string `json:"exercises_in_plan"`
}

// {
// 	"workout" : [
// 		{
// 			"exercise_name" : "push_ups",
// 			"tracker" : [
// 				{
// 					"weight" : 10,
// 					"reps" : 11
// 				},
// 				{
// 					"weight" : 10,
// 					"reps" : 11
// 				},
// 				{
// 					"weight" : 10,
// 					"reps" : 11
// 				}
// 			]
// 		},
// 		{
// 			"exercise_name" : "pull_ups",
// 			"tracker" : [
// 				{
// 					"weight" : 10,
// 					"reps" : 11
// 				},
// 				{
// 					"weight" : 10,
// 					"reps" : 11
// 				},
// 				{
// 					"weight" : 10,
// 					"reps" : 11
// 				}
// 			]
// 		},
// 		{
// 			"exercise_name" : "pull_ups",
// 			"tracker" : [
// 				{
// 					"weight" : 10,
// 					"reps" : 11
// 				},
// 				{
// 					"weight" : 10,
// 					"reps" : 11
// 				},
// 				{
// 					"weight" : 10,
// 					"reps" : 11
// 				}
// 			]
// 		}
// 	]
// }