package app

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "github.com/CodebyTecs/wishlist-service/internal/adapters/http"
	"github.com/CodebyTecs/wishlist-service/internal/adapters/http/middleware"
	jwtadapter "github.com/CodebyTecs/wishlist-service/internal/adapters/jwt"
	"github.com/CodebyTecs/wishlist-service/internal/config"
	"github.com/CodebyTecs/wishlist-service/internal/handlers"
	"github.com/CodebyTecs/wishlist-service/internal/repository"
	"github.com/CodebyTecs/wishlist-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type App struct {
	cfg    *config.Config
	logger *slog.Logger
	server *http.Server
	router http.Handler
	dbPool *pgxpool.Pool
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
		Handler:           a.router,
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
	if a.dbPool != nil {
		a.dbPool.Close()
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

	dbCtx, cancelDB := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelDB()

	dbPool, err := pgxpool.New(dbCtx, a.cfg.DatabaseDSN())
	if err != nil {
		return err
	}
	if err := dbPool.Ping(dbCtx); err != nil {
		dbPool.Close()
		return err
	}
	a.dbPool = dbPool

	migrationsCtx, cancelMigrations := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelMigrations()
	if err := runMigrations(migrationsCtx, a.cfg.DatabaseDSN()); err != nil {
		dbPool.Close()
		return err
	}

	userRepository := repository.NewPostgresUserRepository(dbPool)
	wishlistRepository := repository.NewPostgresWishlistRepository(dbPool)
	wishlistItemRepository := repository.NewPostgresWishlistItemRepository(dbPool)

	tokenService := jwtadapter.NewService(a.cfg.JWT.Secret, a.cfg.JWT.TTL)
	authService := service.NewAuthService(userRepository, tokenService)
	userService := service.NewUserService(userRepository)
	wishlistService := service.NewWishlistService(wishlistRepository)
	wishlistItemService := service.NewWishlistItemService(wishlistItemRepository)
	publicService := service.NewPublicService(wishlistRepository, wishlistItemRepository)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	wishlistHandler := handlers.NewWishlistHandler(wishlistService)
	wishlistItemHandler := handlers.NewWishlistItemHandler(wishlistItemService)
	publicHandler := handlers.NewPublicHandler(publicService)
	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	a.router = httpadapter.NewRouter(
		authHandler,
		userHandler,
		wishlistHandler,
		wishlistItemHandler,
		publicHandler,
		authMiddleware.RequireAuth,
	)
	return nil
}

func runMigrations(ctx context.Context, dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.UpContext(ctx, db, "./migrations")
}
