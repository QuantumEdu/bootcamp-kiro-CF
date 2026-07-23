package ports

import (
	"context"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
)

// ClientRepository defines the contract for client persistence operations.
type ClientRepository interface {
	// Create persists a new client in the store.
	Create(ctx context.Context, client *entities.Client) error

	// List retrieves all clients ordered by name.
	List(ctx context.Context) ([]entities.Client, error)
}
