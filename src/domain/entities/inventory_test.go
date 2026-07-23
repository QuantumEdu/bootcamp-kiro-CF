package entities

import (
	"testing"
	"time"
)

func TestInventoryMovement_Validate_Valid(t *testing.T) {
	refTipo := "venta"
	refID := int64(1)

	tests := []struct {
		name     string
		movement InventoryMovement
	}{
		{
			"valid entrada movement",
			InventoryMovement{
				ID: 1, ProductoID: 10, Tipo: MovimientoEntrada,
				Cantidad: 5, StockResultante: 15, Motivo: "compra proveedor",
				UsuarioID: 1, CreatedAt: time.Now(),
			},
		},
		{
			"valid salida movement with negative cantidad",
			InventoryMovement{
				ID: 2, ProductoID: 10, Tipo: MovimientoSalida,
				Cantidad: -3, StockResultante: 7, Motivo: "venta",
				UsuarioID: 1, ReferenciaTipo: &refTipo, ReferenciaID: &refID,
				CreatedAt: time.Now(),
			},
		},
		{
			"valid ajuste movement",
			InventoryMovement{
				ID: 3, ProductoID: 5, Tipo: MovimientoAjuste,
				Cantidad: -2, StockResultante: 0, Motivo: "inventario fisico",
				UsuarioID: 2, CreatedAt: time.Now(),
			},
		},
		{
			"valid movement with zero resulting stock",
			InventoryMovement{
				ID: 4, ProductoID: 1, Tipo: MovimientoSalida,
				Cantidad: -10, StockResultante: 0, Motivo: "liquidacion",
				UsuarioID: 1, CreatedAt: time.Now(),
			},
		},
		{
			"valid movement with nil references",
			InventoryMovement{
				ID: 5, ProductoID: 3, Tipo: MovimientoEntrada,
				Cantidad: 20, StockResultante: 20, Motivo: "stock inicial",
				UsuarioID: 1, CreatedAt: time.Now(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.movement.Validate()
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestInventoryMovement_Validate_Invalid(t *testing.T) {
	tests := []struct {
		name     string
		movement InventoryMovement
		wantErr  error
	}{
		{
			"zero producto_id",
			InventoryMovement{
				ProductoID: 0, Tipo: MovimientoEntrada,
				Cantidad: 5, StockResultante: 10, UsuarioID: 1,
			},
			ErrMovementNoProduct,
		},
		{
			"negative producto_id",
			InventoryMovement{
				ProductoID: -1, Tipo: MovimientoEntrada,
				Cantidad: 5, StockResultante: 10, UsuarioID: 1,
			},
			ErrMovementNoProduct,
		},
		{
			"invalid movement type",
			InventoryMovement{
				ProductoID: 1, Tipo: "devolucion",
				Cantidad: 5, StockResultante: 10, UsuarioID: 1,
			},
			ErrMovementTypeInvalid,
		},
		{
			"empty movement type",
			InventoryMovement{
				ProductoID: 1, Tipo: "",
				Cantidad: 5, StockResultante: 10, UsuarioID: 1,
			},
			ErrMovementTypeInvalid,
		},
		{
			"zero cantidad",
			InventoryMovement{
				ProductoID: 1, Tipo: MovimientoSalida,
				Cantidad: 0, StockResultante: 10, UsuarioID: 1,
			},
			ErrMovementQtyZero,
		},
		{
			"negative resulting stock",
			InventoryMovement{
				ProductoID: 1, Tipo: MovimientoSalida,
				Cantidad: -5, StockResultante: -1, UsuarioID: 1,
			},
			ErrMovementStockNegative,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.movement.Validate()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err != tt.wantErr {
				t.Errorf("got error %q, want %q", err, tt.wantErr)
			}
		})
	}
}
