package bootstrap

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPanicRecovery_NoPanic(t *testing.T) {
	handler := PanicRecovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != `{"ok":true}` {
		t.Errorf("expected body {\"ok\":true}, got %s", rec.Body.String())
	}
}

func TestPanicRecovery_StringPanic(t *testing.T) {
	handler := PanicRecovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/sales", nil)
	req.Header.Set("X-Amzn-Trace-Id", "Root=1-abc-def")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if body["error"] != "internal server error" {
		t.Errorf("expected error 'internal server error', got %s", body["error"])
	}
}

func TestPanicRecovery_ErrorPanic(t *testing.T) {
	handler := PanicRecovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(42)
	}))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if body["error"] != "internal server error" {
		t.Errorf("expected error 'internal server error', got %s", body["error"])
	}
}

func TestPanicRecovery_NilPanic(t *testing.T) {
	handler := PanicRecovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(nil)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// panic(nil) in Go 1.21+ is recovered as a *runtime.PanicNilError,
	// so we still expect 500.
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
}
