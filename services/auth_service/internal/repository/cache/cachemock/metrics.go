package cachemock

import (
	"context"
	"time"
)

var (
	timeDuration time.Duration = time.Millisecond * 300
)

type MetricsMock struct {
	HasError bool
}

func (m *MetricsMock) GetRespTime(ctx context.Context) *time.Duration {
	if m.HasError{
		return nil
	}

	return &timeDuration
}
