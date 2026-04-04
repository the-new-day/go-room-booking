package main

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
)

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
