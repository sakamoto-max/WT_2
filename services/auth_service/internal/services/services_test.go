package services

import (
	// "auth_service/internal/domain/user"

	"auth_service/internal/repository"
	"auth_service/internal/repository/cache"
	"auth_service/internal/repository/cache/cachemock"
	"auth_service/internal/repository/mock"
	"context"
	"testing"

	"github.com/sakamoto-max/wt_2_proto/shared/auth"
	"github.com/stretchr/testify/assert"
)

func Test_SignUp(t *testing.T) {

	tests := []struct {
		name          string
		emailExits    bool
		userNameExits bool
		wantErr       bool
	}{
		{name: "straight through", emailExits: false, userNameExits: false, wantErr: false},
		{name: "email exits", emailExits: true, userNameExits: false, wantErr: true},
		{name: "user name exits", emailExits: false, userNameExits: true, wantErr: true},
		{name: "both email and user names exits", emailExits: true, userNameExits: true, wantErr: true},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			s := Service{
				db: repository.Db{
					Auth: &mock.AuthMock{
						EmailExits:     test.emailExits,
						UserNameExists: test.userNameExits,
					},
				},
			}

			resp, err := s.UserSignUp(context.Background(),
				&auth.UserSignUpReq{
					Name:     "jon snow",
					Email:    "jonsnow@gmail.com",
					Password: "king_in_the_north",
					Role:     "user",
				},
			)
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, resp.CreatedAt)
			assert.NotZero(t, resp.Email)
			assert.NotZero(t, resp.Name)
			assert.NotZero(t, resp.Role)
			assert.NotZero(t, resp.UserId)
		})

	}
}
func Test_Login(t *testing.T) {

	tests := []struct {
		email             string
		password          string
		name              string
		passWordHasErr    bool
		userDetailsHasErr bool
		TokenMockHasErr   bool
		TokenMockHit      bool
		TokenMockMiss     bool
		uuidMockHasErr    bool
		uuidMockHit       bool
		uuidMockMiss      bool
		wantErr           bool
	}{
		{
			name:              "success",
			email:             "jonsnow@gmail.com",
			password:          "king_in_the_north",
			wantErr:           false,
			passWordHasErr:    false,
			userDetailsHasErr: false,
			TokenMockHasErr:   false,
			TokenMockHit:      true,
			TokenMockMiss:     false,
			uuidMockHasErr:    false,
			uuidMockHit:       true,
			uuidMockMiss:      false,
		},
		{
			name:           "not signed up",
			email:          "jonsnow@gmail.com",
			password:       "king_in_the_north",
			wantErr:        true,
			passWordHasErr: true,
		},
		{
			name:           "incorrect passoword",
			email:          "jonsnow@gmail.com",
			password:       "king_in_the_no",
			wantErr:        true,
			passWordHasErr: false,
		},
		{
			name:              "not first time logging in",
			email:             "jonsnow@gmail.com",
			password:          "king_in_the_north",
			wantErr:           false,
			passWordHasErr:    false,
			userDetailsHasErr: false,
			TokenMockHasErr:   false,
			TokenMockHit:      true,
			TokenMockMiss:     false,
			uuidMockHasErr:    false,
			uuidMockHit:       true,
			uuidMockMiss:      false,
		},
		{
			name:              "first time log in",
			email:             "jonsnow@gmail.com",
			password:          "king_in_the_north",
			wantErr:           false,
			passWordHasErr:    false,
			userDetailsHasErr: false,
			TokenMockHasErr:   false,
			TokenMockHit:      false,
			TokenMockMiss:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s := Service{
				db: repository.Db{
					Password: &mock.PasswordMock{
						HasErr: test.passWordHasErr,
					},
					UserDeatails: &mock.UserMock{
						HasErr: test.userDetailsHasErr,
					},
				},
				cache: &cache.Cache{
					Token: &cachemock.TokenMock{Hit: test.TokenMockHit, Miss: test.TokenMockMiss, HasErr: test.TokenMockHasErr},
					Uuid:  &cachemock.UuidMock{Hit: test.uuidMockHit, Miss: test.uuidMockMiss, HasErr: test.uuidMockHasErr},
				},
			}

			resp, err := s.UserLogin(context.Background(),
				&auth.UserLoginReq{
					Email:    test.email,
					Password: test.password,
				})

			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, resp.UserId)
			assert.NotZero(t, resp.Name)
			assert.NotZero(t, resp.AccessToken)
			assert.NotZero(t, resp.UUID)
			assert.NotZero(t, resp.Email)
			assert.NotZero(t, resp.Message)
		})
	}
}

