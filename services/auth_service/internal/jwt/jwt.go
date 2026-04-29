package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	myerrors "github.com/sakamoto-max/wt_2_pkg/myerrs"
)

type jwtStruct struct {
	claims JwtClaims
}

type JwtClaims struct {
	UserId string
	RoleId string
	jwt.RegisteredClaims
}

func GenerateAccessToken(userId string, roleId string) (string, error) {
	j := jwtStruct{}

	j.claims.UserId = userId
	j.claims.RoleId = roleId
	j.claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    "workout-tracker",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j.claims)

	accessToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", myerrors.InternalServerErrMaker(fmt.Errorf("error signing the refresh token %w\n", err))
	}

	return accessToken, nil
}

func GenerateRefreshToken(userId string, roleId string) (string, error) {

	j := jwtStruct{}

	j.claims.UserId = userId
	j.claims.RoleId = roleId
	j.claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    "workout-tracker",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 15)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j.claims)

	refreshToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", myerrors.InternalServerErrMaker(fmt.Errorf("error signing the refresh token %w\n", err))
	}

	return refreshToken, nil
}

func ValidateToken(myToken string) (*JwtClaims, error) {

	claims := &JwtClaims{}

	token, err := jwt.ParseWithClaims(myToken, claims, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
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

var (
	ErrTokenExpired     = errors.New("token is expired, get a new access token at /refresh")
	ErrTokenMalformed   = errors.New("token is malformed. please check the token again")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrRefreshExpired   = errors.New("referesh token is expired, please login again")
	ErrSignatureInvalid = errors.New("token's signature is invalid")
)
