package repository

import (
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

