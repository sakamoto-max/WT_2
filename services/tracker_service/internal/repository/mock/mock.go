package mock

import (
	"context"
	"fmt"
	"time"
	"tracker_service/internal/domain"
)

type MetricsMock struct {
	PgDown bool
}

func (m *MetricsMock) GetRespTime(ctx context.Context) *time.Duration {
	if m.PgDown {
		return nil
	}

	timeDuration := time.Millisecond * 2342

	return &timeDuration
}

type EndMock struct {
	PgDown         bool
	WorkoutOngoing bool
}

func (e *EndMock) EndWorkout(ctx context.Context, trackerId string, data *domain.Tracker) error {
	if e.PgDown {
		return fmt.Errorf("pg is down")
	}

	if !e.WorkoutOngoing {
		return fmt.Errorf("no workout is going on")
	}

	return nil
}
func (e *EndMock) EndWorkoutWithOutbox(ctx context.Context, userId string, trackerId string, data *domain.Tracker, planName string, newExerciseNames *[]string) error {
	if e.PgDown {
		return fmt.Errorf("pg is down")
	}

	if !e.WorkoutOngoing {
		return fmt.Errorf("no workout is going on")
	}

	return nil
}

type CancelMock struct{
	PgDown bool
	WorkoutOngoing bool
}

func (c *CancelMock) DeleteTrackerIdInPG(ctx context.Context, trackerId string) error {
	if c.PgDown {
		return fmt.Errorf("pg is down")
	}

	if !c.WorkoutOngoing {
		return fmt.Errorf("no workout is ongoing")
	}

	return nil
}

type StartMock struct{
	PgDown bool
	WorkoutOngoing bool
}

func (s *StartMock) StartWorkout(ctx context.Context, paylaod domain.StartWorkout) (string, error) {
	if s.PgDown {
		return "", fmt.Errorf("pg is down")
	}

	if s.WorkoutOngoing {
		return "", fmt.Errorf("already a workout is going on")
	}

	return "123", nil

}
func (s *StartMock) RevertStartWorkout(ctx context.Context, trackerId string) error {
	return nil
}
