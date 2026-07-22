---
inclusion: always
---
# Architecture — Hexagonal Ligera

## Estructura del proyecto

```
src/
├── domain/
│   ├── entities/       # Entidades de negocio con identidad
│   ├── value-objects/  # Objetos inmutables sin identidad
│   └── ports/          # Interfaces (contratos) que el dominio expone/requiere
├── application/
│   ├── use-cases/      # Un archivo por caso de uso, orquesta entidades y puertos
│   ├── services/       # Coordinación entre múltiples use-cases si es necesario
│   └── dtos/           # Estructuras de transferencia entre capas
└── infrastructure/
    ├── adapters/       # Implementaciones de puertos (repos, APIs externas, etc.)
    ├── database/       # Conexión, migraciones, queries (sqlc)
    ├── http/           # Handlers HTTP, middleware, rutas
    └── config/         # Variables de entorno, configuración de app
```

## Reglas no negociables

1. `domain/` NO importa frameworks, ORM, HTTP, filesystem ni servicios externos.
2. `application/` coordina casos de uso y depende SOLO del dominio (puertos).
3. `infrastructure/` implementa adaptadores y detalles técnicos.
4. Las dependencias apuntan hacia adentro: `infrastructure → application → domain`.
5. La lógica de negocio NO vive en handlers/controladores HTTP.
6. Cada puerto en `domain/ports/` es una interfaz Go; la implementación vive en `infrastructure/adapters/`.

## SOLID práctico (sin sobreingeniería)

| Principio | Aplicación concreta |
|-----------|-------------------|
| S — Single Responsibility | Un use-case = un archivo = una responsabilidad |
| O — Open/Closed | Nuevos adaptadores no modifican dominio ni use-cases existentes |
| L — Liskov Substitution | Cualquier implementación de un puerto es intercambiable sin romper tests |
| I — Interface Segregation | Puertos pequeños y específicos, no interfaces gigantes |
| D — Dependency Inversion | Use-cases dependen de puertos (interfaces), no de implementaciones |

## Cuándo NO aplicar hexagonal completa

- Scripts de migración o seed → archivo plano en `infrastructure/database/`
- Utilidades puras sin estado → `src/utils/` o función inline
- Configuración → `infrastructure/config/` directo
- Si un caso de uso solo hace CRUD sin reglas → un handler directo es válido, documentar por qué

## Convención de nombres (Go)

- Interfaces de puertos: `ProductRepository`, `AIQueryService` (sin prefijo `I`)
- Implementaciones: `SQLiteProductRepository`, `OpenRouterQueryService`
- Use-cases: `CreateProduct`, `RegisterSale`, `ProcessNaturalQuery`
- Handlers: `HandleCreateProduct`, `HandleChatQuery`
- Archivos: snake_case (`create_product.go`, `product_repository.go`)
