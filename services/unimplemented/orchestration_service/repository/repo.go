package repository

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type DB struct {
	AuthPg *pgxpool.Pool
	Redis  *redis.Client
}

func initializeAuthPostgres(ctx context.Context) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(ctx, os.Getenv("AUTH_POSTGRES_CONN"))

	if err != nil {
		return pool, fmt.Errorf("error creating the postgres pool : %w\n", err)
	}

	return pool, nil
}

func initializeRedis(ctx context.Context) (*redis.Client, error) {

	database := os.Getenv("REDIS_DB")
	db, _ := strconv.Atoi(database)

	ops := redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       db,
	}
	rdb := redis.NewClient(&ops)

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return rdb, fmt.Errorf("error creating redis client : %w\n", err)
	}

	return rdb, nil
}

func InitializeDBs(ctx context.Context) (*DB, error) {

	poolForAuth, err := initializeAuthPostgres(ctx)
	if err != nil {
		return nil, err
	}

	client, err := initializeRedis(ctx)
	if err != nil {
		return nil, err
	}

	return &DB{AuthPg: poolForAuth, Redis: client}, nil
}

func CloseDBs(pool *pgxpool.Pool, client *redis.Client) {
	pool.Close()
	client.Close()
}
