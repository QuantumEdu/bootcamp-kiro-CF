package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

// Compile-time interface satisfaction check.
var _ ports.ProductRepository = (*SQLiteProductRepository)(nil)

// SQLiteProductRepository implements ports.ProductRepository using SQLite.
type SQLiteProductRepository struct {
	db *sql.DB
}

// NewSQLiteProductRepository creates a new SQLiteProductRepository.
func NewSQLiteProductRepository(db *sql.DB) *SQLiteProductRepository {
	return &SQLiteProductRepository{db: db}
}

// Create persists a new product in the store.
func (r *SQLiteProductRepository) Create(ctx context.Context, product *entities.Product) error {
	query := `INSERT INTO productos (nombre, sku, categoria_id, precio_venta, precio_compra, stock_actual, stock_minimo, unidad, activo)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query,
		product.Nombre,
		nullableString(product.SKU),
		nullableInt64(product.CategoriaID),
		product.PrecioVenta,
		product.PrecioCompra,
		product.StockActual,
		product.StockMinimo,
		string(product.Unidad),
		boolToInt(product.Activo),
	)
	if err != nil {
		return fmt.Errorf("inserting product: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last insert id: %w", err)
	}

	product.ID = id
	return nil
}

// Update modifies an existing product in the store.
func (r *SQLiteProductRepository) Update(ctx context.Context, product *entities.Product) error {
	query := `UPDATE productos
		SET nombre = ?, sku = ?, categoria_id = ?, precio_venta = ?, precio_compra = ?,
		    stock_actual = ?, stock_minimo = ?, unidad = ?, activo = ?, updated_at = datetime('now','localtime')
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		product.Nombre,
		nullableString(product.SKU),
		nullableInt64(product.CategoriaID),
		product.PrecioVenta,
		product.PrecioCompra,
		product.StockActual,
		product.StockMinimo,
		string(product.Unidad),
		boolToInt(product.Activo),
		product.ID,
	)
	if err != nil {
		return fmt.Errorf("updating product: %w", err)
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

// FindByID retrieves a product by its unique identifier.
func (r *SQLiteProductRepository) FindByID(ctx context.Context, id int64) (*entities.Product, error) {
	query := `SELECT id, nombre, sku, categoria_id, precio_venta, precio_compra,
		stock_actual, stock_minimo, unidad, activo, created_at, updated_at
		FROM productos WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanProduct(row)
}

// List retrieves products matching the given filter criteria.
func (r *SQLiteProductRepository) List(ctx context.Context, filter ports.ProductFilter) ([]entities.Product, error) {
	var conditions []string
	var args []interface{}

	if filter.CategoriaID != nil {
		conditions = append(conditions, "categoria_id = ?")
		args = append(args, *filter.CategoriaID)
	}

	if filter.Activo != nil {
		conditions = append(conditions, "activo = ?")
		args = append(args, boolToInt(*filter.Activo))
	}

	if filter.Search != "" {
		conditions = append(conditions, "(nombre LIKE ? OR sku LIKE ?)")
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	query := `SELECT id, nombre, sku, categoria_id, precio_venta, precio_compra,
		stock_actual, stock_minimo, unidad, activo, created_at, updated_at
		FROM productos`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY nombre"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing products: %w", err)
	}
	defer rows.Close()

	var products []entities.Product
	for rows.Next() {
		product, err := r.scanProductFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning product row: %w", err)
		}
		products = append(products, *product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating product rows: %w", err)
	}

	return products, nil
}

// Deactivate marks a product as inactive (soft delete).
func (r *SQLiteProductRepository) Deactivate(ctx context.Context, id int64) error {
	query := `UPDATE productos SET activo = 0, updated_at = datetime('now','localtime') WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deactivating product: %w", err)
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

// FindLowStock retrieves all products whose current stock is at or below the minimum threshold.
func (r *SQLiteProductRepository) FindLowStock(ctx context.Context) ([]entities.Product, error) {
	query := `SELECT id, nombre, sku, categoria_id, precio_venta, precio_compra,
		stock_actual, stock_minimo, unidad, activo, created_at, updated_at
		FROM productos
		WHERE stock_actual <= stock_minimo AND activo = 1
		ORDER BY (stock_minimo - stock_actual) DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("finding low stock products: %w", err)
	}
	defer rows.Close()

	var products []entities.Product
	for rows.Next() {
		product, err := r.scanProductFromRows(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning product row: %w", err)
		}
		products = append(products, *product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating product rows: %w", err)
	}

	return products, nil
}

// scanProduct scans a single product from a *sql.Row.
func (r *SQLiteProductRepository) scanProduct(row *sql.Row) (*entities.Product, error) {
	var product entities.Product
	var sku sql.NullString
	var categoriaID sql.NullInt64
	var activo int
	var unidad string
	var createdAt, updatedAt string

	err := row.Scan(
		&product.ID,
		&product.Nombre,
		&sku,
		&categoriaID,
		&product.PrecioVenta,
		&product.PrecioCompra,
		&product.StockActual,
		&product.StockMinimo,
		&unidad,
		&activo,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	product.SKU = sku.String
	if categoriaID.Valid {
		product.CategoriaID = categoriaID.Int64
	}
	product.Activo = activo == 1
	product.Unidad = entities.Unit(unidad)

	if t, err := parseDateTime(createdAt); err == nil {
		product.CreatedAt = t
	}
	if t, err := parseDateTime(updatedAt); err == nil {
		product.UpdatedAt = t
	}

	return &product, nil
}

// scanProductFromRows scans a single product from *sql.Rows.
func (r *SQLiteProductRepository) scanProductFromRows(rows *sql.Rows) (*entities.Product, error) {
	var product entities.Product
	var sku sql.NullString
	var categoriaID sql.NullInt64
	var activo int
	var unidad string
	var createdAt, updatedAt string

	err := rows.Scan(
		&product.ID,
		&product.Nombre,
		&sku,
		&categoriaID,
		&product.PrecioVenta,
		&product.PrecioCompra,
		&product.StockActual,
		&product.StockMinimo,
		&unidad,
		&activo,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning product: %w", err)
	}

	product.SKU = sku.String
	if categoriaID.Valid {
		product.CategoriaID = categoriaID.Int64
	}
	product.Activo = activo == 1
	product.Unidad = entities.Unit(unidad)

	if t, err := parseDateTime(createdAt); err == nil {
		product.CreatedAt = t
	}
	if t, err := parseDateTime(updatedAt); err == nil {
		product.UpdatedAt = t
	}

	return &product, nil
}

// --- Helper functions ---

// parseDateTime tries common SQLite datetime formats.
func parseDateTime(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05Z",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse datetime: %s", s)
}

// nullableString returns a sql.NullString for empty strings.
func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// nullableInt64 returns nil for zero values (no category).
func nullableInt64(n int64) interface{} {
	if n == 0 {
		return nil
	}
	return n
}

// boolToInt converts a boolean to SQLite integer (0/1).
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
