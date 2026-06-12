package config

import (
	"fmt"
	"os"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	AppPort     string
	DatabaseURL string
	Environment string
}

// Load reads configuration from the environment, applying sensible defaults
// for local development. DatabaseURL can be provided directly or assembled
// from the individual DB_* variables.
func Load() *Config {
	cfg := &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		Environment: getEnv("APP_ENV", "development"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	if cfg.DatabaseURL == "" {
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "postgres")
		name := getEnv("DB_NAME", "user_age_api")
		sslmode := getEnv("DB_SSLMODE", "disable")
		cfg.DatabaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, name, sslmode,
		)
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
