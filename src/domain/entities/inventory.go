package entities

import (
	"errors"
	"time"
)

// Domain errors for InventoryMovement validation.
var (
	ErrMovementTypeInvalid  = errors.New("movement type must be entrada, salida, or ajuste")
	ErrMovementQtyZero      = errors.New("movement quantity cannot be zero")
	ErrMovementStockNegative = errors.New("resulting stock cannot be negative")
	ErrMovementNoProduct    = errors.New("movement must reference a valid product")
)

// MovementType represents the kind of inventory movement.
type MovementType string

const (
	MovimientoEntrada MovementType = "entrada"
	MovimientoSalida  MovementType = "salida"
	MovimientoAjuste  MovementType = "ajuste"
)

// ValidMovementTypes contains all valid movement type values for lookup.
var ValidMovementTypes = map[MovementType]bool{
	MovimientoEntrada: true,
	MovimientoSalida:  true,
	MovimientoAjuste:  true,
}

// InventoryMovement represents a stock change event in the POS system.
// Maps to the "inventario_movimientos" table in the database.
type InventoryMovement struct {
	ID              int64
	ProductoID      int64
	Tipo            MovementType
	Cantidad        float64
	StockResultante float64
	ReferenciaTipo  *string
	ReferenciaID    *int64
	Motivo          string
	UsuarioID       int64
	CreatedAt       time.Time
}

// Validate checks that the InventoryMovement satisfies all domain invariants.
// Returns the first validation error found using guard clauses.
func (m *InventoryMovement) Validate() error {
	if m.ProductoID <= 0 {
		return ErrMovementNoProduct
	}

	if !ValidMovementTypes[m.Tipo] {
		return ErrMovementTypeInvalid
	}

	if m.Cantidad == 0 {
		return ErrMovementQtyZero
	}

	if m.StockResultante < 0 {
		return ErrMovementStockNegative
	}

	return nil
}
