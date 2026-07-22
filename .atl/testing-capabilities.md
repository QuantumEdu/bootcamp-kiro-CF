# Testing Capabilities — pos-ai-first

**Strict TDD Mode**: disabled
**Detected**: 2025-07-22

## Test Runner

- Command: `go test ./...` (planned — not yet available)
- Framework: Go testing (stdlib)

## Test Layers

| Layer       | Available | Tool        |
| ----------- | --------- | ----------- |
| Unit        | ❌        | Go testing (planned) |
| Integration | ❌        | SQLite in-memory (planned) |
| E2E         | ❌        | — |

## Coverage

- Available: ❌ (planned: `go test -cover ./...`)
- Command: —

## Quality Tools

| Tool         | Available | Command        |
| ------------ | --------- | -------------- |
| Linter       | ❌        | `golangci-lint run` (planned) |
| Type checker | ✅        | Go compiler (static typing) |
| Formatter    | ❌        | `gofmt -w .` (planned) |

## Notes

- Project is in pre-implementation phase
- No source code exists yet
- Testing capabilities will be activated after `go mod init` and project scaffolding
- Re-run `/sdd-init` after scaffold to activate Strict TDD Mode
