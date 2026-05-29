package mock

import (
	"auth_service/internal/mappings"
	"context"
	"errors"
	"time"
)

var (
	ErrSomeErrorOccured = errors.New("some error occred")
	Refreshtoken string        = "567033f2-9283-4250-b4a4-ea965524b798"
	Uuid         string        = "875542d9-1da4-456c-9bd7-4087f30f62a7"
	UserId       string        = "0d7da674-66c2-431d-a364-2cc82d480780"
)

type AuthMock struct {
	EmailExits                 bool
	UserNameExists             bool
}

func (m *AuthMock) CreateUser(ctx context.Context, payload mappings.SignUp) (string, time.Time, error) {

	if m.EmailExits || m.UserNameExists {
		return "", time.Now(), ErrSomeErrorOccured
	}

	return "123", time.Now(), nil
}
