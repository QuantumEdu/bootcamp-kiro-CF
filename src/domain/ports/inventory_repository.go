package ports

import (
	"context"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
)

// InventoryRepository defines the contract for inventory movement persistence operations.
type InventoryRepository interface {
	// Create persists a new inventory movement in the store.
	Create(ctx context.Context, movement *entities.InventoryMovement) error

	// FindByProduct retrieves all inventory movements for a given product.
	FindByProduct(ctx context.Context, productID int64) ([]entities.InventoryMovement, error)
}
