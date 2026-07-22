# Mapa: POS AI-First MVP

**Tracker:** local-markdown (`.wayfinder/tickets/`)

---

## Destination

MVP del POS AI-First — app Go + HTMX + SQLite + OpenRouter que entiende lenguaje natural. Un dueño de negocio puede preguntarle "¿qué vendí esta semana?" y obtener respuesta en castellano.

## Notes

- **Stack definido:** Go, HTMX, Tailwind CSS, Alpine.js, SQLite + sqlc, OpenRouter (NL→SQL)
- **Autenticación:** login con PIN multi-usuario (inspirado en POS_Tilapia_Go)
- **Hosting MVP:** localhost
- **Layout:** sistema partido — arriba app tradicional (CRUD + dashboard), abajo barra de chat siempre visible
- **Dashboard:** métricas fijas + posibilidad de generarlas vía SQL
- **NL→SQL:** dinámico desde el MVP, usando OpenRouter
- Skills relevantes: `sdd-*` para implementación, `go-testing`, `brainstorming`, `grilling`
- Proyecto relacionado: `POS_Tilapia_Go` — experiencia previa con POS en Go, sqlc, SQLite

## Decisions so far

_(vacío — estamos charteando)_

## Not yet specified

Niebla que se resolverá cuando el frontier avance:

- **Seguridad NL→SQL:** ¿cómo evitamos queries destructivas (DROP, DELETE, UPDATE)? ¿read-only wrapper, whitelist de tablas, validación post-generación?
- **Estructura exacta del layout:** sidebar vs topbar, transiciones entre pantallas CRUD, comportamiento responsive del chat
- **Testing strategy:** ¿tests unitarios con sqlc mocks? ¿tests de integración con SQLite in-memory? ¿tests del NL→SQL con fixtures?
- **API design:** ¿rutas HTMX directas desde Go net/http o chi? ¿o capa REST intermedia?
- **Alcance de inventario:** ¿solo stock actual o movimientos (entradas/salidas/ajustes)? ¿alertas de stock mínimo?
- **Clientes en el MVP:** ¿datos básicos (nombre, teléfono) o algo más?
- **Dashboard KPIs exactos:** qué métricas fijas, queries, visualización
- **Estructura del proyecto Go:** carpeta layout, patrones (repo, service, handler)
- **CI/CD post-MVP:** Docker image, VPS deploy

## Out of scope

_(vacío — aún no descartamos nada)_

---

## Tickets (hijos del mapa)

| # | Ticket | Tipo | Estado | Bloquea |
|---|--------|------|--------|---------|
| 001 | Modelo de datos del MVP | Research | Abierto | 002, 003, 005 |
| 002 | Setup del proyecto Go | Task | Abierto | — |
| 003 | Layout y navegación UI | Grilling | Abierto | — |
| 004 | OpenRouter + NL→SQL prompt design | Research | Abierto | — |
| 005 | Dashboard de métricas | Research | Abierto | 003 |
| 006 | Seguridad del NL→SQL | Grilling | Abierto | — |
