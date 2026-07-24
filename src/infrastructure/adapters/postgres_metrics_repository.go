package adapters

import (
	"context"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Compile-time interface satisfaction check.
var _ ports.MetricsRepository = (*PostgresMetricsRepository)(nil)

// PostgresMetricsRepository implements ports.MetricsRepository using pgxpool.
type PostgresMetricsRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresMetricsRepository creates a new PostgresMetricsRepository.
func NewPostgresMetricsRepository(pool *pgxpool.Pool) *PostgresMetricsRepository {
	return &PostgresMetricsRepository{pool: pool}
}

// VentasHoy returns the count and total amount of sales for today.
func (r *PostgresMetricsRepository) VentasHoy(ctx context.Context) (int, float64, error) {
	var count int
	var total float64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*), COALESCE(SUM(total), 0)
		 FROM ventas
		 WHERE created_at >= CURRENT_DATE`).Scan(&count, &total)
	return count, total, wrapErr(err, "metrics ventas hoy", 0)
}

// VentasSemana returns the count and total amount of sales for the current week.
func (r *PostgresMetricsRepository) VentasSemana(ctx context.Context) (int, float64, error) {
	var count int
	var total float64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*), COALESCE(SUM(total), 0)
		 FROM ventas
		 WHERE created_at >= date_trunc('week', CURRENT_DATE)`).Scan(&count, &total)
	return count, total, wrapErr(err, "metrics ventas semana", 0)
}

// VentasMes returns the count and total amount of sales for the current month.
func (r *PostgresMetricsRepository) VentasMes(ctx context.Context) (int, float64, error) {
	var count int
	var total float64
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*), COALESCE(SUM(total), 0)
		 FROM ventas
		 WHERE created_at >= date_trunc('month', CURRENT_DATE)`).Scan(&count, &total)
	return count, total, wrapErr(err, "metrics ventas mes", 0)
}

// TopProductos returns the top-selling products ranked by units sold in the last 30 days.
func (r *PostgresMetricsRepository) TopProductos(ctx context.Context, limit int) ([]ports.TopProduct, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT p.nombre, SUM(vi.cantidad) AS unidades
		 FROM venta_items vi
		 JOIN productos p ON p.id = vi.producto_id
		 JOIN ventas v ON v.id = vi.venta_id
		 WHERE v.created_at >= CURRENT_DATE - INTERVAL '30 days'
		 GROUP BY p.nombre
		 ORDER BY unidades DESC
		 LIMIT $1`, limit)
	if err != nil {
		return nil, wrapErr(err, "metrics top productos", 0)
	}
	defer rows.Close()

	var results []ports.TopProduct
	for rows.Next() {
		var tp ports.TopProduct
		if err := rows.Scan(&tp.Nombre, &tp.Unidades); err != nil {
			return nil, wrapErr(err, "metrics top productos scan", 0)
		}
		results = append(results, tp)
	}
	if err := rows.Err(); err != nil {
		return nil, wrapErr(err, "metrics top productos rows", 0)
	}
	return results, nil
}

// StockBajo returns all active products with stock at or below their minimum threshold.
func (r *PostgresMetricsRepository) StockBajo(ctx context.Context) ([]ports.LowStockProduct, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT nombre, stock_actual, stock_minimo
		 FROM productos
		 WHERE stock_actual <= stock_minimo AND activo = true
		 ORDER BY stock_actual ASC
		 LIMIT 10`)
	if err != nil {
		return nil, wrapErr(err, "metrics stock bajo", 0)
	}
	defer rows.Close()

	var results []ports.LowStockProduct
	for rows.Next() {
		var lsp ports.LowStockProduct
		if err := rows.Scan(&lsp.Nombre, &lsp.StockActual, &lsp.StockMinimo); err != nil {
			return nil, wrapErr(err, "metrics stock bajo scan", 0)
		}
		results = append(results, lsp)
	}
	if err := rows.Err(); err != nil {
		return nil, wrapErr(err, "metrics stock bajo rows", 0)
	}
	return results, nil
}

// ClientesFrecuentes returns the most frequent customers ranked by purchase count in the last 30 days.
func (r *PostgresMetricsRepository) ClientesFrecuentes(ctx context.Context, limit int) ([]ports.FrequentClient, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT c.nombre, COUNT(v.id) AS compras, COALESCE(SUM(v.total), 0) AS total_gastado
		 FROM clientes c
		 JOIN ventas v ON v.cliente_id = c.id
		 WHERE v.created_at >= CURRENT_DATE - INTERVAL '30 days'
		 GROUP BY c.nombre
		 ORDER BY compras DESC
		 LIMIT $1`, limit)
	if err != nil {
		return nil, wrapErr(err, "metrics clientes frecuentes", 0)
	}
	defer rows.Close()

	var results []ports.FrequentClient
	for rows.Next() {
		var fc ports.FrequentClient
		if err := rows.Scan(&fc.Nombre, &fc.Compras, &fc.TotalGastado); err != nil {
			return nil, wrapErr(err, "metrics clientes frecuentes scan", 0)
		}
		results = append(results, fc)
	}
	if err := rows.Err(); err != nil {
		return nil, wrapErr(err, "metrics clientes frecuentes rows", 0)
	}
	return results, nil
}
