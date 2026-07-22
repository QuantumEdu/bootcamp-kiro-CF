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

// VentasHoy returns today's sales summary.
func (h *MetricsHandler) VentasHoy(w http.ResponseWriter, r *http.Request) {
	var count int
	var total float64
	h.db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(total), 0) FROM ventas WHERE DATE(created_at) = DATE('now', 'localtime')`).Scan(&count, &total)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<p class="text-sm font-medium text-gray-500">Ventas hoy</p><p class="text-2xl font-bold text-gray-900 mt-1">$%s</p><p class="text-xs text-gray-400 mt-1">%d transacciones</p>`, fmtMoney(total), count)
}

// VentasSemana returns this week's sales.
func (h *MetricsHandler) VentasSemana(w http.ResponseWriter, r *http.Request) {
	var count int
	var total float64
	h.db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(total), 0) FROM ventas WHERE created_at >= datetime('now', '-7 days', 'localtime')`).Scan(&count, &total)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<p class="text-sm font-medium text-gray-500">Ventas semana</p><p class="text-2xl font-bold text-gray-900 mt-1">$%s</p><p class="text-xs text-gray-400 mt-1">%d transacciones (7 dias)</p>`, fmtMoney(total), count)
}

// VentasMes returns this month's sales.
func (h *MetricsHandler) VentasMes(w http.ResponseWriter, r *http.Request) {
	var count int
	var total float64
	h.db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(total), 0) FROM ventas WHERE strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now', 'localtime')`).Scan(&count, &total)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<p class="text-sm font-medium text-gray-500">Ventas mes</p><p class="text-2xl font-bold text-gray-900 mt-1">$%s</p><p class="text-xs text-gray-400 mt-1">%d transacciones</p>`, fmtMoney(total), count)
}

