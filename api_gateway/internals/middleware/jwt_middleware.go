package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sakamoto-max/wt_2/api_gateway/internals/jwt"
	"github.com/sakamoto-max/wt_2/api_gateway/internals/utils"
	myerrors "github.com/sakamoto-max/wt_2_pkg/my_errors"
)

var claimsContextKey contextKey = "get_claims"

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("access-token")
		if token == "" {
			err := myerrors.NewAppErr(jwt.ErrTokenIsMissing, http.StatusUnauthorized)
			err.AppErrWriter(w)
			return
		}
		
		claims, err := jwt.ValidateToken(token)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenMalformed):
				err := myerrors.NewAppErr(jwt.ErrTokenMalformed, http.StatusUnauthorized)
				err.AppErrWriter(w)
			case errors.Is(err, jwt.ErrTokenInvalid):
				err := myerrors.NewAppErr(jwt.ErrTokenInvalid, http.StatusUnauthorized)
				err.AppErrWriter(w)
			case errors.Is(err, jwt.ErrTokenExpired):
				err := myerrors.NewAppErr(jwt.ErrTokenExpired, http.StatusUnauthorized)
				err.AppErrWriter(w)
			case errors.Is(err, jwt.ErrSignatureInvalid):
				err := myerrors.NewAppErr(jwt.ErrSignatureInvalid, http.StatusUnauthorized)
				err.AppErrWriter(w)
			default:
				utils.InternalServerErr(w, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), claimsContextKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func GetClaims(ctx context.Context) (*jwt.JwtClaims, error) {
	claims, ok := ctx.Value(claimsContextKey).(*jwt.JwtClaims)
	if !ok {
		return nil, fmt.Errorf("error getting claims")
	}
	return claims, nil
}
