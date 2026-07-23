package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/use_cases"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/alexedwards/scs/v2"
)

// SaleHandler handles sale registration with HTMX support.
type SaleHandler struct {
	registerUC *use_cases.RegisterSale
	tmpl       *template.Template
	sessions   *scs.SessionManager
}

// NewSaleHandler creates a new SaleHandler.
func NewSaleHandler(
	registerUC *use_cases.RegisterSale,
	tmpl *template.Template,
	sessions *scs.SessionManager,
) *SaleHandler {
	return &SaleHandler{
		registerUC: registerUC,
		tmpl:       tmpl,
		sessions:   sessions,
	}
}

// NewSalePage handles GET /ventas/new — renders the POS-style sale capture page.
func (h *SaleHandler) NewSalePage(w http.ResponseWriter, r *http.Request) {
	data := WithUserContext(r, map[string]interface{}{
		"PageTitle": "Nueva Venta",
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "Error de template: "+err.Error(), http.StatusInternalServerError)
	}
}

// cartItemRequest represents a single item in the cart JSON request.
type cartItemRequest struct {
	ProductoID int64   `json:"producto_id"`
	Cantidad   float64 `json:"cantidad"`
}

// checkoutRequest represents the full checkout JSON request.
type checkoutRequest struct {
	Items      []cartItemRequest `json:"items"`
	MetodoPago string            `json:"metodo_pago"`
}

// CompleteSale handles POST /ventas — receives JSON cart items, calls RegisterSale.
func (h *SaleHandler) CompleteSale(w http.ResponseWriter, r *http.Request) {
	// Parse JSON body.
	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.renderSaleError(w, "Error al procesar los datos del carrito")
		return
	}

	if len(req.Items) == 0 {
		h.renderSaleError(w, "El carrito está vacío")
		return
	}

	// Get user ID from session.
	userID := h.sessions.GetInt64(r.Context(), "user_id")
	if userID == 0 {
		// Fallback: use a default user if session not configured.
		userID = 1
	}

	// Build use-case input.
	var items []use_cases.SaleItemInput
	for _, ci := range req.Items {
		items = append(items, use_cases.SaleItemInput{
			ProductoID: ci.ProductoID,
			Cantidad:   ci.Cantidad,
		})
	}

	metodoPago := entities.PaymentMethod(req.MetodoPago)
	if !entities.ValidPaymentMethods[metodoPago] {
		metodoPago = entities.MetodoEfectivo
	}

	input := use_cases.RegisterSaleInput{
		UsuarioID:  userID,
		MetodoPago: metodoPago,
		Items:      items,
	}

	// Execute use case.
	sale, err := h.registerUC.Execute(r.Context(), input)
	if err != nil {
		h.renderSaleError(w, err.Error())
		return
	}

	// Return success fragment.
	h.renderSaleSuccess(w, sale)
}

// renderSaleError returns an HTMX error fragment.
func (h *SaleHandler) renderSaleError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusUnprocessableEntity)
	fmt.Fprintf(w, `<div class="p-3 bg-red-50 border border-red-200 rounded-lg text-sm text-red-700" role="alert">
		<p class="font-medium">Error en la venta</p>
		<p>%s</p>
	</div>`, template.HTMLEscapeString(msg))
}

// renderSaleSuccess returns an HTMX success fragment.
func (h *SaleHandler) renderSaleSuccess(w http.ResponseWriter, sale *entities.Sale) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("HX-Trigger", "ventaCreada")
	fmt.Fprintf(w, `<div class="p-3 bg-green-50 border border-green-200 rounded-lg text-sm text-green-700" role="alert">
		<p class="font-medium">✅ Venta registrada</p>
		<p>Venta #%d — Total: $%.2f — %d items</p>
	</div>`, sale.ID, sale.Total, len(sale.Items))
}
