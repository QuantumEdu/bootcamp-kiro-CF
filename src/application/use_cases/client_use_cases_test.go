package use_cases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/use_cases"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
)

// mockClientRepo is a simple mock implementing ports.ClientRepository.
type mockClientRepo struct {
	clients []entities.Client
	err     error
}

func (m *mockClientRepo) Create(_ context.Context, client *entities.Client) error {
	if m.err != nil {
		return m.err
	}
	client.ID = int64(len(m.clients) + 1)
	m.clients = append(m.clients, *client)
	return nil
}

func (m *mockClientRepo) List(_ context.Context) ([]entities.Client, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.clients, nil
}

func TestCreateClient_Execute_ValidInput(t *testing.T) {
	repo := &mockClientRepo{}
	uc := use_cases.NewCreateClient(repo)

	input := use_cases.CreateClientInput{
		Nombre:    "Juan Pérez",
		Telefono:  "555-1234",
		Direccion: "Calle 123",
	}

	client, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Nombre != "Juan Pérez" {
		t.Errorf("got Nombre=%q, want %q", client.Nombre, "Juan Pérez")
	}
	if client.Telefono != "555-1234" {
		t.Errorf("got Telefono=%q, want %q", client.Telefono, "555-1234")
	}
	if client.Direccion != "Calle 123" {
		t.Errorf("got Direccion=%q, want %q", client.Direccion, "Calle 123")
	}
	if client.ID == 0 {
		t.Error("expected client ID to be set after creation")
	}
}

func TestCreateClient_Execute_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   use_cases.CreateClientInput
		wantErr bool
	}{
		{"valid name", use_cases.CreateClientInput{Nombre: "María"}, false},
		{"empty name", use_cases.CreateClientInput{Nombre: ""}, true},
		{"whitespace only name", use_cases.CreateClientInput{Nombre: "   "}, true},
		{"tabs only name", use_cases.CreateClientInput{Nombre: "\t\t"}, true},
		{"name with spaces trimmed valid", use_cases.CreateClientInput{Nombre: "  Ana  "}, false},
		{"optional fields empty", use_cases.CreateClientInput{Nombre: "Carlos", Telefono: "", Direccion: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockClientRepo{}
			uc := use_cases.NewCreateClient(repo)
			_, err := uc.Execute(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("got err=%v, wantErr=%v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateClient_Execute_RepoError(t *testing.T) {
	repoErr := errors.New("database error")
	repo := &mockClientRepo{err: repoErr}
	uc := use_cases.NewCreateClient(repo)

	input := use_cases.CreateClientInput{Nombre: "Test"}
	_, err := uc.Execute(context.Background(), input)
	if err == nil {
		t.Fatal("expected error from repo, got nil")
	}
	if !errors.Is(err, repoErr) {
		t.Errorf("expected wrapped repo error, got: %v", err)
	}
}

func TestListClients_Execute_ReturnsClients(t *testing.T) {
	repo := &mockClientRepo{
		clients: []entities.Client{
			{ID: 1, Nombre: "Ana", Telefono: "111"},
			{ID: 2, Nombre: "Bruno", Telefono: "222"},
		},
	}
	uc := use_cases.NewListClients(repo)

	clients, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(clients) != 2 {
		t.Fatalf("got %d clients, want 2", len(clients))
	}
	if clients[0].Nombre != "Ana" {
		t.Errorf("got first client name=%q, want %q", clients[0].Nombre, "Ana")
	}
}

func TestListClients_Execute_EmptyList(t *testing.T) {
	repo := &mockClientRepo{}
	uc := use_cases.NewListClients(repo)

	clients, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(clients) != 0 {
		t.Errorf("got %d clients, want 0", len(clients))
	}
}

func TestListClients_Execute_RepoError(t *testing.T) {
	repoErr := errors.New("connection failed")
	repo := &mockClientRepo{err: repoErr}
	uc := use_cases.NewListClients(repo)

	_, err := uc.Execute(context.Background())
	if err == nil {
		t.Fatal("expected error from repo, got nil")
	}
	if !errors.Is(err, repoErr) {
		t.Errorf("expected wrapped repo error, got: %v", err)
	}
}
