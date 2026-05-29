package cachemock

import (
	"auth_service/internal/jwt"
	"context"
	"fmt"
)

var (
	// Refreshtoken string = "567033f2-9283-4250-b4a4-ea965524b798"
)

type TokenMock struct {
	Hit    bool
	Miss   bool
	HasErr bool
	RefreshExp bool
}

func (t *TokenMock) RefreshExists(ctx context.Context, userId string) (bool, error) {
	if t.Miss {
		return false, nil
	}

	if t.HasErr {
		return false, fmt.Errorf("some error occured")
	}

	return true, nil
}

func (t *TokenMock) SetRefreshTokenAndUUID(ctx context.Context, uuid string, Refreshtoken string, userId string) error {
	if t.HasErr {
		return fmt.Errorf("some error occured")
	}

	return nil
}
func (t *TokenMock) GetRefreshToken(ctx context.Context, uuid string) (string, error) {
	// send a real refresh token
	// send an expired refresh token
	if t.HasErr {
		return "", fmt.Errorf("some error occured")
	}


	if t.Miss {
		return "", nil
	}

	if t.RefreshExp {
		refresh, _ := jwt.GenerateAccessTokenFastExp("123454321", "123456543")
		return refresh, nil
	}

	refresh, _ := jwt.GenerateRefreshToken("123456789765432", "2134567234567")
	return refresh, nil
}
