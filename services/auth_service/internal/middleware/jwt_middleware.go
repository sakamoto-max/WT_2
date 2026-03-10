package middleware

import (
	"auth_service/internal/models"
	"auth_service/internal/utils"
	"context"
	"fmt"
	"net/http"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the access token from header
		token := r.Header.Get("access-token")
		t := utils.JwtToken{}
		claims, err := t.ValidateToken(token)
		if err != nil {
			fmt.Printf("error occured while validating the token : %v\n", err)
			return
		}

		ctx := context.WithValue(r.Context(), models.Claimskey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
