package cachemock

import (
	"context"
	"fmt"
)

type UserMock struct {
	HasErr bool
}

func (u *UserMock) UserLogout(ctx context.Context, userId string, uuid string) error {

	if u.HasErr {
		return fmt.Errorf("some error has occured")
	}

	return nil
}
