package repository

import (
	"context"
	"errors"
	"fmt"
	// myerrors "wt/pkg/my_errors"
	myerrors "github.com/sakamoto-max/wt_2-pkg/my_errors"
	"github.com/redis/go-redis/v9"
)


func (r *repo) RefreshExists(ctx context.Context, userId string) (bool, error) {

	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)

	uuid, err := r.rDB.Get(ctx, uuidKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, myerrors.InternalServerErrMaker(fmt.Errorf("error checking if uuid exists for user : %v", userId))
	}
	
	refreshKey := fmt.Sprintf("%v_refresh", uuid)

	_, err = r.rDB.Get(ctx, refreshKey).Result()
	if err != nil {
		return false, myerrors.InternalServerErrMaker(fmt.Errorf("error getting the refresh token for the user : %v", userId))
	}

	return true, nil
}
func (r *repo) GetUUID(ctx context.Context, userId string) (string, error) {

	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)

	uuid, err := r.rDB.Get(ctx, uuidKey).Result()
	if err != nil {
		return uuid, myerrors.InternalServerErrMaker(fmt.Errorf("error getting the UUID : %w", err))
	}

	return uuid, nil
}
func (r *repo) SetUUID(ctx context.Context, uuid string, userId string) error {

	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)

	err := r.rDB.Set(ctx, uuidKey, uuid, 0).Err()
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error setting the uuid : %w", err))
	}

	return nil
}
func (r *repo) SetRefreshTokenAndUUID(ctx context.Context, uuid string, Refreshtoken string, userId string) error {
	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)
	refreshKey := fmt.Sprintf("%v_refresh", uuid)

	pipe := r.rDB.Pipeline()

	pipe.Set(ctx, uuidKey, uuid, 0)
	pipe.Set(ctx, refreshKey, Refreshtoken, 0)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error setting refresh and uuid for the user %v : %w\n", userId, err))
	}

	return nil
}
func (r *repo) GetRefreshToken(ctx context.Context, uuid string) (string, error) {

	var refreshToken string

	key := fmt.Sprintf("%v_refresh", uuid)

	err := r.rDB.Get(ctx, key).Scan(&refreshToken)
	if err != nil {
		return "", myerrors.InternalServerErrMaker(fmt.Errorf("error getting the refresh token : %w\n", err))
	}

	return refreshToken, nil
}
