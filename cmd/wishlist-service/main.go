package main

import (
	"log/slog"
	"os"

	"github.com/CodebyTecs/wishlist-service/internal/app"
	"github.com/CodebyTecs/wishlist-service/internal/config"
)

// @title Wishlist Service API
// @version 1.0
// @description REST API для вишлистов: авторизация, CRUD вишлистов/подарков, публичный просмотр и резервирование.
// @BasePath /
// @schemes http
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT access token with Bearer prefix. Example: "Bearer <token>"
func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	application := app.New(cfg)
	slog.Info("starting service", "addr", cfg.HTTPAddr())

	if err := application.Run(); err != nil {
		slog.Error("app stopped with error", "error", err)
		os.Exit(1)
	}
}
