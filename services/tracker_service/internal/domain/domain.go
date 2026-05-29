package domain

import trackerpb "github.com/sakamoto-max/wt_2_proto/shared/tracker"

type StartWorkout struct {
	UserId string
	PlanId string
}

type Tracker struct {
	// PlanId int `json:"plan_id"`
	UserResponse bool      `json:"user_response,omitempty"`
	Workout      []Workout `json:"workout"`
}

type Reps struct {
	Weight float32 `json:"weight"`
	Reps   int     `json:"reps"`
}

type Workout struct {
	ExerciseId   string `json:"exercise_id"`
	ExerciseName string `json:"exercise_name"`
	RepsWeight   []Reps `json:"tracker"`
}

func ToTracker(in *trackerpb.EndWorkoutReq) Tracker {
	main := Tracker{}

	if in.UserResponse || !in.UserResponse {
		main.UserResponse = in.UserResponse
		return main
	}

	for _, eachExer := range in.AllExerices {
		w := Workout{
			ExerciseName: eachExer.ExerciseName,
		}

		for _, repsPlusWeight := range eachExer.SetsAndReps {

			rPlusW := Reps{
				Reps:   int(repsPlusWeight.Reps),
				Weight: repsPlusWeight.Weight,
			}

			w.RepsWeight = append(w.RepsWeight, rPlusW)
		}
		main.Workout = append(main.Workout, w)
	}

	return main
}

func (t *Tracker) GetAllExercises() *[]string {
	var allExercises []string

	for _, eachExer := range t.Workout{
		allExercises = append(allExercises, eachExer.ExerciseName)
	}

	return &allExercises
}



type UpdatePlanPayLoad struct {
	UserId string `json:"user_id"`
	PlanName string `json:"plan_name"`
	ExerciseNames *[]string `json:"exercise_names"`
}


func ConvertToLocal(in *trackerpb.EndWorkoutReq) *Tracker {

	main := Tracker{}

	if in.UserResponse {
		main.UserResponse = in.UserResponse
		return &main
	}

	for _, eachExer := range in.AllExerices {
		w := Workout{
			ExerciseName: eachExer.ExerciseName,
		}

		for _, repsPlusWeight := range eachExer.SetsAndReps {

			rPlusW := Reps{
				Reps:   int(repsPlusWeight.Reps),
				Weight: repsPlusWeight.Weight,
			}

			w.RepsWeight = append(w.RepsWeight, rPlusW)
		}
		main.Workout = append(main.Workout, w)
	}

	return &main
}

