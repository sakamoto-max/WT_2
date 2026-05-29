package cache

import (
	"context"
	"time"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Token interface {
		RefreshExists(ctx context.Context, userId string) (bool, error)
		SetRefreshTokenAndUUID(ctx context.Context, uuid string, Refreshtoken string, userId string) error
		GetRefreshToken(ctx context.Context, uuid string) (string, error)
	}
	Metrics interface {
		GetRespTime(ctx context.Context) *time.Duration
	}
	Uuid interface {
		GetUUID(ctx context.Context, userId string) (string, error)
		SetUUID(ctx context.Context, uuid string, userId string) error
	}
	User interface {
		UserLogout(ctx context.Context, userId string, uuid string) error
	}
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		Token:   &tokenCache{client},
		Metrics: &metrics{client},
		Uuid:    &UuidCache{client},
		User:    &userCache{client},
	}
}

