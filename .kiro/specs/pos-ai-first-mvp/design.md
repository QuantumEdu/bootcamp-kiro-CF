# Design: POS AI-First MVP

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    HTMX + Alpine.js + Tailwind              │
│                    (templates/ + static/)                    │
├─────────────────────────────────────────────────────────────┤
│                    Go HTTP Layer (chi router)                │
│                    infrastructure/http/                      │
├─────────────────────────────────────────────────────────────┤
│                    Application Layer                         │
│            use-cases/ (orchestrate domain + ports)           │
├─────────────────────────────────────────────────────────────┤
│                    Domain Layer                              │
│        entities/ + value-objects/ + ports/ (interfaces)      │
├─────────────────────────────────────────────────────────────┤
│                    Infrastructure Layer                      │
│    SQLite (sqlc) │ OpenRouter API │ Config │ Session Store   │
└─────────────────────────────────────────────────────────────┘
```

## Tech Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| Language | Go 1.22+ | Backend, static typing, fast compilation |
| Router | chi v5 | HTTP routing, middleware composition |
| Database | SQLite + sqlc | Persistence, type-safe queries |
| Frontend | HTMX + Alpine.js | Server-driven UI with minimal JS |
| Styling | Tailwind CSS (CDN) | Utility-first CSS, no build step |
| AI | OpenRouter API (GPT-4o-mini) | NL→SQL query generation |
| Auth | bcrypt + cookie sessions | PIN-based authentication |
| Session | alexedwards/scs + SQLite store | Server-side session management |

## Directory Structure

```
pos-ai-first/
├── cmd/
│   └── server/
│       └── main.go                    # Entry point, DI wiring
├── src/
│   ├── domain/
│   │   ├── entities/
│   │   │   ├── product.go            # Product entity + validation
│   │   │   ├── sale.go               # Sale + SaleItem entities
│   │   │   ├── user.go               # User entity + PIN validation
│   │   │   └── inventory.go          # InventoryMovement entity
│   │   ├── value_objects/
│   │   │   ├── money.go              # Money value object (avoid float issues)
│   │   │   ├── pin.go                # PIN hash value object
│   │   │   └── sql_query.go          # Validated SQL query (SELECT-only)
│   │   └── ports/
│   │       ├── product_repository.go # ProductRepository interface
│   │       ├── sale_repository.go    # SaleRepository interface
│   │       ├── user_repository.go    # UserRepository interface
│   │       ├── inventory_repository.go
│   │       └── ai_query_service.go   # AIQueryService interface (NL→SQL)
│   ├── application/
│   │   ├── use_cases/
│   │   │   ├── create_product.go
│   │   │   ├── register_sale.go
│   │   │   ├── authenticate_user.go
│   │   │   ├── process_natural_query.go  # NL→SQL orchestration
│   │   │   └── get_dashboard_metrics.go
│   │   └── dtos/
│   │       ├── product_dto.go
│   │       ├── sale_dto.go
│   │       └── query_response_dto.go
│   └── infrastructure/
│       ├── adapters/
│       │   ├── sqlite_product_repo.go
│       │   ├── sqlite_sale_repo.go
│       │   ├── sqlite_user_repo.go
│       │   ├── sqlite_inventory_repo.go
│       │   └── openrouter_query_service.go  # OpenRouter API adapter
│       ├── database/
│       │   ├── connection.go          # SQLite connection (RW + RO)
│       │   ├── migrations.go          # Auto-run migrations on start
│       │   └── queries/               # sqlc generated code
│       │       ├── products.sql
│       │       ├── sales.sql
│       │       ├── users.sql
│       │       ├── inventory.sql
│       │       └── metrics.sql
│       ├── http/
│       │   ├── router.go             # chi router setup + middleware
│       │   ├── middleware/
│       │   │   ├── auth.go           # Session check middleware
│       │   │   └── logging.go        # Request logging
│       │   └── handlers/
│       │       ├── products.go       # Product CRUD handlers
│       │       ├── sales.go          # Sales handlers
│       │       ├── auth.go           # Login/logout handlers
│       │       ├── chat.go           # NL→SQL chat handler
│       │       ├── dashboard.go      # Dashboard + metrics handlers
│       │       └── metrics.go        # Individual metric fragment handlers
│       └── config/
│           └── config.go             # Env vars loading
├── migrations/
│   └── 001_init.sql                  # Full schema (from wayfinder research)
├── templates/
│   ├── layout.html                   # Base layout (nav + chat bar)
│   ├── login.html
│   ├── dashboard.html
│   ├── products/
│   │   ├── list.html
│   │   ├── form.html
│   │   └── _row.html                # HTMX fragment
│   ├── sales/
│   │   ├── new.html
│   │   └── history.html
│   ├── chat/
│   │   ├── _bar.html                # Always-visible chat bar
│   │   └── _message.html            # Single message fragment
│   └── metrics/
│       ├── _ventas_hoy.html
│       ├── _top_productos.html
│       ├── _stock_bajo.html
│       └── _clientes_frecuentes.html
├── static/
│   └── js/
│       └── app.js                    # Alpine.js components (minimal)
├── testdata/
│   └── seed.sql                      # Realistic demo data
├── .env.example
├── .golangci.yml
├── Makefile
├── go.mod
└── README.md
```

## Data Model (SQLite)

8 tables as defined in `.wayfinder/research/001-init.sql`:
- `usuarios` — PIN auth, roles (admin/cajero)
- `categorias` — product groupings
- `productos` — catalog with stock, pricing, SKU
- `clientes` — customer basic info
- `ventas` — sale headers (total, payment method)
- `venta_items` — sale line items
- `inventario_movimientos` — stock change audit trail
- `configuracion` — key-value system config

Key design decisions:
- Spanish column names (matching user's domain language)
- `created_at` with `datetime('now','localtime')` for readable timestamps
- Trigger `trg_inventario_actualiza_stock` auto-updates product stock on inventory movement
- `REAL` type for quantities (supports kg/litro units)

## NL→SQL Flow

```
User types question
        ↓
