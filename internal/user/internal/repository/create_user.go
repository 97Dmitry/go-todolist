package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/97Dmitry/go-todolist/internal/user/internal/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const usersEmailUniqueIndex = "users_email_unique_idx"

func (r *Repository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	const query = `
	INSERT INTO todolist.users (email, name)
	VALUES ($1, $2)
	RETURNING user_id, version, email, name`

	rows, err := r.db.Query(ctx, query, user.Email, user.Name)
	if err != nil {
		return domain.User{}, insertUserError(err)
	}

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[userModel])
	if err != nil {
		return domain.User{}, insertUserError(err)
	}

	return model.toDomain(), nil
}

func insertUserError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch {
		case pgErr.Code == pgerrcode.UniqueViolation && pgErr.ConstraintName == usersEmailUniqueIndex:
			return domain.ErrEmailTaken
		case pgErr.Code == pgerrcode.CheckViolation:
			return fmt.Errorf("check constraint %s: %w", pgErr.ConstraintName, domain.ErrUserInvalid)
		}
	}

	return fmt.Errorf("insert user: %w", err)
}
