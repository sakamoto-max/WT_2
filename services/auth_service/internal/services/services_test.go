package services

// import (
// 	"auth_service/internal/repository"
// 	"context"
// 	"testing"
// )

// func TestService_Signup(t *testing.T) {
// 	mockDb := repository.NewMockDB()

// 	s := withMockDb(t, mockDb)
// 	tests := []struct{
// 		name string
// 		email string
// 		password string
// 		role string
// 		WantErr bool
// 	}{
// 		{"jon snow", "jonsnow@gmail.com", "password123", "user", false},
// 		{"jon snow 2", "jonsnow2@gmail.com", "password123", "user", true},
// 		{"jon snow 3", "jonsnow3@gmail.com", "password123", "user", true},	
// 	}

// 	// s.SignUp(context.TODO(), t)
// 	// for _, tc := range tests {
// 	// 	userId, time, err := s.SignUp(context.TODO(), tc.name, tc.email, tc.password, tc.role)
// 	// }

// }

// func withMockDb(t *testing.T, r repository.RepoIface) *Service {
// 	s := NewService(r)
// 	return s
// }

