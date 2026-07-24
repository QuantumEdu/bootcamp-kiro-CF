package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/QuantumEdu/bootcamp-kiro-CF/internal/bootstrap"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// Validate SESSION_SECRET
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatalf("SESSION_SECRET environment variable is required but not set")
	}

	queryTimeout, _ := strconv.Atoi(os.Getenv("QUERY_TIMEOUT_SECONDS"))
	if queryTimeout == 0 {
		queryTimeout = 5
	}

	cfg := bootstrap.Config{
		AppEnv:              os.Getenv("APP_ENV"), // defaults to "" which means local
		Port:                getEnvDefault("PORT", "8080"),
		DatabasePath:        getEnvDefault("DATABASE_PATH", "./data/pos.db"),
		SessionSecret:       sessionSecret,
		OpenRouterAPIKey:    os.Getenv("OPENROUTER_API_KEY"),
		OpenRouterModel:     getEnvDefault("OPENROUTER_MODEL", "anthropic/claude-3-haiku"),
		QueryTimeoutSeconds: queryTimeout,
	}

	router, cleanup, err := bootstrap.BuildRouter(cfg)
	if err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}
	defer cleanup()

	addr := ":" + cfg.Port
	log.Printf("POS AI-First starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnvDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
