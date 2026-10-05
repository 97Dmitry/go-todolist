package app

import (
	"context"
	"fmt"

	"github.com/97Dmitry/go-todolist/internal/platform/config"
	"github.com/97Dmitry/go-todolist/internal/platform/httpserver"
	"github.com/97Dmitry/go-todolist/internal/platform/httpserver/middleware"
	"github.com/97Dmitry/go-todolist/internal/platform/logger"
	"github.com/97Dmitry/go-todolist/internal/platform/postgres"
	"github.com/97Dmitry/go-todolist/internal/user"
)

func Run(ctx context.Context, cfg config.Config) error {
	log, err := logger.New(cfg.Logger)
	if err != nil {
		return fmt.Errorf("init application logger: %w", err)
	}
	defer func() {
		// Sync on stdout returns EINVAL on some platforms, so its error is ignored.
		_ = log.Sync()
	}()

	log.Debug("initializing Postgres connection pool")
	pool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("init Postgres connection pool: %w", err)
	}
	defer pool.Close()

	userModule := user.New(pool)

	apiV1Router := httpserver.NewAPIVersionRouter(httpserver.APIVersionV1)
	apiV1Router.RegisterRoutes(userModule.Routes()...)

	server := httpserver.New(
		cfg.HTTP,
		log,
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.AccessLog(),
		middleware.Recover(),
		middleware.BodyLimit(cfg.HTTP.MaxBodyBytes),
	)
	server.RegisterAPIRouters(apiV1Router)

	log.Debug("starting ToDo application")
	if err := server.Run(ctx); err != nil {
		return fmt.Errorf("run HTTP server: %w", err)
	}

	return nil
}
