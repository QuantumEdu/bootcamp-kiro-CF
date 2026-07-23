# Implementation Plan

## Overview

MVP de un sistema POS (Point of Sale) AI-First para el bootcamp Código Facilito + Kiro. Incluye scaffolding del proyecto Go con arquitectura hexagonal, base de datos SQLite, autenticación por PIN, CRUD de productos/ventas, dashboard de métricas, y chat conversacional NL→SQL via OpenRouter API. Desarrollo planificado en 5 días.

## Tasks

- [x] 1. Scaffold Go project with hexagonal structure
  - [x] 1.1 Run `go mod init github.com/QuantumEdu/bootcamp-kiro-CF`
  - [x] 1.2 Create directory structure: `cmd/server/`, `src/domain/entities/`, `src/domain/value_objects/`, `src/domain/ports/`, `src/application/use_cases/`, `src/application/dtos/`, `src/infrastructure/adapters/`, `src/infrastructure/database/`, `src/infrastructure/http/handlers/`, `src/infrastructure/http/middleware/`, `src/infrastructure/config/`
  - [x] 1.3 Create `cmd/server/main.go` with basic chi server (health check endpoint)
  - [x] 1.4 Create `Makefile` with targets: run, test, lint, fmt, build
  - [x] 1.5 Create `.env.example` with all required env vars
  - [x] 1.6 Create `.golangci.yml` with configured linters
  - [x] 1.7 Install dependencies: chi, modernc.org/sqlite, godotenv (x/crypto and scs pending until auth task)
  - [x] 1.8 Verify `make run` starts server on :8080 and `make test` passes
- [x] 2. Database layer — schema + connection + migrations
  - [x] 2.1 Create `migrations/001_init.sql` with full schema (8 tables, indexes, triggers from wayfinder research)
  - [x] 2.2 Create `src/infrastructure/database/connection.go` — open SQLite with WAL mode, create RW and RO connections (using modernc.org/sqlite pure-Go driver)
  - [x] 2.3 Create `src/infrastructure/database/migrations.go` — auto-run migrations on startup (embed SQL files)
  - [x] 2.4 Create `src/infrastructure/config/config.go` — load env vars (PORT, DATABASE_PATH, OPENROUTER_API_KEY, etc.)
  - [x] 2.5 Add session table for `alexedwards/scs` SQLite store
  - [x] 2.6 Wire database initialization in `main.go`
  - [x] 2.7 Write integration test: database creates tables correctly with in-memory SQLite
- [x] 3. Domain entities and value objects
  - [x] 3.1 Create `src/domain/entities/product.go` — Product struct with validation (name not empty, price > 0)
  - [x] 3.2 Create `src/domain/entities/user.go` — User struct with role enum (admin/cajero)
  - [x] 3.3 Create `src/domain/entities/sale.go` — Sale + SaleItem structs with validation (quantity > 0, total >= 0)
  - [x] 3.4 Create `src/domain/entities/inventory.go` — InventoryMovement with type enum (entrada/salida/ajuste)
  - [x] 3.5 Create `src/domain/value_objects/pin.go` — PIN hashing and comparison using bcrypt
  - [x] 3.6 Create `src/domain/value_objects/sql_query.go` — ValidatedQuery type that only holds SELECT statements
  - [x] 3.7 Write unit tests for all entities and value objects (table-driven)
- [x] 4. Seed data for demo
  - [x] 4.1 Create `testdata/seed.sql` with realistic data: 3 users (1 admin, 2 cajeros), 5 categories, 20 products, 5 clients, 50 sales with items, inventory movements
  - [x] 4.2 Create `cmd/seed/main.go` to run seed against database
  - [x] 4.3 Add `seed` target to Makefile
- [x] 5. Domain ports (interfaces)
  - [x] 5.1 Create `src/domain/ports/product_repository.go` — CRUD + filter + low-stock query
  - [x] 5.2 Create `src/domain/ports/sale_repository.go` — create sale + list sales + sale by ID
  - [x] 5.3 Create `src/domain/ports/user_repository.go` — find by PIN hash, find by ID, increment attempts, lock/unlock
  - [x] 5.4 Create `src/domain/ports/inventory_repository.go` — create movement, get movements by product
  - [x] 5.5 Create `src/domain/ports/ai_query_service.go` — GenerateSQL(question string) (sql string, explanation string, err error)
  - [x] 5.6 Create `src/domain/ports/metrics_repository.go` — ventas hoy/semana/mes, top products, low stock, frequent customers
- [x] 6. OpenRouter adapter (NL→SQL)
  - [x] 6.1 Create `src/infrastructure/adapters/openrouter_query_service.go`
  - [x] 6.2 Implement HTTP client calling OpenRouter API with structured outputs (json_schema)
  - [x] 6.3 Include system prompt with schema dump, table descriptions in Spanish, glossary, and few-shot examples
  - [x] 6.4 Parse response JSON into sql + explanation
  - [x] 6.5 Handle errors: timeout (retry with gpt-4o-mini), malformed response, API unavailable
  - [x] 6.6 Write test with mock HTTP server
