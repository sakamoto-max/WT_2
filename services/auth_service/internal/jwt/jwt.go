package jwt

import (
	"auth_service/internal/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

var SECRETKEY string

var (
	ErrTokenExpired     = errors.New("token is expired, get a new access token at /refresh")
	ErrTokenMalformed   = errors.New("token is malformed. please check the token again")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrTokenIsMissing   = errors.New("token is missing, please provide the token")
	ErrRefreshExpired   = errors.New("referesh token is expired, please login again")
	ErrSignatureInvalid = errors.New("token's signature is invalid")
)

type Claims struct {
	UserId string
	RoleId string
	jwt.RegisteredClaims
}

func GenerateAccessToken(userId string, roleId string) (string, error) {

	j := Claims{}

	j.UserId = userId
	j.RoleId = roleId
	j.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    "workout-tracker",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)

	accessToken, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		return "", myerrors.InternalServerErrMaker(fmt.Errorf("error signing the refresh token %w\n", err))
	}

	return accessToken, nil
}

func GenerateRefreshToken(userId string, roleId string) (string, error) {

	j := Claims{}

	j.UserId = userId
	j.RoleId = roleId
	j.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    "workout-tracker",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 15)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)

	refreshToken, err := token.SignedString([]byte(SECRETKEY))
	if err != nil {
		return "", myerrors.InternalServerErrMaker(fmt.Errorf("error signing the refresh token %w\n", err))
	}

	return refreshToken, nil
}

func ValidateToken(myToken string) (*Claims, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(myToken, claims, func(t *jwt.Token) (any, error) {
		return []byte(SECRETKEY), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return claims, ErrTokenMalformed
		case errors.Is(err, jwt.ErrTokenExpired):
			return claims, ErrTokenExpired
		case errors.Is(err, jwt.ErrSignatureInvalid):
			return claims, ErrSignatureInvalid
		}
		return nil, err
	}

	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

func JwtInit(config config.Config) {
	SECRETKEY = config.Auth.SecretKey
}
