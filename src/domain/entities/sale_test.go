package entities

import (
	"errors"
	"testing"
)

func TestSaleItem_Validate_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		item    SaleItem
		wantErr error
	}{
		{
			name: "valid item",
			item: SaleItem{
				ProductoID:     1,
				Cantidad:       2,
				PrecioUnitario: 25.50,
				Subtotal:       51.00,
			},
			wantErr: nil,
		},
		{
			name: "valid item fractional quantity",
			item: SaleItem{
				ProductoID:     1,
				Cantidad:       0.5,
				PrecioUnitario: 100.0,
				Subtotal:       50.0,
			},
			wantErr: nil,
		},
		{
			name: "zero quantity",
			item: SaleItem{
				ProductoID:     1,
				Cantidad:       0,
				PrecioUnitario: 25.0,
				Subtotal:       0,
			},
			wantErr: ErrSaleItemQtyInvalid,
		},
		{
			name: "negative quantity",
			item: SaleItem{
				ProductoID:     1,
				Cantidad:       -3,
				PrecioUnitario: 25.0,
				Subtotal:       -75.0,
			},
			wantErr: ErrSaleItemQtyInvalid,
		},
		{
			name: "zero unit price",
			item: SaleItem{
				ProductoID:     1,
				Cantidad:       2,
				PrecioUnitario: 0,
				Subtotal:       0,
			},
			wantErr: ErrSaleItemPriceInvalid,
		},
		{
			name: "negative unit price",
			item: SaleItem{
				ProductoID:     1,
				Cantidad:       2,
				PrecioUnitario: -10.0,
				Subtotal:       -20.0,
			},
			wantErr: ErrSaleItemPriceInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()

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

func TestSale_Validate_TableDriven(t *testing.T) {
	validItems := []SaleItem{
		{ProductoID: 1, Cantidad: 2, PrecioUnitario: 25.0, Subtotal: 50.0},
		{ProductoID: 2, Cantidad: 1, PrecioUnitario: 30.0, Subtotal: 30.0},
	}

	tests := []struct {
		name    string
		sale    Sale
		wantErr error
	}{
		{
			name: "valid sale with multiple items",
			sale: Sale{
				UsuarioID:  1,
				Total:      80.0,
				MetodoPago: MetodoEfectivo,
				Items:      validItems,
			},
			wantErr: nil,
		},
		{
			name: "valid sale with client",
			sale: Sale{
				UsuarioID:  1,
				ClienteID:  int64Ptr(5),
				Total:      50.0,
				MetodoPago: MetodoTarjeta,
				Items:      []SaleItem{{ProductoID: 1, Cantidad: 2, PrecioUnitario: 25.0, Subtotal: 50.0}},
			},
			wantErr: nil,
		},
		{
			name: "valid sale zero total",
			sale: Sale{
				UsuarioID:  1,
				Total:      0,
				MetodoPago: MetodoTransferencia,
				Items:      []SaleItem{{ProductoID: 1, Cantidad: 1, PrecioUnitario: 10.0, Subtotal: 10.0}},
			},
			wantErr: nil,
		},
		{
			name: "valid sale with mixto payment",
			sale: Sale{
				UsuarioID:  1,
				Total:      100.0,
				MetodoPago: MetodoMixto,
				Items:      []SaleItem{{ProductoID: 1, Cantidad: 4, PrecioUnitario: 25.0, Subtotal: 100.0}},
			},
			wantErr: nil,
		},
		{
			name: "no items",
			sale: Sale{
				UsuarioID:  1,
				Total:      0,
				MetodoPago: MetodoEfectivo,
				Items:      []SaleItem{},
			},
			wantErr: ErrSaleNoItems,
		},
		{
			name: "nil items slice",
			sale: Sale{
				UsuarioID:  1,
				Total:      0,
				MetodoPago: MetodoEfectivo,
				Items:      nil,
			},
			wantErr: ErrSaleNoItems,
		},
		{
			name: "negative total",
			sale: Sale{
				UsuarioID:  1,
				Total:      -10.0,
				MetodoPago: MetodoEfectivo,
				Items:      validItems,
			},
			wantErr: ErrSaleInvalidTotal,
		},
		{
			name: "invalid payment method",
			sale: Sale{
				UsuarioID:  1,
				Total:      80.0,
				MetodoPago: "bitcoin",
				Items:      validItems,
			},
			wantErr: ErrSaleInvalidPayment,
		},
		{
			name: "empty payment method",
			sale: Sale{
				UsuarioID:  1,
				Total:      80.0,
				MetodoPago: "",
				Items:      validItems,
			},
			wantErr: ErrSaleInvalidPayment,
		},
		{
			name: "item with invalid quantity propagates error",
			sale: Sale{
				UsuarioID:  1,
				Total:      50.0,
				MetodoPago: MetodoEfectivo,
				Items:      []SaleItem{{ProductoID: 1, Cantidad: 0, PrecioUnitario: 25.0, Subtotal: 0}},
			},
			wantErr: ErrSaleItemQtyInvalid,
		},
		{
			name: "item with invalid price propagates error",
			sale: Sale{
				UsuarioID:  1,
				Total:      50.0,
				MetodoPago: MetodoEfectivo,
				Items:      []SaleItem{{ProductoID: 1, Cantidad: 2, PrecioUnitario: -5.0, Subtotal: -10.0}},
			},
			wantErr: ErrSaleItemPriceInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sale.Validate()

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

func TestSale_CalculateTotal(t *testing.T) {
	tests := []struct {
		name  string
		items []SaleItem
		want  float64
	}{
		{
			name: "single item",
			items: []SaleItem{
				{ProductoID: 1, Cantidad: 2, PrecioUnitario: 25.0, Subtotal: 50.0},
			},
			want: 50.0,
		},
		{
			name: "multiple items",
			items: []SaleItem{
				{ProductoID: 1, Cantidad: 2, PrecioUnitario: 25.0, Subtotal: 50.0},
				{ProductoID: 2, Cantidad: 1, PrecioUnitario: 30.0, Subtotal: 30.0},
				{ProductoID: 3, Cantidad: 3, PrecioUnitario: 15.0, Subtotal: 45.0},
			},
			want: 125.0,
		},
		{
			name:  "no items returns zero",
			items: []SaleItem{},
			want:  0,
		},
		{
			name: "fractional subtotals",
			items: []SaleItem{
				{ProductoID: 1, Cantidad: 0.5, PrecioUnitario: 100.0, Subtotal: 50.0},
				{ProductoID: 2, Cantidad: 1.5, PrecioUnitario: 20.0, Subtotal: 30.0},
			},
			want: 80.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sale{
				UsuarioID:  1,
				MetodoPago: MetodoEfectivo,
				Items:      tt.items,
			}

			got := s.CalculateTotal()
			if got != tt.want {
				t.Errorf("CalculateTotal() = %v, want %v", got, tt.want)
			}

			// Verify Total field was updated
			if s.Total != tt.want {
				t.Errorf("Sale.Total = %v, want %v after CalculateTotal()", s.Total, tt.want)
			}
		})
	}
}

func TestPaymentMethod_Constants(t *testing.T) {
	// Verify all payment methods are in the valid map
	methods := []PaymentMethod{MetodoEfectivo, MetodoTarjeta, MetodoTransferencia, MetodoMixto}
	for _, m := range methods {
		if !ValidPaymentMethods[m] {
			t.Errorf("expected %q to be a valid payment method", m)
		}
	}

	// Verify invalid method is not in the map
	if ValidPaymentMethods["cheque"] {
		t.Error("expected 'cheque' to NOT be a valid payment method")
	}
}

// int64Ptr is a test helper to create *int64 values.
func int64Ptr(v int64) *int64 {
	return &v
}
