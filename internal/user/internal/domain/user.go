package domain

import (
	"strings"
	"unicode/utf8"
	"uuid"

	"github.com/go-playground/validator/v10"
)

var userValidator = validator.New()

type User struct {
	ID      uuid.UUID
	Version int64
	Email   string
	Name    string
}

// NewUser creates a valid user that is not persisted yet.
func NewUser(email, name string) (User, error) {
	user := User{
		ID:      UninitializedID,
		Version: UninitializedVersion,
		Email:   email,
		Name:    name,
	}

	if err := user.validate(); err != nil {
		return User{}, err
	}

	return user, nil
}

// RehydrateUser restores a user from storage. Stored data is trusted and not
// validated again.
func RehydrateUser(id uuid.UUID, version int64, email, name string) User {
	return User{
		ID:      id,
		Version: version,
		Email:   email,
		Name:    name,
	}
}

func (u User) validate() error {
	if strings.TrimSpace(u.Name) == "" {
		return ErrNameRequired
	}

	if length := utf8.RuneCountInString(u.Name); length < 3 || length > 40 {
		return ErrNameLength
	}

	if err := userValidator.Var(u.Email, "required,email"); err != nil {
		return ErrEmailInvalid
	}

	return nil
}
