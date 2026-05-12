package services

import (
	"context"
	"plan_service/internal/models"
	"plan_service/internal/repository"
	"time"
	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
)

type ServiceIface interface {
	CreatePlan(ctx context.Context, userId string, planName string, exerciseNames *[]string) error
	GetPlans(ctx context.Context, userId string) (int, *[]models.Plan2, error)
	GetPlan(ctx context.Context, userId string, planName string) (string, string, *[]string, error)
	AddExercises(ctx context.Context, userId string, planName string, exerciseNames *[]string) (*models.Plan2, error)
	DeleteExerciseFromPlan(ctx context.Context, userId string, planName string, exerciseNames *[]string) (*models.Plan2, error)
	DeletePlan(ctx context.Context, userId string, planName string) error
	GetEmptyPlanId(ctx context.Context, userId string) (string, error)
	GetHealth(ctx context.Context) (*time.Duration, *time.Duration)
}

type service struct {
	Db      repository.RepoIFace
	GClient exerpb.ExerciseServiceClient
}

func NewService(Db repository.RepoIFace, grpcCli exerpb.ExerciseServiceClient) ServiceIface {
	return &service{Db: Db, GClient: grpcCli}
}
