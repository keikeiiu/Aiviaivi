package config

import (
	"os"
	"testing"
)

func setenv(t *testing.T, k, v string) {
	t.Helper()
	if err := os.Setenv(k, v); err != nil {
		t.Fatalf("setenv: %v", err)
	}
	t.Cleanup(func() { os.Unsetenv(k) })
}

func TestLoadDefaults(t *testing.T) {
	setenv(t, "DATABASE_URL", "postgres://localhost/db")
	setenv(t, "JWT_SECRET", "secret")
	// Unset optional vars
	os.Unsetenv("PORT")
	os.Unsetenv("JWT_EXPIRES_MINUTES")
	os.Unsetenv("MIGRATIONS_DIR")
	os.Unsetenv("REDIS_URL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Port != 8080 {
		t.Fatalf("expected default port 8080, got %d", cfg.Port)
	}
	if cfg.MigrationsDir != "migrations" {
		t.Fatalf("expected default migrations dir, got %s", cfg.MigrationsDir)
	}
	if cfg.JWTExpires.Minutes() != 60 {
		t.Fatalf("expected default expiry 60m, got %v", cfg.JWTExpires)
	}
	if cfg.RedisURL != "" {
		t.Fatalf("expected empty RedisURL, got %s", cfg.RedisURL)
	}
	if cfg.Storage != "local" {
		t.Fatalf("expected default storage local, got %s", cfg.Storage)
	}
}

func TestLoadCustomPort(t *testing.T) {
	setenv(t, "DATABASE_URL", "postgres://localhost/db")
	setenv(t, "JWT_SECRET", "secret")
	setenv(t, "PORT", "3000")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Port != 3000 {
		t.Fatalf("expected port 3000, got %d", cfg.Port)
	}
}

func TestLoadInvalidPort(t *testing.T) {
	setenv(t, "DATABASE_URL", "postgres://localhost/db")
	setenv(t, "JWT_SECRET", "secret")
	setenv(t, "PORT", "99999")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestLoadMissingDB(t *testing.T) {
	os.Unsetenv("DATABASE_URL")
	setenv(t, "JWT_SECRET", "secret")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing DATABASE_URL")
	}
}

func TestLoadMissingJWTSecret(t *testing.T) {
	setenv(t, "DATABASE_URL", "postgres://localhost/db")
	os.Unsetenv("JWT_SECRET")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing JWT_SECRET")
	}
}

func TestAddr(t *testing.T) {
	cfg := Config{Port: 9090}
	if cfg.Addr() != ":9090" {
		t.Fatalf("expected :9090, got %s", cfg.Addr())
	}
}
