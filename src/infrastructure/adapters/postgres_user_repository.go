package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Compile-time interface satisfaction check.
var _ ports.UserRepository = (*PostgresUserRepository)(nil)

// PostgresUserRepository implements ports.UserRepository using pgxpool.
type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresUserRepository creates a new PostgresUserRepository.
func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

// FindByID retrieves a user by their unique identifier.
func (r *PostgresUserRepository) FindByID(ctx context.Context, id int64) (*entities.User, error) {
	const query = `SELECT id, nombre, pin_hash, rol, activo, failed_attempts, locked_until, created_at
		FROM usuarios WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)
	user, err := r.scanUser(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, wrapErr(err, "finding user by id", id)
	}
	return user, nil
}

// FindByPINHash retrieves a user by their hashed PIN value.
func (r *PostgresUserRepository) FindByPINHash(ctx context.Context, pinHash string) (*entities.User, error) {
	const query = `SELECT id, nombre, pin_hash, rol, activo, failed_attempts, locked_until, created_at
		FROM usuarios WHERE pin_hash = $1`

	row := r.pool.QueryRow(ctx, query, pinHash)
	user, err := r.scanUser(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, wrapErr(err, "finding user by pin hash", 0)
	}
	return user, nil
}

// FindAll retrieves all users in the system.
func (r *PostgresUserRepository) FindAll(ctx context.Context) ([]entities.User, error) {
	const query = `SELECT id, nombre, pin_hash, rol, activo, failed_attempts, locked_until, created_at
		FROM usuarios ORDER BY id`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying all users: %w", err)
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		user, err := r.scanUserFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning user row: %w", err)
		}
		users = append(users, *user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating user rows: %w", err)
	}

	return users, nil
}

// IncrementFailedAttempts increases the failed login attempt counter for a user.
func (r *PostgresUserRepository) IncrementFailedAttempts(ctx context.Context, id int64) error {
	const query = `UPDATE usuarios SET failed_attempts = failed_attempts + 1 WHERE id = $1`

	ct, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return wrapErr(err, "incrementing failed attempts", id)
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// Lock sets a lockout period on a user account until the specified time.
func (r *PostgresUserRepository) Lock(ctx context.Context, id int64, until time.Time) error {
	const query = `UPDATE usuarios SET locked_until = $1 WHERE id = $2`

	ct, err := r.pool.Exec(ctx, query, until, id)
	if err != nil {
		return wrapErr(err, "locking user", id)
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// ResetAttempts clears the failed attempt counter and lockout for a user.
func (r *PostgresUserRepository) ResetAttempts(ctx context.Context, id int64) error {
	const query = `UPDATE usuarios SET failed_attempts = 0, locked_until = NULL WHERE id = $1`

	ct, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return wrapErr(err, "resetting attempts", id)
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// scanUser scans a single user from a pgx.Row.
func (r *PostgresUserRepository) scanUser(row pgx.Row) (*entities.User, error) {
	var user entities.User
	var lockedUntil *time.Time

	err := row.Scan(
		&user.ID,
		&user.Nombre,
		&user.PINHash,
		&user.Rol,
		&user.Activo,
		&user.IntentosFallidos,
		&lockedUntil,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if lockedUntil != nil {
		user.BloqueadoHasta = *lockedUntil
	}

	return &user, nil
}

// scanUserFromRows scans a single user from pgx.Rows.
func (r *PostgresUserRepository) scanUserFromRows(rows pgx.Rows) (*entities.User, error) {
	var user entities.User
	var lockedUntil *time.Time

	err := rows.Scan(
		&user.ID,
		&user.Nombre,
		&user.PINHash,
		&user.Rol,
		&user.Activo,
		&user.IntentosFallidos,
		&lockedUntil,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning user: %w", err)
	}

	if lockedUntil != nil {
		user.BloqueadoHasta = *lockedUntil
	}

	return &user, nil
}
