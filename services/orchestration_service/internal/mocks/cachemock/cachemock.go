package cachemock

import (
	"context"
	"fmt"
	"orchestration_service/internal/types"
)

type CacheMock struct {
	Down bool
	Skip bool
	Data []string
}

func (c *CacheMock) SetTaskTimeOut(ctx context.Context, data types.Data) error {
	if c.Down {
		return fmt.Errorf("cache is down")
	}

	c.Data[0] = data.TaskId

	return nil
}

func (c *CacheMock) SkipTask(ctx context.Context, data types.Data) (bool, error) {
	if c.Down {
		return false, fmt.Errorf("cache is down")
	}

	if c.Skip {
		return true, nil
	}

	return false, nil

}
