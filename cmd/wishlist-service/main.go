package main

import (
	"log/slog"
	"os"

	"github.com/CodebyTecs/wishlist-service/internal/config"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		return
	}
	application := app.New(cfg)
	slog.Info("starting service", "port", cfg.HTTPPort)

	if err := application.Run(); err != nil {
		slog.Error("app stopped with error", "error", err)
		os.Exit(1)
	}
}
