package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/97Dmitry/go-todolist/internal/platform/config"
	"github.com/97Dmitry/go-todolist/internal/platform/httpserver/middleware"
	"go.uber.org/zap"
)

type Server struct {
	mux    *http.ServeMux
	config config.HTTP
	log    *zap.Logger

	middleware []middleware.Middleware
}

func New(
	cfg config.HTTP,
	log *zap.Logger,
	middlewares ...middleware.Middleware,
) *Server {
	return &Server{
		mux:        http.NewServeMux(),
		config:     cfg,
		log:        log,
		middleware: middlewares,
	}
}

func (s *Server) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}
}

// Run serves HTTP until ctx is cancelled, then shuts down gracefully.
func (s *Server) Run(ctx context.Context) error {
	errorLog, err := zap.NewStdLogAt(s.log, zap.WarnLevel)
	if err != nil {
		return fmt.Errorf("create HTTP server error log: %w", err)
	}

	listener, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", s.config.Addr, err)
	}

	root := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveMux(s.mux, w, r)
	})

	server := &http.Server{
		Handler:           middleware.Chain(root, s.middleware...),
		ReadHeaderTimeout: s.config.ReadHeaderTimeout,
		ReadTimeout:       s.config.ReadTimeout,
		WriteTimeout:      s.config.WriteTimeout,
		IdleTimeout:       s.config.IdleTimeout,
		ErrorLog:          errorLog,
	}

	serveErr := make(chan error, 1)

	go func() {
		defer close(serveErr)

		s.log.Info("HTTP server listening", zap.String("addr", listener.Addr().String()))

		if err := server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}
	}()

	select {
	case err := <-serveErr:
		return fmt.Errorf("serve HTTP: %w", err)
	case <-ctx.Done():
	}

	s.log.Info("shutting down HTTP server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		_ = server.Close()

		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	s.log.Info("HTTP server stopped")

	return nil
}
