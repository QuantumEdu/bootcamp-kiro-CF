package ports

import (
	"context"
	"time"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
)

// SaleFilter defines the criteria for filtering sales in list queries.
type SaleFilter struct {
	UsuarioID *int64
	Since     *time.Time
	Until     *time.Time
	Limit     int
}

// SaleRepository defines the contract for sale persistence operations.
type SaleRepository interface {
	// Create persists a new sale with its items in the store.
	Create(ctx context.Context, sale *entities.Sale) error

	// FindByID retrieves a sale by its unique identifier, including items.
	FindByID(ctx context.Context, id int64) (*entities.Sale, error)

	// List retrieves sales matching the given filter criteria.
	List(ctx context.Context, filter SaleFilter) ([]entities.Sale, error)
}