- [x] 7. SQL validation and execution
  - [x] 7.1 Implement `ValidateSelectOnly(sql string) error` in `sql_query.go` — check starts with SELECT/WITH, reject dangerous keywords
  - [x] 7.2 Create read-only query executor: opens RO connection, sets PRAGMA query_only, adds LIMIT 500, enforces 5-second timeout
  - [x] 7.3 Create query logging: store every generated query with timestamp, user_id, question, generated_sql, success/failure
  - [x] 7.4 Write unit tests for validation (valid SELECT, CTE, malicious inputs, SQL injection attempts)
  - [x] 7.5 Write integration test: execute validated query against in-memory SQLite
- [x] 8. ProcessNaturalQuery use-case
  - [x] 8.1 Create `src/application/use_cases/process_natural_query.go`
  - [x] 8.2 Orchestrate: receive question → call AIQueryService → validate SQL → execute on RO connection → format results → return response
  - [x] 8.3 Handle all error cases: API failure, invalid SQL, execution error, timeout
  - [x] 8.4 Format results as human-readable Spanish text
  - [x] 8.5 Write tests with mocked ports
- [x] 9. Authentication use-case and handler
  - [x] 9.1 Create `src/application/use_cases/authenticate_user.go` — validate PIN, check lockout, create session
  - [x] 9.2 Create SQLite user repository adapter
  - [x] 9.3 Create `src/infrastructure/http/handlers/auth.go` — login page, login POST, logout
  - [x] 9.4 Configure `alexedwards/scs` with SQLite session store
  - [x] 9.5 Create auth middleware: check session, redirect to login if missing
  - [x] 9.6 Create role middleware: check user role matches required role
  - [x] 9.7 Write tests: valid login, invalid PIN, lockout after 5 attempts, session expiry
- [x] 10. Product CRUD use-cases and handlers
  - [x] 10.1 Create use-cases: CreateProduct, UpdateProduct, ListProducts, DeactivateProduct
  - [x] 10.2 Create SQLite product repository adapter (using sqlc or manual queries)
  - [x] 10.3 Create `src/infrastructure/http/handlers/products.go` — list, create form, create POST, edit form, edit POST, deactivate
  - [x] 10.4 Each handler returns HTMX fragments for seamless updates
  - [x] 10.5 Write tests for product validation rules
- [x] 11. Base layout and templates
  - [x] 11.1 Create `templates/layout.html` — full page layout with Tailwind CDN, HTMX CDN, Alpine.js CDN, sidebar nav, chat bar placeholder
  - [x] 11.2 Create `templates/login.html` — PIN input form
  - [x] 11.3 Create `templates/products/list.html` — product table with edit/deactivate buttons
  - [x] 11.4 Create `templates/products/form.html` — create/edit product form
  - [x] 11.5 Create `templates/products/_row.html` — single product row (HTMX fragment)
  - [x] 11.6 Create `templates/chat/_bar.html` — always-visible chat input at bottom
  - [x] 11.7 Create `templates/chat/_message.html` — single chat message bubble
  - [x] 11.8 Configure Go html/template with layout composition
- [x] 12. Sales registration
  - [x] 12.1 Create `src/application/use_cases/register_sale.go` — validate items, check stock, calculate total, deduct inventory, create sale
  - [x] 12.2 Create SQLite sale repository adapter
  - [x] 12.3 Create SQLite inventory repository adapter
  - [x] 12.4 Create `src/infrastructure/http/handlers/sales.go` — new sale page, add item, complete sale
  - [x] 12.5 Create `templates/sales/new.html` — POS-style sale capture (product search + cart)
  - [x] 12.6 Write tests: successful sale, insufficient stock, negative quantity rejected
- [x] 13. Dashboard metrics
  - [x] 13.1 Create `src/application/use_cases/get_dashboard_metrics.go` — fetch all metrics
  - [x] 13.2 Create SQLite metrics repository with queries from wayfinder research (ventas hoy/semana/mes, top 5, stock bajo, clientes frecuentes)
  - [x] 13.3 Create `src/infrastructure/http/handlers/dashboard.go` — main dashboard page
  - [x] 13.4 Create `src/infrastructure/http/handlers/metrics.go` — individual metric fragment endpoints
  - [x] 13.5 Create `templates/dashboard.html` — full dashboard with HTMX polling
  - [x] 13.6 Create metric fragments: `_ventas_hoy.html`, `_top_productos.html`, `_stock_bajo.html`, `_clientes_frecuentes.html`
  - [x] 13.7 Configure HTMX triggers: ventas-hoy every 30s, others every 60s
- [x] 14. Chat handler and UI integration
  - [x] 14.1 Create `src/infrastructure/http/handlers/chat.go` — POST /chat endpoint
  - [x] 14.2 Wire ProcessNaturalQuery use-case to handler
  - [x] 14.3 Return formatted response as HTMX fragment (_message.html)
  - [x] 14.4 Handle loading state in UI (htmx-indicator)
  - [x] 14.5 Handle error states (friendly messages in Spanish)
  - [x] 14.6 Test 5 key queries: "¿qué vendí hoy?", "¿qué producto se vendió más?", "¿cuántas ventas hubo esta semana?", "¿qué productos tienen stock bajo?", "¿quiénes son mis clientes frecuentes?"
