package mock

import (
	"context"
	"fmt"
	"plan_service/internal/domain"
	"time"
)

type PlanQueryMock struct {
	PgDown    bool
	PlanExits bool
}

func (p *PlanQueryMock) GetAllPlanNamesWithIds(ctx context.Context, userId string) (*[]domain.Plan, error) {
	if p.PgDown {
		return nil, fmt.Errorf("pg is down")
	}

	plan := domain.Plan{
		PlanName: "plan_name",
		PlanId:   "plan_id",
	}

	return &[]domain.Plan{plan, plan, plan}, nil
}

func (p *PlanQueryMock) GetAllExercisesByPlanID(ctx context.Context, planId string) (*[]string, error) {
	if p.PgDown {
		return nil, fmt.Errorf("pg is down")
	}

	if !p.PlanExits {
		return nil, fmt.Errorf("plan doesn't exist")
	}

	exerId := "123"

	return &[]string{exerId, exerId, exerId}, nil
}

func (p *PlanQueryMock) GetPlan(ctx context.Context, payload domain.GetPlan) (string, *[]string, error) {
	if p.PgDown {
		return "", nil, fmt.Errorf("pg is down")
	}

	if !p.PlanExits {
		return "", nil, fmt.Errorf("plan doesn't exist")
	}

	exerId := "123"
	planId := "123"

	return planId, &[]string{exerId, exerId, exerId}, nil
}

func (p *PlanQueryMock) GetPlanId(ctx context.Context, payload domain.GetPlan) (string, error) {
	if p.PgDown {
		return "", fmt.Errorf("pg is down")
	}

	if !p.PlanExits {
		return "", fmt.Errorf("plan doesn't exist")
	}

	return "123", nil
}
func (p *PlanQueryMock) GetEmptyPlanId(ctx context.Context, userId string) (string, error) {
	if p.PgDown {
		return "", fmt.Errorf("pg is down")
	}

	if !p.PlanExits {
		return "", fmt.Errorf("plan doesn't exist")
	}

	return "123", nil
}


type PlanCommandMock struct {
	PgDown     bool
	PlanExists bool
}

func (p *PlanCommandMock) DeletePlan(ctx context.Context, userId string, planId string) error {
	if p.PgDown {
		return fmt.Errorf("pg is down")
	}

	if !p.PlanExists {
		return fmt.Errorf("plan doesn't exist")
	}

	return nil
}
func (p *PlanCommandMock) CreateEmptyPlan(ctx context.Context, userId string) error {
	if p.PgDown {
		return fmt.Errorf("pg is down")
	}

	return nil
}
func (p *PlanCommandMock) CreatePlan(ctx context.Context, payload domain.CreatePlan) error {
	if p.PgDown {
		return fmt.Errorf("pg is down")
	}

	if p.PlanExists {
		return fmt.Errorf("plan already exits")
	}

	return nil
}

type PlanExericseMock struct {
	PgDown bool
}

func (p *PlanExericseMock) RemoveExerciseFromPlan(ctx context.Context, planId string, exerciseIDs *[]string) error {
	if p.PgDown {
		return fmt.Errorf("pg is down")
	}

	return nil

}
func (p *PlanExericseMock) AddExercisesToPlan(ctx context.Context, planId string, exerciseIDs *[]string) error {
	if p.PgDown {
		return fmt.Errorf("pg is down")
	}

	return nil
}

type MetricsMock struct {
	PgDown bool
}

func (m *MetricsMock) GetRespTime(ctx context.Context) *time.Duration {
	if m.PgDown {
		return nil
	}

	timeDuration := time.Millisecond * 272

	return &timeDuration
}
