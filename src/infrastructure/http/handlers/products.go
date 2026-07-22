package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/database"
)

// ProductHandler handles product-related HTTP requests.
type ProductHandler struct {
	repo *database.ProductRepo
	tmpl *template.Template
}

// NewProductHandler creates a new product handler.
func NewProductHandler(repo *database.ProductRepo, tmpl *template.Template) *ProductHandler {
	return &ProductHandler{repo: repo, tmpl: tmpl}
}

// List renders the product list page.
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("q")
	products, err := h.repo.List(search)
	if err != nil {
		http.Error(w, "Error loading products", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"PageTitle": "Productos",
		"Products":  products,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Header.Get("HX-Request") == "true" {
		h.tmpl.ExecuteTemplate(w, "products/list.html", data)
		return
	}
	h.tmpl.ExecuteTemplate(w, "layout.html", data)
}

// Form renders the product create/edit form.
func (h *ProductHandler) Form(w http.ResponseWriter, r *http.Request) {
	categories, _ := h.repo.ListCategories()

	product := &entities.Product{Unidad: "unidad"}

	// Check if editing
	if idStr := chi.URLParam(r, "id"); idStr != "" {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		if p, err := h.repo.GetByID(id); err == nil {
			product = p
		}
	}

	data := map[string]interface{}{
		"PageTitle":  "Producto",
		"Product":    product,
		"Categories": categories,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Header.Get("HX-Request") == "true" {
		h.tmpl.ExecuteTemplate(w, "products/form.html", data)
		return
	}
	h.tmpl.ExecuteTemplate(w, "layout.html", data)
}

// Create handles product creation/update.
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	precioVenta, _ := strconv.ParseFloat(r.FormValue("precio_venta"), 64)
	precioCompra, _ := strconv.ParseFloat(r.FormValue("precio_compra"), 64)
	stockActual, _ := strconv.ParseFloat(r.FormValue("stock_actual"), 64)
	stockMinimo, _ := strconv.ParseFloat(r.FormValue("stock_minimo"), 64)

	var catID *int64
	if cid := r.FormValue("categoria_id"); cid != "" {
		id, _ := strconv.ParseInt(cid, 10, 64)
		if id > 0 {
			catID = &id
		}
	}

	product := &entities.Product{
		Nombre:      r.FormValue("nombre"),
		SKU:         r.FormValue("sku"),
		CategoriaID: catID,
		PrecioVenta: precioVenta,
		PrecioCompra: precioCompra,
		StockActual: stockActual,
		StockMinimo: stockMinimo,
		Unidad:      r.FormValue("unidad"),
	}

	if product.Unidad == "" {
		product.Unidad = "unidad"
	}

	// Check if update
	if idStr := r.FormValue("id"); idStr != "" {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		if id > 0 {
			product.ID = id
			if err := h.repo.Update(product); err != nil {
				http.Error(w, "Error updating product", http.StatusInternalServerError)
				return
			}
		}
	} else {
		if _, err := h.repo.Create(product); err != nil {
			http.Error(w, "Error creating product", http.StatusInternalServerError)
			return
		}
	}

	// Redirect to list
	h.List(w, r)
}

// Search returns product search results as HTML fragments (for sales).
func (h *ProductHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "")
		return
	}

	products, err := h.repo.Search(q)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<p class="text-sm text-red-500">Error buscando productos</p>`)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if len(products) == 0 {
		fmt.Fprint(w, `<p class="text-sm text-gray-400 py-2">No se encontraron productos</p>`)
		return
	}

	html := ""
	for _, p := range products {
		html += fmt.Sprintf(`
			<button type="button" onclick="addToCart(%d, '%s', %.2f)"
			        class="w-full flex items-center justify-between p-2 hover:bg-indigo-50 rounded-lg transition-colors text-left">
				<div>
					<p class="text-sm font-medium text-gray-800">%s</p>
					<p class="text-xs text-gray-500">SKU: %s · Stock: %.0f %s</p>
				</div>
				<span class="text-sm font-bold text-indigo-600">$%.2f</span>
			</button>
		`, p.ID, escapeJS(p.Nombre), p.PrecioVenta, p.Nombre, p.SKU, p.StockActual, p.Unidad, p.PrecioVenta)
	}
	fmt.Fprint(w, html)
}

func escapeJS(s string) string {
	s = fmt.Sprintf("%s", s)
	s = replaceAll(s, "'", "\\'")
	s = replaceAll(s, "\"", "\\\"")
	return s
}

func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old) - 1
		} else {
			result += string(s[i])
		}
	}
	return result
}
