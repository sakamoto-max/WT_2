package mock

import (
	"context"
	"exercise_service/internal/domain"
	"exercise_service/internal/mappings"
	"fmt"
	"time"
)

type MetricsDbMock struct {
	HasErr bool
}

var (
	timeDuration time.Duration = time.Millisecond * 300
)

func (m *MetricsDbMock) GetRespTime(ctx context.Context) *time.Duration {
	if m.HasErr {
		return nil
	}

	return &timeDuration
}

type ExerciseCdMock struct {
	HasErr bool
	ExerciseExists bool
}


func (e *ExerciseCdMock) CreateExercise(ctx context.Context, payload mappings.CreateExercise) (string, error) {
	if e.HasErr {
		return "", fmt.Errorf("some error occured")
	}

	if e.ExerciseExists {
		return "", fmt.Errorf("exercise already exists")
	}

	return "123", nil
}

func (e *ExerciseCdMock) DeleteExecise(ctx context.Context, payload mappings.DeleteExercise) error {

	if e.HasErr {
		return fmt.Errorf("some error occured")
	}

	if !e.ExerciseExists {
		return fmt.Errorf("exercise doesn't exits")
	}

	return nil
}

type ExerciseGetMock struct {
	HasErr bool
	ExerciseExits bool
}

func (e *ExerciseGetMock) GetExerciseByName(ctx context.Context, payload mappings.GetExerciseByName) (*domain.Exercise, error) {
	if e.HasErr{
		return nil, fmt.Errorf("some error occured")
	}

	if !e.ExerciseExits {
		return nil, fmt.Errorf("exercise doesnot exits")
	}

	return &domain.Exercise{
		Id: "123", Name: "name", RestTime: 120, BodyPart: "body", Equipment: "equipment", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}, nil
}
func (e *ExerciseGetMock) GetAllExercises(ctx context.Context, payload mappings.GetAllExercises) (*[]domain.Exercise, error) {
	if e.HasErr{
		return nil, fmt.Errorf("some error occured")
	}

	exer := domain.Exercise{Id : "123",Name: "name"}

	return &[]domain.Exercise{exer, exer, exer, exer}, nil
}
func (e *ExerciseGetMock) GetExerciseNameByID(ctx context.Context, exerciseId string) (string, error) {
	if e.HasErr{
		return "", fmt.Errorf("some error occured")
	}

	if !e.ExerciseExits {
		return "", fmt.Errorf("exercise doesn't exist")
	}

	return "123", nil
}
func (e *ExerciseGetMock) ExerciseExistsReturnId(ctx context.Context, payload mappings.ExerciseExistsReturnId) (string, error) {
	if e.HasErr{
		return "", fmt.Errorf("some error occured")
	}

	if !e.ExerciseExits {
		return "", fmt.Errorf("exercise doesn't exist")
	}

	return "123", nil
}
