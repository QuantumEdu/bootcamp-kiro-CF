package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http/middleware"
	"github.com/alexedwards/scs/v2"
)

func newTestSessionManager() *scs.SessionManager {
	sm := scs.New()
	return sm
}

func TestRequireAuth_NoSession_RedirectsToLogin(t *testing.T) {
	sm := newTestSessionManager()

	handler := sm.LoadAndSave(middleware.RequireAuth(sm)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called when no session")
	})))

	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/login" {
		t.Errorf("expected redirect to /login, got %q", loc)
	}
}

func TestRequireAuth_WithSession_CallsNext(t *testing.T) {
	sm := newTestSessionManager()

	var capturedUserID, capturedUserName, capturedUserRole string

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserID, _ = r.Context().Value(middleware.ContextKeyUserID).(string)
		capturedUserName, _ = r.Context().Value(middleware.ContextKeyUserName).(string)
		capturedUserRole, _ = r.Context().Value(middleware.ContextKeyUserRole).(string)
		w.WriteHeader(http.StatusOK)
	})

	// We need to set up a session with user_id. Use a two-step approach:
	// First request sets the session, second request uses it.
	mux := http.NewServeMux()

	// Setup endpoint to create session data.
	mux.HandleFunc("/setup", func(w http.ResponseWriter, r *http.Request) {
		sm.Put(r.Context(), "user_id", "42")
		sm.Put(r.Context(), "user_name", "Ana")
		sm.Put(r.Context(), "user_role", "admin")
		w.WriteHeader(http.StatusOK)
	})

	// Protected endpoint.
	mux.Handle("/protected", middleware.RequireAuth(sm)(inner))

	handler := sm.LoadAndSave(mux)

	// Step 1: create session.
	setupReq := httptest.NewRequest(http.MethodGet, "/setup", nil)
	setupRec := httptest.NewRecorder()
	handler.ServeHTTP(setupRec, setupReq)

	// Extract session cookie.
	cookies := setupRec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected session cookie after setup")
	}

	// Step 2: access protected endpoint with session cookie.
	protectedReq := httptest.NewRequest(http.MethodGet, "/protected", nil)
	for _, c := range cookies {
		protectedReq.AddCookie(c)
	}
	protectedRec := httptest.NewRecorder()
	handler.ServeHTTP(protectedRec, protectedReq)

	if protectedRec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", protectedRec.Code)
	}
	if capturedUserID != "42" {
		t.Errorf("expected user_id=42, got %q", capturedUserID)
	}
	if capturedUserName != "Ana" {
		t.Errorf("expected user_name=Ana, got %q", capturedUserName)
	}
	if capturedUserRole != "admin" {
		t.Errorf("expected user_role=admin, got %q", capturedUserRole)
	}
}
