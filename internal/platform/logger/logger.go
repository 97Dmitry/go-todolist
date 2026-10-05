package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/97Dmitry/go-todolist/internal/platform/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg config.Logger) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	switch cfg.Format {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return nil, fmt.Errorf("unknown log format %q: want json or console", cfg.Format)
	}

	core := zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), level)

	return zap.New(core, zap.AddCaller()), nil
}

type contextKey struct{}

var nop = zap.NewNop()

func ToContext(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, log)
}

// FromContext returns the request logger, or a no-op logger when the context
// has none, so code reached outside the middleware chain does not crash.
func FromContext(ctx context.Context) *zap.Logger {
	if log, ok := ctx.Value(contextKey{}).(*zap.Logger); ok {
		return log
	}

	return nop
}
