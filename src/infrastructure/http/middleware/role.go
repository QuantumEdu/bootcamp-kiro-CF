package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

// RequireRole returns middleware that checks if the session user has the required role.
// This middleware should be applied AFTER RequireAuth to ensure user_role is in the session.
// If the role does not match, it responds with 403 Forbidden.
func RequireRole(sessions *scs.SessionManager, role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := sessions.GetString(r.Context(), "user_role")
			if userRole != role {
				http.Error(w, "Acceso denegado: no tiene permisos para esta sección.", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r.WithContext(r.Context()))
		})
	}
}
