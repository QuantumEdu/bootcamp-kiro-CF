package use_cases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

// Domain errors for RegisterSale use case.
var (
	ErrSaleInsufficientStock = errors.New("stock insuficiente para el producto")
	ErrSaleProductNotFound   = errors.New("producto no encontrado")
	ErrSaleNoItemsProvided   = errors.New("la venta debe tener al menos un item")
	ErrSaleInvalidQuantity   = errors.New("la cantidad debe ser mayor a cero")
)

// RegisterSaleInput holds the data needed to register a new sale.
type RegisterSaleInput struct {
	UsuarioID  int64
	ClienteID  *int64
	MetodoPago entities.PaymentMethod
	Items      []SaleItemInput
}

// SaleItemInput holds a single item in the sale request.
type SaleItemInput struct {
	ProductoID int64
	Cantidad   float64
}

// RegisterSale handles registering a new sale with stock validation and inventory deduction.
type RegisterSale struct {
	saleRepo      ports.SaleRepository
	productRepo   ports.ProductRepository
	inventoryRepo ports.InventoryRepository
}

// NewRegisterSale creates a new RegisterSale use case.
func NewRegisterSale(
	saleRepo ports.SaleRepository,
	productRepo ports.ProductRepository,
	inventoryRepo ports.InventoryRepository,
) *RegisterSale {
	return &RegisterSale{
		saleRepo:      saleRepo,
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
	}
}

// Execute validates items, checks stock, calculates total, creates the sale, deducts inventory.
func (uc *RegisterSale) Execute(ctx context.Context, input RegisterSaleInput) (*entities.Sale, error) {
	// Validate at least 1 item.
	if len(input.Items) == 0 {
		return nil, ErrSaleNoItemsProvided
	}

	// Validate quantities and build sale items.
	var saleItems []entities.SaleItem
	for _, item := range input.Items {
		if item.Cantidad <= 0 {
			return nil, ErrSaleInvalidQuantity
		}

		// Find product.
		product, err := uc.productRepo.FindByID(ctx, item.ProductoID)
		if err != nil {
			return nil, fmt.Errorf("%w: ID %d", ErrSaleProductNotFound, item.ProductoID)
		}

		// Check stock.
		if product.StockActual < item.Cantidad {
			return nil, fmt.Errorf("%w: %s (disponible: %.2f, solicitado: %.2f)",
				ErrSaleInsufficientStock, product.Nombre, product.StockActual, item.Cantidad)
		}

		// Build SaleItem with current price.
		saleItem := entities.SaleItem{
			ProductoID:     item.ProductoID,
			Cantidad:       item.Cantidad,
			PrecioUnitario: product.PrecioVenta,
			Subtotal:       product.PrecioVenta * item.Cantidad,
		}
		saleItems = append(saleItems, saleItem)
	}

	// Create sale entity.
	sale := &entities.Sale{
		UsuarioID:  input.UsuarioID,
		ClienteID:  input.ClienteID,
		MetodoPago: input.MetodoPago,
		Items:      saleItems,
		CreatedAt:  time.Now(),
	}

	// Calculate total.
	sale.CalculateTotal()

	// Validate the sale entity.
	if err := sale.Validate(); err != nil {
		return nil, fmt.Errorf("validating sale: %w", err)
	}

	// Persist the sale.
	if err := uc.saleRepo.Create(ctx, sale); err != nil {
		return nil, fmt.Errorf("creating sale: %w", err)
	}

	// Deduct stock and create inventory movements.
	for _, item := range input.Items {
		product, err := uc.productRepo.FindByID(ctx, item.ProductoID)
		if err != nil {
			return nil, fmt.Errorf("fetching product for stock deduction: %w", err)
		}

		// Deduct stock.
		product.StockActual -= item.Cantidad
		if err := uc.productRepo.Update(ctx, product); err != nil {
			return nil, fmt.Errorf("updating stock for product %d: %w", item.ProductoID, err)
		}

		// Create inventory movement (salida).
		refTipo := "venta"
		movement := &entities.InventoryMovement{
			ProductoID:      item.ProductoID,
			Tipo:            entities.MovimientoSalida,
			Cantidad:        item.Cantidad,
			StockResultante: product.StockActual,
			ReferenciaTipo:  &refTipo,
			ReferenciaID:    &sale.ID,
			Motivo:          fmt.Sprintf("Venta #%d", sale.ID),
			UsuarioID:       input.UsuarioID,
			CreatedAt:       time.Now(),
		}

		if err := uc.inventoryRepo.Create(ctx, movement); err != nil {
			return nil, fmt.Errorf("creating inventory movement for product %d: %w", item.ProductoID, err)
		}
	}

	return sale, nil
}
