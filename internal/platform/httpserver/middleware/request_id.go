package middleware

import (
	"context"
	"net/http"
	"uuid"
)

const requestIDHeader = "X-Request-Id"

type requestIDKey struct{}

// RequestID takes the request ID from X-Request-Id if it is a valid UUID,
// otherwise generates a new one. The ID is returned in the response in
// canonical form and stored in the request context.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID, err := uuid.Parse(r.Header.Get(requestIDHeader))
			if err != nil {
				requestID = uuid.NewV7()
			}

			w.Header().Set(requestIDHeader, requestID.String())

			ctx := context.WithValue(r.Context(), requestIDKey{}, requestID.String())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey{}).(string)
	return requestID
}
