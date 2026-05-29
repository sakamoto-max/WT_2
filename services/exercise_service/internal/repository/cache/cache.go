package cache

import (
	"context"
	"exercise_service/internal/domain"
	"exercise_service/mappings"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	id        string = "id"
	bodyPart  string = "body_part"
	equipment string = "equipment"
	createdAt string = "created_at"
	updatedAt string = "updated_at"
)


type Cache struct {
	Metrics interface {
		GetRespTime(ctx context.Context) *time.Duration
	}
	CRUD interface {
		GetExerciseByNameR(ctx context.Context, payload mappings.GetExerciseByName) (*domain.Exercise, error)
		SetExerciseByNameR(ctx context.Context, userId string, exerData *domain.Exercise)
		DeleteExerciseByNameR(ctx context.Context, payload mappings.DeleteExercise) error 
	}
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		Metrics: &metricsDb{client},
		CRUD: &crudDB{client},
	}
}
