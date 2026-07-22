---
inclusion: always
---
# Project Conventions — POS AI-First

## Stack

- **Backend:** Go 1.22+ (net/http o chi router)
- **Frontend:** HTMX + Alpine.js + Tailwind CSS
- **Database:** SQLite + sqlc (type-safe queries)
- **AI:** OpenRouter API (Claude/GPT para NL→SQL)
- **Auth:** PIN-based multi-user

## Estructura de directorios

```
pos-ai-first/
├── src/
│   ├── domain/
│   │   ├── entities/
│   │   ├── value-objects/
│   │   └── ports/
│   ├── application/
│   │   ├── use-cases/
│   │   ├── services/
│   │   └── dtos/
│   └── infrastructure/
│       ├── adapters/
│       ├── database/
│       ├── http/
│       └── config/
├── migrations/          # SQL migrations
├── templates/           # HTMX templates
├── static/              # CSS, JS, assets
├── testdata/            # Fixtures para tests
├── cmd/
│   └── server/
│       └── main.go      # Entry point
├── go.mod
├── go.sum
├── .golangci.yml
├── .env.example
├── Makefile
└── README.md
```

## Commits

Conventional commits en inglés:
- `feat: add product CRUD`
- `fix: prevent negative stock on sale`
- `test: add table-driven tests for PIN validation`
- `refactor: extract query validation to value object`
- `docs: add API documentation`
- `chore: add golangci-lint config`

## Branches

- `main` — siempre desplegable
- `feat/xxx` — features nuevas
- `fix/xxx` — correcciones
- `refactor/xxx` — reestructuraciones sin cambio funcional

## Makefile targets

```makefile
.PHONY: run test lint fmt build migrate seed

run:
	go run cmd/server/main.go

test:
	go test ./... -v -cover

lint:
	golangci-lint run

fmt:
	gofmt -w .

build:
	go build -o bin/pos cmd/server/main.go

migrate:
	# sqlite3 migrations tool

seed:
	go run cmd/seed/main.go
```

## Dependencias aprobadas

| Dependencia | Propósito | Justificación |
|-------------|-----------|--------------|
| `github.com/go-chi/chi/v5` | Router HTTP | Ligero, idiomatic Go, middleware composable |
| `github.com/mattn/go-sqlite3` | SQLite driver | Estándar de facto para SQLite en Go |
| `github.com/sqlc-dev/sqlc` | SQL→Go codegen | Type-safe, sin ORM, queries explícitas |
| `golang.org/x/crypto` | bcrypt para PINs | Stdlib extendida, mantenida por Go team |
| `github.com/joho/godotenv` | .env loading | Simple, sin magia |

No agregar dependencias sin documentar la justificación.

## Variables de entorno

```env
PORT=8080
DATABASE_PATH=./data/pos.db
OPENROUTER_API_KEY=sk-...
OPENROUTER_MODEL=anthropic/claude-3-haiku
SESSION_SECRET=...
PIN_MAX_ATTEMPTS=5
PIN_LOCKOUT_MINUTES=5
QUERY_TIMEOUT_SECONDS=5
```

## Decisiones arquitectónicas registradas

1. SQLite sobre Postgres: simplicidad en 5 días, sin servidor externo
2. HTMX sobre SPA: menor complejidad frontend, server-driven
3. OpenRouter sobre Bedrock: acceso inmediato sin config AWS compleja para dev
4. sqlc sobre ORM: queries explícitas, type-safe, sin magia
5. PIN sobre OAuth: contexto POS real (el mesero no tiene email corporativo)
