package config

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Log            Log        `yaml:"log" env-required:"true"`
	Postgres       Postgres   `yaml:"postgres" env-required:"true"`
	HttpServer     HttpServer `yaml:"http_server" env-required:"true"`
	JwtConfig      JwtConfig  `yaml:"jwt" env-required:"true"`
	PasswordHasher PasswordHasher
}

type Log struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"prod"`
}

type Postgres struct {
	Host        string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port        string `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	User        string `yaml:"username" env:"POSTGRES_USER" env-default:"user"`
	Password    string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"pass"`
	Database    string `yaml:"database" env:"POSTGRES_DB" env-default:"booking"`
	SslMode     string `yaml:"ssl_mode" env:"POSTGRES_SSL_MODE" env-default:"disable"`
	MaxPoolSize int    `yaml:"max_pool_size" env:"POSTGRES_MAX_POOL_SIZE" env-default:"20"`
}

type HttpServer struct {
	Port        string        `yaml:"port" env:"HTTP_PORT" env-default:":8080"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"20s"`
}

type JwtConfig struct {
	SignKey    string        `env:"JWT_SIGN_KEY" env-required:"true"`
	AccessTTL  time.Duration `yaml:"access_ttl" env:"JWT_ACCESS_TTL" env-default:"15m"`
	RefreshTTL time.Duration `yaml:"refresh_ttl" env:"JWT_REFRESH_TTL" env-default:"168h"`
}

type PasswordHasher struct {
	Salt string `env:"PASSWORD_HASHER_SALT" env-required:"true"`
}

func NewConfig(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %q", configPath)
	} else if err != nil {
		return nil, fmt.Errorf("failed to check config file: %q", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config file %q: %w", configPath, err)
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
