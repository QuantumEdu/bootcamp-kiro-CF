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
	h.render(w, r, "page.dashboard")
}

// Products renders the products page.
func (h *PageHandler) Products(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "page.products")
}

// Sales renders the sales page.
func (h *PageHandler) Sales(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "page.sales")
}

// Metrics renders the metrics page.
func (h *PageHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, "page.metrics")
}

func (h *PageHandler) render(w http.ResponseWriter, r *http.Request, pageKey string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	contentTmpl := "metrics/dashboard.html"
	switch pageKey {
	case "page.products":
		contentTmpl = "products/list.html"
	case "page.sales":
		contentTmpl = "sales/index.html"
	case "page.metrics", "page.dashboard":
		contentTmpl = "metrics/dashboard.html"
	}

	data := WithUserContext(r, map[string]interface{}{"PageTitle": tr(r, pageKey)})
	if err := RenderPage(w, h.tmpl, contentTmpl, data); err != nil {
		http.Error(w, tr(r, "error.template")+": "+err.Error(), http.StatusInternalServerError)
	}
}
