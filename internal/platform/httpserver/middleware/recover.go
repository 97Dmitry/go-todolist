package middleware

import (
	"net/http"

	"github.com/97Dmitry/go-todolist/internal/platform/httpserver/response"
	"github.com/97Dmitry/go-todolist/internal/platform/logger"
	"go.uber.org/zap"
)

// Recover turns a panic in a handler into a 500 response and logs it with
// the stack trace.
func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := response.NewResponseWriter(w)

			defer func() {
				p := recover()
				if p == nil {
					return
				}

				if p == http.ErrAbortHandler {
					panic(p)
				}

				log := logger.FromContext(r.Context())
				log.Error("panic recovered", zap.Any("panic", p), zap.Stack("stack"))

				if rw.WroteHeader() {
					// Part of the response is already sent: abort the connection
					// so the client does not take it for a complete response.
					panic(http.ErrAbortHandler)
				}

				response.NewResponder(log, rw).Problem(response.NewProblem(http.StatusInternalServerError, ""))
			}()

			next.ServeHTTP(rw, r)
		})
	}
}
