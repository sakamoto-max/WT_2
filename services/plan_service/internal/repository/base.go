package repository

import (
	"context"
	"plan_service/internal/domain"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
)


type Db struct {
	PlanQueryRepo interface {
		GetAllPlanNamesWithIds(ctx context.Context, userId string) (*[]domain.Plan, error)
		GetAllExercisesByPlanID(ctx context.Context, planId string) (*[]string, error)
		GetPlan(ctx context.Context, payload domain.GetPlan) (string, *[]string, error)
		GetPlanId(ctx context.Context, payload domain.GetPlan) (string, error)
		GetEmptyPlanId(ctx context.Context, userId string) (string, error) 
	}
	PlanCommandRepo interface {
		DeletePlan(ctx context.Context, userId string, planId string) error
		CreateEmptyPlan(ctx context.Context, userId string) error
		CreatePlan(ctx context.Context, payload domain.CreatePlan) error
	}
	PlanExericseRepo interface {
		RemoveExerciseFromPlan(ctx context.Context, planId string, exerciseIDs *[]string) error
		AddExercisesToPlan(ctx context.Context, planId string, exerciseIDs *[]string) error
	}
	MetricsRepo interface {
		GetRespTime(ctx context.Context) *time.Duration
	}
	QueueDb interface {
		Insert(data mqTypes.Data) error
		Fetch() (*[]mqTypes.Data, error)
	}
}

func NewDb(pg *pgxpool.Pool) *Db {
	return &Db{
		PlanQueryRepo: &planQueryRepo{pg},
		PlanCommandRepo: &planCommandRepo{pg},
		PlanExericseRepo: &planExerciseRepo{pg},
		MetricsRepo: &metricsRepo{pg},
		QueueDb: &queueDb{pg},
	}
}

