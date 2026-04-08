package middleware

import (
	"context"
	"net/http"
	"wt/pkg/logger"
)

var loggerContextKey contextKey = "my_logger"

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		logger := logger.NewLogger()
		

		ctx := context.WithValue(r.Context(), loggerContextKey, logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetLogger(ctx context.Context) *logger.MyLogger {
	logger := ctx.Value(loggerContextKey).(*logger.MyLogger)

	return logger
}
