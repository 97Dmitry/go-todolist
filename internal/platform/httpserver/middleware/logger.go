package middleware

import (
	"net/http"

	"github.com/97Dmitry/go-todolist/internal/platform/logger"
	"go.uber.org/zap"
)

// Logger stores a request-scoped logger in the request context. It must run
// after RequestID.
func Logger(log *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestLog := log.With(
				zap.String("request_id", RequestIDFromContext(r.Context())),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)

			ctx := logger.ToContext(r.Context(), requestLog)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
