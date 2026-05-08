package repository

import (
	"context"
	"fmt"
	"os"
	"plan_service/internal/models"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type dBs struct {
	pDB *pgxpool.Pool
	rDB *redis.Client
}

type RepoIFace interface {
	CreatePlan(ctx context.Context, userId string, planName string, exerciseIds []string) error
	GetPlans(ctx context.Context, userId string) (*[]models.Plan3, error)
	GetAllExercisesByPlanID(ctx context.Context, planId string) (*[]string, error)
	ReturnsPlanId(ctx context.Context, userId string, planName string) (string, error)
	AddExercisesToPlan(ctx context.Context, planId string, exerciseIDs *[]string) error
	DeleteExerciseFromPlan(ctx context.Context, planId string, exerciseIDs *[]string) error
	DeletePlan(ctx context.Context, userId string, planId string) error
	CreateEmptyPlan(ctx context.Context, userId string) error
	GetPostgresRespTime(ctx context.Context) *time.Duration
	GetRedisRespTime(ctx context.Context) *time.Duration
	SetUserEmptyPlanIdR(ctx context.Context, userId string, planId string) error
	GetUserEmptyPlanIdR(ctx context.Context, userId string) (string, error)
	DelUserEmptyPlanIdR(ctx context.Context, userId string) error
	Close() error
	GetPlan(ctx context.Context, userId string, planName string) (string, *[]string, error)
	SetUserPlanId(ctx context.Context, userId string, planName string, planId string) error
	GetUserPlanId(ctx context.Context, userId string, planName string) (string, error)
	DelUserPlanId(ctx context.Context, userId string, planName string) error
	SetUserPlan(ctx context.Context, userId string, planName string, planId string, exerciseIds *[]string) error 
	GetUserPlan(ctx context.Context, userId string, planName string) (string, *[]string, error) 
	DelUserPlan(ctx context.Context, userId string, planName string) error
}

func NewRepo() (RepoIFace, error) {

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
		return fmt.Errorf("error closing the redis Db : %w", err)
	}

	return nil
}
