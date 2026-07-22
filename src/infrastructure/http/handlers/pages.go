package handlers

import (
	"html/template"
	"net/http"
)

// PageHandler handles page rendering.
type PageHandler struct {
	tmpl *template.Template
}

// NewPageHandler creates a new page handler.
func NewPageHandler(tmpl *template.Template) *PageHandler {
	return &PageHandler{tmpl: tmpl}
}

// Dashboard renders the main dashboard page.
func (h *PageHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, "metrics/dashboard.html", map[string]interface{}{
		"PageTitle": "Dashboard",
	})
}

// Products renders the products list page.
func (h *PageHandler) Products(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, "products/list.html", map[string]interface{}{
		"PageTitle": "Productos",
		"Products":  []interface{}{},
	})
}

// ProductForm renders the product create/edit form.
func (h *PageHandler) ProductForm(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, "products/form.html", map[string]interface{}{
		"PageTitle":  "Nuevo Producto",
		"Product":    map[string]interface{}{"ID": 0, "Nombre": "", "SKU": "", "PrecioVenta": 0, "PrecioCompra": 0, "StockActual": 0, "StockMinimo": 0, "Unidad": "unidad"},
		"Categories": []interface{}{},
	})
}

// Sales renders the sales page.
func (h *PageHandler) Sales(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, "sales/index.html", map[string]interface{}{
		"PageTitle": "Ventas",
	})
}

// Metrics renders the metrics dashboard page.
func (h *PageHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	h.renderPage(w, "metrics/dashboard.html", map[string]interface{}{
		"PageTitle": "Metricas",
	})
}

func (h *PageHandler) renderPage(w http.ResponseWriter, templateName string, data map[string]interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Check if this is an HTMX request (partial) — check request not response header
	// Note: we can't read this here without the request, so full render always.
	// Full page render with layout
	if err := h.tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}
