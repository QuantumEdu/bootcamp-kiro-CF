package database

import (
	"database/sql"
	"fmt"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
)

// SaleRepo handles sale CRUD operations.
type SaleRepo struct {
	db *sql.DB
}

// NewSaleRepo creates a new sale repository.
func NewSaleRepo(db *sql.DB) *SaleRepo {
	return &SaleRepo{db: db}
}

// Create inserts a new sale with its items inside a transaction.
func (r *SaleRepo) Create(req *entities.CreateSaleRequest, usuarioID int64) (*entities.Sale, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	// Calculate total
	var total float64
	for _, item := range req.Items {
		total += item.Cantidad * item.PrecioUnitario
	}

	// Insert sale header
	res, err := tx.Exec(`
		INSERT INTO ventas (usuario_id, cliente_id, total, metodo_pago)
		VALUES (?, ?, ?, ?)
	`, usuarioID, req.ClienteID, total, req.MetodoPago)
	if err != nil {
		return nil, fmt.Errorf("inserting sale: %w", err)
	}

	ventaID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting sale id: %w", err)
	}

	// Insert items and update stock
	for _, item := range req.Items {
		subtotal := item.Cantidad * item.PrecioUnitario
		_, err := tx.Exec(`
			INSERT INTO venta_items (venta_id, producto_id, cantidad, precio_unitario, subtotal)
			VALUES (?, ?, ?, ?, ?)
		`, ventaID, item.ProductoID, item.Cantidad, item.PrecioUnitario, subtotal)
		if err != nil {
			return nil, fmt.Errorf("inserting sale item: %w", err)
		}

		// Update stock directly
		_, err = tx.Exec(`
			UPDATE productos SET stock_actual = stock_actual - ?, updated_at = datetime('now','localtime')
			WHERE id = ?
		`, item.Cantidad, item.ProductoID)
		if err != nil {
			return nil, fmt.Errorf("updating stock for product %d: %w", item.ProductoID, err)
		}

		// Record inventory movement
		var stockResult float64
		err = tx.QueryRow(`SELECT stock_actual FROM productos WHERE id = ?`, item.ProductoID).Scan(&stockResult)
		if err != nil {
			return nil, fmt.Errorf("getting stock result: %w", err)
		}

		_, err = tx.Exec(`
			INSERT INTO inventario_movimientos (producto_id, tipo, cantidad, stock_resultante, referencia_tipo, referencia_id, usuario_id)
			VALUES (?, 'salida', ?, ?, 'venta', ?, ?)
		`, item.ProductoID, item.Cantidad, stockResult, ventaID, usuarioID)
		if err != nil {
			return nil, fmt.Errorf("inserting inventory movement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	sale := &entities.Sale{
		ID:         ventaID,
		UsuarioID:  usuarioID,
		ClienteID:  req.ClienteID,
		Total:      total,
		MetodoPago: req.MetodoPago,
	}
	return sale, nil
}

// ListRecent returns the last N sales.
func (r *SaleRepo) ListRecent(limit int) ([]entities.Sale, error) {
	rows, err := r.db.Query(`
		SELECT v.id, v.total, v.metodo_pago, v.created_at, COALESCE(c.nombre, 'Publico general')
		FROM ventas v
		LEFT JOIN clientes c ON c.id = v.cliente_id
		ORDER BY v.created_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("listing recent sales: %w", err)
	}
	defer rows.Close()

	var sales []entities.Sale
	for rows.Next() {
		var s entities.Sale
		if err := rows.Scan(&s.ID, &s.Total, &s.MetodoPago, &s.CreatedAt, &s.ClienteNombre); err != nil {
			return nil, err
		}
		sales = append(sales, s)
	}
	return sales, rows.Err()
}

// GetByID returns a sale with its items.
func (r *SaleRepo) GetByID(id int64) (*entities.Sale, error) {
	var s entities.Sale
	var clienteID sql.NullInt64
	err := r.db.QueryRow(`
		SELECT v.id, v.usuario_id, v.cliente_id, v.total, v.metodo_pago, v.created_at
		FROM ventas v WHERE v.id = ?
	`, id).Scan(&s.ID, &s.UsuarioID, &clienteID, &s.Total, &s.MetodoPago, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	if clienteID.Valid {
		s.ClienteID = &clienteID.Int64
	}

	// Get items
	rows, err := r.db.Query(`
		SELECT vi.id, vi.producto_id, COALESCE(p.nombre, ''), vi.cantidad, vi.precio_unitario, vi.subtotal
		FROM venta_items vi
		LEFT JOIN productos p ON p.id = vi.producto_id
		WHERE vi.venta_id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entities.SaleItem
		if err := rows.Scan(&item.ID, &item.ProductoID, &item.ProductoNombre, &item.Cantidad, &item.PrecioUnitario, &item.Subtotal); err != nil {
			return nil, err
		}
		s.Items = append(s.Items, item)
	}

	return &s, rows.Err()
}
