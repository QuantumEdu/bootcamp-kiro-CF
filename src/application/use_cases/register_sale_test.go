package use_cases

import (
	"context"
	"errors"
	"testing"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

// --- Mock Repositories ---

// mockSaleRepository implements ports.SaleRepository for testing.
type mockSaleRepository struct {
	sales      []*entities.Sale
	createErr  error
	createCall int
}

func (m *mockSaleRepository) Create(_ context.Context, sale *entities.Sale) error {
	m.createCall++
	if m.createErr != nil {
		return m.createErr
	}
	sale.ID = int64(m.createCall)
	m.sales = append(m.sales, sale)
	return nil
}

func (m *mockSaleRepository) FindByID(_ context.Context, id int64) (*entities.Sale, error) {
	for _, s := range m.sales {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockSaleRepository) List(_ context.Context, _ ports.SaleFilter) ([]entities.Sale, error) {
	var result []entities.Sale
	for _, s := range m.sales {
		result = append(result, *s)
	}
	return result, nil
}

// mockProductRepository implements ports.ProductRepository for testing.
type mockProductRepository struct {
	products   map[int64]*entities.Product
	updates    []*entities.Product
	updateCall int
}

func newMockProductRepo(products ...*entities.Product) *mockProductRepository {
	m := &mockProductRepository{
		products: make(map[int64]*entities.Product),
	}
	for _, p := range products {
		m.products[p.ID] = p
	}
	return m
}

func (m *mockProductRepository) Create(_ context.Context, product *entities.Product) error {
	m.products[product.ID] = product
	return nil
}

func (m *mockProductRepository) Update(_ context.Context, product *entities.Product) error {
	m.updateCall++
	// Clone to capture state at time of update.
	clone := *product
	m.updates = append(m.updates, &clone)
	m.products[product.ID] = product
	return nil
}

func (m *mockProductRepository) FindByID(_ context.Context, id int64) (*entities.Product, error) {
	p, ok := m.products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return p, nil
}

func (m *mockProductRepository) List(_ context.Context, _ ports.ProductFilter) ([]entities.Product, error) {
	return nil, nil
}

func (m *mockProductRepository) Deactivate(_ context.Context, _ int64) error {
	return nil
}

func (m *mockProductRepository) FindLowStock(_ context.Context) ([]entities.Product, error) {
	return nil, nil
}

// mockInventoryRepository implements ports.InventoryRepository for testing.
type mockInventoryRepository struct {
	movements  []*entities.InventoryMovement
	createCall int
}

func (m *mockInventoryRepository) Create(_ context.Context, movement *entities.InventoryMovement) error {
	m.createCall++
	movement.ID = int64(m.createCall)
	m.movements = append(m.movements, movement)
	return nil
}

func (m *mockInventoryRepository) FindByProduct(_ context.Context, productID int64) ([]entities.InventoryMovement, error) {
	var result []entities.InventoryMovement
	for _, mov := range m.movements {
		if mov.ProductoID == productID {
			result = append(result, *mov)
		}
	}
	return result, nil
}

// --- Tests ---

func TestRegisterSale_Execute_SuccessfulSale(t *testing.T) {
	productA := &entities.Product{
		ID: 1, Nombre: "Coca-Cola", PrecioVenta: 25.00, StockActual: 50, Activo: true, Unidad: entities.UnitUnidad,
	}
	productB := &entities.Product{
		ID: 2, Nombre: "Papas Fritas", PrecioVenta: 15.50, StockActual: 30, Activo: true, Unidad: entities.UnitUnidad,
	}

	saleRepo := &mockSaleRepository{}
	productRepo := newMockProductRepo(productA, productB)
	inventoryRepo := &mockInventoryRepository{}

	uc := NewRegisterSale(saleRepo, productRepo, inventoryRepo)

	input := RegisterSaleInput{
		UsuarioID:  1,
		MetodoPago: entities.MetodoEfectivo,
		Items: []SaleItemInput{
			{ProductoID: 1, Cantidad: 3},
			{ProductoID: 2, Cantidad: 2},
		},
	}

	sale, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sale == nil {
		t.Fatal("expected sale, got nil")
	}

	// Verify sale properties.
	expectedTotal := (25.00 * 3) + (15.50 * 2) // 75 + 31 = 106
	if sale.Total != expectedTotal {
		t.Errorf("expected total %.2f, got %.2f", expectedTotal, sale.Total)
	}
	if len(sale.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(sale.Items))
	}
	if sale.MetodoPago != entities.MetodoEfectivo {
		t.Errorf("expected payment method 'efectivo', got '%s'", sale.MetodoPago)
	}
	if sale.ID == 0 {
		t.Error("expected sale ID to be set after creation")
	}

	// Verify stock was deducted.
	if productRepo.products[1].StockActual != 47 { // 50 - 3
		t.Errorf("expected product 1 stock 47, got %.2f", productRepo.products[1].StockActual)
	}
	if productRepo.products[2].StockActual != 28 { // 30 - 2
		t.Errorf("expected product 2 stock 28, got %.2f", productRepo.products[2].StockActual)
	}

	// Verify inventory movements were created.
	if inventoryRepo.createCall != 2 {
		t.Errorf("expected 2 inventory movements, got %d", inventoryRepo.createCall)
	}
	for _, mov := range inventoryRepo.movements {
		if mov.Tipo != entities.MovimientoSalida {
			t.Errorf("expected movement type 'salida', got '%s'", mov.Tipo)
		}
		if mov.UsuarioID != 1 {
			t.Errorf("expected movement user_id 1, got %d", mov.UsuarioID)
		}
	}

	// Verify product update was called (stock deduction).
	if productRepo.updateCall != 2 {
		t.Errorf("expected 2 product updates, got %d", productRepo.updateCall)
	}
}

func TestRegisterSale_Execute_InsufficientStock(t *testing.T) {
	product := &entities.Product{
		ID: 1, Nombre: "Coca-Cola", PrecioVenta: 25.00, StockActual: 2, Activo: true, Unidad: entities.UnitUnidad,
	}

	saleRepo := &mockSaleRepository{}
	productRepo := newMockProductRepo(product)
	inventoryRepo := &mockInventoryRepository{}

	uc := NewRegisterSale(saleRepo, productRepo, inventoryRepo)

	input := RegisterSaleInput{
		UsuarioID:  1,
		MetodoPago: entities.MetodoEfectivo,
		Items: []SaleItemInput{
			{ProductoID: 1, Cantidad: 5}, // Requesting 5, only 2 available.
		},
	}

	sale, err := uc.Execute(context.Background(), input)

	if sale != nil {
		t.Fatalf("expected nil sale, got %+v", sale)
	}
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrSaleInsufficientStock) {
		t.Errorf("expected ErrSaleInsufficientStock, got %v", err)
	}

	// Verify no sale was created.
	if saleRepo.createCall != 0 {
		t.Errorf("expected 0 sale create calls, got %d", saleRepo.createCall)
	}

	// Verify no stock deduction occurred.
	if productRepo.updateCall != 0 {
		t.Errorf("expected 0 product updates, got %d", productRepo.updateCall)
	}

	// Verify no inventory movements.
	if inventoryRepo.createCall != 0 {
		t.Errorf("expected 0 inventory movements, got %d", inventoryRepo.createCall)
	}
}

