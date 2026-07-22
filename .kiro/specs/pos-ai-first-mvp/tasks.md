# Tasks: POS AI-First MVP

## Day 1 — Foundation

### Task 1: Scaffold Go project with hexagonal structure
- [ ] Run `go mod init github.com/QuantumEdu/bootcamp-kiro-CF`
- [ ] Create directory structure: `cmd/server/`, `src/domain/entities/`, `src/domain/value_objects/`, `src/domain/ports/`, `src/application/use_cases/`, `src/application/dtos/`, `src/infrastructure/adapters/`, `src/infrastructure/database/`, `src/infrastructure/http/handlers/`, `src/infrastructure/http/middleware/`, `src/infrastructure/config/`
- [ ] Create `cmd/server/main.go` with basic chi server (health check endpoint)
- [ ] Create `Makefile` with targets: run, test, lint, fmt, build
- [ ] Create `.env.example` with all required env vars
- [ ] Create `.golangci.yml` with configured linters
- [ ] Install dependencies: chi, go-sqlite3, godotenv, x/crypto, scs
- [ ] Verify `make run` starts server on :8080 and `make test` passes

### Task 2: Database layer — schema + connection + migrations
- [ ] Create `migrations/001_init.sql` with full schema (8 tables, indexes, triggers from wayfinder research)
- [ ] Create `src/infrastructure/database/connection.go` — open SQLite with WAL mode, create RW and RO connections
- [ ] Create `src/infrastructure/database/migrations.go` — auto-run migrations on startup (embed SQL files)
- [ ] Create `src/infrastructure/config/config.go` — load env vars (PORT, DATABASE_PATH, OPENROUTER_API_KEY, etc.)
- [ ] Add session table for `alexedwards/scs` SQLite store
- [ ] Wire database initialization in `main.go`
- [ ] Write integration test: database creates tables correctly with in-memory SQLite

### Task 3: Domain entities and value objects
- [ ] Create `src/domain/entities/product.go` — Product struct with validation (name not empty, price > 0)
- [ ] Create `src/domain/entities/user.go` — User struct with role enum (admin/cajero)
- [ ] Create `src/domain/entities/sale.go` — Sale + SaleItem structs with validation (quantity > 0, total >= 0)
- [ ] Create `src/domain/entities/inventory.go` — InventoryMovement with type enum (entrada/salida/ajuste)
- [ ] Create `src/domain/value_objects/pin.go` — PIN hashing and comparison using bcrypt
- [ ] Create `src/domain/value_objects/sql_query.go` — ValidatedQuery type that only holds SELECT statements
- [ ] Write unit tests for all entities and value objects (table-driven)

### Task 4: Seed data for demo
- [ ] Create `testdata/seed.sql` with realistic data: 3 users (1 admin, 2 cajeros), 5 categories, 20 products, 5 clients, 50 sales with items, inventory movements
- [ ] Create `cmd/seed/main.go` to run seed against database
- [ ] Add `seed` target to Makefile

---

## Day 2 — AI Core + Security

### Task 5: Domain ports (interfaces)
- [ ] Create `src/domain/ports/product_repository.go` — CRUD + filter + low-stock query
- [ ] Create `src/domain/ports/sale_repository.go` — create sale + list sales + sale by ID
- [ ] Create `src/domain/ports/user_repository.go` — find by PIN hash, find by ID, increment attempts, lock/unlock
- [ ] Create `src/domain/ports/inventory_repository.go` — create movement, get movements by product
- [ ] Create `src/domain/ports/ai_query_service.go` — GenerateSQL(question string) (sql string, explanation string, err error)
- [ ] Create `src/domain/ports/metrics_repository.go` — ventas hoy/semana/mes, top products, low stock, frequent customers

### Task 6: OpenRouter adapter (NL→SQL)
- [ ] Create `src/infrastructure/adapters/openrouter_query_service.go`
- [ ] Implement HTTP client calling OpenRouter API with structured outputs (json_schema)
- [ ] Include system prompt with schema dump, table descriptions in Spanish, glossary, and few-shot examples
- [ ] Parse response JSON into sql + explanation
- [ ] Handle errors: timeout (retry with gpt-4o-mini), malformed response, API unavailable
- [ ] Write test with mock HTTP server

### Task 7: SQL validation and execution
- [ ] Implement `ValidateSelectOnly(sql string) error` in `sql_query.go` — check starts with SELECT/WITH, reject dangerous keywords
- [ ] Create read-only query executor: opens RO connection, sets PRAGMA query_only, adds LIMIT 500, enforces 5-second timeout
- [ ] Create query logging: store every generated query with timestamp, user_id, question, generated_sql, success/failure
- [ ] Write unit tests for validation (valid SELECT, CTE, malicious inputs, SQL injection attempts)
- [ ] Write integration test: execute validated query against in-memory SQLite

### Task 8: ProcessNaturalQuery use-case
- [ ] Create `src/application/use_cases/process_natural_query.go`
- [ ] Orchestrate: receive question → call AIQueryService → validate SQL → execute on RO connection → format results → return response
- [ ] Handle all error cases: API failure, invalid SQL, execution error, timeout
- [ ] Format results as human-readable Spanish text
- [ ] Write tests with mocked ports

---