POST /chat (HTMX)
        ↓
ProcessNaturalQuery use-case
        ↓
OpenRouter API call (GPT-4o-mini)
  - System prompt with schema + few-shot examples
  - response_format: json_schema (structured outputs)
        ↓
Parse JSON response {sql, explanation}
        ↓
Validate: SELECT-only? (keyword check + starts with SELECT/WITH)
        ↓
Execute on READ-ONLY SQLite connection
  - PRAGMA query_only = ON
  - 5-second timeout
  - LIMIT 500 appended
        ↓
Format results as human-readable Spanish text
        ↓
Return HTMX fragment (_message.html)
```

## Security Layers (NL→SQL)

1. **Prompt-level:** System prompt forbids DDL/DML
2. **Validation-level:** Go code checks query starts with SELECT/WITH, rejects keywords (INSERT, UPDATE, DELETE, DROP, ALTER, TRUNCATE)
3. **Connection-level:** Read-only SQLite connection (separate from app RW connection)
4. **Execution-level:** 5-second timeout, LIMIT 500 enforced
5. **Audit-level:** All generated queries logged with timestamp and user

## Session Management

- Library: `alexedwards/scs` with SQLite session store
- Session lifetime: 8 hours (POS shift duration)
- Cookie: HttpOnly, SameSite=Strict, Secure=false (localhost MVP)
- Session contains: user_id, user_name, user_role

## HTMX Patterns

- Dashboard metrics: `hx-get="/metrics/ventas-hoy" hx-trigger="load, every 30s" hx-target="#ventas-hoy"`
- Product list: `hx-get="/products" hx-trigger="load" hx-target="#product-list"`
- Chat: `hx-post="/chat" hx-target="#chat-messages" hx-swap="beforeend"`
- Forms: `hx-post="/products" hx-target="#product-list" hx-swap="afterbegin"`

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| github.com/go-chi/chi/v5 | v5.1.0 | HTTP router |
| github.com/mattn/go-sqlite3 | v1.14.22 | SQLite driver (CGO) |
| golang.org/x/crypto | latest | bcrypt for PINs |
| github.com/joho/godotenv | v1.5.1 | .env loading |
| github.com/alexedwards/scs/v2 | v2.8.0 | Session management |
| github.com/alexedwards/scs/sqlite3store | latest | SQLite session store |
