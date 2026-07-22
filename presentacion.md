# Presentación: POS AI-First MVP

## Bootcamp Kiro + Código Facilito — Hackathon 2026

---

## Diapositiva 1: Portada

**Título:** POS AI-First MVP — Tu negocio responde preguntas  
**Subtítulo:** Bootcamp Kiro × Código Facilito | Hackathon 2026  
**Comentarios para el presentador:** Saluda al público, presenta tu nombre y el proyecto en una frase: "Un punto de venta donde le preguntas a tus datos y te responden."  
**Propuesta de diseño:** Fondo oscuro con gradiente azul-violeta, logo de Kiro + Código Facilito, mockup del chat POS en el centro.

---

## Diapositiva 2: El Problema

**Título:** ¿Por qué un POS inteligente?  
**Descripción:**
- Los dueños de pequeños negocios no tienen tiempo de revisar reportes
- Las preguntas son simples: "¿Qué vendí hoy?", "¿Qué se está agotando?"
- La respuesta debería ser tan fácil como preguntar en WhatsApp

**Comentarios para el presentador:** Conecta con la audiencia — "Imaginen que el dueño de una taquería quiere saber sus ventas del día sin abrir Excel ni pedir ayuda."  
**Propuesta de diseño:** Ilustración de un dueño de negocio frustrado con hojas de cálculo vs. una burbuja de chat simple.

---

## Diapositiva 3: La Solución

**Título:** Chat NL→SQL: Pregunta en español, obtén respuestas de tus datos  
**Descripción:**
- El usuario escribe: "¿Qué producto se vendió más esta semana?"
- El sistema genera SQL seguro via OpenRouter API
- Ejecuta contra SQLite (read-only, validado, con timeout)
- Devuelve respuesta formateada en español

**Comentarios para el presentador:** Demo rápida mental — "Es literalmente escribir una pregunta y recibir la respuesta como si hablaras con un asistente."  
**Propuesta de diseño:** Screenshot del chat bar con una pregunta y respuesta real. Flechas mostrando el flujo: pregunta → AI → SQL → datos → respuesta.

---

## Diapositiva 4: Objetivo del Hackathon

**Título:** Construir un MVP funcional en 5 días  
**Descripción:**
- Criterios de evaluación: Impacto tecnológico (30%), Innovación (30%), Software funcional (30%), Uso de AWS + Kiro (10%)
- Meta: POS completo con CRUD, ventas, dashboard, y el chat AI como diferenciador
- Stack: Go + HTMX + SQLite + OpenRouter

**Comentarios para el presentador:** Enfatiza que es un proyecto real ejecutable, no un mockup.  
**Propuesta de diseño:** Timeline horizontal de 5 días con iconos representativos de cada fase.

---

## Diapositiva 5: Cómo empecé — Análisis previo

**Título:** Fase 0: Investigación y análisis  
**Descripción:**
- Analicé múltiples proyectos POS existentes para encontrar el diferenciador
- Creé documentos de pre-análisis (`Pre-analisis/bootcamp-analysis.md`, `Deeper-analysis.md`)
- Definí la "constitution" del proyecto: reglas, alcance, no-goals
- Usé Wayfinder para research tickets (modelo de datos, NL→SQL, dashboard, seguridad)

**Comentarios para el presentador:** "Antes de escribir una línea de código, invertí tiempo entendiendo qué construir y por qué. Esto es lo que Kiro te permite hacer de forma estructurada."  
**Propuesta de diseño:** Captura de pantalla de los archivos de pre-análisis y el mapa de wayfinder.

---

## Diapositiva 6: Kiro como copiloto de desarrollo

**Título:** El poder de Kiro: Specs, Steering, Powers y Hooks  
**Descripción:**
- **Specs:** Workflow estructurado Requirements → Design → Tasks
- **Steering:** Reglas persistentes del proyecto (arquitectura, testing, seguridad, quality, convenciones)
- **Powers:** Long-Term Memory para persistencia entre sesiones, Context7 para docs actualizados
- **Hooks:** Automatización (lint on save, test after task, etc.)

**Comentarios para el presentador:** "Kiro no es solo autocomplete. Es un sistema que entiende tu proyecto, mantiene contexto, y trabaja con reglas que tú defines."  
**Propuesta de diseño:** Diagrama de 4 cuadrantes: Specs | Steering | Powers | Hooks, con iconos y 1-liner de cada uno.

---

## Diapositiva 7: Steering — Las reglas del proyecto

**Título:** Steering files: tu proyecto siempre consistente  
**Descripción:**
- `architecture.md` — Hexagonal ligera, reglas de capas, SOLID práctico
- `testing.md` — TDD obligatorio en auth/inventario/NL→SQL, cobertura mínima por capa
- `security.md` — Whitelist SQL, bcrypt PINs, read-only connections
- `quality.md` — golangci-lint, funciones <40 líneas, errores siempre manejados
- `project-conventions.md` — Stack, dependencias aprobadas, commits convencionales
- `design-patterns.md` — Patrones pragmáticos, anti-sobreingeniería

**Comentarios para el presentador:** "Estas reglas se cargan automáticamente en cada sesión. No tengo que repetir 'usa hexagonal' — Kiro ya lo sabe."  
**Propuesta de diseño:** Lista con iconos de candado/check, mostrando fragmentos de código de cada archivo.

---

## Diapositiva 8: Powers — Long-Term Memory

