package ports

import "context"

// Product represents a product entity for port interfaces.
type Product struct {
	ID           int64
	Nombre       string
	SKU          string
	CategoriaID  *int64
	PrecioVenta  float64
	PrecioCompra float64
	StockActual  float64
	StockMinimo  float64
	Unidad       string
	Activo       bool
}

// ProductRepository defines the port for product persistence.
type ProductRepository interface {
	List(ctx context.Context, search string) ([]Product, error)
	GetByID(ctx context.Context, id int64) (*Product, error)
	Create(ctx context.Context, p *Product) (int64, error)
	Update(ctx context.Context, p *Product) error
	Search(ctx context.Context, q string) ([]Product, error)
	Deactivate(ctx context.Context, id int64) error
}

// Sale represents a sale entity for port interfaces.
type Sale struct {
	ID         int64
	UsuarioID  int64
	ClienteID  *int64
	Total      float64
	MetodoPago string
	Items      []SaleItem
}

// SaleItem represents a line item in a sale.
type SaleItem struct {
	ProductoID     int64
	Cantidad       float64
	PrecioUnitario float64
}

// SaleRepository defines the port for sale persistence.
type SaleRepository interface {
	Create(ctx context.Context, sale *Sale) (int64, error)
	GetByID(ctx context.Context, id int64) (*Sale, error)
	ListRecent(ctx context.Context, limit int) ([]Sale, error)
}

// User represents a user entity for port interfaces.
type User struct {
	ID      int64
	Nombre  string
	PinHash string
	Rol     string
	Activo  bool
}

// UserRepository defines the port for user persistence.
type UserRepository interface {
	GetByPinHash(ctx context.Context, pinHash string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}

// InventoryMovement represents a stock change.
type InventoryMovement struct {
	ProductoID      int64
	Tipo            string // entrada, salida, ajuste
	Cantidad        float64
	StockResultante float64
	ReferenciaTipo  string
	ReferenciaID    int64
	UsuarioID       int64
}

// InventoryRepository defines the port for inventory persistence.
type InventoryRepository interface {
	CreateMovement(ctx context.Context, m *InventoryMovement) error
	GetByProduct(ctx context.Context, productoID int64) ([]InventoryMovement, error)
}

// AIQueryService defines the port for NL→SQL generation.
type AIQueryService interface {
	GenerateSQL(ctx context.Context, question, systemPrompt string) (sql string, explanation string, err error)
}

// MetricsRepository defines the port for dashboard metrics.
type MetricsRepository interface {
	VentasHoy(ctx context.Context) (count int, total float64, err error)
	VentasSemana(ctx context.Context) (count int, total float64, err error)
	VentasMes(ctx context.Context) (count int, total float64, err error)
	TopProductos(ctx context.Context, limit int) ([]ProductMetric, error)
	StockBajo(ctx context.Context) ([]Product, error)
	ClientesFrecuentes(ctx context.Context, limit int) ([]ClientMetric, error)
}

// ProductMetric represents a product with sales metrics.
type ProductMetric struct {
	Nombre   string
	Unidades float64
}

// ClientMetric represents a customer with purchase metrics.
type ClientMetric struct {
	Nombre      string
	Compras     int
	TotalGastado float64
}
