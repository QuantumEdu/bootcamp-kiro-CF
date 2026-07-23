package handlers

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/use_cases"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	mw "github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http/middleware"
)

// mockClientRepository is a simple in-memory implementation of ClientRepository for testing.
type mockClientRepository struct {
	clients []entities.Client
}

func (m *mockClientRepository) Create(_ context.Context, client *entities.Client) error {
	client.ID = int64(len(m.clients) + 1)
	m.clients = append(m.clients, *client)
	return nil
}

func (m *mockClientRepository) List(_ context.Context) ([]entities.Client, error) {
	return m.clients, nil
}

// testClientTemplate returns a minimal template for testing ClientHandler.
func testClientTemplate() *template.Template {
	tmpl := template.Must(template.New("layout.html").Parse(
		`{{.PageTitle}}|Clients={{range .Clients}}{{.Nombre}},{{end}}|Error={{index . "Error"}}|UserName={{.UserName}}`))
	return tmpl
}

func TestClientHandler_List_ReturnsClients(t *testing.T) {
	repo := &mockClientRepository{
		clients: []entities.Client{
			{ID: 1, Nombre: "Juan Pérez", Telefono: "555-1234", Direccion: "Calle 1"},
			{ID: 2, Nombre: "María López", Telefono: "555-5678", Direccion: "Calle 2"},
		},
	}
	listUC := use_cases.NewListClients(repo)
	createUC := use_cases.NewCreateClient(repo)
	tmpl := testClientTemplate()
	handler := NewClientHandler(createUC, listUC, tmpl)

	req := httptest.NewRequest(http.MethodGet, "/clientes", nil)
	ctx := context.WithValue(req.Context(), mw.ContextKeyUserName, "TestUser")
	ctx = context.WithValue(ctx, mw.ContextKeyUserRole, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Juan Pérez") {
		t.Errorf("expected body to contain 'Juan Pérez', got: %s", body)
	}
	if !strings.Contains(body, "María López") {
		t.Errorf("expected body to contain 'María López', got: %s", body)
	}
	if !strings.Contains(body, "Clientes") {
		t.Errorf("expected body to contain PageTitle 'Clientes', got: %s", body)
	}
}

func TestClientHandler_Create_ValidInput_Redirects(t *testing.T) {
	repo := &mockClientRepository{}
	listUC := use_cases.NewListClients(repo)
	createUC := use_cases.NewCreateClient(repo)
	tmpl := testClientTemplate()
	handler := NewClientHandler(createUC, listUC, tmpl)

	form := url.Values{}
	form.Set("nombre", "Carlos García")
	form.Set("telefono", "555-9999")
	form.Set("direccion", "Av. Principal 100")
	req := httptest.NewRequest(http.MethodPost, "/clientes", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), mw.ContextKeyUserName, "TestUser")
	ctx = context.WithValue(ctx, mw.ContextKeyUserRole, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/clientes" {
		t.Errorf("expected redirect to /clientes, got %s", loc)
	}
	// Verify client was persisted in the mock repo
	if len(repo.clients) != 1 {
		t.Fatalf("expected 1 client in repo, got %d", len(repo.clients))
	}
	if repo.clients[0].Nombre != "Carlos García" {
		t.Errorf("expected nombre 'Carlos García', got '%s'", repo.clients[0].Nombre)
	}
}

func TestClientHandler_Create_EmptyName_Returns422(t *testing.T) {
	repo := &mockClientRepository{}
	listUC := use_cases.NewListClients(repo)
	createUC := use_cases.NewCreateClient(repo)
	tmpl := testClientTemplate()
	handler := NewClientHandler(createUC, listUC, tmpl)

	form := url.Values{}
	form.Set("nombre", "   ")
	form.Set("telefono", "555-0000")
	form.Set("direccion", "Calle Vacía")
	req := httptest.NewRequest(http.MethodPost, "/clientes", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), mw.ContextKeyUserName, "TestUser")
	ctx = context.WithValue(ctx, mw.ContextKeyUserRole, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
	// Verify no client was persisted
	if len(repo.clients) != 0 {
		t.Errorf("expected 0 clients in repo after invalid input, got %d", len(repo.clients))
	}
}