**Título:** Power: Long-Term Memory — Tu proyecto nunca olvida  
**Descripción:**
- Memoria local persistente entre sesiones
- 3 niveles: archivos recientes (Tier 1) → búsqueda en decisiones (Tier 2) → detalle completo (Tier 3)
- Guarda checkpoints, decisiones, hilos abiertos
- Recall barato: "Pick up where we left off"
- No necesita servicios externos, todo local en `ltm/`

**Comentarios para el presentador:** "Si cierro Kiro hoy y vuelvo mañana, no pierdo contexto. LTM recuerda qué archivos toqué, qué decisiones tomé, y qué queda pendiente."  
**Propuesta de diseño:** Diagrama de 3 tiers con flechas de escalamiento progresivo. Ejemplo de comando "Remember this project."

---

## Diapositiva 9: Sincronización y GitHub

**Título:** Sync + GitHub: Trabajo en cualquier dispositivo  
**Descripción:**
- Kiro Sync Files: workspace local ↔ Kiro cloud (app.dev)
- GitHub: issues, milestones, labels, projects para tracking
- 18 issues creadas automáticamente desde tasks.md
- Milestone "POS AI-First MVP" con dependency tracking

**Comentarios para el presentador:** "Las tareas del spec se convirtieron directamente en issues de GitHub con un comando. Cero trabajo manual de project management."  
**Propuesta de diseño:** Split screen: Kiro IDE a la izquierda, GitHub Issues/Project board a la derecha.

---

## Diapositiva 10: Arquitectura técnica

**Título:** Arquitectura: Hexagonal + AI Layer  
**Descripción:**
```
┌─────────────────────────────────┐
│  HTMX + Alpine.js + Tailwind   │
├─────────────────────────────────┤
│  Go HTTP (chi router)           │
├─────────────────────────────────┤
│  Application (use-cases)        │
├─────────────────────────────────┤
│  Domain (entities + ports)      │
├─────────────────────────────────┤
│  SQLite │ OpenRouter │ Config   │
└─────────────────────────────────┘
```

**Comentarios para el presentador:** "Las dependencias siempre apuntan hacia adentro. El dominio no sabe que existe HTTP ni SQLite."  
**Propuesta de diseño:** Diagrama de capas con colores por nivel y flechas de dependencia.

---

## Diapositiva 11: Seguridad NL→SQL (5 capas)

**Título:** 5 capas de seguridad para queries generadas por AI  
**Descripción:**
1. Prompt: instrucción al LLM de no generar DDL/DML
2. Validación Go: whitelist SELECT/WITH, reject keywords peligrosos
3. Conexión: SQLite read-only separada
4. Ejecución: timeout 5s, LIMIT 500
5. Auditoría: log de toda query generada

**Comentarios para el presentador:** "No confiamos en el LLM. Cada capa es un guardia independiente. Si una falla, las demás atrapan el problema."  
**Propuesta de diseño:** 5 escudos/capas apilados con nombres. Color rojo→verde de más riesgoso a más seguro.

---

## Diapositiva 12: Demo en vivo

**Título:** Demo: Pregúntale a tu POS  
**Descripción:**
- Login con PIN
- CRUD de productos
- Registrar una venta
- Preguntar: "¿Qué vendí hoy?"
- Dashboard con métricas actualizándose

**Comentarios para el presentador:** Preparar la demo con datos seeded. Tener backup en video por si falla la red.  
**Propuesta de diseño:** Pantalla completa del POS funcionando, sin slides — es la demo real.

---

## Diapositiva 13: Resultados y métricas

**Título:** Lo que logramos en 5 días  
**Descripción:**
- 18 tareas completadas
- ~108 sub-tareas implementadas
- 5 capas de seguridad NL→SQL
- Arquitectura hexagonal limpia
- Tests en dominio ≥90%
- Zero lint warnings
- Chat AI funcional con respuestas en español

**Comentarios para el presentador:** Números concretos. Mostrar la barra de progreso del milestone.  
**Propuesta de diseño:** Grid de métricas con iconos y números grandes. Estilo dashboard.

---

## Diapositiva 14: Lecciones aprendidas

**Título:** Lo que aprendí  
**Descripción:**
- Kiro + steering files = consistencia sin esfuerzo
- LTM power = no perder contexto entre sesiones
- Spec workflow = no empezar a codear sin plan
- NL→SQL requiere múltiples capas de defensa, no solo prompt engineering
- HTMX simplifica drasticamente el frontend para MVPs

**Comentarios para el presentador:** Sé honesto sobre qué fue difícil y qué sorprendió.  
**Propuesta de diseño:** Post-its o cards con cada lección, estilo retrospectiva.

---

## Diapositiva 15: Próximos pasos

**Título:** Hacia dónde va el proyecto  
**Descripción:**
- Migrar a AWS (RDS, Lambda o ECS, Bedrock)
- Agregar historial de conversación
- Soporte multi-sucursal
- Modo offline con sync
- Analytics más avanzados con dashboards personalizables

**Comentarios para el presentador:** "Este MVP demuestra el concepto. El siguiente paso es llevarlo a producción con AWS."  
**Propuesta de diseño:** Roadmap horizontal con hitos futuros.

---

## Diapositiva 16: Cierre

**Título:** Gracias — ¿Preguntas?  
**Descripción:** Links al repo, QR code, contacto.  
**Comentarios para el presentador:** Abre a preguntas. Ten preparadas respuestas para: "¿Por qué no usaste un ORM?", "¿Es seguro ejecutar SQL generado por AI?", "¿Por qué SQLite y no Postgres?"  
**Propuesta de diseño:** QR al repo + información de contacto sobre fondo limpio.
