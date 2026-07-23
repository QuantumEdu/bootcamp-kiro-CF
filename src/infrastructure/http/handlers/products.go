package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/use_cases"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
	"github.com/go-chi/chi/v5"
)

// ProductHandler handles product CRUD operations with HTMX support.
type ProductHandler struct {
	createUC     *use_cases.CreateProduct
	updateUC     *use_cases.UpdateProduct
	listUC       *use_cases.ListProducts
	deactivateUC *use_cases.DeactivateProduct
	repo         ports.ProductRepository
	tmpl         *template.Template
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(
	create *use_cases.CreateProduct,
	update *use_cases.UpdateProduct,
	list *use_cases.ListProducts,
	deactivate *use_cases.DeactivateProduct,
	repo ports.ProductRepository,
	tmpl *template.Template,
) *ProductHandler {
	return &ProductHandler{
		createUC:     create,
		updateUC:     update,
		listUC:       list,
		deactivateUC: deactivate,
		repo:         repo,
		tmpl:         tmpl,
	}
}

// List handles GET /productos — renders product list page or HTMX table rows fragment.
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	active := true
	filter := ports.ProductFilter{Activo: &active}

	if search := r.URL.Query().Get("search"); search != "" {
		filter.Search = search
	}

	products, err := h.listUC.Execute(r.Context(), filter)
	if err != nil {
		http.Error(w, "Error al cargar productos", http.StatusInternalServerError)
		return
	}

	// If HTMX request, return only table rows fragment.
	if isHTMX(r) {
		h.renderProductRows(w, products)
		return
	}

	// Full page render.
	data := map[string]interface{}{
		"PageTitle": "Productos",
		"Products":  products,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "Error de template: "+err.Error(), http.StatusInternalServerError)
	}
}

// CreateForm handles GET /productos/new — renders the create product form.
func (h *ProductHandler) CreateForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"PageTitle": "Nuevo Producto",
		"Product":   &entities.Product{Unidad: entities.UnitUnidad, Activo: true},
		"IsEdit":    false,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if t := h.tmpl.Lookup("products/form.html"); t != nil {
		if err := t.Execute(w, data); err != nil {
			http.Error(w, "Error de template", http.StatusInternalServerError)
		}
		return
	}

	// Fallback inline form if template not found.
	h.renderFormFallback(w, data)
}

// Create handles POST /productos — validates and creates a new product.
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderValidationError(w, r, "Error al procesar el formulario")
		return
	}

	input, errMsg := parseProductForm(r)
	if errMsg != "" {
		h.renderValidationError(w, r, errMsg)
		return
	}

	product, err := h.createUC.Execute(r.Context(), input)
	if err != nil {
		h.renderValidationError(w, r, formatUseCaseError(err))
		return
	}

	if isHTMX(r) {
		h.renderProductRow(w, product)
		return
	}

	http.Redirect(w, r, "/productos", http.StatusSeeOther)
}

// EditForm handles GET /productos/{id}/edit — renders the edit form with current values.
func (h *ProductHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	product, err := h.repo.FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"PageTitle": "Editar Producto",
		"Product":   product,
		"IsEdit":    true,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if t := h.tmpl.Lookup("products/form.html"); t != nil {
		if err := t.Execute(w, data); err != nil {
			http.Error(w, "Error de template", http.StatusInternalServerError)
		}
		return
	}

	h.renderFormFallback(w, data)
}

// Edit handles POST /productos/{id} — validates and updates an existing product.
func (h *ProductHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.renderValidationError(w, r, "ID inválido")
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderValidationError(w, r, "Error al procesar el formulario")
		return
	}

	createInput, errMsg := parseProductForm(r)
	if errMsg != "" {
		h.renderValidationError(w, r, errMsg)
		return
	}

	updateInput := use_cases.UpdateProductInput{
		ID:           id,
		Nombre:       createInput.Nombre,
		SKU:          createInput.SKU,
		CategoriaID:  createInput.CategoriaID,
		PrecioVenta:  createInput.PrecioVenta,
		PrecioCompra: createInput.PrecioCompra,
		StockActual:  createInput.StockActual,
		StockMinimo:  createInput.StockMinimo,
		Unidad:       createInput.Unidad,
	}

	product, err := h.updateUC.Execute(r.Context(), updateInput)
	if err != nil {
		h.renderValidationError(w, r, formatUseCaseError(err))
		return
	}

	if isHTMX(r) {
		h.renderProductRow(w, product)
		return
	}

	http.Redirect(w, r, "/productos", http.StatusSeeOther)
}

