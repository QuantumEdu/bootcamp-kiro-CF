package adapters

import (
	"context"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Compile-time interface satisfaction check.
var _ ports.ConfigRepository = (*PostgresConfigRepository)(nil)

// PostgresConfigRepository implements ports.ConfigRepository using pgxpool.
type PostgresConfigRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresConfigRepository creates a new PostgresConfigRepository.
func NewPostgresConfigRepository(pool *pgxpool.Pool) *PostgresConfigRepository {
	return &PostgresConfigRepository{pool: pool}
}

// Get retrieves a configuration value by key. Returns empty string if not found.
func (r *PostgresConfigRepository) Get(ctx context.Context, clave string) (string, error) {
	var valor string
	err := r.pool.QueryRow(ctx,
		`SELECT valor FROM configuracion WHERE clave = $1`, clave).Scan(&valor)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	return valor, wrapErr(err, "getting config", 0)
}

// Set upserts a configuration key-value pair using INSERT ... ON CONFLICT DO UPDATE.
func (r *PostgresConfigRepository) Set(ctx context.Context, clave, valor string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO configuracion (clave, valor) VALUES ($1, $2)
		 ON CONFLICT (clave) DO UPDATE SET valor = EXCLUDED.valor`, clave, valor)
	return wrapErr(err, "setting config", 0)
}
