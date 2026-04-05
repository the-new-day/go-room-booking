package main

import (
	"log"

	"github.com/internships-backend/test-backend-the-new-day/config"
	"github.com/internships-backend/test-backend-the-new-day/internal/app"
)

func main() {
	cfg := loadConfig()
	app.Run(cfg)
}

func loadConfig() *config.Config {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return cfg
}
