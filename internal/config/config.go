package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `env:"ENVIRONMENT"`
	Database    DatabaseConfig
	HTTPServer  HTTPServerConfig
	JWT         JWTConfig
}

type HTTPServerConfig struct {
	Address string        `env:"HTTP_SERVER_ADDRESS"`
	Port    string        `env:"HTTP_SERVER_PORT"`
	Timeout time.Duration `env:"HTTP_SERVER_TIMEOUT"`
}
type DatabaseConfig struct {
	Username string `env:"DB_USER"`
	DBName   string `env:"DB_NAME"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
}

type JWTConfig struct {
	Secret string        `env:"JWT_SECRET"`
	TTL    time.Duration `env:"JWT_TTL"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read config from env: %w", err)
	}

	if cfg.HTTPServer.Timeout == 0 {
		cfg.HTTPServer.Timeout = 15 * time.Second
	}

	return &cfg, nil
}
