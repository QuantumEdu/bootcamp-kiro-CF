package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/adapters"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/config"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/database"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http/handlers"

	usecases "github.com/QuantumEdu/bootcamp-kiro-CF/src/application/use_cases"
)

func main() {
	cfg := config.Load()

	// --- Database ---
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Printf("Warning: migrations error (may already exist): %v", err)
	}

	// Read-only connection for NL->SQL queries
	readDB, err := database.NewReadOnly(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open read-only database: %v", err)
	}
	defer readDB.Close()

	// --- Templates ---
	tmpl, err := loadTemplates("templates")
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)
	}

	// --- Services ---
	openRouter := adapters.NewOpenRouterClient(cfg.OpenRouterAPIKey, cfg.OpenRouterModel)

	schema := getSchema()
	chatService := usecases.NewChatService(openRouter, readDB, schema, cfg.QueryTimeoutSecs)

	// --- Handlers ---
	pageHandler := handlers.NewPageHandler(tmpl)
	metricsHandler := handlers.NewMetricsHandler(db)
	chatHandler := handlers.NewChatHandler(chatService, tmpl)

	// --- Router ---
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	// Static files
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Pages
	r.Get("/", pageHandler.Dashboard)
	r.Get("/productos", pageHandler.Products)
	r.Get("/productos/nuevo", pageHandler.ProductForm)
	r.Get("/ventas", pageHandler.Sales)
	r.Get("/metricas", pageHandler.Metrics)

	// API: Metrics (HTMX fragments)
	r.Get("/api/metrics/ventas-hoy", metricsHandler.VentasHoy)
	r.Get("/api/metrics/ventas-semana", metricsHandler.VentasSemana)
	r.Get("/api/metrics/ventas-mes", metricsHandler.VentasMes)
	r.Get("/api/metrics/top-productos", metricsHandler.TopProductos)
	r.Get("/api/metrics/stock-bajo", metricsHandler.StockBajo)
	r.Get("/api/metrics/clientes-frecuentes", metricsHandler.ClientesFrecuentes)
	r.Get("/api/metrics/ingresos", metricsHandler.Ingresos)

	// API: Chat (NL -> SQL)
	r.Post("/api/chat", chatHandler.HandleChat)

	// Start server
	addr := ":" + cfg.Port
	log.Printf("POS AI-First server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func loadTemplates(dir string) (*template.Template, error) {
	tmpl := template.New("")

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".html" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading template %s: %w", path, err)
		}

		// Use relative path as template name
		relPath, _ := filepath.Rel(dir, path)
		_, err = tmpl.New(relPath).Parse(string(content))
		if err != nil {
			return fmt.Errorf("parsing template %s: %w", relPath, err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func getSchema() string {
	return `CREATE TABLE IF NOT EXISTS usuarios (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    pin_hash TEXT NOT NULL,
    rol TEXT NOT NULL CHECK(rol IN ('admin', 'cajero')),
    activo INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

CREATE TABLE IF NOT EXISTS categorias (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL UNIQUE,
    descripcion TEXT,
    activo INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS productos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    sku TEXT UNIQUE,
    categoria_id INTEGER REFERENCES categorias(id),
    precio_venta REAL NOT NULL CHECK(precio_venta > 0),
    precio_compra REAL NOT NULL DEFAULT 0,
    stock_actual REAL NOT NULL DEFAULT 0,
    stock_minimo REAL NOT NULL DEFAULT 0,
    unidad TEXT NOT NULL DEFAULT 'unidad',
    activo INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

CREATE TABLE IF NOT EXISTS clientes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    telefono TEXT,
    direccion TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

CREATE TABLE IF NOT EXISTS ventas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    usuario_id INTEGER NOT NULL REFERENCES usuarios(id),
    cliente_id INTEGER REFERENCES clientes(id),
    total REAL NOT NULL CHECK(total >= 0),
    metodo_pago TEXT NOT NULL CHECK(metodo_pago IN ('efectivo', 'tarjeta', 'transferencia', 'mixto')),
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

CREATE TABLE IF NOT EXISTS venta_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    venta_id INTEGER NOT NULL REFERENCES ventas(id) ON DELETE CASCADE,
    producto_id INTEGER NOT NULL REFERENCES productos(id),
    cantidad REAL NOT NULL CHECK(cantidad > 0),
    precio_unitario REAL NOT NULL,
    subtotal REAL NOT NULL CHECK(subtotal >= 0)
);

CREATE TABLE IF NOT EXISTS inventario_movimientos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    producto_id INTEGER NOT NULL REFERENCES productos(id),
    tipo TEXT NOT NULL CHECK(tipo IN ('entrada', 'salida', 'ajuste')),
    cantidad REAL NOT NULL CHECK(cantidad > 0),
    stock_resultante REAL NOT NULL,
    referencia_tipo TEXT,
    referencia_id INTEGER,
    motivo TEXT,
    usuario_id INTEGER REFERENCES usuarios(id),
    created_at TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

CREATE TABLE IF NOT EXISTS configuracion (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    clave TEXT NOT NULL UNIQUE,
    valor TEXT NOT NULL
);`
}
