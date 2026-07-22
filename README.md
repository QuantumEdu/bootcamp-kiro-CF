# POS AI-First — Bootcamp Código Facilito + Kiro

Un sistema POS (Point of Sale) con inteligencia artificial que entiende lenguaje natural. El dueño de negocio puede preguntar "¿qué vendí esta semana?" y obtener respuesta inmediata de sus datos reales.

## Stack

| Capa | Tecnología |
|------|-----------|
| Backend | Go 1.22+ (chi router) |
| Frontend | HTMX + Alpine.js + Tailwind CSS (CDN) |
| Base de datos | SQLite + sqlc |
| AI/NL | OpenRouter API (GPT-4o-mini → NL→SQL) |
| Auth | PIN-based multi-user (bcrypt) |
| Sessions | alexedwards/scs + SQLite store |

## Arquitectura

```
src/
├── domain/           # Entidades, value objects, puertos (interfaces)
├── application/      # Use-cases, DTOs
└── infrastructure/   # Adapters (SQLite, OpenRouter), HTTP handlers, config
```

Hexagonal ligera — las dependencias apuntan hacia adentro: `infrastructure → application → domain`.

## Desarrollo rápido

```bash
# Requisitos: Go 1.22+, gcc (para CGO/SQLite)
cp .env.example .env    # Agregar OPENROUTER_API_KEY
make run                # Inicia en http://localhost:8080
make test               # Corre tests
make lint               # golangci-lint
make seed               # Cargar datos demo
```

## Estructura del repositorio

```
├── cmd/server/main.go       # Entry point
├── src/                     # Código fuente (hexagonal)
├── migrations/              # SQL schema
├── templates/               # HTMX templates
├── static/                  # CSS, JS, assets
├── testdata/                # Fixtures y seed data
├── .kiro/                   # Steering, specs, hooks (Kiro IDE config)
├── governance/              # PRD y brief (visión SaaS futura)
├── .wayfinder/              # Tickets y research de planeación
└── Pre-analisis/            # Documentos de evaluación previa
```

## Demo (5 funcionalidades clave)

1. **Login con PIN** — autenticación tipo POS real
2. **CRUD Productos** — catálogo con categorías y stock
3. **Registro de ventas** — carrito + descuento automático de inventario
4. **Dashboard** — métricas en tiempo real con HTMX polling
5. **Chat NL→SQL** — "¿qué vendí ayer?" → respuesta instantánea

## Evaluación Bootcamp

- Impacto tecnológico (30%) — NL→SQL para PYMES sin tech
- Innovación (30%) — POS conversacional desde el diseño
- Software funcional (30%) — CRUD + chat + dashboard demostrable
- Uso de AWS y Kiro (10%) — Kiro para desarrollo, steering, specs
