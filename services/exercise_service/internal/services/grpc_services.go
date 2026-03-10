package services

import (
	"context"
	"exercise_service/internal/repository"

)

type ExerciseService struct {
	DB *repository.Repo
}

func NewExerciseService(repo *repository.Repo) *ExerciseService {
	return &ExerciseService{
		DB: repo,
	}
}

func (o *ExerciseService) ExerciseExistsReturnId(ctx context.Context, exerciseName string) (bool, int32, error) {
	return o.DB.ExerciseExistsReturnId(ctx, exerciseName)
}
