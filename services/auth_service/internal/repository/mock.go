package repository

import (
	"context"
	"time"
)

var (
	dummyEmail string = "jonsnow@gmail.com"
	dummyRefreshToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiJhOWE1OGQ1OS02M2JjLTQ0ODgtYjFjNS0xNmNjNDM4YjEwMTciLCJSb2xlSWQiOiI0YmU1ODVhYS1mZjFhLTQ2NzAtOWE5Mi0xNmVlZjljZGFlMjYiLCJpc3MiOiJ3b3Jrb3V0LXRyYWNrZXIiLCJleHAiOjE3NzcyNzEyMjd9.P32_AubINxG4_evl168OTdZLrEI23uK8D1GeGhX3e3E"
	dummyUUID string = "a9a58d59-63bc-4488-b1c5-16cc438b1017"
	dummyUserId string = "11c62220-2f86-4e73-916b-9cc579cfdab2"
	dummyRoleId string = "11c62220-2f86-4e73-916b-9cc579cfbad2"
	dummyName string = "jon snow"
	dummyPass string = "danerys"
)

type mockDb struct {
}

func NewMockDB() RepoIface {

	a := mockDb{}

	return a
}

func (m mockDb) GetEmail(ctx context.Context, userId string) (string, error) {
	return dummyEmail, nil
}

func (m mockDb) ChangeEmail(ctx context.Context, userId string, newEmail string) error {
	return nil
}
func (m mockDb) GetRefreshToken(ctx context.Context, uuid string) (string, error) {
	return dummyRefreshToken, nil
}
func (m mockDb) SetRefreshTokenAndUUID(ctx context.Context, uuid string, Refreshtoken string, userId string) error {
	return nil
}
func (m mockDb) SetUUID(ctx context.Context, uuid string, userId string) error {
	return nil
}
func (m mockDb) GetUUID(ctx context.Context, userId string) (string, error) {
	return dummyUUID, nil
}
func (m mockDb) RefreshExists(ctx context.Context, userId string) (bool, error) {
	return true, nil
}
func (m mockDb) CreateUser(ctx context.Context, name string, email string, hashedPass string, role string) (string, time.Time, error) {
	return dummyUserId, time.Now(), nil
}
func (m mockDb) FetchUserIdRoleIdName(ctx context.Context, email string) (string, string, string, error) {
	return dummyUserId, dummyRoleId, dummyName, nil
}
func (m mockDb) UserLogout(ctx context.Context, userId string, uuid string) error {
	return nil
}
func (m mockDb) FetchUserPass(ctx context.Context, email string) (string, error) {
	return dummyPass, nil
}
func (m mockDb) ChangePass(ctx context.Context, userId string, newPass string) error {
	return nil
}
func (m mockDb) FetchUserPassById(ctx context.Context, userId string) (string, error) {
	return dummyPass, nil
}
func (m mockDb) GetPostgresRespTime(ctx context.Context) *time.Duration {

	timeStart := time.Now()

	time.Sleep(time.Millisecond * 10)

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
func (m mockDb) GetRedisRespTime(ctx context.Context) *time.Duration {
	timeStart := time.Now()

	time.Sleep(time.Millisecond * 10)

	timeEnd := time.Since(timeStart)

	return &timeEnd
}
