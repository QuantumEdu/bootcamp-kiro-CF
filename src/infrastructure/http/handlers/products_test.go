package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/use_cases"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

// --- Mock repository ---

type mockProductRepo struct {
	products map[int64]*entities.Product
	nextID   int64
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{
		products: make(map[int64]*entities.Product),
		nextID:   1,
	}
}

func (m *mockProductRepo) Create(_ context.Context, p *entities.Product) error {
	p.ID = m.nextID
	m.nextID++
	m.products[p.ID] = p
	return nil
}

func (m *mockProductRepo) Update(_ context.Context, p *entities.Product) error {
	m.products[p.ID] = p
	return nil
}

func (m *mockProductRepo) FindByID(_ context.Context, id int64) (*entities.Product, error) {
	p, ok := m.products[id]
	if !ok {
		return nil, use_cases.ErrProductNotFound
	}
	return p, nil
}

func (m *mockProductRepo) List(_ context.Context, _ ports.ProductFilter) ([]entities.Product, error) {
	var result []entities.Product
	for _, p := range m.products {
		result = append(result, *p)
	}
	return result, nil
}

func (m *mockProductRepo) Deactivate(_ context.Context, id int64) error {
	p, ok := m.products[id]
	if !ok {
		return use_cases.ErrProductNotFound
	}
	p.Activo = false
	return nil
}

func (m *mockProductRepo) FindLowStock(_ context.Context) ([]entities.Product, error) {
	return nil, nil
}

// --- Helper to build handler ---

func setupProductHandler() (*ProductHandler, *mockProductRepo) {
	repo := newMockProductRepo()
	createUC := use_cases.NewCreateProduct(repo)
	updateUC := use_cases.NewUpdateProduct(repo)
	listUC := use_cases.NewListProducts(repo)
	deactivateUC := use_cases.NewDeactivateProduct(repo)

	handler := NewProductHandler(createUC, updateUC, listUC, deactivateUC, repo, nil)
	return handler, repo
}

func makeFormRequest(values url.Values, htmx bool) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/productos", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if htmx {
		req.Header.Set("HX-Request", "true")
	}
	return req
}

// --- Tests ---

func TestProductCreate_Validation_EmptyName(t *testing.T) {
	handler, _ := setupProductHandler()

	form := url.Values{
		"nombre":       {""},
		"precio_venta": {"25.00"},
	}

	rr := httptest.NewRecorder()
	req := makeFormRequest(form, true)

	handler.Create(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "nombre") {
		t.Errorf("expected error about nombre, got: %s", rr.Body.String())
	}
}

func TestProductCreate_Validation_ZeroPrice(t *testing.T) {
	handler, _ := setupProductHandler()

	form := url.Values{
		"nombre":       {"Coca-Cola"},
		"precio_venta": {"0"},
	}

	rr := httptest.NewRecorder()
	req := makeFormRequest(form, true)

	handler.Create(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "precio") {
		t.Errorf("expected error about precio, got: %s", rr.Body.String())
	}
}

func TestProductCreate_Validation_NegativePrice(t *testing.T) {
	handler, _ := setupProductHandler()

	form := url.Values{
		"nombre":       {"Coca-Cola"},
		"precio_venta": {"-5.00"},
	}

	rr := httptest.NewRecorder()
	req := makeFormRequest(form, true)

	handler.Create(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", rr.Code)
	}
}

func TestProductCreate_Validation_NegativeStock(t *testing.T) {
	handler, _ := setupProductHandler()

	form := url.Values{
		"nombre":       {"Coca-Cola"},
		"precio_venta": {"25.00"},
		"stock_actual": {"-10"},
	}

	rr := httptest.NewRecorder()
	req := makeFormRequest(form, true)

	handler.Create(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "stock") {
		t.Errorf("expected error about stock, got: %s", rr.Body.String())
	}
}

func TestProductCreate_Validation_NegativePurchasePrice(t *testing.T) {
	handler, _ := setupProductHandler()

	form := url.Values{
		"nombre":        {"Coca-Cola"},
		"precio_venta":  {"25.00"},
		"precio_compra": {"-3.00"},
	}

	rr := httptest.NewRecorder()
	req := makeFormRequest(form, true)

	handler.Create(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "compra") {
		t.Errorf("expected error about precio compra, got: %s", rr.Body.String())
	}
}

