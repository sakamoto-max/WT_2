package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type DBs struct {
	PDB *pgxpool.Pool
	RDB *redis.Client
}

func NewDBs(pool *pgxpool.Pool, client *redis.Client) *DBs {
	return &DBs{PDB: pool, RDB: client}
}

func (r *DBs) GetPostgresRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := r.PDB.Ping(ctx)
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
func (r *DBs) GetRedisRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := r.RDB.Ping(ctx).Err()
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
