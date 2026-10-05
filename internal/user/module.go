// Package user is the users module.
//
// It is the only package of the module other packages may import:
// everything under internal/ is private to the module and is wired here.
package user

import (
	"github.com/97Dmitry/go-todolist/internal/platform/httpserver"
	"github.com/97Dmitry/go-todolist/internal/platform/postgres"
	"github.com/97Dmitry/go-todolist/internal/user/internal/handler"
	"github.com/97Dmitry/go-todolist/internal/user/internal/repository"
	"github.com/97Dmitry/go-todolist/internal/user/internal/service"
)

type Module struct {
	handler *handler.Handler
}

func New(db postgres.DB) *Module {
	userRepository := repository.New(db)
	userService := service.New(userRepository)

	return &Module{
		handler: handler.New(userService),
	}
}

func (m *Module) Routes() []httpserver.Route {
	return m.handler.Routes()
}
