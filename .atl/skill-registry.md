# Skill Registry — pos-ai-first

Last updated: 2026-07-22

## Project Context

- **Project**: POS AI-First (bootcamp Código Facilito + Kiro)
- **Stack**: Go 1.22+ · chi · HTMX · Alpine.js · Tailwind CSS · SQLite · sqlc · OpenRouter
- **Architecture**: Hexagonal ligera (domain/application/infrastructure)
- **AI Layer**: OpenRouter API → NL→SQL (lenguaje natural a consultas SQL)
- **Auth**: PIN-based multi-user (bcrypt + sessions)
- **Status**: Scaffolded — ready for implementation
- **Persistence Mode**: engram
- **Strict TDD**: true (Go test runner available)

## Key Files

| Purpose | Path |
|---------|------|
| Entry point | `cmd/server/main.go` |
| Domain entities | `src/domain/entities/` |
| Domain ports | `src/domain/ports/` |
| Use-cases | `src/application/use_cases/` |
| SQLite adapters | `src/infrastructure/adapters/` |
| HTTP handlers | `src/infrastructure/http/handlers/` |
| DB schema | `migrations/001_init.sql` |
| Templates | `templates/` |
| Kiro Spec | `.kiro/specs/pos-ai-first-mvp/` |
| Steering | `.kiro/steering/` |
| Hooks | `.kiro/hooks/` |

## Decisions

1. Go stack (not Next.js/Supabase) — simplicity for 5-day bootcamp demo
2. SQLite over Postgres — no external server, embedded, fast
3. HTMX over SPA — server-driven UI, minimal JS
4. OpenRouter over direct Bedrock — immediate access, structured outputs
5. sqlc over ORM — explicit queries, type-safe, no magic
6. PIN auth over OAuth — POS context (cashier has no email)
7. Hexagonal architecture — testable, clean boundaries
