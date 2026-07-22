package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
)

// ProductRepo handles product CRUD operations.
type ProductRepo struct {
	db *sql.DB
}

// NewProductRepo creates a new product repository.
func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

// List returns all active products with optional search.
func (r *ProductRepo) List(search string) ([]entities.Product, error) {
	query := `
		SELECT p.id, p.nombre, COALESCE(p.sku, ''), p.categoria_id, 
		       COALESCE(c.nombre, 'Sin categoria'), p.precio_venta, p.precio_compra,
		       p.stock_actual, p.stock_minimo, p.unidad, p.activo, p.created_at, p.updated_at
		FROM productos p
		LEFT JOIN categorias c ON c.id = p.categoria_id
		WHERE p.activo = 1
	`
	var args []interface{}

	if search != "" {
		query += ` AND (p.nombre LIKE ? OR p.sku LIKE ?)`
		s := "%" + search + "%"
		args = append(args, s, s)
	}

	query += ` ORDER BY p.nombre ASC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing products: %w", err)
	}
	defer rows.Close()

	var products []entities.Product
	for rows.Next() {
		var p entities.Product
		var catID sql.NullInt64
		if err := rows.Scan(&p.ID, &p.Nombre, &p.SKU, &catID, &p.CategoriaNombre,
			&p.PrecioVenta, &p.PrecioCompra, &p.StockActual, &p.StockMinimo,
			&p.Unidad, &p.Activo, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning product: %w", err)
		}
		if catID.Valid {
			p.CategoriaID = &catID.Int64
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// GetByID returns a product by ID.
func (r *ProductRepo) GetByID(id int64) (*entities.Product, error) {
	var p entities.Product
	var catID sql.NullInt64
	err := r.db.QueryRow(`
		SELECT p.id, p.nombre, COALESCE(p.sku, ''), p.categoria_id,
		       COALESCE(c.nombre, 'Sin categoria'), p.precio_venta, p.precio_compra,
		       p.stock_actual, p.stock_minimo, p.unidad, p.activo, p.created_at, p.updated_at
		FROM productos p
		LEFT JOIN categorias c ON c.id = p.categoria_id
		WHERE p.id = ?
	`, id).Scan(&p.ID, &p.Nombre, &p.SKU, &catID, &p.CategoriaNombre,
		&p.PrecioVenta, &p.PrecioCompra, &p.StockActual, &p.StockMinimo,
		&p.Unidad, &p.Activo, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting product %d: %w", id, err)
	}
	if catID.Valid {
		p.CategoriaID = &catID.Int64
	}
	return &p, nil
}

// Create inserts a new product.
func (r *ProductRepo) Create(p *entities.Product) (int64, error) {
	res, err := r.db.Exec(`
		INSERT INTO productos (nombre, sku, categoria_id, precio_venta, precio_compra, stock_actual, stock_minimo, unidad)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, p.Nombre, nullStr(p.SKU), p.CategoriaID, p.PrecioVenta, p.PrecioCompra, p.StockActual, p.StockMinimo, p.Unidad)
	if err != nil {
		return 0, fmt.Errorf("creating product: %w", err)
	}
	return res.LastInsertId()
}

// Update modifies an existing product.
func (r *ProductRepo) Update(p *entities.Product) error {
	_, err := r.db.Exec(`
		UPDATE productos 
		SET nombre = ?, sku = ?, categoria_id = ?, precio_venta = ?, precio_compra = ?, 
		    stock_actual = ?, stock_minimo = ?, unidad = ?, updated_at = datetime('now','localtime')
		WHERE id = ?
	`, p.Nombre, nullStr(p.SKU), p.CategoriaID, p.PrecioVenta, p.PrecioCompra,
		p.StockActual, p.StockMinimo, p.Unidad, p.ID)
	if err != nil {
		return fmt.Errorf("updating product %d: %w", p.ID, err)
	}
	return nil
}

// Search returns products matching query (for sale search).
func (r *ProductRepo) Search(q string) ([]entities.Product, error) {
	query := `
		SELECT id, nombre, COALESCE(sku, ''), precio_venta, stock_actual, unidad
		FROM productos
		WHERE activo = 1 AND (nombre LIKE ? OR sku LIKE ?)
		ORDER BY nombre ASC
		LIMIT 20
	`
	s := "%" + strings.TrimSpace(q) + "%"
	rows, err := r.db.Query(query, s, s)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []entities.Product
	for rows.Next() {
		var p entities.Product
		if err := rows.Scan(&p.ID, &p.Nombre, &p.SKU, &p.PrecioVenta, &p.StockActual, &p.Unidad); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// ListCategories returns all active categories.
func (r *ProductRepo) ListCategories() ([]entities.Category, error) {
	rows, err := r.db.Query(`SELECT id, nombre, COALESCE(descripcion, '') FROM categorias WHERE activo = 1 ORDER BY nombre`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []entities.Category
	for rows.Next() {
		var c entities.Category
		if err := rows.Scan(&c.ID, &c.Nombre, &c.Descripcion); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
