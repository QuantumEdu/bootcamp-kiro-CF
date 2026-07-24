package adapters

import (
	"context"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Compile-time interface satisfaction check.
var _ ports.ClientRepository = (*PostgresClientRepository)(nil)

// PostgresClientRepository implements ports.ClientRepository using pgxpool.
type PostgresClientRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresClientRepository creates a new PostgresClientRepository.
func NewPostgresClientRepository(pool *pgxpool.Pool) *PostgresClientRepository {
	return &PostgresClientRepository{pool: pool}
}

// Create persists a new client in the clientes table and sets the client ID via RETURNING.
func (r *PostgresClientRepository) Create(ctx context.Context, client *entities.Client) error {
	const q = `INSERT INTO clientes (nombre, telefono, direccion)
		VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.pool.QueryRow(ctx, q,
		client.Nombre, client.Telefono, client.Direccion,
	).Scan(&client.ID, &client.CreatedAt)
	return wrapErr(err, "inserting client", 0)
}

// List retrieves all clients ordered by nombre ASC.
func (r *PostgresClientRepository) List(ctx context.Context) ([]entities.Client, error) {
	const q = `SELECT id, nombre, telefono, direccion, created_at
		FROM clientes ORDER BY nombre ASC`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, wrapErr(err, "listing clients", 0)
	}
	defer rows.Close()

	var clients []entities.Client
	for rows.Next() {
		var c entities.Client
		if err := rows.Scan(&c.ID, &c.Nombre, &c.Telefono, &c.Direccion, &c.CreatedAt); err != nil {
			return nil, wrapErr(err, "scanning client row", 0)
		}
		clients = append(clients, c)
	}

	if err := rows.Err(); err != nil {
		return nil, wrapErr(err, "iterating client rows", 0)
	}

	return clients, nil
}
