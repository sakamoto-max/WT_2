package cache

import (
	"context"
	"exercise_service/internal/domain"
	"exercise_service/internal/mappings"

	// "exercise_service/internal/domain/exercise"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type crudDB struct {
	client *redis.Client
}

func (c *crudDB) GetExerciseByNameR(ctx context.Context, payload mappings.GetExerciseByName) (*domain.Exercise, error) {

	key := fmt.Sprintf("user_id:%v:exercise_name:%v", payload.UserId, payload.ExerciseName)

	res, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get exercise by name from cache : %w", err)
	}

	if len(res) == 0 {
		return nil, nil
	}

	layout := "2006-01-02T15:04:05.9999999Z07:00"
	createdAt, err := time.Parse(layout, res[createdAt])
	updatedAt, err := time.Parse(layout, res[updatedAt])

	id := res[id]
	name := payload.ExerciseName
	bodyPart := res[bodyPart]
	equipment := res[equipment]

	data := domain.Exercise{
		Id:        id,
		Name:      name,
		BodyPart:  bodyPart,
		Equipment: equipment,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return &data, nil
}

func (c *crudDB) SetExerciseByNameR(ctx context.Context, userId string, exerData *domain.Exercise) {

	mainKey := fmt.Sprintf("user_id:%v:exercise_name:%v", userId, exerData.Name)
	idKey := "id"
	bodyPartKey := "body_part"
	equipmentKey := "equipment"
	createdAtKey := "created_at"
	updatedAtKey := "updated_at"

	c.client.HSet(ctx, mainKey,
		idKey, exerData.Id,
		bodyPartKey, exerData.BodyPart,
		equipmentKey, exerData.Equipment,
		createdAtKey, exerData.CreatedAt,
		updatedAtKey, exerData.UpdatedAt,
	)

}

func (c *crudDB) DeleteExerciseByNameR(ctx context.Context, payload mappings.DeleteExercise) error {
	key := fmt.Sprintf("user_id:%v:exercise_name:%v", payload.UserId, payload.ExerciseName)

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return myerrs.InternalServerErrMaker(fmt.Errorf("failed to delete the exercise from cache : %w", err))
	}

	return nil
}