// TopProductos returns top 5 selling products.
func (h *MetricsHandler) TopProductos(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT p.nombre, SUM(vi.cantidad) as unidades
		FROM venta_items vi JOIN productos p ON p.id = vi.producto_id
		JOIN ventas v ON v.id = vi.venta_id
		WHERE v.created_at >= datetime('now', '-30 days', 'localtime')
		GROUP BY vi.producto_id ORDER BY unidades DESC LIMIT 5`)
	if err != nil {
		renderErr(w, "Error cargando top productos")
		return
	}
	defer rows.Close()

	html := `<p class="text-sm font-medium text-gray-500 mb-3">Top 5 productos (30 dias)</p><div class="space-y-2">`
	rank := 1
	found := false
	for rows.Next() {
		found = true
		var nombre string
		var unidades float64
		rows.Scan(&nombre, &unidades)
		html += fmt.Sprintf(`<div class="flex items-center justify-between text-sm"><div class="flex items-center gap-2"><span class="w-5 h-5 bg-indigo-100 text-indigo-700 rounded text-xs flex items-center justify-center font-bold">%d</span><span class="text-gray-700">%s</span></div><span class="text-gray-500">%.0f uds</span></div>`, rank, nombre, unidades)
		rank++
	}
	if !found {
		html += `<p class="text-sm text-gray-400 text-center py-4">Sin ventas recientes</p>`
	}
	html += `</div>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// StockBajo returns low stock products.
func (h *MetricsHandler) StockBajo(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`SELECT nombre, stock_actual, stock_minimo FROM productos WHERE stock_actual <= stock_minimo AND activo = 1 ORDER BY stock_actual ASC LIMIT 10`)
	if err != nil {
		renderErr(w, "Error cargando stock bajo")
		return
	}
	defer rows.Close()

	html := `<p class="text-sm font-medium text-gray-500 mb-3">Stock bajo</p>`
	found := false
	for rows.Next() {
		found = true
		var nombre string
		var stockActual, stockMinimo float64
		rows.Scan(&nombre, &stockActual, &stockMinimo)
		icon := "⚠️"
		color := "text-yellow-600"
		if stockActual == 0 {
			icon = "🔴"
			color = "text-red-600"
		}
		html += fmt.Sprintf(`<div class="flex items-center justify-between text-sm py-1"><span class="%s">%s %s</span><span class="%s font-medium">%.0f / %.0f</span></div>`, color, icon, nombre, color, stockActual, stockMinimo)
	}
	if !found {
		html += `<div class="text-sm text-green-600 py-4 text-center">✅ Todo en orden</div>`
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// ClientesFrecuentes returns top customers.
func (h *MetricsHandler) ClientesFrecuentes(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT c.nombre, COUNT(v.id) as compras, COALESCE(SUM(v.total), 0) as total_gastado
		FROM clientes c JOIN ventas v ON v.cliente_id = c.id
		WHERE v.created_at >= datetime('now', '-30 days', 'localtime')
		GROUP BY c.id ORDER BY compras DESC LIMIT 5`)
	if err != nil {
		renderErr(w, "Error cargando clientes")
		return
	}
	defer rows.Close()

	html := `<p class="text-sm font-medium text-gray-500 mb-3">Clientes frecuentes (30d)</p><div class="space-y-2">`
	found := false
	for rows.Next() {
		found = true
		var nombre string
		var compras int
		var total float64
		rows.Scan(&nombre, &compras, &total)
		html += fmt.Sprintf(`<div class="flex items-center justify-between text-sm"><span class="text-gray-700">%s</span><span class="text-gray-500">%d compras · $%s</span></div>`, nombre, compras, fmtMoney(total))
	}
	if !found {
		html += `<p class="text-sm text-gray-400 text-center py-4">Sin datos</p>`
	}
	html += `</div>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// MargenCategoria returns profit margins by category.
func (h *MetricsHandler) MargenCategoria(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT COALESCE(c.nombre, 'Sin categoria') as cat, COUNT(p.id) as prods,
		       COALESCE(AVG(CASE WHEN p.precio_venta > 0 THEN ((p.precio_venta - p.precio_compra) / p.precio_venta) * 100 ELSE 0 END), 0) as pct
		FROM productos p LEFT JOIN categorias c ON c.id = p.categoria_id
		WHERE p.activo = 1 GROUP BY p.categoria_id ORDER BY pct DESC`)
	if err != nil {
		renderErr(w, "Error cargando margenes")
		return
	}
	defer rows.Close()

	html := `<p class="text-sm font-medium text-gray-500 mb-3">Margen por categoria</p><div class="space-y-2">`
	found := false
	for rows.Next() {
		found = true
		var cat string
		var prods int
		var pct float64
		rows.Scan(&cat, &prods, &pct)
		barColor := "bg-green-500"
		if pct < 20 {
			barColor = "bg-red-500"
		} else if pct < 40 {
			barColor = "bg-yellow-500"
		}
		bw := pct
		if bw > 100 {
			bw = 100
		}
		html += fmt.Sprintf(`<div class="text-sm"><div class="flex justify-between mb-1"><span class="text-gray-700">%s <span class="text-gray-400 text-xs">(%d)</span></span><span class="font-medium">%.0f%%</span></div><div class="w-full bg-gray-200 rounded-full h-1.5"><div class="%s h-1.5 rounded-full" style="width:%.0f%%"></div></div></div>`, cat, prods, pct, barColor, bw)
	}
	if !found {
		html += `<p class="text-sm text-gray-400 text-center py-4">Sin datos</p>`
	}
	html += `</div>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// ProductosHTMX returns products table rows.
func (h *MetricsHandler) ProductosHTMX(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`SELECT p.nombre, COALESCE(p.sku,''), p.precio_venta, p.stock_actual FROM productos p WHERE p.activo = 1 ORDER BY p.nombre LIMIT 50`)
	if err != nil {
		renderErr(w, "Error")
		return
	}
	defer rows.Close()

	html := ""
	found := false
	for rows.Next() {
		found = true
		var nombre, sku string
		var precio, stock float64
		rows.Scan(&nombre, &sku, &precio, &stock)
		html += fmt.Sprintf(`<tr class="hover:bg-gray-50"><td class="px-4 py-3 font-medium text-gray-900">%s</td><td class="px-4 py-3 text-gray-500">%s</td><td class="px-4 py-3 text-right text-gray-900">$%.2f</td><td class="px-4 py-3 text-right">%.0f</td></tr>`, nombre, sku, precio, stock)
	}
	if !found {
		html = `<tr><td colspan="4" class="px-4 py-8 text-center text-gray-400">No hay productos</td></tr>`
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// VentasRecientes returns recent sales fragment.
func (h *MetricsHandler) VentasRecientes(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`SELECT v.id, v.total, v.metodo_pago, v.created_at FROM ventas v ORDER BY v.created_at DESC LIMIT 10`)
	if err != nil {
		renderErr(w, "Error")
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
		} else if metodo == "transferencia" {
			badge = "bg-purple-100 text-purple-700"
		}
		display := fecha
		if len(fecha) > 16 {
			display = fecha[:16]
		}
		html += fmt.Sprintf(`<div class="flex items-center justify-between p-2 bg-gray-50 rounded-lg text-sm"><div><span class="font-medium text-gray-800">#%d</span><span class="text-gray-500 ml-2">%s</span></div><div class="flex items-center gap-2"><span class="px-2 py-0.5 rounded text-xs %s">%s</span><span class="font-bold text-gray-900">$%.2f</span></div></div>`, id, display, badge, metodo, total)
	}
	if !found {
		html += `<p class="text-sm text-gray-400 text-center py-4">No hay ventas</p>`
	}
	html += `</div>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// ProductosBuscar returns product search results for the sales cart.
func (h *MetricsHandler) ProductosBuscar(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		return
	}
	s := "%" + q + "%"
	rows, err := h.db.Query(`SELECT id, nombre, precio_venta, stock_actual FROM productos WHERE activo = 1 AND (nombre LIKE ? OR sku LIKE ?) ORDER BY nombre LIMIT 10`, s, s)
	if err != nil {
		return
	}
	defer rows.Close()

	html := ""
	for rows.Next() {
		var id int
		var nombre string
		var precio, stock float64
		rows.Scan(&id, &nombre, &precio, &stock)
		html += fmt.Sprintf(`<button type="button" onclick="addToCart(%d, '%s', %.2f)" class="w-full flex items-center justify-between p-2 hover:bg-indigo-50 rounded-lg text-left"><div><p class="text-sm font-medium text-gray-800">%s</p><p class="text-xs text-gray-500">Stock: %.0f</p></div><span class="text-sm font-bold text-indigo-600">$%.2f</span></button>`, id, escJS(nombre), precio, nombre, stock, precio)
	}
	if html == "" {
		html = `<p class="text-sm text-gray-400 py-2">Sin resultados</p>`
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func renderErr(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<p class="text-sm text-red-500">%s</p>`, msg)
}

func fmtMoney(amount float64) string {
	if amount >= 1000 {
		return fmt.Sprintf("%.0f", amount)
	}
	return fmt.Sprintf("%.2f", amount)
}

func escJS(s string) string {
	out := ""
	for _, c := range s {
		if c == '\'' {
			out += "\\'"
		} else {
			out += string(c)
		}
	}
	return out
}
