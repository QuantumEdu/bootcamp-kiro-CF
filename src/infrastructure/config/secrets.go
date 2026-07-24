package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// AISecretPayload represents the JSON structure stored in the AI secret.
type AISecretPayload struct {
	ModelID     string  `json:"model_id"`
	Region      string  `json:"region"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

// LambdaConfig holds configuration values retrieved from Secrets Manager
// for use in Lambda mode. It provides the fields that the bootstrap package
// needs to wire PostgreSQL, session management, and Bedrock adapters.
type LambdaConfig struct {
	DatabaseURL   string
	SessionSecret string
	BedrockModelID string
	BedrockRegion  string
	MaxTokens      int
	Temperature    float64
}

// SecretsLoader retrieves and caches secrets from AWS Secrets Manager.
type SecretsLoader struct {
	client *secretsmanager.Client
	cache  map[string]string
	mu     sync.RWMutex
}

// NewSecretsLoader creates a SecretsLoader backed by the given Secrets Manager client.
func NewSecretsLoader(client *secretsmanager.Client) *SecretsLoader {
	return &SecretsLoader{client: client, cache: make(map[string]string)}
}

// GetSecret retrieves a secret value by ARN, using an in-memory cache on
// subsequent calls to avoid repeated API requests on warm Lambda starts.
func (s *SecretsLoader) GetSecret(ctx context.Context, arn string) (string, error) {
	s.mu.RLock()
	if val, ok := s.cache[arn]; ok {
		s.mu.RUnlock()
		return val, nil
	}
	s.mu.RUnlock()

	out, err := s.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &arn,
	})
	if err != nil {
		return "", fmt.Errorf("retrieving secret %s: %w", arn, err)
	}

	s.mu.Lock()
	s.cache[arn] = *out.SecretString
	s.mu.Unlock()

	return *out.SecretString, nil
}

// LoadConfig retrieves DB URL, session key, and AI configuration from Secrets
// Manager using ARNs specified in environment variables. It returns a descriptive
// error listing any missing ARNs or failed retrievals, suitable for aborting
// startup with a clear diagnostic message.
func (s *SecretsLoader) LoadConfig(ctx context.Context) (*LambdaConfig, error) {
	dbARN := os.Getenv("SECRET_DB_ARN")
	sessionARN := os.Getenv("SECRET_SESSION_ARN")
	aiARN := os.Getenv("SECRET_AI_ARN")

	// Validate all required ARNs are present before making any API calls.
	var missing []string
	if dbARN == "" {
		missing = append(missing, "SECRET_DB_ARN")
	}
	if sessionARN == "" {
		missing = append(missing, "SECRET_SESSION_ARN")
	}
	if aiARN == "" {
		missing = append(missing, "SECRET_AI_ARN")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required secret ARN environment variables: %v", missing)
	}

	// Retrieve each secret, collecting all errors for a single diagnostic message.
	var errors []string

	dbURL, err := s.GetSecret(ctx, dbARN)
	if err != nil {
		errors = append(errors, fmt.Sprintf("DB secret (%s): %v", dbARN, err))
	}

	sessionKey, err := s.GetSecret(ctx, sessionARN)
	if err != nil {
		errors = append(errors, fmt.Sprintf("session secret (%s): %v", sessionARN, err))
	}

	aiRaw, err := s.GetSecret(ctx, aiARN)
	if err != nil {
		errors = append(errors, fmt.Sprintf("AI secret (%s): %v", aiARN, err))
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("failed to retrieve secrets at startup: %v", errors)
	}

	// Parse AI config as JSON.
	var aiPayload AISecretPayload
	if err := json.Unmarshal([]byte(aiRaw), &aiPayload); err != nil {
		return nil, fmt.Errorf("parsing AI secret JSON: %w", err)
	}

	return &LambdaConfig{
		DatabaseURL:    dbURL,
		SessionSecret:  sessionKey,
		BedrockModelID: aiPayload.ModelID,
		BedrockRegion:  aiPayload.Region,
		MaxTokens:      aiPayload.MaxTokens,
		Temperature:    aiPayload.Temperature,
	}, nil
}
