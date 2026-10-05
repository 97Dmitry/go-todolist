package service

import (
	"context"
	"fmt"

	"github.com/97Dmitry/go-todolist/internal/user/internal/domain"
)

func (s *Service) CreateUser(ctx context.Context, email, name string) (domain.User, error) {
	user, err := domain.NewUser(email, name)
	if err != nil {
		return domain.User{}, fmt.Errorf("new user: %w", err)
	}

	createdUser, err := s.repository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("repository create user: %w", err)
	}

	return createdUser, nil
}
