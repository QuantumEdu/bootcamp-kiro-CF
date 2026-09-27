# Guión de Video: POS AI-First MVP (5 minutos)

## Bootcamp Kiro × Código Facilito — Hackathon 2026

---

## Estructura temporal

| Sección | Duración | Acumulado |
|---------|----------|-----------|
| Intro + Problema | 0:30 | 0:30 |
| Solución + Demo conceptual | 0:45 | 1:15 |
| Cómo lo construí (Kiro completo) | 1:15 | 2:30 |
| Demo en vivo | 1:15 | 3:45 |
| Arquitectura + Seguridad | 0:30 | 4:15 |
| Cierre + Próximos pasos | 0:45 | 5:00 |

---

## Sección 1: Intro + Problema (0:00 – 0:30)

**Visual:** Pantalla con título del proyecto, luego cut a una persona abriendo hojas de cálculo frustrada.

**Narración:**
> "Hola, soy Lupita. Soy dueña de una taquería y se me complica estar registrando y ordenando pedidos y ventas, sobretodo saber cuánto vendí hoy. Actualmente abro Excel, filtro por fecha, sumo columnas... o le pido a alguien. ¿Y si pudiera simplemente preguntar, como en WhatsApp, '¿qué vendí hoy?' y recibir la respuesta al instante?"

**Notas de producción:** Transición rápida, energética. Máximo 2 tomas.

---

## Sección 2: Solución + Demo conceptual (0:30 – 1:15)

**Visual:** Screencast del chat bar del POS. Se escribe "¿Qué producto se vendió más esta semana?" y aparece la respuesta.

**Narración:**
> "Construí un POS que habla. El usuario escribe una pregunta en español, el sistema genera una consulta SQL segura usando AI, la ejecuta contra sus datos reales, y devuelve la respuesta formateada. Todo en menos de 5 segundos."
>
> "Pero no es solo chat. Es un POS completo: productos, ventas, clientes, dashboard con métricas en tiempo real, y el chat como feature diferenciador."

**Notas de producción:** Mostrar el flujo completo con overlay de las 5 capas de seguridad como badges.

---

## Sección 3: Cómo lo construí — Kiro completo (1:15 – 2:30)

**Visual:** Pantalla de Kiro IDE mostrando specs, steering, agents, hooks, MCP y el power empaquetado.

**Narración:**
> "Lo construí en 5 días usando Kiro. Pero no fue 'open IDE y empezar a codear'. Seguí un proceso estructurado con cada feature de Kiro:"

> "**Spec-driven development:** Tres specs completos — requirements, design, tasks — con ejecución en waves paralelas. 100 tareas ejecutadas sin conflictos gracias al grafo de dependencias."

> "**Steering files:** Seis archivos que definen cómo trabaja mi proyecto en cada sesión. Arquitectura hexagonal, testing obligatorio en auth y NL→SQL, seguridad, quality con golangci-lint, convenciones y patrones. Kiro los carga automáticamente — no tengo que repetir las reglas."

> "**Custom agents:** Dos agentes especializados. El NL→SQL Security Reviewer audita las 5 capas del pipeline de seguridad y reporta PASS, WARN o CRITICAL por capa. El Seed Data Generator crea datos realistas de taquería respetando las reglas del dominio."

> "**Hooks:** Siete hooks de automatización. Lint al guardar, tests al completar una tarea, gate de seguridad antes de iniciar cualquier tarea, validación del pipeline NL→SQL cuando cambian esos archivos, y captura automática de memoria para el LTM Power."

> "**MCP:** Dos servidores activos. SQLite MCP para inspeccionar la base de datos en vivo sin salir del IDE — fundamental para debuggear las queries generadas por AI. Context7 para documentación actualizada de chi, pgx y SCS directo en el contexto del agente."

> "**Powers:** Long-Term Memory para persistencia entre sesiones. Y el bonus — empaqué todo esto como un Kiro Power instalable con POWER.md, mcp.json y dos steering files de dominio."

**Notas de producción:** Mostrar brevemente cada feature mientras se menciona. Speed up en la navegación. Highlight en el Power empaquetado al final de esta sección.

---

## Sección 4: Demo en vivo (2:30 – 3:45)

**Visual:** Screencast del POS funcionando en Render (URL real en producción).

**Narración:**
> "Vamos a la demo. El POS está deployado en Render — gratis, con disco persistente para SQLite."

