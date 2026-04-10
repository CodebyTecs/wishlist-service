package main

import (
	"log/slog"
	"os"

	"github.com/CodebyTecs/wishlist-service/internal/app"
	"github.com/CodebyTecs/wishlist-service/internal/config"
)

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
