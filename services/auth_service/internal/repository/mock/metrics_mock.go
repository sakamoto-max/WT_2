package mock

import (
	"context"
	"time"
)

var ()

type MetricsMock struct {
	HasError bool
}

func (m *MetricsMock) GetRespTime(ctx context.Context) *time.Duration {

	if m.HasError {
		return nil
	}

	timeDuration := time.Millisecond * 368

	return &timeDuration
}
