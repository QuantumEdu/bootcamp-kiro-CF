package database

import (
	"database/sql"
	"fmt"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
)

// UserRepo handles user operations.
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo creates a new user repository.
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetByPinHash finds a user by their PIN hash.
func (r *UserRepo) GetByPinHash(pinHash string) (*entities.User, error) {
	var u entities.User
	err := r.db.QueryRow(`
		SELECT id, nombre, pin_hash, rol, activo, created_at
		FROM usuarios WHERE pin_hash = ? AND activo = 1
	`, pinHash).Scan(&u.ID, &u.Nombre, &u.PinHash, &u.Rol, &u.Activo, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("getting user by pin: %w", err)
	}
	return &u, nil
}

// GetByID finds a user by ID.
func (r *UserRepo) GetByID(id int64) (*entities.User, error) {
	var u entities.User
	err := r.db.QueryRow(`
		SELECT id, nombre, pin_hash, rol, activo, created_at
		FROM usuarios WHERE id = ?
	`, id).Scan(&u.ID, &u.Nombre, &u.PinHash, &u.Rol, &u.Activo, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("getting user by id: %w", err)
	}
	return &u, nil
}
