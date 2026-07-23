package ports

import "context"

// TopProduct represents a product ranked by units sold.
type TopProduct struct {
	Nombre   string
	Unidades float64
}

// LowStockProduct represents a product with stock at or below its minimum threshold.
type LowStockProduct struct {
	Nombre      string
	StockActual float64
	StockMinimo float64
}

// FrequentClient represents a customer ranked by purchase frequency.
type FrequentClient struct {
	Nombre       string
	Compras      int
	TotalGastado float64
}

// MetricsRepository defines the contract for dashboard metrics queries.
type MetricsRepository interface {
	// VentasHoy returns the count and total amount of sales for today.
	VentasHoy(ctx context.Context) (count int, total float64, err error)

	// VentasSemana returns the count and total amount of sales for the current week.
	VentasSemana(ctx context.Context) (count int, total float64, err error)

	// VentasMes returns the count and total amount of sales for the current month.
	VentasMes(ctx context.Context) (count int, total float64, err error)

	// TopProductos returns the top-selling products ranked by units sold.
	TopProductos(ctx context.Context, limit int) ([]TopProduct, error)

	// StockBajo returns all products with stock at or below their minimum threshold.
	StockBajo(ctx context.Context) ([]LowStockProduct, error)

	// ClientesFrecuentes returns the most frequent customers ranked by purchase count.
	ClientesFrecuentes(ctx context.Context, limit int) ([]FrequentClient, error)
}
