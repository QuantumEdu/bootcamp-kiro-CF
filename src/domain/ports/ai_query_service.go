package ports

import "context"

// AIQueryService defines the contract for natural language to SQL translation.
type AIQueryService interface {
	// GenerateSQL translates a natural language question into a SQL query.
	// Returns the generated SQL, a human-readable explanation, and any error.
	GenerateSQL(ctx context.Context, question string) (sql string, explanation string, err error)
}
