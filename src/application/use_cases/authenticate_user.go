// Package use_cases contains application use cases for the POS system.
package use_cases

import (
	"context"
	"errors"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/value_objects"
)

// Domain errors for authentication.
var (
	ErrAuthPINInvalid    = errors.New("PIN incorrecto")
	ErrAuthAccountLocked = errors.New("Cuenta bloqueada temporalmente")
)

// AuthenticateUser handles PIN-based authentication for POS users.
type AuthenticateUser struct {
	repo ports.UserRepository
}

// NewAuthenticateUser creates a new AuthenticateUser use case.
func NewAuthenticateUser(repo ports.UserRepository) *AuthenticateUser {
	return &AuthenticateUser{repo: repo}
}

// Execute validates the PIN, checks lockout status, and returns the authenticated user.
// It iterates all active users and compares the PIN against each hash using bcrypt.
func (uc *AuthenticateUser) Execute(ctx context.Context, pin string) (*entities.User, error) {
	if err := value_objects.ValidatePINFormat(pin); err != nil {
		return nil, ErrAuthPINInvalid
	}

	users, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, ErrAuthPINInvalid
	}

	for i := range users {
		user := &users[i]
		if !user.Activo {
			continue
		}

		if err := value_objects.ComparePIN(user.PINHash, pin); err != nil {
			continue
		}

		// PIN matched — check lockout.
		if user.IsLocked() {
			return nil, ErrAuthAccountLocked
		}

		// Successful auth — reset failed attempts.
		_ = uc.repo.ResetAttempts(ctx, user.ID)

		return user, nil
	}

	// No match found — return generic error (don't reveal if user exists).
	return nil, ErrAuthPINInvalid
}
