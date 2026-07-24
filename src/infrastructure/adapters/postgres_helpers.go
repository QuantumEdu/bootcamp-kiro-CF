package adapters

import "fmt"

// wrapErr adds operation context to database errors.
// Returns nil for nil errors. Includes entity ID when > 0.
func wrapErr(err error, operation string, entityID int64) error {
	if err == nil {
		return nil
	}
	if entityID > 0 {
		return fmt.Errorf("%s (id=%d): %w", operation, entityID, err)
	}
	return fmt.Errorf("%s: %w", operation, err)
}
