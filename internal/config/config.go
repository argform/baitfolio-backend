package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppEnv string
	HTTPPort string
	Postgres PostgresConfig
}

type PostgresConfig struct {
	Host string
	Port string
	DB string
	User string
	Password string
	SSLMode string
}

func Load() (*Config, error) {
	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "local"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			DB:       getEnv("POSTGRES_DB", "baitfolio"),
			User:     getEnv("POSTGRES_USER", "baitfolio"),
			Password: getEnv("POSTGRES_PASSWORD", "baitfolio"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
	}

	if cfg.Postgres.Host == "" || cfg.Postgres.DB == "" || cfg.Postgres.User == "" {
		return nil, fmt.Errorf("postgres config is incomplete")
	}

	return cfg, nil
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.DB,
		p.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}