package entities

import (
	"errors"
	"testing"
	"time"
)

func TestUser_Validate_ValidUser(t *testing.T) {
	u := &User{
		ID:     1,
		Nombre: "Carlos Admin",
		Rol:    RoleAdmin,
		Activo: true,
	}

	if err := u.Validate(); err != nil {
		t.Errorf("expected no error for valid user, got: %v", err)
	}
}

func TestUser_Validate_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr error
	}{
		{
			name: "valid admin user",
			user: User{
				Nombre: "Admin Principal",
				Rol:    RoleAdmin,
				Activo: true,
			},
			wantErr: nil,
		},
		{
			name: "valid cajero user",
			user: User{
				Nombre: "Maria Cajera",
				Rol:    RoleCajero,
				Activo: true,
			},
			wantErr: nil,
		},
		{
			name: "valid inactive user",
			user: User{
				Nombre: "Juan Inactivo",
				Rol:    RoleCajero,
				Activo: false,
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			user: User{
				Nombre: "",
				Rol:    RoleAdmin,
			},
			wantErr: ErrUserNameEmpty,
		},
		{
			name: "invalid role empty",
			user: User{
				Nombre: "Test User",
				Rol:    "",
			},
			wantErr: ErrUserRoleInvalid,
		},
		{
			name: "invalid role unknown",
			user: User{
				Nombre: "Test User",
				Rol:    Role("gerente"),
			},
			wantErr: ErrUserRoleInvalid,
		},
		{
			name: "invalid role with wrong case",
			user: User{
				Nombre: "Test User",
				Rol:    Role("Admin"),
			},
			wantErr: ErrUserRoleInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
				return
			}

			if err == nil {
				t.Errorf("expected error %v, got nil", tt.wantErr)
				return
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestUser_IsLocked(t *testing.T) {
	tests := []struct {
		name           string
		bloqueadoHasta time.Time
		want           bool
	}{
		{
			name:           "not locked - zero time",
			bloqueadoHasta: time.Time{},
			want:           false,
		},
		{
			name:           "not locked - lockout expired",
			bloqueadoHasta: time.Now().Add(-1 * time.Minute),
			want:           false,
		},
		{
			name:           "locked - lockout in future",
			bloqueadoHasta: time.Now().Add(5 * time.Minute),
			want:           true,
		},
		{
			name:           "not locked - lockout exactly now (edge)",
			bloqueadoHasta: time.Now(),
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Nombre:         "Test User",
				Rol:            RoleCajero,
				BloqueadoHasta: tt.bloqueadoHasta,
			}

			got := u.IsLocked()
			if got != tt.want {
				t.Errorf("IsLocked() = %v, want %v (BloqueadoHasta: %v, Now: %v)",
					got, tt.want, tt.bloqueadoHasta, time.Now())
			}
		})
	}
}

func TestUser_MaxAttempts(t *testing.T) {
	if MaxAttempts != 5 {
		t.Errorf("MaxAttempts = %d, want 5", MaxAttempts)
	}
}

func TestRole_Constants(t *testing.T) {
	if RoleAdmin != "admin" {
		t.Errorf("RoleAdmin = %q, want %q", RoleAdmin, "admin")
	}
	if RoleCajero != "cajero" {
		t.Errorf("RoleCajero = %q, want %q", RoleCajero, "cajero")
	}
}
