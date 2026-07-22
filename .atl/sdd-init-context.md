# SDD Init Context — pos-ai-first

## Project

- **Name**: pos-ai-first
- **Description**: POS AI-First — Punto de venta con agente conversacional que entiende lenguaje natural
- **Bootcamp**: Código Facilito + Kiro (5-day timeline)
- **Status**: Pre-implementation (planning phase complete, scaffolding pending)
- **Repo root**: `d:\02-A\code\bootcamp\`

## Detected Stack

| Layer | Technology | Status |
|-------|-----------|--------|
| Backend | Go (net/http or chi) | Planned |
| Frontend | HTMX + Alpine.js | Planned |
| Styling | Tailwind CSS | Planned |
| Database | SQLite + sqlc | Planned |
| AI/NL | OpenRouter API (NL→SQL) | Planned |
| Auth | PIN-based multi-user | Planned |

## Architecture (from .wayfinder/map.md)

```
User → HTMX UI → Go handlers → SQLite (products, sales, inventory)
                              → OpenRouter API (NL→SQL→response)
```

Layout: split view — top: traditional app (CRUD + dashboard), bottom: always-visible chat bar.

## Governance Source

- **Primary**: `governance/PRD.md` — full SaaS PRD (Next.js + Supabase version for production)
- **Bootcamp scope**: `.wayfinder/map.md` — Go + HTMX + SQLite MVP for 5-day demo
- **Evaluation**: Impacto tecnológico (30%) + Innovación (30%) + Software funcional (30%) + Uso AWS/Kiro (10%)

## Persistence

- **Mode**: engram
- **Artifact store**: engram (no openspec/ directory)
- **Skill registry**: `.atl/skill-registry.md`

## Strict TDD

- **Status**: false
- **Reason**: No test runner available — project not yet scaffolded
- **Plan**: Activate `strict_tdd: true` when `go mod init` is executed and test runner is available (`go test ./...`)

## Key Decisions

1. **Bootcamp demo uses Go stack** (wayfinder) not Next.js+Supabase (governance PRD)
2. **NL→SQL is dynamic from MVP** using OpenRouter
3. **SQLite** for simplicity in 5-day timeline
4. **HTMX** for minimal JavaScript with server-driven UI
5. **PIN auth** inspired by POS_Tilapia_Go prior project

## Unresolved (from .wayfinder/map.md)

- Security for NL→SQL (prevent destructive queries)
- Testing strategy (unit with sqlc mocks? integration with SQLite in-memory?)
- API design (HTMX direct or REST intermediate layer?)
- Exact project structure
- Dashboard KPIs
- CI/CD post-MVP

## Relevant Skills for This Project

| Skill | Relevance |
|-------|-----------|
| `go-testing` | Go test patterns when scaffolded |
| `grilling` | Stress-test designs before building |
| `implement` | Implementation based on PRD/tickets |
| `to-issues` | Break wayfinder tickets into GitHub issues |
| `codebase-design` | Deep module design |
| `diagnosing-bugs` | Debug during implementation |
| `work-unit-commits` | Reviewable commit units |
| `context7-mcp` | Library/framework reference lookups |

## Next Steps

1. `/sdd-explore` or `/sdd-new` to start implementing ticket 002 (project setup)
2. Scaffold Go project with `go mod init`
3. Once test runner exists, re-run `/sdd-init` to activate Strict TDD
