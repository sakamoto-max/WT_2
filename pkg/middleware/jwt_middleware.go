package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	myerrors "github.com/sakamoto-max/wt_2-pkg/my_errors"
	"github.com/sakamoto-max/wt_2-pkg/jwt"
	"github.com/sakamoto-max/wt_2-pkg/utils"
)

var claimsContextKey contextKey = "get_claims"

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("access-token")
		if token == "" {
			err := myerrors.NewAppErr(jwt.ErrTokenIsMissing, http.StatusBadRequest)
			err.AppErrWriter(w)
			return 
		}
		t := jwt.JwtToken{}
		claims, err := t.ValidateToken(token)
		if err != nil {
			switch{
			case errors.Is(err, jwt.ErrTokenMalformed):
				err := myerrors.NewAppErr(jwt.ErrTokenMalformed, http.StatusBadRequest)
				err.AppErrWriter(w)
			case errors.Is(err, jwt.ErrTokenInvalid):
				fmt.Printf("error : %v\n", err)
				err := myerrors.NewAppErr(jwt.ErrTokenInvalid, http.StatusBadRequest)
				err.AppErrWriter(w)
			case errors.Is(err, jwt.ErrTokenExpired):
				err := myerrors.NewAppErr(jwt.ErrTokenExpired, http.StatusBadRequest)
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

