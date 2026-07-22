package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
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

// ProductosSinRotacion returns products that have never been sold.
func (h *MetricsHandler) ProductosSinRotacion(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT p.nombre, p.sku, p.stock_actual, p.precio_venta,
		       (p.stock_actual * p.precio_venta) AS valor_inmovilizado
		FROM productos p
		WHERE p.activo = 1 AND p.id NOT IN (
			SELECT DISTINCT vi.producto_id FROM venta_items vi
		)
		ORDER BY valor_inmovilizado DESC
		LIMIT 10
	`)
	if err != nil {
		h.renderError(w, "Error cargando productos sin rotacion")
		return
	}
	defer rows.Close()

	html := `<p class="text-sm font-medium text-gray-500 mb-3">Productos sin rotacion</p>`
	html += `<div class="space-y-2">`
	found := false
	for rows.Next() {
		found = true
		var nombre, sku string
		var stockActual, precioVenta, valorInm float64
		if err := rows.Scan(&nombre, &sku, &stockActual, &precioVenta, &valorInm); err != nil {
			continue
		}
		skuLabel := sku
		if skuLabel == "" {
			skuLabel = "sin-sku"
		}
		color := "text-gray-700"
		if valorInm > 1000 {
			color = "text-red-600"
		}
		html += fmt.Sprintf(`
			<div class="flex items-center justify-between text-sm py-1">
				<div class="flex items-center gap-2">
					<span class="text-gray-500 text-xs">%s</span>
					<span class="text-gray-700">%s</span>
				</div>
				<span class="%s font-medium">$%s inmovilizado</span>
			</div>
		`, skuLabel, nombre, color, formatMoney(valorInm))
	}
	if !found {
		html += `<div class="flex items-center gap-2 text-sm text-green-600 py-4 justify-center"><span>✅</span> Todos los productos tienen ventas</div>`
	}
	html += `</div>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// MargenCategoria returns profit margins by category.
func (h *MetricsHandler) MargenCategoria(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		SELECT COALESCE(c.nombre, 'Sin categoria') as categoria,
		       COUNT(DISTINCT p.id) as productos,
		       COALESCE(AVG(p.precio_venta - p.precio_compra), 0) as margen_promedio,
		       COALESCE(AVG(CASE WHEN p.precio_venta > 0 
		           THEN ((p.precio_venta - p.precio_compra) / p.precio_venta) * 100 
		           ELSE 0 END), 0) as margen_pct
		FROM productos p
		LEFT JOIN categorias c ON c.id = p.categoria_id
		WHERE p.activo = 1
		GROUP BY p.categoria_id
		ORDER BY margen_promedio DESC
	`)
	if err != nil {
		h.renderError(w, "Error cargando margen por categoria")
		return
	}
	defer rows.Close()

	html := `<p class="text-sm font-medium text-gray-500 mb-3">Margen por categoria</p>`
	html += `<div class="space-y-2">`
	found := false
	for rows.Next() {
		found = true
		var categoria string
		var productos int
		var margenPromedio, margenPct float64
		if err := rows.Scan(&categoria, &productos, &margenPromedio, &margenPct); err != nil {
			continue
		}
		barColor := "bg-green-500"
		if margenPct < 20 {
			barColor = "bg-red-500"
		} else if margenPct < 40 {
			barColor = "bg-yellow-500"
		}
		barWidth := margenPct
		if barWidth > 100 {
			barWidth = 100
		}
		html += fmt.Sprintf(`
			<div class="text-sm">
				<div class="flex items-center justify-between mb-1">
					<span class="text-gray-700">%s <span class="text-gray-400 text-xs">(%d prod)</span></span>
					<span class="font-medium text-gray-900">%.0f%%</span>
				</div>
				<div class="w-full bg-gray-200 rounded-full h-1.5">
					<div class="%s h-1.5 rounded-full" style="width: %.0f%%"></div>
				</div>
			</div>
		`, categoria, productos, margenPct, barColor, barWidth)
	}
	if !found {
		html += `<p class="text-sm text-gray-400 text-center py-4">Sin productos registrados</p>`
	}
	html += `</div>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// VentasSparkline returns an SVG sparkline for the last 7 days of sales.
func (h *MetricsHandler) VentasSparkline(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`
		WITH RECURSIVE dates(dia) AS (
			SELECT date('now', '-6 days', 'localtime')
			UNION ALL
			SELECT date(dia, '+1 day') FROM dates WHERE dia < date('now', 'localtime')
		)
		SELECT d.dia, COALESCE(SUM(v.total), 0) as total
		FROM dates d
		LEFT JOIN ventas v ON date(v.created_at) = d.dia
		GROUP BY d.dia
		ORDER BY d.dia ASC
	`)
	if err != nil {
		h.renderError(w, "Error cargando sparkline")
		return
	}
	defer rows.Close()

	type dayData struct {
		dia   string
		total float64
	}
	var days []dayData
	var maxTotal float64
	for rows.Next() {
		var d dayData
		if err := rows.Scan(&d.dia, &d.total); err != nil {
			continue
		}
		days = append(days, d)
		if d.total > maxTotal {
			maxTotal = d.total
		}
	}

	// Build SVG sparkline
	svgWidth := 280
	svgHeight := 60
	barWidth := 32
	gap := 8

	html := `<p class="text-sm font-medium text-gray-500 mb-2">Ventas ultimos 7 dias</p>`
	html += fmt.Sprintf(`<svg width="%d" height="%d" class="w-full">`, svgWidth, svgHeight)

	for i, d := range days {
		barHeight := 4 // minimum
		if maxTotal > 0 {
			barHeight = int((d.total / maxTotal) * float64(svgHeight-16))
			if barHeight < 4 {
				barHeight = 4
			}
		}
		x := i * (barWidth + gap)
		y := svgHeight - barHeight - 12
		html += fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" rx="3" class="fill-indigo-500" opacity="0.8"/>`, x, y, barWidth, barHeight)
		// Day label
		dayLabel := d.dia[5:] // MM-DD
		html += fmt.Sprintf(`<text x="%d" y="%d" class="text-[8px] fill-gray-400" text-anchor="middle">%s</text>`, x+barWidth/2, svgHeight-2, dayLabel)
	}
	html += `</svg>`

	if maxTotal > 0 {
		html += fmt.Sprintf(`<p class="text-xs text-gray-400 mt-1">Pico: $%s</p>`, formatMoney(maxTotal))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// SQLLibre executes a user-provided SQL query (read-only, validated).
func (h *MetricsHandler) SQLLibre(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.FormValue("sql_query")
	if query == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<p class="text-sm text-gray-400">Escribe una consulta SQL para ejecutar</p>`)
		return
	}

	// Import validator
	// Validation happens in the handler directly using the same logic
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Basic validation (same as usecases.ValidateSQL but inline here for independence)
	if err := validateSQLQuery(query); err != nil {
		fmt.Fprintf(w, `<div class="p-3 bg-red-50 border border-red-200 rounded-lg"><p class="text-sm text-red-600">⚠️ %s</p></div>`, err.Error())
		return
	}

	// Execute with read-only DB (passed in handler)
	rows, err := h.db.Query(query)
	if err != nil {
		fmt.Fprintf(w, `<div class="p-3 bg-red-50 border border-red-200 rounded-lg"><p class="text-sm text-red-600">Error SQL: %s</p></div>`, err.Error())
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Fprintf(w, `<div class="p-3 bg-red-50 border border-red-200 rounded-lg"><p class="text-sm text-red-600">Error: %s</p></div>`, err.Error())
		return
	}

	html := `<div class="overflow-x-auto"><table class="w-full text-xs border border-gray-200 rounded">`
	html += `<thead class="bg-gray-100"><tr>`
	for _, col := range columns {
		html += fmt.Sprintf(`<th class="px-2 py-1 text-left font-medium text-gray-600 border-b">%s</th>`, col)
	}
	html += `</tr></thead><tbody>`

	rowCount := 0
	maxRows := 100
	for rows.Next() {
		if rowCount >= maxRows {
			break
		}
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			continue
		}
		html += `<tr class="border-b border-gray-100 hover:bg-gray-50">`
		for _, v := range values {
			val := "NULL"
			if v != nil {
				val = fmt.Sprintf("%v", v)
			}
			html += fmt.Sprintf(`<td class="px-2 py-1 text-gray-700">%s</td>`, val)
		}
		html += `</tr>`
		rowCount++
	}
	html += `</tbody></table></div>`
	html += fmt.Sprintf(`<p class="text-xs text-gray-400 mt-2">%d filas devueltas (max %d)</p>`, rowCount, maxRows)

	fmt.Fprint(w, html)
}

func validateSQLQuery(sql string) error {
	upper := strings.ToUpper(strings.TrimSpace(sql))
	if !strings.HasPrefix(upper, "SELECT") && !strings.HasPrefix(upper, "WITH") {
		return fmt.Errorf("solo se permiten consultas SELECT")
	}
	dangerous := []string{"INSERT", "UPDATE", "DELETE", "DROP", "ALTER", "CREATE", "TRUNCATE", "PRAGMA", "ATTACH", "DETACH"}
	for _, kw := range dangerous {
		if strings.Contains(upper, kw) {
			return fmt.Errorf("palabra clave no permitida: %s", kw)
		}
	}
	if strings.Contains(sql, "--") || strings.Contains(sql, "/*") {
		return fmt.Errorf("no se permiten comentarios SQL")
	}
	return nil
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
