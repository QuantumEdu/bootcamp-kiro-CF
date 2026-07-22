package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration.
type Config struct {
	Port               string
	DatabasePath       string
	OpenRouterAPIKey   string
	OpenRouterModel    string
	SessionSecret      string
	PINMaxAttempts     int
	PINLockoutMinutes  int
	QueryTimeoutSecs   int
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		DatabasePath:      getEnv("DATABASE_PATH", "./data/pos.db"),
		OpenRouterAPIKey:  getEnv("OPENROUTER_API_KEY", ""),
		OpenRouterModel:   getEnv("OPENROUTER_MODEL", "openai/gpt-4o-mini"),
		SessionSecret:     getEnv("SESSION_SECRET", "dev-secret-change-in-prod-32ch"),
		PINMaxAttempts:    getEnvInt("PIN_MAX_ATTEMPTS", 5),
		PINLockoutMinutes: getEnvInt("PIN_LOCKOUT_MINUTES", 5),
		QueryTimeoutSecs:  getEnvInt("QUERY_TIMEOUT_SECONDS", 5),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