func TestProductCreate_Validation_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		form       url.Values
		wantStatus int
		wantErrMsg string
	}{
		{
			name: "valid product returns 200",
			form: url.Values{
				"nombre":       {"Coca-Cola 600ml"},
				"precio_venta": {"25.00"},
				"stock_actual": {"50"},
				"stock_minimo": {"10"},
				"unidad":       {"unidad"},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "empty nombre returns 422",
			form: url.Values{
				"nombre":       {""},
				"precio_venta": {"25.00"},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErrMsg: "nombre",
		},
		{
			name: "whitespace-only nombre returns 422",
			form: url.Values{
				"nombre":       {"   "},
				"precio_venta": {"25.00"},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErrMsg: "nombre",
		},
		{
			name: "zero precio_venta returns 422",
			form: url.Values{
				"nombre":       {"Producto"},
				"precio_venta": {"0"},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErrMsg: "precio",
		},
		{
			name: "negative precio_venta returns 422",
			form: url.Values{
				"nombre":       {"Producto"},
				"precio_venta": {"-10"},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErrMsg: "precio",
		},
		{
			name: "non-numeric precio_venta returns 422",
			form: url.Values{
				"nombre":       {"Producto"},
				"precio_venta": {"abc"},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErrMsg: "precio",
		},
		{
			name: "negative stock_actual returns 422",
			form: url.Values{
				"nombre":       {"Producto"},
				"precio_venta": {"25.00"},
				"stock_actual": {"-5"},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErrMsg: "stock",
		},
		{
			name: "negative stock_minimo returns 422",
			form: url.Values{
				"nombre":       {"Producto"},
				"precio_venta": {"25.00"},
				"stock_minimo": {"-1"},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErrMsg: "stock",
		},
		{
			name: "negative precio_compra returns 422",
			form: url.Values{
				"nombre":        {"Producto"},
				"precio_venta":  {"25.00"},
				"precio_compra": {"-2"},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErrMsg: "compra",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, _ := setupProductHandler()
			rr := httptest.NewRecorder()
			req := makeFormRequest(tt.form, true)

			handler.Create(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d. Body: %s", rr.Code, tt.wantStatus, rr.Body.String())
			}
			if tt.wantErrMsg != "" && !strings.Contains(rr.Body.String(), tt.wantErrMsg) {
				t.Errorf("expected body to contain %q, got: %s", tt.wantErrMsg, rr.Body.String())
			}
		})
	}
}

func TestProductCreate_HTMX_ReturnsFragment(t *testing.T) {
	handler, _ := setupProductHandler()

	form := url.Values{
		"nombre":       {"Coca-Cola 600ml"},
		"precio_venta": {"25.00"},
		"stock_actual": {"100"},
		"stock_minimo": {"10"},
		"unidad":       {"unidad"},
	}

	rr := httptest.NewRecorder()
	req := makeFormRequest(form, true)

	handler.Create(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	body := rr.Body.String()
	if !strings.Contains(body, "<tr") {
		t.Errorf("expected HTMX response to contain <tr, got: %s", body)
	}
	if !strings.Contains(body, "Coca-Cola 600ml") {
		t.Errorf("expected product name in response, got: %s", body)
	}
}

func TestProductCreate_NonHTMX_Redirects(t *testing.T) {
	handler, _ := setupProductHandler()

	form := url.Values{
		"nombre":       {"Coca-Cola 600ml"},
		"precio_venta": {"25.00"},
	}

	rr := httptest.NewRecorder()
	req := makeFormRequest(form, false)

	handler.Create(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rr.Code)
	}
	if loc := rr.Header().Get("Location"); loc != "/productos" {
		t.Errorf("expected redirect to /productos, got %s", loc)
	}
}

func TestProductList_HTMX_ReturnsRows(t *testing.T) {
	handler, repo := setupProductHandler()

	// Add a product directly to the repo.
	repo.products[1] = &entities.Product{
		ID:          1,
		Nombre:      "Test Product",
		SKU:         "TST-001",
		PrecioVenta: 15.50,
		StockActual: 20,
		StockMinimo: 5,
		Unidad:      entities.UnitUnidad,
		Activo:      true,
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/productos", nil)
	req.Header.Set("HX-Request", "true")

	handler.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Test Product") {
		t.Errorf("expected product in response, got: %s", body)
	}
	if !strings.Contains(body, "<tr") {
		t.Errorf("expected table rows in HTMX response, got: %s", body)
	}
}
