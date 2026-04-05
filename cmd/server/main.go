package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/router"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
)

func main() {
	cfg := loadConfig()

	logger := setupLogger(cfg.Log.Level)

	logger.Info("starting server")
	logger.Debug("debug messages are enabled")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("connecting to Postgres")
	db := setupDatabase(cfg.Postgres.DSN(), cfg.Postgres.MaxPoolSize)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HttpServer.Port),
		Handler:      router.NewRouter(logger),
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

	logger.Info("server stopped")

	logger.Info("closing Postgres")
	db.Close()
	logger.Info("Postgres closed")
}
