package adapters

import (
	"context"
	"fmt"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Compile-time interface satisfaction check.
var _ ports.InventoryRepository = (*PostgresInventoryRepository)(nil)

// PostgresInventoryRepository implements ports.InventoryRepository using pgxpool.
type PostgresInventoryRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresInventoryRepository creates a new PostgresInventoryRepository.
func NewPostgresInventoryRepository(pool *pgxpool.Pool) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{pool: pool}
}

// Create persists a new inventory movement in the store.
func (r *PostgresInventoryRepository) Create(ctx context.Context, movement *entities.InventoryMovement) error {
	const q = `INSERT INTO inventario_movimientos 
		(producto_id, tipo, cantidad, stock_resultante, referencia_tipo, referencia_id, motivo, usuario_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	err := r.pool.QueryRow(ctx, q,
		movement.ProductoID,
		string(movement.Tipo),
		movement.Cantidad,
		movement.StockResultante,
		movement.ReferenciaTipo,
		movement.ReferenciaID,
		movement.Motivo,
		movement.UsuarioID,
	).Scan(&movement.ID, &movement.CreatedAt)
	if err != nil {
		return fmt.Errorf("inserting inventory movement: %w", err)
	}

	return nil
}

// FindByProduct retrieves all inventory movements for a given product.
func (r *PostgresInventoryRepository) FindByProduct(ctx context.Context, productID int64) ([]entities.InventoryMovement, error) {
	const q = `SELECT id, producto_id, tipo, cantidad, stock_resultante, 
		referencia_tipo, referencia_id, motivo, usuario_id, created_at
		FROM inventario_movimientos
		WHERE producto_id = $1
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, q, productID)
	if err != nil {
		return nil, wrapErr(err, "finding movements by product", productID)
	}
	defer rows.Close()

	var movements []entities.InventoryMovement
	for rows.Next() {
		var m entities.InventoryMovement
		var tipo string

		err := rows.Scan(
			&m.ID,
			&m.ProductoID,
			&tipo,
			&m.Cantidad,
			&m.StockResultante,
			&m.ReferenciaTipo,
			&m.ReferenciaID,
			&m.Motivo,
			&m.UsuarioID,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning inventory movement: %w", err)
		}

		m.Tipo = entities.MovementType(tipo)
		movements = append(movements, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating inventory movements: %w", err)
	}

	return movements, nil
}
