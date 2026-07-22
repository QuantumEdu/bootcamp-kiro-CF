package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/database"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http/middleware"
)

// SaleHandler handles sale-related HTTP requests.
type SaleHandler struct {
	repo *database.SaleRepo
}

// NewSaleHandler creates a new sale handler.
func NewSaleHandler(repo *database.SaleRepo) *SaleHandler {
	return &SaleHandler{repo: repo}
}

// Create processes a new sale from JSON.
func (h *SaleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req entities.CreateSaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Datos de venta invalidos"})
		return
	}

	if len(req.Items) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "La venta debe tener al menos un item"})
		return
	}

	// Validate metodo_pago
	validMethods := map[string]bool{"efectivo": true, "tarjeta": true, "transferencia": true, "mixto": true}
	if !validMethods[req.MetodoPago] {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Metodo de pago invalido"})
		return
	}

	// Get user ID from context (set by auth middleware)
	usuarioID := getUserIDFromContext(r)

	sale, err := h.repo.Create(&req, usuarioID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error creando venta: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, sale)
}

// Recent returns recent sales as an HTML fragment.
func (h *SaleHandler) Recent(w http.ResponseWriter, r *http.Request) {
	sales, err := h.repo.ListRecent(10)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<p class="text-sm text-red-500">Error cargando ventas</p>`)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if len(sales) == 0 {
		fmt.Fprint(w, `<p class="text-sm text-gray-400 text-center py-4">No hay ventas registradas</p>`)
		return
	}

	html := `<div class="space-y-2">`
	for _, s := range sales {
		badgeColor := "bg-green-100 text-green-700"
		switch s.MetodoPago {
		case "tarjeta":
			badgeColor = "bg-blue-100 text-blue-700"
		case "transferencia":
			badgeColor = "bg-purple-100 text-purple-700"
		case "mixto":
			badgeColor = "bg-yellow-100 text-yellow-700"
		}
		html += fmt.Sprintf(`
			<div class="flex items-center justify-between p-2 bg-gray-50 rounded-lg text-sm">
				<div>
					<span class="font-medium text-gray-800">#%d</span>
					<span class="text-gray-500 ml-2">%s</span>
				</div>
				<div class="flex items-center gap-2">
					<span class="px-2 py-0.5 rounded text-xs %s">%s</span>
					<span class="font-bold text-gray-900">$%.2f</span>
				</div>
			</div>
		`, s.ID, s.CreatedAt[:16], badgeColor, s.MetodoPago, s.Total)
	}
	html += `</div>`
	fmt.Fprint(w, html)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func getUserIDFromContext(r *http.Request) int64 {
	// Get user ID from auth middleware context
	if userID, ok := r.Context().Value(middleware.UserIDKey).(int64); ok {
		return userID
	}
	// Default to admin user (ID 1) if no auth context
	return 1
}
