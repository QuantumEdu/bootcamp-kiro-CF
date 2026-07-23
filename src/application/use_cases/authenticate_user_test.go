package use_cases

import (
	"context"
	"testing"
	"time"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/value_objects"
)

// mockUserRepository is a simple mock implementing ports.UserRepository for testing.
type mockUserRepository struct {
	users          []entities.User
	findAllErr     error
	resetCalled    bool
	resetCalledFor int64
}

func (m *mockUserRepository) FindByID(_ context.Context, id int64) (*entities.User, error) {
	for i := range m.users {
		if m.users[i].ID == id {
			return &m.users[i], nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) FindByPINHash(_ context.Context, _ string) (*entities.User, error) {
	return nil, nil
}

func (m *mockUserRepository) FindAll(_ context.Context) ([]entities.User, error) {
	if m.findAllErr != nil {
		return nil, m.findAllErr
	}
	return m.users, nil
}

func (m *mockUserRepository) IncrementFailedAttempts(_ context.Context, _ int64) error {
	return nil
}

func (m *mockUserRepository) Lock(_ context.Context, _ int64, _ time.Time) error {
	return nil
}

func (m *mockUserRepository) ResetAttempts(_ context.Context, id int64) error {
	m.resetCalled = true
	m.resetCalledFor = id
	return nil
}

// helper to create a user with a hashed PIN.
func newTestUser(id int64, nombre string, pin string, rol entities.Role, activo bool) entities.User {
	hash, _ := value_objects.HashPIN(pin)
	return entities.User{
		ID:      id,
		Nombre:  nombre,
		PINHash: hash,
		Rol:     rol,
		Activo:  activo,
	}
}

func TestAuthenticateUser_Execute_ValidPIN(t *testing.T) {
	user := newTestUser(1, "Admin", "1234", entities.RoleAdmin, true)
	repo := &mockUserRepository{users: []entities.User{user}}
	uc := NewAuthenticateUser(repo)

	result, err := uc.Execute(context.Background(), "1234")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected user, got nil")
	}
	if result.ID != 1 {
		t.Errorf("expected user ID 1, got %d", result.ID)
	}
	if result.Nombre != "Admin" {
		t.Errorf("expected nombre Admin, got %s", result.Nombre)
	}
	if !repo.resetCalled {
		t.Error("expected ResetAttempts to be called")
	}
	if repo.resetCalledFor != 1 {
		t.Errorf("expected ResetAttempts called for user 1, got %d", repo.resetCalledFor)
	}
}

func TestAuthenticateUser_Execute_WrongPIN(t *testing.T) {
	user := newTestUser(1, "Admin", "1234", entities.RoleAdmin, true)
	repo := &mockUserRepository{users: []entities.User{user}}
	uc := NewAuthenticateUser(repo)

	result, err := uc.Execute(context.Background(), "9999")

	if result != nil {
		t.Fatalf("expected nil user, got %+v", result)
	}
	if err != ErrAuthPINInvalid {
		t.Errorf("expected ErrAuthPINInvalid, got %v", err)
	}
}

func TestAuthenticateUser_Execute_LockedUser(t *testing.T) {
	user := newTestUser(1, "Cajero", "5678", entities.RoleCajero, true)
	user.BloqueadoHasta = time.Now().Add(5 * time.Minute) // locked for 5 minutes
	repo := &mockUserRepository{users: []entities.User{user}}
	uc := NewAuthenticateUser(repo)

	result, err := uc.Execute(context.Background(), "5678")

	if result != nil {
		t.Fatalf("expected nil user, got %+v", result)
	}
	if err != ErrAuthAccountLocked {
		t.Errorf("expected ErrAuthAccountLocked, got %v", err)
	}
}

func TestAuthenticateUser_Execute_InvalidPINFormat(t *testing.T) {
	tests := []struct {
		name string
		pin  string
	}{
		{"too short", "12"},
		{"too long", "1234567"},
		{"non-numeric", "abcd"},
		{"empty", ""},
	}

	repo := &mockUserRepository{users: []entities.User{}}
	uc := NewAuthenticateUser(repo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := uc.Execute(context.Background(), tt.pin)
			if result != nil {
				t.Fatalf("expected nil user, got %+v", result)
			}
			if err != ErrAuthPINInvalid {
				t.Errorf("expected ErrAuthPINInvalid, got %v", err)
			}
		})
	}
}

func TestAuthenticateUser_Execute_InactiveUserSkipped(t *testing.T) {
	inactiveUser := newTestUser(1, "Inactive", "1234", entities.RoleCajero, false)
	activeUser := newTestUser(2, "Active", "5678", entities.RoleCajero, true)
	repo := &mockUserRepository{users: []entities.User{inactiveUser, activeUser}}
	uc := NewAuthenticateUser(repo)

	// The inactive user has PIN "1234" but should be skipped.
	result, err := uc.Execute(context.Background(), "1234")

	if result != nil {
		t.Fatalf("expected nil user (inactive skipped), got %+v", result)
	}
	if err != ErrAuthPINInvalid {
		t.Errorf("expected ErrAuthPINInvalid, got %v", err)
	}
}

func TestAuthenticateUser_Execute_MultipleUsers_CorrectMatch(t *testing.T) {
	user1 := newTestUser(1, "Admin", "1234", entities.RoleAdmin, true)
	user2 := newTestUser(2, "Cajero1", "5678", entities.RoleCajero, true)
	user3 := newTestUser(3, "Cajero2", "9012", entities.RoleCajero, true)
	repo := &mockUserRepository{users: []entities.User{user1, user2, user3}}
	uc := NewAuthenticateUser(repo)

	result, err := uc.Execute(context.Background(), "5678")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected user, got nil")
	}
	if result.ID != 2 {
		t.Errorf("expected user ID 2, got %d", result.ID)
	}
	if result.Nombre != "Cajero1" {
		t.Errorf("expected nombre Cajero1, got %s", result.Nombre)
	}
}

func TestAuthenticateUser_Execute_RepoError(t *testing.T) {
	repo := &mockUserRepository{
		findAllErr: context.DeadlineExceeded,
	}
	uc := NewAuthenticateUser(repo)

	result, err := uc.Execute(context.Background(), "1234")

	if result != nil {
		t.Fatalf("expected nil user, got %+v", result)
	}
	if err != ErrAuthPINInvalid {
		t.Errorf("expected ErrAuthPINInvalid (generic error), got %v", err)
	}
}
