package cache

import (
	"context"
	"plan_service/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Metrics interface {
		GetRespTime(ctx context.Context) *time.Duration
	}
	EmptyPlan interface {
		SetUserEmptyPlanId(ctx context.Context, userId string, planId string)
		GetUserEmptyPlanId(ctx context.Context, userId string) (string, error)
		DelUserEmptyPlanId(ctx context.Context, userId string) error
	}

	PlanId interface {
		SetUserPlanId(ctx context.Context, payload domain.GetPlan, planId string)
		GetUserPlanId(ctx context.Context, payload domain.GetPlan) (string, error)
		DelUserPlanId(ctx context.Context, payload domain.GetPlan)
	}

	UserPlan interface {
		SetUserPlan(ctx context.Context, payload domain.Plan)
		GetUserPlan(ctx context.Context, payload domain.GetPlan) (string, *[]string, error)
		DelUserPlan(ctx context.Context, payload domain.GetPlan) (error)
	}
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		Metrics:   &metricsCache{client},
		EmptyPlan: &emptyPlanCache{client},
		PlanId:    &planIdCache{client},
		UserPlan:  &userPlanCache{client},
	}
}
