package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func initializePostgres(ctx context.Context) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_CONN"))

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

func RunMigrations() (error) {
	_, err := migrate.New("file://../migrations", os.Getenv("POSTGRES_CONN"))
	
	if err != nil{
		if errors.Is(err, migrate.ErrNoChange) {
			return nil	
		}
		return err
	}

	return nil
}