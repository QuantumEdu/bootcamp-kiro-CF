package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port                string
	DatabasePath        string
	OpenRouterAPIKey    string
	OpenRouterModel     string
	SessionSecret       string
	PINMaxAttempts      int
	PINLockoutMinutes   int
	QueryTimeoutSeconds int
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "8080"),
		DatabasePath:        getEnv("DATABASE_PATH", "./data/pos.db"),
		OpenRouterAPIKey:    getEnv("OPENROUTER_API_KEY", ""),
		OpenRouterModel:     getEnv("OPENROUTER_MODEL", "anthropic/claude-3-haiku"),
		SessionSecret:       getEnv("SESSION_SECRET", "dev-secret-change-me"),
		PINMaxAttempts:      getEnvInt("PIN_MAX_ATTEMPTS", 5),
		PINLockoutMinutes:   getEnvInt("PIN_LOCKOUT_MINUTES", 5),
		QueryTimeoutSeconds: getEnvInt("QUERY_TIMEOUT_SECONDS", 5),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}
