package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTP     HTTP     `envPrefix:"HTTP_"`
	Postgres Postgres `envPrefix:"POSTGRES_"`
	Logger   Logger   `envPrefix:"LOGGER_"`
}

type HTTP struct {
	Addr              string        `env:"ADDR,notEmpty"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT,notEmpty"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"5s"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT" envDefault:"15s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT" envDefault:"15s"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
	MaxBodyBytes      int64         `env:"MAX_BODY_BYTES" envDefault:"1048576"`
}

type Postgres struct {
	Host             string        `env:"HOST,notEmpty"`
	Port             string        `env:"PORT" envDefault:"5432"`
	User             string        `env:"USER,notEmpty"`
	Password         string        `env:"PASSWORD,notEmpty"`
	Database         string        `env:"DB,notEmpty"`
	SSLMode          string        `env:"SSLMODE" envDefault:"prefer"`
	ConnectTimeout   time.Duration `env:"CONNECT_TIMEOUT" envDefault:"5s"`
	StatementTimeout time.Duration `env:"STATEMENT_TIMEOUT" envDefault:"10s"`
}

type Logger struct {
	Level  string `env:"LEVEL" envDefault:"info"`
	Format string `env:"FORMAT" envDefault:"json"`
}

func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("parse env: %w", err)
	}

	return cfg, nil
}
