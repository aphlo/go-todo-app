package config

import (
	"fmt"
	"os"
)

// Config captures runtime configuration needed by the HTTP server.
type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

// Load reads configuration from environment variables, applying sensible defaults
// for local development.
func Load() Config {
	port := getEnv("PORT", "8080")
	addr := fmt.Sprintf(":%s", port)

	cfg := Config{
		HTTPAddr:    addr,
		DatabaseURL: getEnv("DATABASE_URL", "postgres://todo:todo@localhost:5432/todo?sslmode=disable"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
