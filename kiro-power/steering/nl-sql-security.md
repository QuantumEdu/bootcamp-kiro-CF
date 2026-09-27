# NL→SQL Security — Audit Workflow

## When to run this audit

- Before merging any PR that touches `src/application/nlsql/`
- When adding a new AI provider adapter
- When changing the SQL validation or sanitization logic
- When the NL→SQL security hook fires during development

## Layer 1 — Prompt restriction

**File:** `src/application/nlsql/prompt_builder.go` (or equivalent)

**Check:**
- Does the system prompt contain an explicit instruction like "only generate SELECT queries"?
- Is the table whitelist embedded in the prompt (not just in Go validation)?
- Is user input passed to the prompt without interpolating it directly into SQL?

**Pass criteria:** System prompt explicitly restricts output to SELECT. Tables are listed.

---

## Layer 2 — Go-side SQL validation

**File:** `src/application/nlsql/validator.go` (or equivalent)

**Check:**
- Is there a whitelist of allowed statement types (SELECT, WITH)?
- Are these keywords rejected: INSERT, UPDATE, DELETE, DROP, ALTER, TRUNCATE, CREATE, EXEC, EXECUTE?
- Is the check case-insensitive (uppercases the query before checking)?
- Does validation happen BEFORE the query reaches the DB connection?

**Pass criteria:** Whitelist enforced in Go, case-insensitive, pre-execution.

---

## Layer 3 — Read-only connection

**File:** `src/infrastructure/database/` or `internal/bootstrap/bootstrap.go`

**Check:**
- Is there a separate `readDB` / `RO` connection distinct from `writeDB` / `RW`?
- For SQLite: is the read connection opened with `_query_only=true` or equivalent?
- For PostgreSQL: is the connection using a read-only role or `SET TRANSACTION READ ONLY`?
- Does the NL→SQL service receive `readDB`, not `writeDB`?

**Pass criteria:** NL→SQL service only has access to the read-only connection.

---

## Layer 4 — Execution limits

**File:** `src/application/nlsql/` service execution logic

**Check:**
- Is there a context timeout on query execution (max 5 seconds)?
- Is a LIMIT clause injected or enforced if not present in the generated SQL?
- Are query results capped (max 500 rows returned to the LLM)?

**Pass criteria:** Timeout enforced via `context.WithTimeout`. LIMIT present.

---

## Layer 5 — Audit logging

**File:** `src/application/nlsql/logger.go` or equivalent

**Check:**
- Is every generated query logged BEFORE execution (not after)?
- Does the log record: timestamp, user ID, original NL question, generated SQL?
- Are SQL execution errors returned as friendly user messages (no raw DB error exposed)?
- Is the audit log written to the write DB (not just stdout)?

**Pass criteria:** Pre-execution log with full context. Friendly error messages.

---

## Reporting format

```
✅ PASS     — Layer 1: System prompt restricts to SELECT, table whitelist present
✅ PASS     — Layer 2: Go validator whitelist enforced, case-insensitive
⚠️  WARN    — Layer 3: Read-only connection exists but not verified for PostgreSQL mode
❌ CRITICAL — Layer 4: No timeout found in query execution path
✅ PASS     — Layer 5: Pre-execution logging with user ID and NL question

VERDICT: DO NOT MERGE — fix Layer 4 CRITICAL before proceeding.
```
