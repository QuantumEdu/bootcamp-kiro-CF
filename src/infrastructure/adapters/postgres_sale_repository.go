package adapters

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

// Compile-time interface satisfaction check.
var _ ports.SaleRepository = (*PostgresSaleRepository)(nil)

// PostgresSaleRepository implements ports.SaleRepository using pgxpool.
type PostgresSaleRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresSaleRepository creates a new PostgresSaleRepository.
func NewPostgresSaleRepository(pool *pgxpool.Pool) *PostgresSaleRepository {
	return &PostgresSaleRepository{pool: pool}
}

// Create persists a new sale with its items atomically in a transaction.
func (r *PostgresSaleRepository) Create(ctx context.Context, sale *entities.Sale) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning sale transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert sale header with RETURNING to get ID and created_at.
	err = tx.QueryRow(ctx,
		`INSERT INTO ventas (usuario_id, cliente_id, total, metodo_pago)
		 VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		sale.UsuarioID,
		sale.ClienteID,
		sale.Total,
		string(sale.MetodoPago),
	).Scan(&sale.ID, &sale.CreatedAt)
	if err != nil {
		return wrapErr(fmt.Errorf("inserting sale: %w", err), "creating sale", 0)
	}

	// Insert sale items.
	for i := range sale.Items {
		item := &sale.Items[i]
		item.VentaID = sale.ID
		err = tx.QueryRow(ctx,
			`INSERT INTO venta_items (venta_id, producto_id, cantidad, precio_unitario, subtotal)
			 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			item.VentaID,
			item.ProductoID,
			item.Cantidad,
			item.PrecioUnitario,
			item.Subtotal,
		).Scan(&item.ID)
		if err != nil {
			return wrapErr(fmt.Errorf("inserting sale item %d: %w", i, err), "creating sale item", sale.ID)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing sale transaction: %w", err)
	}

	return nil
}

// FindByID retrieves a sale by its unique identifier, including items.
func (r *PostgresSaleRepository) FindByID(ctx context.Context, id int64) (*entities.Sale, error) {
	// Fetch sale header.
	row := r.pool.QueryRow(ctx,
		`SELECT id, usuario_id, cliente_id, total, metodo_pago, created_at
		 FROM ventas WHERE id = $1`, id)

	sale, err := r.scanSaleRow(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("sale not found (id=%d)", id)
		}
		return nil, wrapErr(err, "finding sale", id)
	}

	// Fetch sale items.
	rows, err := r.pool.Query(ctx,
		`SELECT id, venta_id, producto_id, cantidad, precio_unitario, subtotal
		 FROM venta_items WHERE venta_id = $1`, id)
	if err != nil {
		return nil, wrapErr(err, "finding sale items", id)
	}
	defer rows.Close()

	for rows.Next() {
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
			return nil, wrapErr(err, "scanning sale item", id)
		}
		sale.Items = append(sale.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, wrapErr(err, "iterating sale items", id)
	}

	return sale, nil
}

// List retrieves sales matching the given filter criteria.
func (r *PostgresSaleRepository) List(ctx context.Context, filter ports.SaleFilter) ([]entities.Sale, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.UsuarioID != nil {
		conditions = append(conditions, fmt.Sprintf("usuario_id = $%d", argIdx))
		args = append(args, *filter.UsuarioID)
		argIdx++
	}

	if filter.Since != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIdx))
		args = append(args, *filter.Since)
		argIdx++
	}

	if filter.Until != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIdx))
		args = append(args, *filter.Until)
		argIdx++
	}

	query := `SELECT id, usuario_id, cliente_id, total, metodo_pago, created_at FROM ventas`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, wrapErr(err, "listing sales", 0)
	}
	defer rows.Close()

	var sales []entities.Sale
	for rows.Next() {
		var sale entities.Sale
		var clienteID *int64
		var metodoPago string

		err := rows.Scan(
			&sale.ID,
			&sale.UsuarioID,
			&clienteID,
			&sale.Total,
			&metodoPago,
			&sale.CreatedAt,
		)
		if err != nil {
			return nil, wrapErr(err, "scanning sale", 0)
		}

		sale.ClienteID = clienteID
		sale.MetodoPago = entities.PaymentMethod(metodoPago)
		sales = append(sales, sale)
	}

	if err := rows.Err(); err != nil {
		return nil, wrapErr(err, "iterating sales", 0)
	}

	return sales, nil
}

// scanSaleRow scans a single sale from a pgx.Row.
func (r *PostgresSaleRepository) scanSaleRow(row pgx.Row) (*entities.Sale, error) {
	var sale entities.Sale
	var clienteID *int64
	var metodoPago string

	err := row.Scan(
		&sale.ID,
		&sale.UsuarioID,
		&clienteID,
		&sale.Total,
		&metodoPago,
		&sale.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	sale.ClienteID = clienteID
	sale.MetodoPago = entities.PaymentMethod(metodoPago)

	return &sale, nil
}
