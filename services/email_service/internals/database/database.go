package database

import (
	"context"
	"email_service/internals/config"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewPgConn(config config.Config) *pgxpool.Pool {

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		config.Db.UserName,
		config.Db.Pass,
		config.Db.Host,
		config.Db.Port,
		config.Db.DatabaseName,
		config.Db.SSlMode,
	)

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
