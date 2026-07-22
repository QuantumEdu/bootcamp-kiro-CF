# Testing Capabilities — pos-ai-first

**Strict TDD Mode**: enabled (for critical paths)
**Updated**: 2026-07-22

## Test Runner

- Command: `go test ./...`
- Framework: Go testing (stdlib)
- Coverage: `go test -cover ./...`

## TDD Obligatorio en:

- PIN authentication (validation, lockout, session)
- Inventory operations (stock deduction, prevent negative)
- Financial calculations (sale totals, subtotals)
- NL→SQL security (SELECT-only validation, keyword rejection)
- Input validation (all domain entities)

## TDD No Obligatorio en:

- HTMX templates (presentational)
- Seed scripts
- Configuration loading
- CSS/styling changes

## Test Layers

| Layer | Available | Tool |
|-------|-----------|------|
| Unit | ✅ | Go testing stdlib |
| Integration | ✅ | SQLite in-memory (`:memory:`) |
| E2E | ❌ | Not in MVP scope |

## Quality Tools

| Tool | Available | Command |
|------|-----------|---------|
| Linter | ✅ | `golangci-lint run` |
| Type checker | ✅ | Go compiler (static typing) |
| Formatter | ✅ | `gofmt -w .` |
| Vet | ✅ | `go vet ./...` |

## Conventions

- Name: `TestUseCaseName_Scenario_ExpectedResult`
- Table-driven tests for 3+ variants
- Mocks: simple structs implementing port interfaces (no frameworks)
- Fixtures: `testdata/` directory
- Integration: SQLite `:memory:` for DB tests
