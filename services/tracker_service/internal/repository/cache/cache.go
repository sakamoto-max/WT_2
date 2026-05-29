package cache

import (
	"context"
	"time"
	"tracker_service/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Metrics interface { 
		GetRespTime(ctx context.Context) *time.Duration
	}
	UserData interface {
		DelAllUserData(ctx context.Context, userId string, planName string) error
	}
	TrackerId interface { 
		SetTrackerId(ctx context.Context, userId string, trackerId string) error
		GetTrackerId(ctx context.Context, userId string) (string, error)
		DelTrackerId(ctx context.Context, userId string) error
	}
	ExerciseName interface {
		SetExerciseNameById(ctx context.Context, exerciseId string, exerciseName string) error
		GetExerciseNameById(ctx context.Context, exerciseId string) (string, error)
	}
	CurrentPlan interface {
		SetUserCurrentPlanName(ctx context.Context, userId string, planName string) error
		GetUserCurrentPlanName(ctx context.Context, userId string) (string, error)
	}
	Plan interface {
		SetPlanWithExercises(ctx context.Context, userId string, planName string, exerciseNames *[]string) error
		GetPlanWithExercises(ctx context.Context, userId string, planName string) (*[]string, error)
	}
	WithPlan interface {
		SetUserWorkingOutWithPlan(ctx context.Context, userId string, value bool) error
		GetUserWorkingOutWithPlan(ctx context.Context, userId string) (bool, error)
	}
	Conflict interface {
		SetConflictLevel(ctx context.Context, userId string, conflictLevel int) error
		GetConflictLevel(ctx context.Context, userId string) (int, error)
	}
	TrackerData interface {
		SetUserTrackerData(ctx context.Context, userId string, data *domain.Tracker) error
		GetUserTrackerData(ctx context.Context, userId string) (*domain.Tracker, error)
	}
	NewExercises interface {
		SetUserNewExercises(ctx context.Context, userId string, exerciseNames *[]string) error
		GetUserNewExercises(ctx context.Context, userId string) (*[]string, error)
	}
}



func NewCache(client *redis.Client) *Cache {
	return &Cache{
		Metrics: &metricsCache{client},
		NewExercises: &newExerCache{client},
		TrackerData: &trackerDataCache{client},
		Conflict: &conflictCache{client},
		WithPlan: &withPlanCache{client},
		Plan: &planCache{client},
		CurrentPlan: &CurrentPlanCache{client},
		UserData: &userDataCache{client},
		ExerciseName: &exerciseNameCache{client},
		TrackerId: &trackerIdCache{client},
	}
}













