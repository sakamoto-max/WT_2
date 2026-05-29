package mock

import (
	"context"
	"fmt"
)

type EmailMock struct {
	HasError bool
	OldEmail string
}

func (e *EmailMock) GetEmail(ctx context.Context, userId string) (string, error)           {
	if e.HasError {
		return "", fmt.Errorf("some error occured")
	}

	return e.OldEmail, nil
}
func (e *EmailMock) ChangeEmail(ctx context.Context, userId string, newEmail string) error {
	if e.HasError {
		return fmt.Errorf("some error occured")
	}

	return nil
}
