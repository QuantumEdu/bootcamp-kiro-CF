---
title: "Setup del proyecto Go"
labels: wayfinder:task
blocking: []
---

## Question

Crear el scaffolding del proyecto Go + HTMX + Tailwind + Alpine.js + sqlc.

Incluir:
- `go mod init`
- Estructura de carpetas: `cmd/`, `internal/handler/`, `internal/service/`, `internal/repository/`, `internal/model/`, `internal/middleware/`, `web/templates/`, `web/static/`
- Configuración de sqlc con SQLite driver
- Configuración de Tailwind CSS + Alpine.js (CDN o build)
- Template base HTML con HTMX + Alpine.js cargados
- net/http o chi como router
- Middleware básico de autenticación por PIN
- Makefile o taskfile con comandos: `run`, `build`, `test`, `sqlc-generate`
