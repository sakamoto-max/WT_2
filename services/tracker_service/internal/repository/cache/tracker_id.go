package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type trackerIdCache struct {
	client *redis.Client
}

func (c *trackerIdCache) SetTrackerId(ctx context.Context, userId string, trackerId string) error {
	keyforTrackId := fmt.Sprintf("user:%v:tracker_id", userId)

	if err := c.client.Set(ctx, keyforTrackId, trackerId, 0).Err(); err != nil {
		return fmt.Errorf("error setting the tracker id : %w", err)
	}

	return nil

}

func (c *trackerIdCache) GetTrackerId(ctx context.Context, userId string) (string, error) {
	var id string
	key := fmt.Sprintf("user:%v:tracker_id", userId)
	id, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return id, nil
		}
		return id, myerrors.InternalServerErrMaker(fmt.Errorf("error in getting the tracker Id of the user with id %v : %w", userId, err))
	}

	return id, nil

}

func (c *trackerIdCache) DelTrackerId(ctx context.Context, userId string) error {
	key := fmt.Sprintf("user:%v:tracker_id", userId)

	if err := c.client.Del(ctx, key).Err(); err != nil {
		return myerrors.InternalServerErrMaker(fmt.Errorf("error deleting tracker id in redis : %w", err))
	}

	return nil
}
