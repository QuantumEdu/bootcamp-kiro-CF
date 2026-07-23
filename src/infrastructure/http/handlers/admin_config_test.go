package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"html/template"
	"testing"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/services"
	mw "github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http/middleware"
)

// mockConfigRepository is a simple in-memory implementation of ConfigRepository for testing.
type mockConfigRepository struct {
	store map[string]string
}

func newMockConfigRepository() *mockConfigRepository {
	return &mockConfigRepository{store: make(map[string]string)}
}

func (m *mockConfigRepository) Get(_ context.Context, clave string) (string, error) {
	return m.store[clave], nil
}

func (m *mockConfigRepository) Set(_ context.Context, clave, valor string) error {
	m.store[clave] = valor
	return nil
}

// testAdminConfigTemplate returns a minimal template for testing.
func testAdminConfigTemplate() *template.Template {
	tmpl := template.Must(template.New("layout.html").Parse(
		`{{.PageTitle}}|MaskedKey={{.MaskedKey}}|HasKey={{.HasKey}}|Error={{index . "Error"}}`))
	return tmpl
}

func TestAdminConfigHandler_Show_NoKey(t *testing.T) {
	repo := newMockConfigRepository()
	crypto := services.NewCryptoService("test-secret")
	tmpl := testAdminConfigTemplate()
	handler := NewAdminConfigHandler(repo, crypto, tmpl)

	req := httptest.NewRequest(http.MethodGet, "/admin/config", nil)
	ctx := context.WithValue(req.Context(), mw.ContextKeyUserName, "Admin")
	ctx = context.WithValue(ctx, mw.ContextKeyUserRole, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.Show(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "HasKey=false") {
		t.Errorf("expected HasKey=false, got: %s", body)
	}
	if !strings.Contains(body, "MaskedKey=") {
		t.Errorf("expected empty MaskedKey, got: %s", body)
	}
}

func TestAdminConfigHandler_Show_WithKey(t *testing.T) {
	repo := newMockConfigRepository()
	crypto := services.NewCryptoService("test-secret")
	tmpl := testAdminConfigTemplate()
	handler := NewAdminConfigHandler(repo, crypto, tmpl)

	// Store an encrypted key
	encrypted, err := crypto.Encrypt("sk-test-key-1234")
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}
	repo.store["openrouter_api_key"] = encrypted

	req := httptest.NewRequest(http.MethodGet, "/admin/config", nil)
	ctx := context.WithValue(req.Context(), mw.ContextKeyUserName, "Admin")
	ctx = context.WithValue(ctx, mw.ContextKeyUserRole, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.Show(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "HasKey=true") {
		t.Errorf("expected HasKey=true, got: %s", body)
	}
	// Last 4 chars of "sk-test-key-1234" is "1234"
	if !strings.Contains(body, "1234") {
		t.Errorf("expected masked key to contain last 4 chars '1234', got: %s", body)
	}
	// Should have asterisks for the masked portion
	if !strings.Contains(body, "****") {
		t.Errorf("expected masked key to contain asterisks, got: %s", body)
	}
}

func TestAdminConfigHandler_Update_Success(t *testing.T) {
	repo := newMockConfigRepository()
	crypto := services.NewCryptoService("test-secret")
	tmpl := testAdminConfigTemplate()
	handler := NewAdminConfigHandler(repo, crypto, tmpl)

	form := url.Values{}
	form.Set("api_key", "sk-new-api-key-5678")
	req := httptest.NewRequest(http.MethodPost, "/admin/config", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), mw.ContextKeyUserName, "Admin")
	ctx = context.WithValue(ctx, mw.ContextKeyUserRole, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/admin/config" {
		t.Errorf("expected redirect to /admin/config, got %s", loc)
	}
	// Verify the key was stored encrypted
	stored := repo.store["openrouter_api_key"]
	if stored == "" {
		t.Fatal("expected key to be stored")
	}
	// Decrypt and verify round-trip
	decrypted, err := crypto.Decrypt(stored)
	if err != nil {
		t.Fatalf("failed to decrypt stored key: %v", err)
	}
	if decrypted != "sk-new-api-key-5678" {
		t.Errorf("expected 'sk-new-api-key-5678', got '%s'", decrypted)
	}
}

func TestAdminConfigHandler_Update_EmptyKey(t *testing.T) {
	repo := newMockConfigRepository()
	crypto := services.NewCryptoService("test-secret")
	tmpl := testAdminConfigTemplate()
	handler := NewAdminConfigHandler(repo, crypto, tmpl)

	form := url.Values{}
	form.Set("api_key", "   ")
	req := httptest.NewRequest(http.MethodPost, "/admin/config", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), mw.ContextKeyUserName, "Admin")
	ctx = context.WithValue(ctx, mw.ContextKeyUserRole, "admin")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}
