package ports

import "context"

// NLSQLResult represents the result of a natural language to SQL conversion.
type NLSQLResult struct {
	SQL         string
	Explanation string
	Columns     []string
	Rows        [][]string
	Error       string
}

// NLSQLService defines the port for natural language to SQL conversion.
type NLSQLService interface {
	// ProcessQuery takes a natural language query and returns SQL results.
	ProcessQuery(ctx context.Context, query string) (*NLSQLResult, error)
}
