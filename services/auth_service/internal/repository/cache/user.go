package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	myErrs "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type userCache struct {
	client *redis.Client
}

func (c *userCache) UserLogout(ctx context.Context, userId string, uuid string) error {
	// del refresh
	refreshKey := fmt.Sprintf("%v_refresh", uuid)
	uuidKey := fmt.Sprintf("user_id:%v:uuid", userId)

	pipe := c.client.Pipeline()

	pipe.Del(ctx, refreshKey)
	pipe.Del(ctx, uuidKey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return myErrs.InternalServerErrMaker(fmt.Errorf("error deleting the refresh token after logout : %w\n", err))
	}

	return nil
}
