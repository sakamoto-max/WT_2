package repository

import (

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


