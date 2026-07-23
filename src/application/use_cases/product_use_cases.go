package use_cases

import (
	"context"
	"errors"
	"fmt"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

// Domain errors for product use cases.
var (
	ErrProductNotFound = errors.New("producto no encontrado")
)

// --- CreateProduct ---

// CreateProductInput holds the data needed to create a new product.
type CreateProductInput struct {
	Nombre       string
	SKU          string
	CategoriaID  int64
	PrecioVenta  float64
	PrecioCompra float64
	StockActual  float64
	StockMinimo  float64
	Unidad       entities.Unit
}

// CreateProduct handles creating a new product in the catalog.
type CreateProduct struct {
	repo ports.ProductRepository
}

// NewCreateProduct creates a new CreateProduct use case.
func NewCreateProduct(repo ports.ProductRepository) *CreateProduct {
	return &CreateProduct{repo: repo}
}

// Execute validates the product data and persists it.
func (uc *CreateProduct) Execute(ctx context.Context, input CreateProductInput) (*entities.Product, error) {
	product := &entities.Product{
		Nombre:       input.Nombre,
		SKU:          input.SKU,
		CategoriaID:  input.CategoriaID,
		PrecioVenta:  input.PrecioVenta,
		PrecioCompra: input.PrecioCompra,
		StockActual:  input.StockActual,
		StockMinimo:  input.StockMinimo,
		Unidad:       input.Unidad,
		Activo:       true,
	}

	if err := product.Validate(); err != nil {
		return nil, fmt.Errorf("validating product: %w", err)
	}

	if err := uc.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("creating product: %w", err)
	}

	return product, nil
}

// --- UpdateProduct ---

// UpdateProductInput holds the data needed to update an existing product.
type UpdateProductInput struct {
	ID           int64
	Nombre       string
	SKU          string
	CategoriaID  int64
	PrecioVenta  float64
	PrecioCompra float64
	StockActual  float64
	StockMinimo  float64
	Unidad       entities.Unit
}

// UpdateProduct handles updating an existing product.
type UpdateProduct struct {
	repo ports.ProductRepository
}

// NewUpdateProduct creates a new UpdateProduct use case.
func NewUpdateProduct(repo ports.ProductRepository) *UpdateProduct {
	return &UpdateProduct{repo: repo}
}

// Execute finds the product, validates new data, and persists the update.
func (uc *UpdateProduct) Execute(ctx context.Context, input UpdateProductInput) (*entities.Product, error) {
	existing, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	existing.Nombre = input.Nombre
	existing.SKU = input.SKU
	existing.CategoriaID = input.CategoriaID
	existing.PrecioVenta = input.PrecioVenta
	existing.PrecioCompra = input.PrecioCompra
	existing.StockActual = input.StockActual
	existing.StockMinimo = input.StockMinimo
	existing.Unidad = input.Unidad

	if err := existing.Validate(); err != nil {
		return nil, fmt.Errorf("validating product: %w", err)
	}

	if err := uc.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("updating product: %w", err)
	}

	return existing, nil
}

// --- ListProducts ---

// ListProducts handles retrieving products with optional filtering.
type ListProducts struct {
	repo ports.ProductRepository
}

// NewListProducts creates a new ListProducts use case.
func NewListProducts(repo ports.ProductRepository) *ListProducts {
	return &ListProducts{repo: repo}
}

// Execute retrieves products matching the given filter.
func (uc *ListProducts) Execute(ctx context.Context, filter ports.ProductFilter) ([]entities.Product, error) {
	products, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("listing products: %w", err)
	}

	return products, nil
}

// --- DeactivateProduct ---

// DeactivateProduct handles soft-deleting a product.
type DeactivateProduct struct {
	repo ports.ProductRepository
}

// NewDeactivateProduct creates a new DeactivateProduct use case.
func NewDeactivateProduct(repo ports.ProductRepository) *DeactivateProduct {
	return &DeactivateProduct{repo: repo}
}

// Execute marks a product as inactive.
func (uc *DeactivateProduct) Execute(ctx context.Context, id int64) error {
	if err := uc.repo.Deactivate(ctx, id); err != nil {
		return fmt.Errorf("deactivating product: %w", err)
	}

	return nil
}
