package config

import (
	"fmt"
	"net"
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
	SSLMode  string `env:"DB_SSLMODE"`
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
	if cfg.JWT.TTL == 0 {
		cfg.JWT.TTL = 30 * time.Minute
	}
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "test_secret"
	}
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}

	return &cfg, nil
}

func (c *Config) HTTPAddr() string {
	return net.JoinHostPort(c.HTTPServer.Address, c.HTTPServer.Port)
}

func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
