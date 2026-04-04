package app

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/internships-backend/test-backend-the-new-day/internal/config"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
)

func Run(configPath string) {
	cfg, err := config.NewConfig("./config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger, err := setupLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("failed to setup logger: %v", err)
	}

	logger.Info("starting server")
	logger.Debug("debug messages are enabled")

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done
	logger.Info("stopping server")

	logger.Info("server stopped")
}

func setupLogger(level string) (*slog.Logger, error) {
	var logger *slog.Logger

	switch level {
	case LogLevelDebug:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case LogLevelInfo:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return nil, fmt.Errorf("login setup failed: insupported level %q", level)
	}

	return logger, nil
}
