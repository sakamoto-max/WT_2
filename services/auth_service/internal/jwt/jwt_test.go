package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateAccessToken(t *testing.T) {

	accessToken, err := GenerateAccessToken("123", "456")
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
}

func Test_GenerateAccessTokenFastExp(t *testing.T) {
	accessToken, err := GenerateAccessTokenFastExp("123", "456")
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
}

func Test_GenerateRefreshToken(t *testing.T) {
	refreshToken, err := GenerateRefreshToken("123", "456")
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)
}

func Test_ValidateToken(t *testing.T) {

	accessToken, _ := GenerateAccessToken("123", "456")

	claims, err := ValidateToken(accessToken)
	assert.NoError(t, err)

	assert.Equal(t, claims.UserId, "123")
	assert.Equal(t, claims.RoleId, "456")
}
