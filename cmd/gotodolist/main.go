package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/97Dmitry/go-todolist/internal/app"
	"github.com/97Dmitry/go-todolist/internal/platform/config"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "gotodolist: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return app.Run(ctx, cfg)
}
