package entities

import (
	"errors"
	"time"
)

// Domain errors for Sale validation.
var (
	ErrSaleNoItems        = errors.New("sale must have at least one item")
	ErrSaleInvalidTotal   = errors.New("sale total must be greater than or equal to zero")
	ErrSaleInvalidPayment = errors.New("sale payment method is invalid")
	ErrSaleItemQtyInvalid   = errors.New("sale item quantity must be greater than zero")
	ErrSaleItemPriceInvalid = errors.New("sale item unit price must be greater than zero")
)

// PaymentMethod represents how a sale was paid.
type PaymentMethod string

const (
	MetodoEfectivo      PaymentMethod = "efectivo"
	MetodoTarjeta       PaymentMethod = "tarjeta"
	MetodoTransferencia PaymentMethod = "transferencia"
	MetodoMixto         PaymentMethod = "mixto"
)

// ValidPaymentMethods contains all valid payment method values for lookup.
var ValidPaymentMethods = map[PaymentMethod]bool{
	MetodoEfectivo:      true,
	MetodoTarjeta:       true,
	MetodoTransferencia: true,
	MetodoMixto:         true,
}

// SaleItem represents a single line item within a sale.
// Maps to the "venta_items" table in the database.
type SaleItem struct {
	ID             int64
	VentaID        int64
	ProductoID     int64
	Cantidad       float64
	PrecioUnitario float64
	Subtotal       float64
}

// Validate checks that the SaleItem satisfies all domain invariants.
func (si *SaleItem) Validate() error {
	if si.Cantidad <= 0 {
		return ErrSaleItemQtyInvalid
	}

	if si.PrecioUnitario <= 0 {
		return ErrSaleItemPriceInvalid
	}

	return nil
}

// Sale represents a completed sale transaction in the POS system.
// Maps to the "ventas" table in the database.
type Sale struct {
	ID         int64
	UsuarioID  int64
	ClienteID  *int64
	Total      float64
	MetodoPago PaymentMethod
	Items      []SaleItem
	CreatedAt  time.Time
}

// Validate checks that the Sale satisfies all domain invariants.
// Returns the first validation error found using guard clauses.
func (s *Sale) Validate() error {
	if len(s.Items) == 0 {
		return ErrSaleNoItems
	}

	if s.Total < 0 {
		return ErrSaleInvalidTotal
	}

	if !ValidPaymentMethods[s.MetodoPago] {
		return ErrSaleInvalidPayment
	}

	for i := range s.Items {
		if err := s.Items[i].Validate(); err != nil {
			return err
		}
	}

	return nil
}

// CalculateTotal sums all item subtotals and updates the Sale's Total field.
func (s *Sale) CalculateTotal() float64 {
	var total float64
	for _, item := range s.Items {
		total += item.Subtotal
	}
	s.Total = total
	return total
}
