package app

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	authapp "github.com/mohit838/olario-platform-backend/internal/application/auth"
	devapp "github.com/mohit838/olario-platform-backend/internal/application/dev"
	"github.com/mohit838/olario-platform-backend/internal/config"
	httpapi "github.com/mohit838/olario-platform-backend/internal/http"
	"github.com/mohit838/olario-platform-backend/internal/platform/logger"
	platformpostgres "github.com/mohit838/olario-platform-backend/internal/platform/postgres"
	platformredis "github.com/mohit838/olario-platform-backend/internal/platform/redis"
	"github.com/mohit838/olario-platform-backend/internal/security/token"
)

// Run owns the API application lifecycle.
// It solves the wiring problem for the process: load config, create shared
// infrastructure, start HTTP, and stop the server gracefully when ctx is done.
func Run(ctx context.Context, args []string) error {
	cfg, err := config.Load(args)
	if err != nil {
		return err
	}

	log := logger.New(cfg.Env)

	deps, cleanup, err := buildDependencies(ctx, cfg, log)
	if err != nil {
		return err
	}
	defer cleanup()

	router := httpapi.NewRouter(log, cfg, deps)

	server := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
	}

	errCh := make(chan error, 1)

	go func() {
		log.Info("http server starting", slog.String("addr", cfg.HTTP.Addr), slog.String("env", cfg.Env))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()

	startedAt := time.Now()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("http server shutdown failed", slog.String("error", err.Error()))
		return err
	}

	log.Info("http server stopped", slog.Duration("duration", time.Since(startedAt)))
	return nil
}

// buildDependencies connects infrastructure adapters and builds application
// services. Keeping this wiring in app avoids coupling cmd/api or HTTP handlers
// to concrete database and Redis setup.
func buildDependencies(ctx context.Context, cfg config.Config, log *slog.Logger) (httpapi.Dependencies, func(), error) {
	var (
		db      *sql.DB
		cleanup = func() {}
		deps    httpapi.Dependencies
	)

	if databaseConfigured(cfg.Database.DSN) {
		openedDB, err := platformpostgres.Open(ctx, cfg.Database)
		if err != nil {
			return deps, cleanup, err
		}
		db = openedDB
		cleanup = func() {
			if err := db.Close(); err != nil {
				log.Error("database close failed", slog.String("error", err.Error()))
			}
		}
	}

	if db == nil {
		log.Warn("database is not configured; database-backed routes are disabled")
		return deps, cleanup, nil
	}

	redisClient, err := platformredis.Open(ctx, cfg.Redis)
	if err != nil {
		cleanup()
		return deps, func() {}, err
	}

	previousCleanup := cleanup
	cleanup = func() {
		if err := redisClient.Close(); err != nil {
			log.Error("redis close failed", slog.String("error", err.Error()))
		}
		previousCleanup()
	}

	devRepo := platformpostgres.NewDevRepository(db)
	devCache := platformredis.NewDevCache(redisClient)
	deps.DevService = devapp.NewService(devRepo, devCache)

	authRepo := platformpostgres.NewAuthRepository(db)
	refreshStore := platformredis.NewRefreshStore(redisClient)
	tokenManager := token.NewManager(cfg.Auth.AccessTokenSecret, cfg.Auth.RefreshTokenSecret)
	deps.Auth = authapp.NewService(authRepo, refreshStore, tokenManager, cfg.Auth.AccessTokenTTL, cfg.Auth.RefreshTokenTTL)

	return deps, cleanup, nil
}

func databaseConfigured(dsn string) bool {
	return dsn != "" && !strings.Contains(dsn, "<") && !strings.Contains(dsn, ">")
}