func TestRegisterSale_Execute_NegativeQuantity(t *testing.T) {
	product := &entities.Product{
		ID: 1, Nombre: "Coca-Cola", PrecioVenta: 25.00, StockActual: 50, Activo: true, Unidad: entities.UnitUnidad,
	}

	saleRepo := &mockSaleRepository{}
	productRepo := newMockProductRepo(product)
	inventoryRepo := &mockInventoryRepository{}

	uc := NewRegisterSale(saleRepo, productRepo, inventoryRepo)

	tests := []struct {
		name     string
		cantidad float64
	}{
		{"negative quantity", -5},
		{"zero quantity", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := RegisterSaleInput{
				UsuarioID:  1,
				MetodoPago: entities.MetodoEfectivo,
				Items: []SaleItemInput{
					{ProductoID: 1, Cantidad: tt.cantidad},
				},
			}

			sale, err := uc.Execute(context.Background(), input)

			if sale != nil {
				t.Fatalf("expected nil sale, got %+v", sale)
			}
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, ErrSaleInvalidQuantity) {
				t.Errorf("expected ErrSaleInvalidQuantity, got %v", err)
			}
		})
	}
}

func TestRegisterSale_Execute_NoItems(t *testing.T) {
	saleRepo := &mockSaleRepository{}
	productRepo := newMockProductRepo()
	inventoryRepo := &mockInventoryRepository{}

	uc := NewRegisterSale(saleRepo, productRepo, inventoryRepo)

	input := RegisterSaleInput{
		UsuarioID:  1,
		MetodoPago: entities.MetodoEfectivo,
		Items:      []SaleItemInput{},
	}

	sale, err := uc.Execute(context.Background(), input)

	if sale != nil {
		t.Fatalf("expected nil sale, got %+v", sale)
	}
	if !errors.Is(err, ErrSaleNoItemsProvided) {
		t.Errorf("expected ErrSaleNoItemsProvided, got %v", err)
	}
}

func TestRegisterSale_Execute_ProductNotFound(t *testing.T) {
	saleRepo := &mockSaleRepository{}
	productRepo := newMockProductRepo() // Empty — no products.
	inventoryRepo := &mockInventoryRepository{}

	uc := NewRegisterSale(saleRepo, productRepo, inventoryRepo)

	input := RegisterSaleInput{
		UsuarioID:  1,
		MetodoPago: entities.MetodoEfectivo,
		Items: []SaleItemInput{
			{ProductoID: 999, Cantidad: 1},
		},
	}

	sale, err := uc.Execute(context.Background(), input)

	if sale != nil {
		t.Fatalf("expected nil sale, got %+v", sale)
	}
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrSaleProductNotFound) {
		t.Errorf("expected ErrSaleProductNotFound, got %v", err)
	}
}
