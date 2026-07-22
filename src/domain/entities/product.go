package entities

// Product represents a product in the catalog.
type Product struct {
	ID             int64   `json:"id"`
	Nombre         string  `json:"nombre"`
	SKU            string  `json:"sku"`
	CategoriaID    *int64  `json:"categoria_id"`
	CategoriaNombre string `json:"categoria_nombre,omitempty"`
	PrecioVenta    float64 `json:"precio_venta"`
	PrecioCompra   float64 `json:"precio_compra"`
	StockActual    float64 `json:"stock_actual"`
	StockMinimo    float64 `json:"stock_minimo"`
	Unidad         string  `json:"unidad"`
	Activo         bool    `json:"activo"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// Category represents a product category.
type Category struct {
	ID          int64  `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Activo      bool   `json:"activo"`
}
