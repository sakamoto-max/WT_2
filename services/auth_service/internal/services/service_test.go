package services

// import (
	// "auth_service/internal/repository"
	// "auth_service/internal/repository/mock"
	// "context"
// 	// "testing"
// 	// "github.com/stretchr/testify/assert"
// )


// var (
// 	name string = "jon snow"
// 	email string = "jonsnow@gmail.com"
// 	password string = "king in the north"
// 	roleNill string = ""
// )



// func TestSignUp(t *testing.T) {

// 	type testCase struct {
// 		name string
// 		email string
// 		passowrd string
// 		role string
// 		wantErr bool 
// 	}

// 	tests := []testCase{
// 		{"likith", "likith1@gmail.com", "likithlikith", "user", false},
// 		{"likith", "likith2@gmail.com", "likithlikith", "user", false},
// 		{"likith", "likith3@gmail.com", "likithlikith", "user", false},
// 		{"likith", "likith4@gmail.com", "likithlikith", "user", false},
// 		{"likith", "likith5@gmail.com", "likithlikith", "user", false},
// 	}



// 	mockDb := repository.NewMockDB()

// 	service := NewService(mockDb)

// 	userId, createdAt, err := service.SignUp(context.Background(), name, email, password, roleNill)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, userId)
// 	assert.NotEmpty(t, createdAt)
// }













