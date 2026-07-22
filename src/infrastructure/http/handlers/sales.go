package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// SaleHandler handles sale-related HTTP requests.
type SaleHandler struct {
	db *sql.DB
}

// NewSaleHandler creates a new sale handler.
func NewSaleHandler(db *sql.DB) *SaleHandler {
	return &SaleHandler{db: db}
}

type createSaleReq struct {
	Items      []saleItemReq `json:"items"`
	MetodoPago string        `json:"metodo_pago"`
}

type saleItemReq struct {
	ProductoID     int64   `json:"producto_id"`
	Cantidad       float64 `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
}

// Create processes a new sale from JSON.
func (h *SaleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createSaleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Datos invalidos"})
		return
	}
	if len(req.Items) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Agrega al menos un item"})
		return
	}

	valid := map[string]bool{"efectivo": true, "tarjeta": true, "transferencia": true, "mixto": true}
	if !valid[req.MetodoPago] {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Metodo de pago invalido"})
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error interno"})
		return
	}
	defer tx.Rollback()

	var total float64
	for _, item := range req.Items {
		total += item.Cantidad * item.PrecioUnitario
	}

	res, err := tx.Exec(`INSERT INTO ventas (usuario_id, total, metodo_pago) VALUES (1, ?, ?)`, total, req.MetodoPago)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error creando venta"})
		return
	}
	ventaID, _ := res.LastInsertId()

	for _, item := range req.Items {
		subtotal := item.Cantidad * item.PrecioUnitario
		_, err := tx.Exec(`INSERT INTO venta_items (venta_id, producto_id, cantidad, precio_unitario, subtotal) VALUES (?, ?, ?, ?, ?)`,
			ventaID, item.ProductoID, item.Cantidad, item.PrecioUnitario, subtotal)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error creando item"})
			return
		}
		tx.Exec(`UPDATE productos SET stock_actual = stock_actual - ?, updated_at = datetime('now','localtime') WHERE id = ?`, item.Cantidad, item.ProductoID)
	}

	if err := tx.Commit(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error finalizando venta"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"id": ventaID, "total": total})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// Recent returns recent sales as an HTML fragment.
func (h *SaleHandler) Recent(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`SELECT id, total, metodo_pago, created_at FROM ventas ORDER BY created_at DESC LIMIT 10`)
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<p class="text-sm text-red-500">Error</p>`)
		return
	}
	defer rows.Close()

	html := `<div class="space-y-2">`
	found := false
	for rows.Next() {
		found = true
		var id int
		var total float64
		var metodo, fecha string
		rows.Scan(&id, &total, &metodo, &fecha)
		badge := "bg-green-100 text-green-700"
		if metodo == "tarjeta" {
			badge = "bg-blue-100 text-blue-700"
		}
		display := fecha
		if len(fecha) > 16 {
			display = fecha[:16]
		}
		html += fmt.Sprintf(`<div class="flex items-center justify-between p-2 bg-gray-50 rounded-lg text-sm"><div><span class="font-medium">#%d</span><span class="text-gray-500 ml-2">%s</span></div><div class="flex items-center gap-2"><span class="px-2 py-0.5 rounded text-xs %s">%s</span><span class="font-bold">$%.2f</span></div></div>`, id, display, badge, metodo, total)
	}
	if !found {
		html += `<p class="text-sm text-gray-400 text-center py-4">No hay ventas</p>`
	}
	html += `</div>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}
