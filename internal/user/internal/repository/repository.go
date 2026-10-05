package repository

import "github.com/97Dmitry/go-todolist/internal/platform/postgres"

type Repository struct {
	db postgres.DB
}

func New(db postgres.DB) *Repository {
	return &Repository{
		db: db,
	}
}
