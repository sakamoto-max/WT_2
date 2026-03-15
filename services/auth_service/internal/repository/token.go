package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func (r *Repo) RefreshExists(ctx context.Context, userId int) (bool, error) {

	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)
	// get uuid

	uuid, err := r.rDB.Get(ctx, uuidKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, fmt.Errorf("error checking if uuid exists for user : %v", userId)
	}

	refreshKey := fmt.Sprintf("%v_refresh", uuid)

	_, err = r.rDB.Get(ctx, refreshKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, fmt.Errorf("error getting the refresh token for the user : %v", userId)
	}

	return true, nil
}

func (r *Repo) GetUUID(ctx context.Context, userId int) (string, error) {

	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)

	uuid, err := r.rDB.Get(ctx, uuidKey).Result()
	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func (r *Repo) SetUUID(ctx context.Context, uuid uuid.UUID, userId int) error {

	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)

	err := r.rDB.Set(ctx, uuidKey, uuid, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting the uuid : %w", err)
	}

	return nil
}

func (r *Repo) SetRefreshTokenAndUUID(ctx context.Context, uuid string, Refreshtoken string, userId int) error {
	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)
	refreshKey := fmt.Sprintf("%v_refresh", uuid)

	pipe := r.rDB.Pipeline()

	pipe.Set(ctx, uuidKey, uuid, 0)
	pipe.Set(ctx, refreshKey, Refreshtoken, 0)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error setting refresh and uuid for the user %v : %w\n", userId, err)
	}

	return nil
}


func (r *Repo) GetRefreshToken(ctx context.Context, uuid uuid.UUID) (string, error) {

	var refreshToken string

	key := fmt.Sprintf("%v_refresh", uuid)

	err := r.rDB.Get(ctx, key).Scan(&refreshToken)
	if err != nil {
		return "", fmt.Errorf("error getting the refresh token : %w\n", err)
	}

	return refreshToken, nil
}
