package handlers

import (
	"net/http"

	mw "github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http/middleware"
)

// WithUserContext enriches template data with user information from the request context.
// It adds UserName, UserRole, and UserInitial fields used by layout.html sidebar.
func WithUserContext(r *http.Request, data map[string]interface{}) map[string]interface{} {
	userName, _ := r.Context().Value(mw.ContextKeyUserName).(string)
	userRole, _ := r.Context().Value(mw.ContextKeyUserRole).(string)

	initial := "U"
	if len(userName) > 0 {
		initial = string([]rune(userName)[:1])
	}

	data["UserName"] = userName
	data["UserRole"] = userRole
	data["UserInitial"] = initial
	return data
}
