package apperrors

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrConflict         = errors.New("conflict")
	ErrInvalidArguments = errors.New("invalid arguments")
)

// Error is an application error with a message that is safe to return to
// API clients. errors.Is matches it against its kind, one of the Err* values.
type Error struct {
	kind    error
	message string
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Unwrap() error {
	return e.kind
}

func (e *Error) Message() string {
	return e.message
}

func NotFound(message string) error {
	return &Error{kind: ErrNotFound, message: message}
}

func Conflict(message string) error {
	return &Error{kind: ErrConflict, message: message}
}

func InvalidArguments(message string) error {
	return &Error{kind: ErrInvalidArguments, message: message}
}
