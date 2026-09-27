# English demo completion
## Objective
Complete interrupted ES/EN app with English default and preserve redesigned 22-slide deck, translating presentation and video to English.
## Authorized scope
Local app, tests, presentacion.md, video.md, existing PPTX. No deploy, push, PR or remote sessions. Preserve existing dirty changes and unrelated registry edits.
## Configuration
Strict TDD: enabled by current user AGENTS instructions. Runner: go test ./... -v -cover; build: go build -o bin/pos cmd/server/main.go; lint: golangci-lint run.
RDD: unavailable (mode status printed off then failed unsafe authority ownership); do not enable or change ownership.
Delivery: ask-on-risk; user selected feature-branch-chain (integrator). Local commits only. Base: main at current HEAD.
Forecast: 1900–2650 new authored lines, generated PPTX excluded. Cohesive slices: foundation, remaining localized UI, English demo artifacts; record actual size and overage without code-golf. No PR creation authorized.
## Tasks
- [ ] T1 Restore bilingual foundation, remove duplicate metrics body, English default, safe language switch, login context, regression tests. Route delegated: multiple nontrivial files and read-before-write trigger.
- [ ] T2 Finish localized pages/fragments/errors/chat, persisted payment enum display labels, translation parity and rendering tests. Route delegated: multiple nontrivial files.
- [ ] T3 Translate presentacion.md/video.md and preserve 22-slide PPTX redesign in English. Remove unverified deployment/performance claims. Route delegated: document/deck specialist work.
## Acceptance and checks
Observed RED/GREEN/REFACTOR for new behavior; no changing persisted enum/database identifiers/user content. ES/EN both render without template errors; app defaults English. Docs match demonstrated implementation, Render deployment described as pending. Render deck and inspect layout/text.
Run focused tests, complete runner, build, lint; record unavailable/failed/pending honestly. Each task has Conventional work-unit commit, exact evidence, rollback boundary, authored line count. Independent verification if risk assessment fails. Parent spot check.
## Progress
Read-only audit complete. Default cache denied; prior TEMP-cache process unavailable, rerun tests. Source compile blocker: duplicated metrics implementation. No implementation performed yet.
## Next step
Delegate bounded writer sequential tasks, record commits and checks after each.

### T1 implementation evidence (commit pending)
Removed second appended metrics body; English fallback for missing/invalid cookies; explicit ES remains; login templates receive request Lang/T; language switch rejects unsupported values, stores HttpOnly SameSite=Lax cookie, rejects external/network-path/mismatched-origin redirects. Kept request-local redirect fallback `/`.
RED: focused regression tests failed default/invalid cookie, missing login .T, open redirect and unsafe cookie cases. Removed duplicate compile blocker before behavior tests. Initial Spanish test expectation corrected to existing `Entrar` label, not a source behavior change.
GREEN/REFACTOR: gofmt; `go test ./src/infrastructure/http/handlers ./src/infrastructure/i18n -v -cover` passed (25.0% and 61.5%); `go test ./... -v -cover` exit 0; `go build -o bin/pos cmd/server/main.go` exit 0. httptest smoke renders real login template EN/ES and exercises switch return paths/cookies. `git diff --check` passed.
`golangci-lint run` exit 1, 12 findings in existing health, postgres adapters, admin_config, sales, bootstrap config and lambda unused helper. Initial cache access warnings resolved by GOLANGCI_LINT_CACHE=$env:TEMP/bootcamp-lint-cache; findings unchanged. GOCACHE=$env:TEMP/bootcamp-go-cache. No claim all lint is clean.
Rollback boundary: foundation changes in handlers metrics/auth/lang, request template context, i18n package, route bootstrap and existing login/layout switch; do not revert unrelated registry/deck. Parent chooses cohesive inherited-change commit scope. T1 scoped HEAD diff metrics/auth = 78 authored lines, plus new lang/lang tests/i18n package = 486 lines; prior interrupted translations included, templates/bootstrap/helper count pending parent staging. New work also removed ~335 duplicate baseline lines; full tracked diff remains contaminated by inherited/unrelated work, so report staged exact count before commit.
Next: parent commits T1; resume T2 only after parent followup. T1 checkbox intentionally pending commit identity.
