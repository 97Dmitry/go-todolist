package repository

import (
	"uuid"

	"github.com/97Dmitry/go-todolist/internal/user/internal/domain"
)

type userModel struct {
	UserID  uuid.UUID `db:"user_id"`
	Version int64     `db:"version"`
	Email   string    `db:"email"`
	Name    string    `db:"name"`
}

func (m userModel) toDomain() domain.User {
	return domain.RehydrateUser(m.UserID, m.Version, m.Email, m.Name)
}
