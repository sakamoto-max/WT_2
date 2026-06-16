package database

import (
	"context"
	"fmt"
	"orchestration_service/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewPgConn(connStr string, config config.Config) *pgxpool.Pool {

	pgConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		config.Logger.Log.Fatalw("failed to parse pgConfig", zap.Error(err))
	}

	pgConfig.MaxConns = 10
	pgConfig.MaxConnLifetime = time.Duration(time.Minute * 10)

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		config.Logger.Log.Fatalw("error creating the postgres pool", zap.Error(err))
	}

	ctxForPing, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err = pool.Ping(ctxForPing)
	if err != nil {
		config.Logger.Log.Fatalw("failed to ping to pg", zap.Error(err))
	}

	return pool
}

func NewRedisConn(config config.Config) (*redis.Client, error) {

	redisConnStr := fmt.Sprintf("redis://%s:%s@%s:%s/%s",
		config.Cache.UserName,
		config.Cache.Password,
		config.Cache.Host,
		config.Cache.Port,
		config.Cache.Db,
	)

	ops, err := redis.ParseURL(redisConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis conn url : %w", err)
	}

	client := redis.NewClient(ops)

	ctxForPing, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	err = client.Ping(ctxForPing).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis : %w", err)
	}

	return client, nil
}
