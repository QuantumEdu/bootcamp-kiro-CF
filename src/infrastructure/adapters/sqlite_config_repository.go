package adapters

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

// Compile-time interface satisfaction check.
var _ ports.ConfigRepository = (*SQLiteConfigRepository)(nil)

// SQLiteConfigRepository implements ports.ConfigRepository using SQLite.
type SQLiteConfigRepository struct {
	db *sql.DB
}

// NewSQLiteConfigRepository creates a new SQLiteConfigRepository.
func NewSQLiteConfigRepository(db *sql.DB) *SQLiteConfigRepository {
	return &SQLiteConfigRepository{db: db}
}

// Get retrieves a configuration value by key. Returns empty string if not found.
func (r *SQLiteConfigRepository) Get(ctx context.Context, clave string) (string, error) {
	var valor string
	err := r.db.QueryRowContext(ctx,
		`SELECT valor FROM configuracion WHERE clave = ?`, clave).Scan(&valor)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("getting config key %q: %w", clave, err)
	}
	return valor, nil
}

// Set upserts a configuration key-value pair using INSERT ... ON CONFLICT DO UPDATE.
func (r *SQLiteConfigRepository) Set(ctx context.Context, clave, valor string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO configuracion (clave, valor) VALUES (?, ?)
		 ON CONFLICT(clave) DO UPDATE SET valor = excluded.valor`,
		clave, valor)
	if err != nil {
		return fmt.Errorf("setting config key %q: %w", clave, err)
	}
	return nil
}
