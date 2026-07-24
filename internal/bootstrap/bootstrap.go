package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/nlsql"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/services"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/application/use_cases"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/adapters"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/database"
	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http/handlers"
	infrahttp "github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http"
	mw "github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/http/middleware"
)

// TemplatesDir is the path to the templates directory.
// Can be overridden for Lambda (where templates are bundled at a fixed path).
var TemplatesDir = "templates"

// StaticDir is the path to the static assets directory for local mode.
var StaticDir = "static"

// BuildRouter constructs the chi router with all dependencies wired.
// Returns the router, a cleanup function, and any initialization error.
// It switches on cfg.AppEnv to select the appropriate adapters:
//   - "lambda" → PostgreSQL + Bedrock + pgx session store
//   - default  → SQLite + OpenRouter + SQLite session store
func BuildRouter(cfg Config) (http.Handler, func(), error) {
	if err := ValidateConfig(&cfg); err != nil {
		return nil, nil, fmt.Errorf("config validation failed: %w", err)
	}

	var (
		readDB         *sql.DB
		writeDB        *sql.DB
		sessionManager *scs.SessionManager
		cleanup        func()
		pool           *pgxpool.Pool // nil in local/SQLite mode
	)

	// Templates
	tmpl, err := loadTemplates(TemplatesDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load templates: %w", err)
	}

	// Services
	cryptoService := services.NewCryptoService(cfg.SessionSecret)

	// Schema for NL-SQL service
	schema := getSchema()

	// Declare handlers outside switch so they're accessible for routing
	var (
		chatHandler        *handlers.ChatHandler
		metricsHandler     *handlers.MetricsHandler
		authHandler        *handlers.AuthHandler
		productHandler     *handlers.ProductHandler
		saleHandler        *handlers.SaleHandler
		adminConfigHandler *handlers.AdminConfigHandler
		clientHandler      *handlers.ClientHandler
		pageHandler        *handlers.PageHandler
	)

	switch cfg.AppEnv {
	case "lambda":
		// PostgreSQL path
		var poolErr error
		pool, poolErr = newPgxPool(cfg.DatabaseURL)
		if poolErr != nil {
			return nil, nil, fmt.Errorf("creating pgx pool: %w", poolErr)
		}

		// Verify connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			return nil, nil, fmt.Errorf("pinging PostgreSQL: %w", err)
		}

		// stdlib *sql.DB for NL-SQL read-only queries and metrics handler
		readDB = stdlib.OpenDBFromPool(pool)
		writeDB = readDB // same pool in postgres mode

		// PostgreSQL repositories
		userRepo := adapters.NewPostgresUserRepository(pool)
		productRepo := adapters.NewPostgresProductRepository(pool)
		saleRepo := adapters.NewPostgresSaleRepository(pool)
		inventoryRepo := adapters.NewPostgresInventoryRepository(pool)
		clientRepo := adapters.NewPostgresClientRepository(pool)
		configRepo := adapters.NewPostgresConfigRepository(pool)

		// Bedrock AI service
		bedrockClient := bedrockruntime.New(bedrockruntime.Options{
			Region: cfg.BedrockRegion,
		})
		bedrockCfg := adapters.BedrockConfig{
			ModelID:     cfg.BedrockModelID,
			Region:      cfg.BedrockRegion,
			MaxTokens:   cfg.MaxTokens,
			Temperature: cfg.Temperature,
		}
		_ = adapters.NewBedrockQueryService(bedrockClient, bedrockCfg, schema)

		// NL-SQL service for Lambda mode — uses Bedrock via the nlsql service
		// The nlsql.Service currently takes *adapters.OpenRouterClient directly.
		// For Lambda, we create a nil OpenRouter client and use a Bedrock-backed service instead.
		// TODO: Refactor nlsql.Service to accept ports.AIQueryService interface
		nlsqlService := nlsql.NewService(nil, readDB, schema, cfg.QueryTimeoutSeconds)
		nlsqlService.SetLogger(nlsql.NewQueryLogger(writeDB))

		// pgx session store
		sessionManager = newPgxSessionManager(pool)

		// Use cases
		authUC := use_cases.NewAuthenticateUser(userRepo)
		createProductUC := use_cases.NewCreateProduct(productRepo)
		updateProductUC := use_cases.NewUpdateProduct(productRepo)
		listProductsUC := use_cases.NewListProducts(productRepo)
		deactivateProductUC := use_cases.NewDeactivateProduct(productRepo)
		registerSaleUC := use_cases.NewRegisterSale(saleRepo, productRepo, inventoryRepo)
		createClientUC := use_cases.NewCreateClient(clientRepo)
		listClientsUC := use_cases.NewListClients(clientRepo)

		// Handlers
		pageHandler = handlers.NewPageHandler(tmpl)
		chatHandler = handlers.NewChatHandler(nlsqlService, tmpl)
		metricsHandler = handlers.NewMetricsHandler(writeDB)
		authHandler = handlers.NewAuthHandler(authUC, tmpl, sessionManager)
		productHandler = handlers.NewProductHandler(createProductUC, updateProductUC, listProductsUC, deactivateProductUC, productRepo, tmpl)
		saleHandler = handlers.NewSaleHandler(registerSaleUC, tmpl, sessionManager)
		adminConfigHandler = handlers.NewAdminConfigHandler(configRepo, cryptoService, tmpl)
		clientHandler = handlers.NewClientHandler(createClientUC, listClientsUC, tmpl)

		cleanup = func() {
			if readDB != nil {
				readDB.Close()
			}
			pool.Close()
		}

	default: // "local" or unset
		// SQLite path (existing logic)
		db, err := database.New(cfg.DatabasePath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to initialize database: %w", err)
		}

		readDB = db.RO
		writeDB = db.RW

		// OpenRouter adapter
		openRouter := adapters.NewOpenRouterClient(cfg.OpenRouterAPIKey, cfg.OpenRouterModel)
		nlsqlService := nlsql.NewService(openRouter, db.RO, schema, cfg.QueryTimeoutSeconds)
		nlsqlService.SetLogger(nlsql.NewQueryLogger(db.RW))

		// SQLite session manager
		sessionManager = infrahttp.NewSessionManager(db.RW)

		// SQLite repositories
		userRepo := adapters.NewSQLiteUserRepository(db.RW)
		productRepo := adapters.NewSQLiteProductRepository(db.RW)
		saleRepo := adapters.NewSQLiteSaleRepository(db.RW)
		inventoryRepo := adapters.NewSQLiteInventoryRepository(db.RW)
		clientRepo := adapters.NewSQLiteClientRepository(db.RW)
		configRepo := adapters.NewSQLiteConfigRepository(db.RW)

		// Use cases
		authUC := use_cases.NewAuthenticateUser(userRepo)
		createProductUC := use_cases.NewCreateProduct(productRepo)
		updateProductUC := use_cases.NewUpdateProduct(productRepo)
		listProductsUC := use_cases.NewListProducts(productRepo)
		deactivateProductUC := use_cases.NewDeactivateProduct(productRepo)
		registerSaleUC := use_cases.NewRegisterSale(saleRepo, productRepo, inventoryRepo)
		createClientUC := use_cases.NewCreateClient(clientRepo)
		listClientsUC := use_cases.NewListClients(clientRepo)

		// Handlers
		pageHandler = handlers.NewPageHandler(tmpl)
		chatHandler = handlers.NewChatHandler(nlsqlService, tmpl)
		metricsHandler = handlers.NewMetricsHandler(db.RW)
		authHandler = handlers.NewAuthHandler(authUC, tmpl, sessionManager)
		productHandler = handlers.NewProductHandler(createProductUC, updateProductUC, listProductsUC, deactivateProductUC, productRepo, tmpl)
		saleHandler = handlers.NewSaleHandler(registerSaleUC, tmpl, sessionManager)
		adminConfigHandler = handlers.NewAdminConfigHandler(configRepo, cryptoService, tmpl)
		clientHandler = handlers.NewClientHandler(createClientUC, listClientsUC, tmpl)

		cleanup = func() { db.Close() }
	}

	// Build router (shared between both modes)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Only compress in local mode — API Gateway handles compression for Lambda
	if cfg.AppEnv != "lambda" {
		r.Use(middleware.Compress(5))
	}
	r.Use(sessionManager.LoadAndSave)

	// Static files — only in local mode (CloudFront handles in Lambda)
	if cfg.AppEnv != "lambda" {
		fileServer := http.FileServer(http.Dir(StaticDir))
		r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	}

	// Health
	r.Get("/health", healthHandler(pool, cfg.AppEnv))

	// Public routes
	r.Get("/login", authHandler.LoginPage)
	r.Post("/login", authHandler.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(mw.RequireAuth(sessionManager))

		// Pages
		r.Get("/", pageHandler.Dashboard)
		r.Get("/productos", pageHandler.Products)
		r.Get("/ventas", pageHandler.Sales)
		r.Get("/metricas", pageHandler.Metrics)
		r.Post("/logout", authHandler.Logout)

		// Client routes
		r.Get("/clientes", clientHandler.List)
		r.Get("/clientes/new", clientHandler.CreateForm)
		r.Post("/clientes", clientHandler.Create)

		// API: Sales
		r.Post("/api/ventas", saleHandler.CompleteSale)
		r.Get("/api/ventas/recientes", metricsHandler.VentasRecientes)

		// HTMX fragment routes — NoCache to prevent stale responses
		r.Group(func(r chi.Router) {
			r.Use(mw.NoCache)

			// API: Metrics (HTMX fragments)
			r.Get("/api/metrics/ventas-hoy", metricsHandler.VentasHoy)
			r.Get("/api/metrics/ventas-semana", metricsHandler.VentasSemana)
			r.Get("/api/metrics/ventas-mes", metricsHandler.VentasMes)
			r.Get("/api/metrics/top-productos", metricsHandler.TopProductos)
			r.Get("/api/metrics/stock-bajo", metricsHandler.StockBajo)
			r.Get("/api/metrics/clientes-frecuentes", metricsHandler.ClientesFrecuentes)
			r.Get("/api/metrics/margen-categoria", metricsHandler.MargenCategoria)

			// API: Products (HTMX fragments)
			r.Get("/api/productos", metricsHandler.ProductosHTMX)
			r.Get("/api/productos/buscar", metricsHandler.ProductosBuscar)
		})
		r.Get("/productos/new", productHandler.CreateForm)
		r.Post("/productos", productHandler.Create)
		r.Get("/productos/{id}/edit", productHandler.EditForm)
		r.Post("/productos/{id}", productHandler.Edit)
		r.Delete("/productos/{id}", productHandler.Deactivate)

		// API: Chat NL→SQL
		r.Post("/api/chat", chatHandler.HandleChat)

		// Admin-only routes
		r.Group(func(r chi.Router) {
			r.Use(mw.RequireRole(sessionManager, "admin"))
			r.Get("/admin/config", adminConfigHandler.Show)
			r.Post("/admin/config", adminConfigHandler.Update)
		})
	})

	return r, cleanup, nil
}

// newPgxPool creates a pgxpool.Pool with Lambda-optimized settings.
func newPgxPool(databaseURL string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing database URL: %w", err)
	}

	poolConfig.MaxConns = 5
	poolConfig.MinConns = 1
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute
	poolConfig.ConnConfig.ConnectTimeout = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("creating pool: %w", err)
	}

	return pool, nil
}

// newPgxSessionManager creates an scs.SessionManager backed by PostgreSQL via pgxstore.
func newPgxSessionManager(pool *pgxpool.Pool) *scs.SessionManager {
	sm := scs.New()
	sm.Store = pgxstore.New(pool)
	sm.Lifetime = 8 * time.Hour
	sm.Cookie.Name = "pos_session"
	sm.Cookie.HttpOnly = true
	sm.Cookie.SameSite = http.SameSiteLaxMode
	sm.Cookie.Secure = true // HTTPS via API Gateway
	return sm
}

// loadTemplates walks the templates directory and parses all .html files.
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

// getSchema returns the database schema string used for NL-SQL system prompts.
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
