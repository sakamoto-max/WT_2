package services

import (
	"context"
	"exercise_service/internal/models"
	"exercise_service/internal/repository"
	// "exercise_service/internal/user"
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

func (s *Service) CreateExerciseSer(ctx context.Context, exerciseName string, bodyPartName string, equipmentName string, restTime int) (*models.CreateExerciseResp, error) {
	resp := models.CreateExerciseResp{}

	if restTime == 0 {
		restTime = 120
	}

	err := s.DB.CreateExercise(ctx, exerciseName, bodyPartName, equipmentName, restTime)
	if err != nil {
		return &resp, err
	}

	resp.Message = fmt.Sprintf("Exercise %v has successfully been created ", exerciseName)
	resp.Exercise.Name = exerciseName
	resp.Exercise.RestTime = restTime
	resp.Exercise.BodyPart = bodyPartName
	resp.Exercise.Equipment = equipmentName
	resp.Exercise.CreatedAt = time.Now()

	return &resp, nil
}

func (s *Service) ExerciseExistsReturnId(ctx context.Context, exerciseName string) (bool, int32, error) {
	return s.DB.ExerciseExistsReturnId(ctx, exerciseName)
}

func (s *Service)  GetExerciseNameByID(ctx context.Context, exerciseId int) (string, error) {
	return s.DB.GetExerciseNameByID(ctx, exerciseId)
}

