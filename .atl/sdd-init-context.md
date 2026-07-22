# SDD Init Context — pos-ai-first

## Project

- **Name**: pos-ai-first
- **Description**: POS AI-First — Punto de venta con agente conversacional que entiende lenguaje natural
- **Bootcamp**: Código Facilito + Kiro (5-day timeline)
- **Status**: Scaffolded — structure created, implementation pending
- **Repo**: `https://github.com/QuantumEdu/bootcamp-kiro-CF`

## Stack

| Layer | Technology | Status |
|-------|-----------|--------|
| Backend | Go 1.22+ (chi router) | Ready |
| Frontend | HTMX + Alpine.js | Ready |
| Styling | Tailwind CSS (CDN) | Ready |
| Database | SQLite + sqlc | Ready |
| AI/NL | OpenRouter API (NL→SQL) | Ready |
| Auth | PIN-based (bcrypt + scs sessions) | Ready |

## Architecture

```
cmd/server/main.go → DI wiring → chi router
                                      ↓
src/infrastructure/http/handlers/ → use-cases → domain (ports)
                                      ↓
src/infrastructure/adapters/ → SQLite (sqlc) + OpenRouter API
```

Layout: split view — top: traditional app (CRUD + dashboard), bottom: always-visible chat bar.

## Directory Structure

```
pos-ai-first/
├── cmd/server/main.go          # Entry point
├── src/
│   ├── domain/
│   │   ├── entities/           # Product, Sale, User, Inventory
│   │   ├── value_objects/      # PIN, Money, ValidatedQuery
│   │   └── ports/              # Repository interfaces, AIQueryService
│   ├── application/
│   │   ├── use_cases/          # CreateProduct, RegisterSale, ProcessNaturalQuery...
│   │   └── dtos/               # Transfer objects between layers
│   └── infrastructure/
│       ├── adapters/           # SQLite repos, OpenRouter adapter
│       ├── database/           # Connection, migrations
│       ├── http/handlers/      # HTTP handlers (thin)
│       ├── http/middleware/    # Auth, logging
│       └── config/             # Env vars
├── migrations/001_init.sql     # SQLite schema (8 tables)
├── templates/                  # HTMX templates
├── static/                     # Tailwind CDN, Alpine.js, minimal JS
├── testdata/                   # Seed and fixtures
├── .kiro/                      # Kiro config (steering, specs, hooks)
├── Makefile                    # run, test, lint, build, seed
├── .golangci.yml               # Linter config
└── .env.example                # Environment template
```

## Kiro Spec

- **Location**: `.kiro/specs/pos-ai-first-mvp/`
- **Requirements**: 7 functional + 4 non-functional
- **Design**: hexagonal architecture, NL→SQL flow, security layers
- **Tasks**: 18 tasks across 5 days (Run All Tasks compatible)

## Governance

- `governance/PRD.md` — Full SaaS PRD (Next.js + Supabase for PRODUCTION, not bootcamp)
- `governance/project-brief.yaml` — Business context and vision
- **IMPORTANT**: Governance describes the FUTURE SaaS product. The bootcamp demo uses Go+SQLite stack.

## Wayfinder

- `.wayfinder/map.md` — Project map with 6 tickets
- `.wayfinder/research/` — Completed research (schema, NL→SQL, dashboard)
- `.wayfinder/tickets/` — 6 implementation tickets

## Testing

- **Strict TDD**: enabled for auth, inventory, finance, NL→SQL security
- **Runner**: `go test ./...`
- **Coverage**: `go test -cover ./...`
- **Linter**: `golangci-lint run`
- **Formatter**: `gofmt -w .`

## Next Steps

1. Execute Kiro Spec tasks (Run All Tasks in app.kiro.dev)
2. Day 1: Tasks 1-4 (scaffold, DB, entities, seed)
3. Day 2: Tasks 5-8 (ports, OpenRouter adapter, NL→SQL)
4. Day 3: Tasks 9-12 (auth, CRUD, UI, sales)
5. Day 4: Tasks 13-15 (dashboard, chat, router)
6. Day 5: Tasks 16-18 (polish, demo, testing)
