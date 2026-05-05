package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port          int
	DatabaseURL   string
	JWTSecret     string
	JWTExpires    time.Duration
	MigrationsDir string
}

func Load() (Config, error) {
	var cfg Config

	port, err := parsePort(getenv("PORT", "8080"))
	if err != nil {
		return Config{}, errors.New("PORT is invalid")
	}
	cfg.Port = port

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	cfg.MigrationsDir = getenv("MIGRATIONS_DIR", "migrations")

	expiresMin, err := strconv.Atoi(getenv("JWT_EXPIRES_MINUTES", "60"))
	if err != nil {
		return Config{}, errors.New("JWT_EXPIRES_MINUTES is invalid")
	}
	cfg.JWTExpires = time.Duration(expiresMin) * time.Minute

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT_SECRET is required")
	}

	return cfg, nil
}

func parsePort(s string) (int, error) {
	port, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if port < 1 || port > 65535 {
		return 0, errors.New("port out of range")
	}
	return port, nil
}

func (c Config) Addr() string {
	return ":" + strconv.Itoa(c.Port)
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
