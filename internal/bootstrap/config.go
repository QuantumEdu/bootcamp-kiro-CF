package bootstrap

import (
	"fmt"
	"strings"
)

// Config holds bootstrap configuration resolved from env or Secrets Manager.
type Config struct {
	AppEnv              string
	Port                string
	DatabaseURL         string // PostgreSQL connection string (lambda)
	DatabasePath        string // SQLite path (local)
	SessionSecret       string
	BedrockModelID      string
	BedrockRegion       string
	MaxTokens           int
	Temperature         float64
	OpenRouterAPIKey    string
	OpenRouterModel     string
	QueryTimeoutSeconds int
	PINMaxAttempts      int
	PINLockoutMinutes   int
}

// ValidateConfig checks that all required fields are present for the given mode.
func ValidateConfig(cfg *Config) error {
	var missing []string

	if cfg.SessionSecret == "" {
		missing = append(missing, "SESSION_SECRET")
	}

	if cfg.AppEnv == "lambda" {
		if cfg.DatabaseURL == "" {
			missing = append(missing, "DATABASE_URL (from SECRET_DB_ARN)")
		}
		if cfg.BedrockModelID == "" {
			missing = append(missing, "BEDROCK_MODEL_ID (from SECRET_AI_ARN)")
		}
	} else {
		if cfg.DatabasePath == "" {
			missing = append(missing, "DATABASE_PATH")
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config: %s", strings.Join(missing, ", "))
	}
	return nil
}
