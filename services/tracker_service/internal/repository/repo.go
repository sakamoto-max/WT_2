package repository

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
	"tracker_service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type RepoIface interface {
	GetPostgresRespTime(ctx context.Context) *time.Duration
	GetRedisRespTime(ctx context.Context) *time.Duration
	StartWorkout(ctx context.Context, userId string, planId string) (string, error)
	DeleteTrackerIdInPG(ctx context.Context, trackerId string) error
	RevertStartWorkout(ctx context.Context, trackerId string) error
	SetTrackerId(ctx context.Context, userId string, trackerId string) error
	GetTrackerId(ctx context.Context, userId string) (string, error)
	DelTrackerId(ctx context.Context, userId string) error
	EndWorkout(ctx context.Context, trackerId string, data *models.Tracker) error
	EndWorkoutWithOutbox(ctx context.Context, userId string, trackerId string, data *models.Tracker, planName string, newExerciseNames *[]string) error
	SetExerciseNameById(ctx context.Context, exerciseId string, exerciseName string) error
	GetExerciseNameById(ctx context.Context, exerciseId string) (string, error)
	SetUserCurrentPlanName(ctx context.Context, userId string, planName string) error
	GetUserCurrentPlanName(ctx context.Context, userId string) (string, error)
	SetPlanWithExercises(ctx context.Context, userId string, planName string, exerciseNames *[]string) error
	GetPlanWithExercises(ctx context.Context, userId string, planName string) (*[]string, error)
	SetUserWorkingOutWithPlan(ctx context.Context, userId string, value bool) error
	GetUserWorkingOutWithPlan(ctx context.Context, userId string) (bool, error)
	SetConflictLevel(ctx context.Context, userId string, conflictLevel int) error
	GetConflictLevel(ctx context.Context, userId string) (int, error)
	SetUserTrackerData(ctx context.Context, userId string, data *models.Tracker) error
	GetUserTrackerData(ctx context.Context, userId string) (*models.Tracker, error)
	SetUserNewExercises(ctx context.Context, userId string, exerciseNames *[]string) error
	GetUserNewExercises(ctx context.Context, userId string) (*[]string, error)
	DelAllUserData(ctx context.Context, userId string, planName string) error
}

type dBs struct {
	pDB *pgxpool.Pool
	rDB *redis.Client
}

func NewRepo() (RepoIface, error) {

	pool, err := newPgConn()
	if err != nil {
		return nil, err
	}

	client, err := newRedisConn()
	if err != nil {
		pool.Close()
		return nil, err
	}

	return &dBs{pDB: pool, rDB: client}, nil
}

func newPgConn() (*pgxpool.Pool, error) {

	pgConfig, err := pgxpool.ParseConfig(os.Getenv("POSTGRES_CONN"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgConfig : %w", err)
	}

	pgConfig.MaxConns = 10
	pgConfig.MaxConnLifetime = time.Duration(time.Minute * 10)

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating the postgres pool : %w\n", err)
	}

	ctxForPing, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	err = pool.Ping(ctxForPing)
	if err != nil {
		return nil, fmt.Errorf("Pg conn failed : %w", err)
	}

	return pool, nil
}

func newRedisConn() (*redis.Client, error) {

	database := os.Getenv("REDIS_DB")
	db, _ := strconv.Atoi(database)

	ops := redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       db,
	}
	client := redis.NewClient(&ops)

	ctxForPing, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	err := client.Ping(ctxForPing).Err()
	if err != nil {
		return nil, fmt.Errorf("error creating redis client : %w\n", err)
	}

	return client, nil

}

func (r *dBs) Close() error {
	r.pDB.Close()

	err := r.rDB.Close()
	if err != nil {
		return fmt.Errorf("error while closing redis db : %w", err)
	}

	return nil
}
