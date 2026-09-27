package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations applies embedded SQL once. Only fresh or ledger-managed databases
// are supported; legacy databases need an explicit migration/import procedure.
func (db *DB) RunMigrations() error {
	return db.runMigrations(migrationsFS)
}

func (db *DB) runMigrations(source fs.FS) error {
	var hasLedger bool
	if err := db.RW.QueryRow(`SELECT EXISTS(SELECT 1 FROM sqlite_master WHERE type='table' AND name='schema_migrations')`).Scan(&hasLedger); err != nil {
		return fmt.Errorf("checking migration ledger: %w", err)
	}
	if !hasLedger {
		var tables int
		if err := db.RW.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`).Scan(&tables); err != nil {
			return fmt.Errorf("checking existing schema: %w", err)
		}
		if tables > 0 {
			return fmt.Errorf("legacy database has no migration ledger; use a fresh demo database or an explicit migration/import procedure")
		}
	}
	if _, err := db.RW.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (name TEXT PRIMARY KEY, applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP)`); err != nil {
		return fmt.Errorf("creating migration ledger: %w", err)
	}
	entries, err := fs.ReadDir(source, "migrations")
	if err != nil {
		return fmt.Errorf("reading migrations directory: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		content, err := fs.ReadFile(source, "migrations/"+entry.Name())
		if err != nil {
			return fmt.Errorf("reading migration %s: %w", entry.Name(), err)
		}
		if err := db.applyMigration(entry.Name(), string(content)); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) applyMigration(name, content string) error {
	tx, err := db.RW.Begin()
	if err != nil {
		return fmt.Errorf("beginning migration %s: %w", name, err)
	}
	defer tx.Rollback()
	var applied string
	err = tx.QueryRow(`SELECT name FROM schema_migrations WHERE name = ?`, name).Scan(&applied)
	if err == nil {
		return nil
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("checking migration %s: %w", name, err)
	}
	log.Printf("Running migration: %s", name)
	if _, err := tx.Exec(content); err != nil {
		return fmt.Errorf("executing migration %s: %w", name, err)
	}
	if _, err := tx.Exec(`INSERT INTO schema_migrations (name) VALUES (?)`, name); err != nil {
		return fmt.Errorf("recording migration %s: %w", name, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing migration %s: %w", name, err)
	}
	return nil
}