// Deactivate handles DELETE /productos/{id} — deactivates a product.
func (h *ProductHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.deactivateUC.Execute(r.Context(), id); err != nil {
		http.Error(w, "Error al desactivar producto", http.StatusInternalServerError)
		return
	}

	if isHTMX(r) {
		// Return empty response — HTMX will remove the row via hx-swap="delete".
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/productos", http.StatusSeeOther)
}

// --- Helpers ---

// isHTMX checks if the request comes from HTMX.
func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// parseProductForm extracts and validates form input into a CreateProductInput.
// Returns the input and an error message (empty if valid).
func parseProductForm(r *http.Request) (use_cases.CreateProductInput, string) {
	nombre := strings.TrimSpace(r.FormValue("nombre"))
	if nombre == "" {
		return use_cases.CreateProductInput{}, "El nombre del producto es obligatorio"
	}

	precioVenta, err := strconv.ParseFloat(r.FormValue("precio_venta"), 64)
	if err != nil || precioVenta <= 0 {
		return use_cases.CreateProductInput{}, "El precio de venta debe ser mayor a cero"
	}

	precioCompra, _ := strconv.ParseFloat(r.FormValue("precio_compra"), 64)
	if precioCompra < 0 {
		return use_cases.CreateProductInput{}, "El precio de compra no puede ser negativo"
	}

	stockActual, _ := strconv.ParseFloat(r.FormValue("stock_actual"), 64)
	if stockActual < 0 {
		return use_cases.CreateProductInput{}, "El stock actual no puede ser negativo"
	}

	stockMinimo, _ := strconv.ParseFloat(r.FormValue("stock_minimo"), 64)
	if stockMinimo < 0 {
		return use_cases.CreateProductInput{}, "El stock mínimo no puede ser negativo"
	}

	categoriaID, _ := strconv.ParseInt(r.FormValue("categoria_id"), 10, 64)

	unidad := entities.Unit(r.FormValue("unidad"))
	if unidad == "" {
		unidad = entities.UnitUnidad
	}

	return use_cases.CreateProductInput{
		Nombre:       nombre,
		SKU:          strings.TrimSpace(r.FormValue("sku")),
		CategoriaID:  categoriaID,
		PrecioVenta:  precioVenta,
		PrecioCompra: precioCompra,
		StockActual:  stockActual,
		StockMinimo:  stockMinimo,
		Unidad:       unidad,
	}, ""
}

// formatUseCaseError extracts a user-friendly message from use-case errors.
func formatUseCaseError(err error) string {
	msg := err.Error()
	// Strip "validating product: " prefix for cleaner display.
	if strings.Contains(msg, "validating product: ") {
		parts := strings.SplitN(msg, "validating product: ", 2)
		if len(parts) == 2 {
			return parts[1]
		}
	}
	return msg
}

// renderValidationError returns a 422 response with the error message.
// For HTMX requests, returns an inline error fragment.
func (h *ProductHandler) renderValidationError(w http.ResponseWriter, r *http.Request, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusUnprocessableEntity)

	if isHTMX(r) {
		fmt.Fprintf(w, `<div class="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700" role="alert">%s</div>`, template.HTMLEscapeString(msg))
		return
	}

	fmt.Fprintf(w, `<p class="text-red-600">%s</p>`, template.HTMLEscapeString(msg))
}

// renderProductRows renders table rows for a list of products (HTMX fragment).
func (h *ProductHandler) renderProductRows(w http.ResponseWriter, products []entities.Product) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if len(products) == 0 {
		fmt.Fprint(w, `<tr><td colspan="5" class="px-4 py-8 text-center text-gray-400">No hay productos</td></tr>`)
		return
	}

	for i := range products {
		h.writeProductRow(w, &products[i])
	}
}

// renderProductRow renders a single product table row (HTMX fragment).
func (h *ProductHandler) renderProductRow(w http.ResponseWriter, product *entities.Product) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.writeProductRow(w, product)
}

