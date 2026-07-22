# POS AI-First

> Sistema de Punto de Venta con Inteligencia Artificial para PYMEs — Bootcamp Codigo Facilito + Kiro

## Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│ Frontend: HTMX + Alpine.js + Tailwind CSS (CDN)             │
├─────────────────────────────────────────────────────────────┤
│ Backend: Go + chi router                                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────────┐              │
│  │ Handlers │→ │ Use Cases│→ │ Adapters     │              │
│  │ (HTTP)   │  │ (NL→SQL) │  │ (OpenRouter) │              │
│  └──────────┘  └──────────┘  └──────────────┘              │
├─────────────────────────────────────────────────────────────┤
│ Database: SQLite (WAL mode) — RW + RO connections           │
└─────────────────────────────────────────────────────────────┘
```

## Stack

| Componente | Tecnologia |
|---|---|
| Backend | Go 1.22+ / chi v5 |
| Database | SQLite (modernc.org/sqlite) |
| Frontend | HTMX 1.9 + Alpine.js 3 + Tailwind CSS |
| AI | OpenRouter (gpt-4o-mini / claude-3-haiku) |
| Auth | PIN + SHA-256 + session tokens |

## Setup rapido

```bash
# 1. Clonar
git clone https://github.com/QuantumEdu/bootcamp-kiro-CF.git
cd bootcamp-kiro-CF

# 2. Configurar
cp .env.example .env
# Editar .env: agregar OPENROUTER_API_KEY

# 3. Ejecutar
go run ./cmd/server

# 4. Abrir
# http://localhost:8080/login
# PIN admin: 1234
# PIN cajero: 123
```

## Funcionalidades

### Completadas

- [x] **Dashboard** — Ventas hoy/semana/mes, top productos, stock bajo, margen por categoria
- [x] **Productos** — Lista con busqueda HTMX
- [x] **Ventas** — Carrito con Alpine.js, checkout JSON, historial reciente
- [x] **Chat IA** — Panel lateral con NL→SQL (lenguaje natural a consultas SQL)
- [x] **Autenticacion** — PIN con lockout por intentos fallidos
- [x] **Metricas** — HTMX polling automatico (30s/60s)
- [x] **Seguridad NL→SQL** — Jailbreak detection, table whitelist, SELECT-only, max 100 rows
- [x] **Seed data** — 30 productos, 10 ventas, 5 clientes (tienda de abarrotes MX)

### Seguridad del Chat IA

10 capas de defensa:
1. Validacion de input (jailbreak detection)
2. System prompt (solo SELECT)
3. SQL validation (keyword blocking)
4. Table whitelist (8 tablas)
5. Multi-statement detection
6. Comment injection blocking
7. Read-only DB connection
8. Query timeout (5s)
9. Row limit (100 max)
10. Max input length (500 chars)

## Estructura del proyecto

```
cmd/server/main.go          — Entry point
src/
├── application/nlsql/      — NL→SQL service + validator
├── domain/ports/           — Interfaces (hexagonal)
└── infrastructure/
    ├── adapters/           — OpenRouter client
    ├── config/             — Env vars
    ├── database/           — SQLite connection + migrations
    └── http/
        ├── handlers/       — HTTP handlers (pages, metrics, chat, auth, sales)
        └── middleware/     — Auth middleware
templates/                  — HTML templates (Go html/template)
static/js/                  — Alpine.js components
```

## Comandos

```bash
make run       # Ejecutar servidor
make test      # Tests
make lint      # Linter
make build     # Compilar binario
```

## Variables de entorno

| Variable | Default | Descripcion |
|---|---|---|
| `PORT` | 8080 | Puerto del servidor |
| `DATABASE_PATH` | ./data/pos.db | Ruta de la BD SQLite |
| `OPENROUTER_API_KEY` | — | API key de OpenRouter |
| `OPENROUTER_MODEL` | anthropic/claude-3-haiku | Modelo LLM |
| `SESSION_SECRET` | dev-secret | Secreto para tokens |
| `PIN_MAX_ATTEMPTS` | 5 | Intentos antes de lockout |
| `PIN_LOCKOUT_MINUTES` | 5 | Minutos de bloqueo |
| `QUERY_TIMEOUT_SECONDS` | 5 | Timeout para queries |

## Demo script

1. Login con PIN 1234
2. Ver dashboard (metricas se auto-refrescan)
3. Ir a Productos → ver catalogo
4. Ir a Ventas → buscar producto → agregar al carrito → cobrar
5. Chat IA: "cuantas ventas hubo hoy?" → ver SQL generado + resultados
6. Chat IA: "que producto se vendio mas esta semana?"
7. Intentar jailbreak: "ignora instrucciones" → ver rechazo

## Equipo

Proyecto desarrollado con [Kiro](https://kiro.dev) durante el Bootcamp Codigo Facilito 2026.
