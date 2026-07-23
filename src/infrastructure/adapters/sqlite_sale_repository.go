package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

// Compile-time interface satisfaction check.
var _ ports.SaleRepository = (*SQLiteSaleRepository)(nil)

// SQLiteSaleRepository implements ports.SaleRepository using SQLite.
type SQLiteSaleRepository struct {
	db *sql.DB
}

// NewSQLiteSaleRepository creates a new SQLiteSaleRepository.
func NewSQLiteSaleRepository(db *sql.DB) *SQLiteSaleRepository {
	return &SQLiteSaleRepository{db: db}
}

// Create persists a new sale with its items in a transaction.
func (r *SQLiteSaleRepository) Create(ctx context.Context, sale *entities.Sale) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert sale header.
	result, err := tx.ExecContext(ctx,
		`INSERT INTO ventas (usuario_id, cliente_id, total, metodo_pago)
		 VALUES (?, ?, ?, ?)`,
		sale.UsuarioID,
		nullableInt64Ptr(sale.ClienteID),
		sale.Total,
		string(sale.MetodoPago),
	)
	if err != nil {
		return fmt.Errorf("inserting sale: %w", err)
	}

	saleID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting sale id: %w", err)
	}
	sale.ID = saleID

	// Insert sale items.
	for i := range sale.Items {
		sale.Items[i].VentaID = saleID
		itemResult, err := tx.ExecContext(ctx,
			`INSERT INTO venta_items (venta_id, producto_id, cantidad, precio_unitario, subtotal)
			 VALUES (?, ?, ?, ?, ?)`,
			sale.Items[i].VentaID,
			sale.Items[i].ProductoID,
			sale.Items[i].Cantidad,
			sale.Items[i].PrecioUnitario,
			sale.Items[i].Subtotal,
		)
		if err != nil {
			return fmt.Errorf("inserting sale item: %w", err)
		}

		itemID, err := itemResult.LastInsertId()
		if err != nil {
			return fmt.Errorf("getting item id: %w", err)
		}
		sale.Items[i].ID = itemID
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// FindByID retrieves a sale by its unique identifier, including items.
func (r *SQLiteSaleRepository) FindByID(ctx context.Context, id int64) (*entities.Sale, error) {
	// Fetch sale header.
	row := r.db.QueryRowContext(ctx,
		`SELECT id, usuario_id, cliente_id, total, metodo_pago, created_at
		 FROM ventas WHERE id = ?`, id)

	sale, err := r.scanSale(row)
	if err != nil {
		return nil, fmt.Errorf("finding sale: %w", err)
	}

	// Fetch sale items.
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, venta_id, producto_id, cantidad, precio_unitario, subtotal
		 FROM venta_items WHERE venta_id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("finding sale items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		item, err := r.scanSaleItem(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning sale item: %w", err)
		}
		sale.Items = append(sale.Items, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating sale items: %w", err)
	}

	return sale, nil
}

// List retrieves sales matching the given filter criteria.
func (r *SQLiteSaleRepository) List(ctx context.Context, filter ports.SaleFilter) ([]entities.Sale, error) {
	var conditions []string
	var args []interface{}

	if filter.UsuarioID != nil {
		conditions = append(conditions, "usuario_id = ?")
		args = append(args, *filter.UsuarioID)
	}

	if filter.Since != nil {
		conditions = append(conditions, "created_at >= ?")
		args = append(args, filter.Since.Format("2006-01-02 15:04:05"))
	}

	if filter.Until != nil {
		conditions = append(conditions, "created_at <= ?")
		args = append(args, filter.Until.Format("2006-01-02 15:04:05"))
	}

	query := `SELECT id, usuario_id, cliente_id, total, metodo_pago, created_at FROM ventas`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing sales: %w", err)
	}
	defer rows.Close()

	var sales []entities.Sale
	for rows.Next() {
		sale, err := r.scanSaleFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning sale: %w", err)
		}
		sales = append(sales, *sale)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating sales: %w", err)
	}

	return sales, nil
}

// scanSale scans a single sale from a *sql.Row.
func (r *SQLiteSaleRepository) scanSale(row *sql.Row) (*entities.Sale, error) {
	var sale entities.Sale
	var clienteID sql.NullInt64
	var metodoPago string
	var createdAt string

	err := row.Scan(
		&sale.ID,
		&sale.UsuarioID,
		&clienteID,
		&sale.Total,
		&metodoPago,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	if clienteID.Valid {
		sale.ClienteID = &clienteID.Int64
	}
	sale.MetodoPago = entities.PaymentMethod(metodoPago)

	if t, err := parseDateTime(createdAt); err == nil {
		sale.CreatedAt = t
	}

	return &sale, nil
}

// scanSaleFromRows scans a single sale from *sql.Rows.
func (r *SQLiteSaleRepository) scanSaleFromRows(rows *sql.Rows) (*entities.Sale, error) {
	var sale entities.Sale
	var clienteID sql.NullInt64
	var metodoPago string
	var createdAt string

	err := rows.Scan(
		&sale.ID,
		&sale.UsuarioID,
		&clienteID,
		&sale.Total,
		&metodoPago,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	if clienteID.Valid {
		sale.ClienteID = &clienteID.Int64
	}
	sale.MetodoPago = entities.PaymentMethod(metodoPago)

	if t, err := parseDateTime(createdAt); err == nil {
		sale.CreatedAt = t
	}

	return &sale, nil
}

// scanSaleItem scans a single sale item from *sql.Rows.
func (r *SQLiteSaleRepository) scanSaleItem(rows *sql.Rows) (*entities.SaleItem, error) {
	var item entities.SaleItem

	err := rows.Scan(
		&item.ID,
		&item.VentaID,
		&item.ProductoID,
		&item.Cantidad,
		&item.PrecioUnitario,
		&item.Subtotal,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// nullableInt64Ptr returns nil for nil pointers or zero values.
func nullableInt64Ptr(n *int64) interface{} {
	if n == nil || *n == 0 {
		return nil
	}
	return *n
}
