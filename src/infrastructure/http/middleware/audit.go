package middleware

import (
	"database/sql"
	"log"
	"time"
)

// AuditLogger logs all NL→SQL queries for auditing purposes.
type AuditLogger struct {
	db *sql.DB
}

// NewAuditLogger creates a new audit logger.
// It also creates the audit table if it doesn't exist.
func NewAuditLogger(db *sql.DB) *AuditLogger {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS audit_queries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			source TEXT NOT NULL,
			client_ip TEXT NOT NULL,
			user_query TEXT NOT NULL,
			generated_sql TEXT,
			was_allowed INTEGER NOT NULL DEFAULT 1,
			rejection_reason TEXT,
			execution_time_ms INTEGER,
			result_rows INTEGER,
			created_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
		)
	`)
	if err != nil {
		log.Printf("Warning: could not create audit_queries table: %v", err)
	}

	return &AuditLogger{db: db}
}

// AuditEntry represents a single audit log entry.
type AuditEntry struct {
	Source          string
	ClientIP        string
	UserQuery       string
	GeneratedSQL    string
	WasAllowed      bool
	RejectionReason string
	ExecutionTimeMs int64
	ResultRows      int
}

// Log records an audit entry.
func (a *AuditLogger) Log(entry AuditEntry) {
	allowed := 0
	if entry.WasAllowed {
		allowed = 1
	}

	_, err := a.db.Exec(`
		INSERT INTO audit_queries (source, client_ip, user_query, generated_sql, was_allowed, rejection_reason, execution_time_ms, result_rows)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, entry.Source, entry.ClientIP, entry.UserQuery, entry.GeneratedSQL, allowed, entry.RejectionReason, entry.ExecutionTimeMs, entry.ResultRows)

	if err != nil {
		log.Printf("Audit log error: %v", err)
	}
}

// LogQuery is a convenience method for logging allowed queries.
func (a *AuditLogger) LogQuery(source, clientIP, userQuery, generatedSQL string, execTime time.Duration, resultRows int) {
	a.Log(AuditEntry{
		Source:          source,
		ClientIP:        clientIP,
		UserQuery:       userQuery,
		GeneratedSQL:    generatedSQL,
		WasAllowed:      true,
		ExecutionTimeMs: execTime.Milliseconds(),
		ResultRows:      resultRows,
	})
}

// LogRejection logs a rejected query attempt.
func (a *AuditLogger) LogRejection(source, clientIP, userQuery, generatedSQL, reason string) {
	a.Log(AuditEntry{
		Source:          source,
		ClientIP:        clientIP,
		UserQuery:       userQuery,
		GeneratedSQL:    generatedSQL,
		WasAllowed:      false,
		RejectionReason: reason,
	})
}
