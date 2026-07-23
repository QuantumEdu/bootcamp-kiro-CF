package middleware_test

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http/middleware"
)

// **Validates: Requirements 6.2**
// Property 4: Non-admin users cannot access admin routes
// For any authenticated user whose role is not "admin", requests to the admin config
// route SHALL receive a 403 Forbidden response from the RequireRole("admin") middleware.

// generateRandomRole creates a random non-admin role string.
func generateRandomRole(rng *rand.Rand) string {
	// Pool of realistic and random roles to test.
	knownRoles := []string{
		"cashier", "viewer", "manager", "supervisor",
		"cajero", "mesero", "gerente", "cocinero",
		"guest", "user", "moderator", "support",
		"readonly", "editor", "operator", "analyst",
		"", // empty role
	}

	// 50% chance to use a known role, 50% to generate a random string.
	if rng.Intn(2) == 0 {
		return knownRoles[rng.Intn(len(knownRoles))]
	}

	// Generate a random string of length 1-20.
	length := rng.Intn(20) + 1
	chars := "abcdefghijklmnopqrstuvwxyz0123456789_-"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rng.Intn(len(chars))]
	}
	role := string(result)

	// Ensure we never accidentally generate "admin".
	if role == "admin" {
		return "not-admin"
	}
	return role
}

func TestRequireRole_Property_NonAdminGets403(t *testing.T) {
	const iterations = 150

	rng := rand.New(rand.NewSource(42))

	for i := 0; i < iterations; i++ {
		role := generateRandomRole(rng)

		t.Run("non-admin-role", func(t *testing.T) {
			sm := newTestSessionManager()

			inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Fatalf("handler should NOT be called for non-admin role %q", role)
			})

			mux := http.NewServeMux()
			mux.HandleFunc("/setup", func(w http.ResponseWriter, r *http.Request) {
				sm.Put(r.Context(), "user_id", "99")
				sm.Put(r.Context(), "user_name", "TestUser")
				sm.Put(r.Context(), "user_role", role)
				w.WriteHeader(http.StatusOK)
			})
			mux.Handle("/admin/config", middleware.RequireRole(sm, "admin")(inner))

			handler := sm.LoadAndSave(mux)

			// Step 1: Set up session with non-admin role.
			setupReq := httptest.NewRequest(http.MethodGet, "/setup", nil)
			setupRec := httptest.NewRecorder()
			handler.ServeHTTP(setupRec, setupReq)

			cookies := setupRec.Result().Cookies()
			if len(cookies) == 0 {
				t.Fatal("expected session cookie")
			}

			// Step 2: Attempt to access admin route.
			req := httptest.NewRequest(http.MethodGet, "/admin/config", nil)
			for _, c := range cookies {
				req.AddCookie(c)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusForbidden {
				t.Errorf("iteration %d: role=%q expected 403, got %d", i, role, rec.Code)
			}
		})
	}
}

func TestRequireRole_Property_AdminGets200(t *testing.T) {
	// Positive case: admin role should always get 200.
	const iterations = 50

	for i := 0; i < iterations; i++ {
		t.Run("admin-role", func(t *testing.T) {
			sm := newTestSessionManager()

			inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			mux := http.NewServeMux()
			mux.HandleFunc("/setup", func(w http.ResponseWriter, r *http.Request) {
				sm.Put(r.Context(), "user_id", "1")
				sm.Put(r.Context(), "user_name", "AdminUser")
				sm.Put(r.Context(), "user_role", "admin")
				w.WriteHeader(http.StatusOK)
			})
			mux.Handle("/admin/config", middleware.RequireRole(sm, "admin")(inner))

			handler := sm.LoadAndSave(mux)

			// Step 1: Set up session with admin role.
			setupReq := httptest.NewRequest(http.MethodGet, "/setup", nil)
			setupRec := httptest.NewRecorder()
			handler.ServeHTTP(setupRec, setupReq)

			cookies := setupRec.Result().Cookies()
			if len(cookies) == 0 {
				t.Fatal("expected session cookie")
			}

			// Step 2: Access admin route.
			req := httptest.NewRequest(http.MethodGet, "/admin/config", nil)
			for _, c := range cookies {
				req.AddCookie(c)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("iteration %d: admin role expected 200, got %d", i, rec.Code)
			}
		})
	}
}
