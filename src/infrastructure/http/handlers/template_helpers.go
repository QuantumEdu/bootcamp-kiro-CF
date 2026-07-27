package handlers

import (
	"html/template"
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

// RenderPage renders a full page with layout.html wrapping the specified content template.
// It clones the template set and overrides the "content" block with the specified template.
func RenderPage(w http.ResponseWriter, tmpl *template.Template, contentName string, data map[string]interface{}) error {
	// Clone so we can override "content" without affecting other requests.
	t, err := tmpl.Clone()
	if err != nil {
		return err
	}

	// Look up the content template (try both slash styles for cross-platform).
	ct := tmpl.Lookup(contentName)
	if ct == nil {
		// Content templates define a "content" block — use that directly.
		// If not found, layout.html will render with empty content.
		return t.ExecuteTemplate(w, "layout.html", data)
	}

	// Override "content" in the cloned template with the specific page's tree.
	t.AddParseTree("content", ct.Tree)
	return t.ExecuteTemplate(w, "layout.html", data)
}
