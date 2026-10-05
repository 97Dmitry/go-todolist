package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/97Dmitry/go-todolist/internal/apperrors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Responder struct {
	log *zap.Logger
	w   http.ResponseWriter
}

func NewResponder(log *zap.Logger, w http.ResponseWriter) *Responder {
	return &Responder{
		log: log,
		w:   w,
	}
}

func (r *Responder) JSON(statusCode int, body any) {
	r.write(statusCode, "application/json", body)
}

func (r *Responder) Problem(problem Problem) {
	r.write(problem.Status, problemContentType, problem)
}

// Error logs err with msg and responds with problem details. Only messages of
// apperrors.Error reach the client; internal error text stays in the logs.
func (r *Responder) Error(err error, msg string) {
	statusCode, level := classify(err)

	r.log.Log(level, msg, zap.Error(err))

	var detail string
	if statusCode < http.StatusInternalServerError {
		detail = publicMessage(err)
	}

	r.Problem(NewProblem(statusCode, detail))
}

func (r *Responder) write(statusCode int, contentType string, body any) {
	r.w.Header().Set("Content-Type", contentType)
	r.w.WriteHeader(statusCode)

	if err := json.NewEncoder(r.w).Encode(body); err != nil {
		r.log.Error("encode response body", zap.Error(err))
	}
}

func classify(err error) (int, zapcore.Level) {
	switch {
	case isMaxBytesError(err):
		return http.StatusRequestEntityTooLarge, zapcore.InfoLevel
	case errors.Is(err, apperrors.ErrInvalidArguments):
		return http.StatusBadRequest, zapcore.InfoLevel
	case errors.Is(err, apperrors.ErrNotFound):
		return http.StatusNotFound, zapcore.DebugLevel
	case errors.Is(err, apperrors.ErrConflict):
		return http.StatusConflict, zapcore.InfoLevel
	default:
		return http.StatusInternalServerError, zapcore.ErrorLevel
	}
}

func publicMessage(err error) string {
	if appErr, ok := errors.AsType[*apperrors.Error](err); ok {
		return appErr.Message()
	}

	if isMaxBytesError(err) {
		return "request body is too large"
	}

	return ""
}

func isMaxBytesError(err error) bool {
	_, ok := errors.AsType[*http.MaxBytesError](err)
	return ok
}