> "Login con PIN — autenticación bcrypt con lockout por intentos fallidos."

*[Muestra login con PIN]*

> "Dashboard: ventas de hoy, productos más vendidos, alertas de stock bajo. Se actualiza solo cada 30 segundos con HTMX."

*[Muestra dashboard con métricas]*

> "Registro de venta: selecciono productos, agrego al carrito, completo. El inventario se actualiza automáticamente."

*[Muestra flujo de venta]*

> "Y ahora, la estrella: '¿Cuántas ventas hubo esta semana?'"

*[Escribe en chat, espera respuesta]*

> "Respuesta en 3 segundos. SQL generado, validado, ejecutado en read-only con timeout. 5 capas de seguridad entre el LLM y mi base de datos."

**Notas de producción:** Pregrabar la demo como backup. Si es en vivo, tener datos seeded. La parte del chat es el clímax del video.

---

## Sección 5: Arquitectura + Seguridad (3:45 – 4:15)

**Visual:** Diagrama de arquitectura hexagonal. Luego los 5 escudos de seguridad NL→SQL.

**Narración:**
> "La arquitectura es hexagonal: el dominio no importa frameworks, los use-cases orquestan, y la infraestructura implementa los adaptadores. Gracias a esto, migré de SQLite a PostgreSQL para AWS con 7 archivos nuevos — cero cambios en el dominio."
>
> "El mismo código corre local con SQLite, en Render con SQLite persistente, o en Lambda con PostgreSQL. Un switch en el bootstrap según APP_ENV."
>
> "Seguridad NL→SQL: 5 capas. Prompt, validación Go, conexión read-only, timeout con LIMIT, y auditoría. No confiamos en el LLM — cada capa es un guardia independiente."

**Notas de producción:** Diagrama animado con los adaptadores intercambiables. 30 segundos máximo.

---

## Sección 6: Cierre (4:15 – 5:00)

**Visual:** Resumen de métricas + pantalla del Power empaquetado + QR al repo.

**Narración:**
> "En 5 días: 3 specs, 100+ tareas en paralelo, arquitectura hexagonal limpia, chat AI en español, deploy en Render con $0 de costo, y un Kiro Power empaquetado e instalable."
>
> "Todos los requisitos del hackathon cubiertos: spec-driven, steering, hooks, property tests, powers, MCP, custom agents — y el bonus del power empaquetado."
>
> "Kiro no es solo autocomplete. Es un sistema que entiende tu proyecto, mantiene contexto, ejecuta en paralelo respetando dependencias, y trabaja con las reglas que vos definís. Así se construye software en 2026."
>
> "Gracias."

**Notas de producción:** Grid animado de métricas. Screenshot del Power instalado en Kiro. QR al repo. Terminar en exactamente 5:00.

---

## Checklist de producción del video

- [ ] Grabar screencast de demo con datos seeded en Render
- [ ] Preparar backup de demo en caso de fallo de red
- [ ] Mostrar Kiro IDE con specs, steering, agents, hooks, MCP en la sección 3
- [ ] Mostrar el Power empaquetado (`kiro-power/POWER.md`) brevemente
- [ ] Grabar narración por separado (mejor audio)
- [ ] Editar con overlays de texto para puntos clave de cada feature
- [ ] Verificar que el video dura exactamente ≤ 5:00
- [ ] Exportar en 1080p mínimo
- [ ] Subir copia a Google Drive como backup
- [ ] Probar reproducción antes de la presentación

---

## Mapa de features Kiro → sección del video

| Feature Kiro | Minuto | Evidencia visual |
|-------------|--------|-----------------|
| Spec-driven development | 1:15 | `.kiro/specs/` con 3 specs abiertos |
| Steering documents | 1:25 | `.kiro/steering/` — 8 archivos |
| Custom agents | 1:35 | `.kiro/agents/nl-sql-security-reviewer.md` |
| Hooks | 1:45 | `.kiro/hooks/` — 7 hooks, uno disparándose |
| MCP | 1:55 | sqlite-pos conectado, query en vivo |
| Powers (LTM) | 2:05 | `ltm/runtime/active-context.json` |
| Power empaquetado | 2:15 | `kiro-power/POWER.md` con frontmatter |
| Property-based tests | — | Mencionar en arquitectura: `TestRequireRole_Property_*` |
