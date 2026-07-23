package adapters_test

import (
	"context"
	"testing"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/adapters"
)

func TestSQLiteProductRepository_Create_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteProductRepository(db.RW)
	ctx := context.Background()

	product := &entities.Product{
		Nombre:       "Producto Test",
		SKU:          "TST-001",
		CategoriaID:  1,
		PrecioVenta:  25.50,
		PrecioCompra: 18.00,
		StockActual:  100,
		StockMinimo:  10,
		Unidad:       entities.UnitUnidad,
		Activo:       true,
	}

	// Create
	err := repo.Create(ctx, product)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if product.ID == 0 {
		t.Fatal("Create() did not set product ID")
	}

	// FindByID
	found, err := repo.FindByID(ctx, product.ID)
	if err != nil {
		t.Fatalf("FindByID(%d) error = %v", product.ID, err)
	}

	if found.Nombre != "Producto Test" {
		t.Errorf("found.Nombre = %q, want %q", found.Nombre, "Producto Test")
	}
	if found.SKU != "TST-001" {
		t.Errorf("found.SKU = %q, want %q", found.SKU, "TST-001")
	}
	if found.CategoriaID != 1 {
		t.Errorf("found.CategoriaID = %d, want 1", found.CategoriaID)
	}
	if found.PrecioVenta != 25.50 {
		t.Errorf("found.PrecioVenta = %f, want 25.50", found.PrecioVenta)
	}
	if found.PrecioCompra != 18.00 {
		t.Errorf("found.PrecioCompra = %f, want 18.00", found.PrecioCompra)
	}
	if found.StockActual != 100 {
		t.Errorf("found.StockActual = %f, want 100", found.StockActual)
	}
	if found.StockMinimo != 10 {
		t.Errorf("found.StockMinimo = %f, want 10", found.StockMinimo)
	}
	if found.Unidad != entities.UnitUnidad {
		t.Errorf("found.Unidad = %q, want %q", found.Unidad, entities.UnitUnidad)
	}
	if !found.Activo {
		t.Error("found.Activo = false, want true")
	}
}

func TestSQLiteProductRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteProductRepository(db.RW)

	_, err := repo.FindByID(context.Background(), 99999)
	if err == nil {
		t.Fatal("FindByID(99999) expected error, got nil")
	}
}

func TestSQLiteProductRepository_List_NoFilter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteProductRepository(db.RW)
	ctx := context.Background()

	products, err := repo.List(ctx, ports.ProductFilter{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	// Seed has 30 products
	if len(products) < 30 {
		t.Errorf("List() got %d products, want at least 30", len(products))
	}
}

func TestSQLiteProductRepository_List_FilterByCategory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteProductRepository(db.RW)
	ctx := context.Background()

	catID := int64(1) // Bebidas
	products, err := repo.List(ctx, ports.ProductFilter{CategoriaID: &catID})
	if err != nil {
		t.Fatalf("List(CategoriaID=1) error = %v", err)
	}

	// Seed has 5 bebidas
	if len(products) != 5 {
		t.Errorf("List(CategoriaID=1) got %d products, want 5", len(products))
	}

	for _, p := range products {
		if p.CategoriaID != 1 {
			t.Errorf("product %q has CategoriaID = %d, want 1", p.Nombre, p.CategoriaID)
		}
	}
}

func TestSQLiteProductRepository_List_FilterByActivo(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteProductRepository(db.RW)
	ctx := context.Background()

	// First deactivate one product
	err := repo.Deactivate(ctx, 1)
	if err != nil {
		t.Fatalf("Deactivate(1) error = %v", err)
	}

	activo := true
	products, err := repo.List(ctx, ports.ProductFilter{Activo: &activo})
	if err != nil {
		t.Fatalf("List(Activo=true) error = %v", err)
	}

	// Should be 29 (30 - 1 deactivated)
	if len(products) != 29 {
		t.Errorf("List(Activo=true) got %d products, want 29", len(products))
	}
}

func TestSQLiteProductRepository_List_FilterBySearch(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteProductRepository(db.RW)
	ctx := context.Background()

	products, err := repo.List(ctx, ports.ProductFilter{Search: "Coca"})
	if err != nil {
		t.Fatalf("List(Search=Coca) error = %v", err)
	}

	if len(products) != 1 {
		t.Errorf("List(Search=Coca) got %d products, want 1", len(products))
	}

	if len(products) > 0 && products[0].Nombre != "Coca Cola 600ml" {
		t.Errorf("products[0].Nombre = %q, want %q", products[0].Nombre, "Coca Cola 600ml")
	}
}

func TestSQLiteProductRepository_Deactivate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteProductRepository(db.RW)
	ctx := context.Background()

	// Deactivate product ID 1
	err := repo.Deactivate(ctx, 1)
	if err != nil {
		t.Fatalf("Deactivate(1) error = %v", err)
	}

	// Verify it's inactive
	product, err := repo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("FindByID(1) after deactivate error = %v", err)
	}

	if product.Activo {
		t.Error("product.Activo = true after Deactivate(), want false")
	}
}

func TestSQLiteProductRepository_Deactivate_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteProductRepository(db.RW)

	err := repo.Deactivate(context.Background(), 99999)
	if err == nil {
		t.Fatal("Deactivate(99999) expected error, got nil")
	}
}

func TestSQLiteProductRepository_FindLowStock(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteProductRepository(db.RW)
	ctx := context.Background()

	// Create a product with stock at minimum
	lowStockProduct := &entities.Product{
		Nombre:       "Low Stock Item",
		SKU:          "LOW-001",
		CategoriaID:  1,
		PrecioVenta:  10.00,
		PrecioCompra: 5.00,
		StockActual:  2,
		StockMinimo:  5,
		Unidad:       entities.UnitUnidad,
		Activo:       true,
	}

	err := repo.Create(ctx, lowStockProduct)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	products, err := repo.FindLowStock(ctx)
	if err != nil {
		t.Fatalf("FindLowStock() error = %v", err)
	}

	// Should find our low-stock product
	if len(products) == 0 {
		t.Fatal("FindLowStock() returned 0 products, want at least 1")
	}

	// Verify our product is in the list
	found := false
	for _, p := range products {
		if p.ID == lowStockProduct.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("FindLowStock() did not include product ID %d", lowStockProduct.ID)
	}

	// Verify all returned products actually have low stock
	for _, p := range products {
		if p.StockActual > p.StockMinimo {
			t.Errorf("product %q has StockActual=%f > StockMinimo=%f, should not be in low stock list",
				p.Nombre, p.StockActual, p.StockMinimo)
		}
	}
}

func TestSQLiteProductRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := adapters.NewSQLiteProductRepository(db.RW)
	ctx := context.Background()

	// Get an existing product from seed
	product, err := repo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("FindByID(1) error = %v", err)
	}

	// Update fields
	product.Nombre = "Coca Cola 600ml UPDATED"
	product.PrecioVenta = 25.00

	err = repo.Update(ctx, product)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	// Verify update
	updated, err := repo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("FindByID(1) after update error = %v", err)
	}

	if updated.Nombre != "Coca Cola 600ml UPDATED" {
		t.Errorf("updated.Nombre = %q, want %q", updated.Nombre, "Coca Cola 600ml UPDATED")
	}
	if updated.PrecioVenta != 25.00 {
		t.Errorf("updated.PrecioVenta = %f, want 25.00", updated.PrecioVenta)
	}
}
