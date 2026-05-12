package mock

import (
	"auth_service/internal/repository"
	"context"
	"time"
)

type mockDb struct{}

var (
	email        string = "jonsnow@gmail.com"
	refreshtoken string = "567033f2-9283-4250-b4a4-ea965524b798"
	uuid         string = "875542d9-1da4-456c-9bd7-4087f30f62a7"
	userId string = "0d7da674-66c2-431d-a364-2cc82d480780"
	name string = "jon snow"
	RoleId string = "4524f5eb-4ce2-458c-b3ac-0bd32781718b"
	password string = "king_in_the_north"
	timeDuration time.Duration = time.Second
	userCreatedAt time.Time = time.Now()
)

func NewMockDb() repository.RepoIface {
	return &mockDb{}
}

func (m *mockDb) CreateUser(ctx context.Context, name string, email string, hashedPass string, role string) (string, time.Time, error) {
	return userId, userCreatedAt, nil
}

func (m *mockDb) GetEmail(ctx context.Context, UserId string) (string, error) {
	return email, nil
}

func (m *mockDb) ChangeEmail(ctx context.Context, UserId string, newEmail string) error {
	return nil
}

func (m *mockDb) GetRefreshToken(ctx context.Context, uuid string) (string, error) {
	return refreshtoken, nil
}

func (m *mockDb) SetRefreshTokenAndUUID(ctx context.Context, uuid string, Refreshtoken string, UserId string) error {
	return nil
}

func (m *mockDb) SetUUID(ctx context.Context, uuid string, UserId string) error {
	return nil
}

func (m *mockDb) GetUUID(ctx context.Context, UserId string) (string, error) {
	return uuid, nil
}

func (m *mockDb) RefreshExists(ctx context.Context, UserId string) (bool, error) {
	return true, nil
}
func (m *mockDb) FetchUserIdRoleIdName(ctx context.Context, email string) (string, string, string, error) {
	return userId, RoleId, name, nil
}

func (m *mockDb) UserLogout(ctx context.Context, UserId string, uuid string) error     {
	return nil
}
func (m *mockDb) FetchUserPass(ctx context.Context, email string) (string, error)      {
	return password, nil
}
func (m *mockDb) ChangePass(ctx context.Context, UserId string, newPass string) error  {
	return nil
}
func (m *mockDb) FetchUserPassById(ctx context.Context, UserId string) (string, error) {
	return password, nil
}
func (m *mockDb) GetPostgresRespTime(ctx context.Context) *time.Duration               {
	return &timeDuration
}
func (m *mockDb) GetRedisRespTime(ctx context.Context) *time.Duration                  {
	return &timeDuration
}
