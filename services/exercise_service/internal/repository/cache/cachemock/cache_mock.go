package cachemock

import (
	"context"
	"exercise_service/internal/domain"
	"exercise_service/internal/mappings"

	// "exercise_service/internal/domain/exercise"
	"fmt"
	"time"
)

type MetricsMock struct {
	Down bool
}

var (
	timeDuration time.Duration = time.Millisecond * 384
)

func (m *MetricsMock) GetRespTime(ctx context.Context) *time.Duration {
	if m.Down {
		return nil
	}

	return &timeDuration
}

type CrudMock struct {
	Down bool
	Hit  bool
	Miss bool
}

func (c *CrudMock) GetExerciseByNameR(ctx context.Context, payload mappings.GetExerciseByName) (*domain.Exercise, error) {
	if c.Down {
		return nil, fmt.Errorf("cache is down")
	}

	if c.Miss {
		return nil, nil
	}

	return &domain.Exercise{
		Id: "123", Name: "name", RestTime: 120, BodyPart: "body", Equipment: "equipment", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}, nil
}
func (c *CrudMock) SetExerciseByNameR(ctx context.Context, userId string, exerData *domain.Exercise) {
}

func (c *CrudMock) DeleteExerciseByNameR(ctx context.Context, payload mappings.DeleteExercise) error {

	if c.Down {
		return fmt.Errorf("redis is down")
	}

	return nil
}
