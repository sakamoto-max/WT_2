package dbmock

import (
	"email_service/internals/types"
	"fmt"
)

type DbMock struct {
	Down bool
}

func (d *DbMock) PushToFailed(data types.Data, numberOfTries int, status string, Err error) error {
	if d.Down {
		return fmt.Errorf("db is down")
	}

	return nil
}
