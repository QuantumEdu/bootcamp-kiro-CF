package adapters

import (
	"context"
	"fmt"
	"strings"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Compile-time interface satisfaction check.
var _ ports.ProductRepository = (*PostgresProductRepository)(nil)

// PostgresProductRepository implements ports.ProductRepository using pgxpool.
type PostgresProductRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresProductRepository creates a new PostgresProductRepository.
func NewPostgresProductRepository(pool *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{pool: pool}
}

// Create persists a new product in the store.
func (r *PostgresProductRepository) Create(ctx context.Context, product *entities.Product) error {
	const q = `INSERT INTO productos (nombre, sku, categoria_id, precio_venta, precio_compra, stock_actual, stock_minimo, unidad, activo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, q,
		product.Nombre,
		nullableString(product.SKU),
		nullableInt64(product.CategoriaID),
		product.PrecioVenta,
		product.PrecioCompra,
		product.StockActual,
		product.StockMinimo,
		string(product.Unidad),
		product.Activo,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	return wrapErr(err, "creating product", 0)
}

// Update modifies an existing product in the store.
func (r *PostgresProductRepository) Update(ctx context.Context, product *entities.Product) error {
	const q = `UPDATE productos SET nombre=$1, sku=$2, categoria_id=$3,
		precio_venta=$4, precio_compra=$5, stock_actual=$6, stock_minimo=$7,
		unidad=$8, activo=$9, updated_at=NOW()
		WHERE id=$10`

	tag, err := r.pool.Exec(ctx, q,
		product.Nombre,
		nullableString(product.SKU),
		nullableInt64(product.CategoriaID),
		product.PrecioVenta,
		product.PrecioCompra,
		product.StockActual,
		product.StockMinimo,
		string(product.Unidad),
		product.Activo,
		product.ID,
	)
	if err != nil {
		return wrapErr(err, "updating product", product.ID)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("updating product (id=%d): %w", product.ID, pgx.ErrNoRows)
	}

	return nil
}

// FindByID retrieves a product by its unique identifier.
func (r *PostgresProductRepository) FindByID(ctx context.Context, id int64) (*entities.Product, error) {
	const q = `SELECT id, nombre, sku, categoria_id, precio_venta, precio_compra,
		stock_actual, stock_minimo, unidad, activo, created_at, updated_at
		FROM productos WHERE id = $1`

	row := r.pool.QueryRow(ctx, q, id)
	product, err := r.scanProduct(row)
	if err != nil {
		return nil, wrapErr(err, "finding product", id)
	}
	return product, nil
}

// List retrieves products matching the given filter criteria.
func (r *PostgresProductRepository) List(ctx context.Context, filter ports.ProductFilter) ([]entities.Product, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filter.CategoriaID != nil {
		conditions = append(conditions, fmt.Sprintf("categoria_id = $%d", argIdx))
		args = append(args, *filter.CategoriaID)
		argIdx++
	}

	if filter.Activo != nil {
		conditions = append(conditions, fmt.Sprintf("activo = $%d", argIdx))
		args = append(args, *filter.Activo)
		argIdx++
	}

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(nombre ILIKE $%d OR sku ILIKE $%d)", argIdx, argIdx+1))
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argIdx += 2
	}

	query := `SELECT id, nombre, sku, categoria_id, precio_venta, precio_compra,
		stock_actual, stock_minimo, unidad, activo, created_at, updated_at
		FROM productos`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY nombre"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, wrapErr(err, "listing products", 0)
	}
	defer rows.Close()

	var products []entities.Product
	for rows.Next() {
		product, err := r.scanProductFromRows(rows)
		if err != nil {
			return nil, wrapErr(err, "scanning product row", 0)
		}
		products = append(products, *product)
	}

	if err := rows.Err(); err != nil {
		return nil, wrapErr(err, "iterating product rows", 0)
	}

	return products, nil
}

// Deactivate marks a product as inactive (soft delete).
func (r *PostgresProductRepository) Deactivate(ctx context.Context, id int64) error {
	const q = `UPDATE productos SET activo = false, updated_at = NOW() WHERE id = $1`

	tag, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return wrapErr(err, "deactivating product", id)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("deactivating product (id=%d): %w", id, pgx.ErrNoRows)
	}

	return nil
}

// FindLowStock retrieves all products whose current stock is at or below the minimum threshold.
func (r *PostgresProductRepository) FindLowStock(ctx context.Context) ([]entities.Product, error) {
	const q = `SELECT id, nombre, sku, categoria_id, precio_venta, precio_compra,
		stock_actual, stock_minimo, unidad, activo, created_at, updated_at
		FROM productos
		WHERE stock_actual <= stock_minimo AND activo = true
		ORDER BY (stock_minimo - stock_actual) DESC`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, wrapErr(err, "finding low stock products", 0)
	}
	defer rows.Close()

	var products []entities.Product
	for rows.Next() {
		product, err := r.scanProductFromRows(rows)
		if err != nil {
			return nil, wrapErr(err, "scanning low stock product row", 0)
		}
		products = append(products, *product)
	}

	if err := rows.Err(); err != nil {
		return nil, wrapErr(err, "iterating low stock rows", 0)
	}

	return products, nil
}

// scanProduct scans a single product from a pgx.Row.
func (r *PostgresProductRepository) scanProduct(row pgx.Row) (*entities.Product, error) {
	var product entities.Product
	var sku *string
	var categoriaID *int64
	var unidad string

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
		&product.Activo,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if sku != nil {
		product.SKU = *sku
	}
	if categoriaID != nil {
		product.CategoriaID = *categoriaID
	}
	product.Unidad = entities.Unit(unidad)

	return &product, nil
}

// scanProductFromRows scans a single product from pgx.Rows.
func (r *PostgresProductRepository) scanProductFromRows(rows pgx.Rows) (*entities.Product, error) {
	var product entities.Product
	var sku *string
	var categoriaID *int64
	var unidad string

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
		&product.Activo,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if sku != nil {
		product.SKU = *sku
	}
	if categoriaID != nil {
		product.CategoriaID = *categoriaID
	}
	product.Unidad = entities.Unit(unidad)

	return &product, nil
}
