package middleware

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

// contextKey is an unexported type for context keys in this package.
type contextKey string

// Context keys for extracting user information in handlers.
const (
	ContextKeyUserID   contextKey = "user_id"
	ContextKeyUserName contextKey = "user_name"
	ContextKeyUserRole contextKey = "user_role"
)

// RequireAuth returns middleware that checks for an authenticated session.
// If user_id is missing from the session, it redirects to /login.
// If present, it adds user_id, user_name, and user_role to the request context.
func RequireAuth(sessions *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := sessions.GetString(r.Context(), "user_id")
			if userID == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Add user info to context for downstream handlers.
			userName := sessions.GetString(r.Context(), "user_name")
			userRole := sessions.GetString(r.Context(), "user_role")

			ctx := context.WithValue(r.Context(), ContextKeyUserID, userID)
			ctx = context.WithValue(ctx, ContextKeyUserName, userName)
			ctx = context.WithValue(ctx, ContextKeyUserRole, userRole)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
