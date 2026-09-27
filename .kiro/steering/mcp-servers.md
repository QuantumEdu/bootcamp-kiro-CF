---
inclusion: manual
---

# MCP Servers — Model Context Protocol

## Configured servers

This project uses two MCP servers defined in `.kiro/settings/mcp.json`.

### sqlite-pos

```json
{
  "command": "uvx",
  "args": ["mcp-server-sqlite", "--db-path", "./data/pos.db"]
}
```

**Purpose:** Gives the agent direct read access to the live SQLite database during development.

**Used for:**
- Inspecting actual data while debugging NL→SQL queries ("show me the last 5 sales rows")
- Verifying seed data was inserted correctly without leaving the IDE
- Checking schema consistency between `migrations/001_init.sql` and the live DB
- Validating that stock deductions are correct after a sale is registered

**When to activate:** `#mcp-servers` in chat, or when debugging data-layer issues.

**Setup requirement:** `uv` must be installed (`pip install uv` or `brew install uv`).
Once installed, `uvx` downloads and runs `mcp-server-sqlite` automatically — no separate install needed.

---

### context7

```json
{
  "command": "npx",
  "args": ["-y", "@upstash/context7-mcp@latest"]
}
```

**Purpose:** Provides up-to-date documentation for libraries used in this project directly inside the agent context.

**Used for:**
- `chi/v5` router middleware patterns
- `pgx/v5` PostgreSQL connection pool configuration
- `alexedwards/scs` session management
- `aws-sdk-go-v2` Bedrock and Secrets Manager API
- `modernc.org/sqlite` pure-Go SQLite driver options

**When to activate:** When working with external library APIs or when official docs are needed without leaving the IDE.

**Setup requirement:** Node.js with `npx` available (standard in most dev environments).

---

## Why MCP matters in this project

The NL→SQL feature generates SQL from user questions in Spanish. During development,
`sqlite-pos` allowed direct inspection of the live database to verify that:

1. Generated queries returned the expected rows
2. Schema changes in migrations propagated correctly
3. The read-only connection truly blocked write operations

Without MCP, this verification required switching to a separate DB client (TablePlus, DBeaver)
or writing ad-hoc debug scripts. MCP keeps the full development loop inside Kiro.
