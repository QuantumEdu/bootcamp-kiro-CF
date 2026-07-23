package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

// Compile-time interface satisfaction check.
var _ ports.UserRepository = (*SQLiteUserRepository)(nil)

// SQLiteUserRepository implements ports.UserRepository using SQLite.
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository creates a new SQLiteUserRepository.
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// FindByID retrieves a user by their unique identifier.
func (r *SQLiteUserRepository) FindByID(ctx context.Context, id int64) (*entities.User, error) {
	query := `SELECT id, nombre, pin_hash, rol, activo, intentos_fallidos, bloqueado_hasta, created_at
		FROM usuarios WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanUser(row)
}

// FindByPINHash retrieves a user by their hashed PIN value.
func (r *SQLiteUserRepository) FindByPINHash(ctx context.Context, pinHash string) (*entities.User, error) {
	query := `SELECT id, nombre, pin_hash, rol, activo, intentos_fallidos, bloqueado_hasta, created_at
		FROM usuarios WHERE pin_hash = ?`

	row := r.db.QueryRowContext(ctx, query, pinHash)
	return r.scanUser(row)
}

// FindAll retrieves all users in the system.
func (r *SQLiteUserRepository) FindAll(ctx context.Context) ([]entities.User, error) {
	query := `SELECT id, nombre, pin_hash, rol, activo, intentos_fallidos, bloqueado_hasta, created_at
		FROM usuarios ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
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
func (r *SQLiteUserRepository) IncrementFailedAttempts(ctx context.Context, id int64) error {
	query := `UPDATE usuarios SET intentos_fallidos = intentos_fallidos + 1 WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("incrementing failed attempts: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Lock sets a lockout period on a user account until the specified time.
func (r *SQLiteUserRepository) Lock(ctx context.Context, id int64, until time.Time) error {
	query := `UPDATE usuarios SET bloqueado_hasta = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, until.Format(time.RFC3339), id)
	if err != nil {
		return fmt.Errorf("locking user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ResetAttempts clears the failed attempt counter and lockout for a user.
func (r *SQLiteUserRepository) ResetAttempts(ctx context.Context, id int64) error {
	query := `UPDATE usuarios SET intentos_fallidos = 0, bloqueado_hasta = NULL WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("resetting attempts: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// scanUser scans a single user from a *sql.Row.
func (r *SQLiteUserRepository) scanUser(row *sql.Row) (*entities.User, error) {
	var user entities.User
	var activo int
	var bloqueadoHasta sql.NullString
	var createdAt string

	err := row.Scan(
		&user.ID,
		&user.Nombre,
		&user.PINHash,
		&user.Rol,
		&activo,
		&user.IntentosFallidos,
		&bloqueadoHasta,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	user.Activo = activo == 1

	if bloqueadoHasta.Valid {
		t, err := time.Parse(time.RFC3339, bloqueadoHasta.String)
		if err != nil {
			return nil, fmt.Errorf("parsing bloqueado_hasta: %w", err)
		}
		user.BloqueadoHasta = t
	}

	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return nil, fmt.Errorf("parsing created_at: %w", err)
	}
	user.CreatedAt = t

	return &user, nil
}

// scanUserFromRows scans a single user from *sql.Rows.
func (r *SQLiteUserRepository) scanUserFromRows(rows *sql.Rows) (*entities.User, error) {
	var user entities.User
	var activo int
	var bloqueadoHasta sql.NullString
	var createdAt string

	err := rows.Scan(
		&user.ID,
		&user.Nombre,
		&user.PINHash,
		&user.Rol,
		&activo,
		&user.IntentosFallidos,
		&bloqueadoHasta,
		&createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning user: %w", err)
	}

	user.Activo = activo == 1

	if bloqueadoHasta.Valid {
		t, err := time.Parse(time.RFC3339, bloqueadoHasta.String)
		if err != nil {
			return nil, fmt.Errorf("parsing bloqueado_hasta: %w", err)
		}
		user.BloqueadoHasta = t
	}

	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return nil, fmt.Errorf("parsing created_at: %w", err)
	}
	user.CreatedAt = t

	return &user, nil
}
