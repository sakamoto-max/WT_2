package servicesmock

import (
	"auth_service/internal/domain/user"
	"context"
	"time"
)

type ServiceMock struct {
	Err error
	PgDown bool
	RedisDown bool
}

func (s *ServiceMock) SignUp(ctx context.Context, payload user.SignUpPayload) (string, time.Time, error) {
	if s.Err != nil {
		return "", time.Now(), s.Err
	}

	return "123", time.Now(), nil
}
func (s *ServiceMock) Login(ctx context.Context, payload user.LoginPayload) (string, string, string, string, error) {
	if s.Err != nil {
		return "", "", "", "", s.Err
	}

	return "123", "jon snow", "king_in_the_north", "123454323454", nil
}
func (s *ServiceMock) Logout(ctx context.Context, userId string) error {
	if s.Err != nil {
		return s.Err
	}

	return nil
}
func (s *ServiceMock) GetNewAccessTokenSer(ctx context.Context, UUID string) (string, error) {
	if s.Err != nil {
		return "", s.Err
	}

	return "new_access_token", nil
}
func (s *ServiceMock) ChangePass(ctx context.Context, payload user.ChangePassPayload) error {
	if s.Err != nil {
		return s.Err
	}

	return nil
}
func (s *ServiceMock) ChangeEmail(ctx context.Context, payload user.ChangeEmailPayload) error {
	if s.Err != nil {
		return s.Err
	}

	return nil
}
func (s *ServiceMock) GetHealth(ctx context.Context) (*time.Duration, *time.Duration) {
	timeDuration := time.Millisecond * 458
	if s.PgDown {
		return nil, &timeDuration
	}

	if s.RedisDown {
		return &timeDuration, nil
	}

	if s.PgDown && s.RedisDown {
		return nil, nil
	}

	return &timeDuration, &timeDuration
}
