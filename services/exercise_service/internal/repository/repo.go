package repository

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repo struct {
	pDB *pgxpool.Pool
	rDB *redis.Client
}

func NewRepo() (*Repo, error) {

	pool, err := newPgConn()
	if err != nil {
		return nil, err
	}

	client, err := newRedisConn()
	if err != nil {
		pool.Close()
		return nil, err
	}

	return &Repo{pDB: pool, rDB: client}, nil
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

func (r *Repo) Close() error {
	r.pDB.Close()

	if err := r.rDB.Close(); err != nil {
		return fmt.Errorf("error closing the redis Db : %w", err)
	}

	return nil
}
