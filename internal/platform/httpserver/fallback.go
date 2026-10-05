package httpserver

import (
	"maps"
	"net/http"

	"github.com/97Dmitry/go-todolist/internal/platform/httpserver/response"
	"github.com/97Dmitry/go-todolist/internal/platform/logger"
)

// serveMux serves r with mux, but answers unmatched requests with problem
// details instead of ServeMux's plain-text 404 and 405 responses.
func serveMux(mux *http.ServeMux, w http.ResponseWriter, r *http.Request) {
	fallback, pattern := mux.Handler(r)
	if pattern != "" {
		// Dispatch through the mux so that r.Pattern and path values are set.
		mux.ServeHTTP(w, r)
		return
	}

	// ServeMux has no hooks for its fallback handlers, so run the fallback
	// against a recorder to learn whether it is a 404, a 405 or a redirect.
	recorder := &headerRecorder{header: make(http.Header)}
	fallback.ServeHTTP(recorder, r)

	if recorder.statusCode >= http.StatusMultipleChoices && recorder.statusCode < http.StatusBadRequest {
		maps.Copy(w.Header(), recorder.header)
		w.WriteHeader(recorder.statusCode)
		return
	}

	statusCode := http.StatusNotFound
	if recorder.statusCode == http.StatusMethodNotAllowed {
		statusCode = http.StatusMethodNotAllowed
		w.Header().Set("Allow", recorder.header.Get("Allow"))
	}

	response.NewResponder(logger.FromContext(r.Context()), w).Problem(response.NewProblem(statusCode, ""))
}

type headerRecorder struct {
	header     http.Header
	statusCode int
}

func (h *headerRecorder) Header() http.Header {
	return h.header
}

func (h *headerRecorder) Write(b []byte) (int, error) {
	return len(b), nil
}

func (h *headerRecorder) WriteHeader(statusCode int) {
	h.statusCode = statusCode
}
