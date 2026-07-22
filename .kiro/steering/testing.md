---
inclusion: always
---
# Testing — Estrategia Pragmática

## TDD obligatorio en:

- Autenticación (PIN validation, session management)
- Autorización (permisos por rol)
- Inventario crítico (stock updates, prevent negative)
- Cálculos financieros (totales, impuestos, cambio)
- NL→SQL (validación y sanitización de queries generadas)
- Validaciones de seguridad (input sanitization, SQL injection prevention)
- Lógica con muchos casos borde (descuentos, promociones)

## TDD no obligatorio en:

- UI presentacional (templates HTMX)
- Scripts temporales o de seed
- Prototipos desechables
- Configuración
- Cambios cosméticos (CSS, copy)

## Tests esperados — prioridad

1. Casos felices (happy path)
2. Casos borde (límites, cero, máximos)
3. Entradas inválidas (tipos incorrectos, strings vacíos, negativos)
4. Permisos denegados (PIN incorrecto, rol sin acceso)
5. Errores esperados (DB down, API timeout, input malformado)
6. Regresiones de bugs (todo bug corregido lleva test que lo reproduce)
7. Reglas críticas de dominio (invariantes de negocio)

## Estructura de tests

```
src/
├── domain/
│   ├── entities/
│   │   └── product_test.go          # Tests de entidad pura
│   └── value-objects/
│       └── money_test.go            # Tests de value objects
├── application/
│   └── use-cases/
│       ├── create_product_test.go   # Tests con mocks de puertos
│       └── process_query_test.go
└── infrastructure/
    ├── adapters/
    │   └── sqlite_product_repo_test.go  # Tests de integración con SQLite in-memory
    └── http/
        └── handlers_test.go         # Tests de handlers (HTTP)
```

## Convenciones Go para tests

- Nombre: `TestNombreUseCase_Escenario_ResultadoEsperado`
- Table-driven tests cuando hay 3+ variantes del mismo caso
- Mocks: interfaces de puertos se mockean con structs simples, no frameworks pesados
- Fixtures: usar `testdata/` si se necesitan archivos de prueba
- SQLite in-memory para integration tests: `sql.Open("sqlite3", ":memory:")`

## Ejemplo de test table-driven

```go
func TestCreateProduct_Validation(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateProductInput
        wantErr bool
    }{
        {"valid product", CreateProductInput{Name: "Coca", Price: 20}, false},
        {"empty name", CreateProductInput{Name: "", Price: 20}, true},
        {"negative price", CreateProductInput{Name: "Coca", Price: -1}, true},
        {"zero price", CreateProductInput{Name: "Coca", Price: 0}, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := usecase.Execute(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("got err=%v, wantErr=%v", err, tt.wantErr)
            }
        })
    }
}
```

## Cobertura mínima

- Dominio: 90%+ (es lógica pura, no hay excusa)
- Application/use-cases: 80%+
- Infrastructure: 60%+ (los adapters se testean con integración)
- HTTP handlers: testear contract (status codes, response shape), no lógica
