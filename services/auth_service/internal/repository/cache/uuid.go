package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	myErrs "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type UuidCache struct {
	client *redis.Client
}

func (c *UuidCache) GetUUID(ctx context.Context, userId string) (string, error) {

	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)

	uuid, err := c.client.Get(ctx, uuidKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", myErrs.BadReqErrMaker(fmt.Errorf("please login first"))
		}
		return uuid, myErrs.InternalServerErrMaker(fmt.Errorf("error getting the UUID : %w", err))
	}

	return uuid, nil
}
func (c *UuidCache) SetUUID(ctx context.Context, uuid string, userId string) error {

	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)

	err := c.client.Set(ctx, uuidKey, uuid, 0).Err()
	if err != nil {
		return myErrs.InternalServerErrMaker(fmt.Errorf("error setting the uuid : %w", err))
	}

	return nil
}
