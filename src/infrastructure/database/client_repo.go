package database

import (
	"database/sql"
	"fmt"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
)

// ClientRepo handles client CRUD operations.
type ClientRepo struct {
	db *sql.DB
}

// NewClientRepo creates a new client repository.
func NewClientRepo(db *sql.DB) *ClientRepo {
	return &ClientRepo{db: db}
}

// List returns all clients with optional search.
func (r *ClientRepo) List(search string) ([]entities.Client, error) {
	query := `SELECT id, nombre, COALESCE(telefono, ''), COALESCE(direccion, ''), created_at FROM clientes`
	var args []interface{}

	if search != "" {
		query += ` WHERE nombre LIKE ? OR telefono LIKE ?`
		s := "%" + search + "%"
		args = append(args, s, s)
	}
	query += ` ORDER BY nombre ASC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing clients: %w", err)
	}
	defer rows.Close()

	var clients []entities.Client
	for rows.Next() {
		var c entities.Client
		if err := rows.Scan(&c.ID, &c.Nombre, &c.Telefono, &c.Direccion, &c.CreatedAt); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, rows.Err()
}

// Create inserts a new client.
func (r *ClientRepo) Create(c *entities.Client) (int64, error) {
	res, err := r.db.Exec(`
		INSERT INTO clientes (nombre, telefono, direccion) VALUES (?, ?, ?)
	`, c.Nombre, nullStr(c.Telefono), nullStr(c.Direccion))
	if err != nil {
		return 0, fmt.Errorf("creating client: %w", err)
	}
	return res.LastInsertId()
}

// GetByID returns a client by ID.
func (r *ClientRepo) GetByID(id int64) (*entities.Client, error) {
	var c entities.Client
	err := r.db.QueryRow(`
		SELECT id, nombre, COALESCE(telefono, ''), COALESCE(direccion, ''), created_at
		FROM clientes WHERE id = ?
	`, id).Scan(&c.ID, &c.Nombre, &c.Telefono, &c.Direccion, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
