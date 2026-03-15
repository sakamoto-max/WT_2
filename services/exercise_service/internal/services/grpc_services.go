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

func (e *ExerciseService) ExerciseExistsReturnId(ctx context.Context, exerciseName string) (bool, int32, error) {
	return e.DB.ExerciseExistsReturnId(ctx, exerciseName)
}


func (e *ExerciseService) GetExerciseName(ctx context.Context, exerciseId int) (string, error) {
	return e.DB.GetExerciseNameByID(ctx, exerciseId)
}