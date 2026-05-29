package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	myErrs "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type tokenCache struct {
	client *redis.Client
}


func (c *tokenCache) RefreshExists(ctx context.Context, userId string) (bool, error) {

	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)

	uuid, err := c.client.Get(ctx, uuidKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, myErrs.InternalServerErrMaker(fmt.Errorf("error checking if uuid exists for user : %v", userId))
	}

	refreshKey := fmt.Sprintf("%v_refresh", uuid)

	_, err = c.client.Get(ctx, refreshKey).Result()
	if err != nil {
		return false, myErrs.InternalServerErrMaker(fmt.Errorf("error getting the refresh token for the user : %v", userId))
	}

	return true, nil
}

func (c *tokenCache) SetRefreshTokenAndUUID(ctx context.Context, uuid string, Refreshtoken string, userId string) error {
	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)
	refreshKey := fmt.Sprintf("%v_refresh", uuid)

	pipe := c.client.Pipeline()

	pipe.Set(ctx, uuidKey, uuid, time.Hour*24*30)
	pipe.Set(ctx, refreshKey, Refreshtoken, time.Hour*24*30)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return myErrs.InternalServerErrMaker(fmt.Errorf("error setting refresh and uuid for the user %v : %w\n", userId, err))
	}

	return nil
}
func (c *tokenCache) GetRefreshToken(ctx context.Context, uuid string) (string, error) {

	var refreshToken string

	key := fmt.Sprintf("%v_refresh", uuid)

	err := c.client.Get(ctx, key).Scan(&refreshToken)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", myErrs.BadReqErrMaker(fmt.Errorf("refresh token does not exist"))
		}
		return "", myErrs.InternalServerErrMaker(fmt.Errorf("error getting the refresh token : %w\n", err))
	}

	return refreshToken, nil
}
