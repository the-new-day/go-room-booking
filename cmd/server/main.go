package main

import (
	"os"
	"os/signal"
	"syscall"
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

	<-interrupt
	logger.Info("stopping server")

	logger.Info("server stopped")

	logger.Info("closing Postgres")
	db.Close()
	logger.Info("Postgres closed")
}
