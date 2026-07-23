package entities

import (
	"errors"
	"testing"
)

func TestProduct_Validate_ValidProduct(t *testing.T) {
	p := &Product{
		Nombre:       "Coca-Cola 600ml",
		SKU:          "CC-600",
		CategoriaID:  1,
		PrecioVenta:  25.50,
		PrecioCompra: 18.00,
		StockActual:  100,
		StockMinimo:  10,
		Unidad:       UnitUnidad,
		Activo:       true,
	}

	if err := p.Validate(); err != nil {
		t.Errorf("expected no error for valid product, got: %v", err)
	}
}

func TestProduct_Validate_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		product Product
		wantErr error
	}{
		{
			name: "valid product with minimum fields",
			product: Product{
				Nombre:      "Agua Natural",
				PrecioVenta: 10.0,
			},
			wantErr: nil,
		},
		{
			name: "valid product with all fields",
			product: Product{
				Nombre:       "Sabritas Original",
				SKU:          "SAB-001",
				CategoriaID:  2,
				PrecioVenta:  22.00,
				PrecioCompra: 15.00,
				StockActual:  50,
				StockMinimo:  5,
				Unidad:       UnitPaquete,
				Activo:       true,
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			product: Product{
				Nombre:      "",
				PrecioVenta: 20.0,
			},
			wantErr: ErrProductNameEmpty,
		},
		{
			name: "negative sale price",
			product: Product{
				Nombre:      "Producto X",
				PrecioVenta: -5.0,
			},
			wantErr: ErrProductPriceInvalid,
		},
		{
			name: "zero sale price",
			product: Product{
				Nombre:      "Producto Y",
				PrecioVenta: 0,
			},
			wantErr: ErrProductPriceInvalid,
		},
		{
			name: "negative purchase price",
			product: Product{
				Nombre:       "Producto Z",
				PrecioVenta:  20.0,
				PrecioCompra: -1.0,
			},
			wantErr: ErrProductCostNegative,
		},
		{
			name: "negative stock actual",
			product: Product{
				Nombre:      "Producto W",
				PrecioVenta: 20.0,
				StockActual: -10,
			},
			wantErr: ErrProductStockNegative,
		},
		{
			name: "negative stock minimo",
			product: Product{
				Nombre:      "Producto V",
				PrecioVenta: 20.0,
				StockMinimo: -5,
			},
			wantErr: ErrProductStockNegative,
		},
		{
			name: "zero stock is valid",
			product: Product{
				Nombre:      "Producto Agotado",
				PrecioVenta: 15.0,
				StockActual: 0,
				StockMinimo: 0,
			},
			wantErr: nil,
		},
		{
			name: "zero purchase price is valid",
			product: Product{
				Nombre:       "Regalo",
				PrecioVenta:  10.0,
				PrecioCompra: 0,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()

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

func TestProduct_IsLowStock(t *testing.T) {
	tests := []struct {
		name        string
		stockActual float64
		stockMinimo float64
		want        bool
	}{
		{"stock below minimum", 3, 10, true},
		{"stock equals minimum", 10, 10, true},
		{"stock above minimum", 50, 10, false},
		{"both zero", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Product{
				Nombre:      "Test",
				PrecioVenta: 10,
				StockActual: tt.stockActual,
				StockMinimo: tt.stockMinimo,
			}

			got := p.IsLowStock()
			if got != tt.want {
				t.Errorf("IsLowStock() = %v, want %v", got, tt.want)
			}
		})
	}
}
