package mock

import (
	"context"
	"time"
	"tracker_service/internal/models"
)

var (
	timeDuration time.Duration = time.Second
	trackerId string = "170ff3fb-cc2a-4001-8ae2-ff02acf4dad5"
	exerciseName string = "mock_exercise"
	exerciseId string = "29433d8-a9fb-4643-82aa-2258b1387e49"
	currentPlanName string = "mock_plan_name"
)

type mockDb struct{}

func (m *mockDb) GetConflictLevel(ctx context.Context, userId string) (int, error)                  {
	return 2, nil
}

func (m *mockDb) GetUserWorkingOutWithPlan(ctx context.Context, userId string) (bool, error)        {
	return true, nil
}

func (m *mockDb) GetUserTrackerData(ctx context.Context, userId string) (*models.Tracker, error)    {
	workout := models.Workout{
		ExerciseId: exerciseId,
		RepsWeight: []models.Reps{
			{10, 10},
			{10, 9},
			{10, 8},
		},
	}
	resp := models.Tracker{
		Workout: []models.Workout{
			workout,
			workout,
			workout,
		},
	}

	return &resp,nil
}

func (m *mockDb) GetPostgresRespTime(ctx context.Context) *time.Duration                         {
	return &timeDuration
}
func (m *mockDb) GetRedisRespTime(ctx context.Context) *time.Duration                            {
	return &timeDuration
}
func (m *mockDb) StartWorkout(ctx context.Context, userId string, planId string) (string, error) {
	return trackerId, nil
}
func (m *mockDb) DeleteTrackerIdInPG(ctx context.Context, trackerId string) error                {
	return nil
}
func (m *mockDb) RevertStartWorkout(ctx context.Context, trackerId string) error                 {
	return nil
}
func (m *mockDb) SetTrackerId(ctx context.Context, userId string, trackerId string) error        {
	return nil
}
func (m *mockDb) GetTrackerId(ctx context.Context, userId string) (string, error)                {
	return trackerId, nil
}
func (m *mockDb) DelTrackerId(ctx context.Context, userId string) error                          {
	return nil
}
func (m *mockDb) EndWorkout(ctx context.Context, trackerId string, data *models.Tracker) error   {
	return nil
}
func (m *mockDb) EndWorkoutWithOutbox(ctx context.Context, userId string, trackerId string, data *models.Tracker, planName string, newExerciseNames *[]string) error {
	return nil
}
func (m *mockDb) SetExerciseNameById(ctx context.Context, exerciseId string, exerciseName string) error {
	return nil
}
func (m *mockDb) GetExerciseNameById(ctx context.Context, exerciseId string) (string, error)       {
	return exerciseName, nil
}
func (m *mockDb) SetUserCurrentPlanName(ctx context.Context, userId string, planName string) error {
	return nil
}
func (m *mockDb) GetUserCurrentPlanName(ctx context.Context, userId string) (string, error)        {
	return currentPlanName, nil
}
func (m *mockDb) SetPlanWithExercises(ctx context.Context, userId string, planName string, exerciseNames *[]string) error {
	return nil
}
func (m *mockDb) GetPlanWithExercises(ctx context.Context, userId string, planName string) (*[]string, error) {
	return &[]string{exerciseId, exerciseId, exerciseId}, nil
}
func (m *mockDb) SetUserWorkingOutWithPlan(ctx context.Context, userId string, value bool) error    {
	return nil
}

func (m *mockDb) SetConflictLevel(ctx context.Context, userId string, conflictLevel int) error      {
	return nil
}

func (m *mockDb) SetUserTrackerData(ctx context.Context, userId string, data *models.Tracker) error {
	return nil
}

func (m *mockDb) SetUserNewExercises(ctx context.Context, userId string, exerciseNames *[]string) error {
	return nil
}
func (m *mockDb) GetUserNewExercises(ctx context.Context, userId string) (*[]string, error) {
	return &[]string{exerciseId, exerciseId, exerciseId}, nil
}
func (m *mockDb) DelAllUserData(ctx context.Context, userId string, planName string) error  {
	return nil
}
