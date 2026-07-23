package ports

import (
	"context"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
)

// ProductFilter defines the criteria for filtering products in list queries.
type ProductFilter struct {
	CategoriaID *int64
	Activo      *bool
	Search      string
}

// ProductRepository defines the contract for product persistence operations.
type ProductRepository interface {
	// Create persists a new product in the store.
	Create(ctx context.Context, product *entities.Product) error

	// Update modifies an existing product in the store.
	Update(ctx context.Context, product *entities.Product) error

	// FindByID retrieves a product by its unique identifier.
	FindByID(ctx context.Context, id int64) (*entities.Product, error)

	// List retrieves products matching the given filter criteria.
	List(ctx context.Context, filter ProductFilter) ([]entities.Product, error)

	// Deactivate marks a product as inactive (soft delete).
	Deactivate(ctx context.Context, id int64) error

	// FindLowStock retrieves all products whose current stock is at or below the minimum threshold.
	FindLowStock(ctx context.Context) ([]entities.Product, error)
}
