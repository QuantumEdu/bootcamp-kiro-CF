package handlers

import (
	"html/template"
	"net/http"
)

// PageHandler renders full pages.
type PageHandler struct {
	tmpl *template.Template
}

// NewPageHandler creates a new page handler.
func NewPageHandler(tmpl *template.Template) *PageHandler {
	return &PageHandler{tmpl: tmpl}
}

// Dashboard renders the dashboard.
func (h *PageHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	h.render(w, "Dashboard")
}

// Products renders the products page.
func (h *PageHandler) Products(w http.ResponseWriter, r *http.Request) {
	h.render(w, "Productos")
}

// Sales renders the sales page.
func (h *PageHandler) Sales(w http.ResponseWriter, r *http.Request) {
	h.render(w, "Ventas")
}

// Metrics renders the metrics page.
func (h *PageHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	h.render(w, "Metricas")
}

func (h *PageHandler) render(w http.ResponseWriter, title string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Determine which content template to use
	contentTmpl := "metrics/dashboard.html"
	switch title {
	case "Productos":
		contentTmpl = "products/list.html"
	case "Ventas":
		contentTmpl = "sales/index.html"
	case "Metricas", "Dashboard":
		contentTmpl = "metrics/dashboard.html"
	}
	_ = contentTmpl

	data := map[string]interface{}{"PageTitle": title}
	if err := h.tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}
