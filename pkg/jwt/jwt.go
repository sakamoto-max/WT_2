package jwt

import (
	"context"
	"errors"
	myerrors "wt/pkg/my_errors"
	"fmt"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type ClaimsContext string

var Claimskey ClaimsContext

type JwtToken struct {
	claims jwtClaims
}

type jwtClaims struct {
	UserId string
	RoleId string
	jwt.RegisteredClaims
}

type Token interface {
	GenerateAccessToken(userId int, roleId int) (string, error)
	GenerateRefreshToken(userId int, roleId int) (string, error)
	ValidateToken(myToken string) (*jwtClaims, error)
}

func (j *JwtToken) GenerateAccessToken(userId string, roleId string) (string, error) {

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
		return accessToken, myerrors.InternalServerErrMaker(fmt.Errorf("error signing the refresh token %w\n", err))
	}

	return accessToken, nil
}

func (j *JwtToken) GenerateRefreshToken(userId string, roleId string) (string, error) {

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
		return refreshToken, myerrors.InternalServerErrMaker(fmt.Errorf("error signing the refresh token %w\n", err))
	}

	return refreshToken, nil
}

func (j *JwtToken) ValidateToken(myToken string) (*jwtClaims, error) {

	claims := &jwtClaims{}
	
	token, err := jwt.ParseWithClaims(myToken, claims, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	
	if err != nil {
		switch{
		case errors.Is(err, jwt.ErrTokenMalformed):
			return claims, ErrTokenMalformed
		case errors.Is(err, jwt.ErrTokenExpired):
			return claims, ErrTokenExpired
		case errors.Is(err, jwt.ErrSignatureInvalid):
			return claims, ErrSignatureInvalid
		}
		return claims, err
	}
	
	if !token.Valid {
		return claims, ErrTokenInvalid
	}
	
	return claims, nil
}

func (j *JwtToken) GetClaimsFromContext(ctx context.Context) (*jwtClaims, bool) {
	claims, ok := ctx.Value(Claimskey).(*jwtClaims)
	return claims, ok
}

var (
	ErrTokenExpired     = errors.New("token is expired, get a new access token at /refresh")
	// ErrUserExits2       = status.Error(codes.AlreadyExists, "user already exits")
	ErrTokenMalformed   = errors.New("token is malformed. please check the token again")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrTokenIsMissing   = errors.New("token is missing, please provide the token")
	ErrRefreshExpired   = errors.New("referesh token is expired, please login again")
	ErrSignatureInvalid = errors.New("token's signature is invalid")
	ErrOldPassNewPassSame = errors.New("the old pass and new pass cannot be the same")
)
