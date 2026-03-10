package database

import (
	"context"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func initializePostgres(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONN_STR"))
	if err != nil {
		return pool, err
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

	client := redis.NewClient(&ops)

	err := client.Ping(ctx).Err()
	if err != nil {
		return client, err
	}

	return client, nil
}

func InitializeDBs(ctx context.Context) (*pgxpool.Pool, *redis.Client, error) {

	var pool *pgxpool.Pool
	var client *redis.Client

	pool, err := initializePostgres(ctx)
	if err != nil {
		return pool, client, err
	}

	client, err = initializeRedis(ctx)
	if err != nil {
		return pool, client, err
	}

	return pool, client, nil
}

func CloseDBs(pool *pgxpool.Pool, client *redis.Client) {
	pool.Close()
	client.Close()

}
