---
name: "pos-ai-first"
displayName: "POS AI-First"
description: "Complete development power for the POS AI-First system. Packages NL→SQL security agents, SQLite database inspection via MCP, and Go POS domain expertise into a single installable Kiro Power."
keywords: ["pos", "point-of-sale", "nl-sql", "sqlite", "go", "htmx", "hexagonal-architecture", "ai-first", "security", "mexico"]
author: "Gabriel Magallón"
---

# POS AI-First Power

## Overview

This power bundles everything a Kiro agent needs to work effectively on the **POS AI-First** codebase — a Go-based point-of-sale system with an AI-powered NL→SQL chat interface.

**What this power provides:**

- 🔒 **NL→SQL security expertise** — 5-layer validation pipeline for AI-generated SQL
- 🗃️ **Live database access** — SQLite MCP server for real-time data inspection
- 🌐 **Up-to-date library docs** — Context7 MCP for chi, pgx, scs, and AWS SDK
- 🤖 **Domain agents** — Pre-configured agents for security review and seed generation
- 📐 **Architecture guidance** — Hexagonal architecture patterns for Go POS systems

## Available MCP Servers

### sqlite-pos

Connects the agent to the live SQLite database at `./data/pos.db`.

**Tools available:**
- `read_query` — Execute SELECT queries against the POS database
- `list_tables` — List all tables in the schema
- `describe_table` — Show column definitions for any table
- `write_query` — Execute INSERT/UPDATE (use with caution — prefer read_query)

**Use cases in this project:**
- Inspect actual sales data while debugging NL→SQL queries
- Verify seed data was inserted correctly without leaving the IDE
- Check schema consistency between migrations and the live DB
- Validate that stock deductions are correct after a sale

**Setup:** Requires `uv` (`pip install uv`). Run `uvx mcp-server-sqlite` — no separate install needed.

### context7

Provides up-to-date documentation for Go libraries used in this project.

**Use cases:**
- `chi/v5` router middleware patterns
- `pgx/v5` PostgreSQL connection pool configuration
- `alexedwards/scs` session management API
- `modernc.org/sqlite` pure-Go driver options
- `aws-sdk-go-v2` Bedrock and Secrets Manager

**Setup:** Requires Node.js with `npx` (standard in most dev environments).

## Available Steering Files

- **`nl-sql-security.md`** — 5-layer NL→SQL security checklist and audit workflow
- **`domain-guide.md`** — POS domain concepts, entity rules, and business invariants

## Available Agents

### @nl-sql-security-reviewer

Audits the NL→SQL pipeline for all 5 security layers. Reports PASS / WARN / CRITICAL per layer.

**Invoke:** `@nl-sql-security-reviewer review the query validator`

### @seed-data-generator

Generates realistic Mexican food business seed data (taquerías, fondas). Respects domain rules: bcrypt PINs, consistent stock, coherent sale totals.

**Invoke:** `@seed-data-generator generate fresh demo data for SQLite`

## Architecture Quick Reference

```
src/
├── domain/          ← Pure Go. Zero external deps. Entities + ports.
├── application/     ← Use cases + NL→SQL service. Depends only on domain.
└── infrastructure/  ← Adapters (SQLite, PostgreSQL, OpenRouter, Bedrock)

cmd/
├── server/main.go   ← Local entry point (SQLite + OpenRouter)
└── lambda/main.go   ← AWS Lambda entry point (PostgreSQL + Bedrock)
```

**Dual-mode bootstrap:**
```go
switch cfg.AppEnv {
case "lambda":  // PostgreSQL + Bedrock + pgx sessions
default:        // SQLite + OpenRouter + SQLite sessions
}
```

## NL→SQL Security — 5 Layers

| Layer | What it does |
|-------|-------------|
| 1. Prompt | System prompt restricts LLM to SELECT-only |
| 2. Go validation | Whitelist SELECT/WITH, reject DDL/DML |
| 3. Connection | Separate read-only DB connection |
| 4. Execution | 5s timeout + LIMIT 500 |
| 5. Audit | Every generated query is logged before execution |

## Key Files Reference

| File | Purpose |
|------|---------|
| `src/application/nlsql/` | NL→SQL service, validator, formatter |
| `src/domain/ports/` | Repository and service interfaces |
| `src/infrastructure/adapters/` | SQLite + PostgreSQL + AI implementations |
| `internal/bootstrap/bootstrap.go` | Dual-mode dependency wiring |
| `migrations/001_init.sql` | SQLite schema (source of truth) |
| `migrations/postgres/001_init.sql` | PostgreSQL schema |
| `.kiro/steering/` | Architecture, testing, security, quality rules |
| `.kiro/hooks/` | 7 automation hooks (lint, tests, security gates) |

## Troubleshooting

### sqlite-pos MCP not connecting
- Verify `uv` is installed: `pip install uv`
- Check that `./data/pos.db` exists: run `make migrate` then `make seed` first
- Confirm the path in `mcp.json` matches your workspace root

### NL→SQL returns unexpected results
1. Invoke `@nl-sql-security-reviewer` to audit the validation layer
2. Check `src/application/nlsql/validator.go` for keyword whitelist
3. Verify the schema string in `internal/bootstrap/bootstrap.go → getSchema()`
4. Use `sqlite-pos` MCP to run the generated SQL manually and inspect results

### Build fails on Windows (CGO)
The project uses `modernc.org/sqlite` (pure Go, no CGO). If you see CGO errors, ensure you are NOT using `github.com/mattn/go-sqlite3` — check `go.mod`.
