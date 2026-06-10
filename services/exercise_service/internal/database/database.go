package database

import (
	"context"
	"exercise_service/internal/config"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewPgConn(config config.Config) *pgxpool.Pool {

	// "postgresql://postgres:root@localhost:5432/WT_EXERCISES?sslmode=disable"

	dbConnStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		config.Db.PgUser,
		config.Db.PgPass,
		config.Db.PgHost,
		config.Db.PgPort,
		config.Db.PgDatabaseName,
		config.Db.PgSSLMode,
	)

	pgConfig, err := pgxpool.ParseConfig(dbConnStr)
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
		config.Logger.Log.Fatalw("failed to ping pg", zap.Error(err))
	}

	return pool
}

func NewRedisConn(config config.Config) (*redis.Client) {

	redisConnStr := fmt.Sprintf("redis://%s:%s@%s:%s/%s",
		config.Cache.RedisUserName,
		config.Cache.RedisPass,
		config.Cache.RedisHost,
		config.Cache.RedisPort,
		config.Cache.RedisDb,
	)

	ops, err := redis.ParseURL(redisConnStr)
	if err != nil {
		config.Logger.Log.Fatalw("failed to parse redis connection url", zap.Error(err))
	}

	client := redis.NewClient(ops)

	ctxForPing, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	err = client.Ping(ctxForPing).Err()
	if err != nil {
		config.Logger.Log.Fatalw("failed to ping redis", zap.Error(err))
	}

	return client
}
