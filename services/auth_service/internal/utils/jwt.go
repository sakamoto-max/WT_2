package utils

import (
	"auth_service/internal/models"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token interface {
	GenerateAccessToken(userId int, roleId int) (string, error)
	GenerateRefreshToken(userId int, roleId int) (string, error)
	ValidateToken(myToken string) (*jwtClaims, error)
}
type JwtToken struct {
	claims jwtClaims
}
type jwtClaims struct {
	UserId int
	RoleId int
	jwt.RegisteredClaims
}

func (j *JwtToken) GenerateAccessToken(userId int, roleId int) (string, error) {

	var accessToken string

	j.claims.UserId = userId
	j.claims.RoleId = roleId
	j.claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    "workout-tracker",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j.claims)

	accessToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return accessToken, fmt.Errorf("error signing the refresh token %w\n", err)
	}

	return accessToken, nil
}

func (j *JwtToken) GenerateRefreshToken(userId int, roleId int) (string, error) {

	var refreshToken string

	j.claims.UserId = userId
	j.claims.RoleId = roleId
	j.claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    "workout-tracker",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 15)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j.claims)

	refreshToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return refreshToken, fmt.Errorf("error signing the refresh token %w\n", err)
	}

	return refreshToken, nil
}

func (j *JwtToken) ValidateToken(myToken string) (*jwtClaims, error) {
	claims := &jwtClaims{}

	token, err := jwt.ParseWithClaims(myToken, claims, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return claims, err
	}

	if !token.Valid {
		return claims, fmt.Errorf("token is not valid")
	}

	return claims, nil
}

func (j *JwtToken) GetClaimsFromContext(ctx context.Context) (*jwtClaims, bool) {
	claims, ok := ctx.Value(models.Claimskey).(*jwtClaims)
	return claims, ok
}
