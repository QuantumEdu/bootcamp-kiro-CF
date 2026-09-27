package i18n

// translationsES contains all Spanish UI strings.
var translationsES = map[string]string{
	// App branding
	"app.subtitle": "Punto de Venta Inteligente",

	// Navigation
	"nav.dashboard": "Dashboard",
	"nav.products":  "Productos",
	"nav.sales":     "Ventas",
	"nav.clients":   "Clientes",
	"nav.metrics":   "Métricas",
	"nav.config":    "Configuración",
	"nav.logout":    "Cerrar Sesión",

	// Login page
	"login.title":     "Iniciar Sesión — POS AI-First",
	"login.subtitle":  "Punto de Venta Inteligente",
	"login.pin_label": "PIN de acceso",
	"login.button":    "Entrar",

	// Chat panel
	"chat.title":       "Asistente IA",
	"chat.placeholder": "Ej: ¿Cuánto vendí hoy?",
	"chat.welcome":     "¡Hola! Pregúntame sobre ventas, productos, inventario o cualquier métrica de tu negocio.",
	"chat.button":      "Chat IA",

	// Page titles
	"page.dashboard":     "Dashboard",
	"page.products":      "Productos",
	"page.products.new":  "Nuevo Producto",
	"page.products.edit": "Editar Producto",
	"page.sales":         "Ventas",
	"page.sales.new":     "Nueva Venta",
	"page.clients":       "Clientes",
	"page.clients.new":   "Nuevo Cliente",
	"page.metrics":       "Métricas",
	"page.config":        "Configuración",
	"page.login":         "Iniciar Sesión",

	// Products list
	"products.new":                "Nuevo Producto",
	"products.col.product":        "Producto",
	"products.col.sku":            "SKU",
	"products.col.price":          "Precio",
	"products.col.stock":          "Stock",
	"products.action.edit":        "Editar",
	"products.action.deactivate":  "Desactivar",
	"products.confirm.deactivate": "¿Desactivar este producto?",
	"products.loading":            "Cargando...",
	"products.empty":              "No hay productos",

	// Product form
	"product.form.name":           "Nombre",
	"product.form.sku":            "SKU",
	"product.form.sale_price":     "Precio venta",
	"product.form.purchase_price": "Precio compra",
	"product.form.stock":          "Stock actual",
	"product.form.min_stock":      "Stock mínimo",
	"product.form.unit":           "Unidad",
	"product.form.category":       "Categoría ID",
	"product.form.save":           "Guardar",
	"product.form.cancel":         "Cancelar",
	"product.unit.unit":           "Unidad",
	"product.unit.kg":             "Kg",
	"product.unit.liter":          "Litro",
	"product.unit.package":        "Paquete",

	// Sales
	"sales.new":              "Nueva Venta",
	"sales.search":           "Buscar producto por nombre o SKU...",
	"sales.recent":           "Ventas recientes",
	"sales.available":        "Productos disponibles",
	"sales.cart":             "Carrito",
	"sales.cart.empty":       "Agrega productos al carrito",
	"sales.total":            "Total:",
	"sales.checkout":         "Cobrar",
	"sales.processing":       "Procesando...",
	"sales.payment.cash":     "Efectivo",
	"sales.payment.card":     "Tarjeta",
	"sales.payment.transfer": "Transferencia",
	"sales.error.connection": "Error de conexión",

	// Clients
	"clients.title":                    "Clientes",
	"clients.new":                      "Nuevo Cliente",
	"clients.col.name":                 "Nombre",
	"clients.col.phone":                "Teléfono",
	"clients.col.address":              "Dirección",
	"clients.empty":                    "No hay clientes registrados",
	"clients.form.name":                "Nombre",
	"clients.form.name.placeholder":    "Nombre del cliente",
	"clients.form.phone":               "Teléfono",
	"clients.form.phone.placeholder":   "Número de teléfono",
	"clients.form.address":             "Dirección",
	"clients.form.address.placeholder": "Dirección del cliente",
	"clients.form.save":                "Guardar",
	"clients.form.cancel":              "Cancelar",

	// Admin config
	"config.title":       "Configuración — API Key",
	"config.current_key": "API Key actual",
	"config.update":      "Actualizar",
	"config.set":         "Configurar",
	"config.key_label":   "Nueva API Key",
	"config.save":        "Guardar API Key",
	"config.no_key":      "No hay API key configurada. Ingresa una clave para habilitar la integración con OpenRouter.",

	// Metrics dashboard
	"metrics.recent_sales":     "Ventas recientes",
	"metrics.sales_today":      "Ventas hoy",
	"metrics.sales_week":       "Ventas semana",
	"metrics.sales_month":      "Ventas mes",
	"metrics.transactions":     "transacciones",
	"metrics.transactions_7d":  "transacciones (7 días)",
	"metrics.top_products":     "Top 5 productos (30 días)",
	"metrics.units":            "uds",
	"metrics.no_recent_sales":  "Sin ventas recientes",
	"metrics.low_stock":        "Stock bajo",
	"metrics.stock_ok":         "Todo en orden",
	"metrics.frequent_clients": "Clientes frecuentes (30d)",
	"metrics.purchases":        "compras",
	"metrics.no_data":          "Sin datos",
	"metrics.margin_category":  "Margen por categoría",
	"metrics.error_loading":    "Error cargando datos",
	"metrics.no_sales":         "No hay ventas",
	"metrics.stock_label":      "Stock",
	"metrics.no_results":       "Sin resultados",

	// Language switch
	"lang.switch.en": "English",
	"lang.switch.es": "Español",
}
