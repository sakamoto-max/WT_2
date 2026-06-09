package cachemock

import (
	"context"
	"fmt"
	"plan_service/internal/domain"
	"time"
)

type MetricsMock struct {
	Down bool
}

func (m *MetricsMock) GetRespTime(ctx context.Context) *time.Duration {
	if m.Down {
		return nil
	}

	timeDuration := time.Millisecond * 234

	return &timeDuration
}

type EmptyPlan struct {
	Down bool
	Hit  bool
}

func (e *EmptyPlan) SetUserEmptyPlanId(ctx context.Context, userId string, planId string) {}

func (e *EmptyPlan) GetUserEmptyPlanId(ctx context.Context, userId string) (string, error) {
	if e.Down {
		return "", fmt.Errorf("redis is down")
	}

	if !e.Hit {
		return "", nil
	}

	return "123", nil
}
func (e *EmptyPlan) DelUserEmptyPlanId(ctx context.Context, userId string) error {
	if e.Down {
		return fmt.Errorf("redis is down")
	}

	return nil
}

type PlanId struct {
	Down bool
	Hit  bool
}

func (p *PlanId) SetUserPlanId(ctx context.Context, payload domain.GetPlan, planId string) {}

func (p *PlanId) GetUserPlanId(ctx context.Context, payload domain.GetPlan) (string, error) {
	if p.Down {
		return "", fmt.Errorf("redis is down")
	}
	
	if !p.Hit {
		return "", nil
	}
	
	return "123", nil
}
func (p *PlanId) DelUserPlanId(ctx context.Context, payload domain.GetPlan) {}


type UserPlan struct{
	Down bool
	Hit bool
}

func (u *UserPlan) SetUserPlan(ctx context.Context, payload domain.Plan) {}

func (u *UserPlan) GetUserPlan(ctx context.Context, payload domain.GetPlan) (string, *[]string, error) {
	if u.Down {
		return "", nil, fmt.Errorf("redis is down")
	}

	if !u.Hit {
		return "", nil, nil
	}

	return "123", &[]string{"123", "123", "123"}, nil
}
func (u *UserPlan) DelUserPlan(ctx context.Context, payload domain.GetPlan) error  {
	return nil
}
