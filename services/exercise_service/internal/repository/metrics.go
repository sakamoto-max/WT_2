package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type metricsDb struct {
	pg *pgxpool.Pool
}

func (d *metricsDb) GetRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()
	err := d.pg.Ping(ctx)
	if err != nil {
		return nil
	}

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
