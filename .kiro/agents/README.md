# Custom Agents — POS AI-First

Custom agents extend Kiro with project-specific expertise. Each agent is a specialized
reviewer or generator that understands the domain and architecture of this POS system.

## Available agents

### nl-sql-security-reviewer

**File:** `nl-sql-security-reviewer.md`

Audits the NL→SQL pipeline for security vulnerabilities. Checks all 5 security layers:
prompt restriction, Go-side SQL validation, read-only connection, execution timeout + LIMIT,
and query audit logging.

**Invoke when:** modifying `src/application/nlsql/`, AI adapters, or SQL validation logic.

**Output:** Per-layer PASS / WARN / CRITICAL report with actionable fixes.

---

### seed-data-generator

**File:** `seed-data-generator.md`

Generates realistic Mexican food business seed data for local (SQLite) and production
(PostgreSQL on Render or AWS RDS) environments. Enforces domain rules: bcrypt PINs,
consistent stock levels, coherent sale totals, spread-out dates for dashboard metrics.

**Invoke when:** setting up a fresh environment, preparing a demo, or adding new entities
that need seed coverage.

**Output:** Complete Go seed file ready to run with `go run cmd/seed/main.go`.

---

## How to invoke a custom agent

In Kiro chat, mention the agent by name:

```
@nl-sql-security-reviewer review the latest changes to the query validator
@seed-data-generator generate fresh demo data for the hackathon presentation
```

Or open the agent file directly and use it as context for your request.
