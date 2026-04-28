package repository

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
	"github.com/sakamoto-max/wt_2-pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const queryExecutionTime = time.Second * 3

type DB struct {
	AuthPg *pgxpool.Pool
	TrackerPg *pgxpool.Pool
	Redis  *redis.Client
}

func NewDBs(logger *logger.MyLogger) (*DB, error) {

	authPool, err := newAuthPgConn()
	if err != nil {
		return nil, err
	}

	trackerPool, err := newTrackerPgConn()
	if err != nil {
		return nil, err
	}

	logger.Log.Infoln("connected to postgres")
	
	client, err := newRedisConn()
	if err != nil{
		authPool.Close()
		trackerPool.Close()
		return nil, err
	}
	
	logger.Log.Infoln("connected to redis")

	return &DB{AuthPg: authPool, TrackerPg: trackerPool, Redis: client}, nil
}

func CloseDBs(pool *pgxpool.Pool, client *redis.Client) {
	pool.Close()
	client.Close()
}

func newAuthPgConn() (*pgxpool.Pool, error) {
	pgConfig, err := pgxpool.ParseConfig(os.Getenv("AUTH_POSTGRES_CONN"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgConfig : %w", err)
	}

	pgConfig.MaxConns = 10
	pgConfig.MaxConnLifetime = time.Duration(time.Minute * 10)

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating the postgres pool for auth : %w\n", err)
	}

	ctxForPing, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	err = pool.Ping(ctxForPing)
	if err != nil {
		return nil, fmt.Errorf("Pg conn failed : %w", err)
	}

	return pool, nil
}
func newTrackerPgConn() (*pgxpool.Pool, error) {
	pgConfig, err := pgxpool.ParseConfig(os.Getenv("TRACKER_POSTGRES_CONN"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgConfig : %w", err)
	}

	pgConfig.MaxConns = 10
	pgConfig.MaxConnLifetime = time.Duration(time.Minute * 10)

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating the postgres pool for tracker : %w\n", err)
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

	ctxForPing, cancel := context.WithTimeout(context.Background(), time.Second * 1)
	defer cancel()

	err := client.Ping(ctxForPing).Err()
	if err != nil {
		return nil, fmt.Errorf("error creating redis client : %w\n", err)
	}

	return client, nil
}
