package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

// MetricsHandler handles metric fragment rendering for HTMX polling.
type MetricsHandler struct {
	db *sql.DB
}

// NewMetricsHandler creates a new metrics handler.
func NewMetricsHandler(db *sql.DB) *MetricsHandler {
	return &MetricsHandler{db: db}
}

// VentasHoy returns today's sales summary as an HTML fragment.
func (h *MetricsHandler) VentasHoy(w http.ResponseWriter, r *http.Request) {
	var count int
	var total float64

	err := h.db.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(total), 0) 
		FROM ventas 
		WHERE DATE(created_at) = DATE('now', 'localtime')
	`).Scan(&count, &total)

	if err != nil {
		h.renderError(w, "Error cargando ventas de hoy")
		return
	}

	html := fmt.Sprintf(`
		<p class="text-sm font-medium text-gray-500">Ventas hoy</p>
		<p class="text-2xl font-bold text-gray-900 mt-1">$%s</p>
		<p class="text-xs text-gray-400 mt-1">%d transacciones</p>
	`, formatMoney(total), count)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// VentasSemana returns this week's sales summary.
func (h *MetricsHandler) VentasSemana(w http.ResponseWriter, r *http.Request) {
	var count int
	var total float64

	err := h.db.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(total), 0) 
		FROM ventas 
		WHERE created_at >= datetime('now', '-7 days', 'localtime')
	`).Scan(&count, &total)

	if err != nil {
		h.renderError(w, "Error cargando ventas de la semana")
		return
	}

	html := fmt.Sprintf(`
		<p class="text-sm font-medium text-gray-500">Ventas esta semana</p>
		<p class="text-2xl font-bold text-gray-900 mt-1">$%s</p>
		<p class="text-xs text-gray-400 mt-1">%d transacciones (ultimos 7 dias)</p>
	`, formatMoney(total), count)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// VentasMes returns this month's sales summary.
