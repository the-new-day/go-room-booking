package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/internships-backend/test-backend-the-new-day/internal/config"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
	"github.com/internships-backend/test-backend-the-new-day/pkg/postgres"
)

const (
	logLevelDebug   = "debug"
	logLevelInfo    = "info"
	logLevelDiscard = "discard"
)

const defaultConfigPath = "./config/config.yaml"

func loadConfig() *config.Config {
	configPath := flag.String("config", defaultConfigPath, "Path to the config file")

	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return cfg
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
