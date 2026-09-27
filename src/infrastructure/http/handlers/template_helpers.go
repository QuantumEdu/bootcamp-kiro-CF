package handlers

import (
	"html/template"
	"net/http"
	"strings"

	mw "github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http/middleware"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/i18n"
)

// WithUserContext enriches template data with user information from the request context.
// It adds UserName, UserRole, UserInitial, Lang, and T (translation func) used by layout.html.
func WithUserContext(r *http.Request, data map[string]interface{}) map[string]interface{} {
	userName, _ := r.Context().Value(mw.ContextKeyUserName).(string)
	userRole, _ := r.Context().Value(mw.ContextKeyUserRole).(string)

	initial := "U"
	if len(userName) > 0 {
		initial = string([]rune(userName)[:1])
	}

	tr := i18n.FromRequest(r)

	data["UserName"] = userName
	data["UserRole"] = userRole
	data["UserRoleLabel"] = userRole
	if userRole == "admin" || userRole == "cajero" {
		data["UserRoleLabel"] = tr.Translate("role." + userRole)
	}
	data["UserInitial"] = initial
	data["Lang"] = tr.Lang()
	data["T"] = tr.Func()
	return data
}

// RenderPage renders a full page with layout.html wrapping the specified content template.
// Content templates are standalone HTML fragments (no {{define}} blocks).
// We render them into a buffer and pass as safe HTML to layout.html via .Content field.
func RenderPage(w http.ResponseWriter, tmpl *template.Template, contentName string, data map[string]interface{}) error {
	ct := tmpl.Lookup(contentName)
	if ct == nil {
		return tmpl.ExecuteTemplate(w, "layout.html", data)
	}

	// Render the content template to a buffer
	var buf strings.Builder
	if err := ct.Execute(&buf, data); err != nil {
		return err
	}

	// Pass rendered content as safe HTML to layout
	data["Content"] = template.HTML(buf.String())
	return tmpl.ExecuteTemplate(w, "layout.html", data)
}