func (h *MetricsHandler) VentasMes(w http.ResponseWriter, r *http.Request) {
	var count int
	var total float64

	err := h.db.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(total), 0) 
		FROM ventas 
		WHERE strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now', 'localtime')
	`).Scan(&count, &total)

	if err != nil {
		h.renderError(w, "Error cargando ventas del mes")
		return
	}

	html := fmt.Sprintf(`
		<p class="text-sm font-medium text-gray-500">Ventas este mes</p>
		<p class="text-2xl font-bold text-gray-900 mt-1">$%s</p>
		<p class="text-xs text-gray-400 mt-1">%d transacciones</p>
	`, formatMoney(total), count)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// TopProductos returns top 5 selling products as an HTML fragment.
func (h *MetricsHandler) TopProductos(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT p.nombre, SUM(vi.cantidad) as unidades, 
		       SUM(vi.cantidad * vi.precio_unitario) as total_venta
		FROM venta_items vi
		JOIN productos p ON p.id = vi.producto_id
		JOIN ventas v ON v.id = vi.venta_id
		WHERE v.created_at >= datetime('now', '-30 days', 'localtime')
		GROUP BY vi.producto_id
		ORDER BY unidades DESC
		LIMIT 5
	`)
	if err != nil {
		h.renderError(w, "Error cargando top productos")
		return
	}
	defer rows.Close()

	html := `<p class="text-sm font-medium text-gray-500 mb-3">Top 5 productos (30 dias)</p>`
	html += `<div class="space-y-2">`
	found := false
	rank := 1
	for rows.Next() {
		found = true
		var nombre string
		var unidades float64
		var totalVenta float64
		if err := rows.Scan(&nombre, &unidades, &totalVenta); err != nil {
			continue
		}
		html += fmt.Sprintf(`
			<div class="flex items-center justify-between text-sm">
				<div class="flex items-center gap-2">
					<span class="w-5 h-5 bg-indigo-100 text-indigo-700 rounded text-xs flex items-center justify-center font-bold">%d</span>
					<span class="text-gray-700">%s</span>
				</div>
				<span class="text-gray-500">%.0f uds · $%s</span>
			</div>
		`, rank, nombre, unidades, formatMoney(totalVenta))
		rank++
	}
	if !found {
		html += `<p class="text-sm text-gray-400 text-center py-4">Sin ventas en los ultimos 30 dias</p>`
	}
	html += `</div>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// StockBajo returns low stock products.
func (h *MetricsHandler) StockBajo(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT nombre, stock_actual, stock_minimo
		FROM productos
		WHERE stock_actual <= stock_minimo AND activo = 1
		ORDER BY stock_actual ASC
		LIMIT 10
	`)
	if err != nil {
		h.renderError(w, "Error cargando stock bajo")
		return
	}
	defer rows.Close()

	html := `<p class="text-sm font-medium text-gray-500 mb-3">Stock bajo</p>`
	found := false
	for rows.Next() {
		found = true
		var nombre string
		var stockActual, stockMinimo float64
		if err := rows.Scan(&nombre, &stockActual, &stockMinimo); err != nil {
			continue
		}
		color := "text-yellow-600"
		icon := "⚠️"
		if stockActual == 0 {
			color = "text-red-600"
			icon = "🔴"
		}
		html += fmt.Sprintf(`
			<div class="flex items-center justify-between text-sm py-1">
				<span class="%s">%s %s</span>
				<span class="%s font-medium">%.0f / %.0f</span>
			</div>
		`, color, icon, nombre, color, stockActual, stockMinimo)
	}
	if !found {
		html += `<div class="flex items-center gap-2 text-sm text-green-600 py-4 justify-center"><span>✅</span> Todo en orden</div>`
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// ClientesFrecuentes returns top frequent customers.
func (h *MetricsHandler) ClientesFrecuentes(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT c.nombre, COUNT(v.id) as compras, COALESCE(SUM(v.total), 0) as total_gastado
		FROM clientes c
		JOIN ventas v ON v.cliente_id = c.id
		WHERE v.created_at >= datetime('now', '-30 days', 'localtime')
		GROUP BY c.id
		ORDER BY compras DESC
		LIMIT 5
	`)
	if err != nil {
		h.renderError(w, "Error cargando clientes frecuentes")
		return
	}
	defer rows.Close()

	html := `<p class="text-sm font-medium text-gray-500 mb-3">Clientes frecuentes (30 dias)</p>`
	html += `<div class="space-y-2">`
	found := false
	for rows.Next() {
		found = true
		var nombre string
		var compras int
		var totalGastado float64
		if err := rows.Scan(&nombre, &compras, &totalGastado); err != nil {
			continue
		}
		html += fmt.Sprintf(`
			<div class="flex items-center justify-between text-sm">
				<span class="text-gray-700">%s</span>
				<span class="text-gray-500">%d compras · $%s</span>
			</div>
		`, nombre, compras, formatMoney(totalGastado))
	}
	if !found {
		html += `<p class="text-sm text-gray-400 text-center py-4">Sin clientes registrados</p>`
	}
	html += `</div>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// Ingresos returns revenue for the last 14 days.
func (h *MetricsHandler) Ingresos(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT 'hoy' AS periodo, COUNT(*) AS ventas, COALESCE(SUM(total), 0) AS ingresos FROM ventas WHERE DATE(created_at) = DATE('now', 'localtime')
		UNION ALL
		SELECT 'semana', COUNT(*), COALESCE(SUM(total), 0) FROM ventas WHERE created_at >= datetime('now', '-7 days', 'localtime')
		UNION ALL
		SELECT 'mes', COUNT(*), COALESCE(SUM(total), 0) FROM ventas WHERE strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now', 'localtime')
	`)
	if err != nil {
		h.renderError(w, "Error cargando ingresos")
		return
	}
	defer rows.Close()

	html := `<p class="text-sm font-medium text-gray-500 mb-3">Resumen de ingresos</p>`
	html += `<div class="space-y-2">`
	for rows.Next() {
		var periodo string
		var ventas int
		var ingresos float64
		if err := rows.Scan(&periodo, &ventas, &ingresos); err != nil {
			continue
		}
		html += fmt.Sprintf(`
			<div class="flex items-center justify-between text-sm">
				<span class="text-gray-600 capitalize">%s</span>
				<span class="text-gray-900 font-medium">$%s <span class="text-gray-400 text-xs">(%d ventas)</span></span>
			</div>
		`, periodo, formatMoney(ingresos), ventas)
	}
	html += `</div>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func (h *MetricsHandler) renderError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<p class="text-sm text-red-500">%s</p>`, msg)
}

func formatMoney(amount float64) string {
	if amount >= 1000 {
		return fmt.Sprintf("%.0f", amount)
	}
	return fmt.Sprintf("%.2f", amount)
}
