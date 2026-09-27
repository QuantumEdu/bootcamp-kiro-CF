package i18n

// translationsEN contains all English UI strings.
var translationsEN = map[string]string{
	// App branding
	"app.subtitle": "AI-Powered Point of Sale",

	// Navigation
	"nav.dashboard": "Dashboard",
	"nav.products":  "Products",
	"nav.sales":     "Sales",
	"nav.clients":   "Clients",
	"nav.metrics":   "Metrics",
	"nav.config":    "Settings",
	"nav.logout":    "Sign Out",

	// Login page
	"login.title":     "Sign In — POS AI-First",
	"login.subtitle":  "AI-Powered Point of Sale",
	"login.pin_label": "Access PIN",
	"login.button":    "Sign In",

	// Chat panel
	"chat.title":       "AI Assistant",
	"chat.placeholder": "e.g. How much did I sell today?",
	"chat.welcome":     "Hi! Ask me anything about sales, products, inventory, or any business metric.",
	"chat.button":      "AI Chat",

	// Page titles
	"page.dashboard":     "Dashboard",
	"page.products":      "Products",
	"page.products.new":  "New Product",
	"page.products.edit": "Edit Product",
	"page.sales":         "Sales",
	"page.sales.new":     "New Sale",
	"page.clients":       "Clients",
	"page.clients.new":   "New Client",
	"page.metrics":       "Metrics",
	"page.config":        "Settings",
	"page.login":         "Sign In",

	// Products list
	"products.new":                "New Product",
	"products.col.product":        "Product",
	"products.col.sku":            "SKU",
	"products.col.price":          "Price",
	"products.col.stock":          "Stock",
	"products.action.edit":        "Edit",
	"products.action.deactivate":  "Deactivate",
	"products.confirm.deactivate": "Deactivate this product?",
	"products.loading":            "Loading...",
	"products.empty":              "No products found",

	// Product form
	"product.form.name":           "Name",
	"product.form.sku":            "SKU",
	"product.form.sale_price":     "Sale price",
	"product.form.purchase_price": "Purchase price",
	"product.form.stock":          "Current stock",
	"product.form.min_stock":      "Minimum stock",
	"product.form.unit":           "Unit",
	"product.form.category":       "Category ID",
	"product.form.save":           "Save",
	"product.form.cancel":         "Cancel",
	"product.unit.unit":           "Unit",
	"product.unit.kg":             "Kg",
	"product.unit.liter":          "Liter",
	"product.unit.package":        "Package",

	// Sales
	"sales.new":              "New Sale",
	"sales.search":           "Search product by name or SKU...",
	"sales.recent":           "Recent sales",
	"sales.available":        "Available products",
	"sales.cart":             "Cart",
	"sales.cart.empty":       "Add products to cart",
	"sales.total":            "Total:",
	"sales.checkout":         "Charge",
	"sales.processing":       "Processing...",
	"sales.payment.cash":     "Cash",
	"sales.payment.card":     "Card",
	"sales.payment.transfer": "Transfer",
	"sales.error.connection": "Connection error",

	// Clients
	"clients.title":                    "Clients",
	"clients.new":                      "New Client",
	"clients.col.name":                 "Name",
	"clients.col.phone":                "Phone",
	"clients.col.address":              "Address",
	"clients.empty":                    "No clients registered",
	"clients.form.name":                "Name",
	"clients.form.name.placeholder":    "Client name",
	"clients.form.phone":               "Phone",
	"clients.form.phone.placeholder":   "Phone number",
	"clients.form.address":             "Address",
	"clients.form.address.placeholder": "Client address",
	"clients.form.save":                "Save",
	"clients.form.cancel":              "Cancel",

	// Admin config
	"config.title":       "Settings — API Key",
	"config.current_key": "Current API Key",
	"config.update":      "Update",
	"config.set":         "Configure",
	"config.key_label":   "New API Key",
	"config.save":        "Save API Key",
	"config.no_key":      "No API key configured. Enter a key to enable OpenRouter integration.",

	// Metrics dashboard
	"metrics.recent_sales":     "Recent sales",
	"metrics.sales_today":      "Sales today",
	"metrics.sales_week":       "Sales this week",
	"metrics.sales_month":      "Sales this month",
	"metrics.transactions":     "transactions",
	"metrics.transactions_7d":  "transactions (7 days)",
	"metrics.top_products":     "Top 5 products (30 days)",
	"metrics.units":            "units",
	"metrics.no_recent_sales":  "No recent sales",
	"metrics.low_stock":        "Low stock",
	"metrics.stock_ok":         "All good",
	"metrics.frequent_clients": "Frequent clients (30d)",
	"metrics.purchases":        "purchases",
	"metrics.no_data":          "No data",
	"metrics.margin_category":  "Margin by category",
	"metrics.error_loading":    "Error loading data",
	"metrics.no_sales":         "No sales yet",
	"metrics.stock_label":      "Stock",
	"metrics.no_results":       "No results",

	// Language switch
	"lang.switch.en": "English",
	"lang.switch.es": "Español",
}
