package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/nlsql"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/adapters"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/config"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/database"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http/handlers"
	mw "github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	// Database
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Templates
	tmpl, err := loadTemplates("templates")
	if err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}

	// Services
	openRouter := adapters.NewOpenRouterClient(cfg.OpenRouterAPIKey, cfg.OpenRouterModel)
	nlsqlService := nlsql.NewService(openRouter, db.RO, getSchema(), cfg.QueryTimeoutSeconds)

	// Handlers
	pageHandler := handlers.NewPageHandler(tmpl)
	chatHandler := handlers.NewChatHandler(nlsqlService, tmpl)
	metricsHandler := handlers.NewMetricsHandler(db.RW)
	saleHandler := handlers.NewSaleHandler(db.RW)
	authMW := mw.NewAuthMiddleware(cfg.SessionSecret, cfg.PINMaxAttempts, cfg.PINLockoutMinutes)
	authHandler := handlers.NewAuthHandler(authMW, tmpl)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	// Static files
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","port":%q}`, cfg.Port)
	})

	// Public routes
	r.Get("/login", authHandler.LoginPage)
	r.Post("/login", authHandler.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMW.RequireAuth)

		// Pages
		r.Get("/", pageHandler.Dashboard)
		r.Get("/productos", pageHandler.Products)
		r.Get("/ventas", pageHandler.Sales)
		r.Get("/metricas", pageHandler.Metrics)
		r.Get("/logout", authHandler.Logout)

		// API: Sales
		r.Post("/api/ventas", saleHandler.Create)
		r.Get("/api/ventas/recientes", saleHandler.Recent)

		// API: Metrics (HTMX fragments)
		r.Get("/api/metrics/ventas-hoy", metricsHandler.VentasHoy)
		r.Get("/api/metrics/ventas-semana", metricsHandler.VentasSemana)
		r.Get("/api/metrics/ventas-mes", metricsHandler.VentasMes)
		r.Get("/api/metrics/top-productos", metricsHandler.TopProductos)
		r.Get("/api/metrics/stock-bajo", metricsHandler.StockBajo)
		r.Get("/api/metrics/clientes-frecuentes", metricsHandler.ClientesFrecuentes)
		r.Get("/api/metrics/margen-categoria", metricsHandler.MargenCategoria)

		// API: Products
		r.Get("/api/productos", metricsHandler.ProductosHTMX)
		r.Get("/api/productos/buscar", metricsHandler.ProductosBuscar)

		// API: Chat NL→SQL
		r.Post("/api/chat", chatHandler.HandleChat)
	})

	// Start
	addr := ":" + cfg.Port
	log.Printf("POS AI-First starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func loadTemplates(dir string) (*template.Template, error) {
	tmpl := template.New("")
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		relPath, _ := filepath.Rel(dir, path)
		_, err = tmpl.New(relPath).Parse(string(content))
		if err != nil {
			return fmt.Errorf("parsing %s: %w", relPath, err)
		}
		return nil
	})
	return tmpl, err
}

func getSchema() string {
	return `Tables:
- usuarios (id, nombre, pin_hash, rol, activo, created_at)
- categorias (id, nombre, descripcion, activo)
- productos (id, nombre, sku, categoria_id, precio_venta, precio_compra, stock_actual, stock_minimo, unidad, activo, created_at, updated_at)
- clientes (id, nombre, telefono, direccion, created_at)
- ventas (id, usuario_id, cliente_id, total, metodo_pago, created_at)
- venta_items (id, venta_id, producto_id, cantidad, precio_unitario, subtotal)
- inventario_movimientos (id, producto_id, tipo, cantidad, stock_resultante, referencia_tipo, referencia_id, motivo, usuario_id, created_at)
- configuracion (id, clave, valor)`
}
