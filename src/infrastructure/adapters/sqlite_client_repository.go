package adapters

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

// Compile-time interface satisfaction check.
var _ ports.ClientRepository = (*SQLiteClientRepository)(nil)

// SQLiteClientRepository implements ports.ClientRepository using SQLite.
type SQLiteClientRepository struct {
	db *sql.DB
}

// NewSQLiteClientRepository creates a new SQLiteClientRepository.
func NewSQLiteClientRepository(db *sql.DB) *SQLiteClientRepository {
	return &SQLiteClientRepository{db: db}
}

// Create persists a new client in the clientes table and sets the client ID from LastInsertId.
func (r *SQLiteClientRepository) Create(ctx context.Context, client *entities.Client) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO clientes (nombre, telefono, direccion) VALUES (?, ?, ?)`,
		client.Nombre, client.Telefono, client.Direccion)
	if err != nil {
		return fmt.Errorf("inserting client: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last insert id: %w", err)
	}

	client.ID = id
	return nil
}

// List retrieves all clients ordered by nombre ASC.
func (r *SQLiteClientRepository) List(ctx context.Context) ([]entities.Client, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, nombre, telefono, direccion, created_at FROM clientes ORDER BY nombre ASC`)
	if err != nil {
		return nil, fmt.Errorf("listing clients: %w", err)
	}
	defer rows.Close()

	var clients []entities.Client
	for rows.Next() {
		var c entities.Client
		var createdAt string
		if err := rows.Scan(&c.ID, &c.Nombre, &c.Telefono, &c.Direccion, &createdAt); err != nil {
			return nil, fmt.Errorf("scanning client row: %w", err)
		}
		if t, err := parseDateTime(createdAt); err == nil {
			c.CreatedAt = t
		}
		clients = append(clients, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating client rows: %w", err)
	}

	return clients, nil
}