## Day 3 — Authentication + Product CRUD + UI Shell

### Task 9: Authentication use-case and handler
- [ ] Create `src/application/use_cases/authenticate_user.go` — validate PIN, check lockout, create session
- [ ] Create SQLite user repository adapter
- [ ] Create `src/infrastructure/http/handlers/auth.go` — login page, login POST, logout
- [ ] Configure `alexedwards/scs` with SQLite session store
- [ ] Create auth middleware: check session, redirect to login if missing
- [ ] Create role middleware: check user role matches required role
- [ ] Write tests: valid login, invalid PIN, lockout after 5 attempts, session expiry

### Task 10: Product CRUD use-cases and handlers
- [ ] Create use-cases: CreateProduct, UpdateProduct, ListProducts, DeactivateProduct
- [ ] Create SQLite product repository adapter (using sqlc or manual queries)
- [ ] Create `src/infrastructure/http/handlers/products.go` — list, create form, create POST, edit form, edit POST, deactivate
- [ ] Each handler returns HTMX fragments for seamless updates
- [ ] Write tests for product validation rules

### Task 11: Base layout and templates
- [ ] Create `templates/layout.html` — full page layout with Tailwind CDN, HTMX CDN, Alpine.js CDN, sidebar nav, chat bar placeholder
- [ ] Create `templates/login.html` — PIN input form
- [ ] Create `templates/products/list.html` — product table with edit/deactivate buttons
- [ ] Create `templates/products/form.html` — create/edit product form
- [ ] Create `templates/products/_row.html` — single product row (HTMX fragment)
- [ ] Create `templates/chat/_bar.html` — always-visible chat input at bottom
- [ ] Create `templates/chat/_message.html` — single chat message bubble
- [ ] Configure Go html/template with layout composition

### Task 12: Sales registration
- [ ] Create `src/application/use_cases/register_sale.go` — validate items, check stock, calculate total, deduct inventory, create sale
- [ ] Create SQLite sale repository adapter
- [ ] Create SQLite inventory repository adapter
- [ ] Create `src/infrastructure/http/handlers/sales.go` — new sale page, add item, complete sale
- [ ] Create `templates/sales/new.html` — POS-style sale capture (product search + cart)
- [ ] Write tests: successful sale, insufficient stock, negative quantity rejected

---

## Day 4 — Dashboard + Chat Integration

### Task 13: Dashboard metrics
- [ ] Create `src/application/use_cases/get_dashboard_metrics.go` — fetch all metrics
- [ ] Create SQLite metrics repository with queries from wayfinder research (ventas hoy/semana/mes, top 5, stock bajo, clientes frecuentes)
- [ ] Create `src/infrastructure/http/handlers/dashboard.go` — main dashboard page
- [ ] Create `src/infrastructure/http/handlers/metrics.go` — individual metric fragment endpoints
- [ ] Create `templates/dashboard.html` — full dashboard with HTMX polling
- [ ] Create metric fragments: `_ventas_hoy.html`, `_top_productos.html`, `_stock_bajo.html`, `_clientes_frecuentes.html`
- [ ] Configure HTMX triggers: ventas-hoy every 30s, others every 60s

### Task 14: Chat handler and UI integration
- [ ] Create `src/infrastructure/http/handlers/chat.go` — POST /chat endpoint
- [ ] Wire ProcessNaturalQuery use-case to handler
- [ ] Return formatted response as HTMX fragment (_message.html)
- [ ] Handle loading state in UI (htmx-indicator)
- [ ] Handle error states (friendly messages in Spanish)
- [ ] Test 5 key queries: "¿qué vendí hoy?", "¿qué producto se vendió más?", "¿cuántas ventas hubo esta semana?", "¿qué productos tienen stock bajo?", "¿quiénes son mis clientes frecuentes?"

### Task 15: Router wiring and middleware stack
- [ ] Wire all handlers in `src/infrastructure/http/router.go`
- [ ] Apply middleware chain: logging → session → auth (on protected routes)
- [ ] Configure static file serving for `/static/`
- [ ] Configure template rendering
- [ ] Verify all routes work end-to-end

---

## Day 5 — Polish + Demo Prep

### Task 16: UI polish and responsive design
- [ ] Polish Tailwind styles: consistent spacing, colors, typography
- [ ] Add empty states for all lists (products, sales, metrics)
- [ ] Add loading skeletons for HTMX fragments
- [ ] Ensure chat bar works on mobile (responsive)
- [ ] Add toast/notification for successful operations (Alpine.js)

### Task 17: Demo preparation
- [ ] Run seed with 20+ products, 50+ sales across multiple days
- [ ] Test all 5 NL→SQL demo queries with real data
- [ ] Create `README.md` with setup instructions, architecture diagram, demo script
- [ ] Ensure `make run` works cleanly from fresh clone
- [ ] Record backup video of key flows (in case of connectivity issues at demo)

### Task 18: Final testing and edge cases
- [ ] Run full test suite: `make test` passes
- [ ] Run linter: `make lint` passes with zero warnings
- [ ] Test error cases: OpenRouter timeout, invalid queries, empty database
- [ ] Test auth: wrong PIN, lockout, expired session
- [ ] Test concurrent access: multiple browser tabs
- [ ] Final commit with clean git history
