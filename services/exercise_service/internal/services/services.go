package services

import (
	"context"
	"exercise_service/internal/models"
	"exercise_service/internal/repository"
	"exercise_service/internal/user"
	"fmt"
	"time"
)

type Service struct {
	DB *repository.Repo
}

func NewService(r *repository.Repo) *Service {
	return &Service{
		DB: r,
	}
}

func (s *Service) GetExerciseByNameSer(ctx context.Context, exerciseName string) (*models.Exercise, error) {
	// check if the execise exists
	// transform the exercises

	exercise, err := s.DB.GetExerciseByName(ctx, exerciseName)
	if err != nil {
		return exercise, err
	}
	return exercise, nil
}

func (s *Service) GetAllExercisesSer(ctx context.Context) (*[]models.Exercise, error) {
	allExercises, err := s.DB.GetAllExercises(ctx)
	if err != nil {
		return allExercises, err
	}

	return allExercises, nil
}

func (s *Service) DeleteExeciseSer(ctx context.Context, exerciseName string) error {
	// check if exercise exists
	err := s.DB.DeleteExecise(ctx, exerciseName)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateExerciseSer(ctx context.Context, exercise *user.Exercise) (*models.CreateExerciseResp, error) {
	resp := models.CreateExerciseResp{}

	err := s.DB.CreateExercise(ctx, exercise)
	if err != nil {
		return &resp, err
	}

	resp.Message = fmt.Sprintf("Exercise %v has successfully been created ", exercise.Name)
	resp.Exercise.Name = exercise.Name
	resp.Exercise.RestTime = exercise.RestTime
	resp.Exercise.BodyPart = exercise.BodyPart
	resp.Exercise.Equipment = exercise.Equipment
	resp.Exercise.CreatedAt = time.Now()

	return &resp, nil
}

