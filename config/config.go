package config

import (
	"fmt"
	"net"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Log        Log
	Postgres   Postgres
	HttpServer HttpServer
	JwtConfig  Jwt
}

type Log struct {
	Level string `env:"LOG_LEVEL" env-default:"info"`
}

type Postgres struct {
	Host        string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port        string `env:"POSTGRES_PORT" env-default:"5432"`
	User        string `env:"POSTGRES_USER" env-default:"user"`
	Password    string `env:"POSTGRES_PASSWORD" env-default:"pass"`
	Database    string `env:"POSTGRES_DB" env-default:"meeting_booking"`
	SslMode     string `env:"POSTGRES_SSLMODE" env-default:"disable"`
	MaxPoolSize int    `env:"POSTGRES_MAX_POOL_SIZE" env-default:"20"`
}

type HttpServer struct {
	Port        string        `env:"HTTP_PORT" env-default:":8080"`
	Timeout     time.Duration `env:"HTTP_SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"20s"`
}

type Jwt struct {
	SignKey    string        `env:"JWT_SIGN_KEY" env-default:"secret"`
	AccessTTL  time.Duration `env:"JWT_ACCESS_TTL" env-default:"15m"`
	RefreshTTL time.Duration `env:"JWT_REFRESH_TTL" env-default:"168h"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("cannot read env: %w", err)
	}

	return &cfg, nil
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		p.User,
		p.Password,
		net.JoinHostPort(p.Host, p.Port),
		p.Database,
		p.SslMode,
	)
}
