package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB holds read-write and read-only database connections.
type DB struct {
	RW *sql.DB // Read-write connection for mutations
	RO *sql.DB // Read-only connection for NL→SQL queries
}

// New creates a new DB instance with RW and RO connections.
// It ensures the data directory exists and configures WAL mode.
func New(dbPath string) (*DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating database directory: %w", err)
	}

	// Open read-write connection
	rw, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("opening RW connection: %w", err)
	}

	// Verify connection works
	if err := rw.Ping(); err != nil {
		rw.Close()
		return nil, fmt.Errorf("pinging RW connection: %w", err)
	}

	// Open read-only connection for NL→SQL queries
	ro, err := sql.Open("sqlite", dbPath+"?mode=ro&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)")
	if err != nil {
		rw.Close()
		return nil, fmt.Errorf("opening RO connection: %w", err)
	}

	if err := ro.Ping(); err != nil {
		rw.Close()
		ro.Close()
		return nil, fmt.Errorf("pinging RO connection: %w", err)
	}

	// Set query_only on RO connection
	if _, err := ro.Exec("PRAGMA query_only=ON"); err != nil {
		rw.Close()
		ro.Close()
		return nil, fmt.Errorf("setting query_only on RO: %w", err)
	}

	db := &DB{RW: rw, RO: ro}

	// Run migrations automatically
	if err := db.RunMigrations(); err != nil {
		db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return db, nil
}

// NewInMemory creates an in-memory database for testing.
func NewInMemory() (*DB, error) {
	rw, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("opening in-memory RW: %w", err)
	}

	if err := rw.Ping(); err != nil {
		rw.Close()
		return nil, fmt.Errorf("pinging in-memory RW: %w", err)
	}

	db := &DB{RW: rw, RO: rw} // In-memory: share same connection

	if err := db.RunMigrations(); err != nil {
		db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return db, nil
}

// Close closes both database connections.
func (db *DB) Close() error {
	var firstErr error
	if db.RO != nil && db.RO != db.RW {
		if err := db.RO.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if db.RW != nil {
		if err := db.RW.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
