package ports

import (
	"context"
	"time"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
)

// UserRepository defines the contract for user persistence and authentication state operations.
type UserRepository interface {
	// FindByID retrieves a user by their unique identifier.
	FindByID(ctx context.Context, id int64) (*entities.User, error)

	// FindByPINHash retrieves a user by their hashed PIN value.
	FindByPINHash(ctx context.Context, pinHash string) (*entities.User, error)

	// FindAll retrieves all users in the system.
	FindAll(ctx context.Context) ([]entities.User, error)

	// IncrementFailedAttempts increases the failed login attempt counter for a user.
	IncrementFailedAttempts(ctx context.Context, id int64) error

	// Lock sets a lockout period on a user account until the specified time.
	Lock(ctx context.Context, id int64, until time.Time) error

	// ResetAttempts clears the failed attempt counter for a user after successful login.
	ResetAttempts(ctx context.Context, id int64) error
}
