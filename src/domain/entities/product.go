// Package entities contains the core business objects of the POS system.
package entities

import (
	"errors"
	"time"
)

// Domain errors for Product validation.
var (
	ErrProductNameEmpty     = errors.New("product name cannot be empty")
	ErrProductPriceInvalid  = errors.New("product sale price must be greater than zero")
	ErrProductCostNegative  = errors.New("product purchase price cannot be negative")
	ErrProductStockNegative = errors.New("product stock cannot be negative")
)

// Unit represents the measurement unit for a product.
type Unit string

const (
	UnitUnidad  Unit = "unidad"
	UnitKg      Unit = "kg"
	UnitLitro   Unit = "litro"
	UnitPaquete Unit = "paquete"
)

// Product represents a catalog item in the POS system.
// Maps to the "productos" table in the database.
type Product struct {
	ID           int64
	Nombre       string
	SKU          string
	CategoriaID  int64
	PrecioVenta  float64
	PrecioCompra float64
	StockActual  float64
	StockMinimo  float64
	Unidad       Unit
	Activo       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Validate checks that the Product satisfies all domain invariants.
// Returns the first validation error found using guard clauses.
func (p *Product) Validate() error {
	if p.Nombre == "" {
		return ErrProductNameEmpty
	}

	if p.PrecioVenta <= 0 {
		return ErrProductPriceInvalid
	}

	if p.PrecioCompra < 0 {
		return ErrProductCostNegative
	}

	if p.StockActual < 0 {
		return ErrProductStockNegative
	}

	if p.StockMinimo < 0 {
		return ErrProductStockNegative
	}

	return nil
}

// IsLowStock returns true when the product's current stock is at or below the minimum threshold.
func (p *Product) IsLowStock() bool {
	return p.StockActual <= p.StockMinimo
}
