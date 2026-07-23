package adapters

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

// Compile-time interface satisfaction check.
var _ ports.InventoryRepository = (*SQLiteInventoryRepository)(nil)

// SQLiteInventoryRepository implements ports.InventoryRepository using SQLite.
type SQLiteInventoryRepository struct {
	db *sql.DB
}

// NewSQLiteInventoryRepository creates a new SQLiteInventoryRepository.
func NewSQLiteInventoryRepository(db *sql.DB) *SQLiteInventoryRepository {
	return &SQLiteInventoryRepository{db: db}
}

// Create persists a new inventory movement in the store.
func (r *SQLiteInventoryRepository) Create(ctx context.Context, movement *entities.InventoryMovement) error {
	query := `INSERT INTO inventario_movimientos (producto_id, tipo, cantidad, stock_resultante, referencia_tipo, referencia_id, motivo, usuario_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query,
		movement.ProductoID,
		string(movement.Tipo),
		movement.Cantidad,
		movement.StockResultante,
		nullableStringPtr(movement.ReferenciaTipo),
		nullableInt64PtrField(movement.ReferenciaID),
		movement.Motivo,
		movement.UsuarioID,
	)
	if err != nil {
		return fmt.Errorf("inserting inventory movement: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last insert id: %w", err)
	}

	movement.ID = id
	return nil
}

// FindByProduct retrieves all inventory movements for a given product.
func (r *SQLiteInventoryRepository) FindByProduct(ctx context.Context, productID int64) ([]entities.InventoryMovement, error) {
	query := `SELECT id, producto_id, tipo, cantidad, stock_resultante, referencia_tipo, referencia_id, motivo, usuario_id, created_at
		FROM inventario_movimientos
		WHERE producto_id = ?
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("finding movements by product: %w", err)
	}
	defer rows.Close()

	var movements []entities.InventoryMovement
	for rows.Next() {
		movement, err := r.scanMovement(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning movement: %w", err)
		}
		movements = append(movements, *movement)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating movements: %w", err)
	}

	return movements, nil
}

// scanMovement scans a single inventory movement from *sql.Rows.
func (r *SQLiteInventoryRepository) scanMovement(rows *sql.Rows) (*entities.InventoryMovement, error) {
	var m entities.InventoryMovement
	var tipo string
	var refTipo sql.NullString
	var refID sql.NullInt64
	var createdAt string

	err := rows.Scan(
		&m.ID,
		&m.ProductoID,
		&tipo,
		&m.Cantidad,
		&m.StockResultante,
		&refTipo,
		&refID,
		&m.Motivo,
		&m.UsuarioID,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	m.Tipo = entities.MovementType(tipo)

	if refTipo.Valid {
		m.ReferenciaTipo = &refTipo.String
	}
	if refID.Valid {
		m.ReferenciaID = &refID.Int64
	}

	if t, err := parseDateTime(createdAt); err == nil {
		m.CreatedAt = t
	}

	return &m, nil
}

// nullableStringPtr returns nil for nil or empty strings.
func nullableStringPtr(s *string) interface{} {
	if s == nil || *s == "" {
		return nil
	}
	return *s
}

// nullableInt64PtrField returns nil for nil pointers.
func nullableInt64PtrField(n *int64) interface{} {
	if n == nil {
		return nil
	}
	return *n
}
