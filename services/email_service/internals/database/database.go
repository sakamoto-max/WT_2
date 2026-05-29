package database

import (
	"context"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"go.uber.org/zap"
)

func NewDb(connString string, logger *logger.MyLogger) *pgxpool.Pool {
	pgConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		logger.Log.Fatalw("failed to parse the connection string", zap.Error(err))
		return nil
	}

	pgConfig.MinConns = 10
	pgConfig.MaxConnLifetime = time.Minute * 10

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		logger.Log.Fatalw("failed to open pool connections", zap.Error(err))
		return nil
	}

	ctxForPing, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err = pool.Ping(ctxForPing)
	if err != nil {
		return nil
	}

	return pool
}
