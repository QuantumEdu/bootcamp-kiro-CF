---
inclusion: always
---
# Design Patterns — Uso Pragmático

## Regla base

No introducir patrones por default. Cada patrón debe responder:

1. ¿Qué problema resuelve?
2. ¿Qué complejidad evita?
3. ¿Qué pasaría si no se usa?
4. ¿Hay una alternativa más simple?

Si no puedes responder las 4 preguntas, no uses el patrón.

---

## Repository

**Usar si:**
- Hay persistencia (SQLite, Postgres, etc.)
- Quieres aislar DB/ORM del dominio
- Quieres testear use-cases sin DB real (mock del puerto)

**No usar si:**
- Es un script simple o one-off
- Es prototipo desechable
- No hay lógica de negocio (CRUD puro sin reglas)

**En este proyecto:** SÍ — ProductRepository, SaleRepository, InventoryRepository como puertos en `domain/ports/`.

---

## Factory

**Usar si:**
- La creación del objeto tiene reglas de validación
- Hay variantes (ej: distintos tipos de venta)
- Hay invariantes de dominio que deben cumplirse al crear

**No usar si:**
- Solo se instancia un struct simple sin validación

**En este proyecto:** Considerar para `Sale` (tiene reglas: stock suficiente, precio positivo, items no vacíos).

---

## Strategy

**Usar si:**
- Hay 3 o más algoritmos intercambiables
- Se espera variabilidad real (no imaginaria)

**No usar si:**
- Solo existe una implementación
- La "variabilidad" es especulativa

**En este proyecto:** NO por ahora. Solo hay un proveedor de AI (OpenRouter). Si se agrega Bedrock como alternativa, entonces sí.

---

## Adapter

**Usar si:**
- Hay APIs externas (OpenRouter, Bedrock)
- Hay DB (SQLite ahora, Postgres después)
- Hay servicios reemplazables
- Hay filesystem, colas, webhooks o integraciones

**No usar si:**
- La integración es trivial y no se va a cambiar nunca

**En este proyecto:** SÍ — OpenRouterAdapter (implementa AIQueryPort), SQLiteAdapter (implementa repos).

---

## Patrones que NO usar en este proyecto

| Patrón | Razón para no usarlo |
|--------|---------------------|
| Singleton | Go usa dependency injection natural via constructores |
| Observer | No hay eventos complejos; HTMX maneja la reactividad |
| Abstract Factory | Un solo tipo de DB y un solo proveedor AI por ahora |
| Builder | Los DTOs y entities son simples, no necesitan builder |
| Decorator | Sin necesidad de wrapping dinámico en MVP |

---

## Señales de sobreingeniería

Si ves alguna de estas, simplifica:

- Una interfaz con una sola implementación que nunca va a cambiar
- Un package con un solo archivo de 20 líneas
- Más de 3 niveles de indirección para llegar a la lógica real
- Tests que prueban el framework, no tu lógica
- DTOs que son copia exacta de la entidad sin transformación
- Un "service" que solo delega al repository sin agregar lógica

> La mejor arquitectura es la más simple que cumple los requisitos actuales
> y permite crecer sin reescribir.
