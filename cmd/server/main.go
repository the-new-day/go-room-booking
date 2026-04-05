package main

import (
	"github.com/internships-backend/test-backend-the-new-day/internal/app"
)

func main() {
	app.Run(app.LoadConfig())
}
