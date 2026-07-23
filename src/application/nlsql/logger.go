package nlsql

import (
	"context"
	"database/sql"
	"fmt"
)

// QueryLogEntry holds the data for a single audit log entry.
type QueryLogEntry struct {
	UserID          *int64
	Question        string
	GeneratedSQL    string
	Success         bool
	ErrorMessage    string
	ExecutionTimeMs int64
}

// QueryLogger writes NL→SQL query audit logs to the database.
type QueryLogger struct {
	db *sql.DB
}

// NewQueryLogger creates a new QueryLogger with the given write connection.
func NewQueryLogger(db *sql.DB) *QueryLogger {
	return &QueryLogger{db: db}
}

// Log inserts a query log entry into the query_log table.
func (l *QueryLogger) Log(ctx context.Context, entry QueryLogEntry) error {
	if l == nil || l.db == nil {
		return nil
	}

	successInt := 0
	if entry.Success {
		successInt = 1
	}

	_, err := l.db.ExecContext(ctx,
		`INSERT INTO query_log (user_id, question, generated_sql, success, error_message, execution_time_ms)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		entry.UserID,
		entry.Question,
		entry.GeneratedSQL,
		successInt,
		entry.ErrorMessage,
		entry.ExecutionTimeMs,
	)
	if err != nil {
		return fmt.Errorf("logging query: %w", err)
	}
	return nil
}
