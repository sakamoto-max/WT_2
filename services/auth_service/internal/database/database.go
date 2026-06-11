package database

import (
	"auth_service/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func NewPgConn(config config.Config) (*pgxpool.Pool, error) {

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		config.Db.PgUser,
		config.Db.PgPass,
		config.Db.PgHost,
		config.Db.PgPort,
		config.Db.PgDatabaseName,
		config.Db.PgSSLMode,
	)

	pgConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres url : %w", err)
	}

	pgConfig.MaxConns = 10
	pgConfig.MaxConnLifetime = time.Duration(time.Minute * 10)

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool : %w", err)
	}

	ctxForPing, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err = pool.Ping(ctxForPing)
	if err != nil {
		return nil, fmt.Errorf("failed to ping to pg : %w", err)
	}

	return pool, nil
}

func NewRedisConn(config config.Config) (*redis.Client, error) {

	redisConnStr := fmt.Sprintf("redis://%s:%s@%s:%s/%s",
		config.Cache.RedisUserName,
		config.Cache.RedisPass,
		config.Cache.RedisHost,
		config.Cache.RedisPort,
		config.Cache.RedisDb,
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

// docker run -p 6001:6001 -e POSTGRES_CONN="postgresql://postgres:root@host.docker.internal:5432/auth?sslmode=disable" -e REDIS_ADDR="host.docker.internal:6379" -e REDIS_DB="0" -e SERVICE_NAME="auth_service" -e REDIS_PASS="" -e SECRET_KEY="asdfghjklazsxdc" -e GRPC_SERVER_ADDR="6001" -it auth_service
