package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http/middleware"
)

func TestRequireRole_MatchingRole_CallsNext(t *testing.T) {
	sm := newTestSessionManager()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/setup", func(w http.ResponseWriter, r *http.Request) {
		sm.Put(r.Context(), "user_id", "1")
		sm.Put(r.Context(), "user_name", "Admin")
		sm.Put(r.Context(), "user_role", "admin")
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("/admin", middleware.RequireRole(sm, "admin")(inner))

	handler := sm.LoadAndSave(mux)

	// Setup session.
	setupReq := httptest.NewRequest(http.MethodGet, "/setup", nil)
	setupRec := httptest.NewRecorder()
	handler.ServeHTTP(setupRec, setupReq)

	cookies := setupRec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected session cookie")
	}

	// Access admin endpoint.
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestRequireRole_MismatchedRole_Returns403(t *testing.T) {
	sm := newTestSessionManager()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called for wrong role")
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/setup", func(w http.ResponseWriter, r *http.Request) {
		sm.Put(r.Context(), "user_id", "2")
		sm.Put(r.Context(), "user_name", "Cajero")
		sm.Put(r.Context(), "user_role", "cajero")
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("/admin", middleware.RequireRole(sm, "admin")(inner))

	handler := sm.LoadAndSave(mux)

	// Setup session with cajero role.
	setupReq := httptest.NewRequest(http.MethodGet, "/setup", nil)
	setupRec := httptest.NewRecorder()
	handler.ServeHTTP(setupRec, setupReq)

	cookies := setupRec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected session cookie")
	}

	// Access admin endpoint as cajero.
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", rec.Code)
	}
}

func TestRequireRole_NoRole_Returns403(t *testing.T) {
	sm := newTestSessionManager()

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called when no role in session")
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/setup", func(w http.ResponseWriter, r *http.Request) {
		// No role stored, only user_id.
		sm.Put(r.Context(), "user_id", "3")
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("/admin", middleware.RequireRole(sm, "admin")(inner))

	handler := sm.LoadAndSave(mux)

	setupReq := httptest.NewRequest(http.MethodGet, "/setup", nil)
	setupRec := httptest.NewRecorder()
	handler.ServeHTTP(setupRec, setupReq)

	cookies := setupRec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected session cookie")
	}

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", rec.Code)
	}
}
