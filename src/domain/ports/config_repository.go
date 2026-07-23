package ports

import (
	"context"
	"strings"
)

// ConfigRepository defines the contract for key-value configuration persistence.
type ConfigRepository interface {
	// Get retrieves a configuration value by key. Returns empty string if not found.
	Get(ctx context.Context, clave string) (string, error)

	// Set upserts a configuration key-value pair.
	Set(ctx context.Context, clave, valor string) error
}

// MaskAPIKey returns a masked version of an API key, showing only the last 4 characters.
// For strings of length >= 4, it replaces all but the last 4 chars with asterisks.
// For empty strings, it returns an empty string.
// For strings shorter than 4 chars, it returns all asterisks.
func MaskAPIKey(key string) string {
	if key == "" {
		return ""
	}
	if len(key) < 4 {
		return strings.Repeat("*", len(key))
	}
	return strings.Repeat("*", len(key)-4) + key[len(key)-4:]
}