// writeProductRow writes a single <tr> for a product.
func (h *ProductHandler) writeProductRow(w http.ResponseWriter, p *entities.Product) {
	stockClass := "text-gray-900"
	if p.IsLowStock() {
		stockClass = "text-red-600 font-medium"
	}

	fmt.Fprintf(w, `<tr id="product-%d" class="hover:bg-gray-50">
<td class="px-4 py-3 font-medium text-gray-900">%s</td>
<td class="px-4 py-3 text-gray-500">%s</td>
<td class="px-4 py-3 text-right text-gray-900">$%.2f</td>
<td class="px-4 py-3 text-right %s">%.0f</td>
<td class="px-4 py-3 text-right">
<button hx-get="/productos/%d/edit" hx-target="#product-form-area" hx-swap="innerHTML" class="text-indigo-600 hover:text-indigo-800 text-sm mr-2">Editar</button>
<button hx-delete="/productos/%d" hx-target="#product-%d" hx-swap="outerHTML" hx-confirm="¿Desactivar este producto?" class="text-red-500 hover:text-red-700 text-sm">Desactivar</button>
</td>
</tr>`,
		p.ID,
		template.HTMLEscapeString(p.Nombre),
		template.HTMLEscapeString(p.SKU),
		p.PrecioVenta,
		stockClass, p.StockActual,
		p.ID,
		p.ID, p.ID,
	)
}

// renderFormFallback renders a minimal product form when the template is not found.
func (h *ProductHandler) renderFormFallback(w http.ResponseWriter, data map[string]interface{}) {
	product := data["Product"].(*entities.Product)
	isEdit := data["IsEdit"].(bool)

	action := "/productos"
	title := "Nuevo Producto"
	if isEdit {
		action = fmt.Sprintf("/productos/%d", product.ID)
		title = "Editar Producto"
	}

	html := fmt.Sprintf(`<form method="POST" action="%s" hx-post="%s" hx-target="#product-list" hx-swap="afterbegin" class="space-y-4 p-4 bg-white rounded-xl shadow-sm border border-gray-100">
<h3 class="text-lg font-bold text-gray-800">%s</h3>
<div id="form-errors"></div>
<div class="grid grid-cols-2 gap-4">
<div><label class="block text-sm font-medium text-gray-700 mb-1">Nombre *</label><input type="text" name="nombre" value="%s" required class="w-full rounded-lg border-gray-300 shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm p-2 border"/></div>
<div><label class="block text-sm font-medium text-gray-700 mb-1">SKU</label><input type="text" name="sku" value="%s" class="w-full rounded-lg border-gray-300 shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm p-2 border"/></div>
<div><label class="block text-sm font-medium text-gray-700 mb-1">Precio venta *</label><input type="number" name="precio_venta" value="%.2f" step="0.01" min="0.01" required class="w-full rounded-lg border-gray-300 shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm p-2 border"/></div>
<div><label class="block text-sm font-medium text-gray-700 mb-1">Precio compra</label><input type="number" name="precio_compra" value="%.2f" step="0.01" min="0" class="w-full rounded-lg border-gray-300 shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm p-2 border"/></div>
<div><label class="block text-sm font-medium text-gray-700 mb-1">Stock actual</label><input type="number" name="stock_actual" value="%.0f" step="0.01" min="0" class="w-full rounded-lg border-gray-300 shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm p-2 border"/></div>
<div><label class="block text-sm font-medium text-gray-700 mb-1">Stock mínimo</label><input type="number" name="stock_minimo" value="%.0f" step="0.01" min="0" class="w-full rounded-lg border-gray-300 shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm p-2 border"/></div>
<div><label class="block text-sm font-medium text-gray-700 mb-1">Unidad</label><select name="unidad" class="w-full rounded-lg border-gray-300 shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm p-2 border"><option value="unidad"%s>Unidad</option><option value="kg"%s>Kg</option><option value="litro"%s>Litro</option><option value="paquete"%s>Paquete</option></select></div>
<div><label class="block text-sm font-medium text-gray-700 mb-1">Categoría ID</label><input type="number" name="categoria_id" value="%d" min="0" class="w-full rounded-lg border-gray-300 shadow-sm focus:ring-indigo-500 focus:border-indigo-500 text-sm p-2 border"/></div>
</div>
<div class="flex gap-2"><button type="submit" class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 text-sm font-medium">Guardar</button><button type="button" onclick="document.getElementById('product-form-area').innerHTML=''" class="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 text-sm font-medium">Cancelar</button></div>
</form>`,
		action, action, title,
		template.HTMLEscapeString(product.Nombre),
		template.HTMLEscapeString(product.SKU),
		product.PrecioVenta, product.PrecioCompra,
		product.StockActual, product.StockMinimo,
		selected(product.Unidad, entities.UnitUnidad),
		selected(product.Unidad, entities.UnitKg),
		selected(product.Unidad, entities.UnitLitro),
		selected(product.Unidad, entities.UnitPaquete),
		product.CategoriaID,
	)
	fmt.Fprint(w, html)
}

// selected returns ` selected` if the values match.
func selected(current, target entities.Unit) string {
	if current == target {
		return " selected"
	}
	return ""
}
