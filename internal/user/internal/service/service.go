package service

import (
	"context"

	"github.com/97Dmitry/go-todolist/internal/user/internal/domain"
)

type Repository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}
