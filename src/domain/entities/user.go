package entities

import (
	"errors"
	"time"
)

// Domain errors for User validation.
var (
	ErrUserNameEmpty    = errors.New("user name cannot be empty")
	ErrUserRoleInvalid = errors.New("user role must be admin or cajero")
)

// MaxAttempts is the maximum number of failed PIN attempts before lockout.
const MaxAttempts = 5

// Role represents the access level of a user in the POS system.
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleCajero Role = "cajero"
)

// ValidRoles contains all valid role values for lookup.
var ValidRoles = map[Role]bool{
	RoleAdmin:  true,
	RoleCajero: true,
}

// User represents a system operator in the POS system.
// Maps to the "usuarios" table in the database.
type User struct {
	ID               int64
	Nombre           string
	PINHash          string
	Rol              Role
	Activo           bool
	IntentosFallidos int
	BloqueadoHasta   time.Time
	CreatedAt        time.Time
}

// Validate checks that the User satisfies all domain invariants.
// Returns the first validation error found using guard clauses.
func (u *User) Validate() error {
	if u.Nombre == "" {
		return ErrUserNameEmpty
	}

	if !ValidRoles[u.Rol] {
		return ErrUserRoleInvalid
	}

	return nil
}

// IsLocked returns true when the user account is currently locked due to
// excessive failed PIN attempts. A user is locked if BloqueadoHasta is
// after the current time.
func (u *User) IsLocked() bool {
	if u.BloqueadoHasta.IsZero() {
		return false
	}
	return time.Now().Before(u.BloqueadoHasta)
}
