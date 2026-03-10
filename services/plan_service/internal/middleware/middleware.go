package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	UserId int
	RoleId int
	jwt.RegisteredClaims
}

type ClaimsContext string

var Claimskey ClaimsContext

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the access token from header
		token := r.Header.Get("access-token")
		claims, err := validateToken(token)
		if err != nil {
			fmt.Printf("error occured while validating the token : %v\n", err)
			return
		}

		ctx := context.WithValue(r.Context(), Claimskey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func GetClaimsFromContext(ctx context.Context) (*JwtClaims, bool) {
	claims, ok := ctx.Value(Claimskey).(*JwtClaims)

	return claims, ok
}

func validateToken(myToken string) (*JwtClaims, error) {
	claims := &JwtClaims{}

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
