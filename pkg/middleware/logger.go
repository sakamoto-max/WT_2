package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sakamoto-max/wt_2-pkg/logger"
)

var loggerContextKey contextKey = "my_logger"

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		logger := logger.NewLogger()
		

		ctx := context.WithValue(r.Context(), loggerContextKey, logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetLogger(ctx context.Context) (*logger.MyLogger, error) {
	logger, ok := ctx.Value(loggerContextKey).(*logger.MyLogger)
	if !ok {
		return nil, fmt.Errorf("error getting the logger")
	}

	return logger, nil
}
