package cachemock

import (
	"context"
	"fmt"
	"time"
	"tracker_service/internal/domain"
)

type MetricsMock struct {
	Down bool
}

func (m *MetricsMock) GetRespTime(ctx context.Context) *time.Duration {
	if m.Down {
		return nil
	}

	timeDuration := time.Millisecond * 234

	return &timeDuration
}

type UserDataMock struct {
	Down bool
}

func (u *UserDataMock) DelAllUserData(ctx context.Context, userId string, planName string) error {
	if u.Down {
		return fmt.Errorf("redis is down")
	}

	return nil
}

type TrackerIdMock struct {
	Down           bool
	WorkoutOngoing bool
}

func (t *TrackerIdMock) SetTrackerId(ctx context.Context, userId string, trackerId string) error {
	if t.Down {
		return fmt.Errorf("redis is dowm")
	}
	return nil
}

func (t *TrackerIdMock) GetTrackerId(ctx context.Context, userId string) (string, error) {
	if t.Down {
		return "", fmt.Errorf("redis is dowm")
	}

	if t.WorkoutOngoing {
		return "123", nil
	}

	return "", nil

}

func (t *TrackerIdMock) DelTrackerId(ctx context.Context, userId string) error {
	if t.Down {
		return fmt.Errorf("redis is dowm")
	}

	if !t.WorkoutOngoing {
		return fmt.Errorf("no workout is ongoing")
	}

	return nil
}

type ExerciseNameMock struct {
	Down bool
}

func (e *ExerciseNameMock) SetExerciseNameById(ctx context.Context, exerciseId string, exerciseName string) error {
	if e.Down {
		return fmt.Errorf("cache is down")
	}

	return nil
}
func (e *ExerciseNameMock) GetExerciseNameById(ctx context.Context, exerciseId string) (string, error) {
	if e.Down {
		return "", fmt.Errorf("cache is down")
	}

	return "123", nil
}

type CurrentPlanMock struct {
	Down           bool
	WorkoutOngoing bool
	WithPlan       bool
}

func (c *CurrentPlanMock) SetUserCurrentPlanName(ctx context.Context, userId string, planName string) error {
	if c.Down {
		return fmt.Errorf("cache is down")
	}

	return nil
}

func (c *CurrentPlanMock) GetUserCurrentPlanName(ctx context.Context, userId string) (string, error) {
	if c.Down {
		return "", fmt.Errorf("cache is down")
	}

	if !c.WorkoutOngoing {
		return "", fmt.Errorf("no workout is ongoing")
	}

	if c.WithPlan {
		return "plan", nil
	}

	return "", nil
}

type PlanMock struct {
	Down           bool
	WorkoutOngoing bool
}

func (p *PlanMock) SetPlanWithExercises(ctx context.Context, userId string, planName string, exerciseNames *[]string) error {
	if p.Down {
		return fmt.Errorf("redis is down")
	}

	return nil
}

func (p *PlanMock) GetPlanWithExercises(ctx context.Context, userId string, planName string) (*[]string, error) {
	if p.Down {
		return nil, fmt.Errorf("redis is down")
	}

	if !p.WorkoutOngoing {
		return nil, fmt.Errorf("no workout is ongoing")
	}

	return &[]string{"exer_1", "exer_2", "exer_3"}, nil
}

type WithPlanMock struct {
	Down     bool
	WithPlan bool
}

func (w *WithPlanMock) SetUserWorkingOutWithPlan(ctx context.Context, userId string, value bool) error {
	if w.Down {
		return fmt.Errorf("cache is down")
	}

	return nil

}
func (w *WithPlanMock) GetUserWorkingOutWithPlan(ctx context.Context, userId string) (bool, error) {
	if w.Down {
		return false, fmt.Errorf("redis is down")
	}

	if w.WithPlan {
		return true, nil
	}

	return false, nil
}

type ConflictMock struct {
	Down          bool
	ConflictLevel int
}

func (c *ConflictMock) SetConflictLevel(ctx context.Context, userId string, conflictLevel int) error {
	if c.Down {
		return fmt.Errorf("cache is down")
	}

	return nil
}
func (c *ConflictMock) GetConflictLevel(ctx context.Context, userId string) (int, error) {
	if c.Down {
		return 0, fmt.Errorf("cache is down")
	}

	return c.ConflictLevel, nil
}

type TrackerDataMock struct {
	Down bool
	Data *domain.Tracker
}

func (t *TrackerDataMock) SetUserTrackerData(ctx context.Context, userId string, data *domain.Tracker) error {
	if t.Down {
		return fmt.Errorf("cache is down")
	}

	return nil
}
func (t *TrackerDataMock) GetUserTrackerData(ctx context.Context, userId string) (*domain.Tracker, error) {
	if t.Down {
		return nil, fmt.Errorf("cache is down")
	}
	return t.Data, nil
}

type NewExercisesMock struct{
	Down bool
}

func (n *NewExercisesMock) SetUserNewExercises(ctx context.Context, userId string, exerciseNames *[]string) error {
	if n.Down {
		return fmt.Errorf("cache is down")
	}

	return nil
}
func (n *NewExercisesMock) GetUserNewExercises(ctx context.Context, userId string) (*[]string, error) {
	if n.Down {
		return nil, fmt.Errorf("cache is down")
	}

	return &[]string{"exer_1", "exer_2", "exer_3"}, nil
}
