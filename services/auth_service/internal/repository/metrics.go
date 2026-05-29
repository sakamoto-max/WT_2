package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type metricsDb struct {
	pg *pgxpool.Pool
}

type MetricsDbIface interface {
	GetRespTime(ctx context.Context) *time.Duration
}

func NewMetricsDb(pg *pgxpool.Pool) MetricsDbIface {
	return &metricsDb{pg: pg}
}

func (m *metricsDb) GetRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := m.pg.Ping(ctx)
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
