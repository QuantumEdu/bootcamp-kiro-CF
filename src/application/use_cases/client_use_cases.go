package use_cases

import (
	"context"
	"fmt"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

// --- CreateClient ---

// CreateClientInput holds the data needed to create a new client.
type CreateClientInput struct {
	Nombre    string
	Telefono  string
	Direccion string
}

// CreateClient handles creating a new client in the system.
type CreateClient struct {
	repo ports.ClientRepository
}

// NewCreateClient creates a new CreateClient use case.
func NewCreateClient(repo ports.ClientRepository) *CreateClient {
	return &CreateClient{repo: repo}
}

// Execute validates the client entity and persists it.
func (uc *CreateClient) Execute(ctx context.Context, input CreateClientInput) (*entities.Client, error) {
	client := &entities.Client{
		Nombre:    input.Nombre,
		Telefono:  input.Telefono,
		Direccion: input.Direccion,
	}

	if err := client.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, client); err != nil {
		return nil, fmt.Errorf("creating client: %w", err)
	}

	return client, nil
}

// --- ListClients ---

// ListClients handles retrieving all clients.
type ListClients struct {
	repo ports.ClientRepository
}

// NewListClients creates a new ListClients use case.
func NewListClients(repo ports.ClientRepository) *ListClients {
	return &ListClients{repo: repo}
}

// Execute retrieves all clients ordered by name.
func (uc *ListClients) Execute(ctx context.Context) ([]entities.Client, error) {
	clients, err := uc.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing clients: %w", err)
	}

	return clients, nil
}
