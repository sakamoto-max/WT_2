package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repo struct {
	pDB *pgxpool.Pool
	rDB *redis.Client
}

func NewRepo(pool *pgxpool.Pool, client *redis.Client) *Repo {
	return &Repo{pDB: pool, rDB: client}
}



func (r *Repo) GetPostgresRespTime(ctx context.Context) (*time.Duration) {
	timeStart := time.Now()
	err := r.pDB.Ping(ctx)
	if err != nil{
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
func (r *Repo) GetRedisRespTime(ctx context.Context) (*time.Duration) {
	timeStart := time.Now()
	err := r.rDB.Ping(ctx).Err()
	if err != nil{
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}

