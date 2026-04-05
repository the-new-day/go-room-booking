package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/internships-backend/test-backend-the-new-day/pkg/hasher"
	"github.com/internships-backend/test-backend-the-new-day/pkg/postgres"
)

type postgresConfig struct {
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	User     string `env:"POSTGRES_USER" env-default:"user"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"pass"`
	Database string `env:"POSTGRES_DB" env-default:"meeting_booking"`
	SslMode  string `env:"POSTGRES_SSLMODE" env-default:"disable"`
}

func (p *postgresConfig) dsn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		p.User,
		p.Password,
		net.JoinHostPort(p.Host, p.Port),
		p.Database,
		p.SslMode,
	)
}

func main() {
	ctx := context.Background()
	var cfg postgresConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := postgres.New(cfg.dsn(), postgres.MaxPoolSize(1))
	if err != nil {
		log.Fatalf("database setup failed: %v", err)
	}
	defer db.Close()

	passwordHasher := hasher.NewBcryptHasher()

	adminHash, _ := passwordHasher.Hash("admin123")
	userHash, _ := passwordHasher.Hash("user123")

	users := []struct {
		id       string
		email    string
		password string
		role     string
	}{
		{"00000000-0000-0000-0000-000000000001", "admin@example.com", string(adminHash), "admin"},
		{"00000000-0000-0000-0000-000000000002", "user@example.com", string(userHash), "user"},
	}

	for _, u := range users {
		_, err := db.Pool.Exec(ctx,
			`INSERT INTO users (id, email, password_hash, role) 
             VALUES ($1, $2, $3, $4)
             ON CONFLICT (email) DO NOTHING`,
			u.id, u.email, u.password, u.role,
		)
		if err != nil {
			log.Printf("failed to insert user %s: %v", u.email, err)
		}
	}

	log.Println("seeding completed successfully")
}
