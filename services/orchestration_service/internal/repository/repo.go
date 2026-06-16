package repository

import (
	"context"
	"orchestration_service/internal/repository/cache"
	"orchestration_service/internal/types"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const queryExecutionTime = time.Second * 3

type Db struct {
	Auth interface {
		FetchData(ctx context.Context, wg *sync.WaitGroup, dataChan chan<- *[]types.Data)
		UpdateTaskStatus(ctx context.Context, dbIndex string, updateValue string) error
		UpdateTaskStatusWithNumberOfTries(ctx context.Context, dbIndex string, updateValue string) error
	}
	Tracker interface {
		FetchData(ctx context.Context, wg *sync.WaitGroup, dataChan chan<- *[]types.Data)
		UpdateTaskStatus(ctx context.Context, dbIndex string, updateValue string) error
		UpdateTaskStatusWithNumberOfTries(ctx context.Context, dbIndex string, updateValue string) error
	}
	Cache interface {
		SetTaskTimeOut(ctx context.Context, data types.Data) error
		SkipTask(ctx context.Context, data types.Data) (bool, error)
	}
}

func NewDb(auth *database, tracker *database, redisClient *redis.Client) Db {
	return Db{
		Auth:    auth,
		Tracker: tracker,
		Cache:   cache.NewCache(redisClient),
	}
}

type database struct {
	pg     *pgxpool.Pool
	dbName string
}

func RegisterDb(pool *pgxpool.Pool, dbName string) *database {
	return &database{
		pg:     pool,
		dbName: dbName,
	}
}
