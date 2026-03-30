package services

import (
	"context"
	"exercise_service/internal/models"
	"exercise_service/internal/repository"
	// "exercise_service/internal/user"
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

func (s *Service) GetExerciseByNameSer(ctx context.Context, userID string, exerciseName string) (*models.Exercise2, error) {
	// check if the execise exists
	// transform the exercises
	exercise, err := s.DB.GetExerciseByName(ctx, userID, exerciseName)
	if err != nil {
		return exercise, err
	}
	return exercise, nil
}

func (s *Service) GetAllExercisesSer(ctx context.Context, userId string) (*[]models.Exercise2, error) {
	allExercises, err := s.DB.GetAllExercises(ctx, userId)
	if err != nil {
		return allExercises, err
	}

	return allExercises, nil
}

func (s *Service) DeleteExeciseSer(ctx context.Context, userId string, exerciseName string) error {
	err := s.DB.DeleteExecise(ctx, userId, exerciseName)
	if err != nil{
		return err
	}

	return nil
}


func (s *Service) CreateExerciseSer(ctx context.Context, userId string, exerciseName string, bodyPartName string, equipmentName string) (string, error) {

	UUId, err := s.DB.CreateExercise(ctx, userId ,exerciseName, bodyPartName, equipmentName)
	if err != nil {
		return "", err
	}

	return UUId, nil
}

func (s *Service) ExerciseExistsReturnId(ctx context.Context, userId string, exerciseName string) (string, error) {
	return s.DB.ExerciseExistsReturnId(ctx, userId, exerciseName)
}

func (s *Service)  GetExerciseNameByID(ctx context.Context, exerciseId string) (string, error) {
	return s.DB.GetExerciseNameByID(ctx, exerciseId)
}


func (s *Service) GetHealth(ctx context.Context) (*time.Duration, *time.Duration) {

	// check resp time of pg

	pgRespTime := s.DB.GetPostgresRespTime(ctx)
	redisRespTime := s.DB.GetRedisRespTime(ctx)

	return pgRespTime, redisRespTime
}