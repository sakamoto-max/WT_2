package repository

import (
	"context"
	"exercise_service/internal/domain"
	"exercise_service/internal/mappings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	Metrics interface {
		GetRespTime(ctx context.Context) *time.Duration
	}
	ExerciseCD interface {
		CreateExercise(ctx context.Context, payload mappings.CreateExercise) (string, error)
		DeleteExecise(ctx context.Context, payload mappings.DeleteExercise) error
	}
	ExerciseGet interface {
		GetExerciseByName(ctx context.Context, payload mappings.GetExerciseByName) (*domain.Exercise, error)
		GetAllExercises(ctx context.Context, paylaod mappings.GetAllExercises) (*[]domain.Exercise, error)
		GetExerciseNameByID(ctx context.Context, exerciseId string) (string, error)
		ExerciseExistsReturnId(ctx context.Context, payload mappings.ExerciseExistsReturnId) (string, error)
	}
}

func NewDb(pg *pgxpool.Pool) *Db {
	return &Db{
		Metrics:     &metricsDb{pg},
		ExerciseCD:  &exerciseCD{pg},
		ExerciseGet: &exerciseGetDB{pg},
	}
}
