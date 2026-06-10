package database

import (
	"context"
	"orchestration_service/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewPgConn(connStr string, config config.Config) *pgxpool.Pool {

	pgConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		config.Logger.Log.Fatalw("failed to parse pgConfig", zap.Error(err))
	}

	pgConfig.MaxConns = 10
	pgConfig.MaxConnLifetime = time.Duration(time.Minute * 10)

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		config.Logger.Log.Fatalw("error creating the postgres pool", zap.Error(err))
	}

	ctxForPing, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err = pool.Ping(ctxForPing)
	if err != nil {
		config.Logger.Log.Fatalw("failed to ping to pg", zap.Error(err))
	}

	return pool
}
