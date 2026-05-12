package services

import (
	"context"
	"time"
	"tracker_service/internal/models"
	"tracker_service/internal/repository"

	exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"
	planpb "github.com/sakamoto-max/wt_2_proto/shared/plan"
)

type service struct {
	db      repository.RepoIface
	pClient planpb.PlanServiceClient
	eClient exerpb.ExerciseServiceClient
}

func NewService(Db repository.RepoIface, planClient planpb.PlanServiceClient, exerClient exerpb.ExerciseServiceClient) ServiceIface {
	return &service{db: Db, pClient: planClient, eClient: exerClient}
}

type ServiceIface interface {
	GetHealth(ctx context.Context) (*time.Duration, *time.Duration)
	StartEmptyWorkoutSer(ctx context.Context, userID string) error
	StartWorkoutWithPlanSer(ctx context.Context, userId string, planName string) (*[]string, error)
	EndWorkoutSer(ctx context.Context, userId string, data *models.Tracker) (*string, error)
	CancelWorkout(ctx context.Context, userID string) error
}
