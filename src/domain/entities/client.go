package entities

import (
	"errors"
	"strings"
	"time"
)

// Domain errors for Client validation.
var ErrClientNameRequired = errors.New("client name is required")

// Client represents a customer in the POS system.
// Maps to the "clientes" table in the database.
type Client struct {
	ID        int64
	Nombre    string
	Telefono  string
	Direccion string
	CreatedAt time.Time
}

// Validate checks domain invariants for a client.
// Returns an error if the Nombre is empty or whitespace-only.
func (c *Client) Validate() error {
	if strings.TrimSpace(c.Nombre) == "" {
		return ErrClientNameRequired
	}
	return nil
}
