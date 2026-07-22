---
inclusion: always
---
# Quality & Linting — Estándares de Código

## Herramientas requeridas

| Herramienta | Propósito | Comando |
|-------------|-----------|---------|
| `gofmt` | Formateo estándar Go | `gofmt -w .` |
| `go vet` | Análisis estático básico | `go vet ./...` |
| `staticcheck` | Linter avanzado | `staticcheck ./...` |
| `golangci-lint` | Meta-linter (agrupa varios) | `golangci-lint run` |

## Configuración golangci-lint

Crear `.golangci.yml` en raíz:

```yaml
run:
  timeout: 3m

linters:
  enable:
    - errcheck       # Errores no manejados
    - govet          # Análisis estático
    - staticcheck    # Bugs comunes
    - unused         # Código muerto
    - gosimple       # Simplificaciones
    - ineffassign    # Asignaciones inútiles
    - gocritic       # Style + performance
    - revive         # Reglas de estilo configurables

linters-settings:
  revive:
    rules:
      - name: exported
        arguments: [checkPrivateReceivers]
      - name: unused-parameter

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

## Reglas de código

### Errores
- Todo error se maneja o se propaga con contexto: `fmt.Errorf("creating product: %w", err)`
- Nunca `_ = someFunction()` si retorna error (excepto `defer file.Close()`)
- Errores de dominio son tipos propios: `ErrProductNotFound`, `ErrInsufficientStock`

### Funciones
- Máximo 40 líneas por función (señal de que necesita split)
- Máximo 3 parámetros; si hay más, usar struct de opciones
- Retornar early (guard clauses), no anidar ifs

### Packages
- Un package = una responsabilidad coherente
- No packages circulares (Go no lo permite, pero diseñar para evitarlo)
- Exports mínimos: solo lo que otra capa necesita

### Naming
- Variables: cortas en scope corto (`p` para product en un loop), descriptivas en scope largo
- Errores: `Err` + sustantivo (`ErrNotFound`, `ErrInvalidPIN`)
- Interfaces: nombre del comportamiento (`Reader`, `ProductRepository`), no `IProduct`

### Comments
- Toda función exportada lleva comment godoc (una línea mínimo)
- Comments explican el POR QUÉ, no el QUÉ (el código dice el qué)
- TODO/FIXME con contexto: `// TODO(ticket-003): agregar paginación`

## Pre-commit checklist (manual o hook)

```bash
gofmt -w .
go vet ./...
golangci-lint run
go test ./...
```

## Métricas de calidad aceptables

| Métrica | Umbral |
|---------|--------|
| Cobertura dominio | ≥ 90% |
| Cobertura application | ≥ 80% |
| Cobertura infrastructure | ≥ 60% |
| Lint warnings | 0 (CI falla si hay) |
| Funciones > 40 líneas | 0 |
| Errores no manejados | 0 |