func Test_GetHealth(t *testing.T) {

	tests := []struct {
		name     string
		pgErr    bool
		redisErr bool
		pgNil    bool
		redisNil bool
		wantErr  bool
	}{
		{name: "both working well", pgErr: false, redisErr: false, pgNil: false, redisNil: false},
		{name: "pg is working well and redis is down", pgErr: false, redisErr: true, pgNil: false, redisNil: true},
		{name: "redis is working well and pg is down", pgErr: true, redisErr: false, pgNil: true, redisNil: false},
		{name: "both are down", pgErr: true, redisErr: true, pgNil: true, redisNil: true},
	}

	for _, test := range tests {

		s := Service{
			db: repository.Db{
				Metrics: &mock.MetricsMock{
					HasError: test.pgErr,
				},
			},
			cache: &cache.Cache{
				Metrics: &cachemock.MetricsMock{HasError: test.redisErr},
			},
		}

		resp, err := s.GetHealth(context.Background(), &auth.GetHealthReq{})
		if test.wantErr {
			assert.Error(t, err)
			return
		}

		if test.pgNil && test.redisNil {
			assert.Nil(t, resp.PostgresRespTime)
			assert.Nil(t, resp.RedisRespTime)
			return
		}

		if test.redisNil {
			assert.Nil(t, resp.RedisRespTime)
			assert.NotNil(t, resp.PostgresRespTime)
			return
		}

		if test.pgNil {
			assert.Nil(t, resp.PostgresRespTime)
			assert.NotNil(t, resp.RedisRespTime)
			return
		}

		assert.NotNil(t, resp.PostgresRespTime)
		assert.NotNil(t, resp.RedisRespTime)

	}

}
func Test_LogOut(t *testing.T) {

	tests := []struct {
		name           string
		uuidHit        bool
		uuidMiss       bool
		uuidHasErr     bool
		userMockHasErr bool
		wantErr        bool
	}{
		{name: "success", uuidHit: true},
		{name: "uuid miss", uuidMiss: true, wantErr: true},
		{name: "redis went down", uuidHasErr: true, userMockHasErr: true, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s := Service{
				cache: &cache.Cache{
					Uuid: &cachemock.UuidMock{
						Hit:    test.uuidHit,
						Miss:   test.uuidMiss,
						HasErr: test.uuidHasErr,
					},
					User: &cachemock.UserMock{
						HasErr: test.userMockHasErr,
					},
				},
			}

			resp, err := s.UserLogOut(context.Background(), &auth.SendUserId{UserId: "123"})
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, resp.Message)
		})
	}
}

func Test_GetNewAccessToken(t *testing.T) {

	tests := []struct {
		name         string
		hit          bool
		miss         bool
		redisHasErr  bool
		tokenExpired bool
		wantErr      bool
	}{
		{name: "success", hit: true, tokenExpired: false, wantErr: false},
		{name: "redis down", redisHasErr: true, wantErr: true},
		{name: "refresh expired", hit: true, tokenExpired: true, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Service{
				cache: &cache.Cache{
					Token: &cachemock.TokenMock{
						Hit:        test.hit,
						HasErr:     test.wantErr,
						RefreshExp: test.tokenExpired,
					},
				},
			}

			resp, err := s.GetNewAccessToken(context.Background(), &auth.SendUUID{UUID: "123"})

			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, resp.NewAccessToken)

		})
	}
}

func Test_ChangePass(t *testing.T) {

	tests := []struct {
		name               string
		oldPass            string
		newPass            string
		passWordMockHasErr bool
		wantErr            bool
	}{
		{name: "success", oldPass: "king_in_the_north", newPass: "password123", passWordMockHasErr: false, wantErr: false},
		{name: "old password same as the new password", oldPass: "king_in_the_north", newPass: "king_in_the_north", passWordMockHasErr: false, wantErr: true},
		{name: "incorrect password", oldPass: "king_of_the_seven_kingdoms", newPass: "password123", passWordMockHasErr: true, wantErr: true},
		{name: "db down", oldPass: "king_in_the_north", newPass: "password123", passWordMockHasErr: true, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Service{
				db: repository.Db{
					Password: &mock.PasswordMock{
						HasErr: test.passWordMockHasErr,
					},
				},
			}

			resp, err := s.ChangePass(context.Background(), &auth.ChangePassReq{
				UserId:  "1234543212345",
				OldPass: test.oldPass,
				NewPass: test.newPass,
			})

			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, resp.Message)
		})
	}
}

func Test_ChangeEmail(t *testing.T) {

	tests := []struct {
		name     string
		oldEmail string
		NewEmail string
		DbErr    bool
		wantErr  bool
	}{
		{name: "success", oldEmail: "jonsnow@gmail.com", NewEmail: "kinginthenorth@gmail.com", DbErr: false, wantErr: false},
		{name: "success", oldEmail: "jonsnow@gmail.com", NewEmail: "jonsnow@gmail.com", DbErr: false, wantErr: true},
		{name: "success", oldEmail: "jonsnow@gmail.com", NewEmail: "kinginthenorth@gmail.com", DbErr: true, wantErr: true},
		{name: "success", oldEmail: "bastardking@gmail.com", NewEmail: "kinginthenorth@gmail.com", DbErr: false, wantErr: true},
	}

	for _, test := range tests {
		s := Service{
			db: repository.Db{
				Email: &mock.EmailMock{
					HasError: test.DbErr,
					OldEmail: test.oldEmail,
				},
			},
		}

		resp, err := s.ChangeEmail(context.Background(), &auth.ChangeEmailReq{
			UserId:   "12345432123456",
			OldEmail: test.oldEmail,
			NewEmail: test.NewEmail,
		})

		if test.wantErr {
			assert.Error(t, err)
			return
		}

		assert.NoError(t, err)
		assert.NotZero(t, resp)
	}
}
