package repository

import (
	"context"
	"time"
	"tracker_service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	Metrics interface {
		GetRespTime(ctx context.Context) *time.Duration
	}
	End interface {
		EndWorkout(ctx context.Context, trackerId string, data *domain.Tracker) error
		EndWorkoutWithOutbox(ctx context.Context, userId string, trackerId string, data *domain.Tracker, planName string, newExerciseNames *[]string) error
	}
	Cancel interface {
		DeleteTrackerIdInPG(ctx context.Context, trackerId string) error
	}
	Start interface {
		StartWorkout(ctx context.Context, payload domain.StartWorkout) (string, error)
		RevertStartWorkout(ctx context.Context, trackerId string) error
	}
}

func NewDb(pool *pgxpool.Pool) *Db {
	return &Db{
		Metrics: &metricsRepo{pool},
		End:     &endWorkoutRepo{pool},
		Cancel:  &cancelRepo{pool},
		Start:   &startWorkoutRepo{pool},
	}
}
