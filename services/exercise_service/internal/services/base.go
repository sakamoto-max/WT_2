package services

import (
	"context"
	"exercise_service/internal/domain"
	"exercise_service/internal/repository"
	"time"
)

type service struct {
	DB repository.RepoIface
}

func NewService(r repository.RepoIface) *service {
	return &service{
		DB: r,
	}
}

type ServiceIface interface {
	GetExerciseByName(ctx context.Context, userID string, exerciseName string) (*domain.Exercise, error)
	GetAllExercisesSer(ctx context.Context, userId string) (*[]domain.Exercise, error)
	DeleteExeciseSer(ctx context.Context, userId string, exerciseName string) error
	CreateExerciseSer(ctx context.Context, userId string, exerciseName string, bodyPartName string, equipmentName string) (string, error)
	ExerciseExistsReturnId(ctx context.Context, userId string, exerciseName string) (string, error)
	GetExerciseNameByID(ctx context.Context, exerciseId string) (string, error)
	GetHealth(ctx context.Context) (*time.Duration, *time.Duration)
}
