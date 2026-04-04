package main

import (
	"flag"

	"github.com/internships-backend/test-backend-the-new-day/internal/app"
)

const DefaultConfigPath = "./config/config.yaml"

func main() {
	configPath := flag.String("config", DefaultConfigPath, "Path to the config file")
	app.Run(*configPath)
}
