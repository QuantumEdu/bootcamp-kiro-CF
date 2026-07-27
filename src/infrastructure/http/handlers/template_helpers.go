package handlers

import (
	"html/template"
	"net/http"
	"strings"

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
// Instead of cloning (which fails after first execution), we execute the content template
// directly by looking it up and rendering it within the layout context.
func RenderPage(w http.ResponseWriter, tmpl *template.Template, contentName string, data map[string]interface{}) error {
	// Execute the layout template. The {{template "content" .}} call inside layout.html
	// will look for a "content" definition. We need to ensure the right one is active.
	//
	// Strategy: Execute the specific content template first into a buffer,
	// then pass it as HTML data to the layout.
	ct := tmpl.Lookup(contentName)
	if ct == nil {
		// Fallback: try the "content" defined block directly
		return tmpl.ExecuteTemplate(w, "layout.html", data)
	}

	// Render the content template to get its HTML
	var buf strings.Builder
	if err := ct.Execute(&buf, data); err != nil {
		return err
	}

	// Add rendered content as safe HTML to the data
	data["Content"] = template.HTML(buf.String())
	return tmpl.ExecuteTemplate(w, "layout.html", data)
}
