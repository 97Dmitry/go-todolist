package handler

import (
	"context"
	"net/http"

	"github.com/97Dmitry/go-todolist/internal/platform/httpserver"
	"github.com/97Dmitry/go-todolist/internal/user/internal/domain"
)

type Service interface {
	CreateUser(ctx context.Context, email, name string) (domain.User, error)
}

type Handler struct {
	service Service
}

func New(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Routes() []httpserver.Route {
	return []httpserver.Route{
		httpserver.NewRoute(http.MethodPost, "/users", h.CreateUser),
	}
}
