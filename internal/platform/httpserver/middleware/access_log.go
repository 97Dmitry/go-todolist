package middleware

import (
	"net/http"
	"time"

	"github.com/97Dmitry/go-todolist/internal/platform/httpserver/response"
	"github.com/97Dmitry/go-todolist/internal/platform/logger"
	"go.uber.org/zap"
)

// AccessLog writes one log entry per request. It must run after Logger and
// before Recover, so that requests ending in a panic are logged too.
func AccessLog() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := response.NewResponseWriter(w)
			start := time.Now()

			defer func() {
				logger.FromContext(r.Context()).Info(
					"HTTP request",
					zap.Int("status", rw.StatusCode()),
					zap.Int("bytes", rw.BytesWritten()),
					zap.Duration("duration", time.Since(start)),
				)
			}()

			next.ServeHTTP(rw, r)
		})
	}
}
