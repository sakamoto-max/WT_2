package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type conflictCache struct {
	client *redis.Client
}

func (c *conflictCache) SetConflictLevel(ctx context.Context, userId string, conflictLevel int) error {
	key := fmt.Sprintf("user_id:%v:conflict_level", userId)

	err := c.client.Set(ctx, key, conflictLevel, 0).Err()
	if err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error setting conflict status : %w", err))
	}

	return nil
}
func (c *conflictCache) GetConflictLevel(ctx context.Context, userId string) (int, error) {
	key := fmt.Sprintf("user_id:%v:conflict_level", userId)

	var conflictLevel int

	err := c.client.Get(ctx, key).Scan(&conflictLevel)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}

		return 0, myerrors.InternalServerErrMaker(fmt.Errorf("error getting conflict status : %w", err))
	}

	return conflictLevel, nil
}
