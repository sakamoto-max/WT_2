package models

type PlanName struct {
	Name string `json:"plan_name"`
}

type Reps struct {
	Weight float32 `json:"weight"`
	Reps int `json:"reps"`
}

type Workout struct {
	ExerciseId string `json:"exercise_id"`
	ExerciseName string `json:"exercise_name"`
	RepsWeight []Reps `json:"tracker"`
}

type Tracker struct {
	// PlanId int `json:"plan_id"`
	UserResponse bool `json:"user_response"`
	Workout []Workout `json:"workout"`
}

func (t *Tracker) GetAllExercises() *[]string {
	var allExercises []string

	for _, eachExer := range t.Workout{
		allExercises = append(allExercises, eachExer.ExerciseName)
	}

	return &allExercises
}

type GeneralResp struct {
	Message string `json:"message"`
}

type Plan struct{
	Message string `json:"message"`
	PlanName string `json:"plan_name"`
	Exercises []string `json:"exercises_in_plan"`
}

type ExercisesNotPerformed struct {
	Exercises []string `json:"exercises_not_performed"`
}

type UpdatePlanPayLoad struct {
	UserId string `json:"user_id"`
	PlanName string `json:"plan_name"`
	ExerciseNames *[]string `json:"exercise_names"`
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