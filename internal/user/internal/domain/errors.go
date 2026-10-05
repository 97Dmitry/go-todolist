package domain

import "github.com/97Dmitry/go-todolist/internal/apperrors"

var (
	ErrNameRequired = apperrors.InvalidArguments("user name is required")
	ErrNameLength   = apperrors.InvalidArguments("user name must contain between 3 and 40 characters")
	ErrEmailInvalid = apperrors.InvalidArguments("user email is invalid")
	ErrEmailTaken   = apperrors.Conflict("user with this email already exists")
	ErrUserInvalid  = apperrors.InvalidArguments("user data is invalid")
)
