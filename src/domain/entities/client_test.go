package entities

import (
	"errors"
	"testing"
)

func TestClient_Validate_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		client  Client
		wantErr error
	}{
		{
			name:    "valid name",
			client:  Client{Nombre: "Juan Pérez"},
			wantErr: nil,
		},
		{
			name:    "valid name with optional fields",
			client:  Client{Nombre: "María López", Telefono: "555-1234", Direccion: "Calle 5"},
			wantErr: nil,
		},
		{
			name:    "empty string",
			client:  Client{Nombre: ""},
			wantErr: ErrClientNameRequired,
		},
		{
			name:    "whitespace only spaces",
			client:  Client{Nombre: "   "},
			wantErr: ErrClientNameRequired,
		},
		{
			name:    "whitespace only tabs",
			client:  Client{Nombre: "\t\t"},
			wantErr: ErrClientNameRequired,
		},
		{
			name:    "whitespace mixed",
			client:  Client{Nombre: " \t \n "},
			wantErr: ErrClientNameRequired,
		},
		{
			name:    "name with leading and trailing spaces is valid",
			client:  Client{Nombre: "  Ana  "},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.Validate()

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
