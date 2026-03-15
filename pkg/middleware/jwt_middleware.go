package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	myerrors "wt/pkg/my_errors"
	"wt/pkg/shared"
	"wt/pkg/utils"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the access token from header
		token := r.Header.Get("access-token")
		if token == "" {
			err := myerrors.NewAppErr(myerrors.ErrTokenIsMissing, http.StatusBadRequest)
			err.AppErrWriter(w)
			return 
		}
		t := shared.JwtToken{}
		claims, err := t.ValidateToken(token)
		if err != nil {
			switch{
			case errors.Is(err, myerrors.ErrTokenMalformed):
				err := myerrors.NewAppErr(myerrors.ErrTokenMalformed, http.StatusBadRequest)
				err.AppErrWriter(w)
			case errors.Is(err, myerrors.ErrTokenInvalid):
				fmt.Printf("error : %v\n", err)
				err := myerrors.NewAppErr(myerrors.ErrTokenInvalid, http.StatusBadRequest)
				err.AppErrWriter(w)
			case errors.Is(err, myerrors.ErrTokenExpired):
				err := myerrors.NewAppErr(myerrors.ErrTokenExpired, http.StatusBadRequest)
				err.AppErrWriter(w)
			default:
				utils.InternalServerErr(w, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), shared.Claimskey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}