- [x] 15. Router wiring and middleware stack
  - [x] 15.1 Wire all handlers in `src/infrastructure/http/router.go`
  - [x] 15.2 Apply middleware chain: logging → session → auth (on protected routes)
  - [x] 15.3 Configure static file serving for `/static/`
  - [x] 15.4 Configure template rendering
  - [x] 15.5 Verify all routes work end-to-end
- [x] 16. UI polish and responsive design
  - [x] 16.1 Polish Tailwind styles: consistent spacing, colors, typography
  - [x] 16.2 Add empty states for all lists (products, sales, metrics)
  - [x] 16.3 Add loading skeletons for HTMX fragments
  - [x] 16.4 Ensure chat bar works on mobile (responsive)
  - [x] 16.5 Add toast/notification for successful operations (Alpine.js)
- [ ] 17. Demo preparation
  - [x] 17.1 Run seed with 20+ products, 50+ sales across multiple days
  - [x] 17.2 Test all 5 NL→SQL demo queries with real data
  - [x] 17.3 Create `README.md` with setup instructions, architecture diagram, demo script
  - [x] 17.4 Ensure `make run` works cleanly from fresh clone
  - [~] 17.5 Record backup video of key flows (in case of connectivity issues at demo)
- [ ] 18. Final testing and edge cases
  - [x] 18.1 Run full test suite: `make test` passes
  - [x] 18.2 Run linter: `make lint` passes with zero warnings
  - [x] 18.3 Test error cases: OpenRouter timeout, invalid queries, empty database
  - [x] 18.4 Test auth: wrong PIN, lockout, expired session
  - [x] 18.5 Test concurrent access: multiple browser tabs
  - [~] 18.6 Final commit with clean git history
- [ ] 19. Presentation and video preparation
  - [~] 19.1 Complete `presentacion.md` with final metrics, screenshots, and demo captures
  - [~] 19.2 Complete `video.md` script with final narration text and timing
  - [~] 19.3 Record 5-minute demo video following the script
  - [~] 19.4 Export video in 1080p and upload backup to Google Drive
  - [~] 19.5 Test presentation flow end-to-end (slide transitions, timing)
- [ ] 20. GitHub Projects and documentation setup
  - [~] 20.1 Create GitHub Project board with columns: Backlog, In Progress, Review, Done
  - [~] 20.2 Link all 18 task issues to the project board
  - [~] 20.3 Add custom fields: Day, Priority, Estimation
  - [~] 20.4 Update README.md with project board link, architecture diagram, and setup instructions
  - [~] 20.5 Final documentation review: verify all links, badges, and references are correct

## Task Dependency Graph

```json
{
  "waves": [
    {
      "wave": 1,
      "tasks": [1],
      "description": "Project scaffolding and initial setup"
    },
    {
      "wave": 2,
      "tasks": [2, 3],
      "description": "Database layer and domain entities (depend on task 1)"
    },
    {
      "wave": 3,
      "tasks": [4, 5],
      "description": "Seed data and domain ports (depend on tasks 2, 3)"
    },
    {
      "wave": 4,
      "tasks": [6, 7, 9, 10, 11],
      "description": "Adapters, validation, auth, product CRUD, and templates (depend on tasks 2, 5)"
    },
    {
      "wave": 5,
      "tasks": [8, 12],
      "description": "NL→SQL use-case and sales registration (depend on tasks 6, 7, 9, 10)"
    },
    {
      "wave": 6,
      "tasks": [13, 14],
      "description": "Dashboard metrics and chat integration (depend on tasks 8, 10, 11, 12)"
    },
    {
      "wave": 7,
      "tasks": [15],
      "description": "Router wiring connecting all handlers (depends on tasks 9–14)"
    },
    {
      "wave": 8,
      "tasks": [16],
      "description": "UI polish and responsive design (depends on tasks 11, 15)"
    },
    {
      "wave": 9,
      "tasks": [17],
      "description": "Demo preparation (depends on tasks 4, 15, 16)"
    },
    {
      "wave": 10,
      "tasks": [18],
      "description": "Final testing and edge cases (depends on tasks 15, 17)"
    },
    {
      "wave": 11,
      "tasks": [19, 20],
      "description": "Presentation, video, and project board (depend on tasks 17, 18)"
    }
  ]
}
```

## Notes

### Day Groupings

The tasks are organized across a 5-day bootcamp schedule:

- **Day 1 — Foundation** (Tasks 1–4): Project scaffolding, database, domain entities, seed data
- **Day 2 — AI Core + Security** (Tasks 5–8): Domain ports, OpenRouter adapter, SQL validation, NL→SQL use-case
- **Day 3 — Authentication + Product CRUD + UI Shell** (Tasks 9–12): Auth, product management, templates, sales
- **Day 4 — Dashboard + Chat Integration** (Tasks 13–15): Metrics, chat UI, router wiring
- **Day 5 — Polish + Demo Prep** (Tasks 16–20): UI polish, demo preparation, final testing, presentation, project board
