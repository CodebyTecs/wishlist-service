package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "github.com/CodebyTecs/wishlist-service/internal/adapters/http"
	"github.com/CodebyTecs/wishlist-service/internal/config"
)

type App struct {
	cfg    *config.Config
	logger *slog.Logger
	server *http.Server
}

func New(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
	}
}

func (a *App) Run() error {
	baseCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := a.wireDependencies(baseCtx); err != nil {
		return err
	}

	a.server = &http.Server{
		Addr:              a.cfg.HTTPAddr(),
		Handler:           httpadapter.NewRouter(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       a.cfg.HTTPServer.Timeout,
		WriteTimeout:      a.cfg.HTTPServer.Timeout,
		IdleTimeout:       60 * time.Second,
	}

	a.logger.Info("starting http server", "addr", a.server.Addr)
	serverErrCh := make(chan error, 1)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
	}()

	select {
	case <-baseCtx.Done():
		a.logger.Info("shutdown signal received")
	case err := <-serverErrCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	a.logger.Info("server stopped gracefully")
	return nil
}

func (a *App) wireDependencies(_ context.Context) error {
	a.logger.Info(
		"wiring dependencies",
		"http_addr", a.cfg.HTTPAddr(),
		"db_host", a.cfg.Database.Host,
		"db_port", a.cfg.Database.Port,
		"db_name", a.cfg.Database.DBName,
	)

	// TODO: connect DB pool and run migrations
	// TODO: initialize repositories
	// TODO: initialize use cases (auth, wishlist, wishlist-item, public reservation)
	// TODO: initialize HTTP handlers and auth middleware
	return nil
}
