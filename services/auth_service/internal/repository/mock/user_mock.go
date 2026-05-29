package mock

import (
	"context"
	"fmt"
)

type UserMock struct {
	HasErr bool
}


func (u *UserMock) FetchUserIdRoleIdName(ctx context.Context, email string) (string, string, string, error) {
	if u.HasErr {
		return "", "", "", fmt.Errorf("some error occured")
	}

	return "123", "345", "name", nil
}
