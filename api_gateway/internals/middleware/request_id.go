package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

var(
	ReqIdContextKey contextKey = "req_id"
)


func ReqIdGenerator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := uuid.NewString()

		ctx := context.WithValue(r.Context(), ReqIdContextKey, reqId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetReqId(ctx context.Context) (string, error) {

	reqId, ok := ctx.Value(ReqIdContextKey).(string)
	if !ok {
		return "", fmt.Errorf("error getting the req id")
	}
	return reqId, nil

}

// reqID middleware -> generate a uuid for the request
// logger middleware -> creates and sends the logger 
// jwt middleware -> validates the token