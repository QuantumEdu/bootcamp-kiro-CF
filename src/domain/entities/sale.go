package entities

// Sale represents a sales transaction header.
type Sale struct {
	ID          int64      `json:"id"`
	UsuarioID   int64      `json:"usuario_id"`
	ClienteID   *int64     `json:"cliente_id"`
	Total       float64    `json:"total"`
	MetodoPago  string     `json:"metodo_pago"`
	CreatedAt   string     `json:"created_at"`
	Items       []SaleItem `json:"items,omitempty"`
	ClienteNombre string  `json:"cliente_nombre,omitempty"`
}

// SaleItem represents a line item in a sale.
type SaleItem struct {
	ID             int64   `json:"id"`
	VentaID        int64   `json:"venta_id"`
	ProductoID     int64   `json:"producto_id"`
	ProductoNombre string  `json:"producto_nombre,omitempty"`
	Cantidad       float64 `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Subtotal       float64 `json:"subtotal"`
}

// CreateSaleRequest is the input for creating a new sale.
type CreateSaleRequest struct {
	Items      []CreateSaleItem `json:"items"`
	MetodoPago string           `json:"metodo_pago"`
	ClienteID  *int64           `json:"cliente_id"`
}

// CreateSaleItem is an item in the create sale request.
type CreateSaleItem struct {
	ProductoID     int64   `json:"producto_id"`
	Cantidad       float64 `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
}
