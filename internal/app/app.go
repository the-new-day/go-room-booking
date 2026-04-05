package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/internships-backend/test-backend-the-new-day/config"
	"github.com/internships-backend/test-backend-the-new-day/internal/auth"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/router"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
	"github.com/internships-backend/test-backend-the-new-day/pkg/postgres"
)

const (
	logLevelDebug   = "debug"
	logLevelInfo    = "info"
	logLevelDiscard = "discard"
)

func Run(cfg *config.Config) {
	logger := setupLogger(cfg.Log.Level)

	logger.Info("starting server")
	logger.Debug("debug messages are enabled")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("connecting to Postgres")
	db := setupDatabase(cfg.Postgres.DSN(), cfg.Postgres.MaxPoolSize)

	jwtManager := auth.NewJwtManager(cfg.JwtConfig.SignKey, cfg.JwtConfig.AccessTTL)

	router := router.NewRouter(logger, jwtManager)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HttpServer.Port),
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to start server", sl.Err(err))
		}
	}()

	<-interrupt
	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server", sl.Err(err))
	} else {
		logger.Info("server stopped")
	}

	logger.Info("closing Postgres")
	db.Close()
	logger.Info("Postgres closed")
}

func setupLogger(level string) *slog.Logger {
	var logger *slog.Logger

	switch level {
	case logLevelDebug:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case logLevelInfo:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case logLevelDiscard:
		logger = sl.NewDiscardLogger()
	default:
		log.Fatalf("logger setup failed: unsupported level %q", level)
	}

	return logger
}

func setupDatabase(dsn string, maxPoolSize int) *postgres.Postgres {
	pg, err := postgres.New(dsn, postgres.MaxPoolSize(maxPoolSize))
	if err != nil {
		log.Fatalf("database setup failed: %v", err)
	}
	return pg
}
