package mock

import (
	"context"
	"fmt"
)

var (
	// password string = "king_in_the_north"
	HashedPass string = "$2a$10$RcxRiWy8HE7YBfVCCZUR8.wtf/Q3AFgVgSmXnI8mMBzs0.H83H9.."
)

type PasswordMock struct {
	// password  string
	HasErr bool
}

func (p *PasswordMock) FetchUserPass(ctx context.Context, email string) (string, error) {
	if p.HasErr {
		return "", fmt.Errorf("some error occured")
	}
	return HashedPass, nil
}

func (p *PasswordMock) FetchUserPassById(ctx context.Context, userId string) (string, error) {
	if p.HasErr {
		return "", fmt.Errorf("some error occured")
	}
	return HashedPass, nil
}

func (p *PasswordMock) ChangePass(ctx context.Context, userId string, newPass string) error {
	if p.HasErr {
		return fmt.Errorf("some error occured")
	}

	return nil
}
