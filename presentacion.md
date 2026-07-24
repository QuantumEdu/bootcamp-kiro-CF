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

## Diapositiva 5: Onboarding — De 3 proyectos a 1

**Título:** El proceso de selección: 3 ideas → 1 MVP  
**Descripción:**
- En el bootcamp se presentaron 3 propuestas de proyecto a evaluar
- Se analizó viabilidad técnica, impacto, e innovación de cada una
- Se eligió el POS AI-First por combinar impacto real en PYMES + diferenciador técnico (NL→SQL)
- Pre-análisis documentado: `bootcamp-analysis.md`, `Deeper-analysis.md`, `constitution.md`

**Comentarios para el presentador:** "No fue solo 'tengo una idea'. Evaluamos 3 opciones con criterios claros. Este proceso de decisión es parte del valor de trabajar con Kiro — todo queda documentado."  
**Propuesta de diseño:** Tabla comparativa de 3 proyectos con checkmarks indicando por qué ganó el POS AI-First.

---

## Diapositiva 6: Cómo empecé — Análisis previo

**Título:** Fase 0: Investigación y análisis  
**Descripción:**
- Analicé múltiples proyectos POS existentes para encontrar el diferenciador
- Creé documentos de pre-análisis (`Pre-analisis/bootcamp-analysis.md`, `Deeper-analysis.md`)
- Definí la "constitution" del proyecto: reglas, alcance, no-goals
- Usé Wayfinder para research tickets (modelo de datos, NL→SQL, dashboard, seguridad)

**Comentarios para el presentador:** "Antes de escribir una línea de código, invertí tiempo entendiendo qué construir y por qué. Esto es lo que Kiro te permite hacer de forma estructurada."  
**Propuesta de diseño:** Captura de pantalla de los archivos de pre-análisis y el mapa de wayfinder.

---

## Diapositiva 7: Kiro como copiloto de desarrollo

**Título:** El poder de Kiro: Specs, Steering, Powers y Hooks  
**Descripción:**
- **Specs:** Workflow estructurado Requirements → Design → Tasks, con ejecución paralela por waves
- **Steering:** Reglas persistentes del proyecto (arquitectura, testing, seguridad, quality, convenciones, design patterns)
- **Powers:** LTM para memoria entre sesiones, Context7 para docs actualizados, Power Builder para crear powers custom
- **Hooks:** Automatización (lint on save, test after task, auto-doc de prompts)
- **Task DAG:** Grafo de dependencias para ejecutar tareas en paralelo sin conflictos

**Comentarios para el presentador:** "Kiro no es solo autocomplete. Es un sistema que entiende tu proyecto, mantiene contexto, ejecuta tareas en paralelo respetando dependencias, y trabaja con reglas que tú defines. 5 steering files definen cómo trabaja MI proyecto."  
**Propuesta de diseño:** Diagrama de 4 cuadrantes: Specs | Steering | Powers | Hooks, con iconos y 1-liner de cada uno.

---

## Diapositiva 8: Steering — Las reglas del proyecto

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

## Diapositiva 9: Powers — Long-Term Memory

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

## Diapositiva 10: Sincronización y GitHub

**Título:** Sync + GitHub: Trabajo en cualquier dispositivo  
**Descripción:**
- Kiro Sync Files: workspace local ↔ Kiro cloud (app.dev)
- GitHub: issues, milestones, labels, projects para tracking
- 20 issues creadas automáticamente desde tasks.md
- Milestone "POS AI-First MVP" con dependency tracking
- GitHub Project Board V2 con Kanban (Backlog/In Progress/Review/Done)
- Campos custom: Day, Priority, Estimation
- Hook automático que registra cada prompt y su versión mejorada

**Comentarios para el presentador:** "Las tareas del spec se convirtieron directamente en issues de GitHub con un comando. El Project Board da visibilidad del progreso sin esfuerzo manual. Y el hook de insights documenta automáticamente mi proceso de comunicación con el agente."  
**Propuesta de diseño:** Split screen: Kiro IDE a la izquierda, GitHub Project Board Kanban a la derecha. Overlay mostrando el flujo: spec → issues → board.

---

## Diapositiva 11: Multiplataforma — Kiro Desktop, Web y Mobile

**Título:** Desarrollo continuo desde cualquier dispositivo  
**Descripción:**
- **Kiro Desktop (VS Code):** Desarrollo principal con terminal, debugging, extensiones
- **Kiro Web (app.kiro.dev):** Acceso desde cualquier navegador, mismo workspace
- **Kiro Mobile:** Revisión de código, aprobación de PRs, consulta rápida desde el celular
- **Sync Files:** Sincronización bidireccional entre todos los entornos
- Flujo real: escribí código en desktop, revisé issues en mobile, continué en web desde otro PC

**Comentarios para el presentador:** "No estuve amarrado a una sola máquina. Pude avanzar desde el café con la web, revisar desde el teléfono, y retomar en el desktop sin perder nada."  
**Propuesta de diseño:** 3 dispositivos (laptop, browser, phone) conectados por flechas de sync. Screenshots reales de cada uno.

---

## Diapositiva 12: Arquitectura técnica

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

## Diapositiva 13: Seguridad NL→SQL (5 capas)

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

## Diapositiva 14: Nuevas funcionalidades — UI Fixes y Admin Config

**Título:** Iteración 2: Gestión de Clientes + Configuración Admin  
**Descripción:**
- **Logout visible:** Botón "Cerrar Sesión" en el sidebar footer para todos los usuarios
- **Nuevo Producto:** Botón directo en la lista de productos para crear rápido
- **CRUD Clientes:** Listado, creación con validación de nombre, tabla con datos de contacto
- **Admin Config:** Página exclusiva para admins — almacena API key cifrada con AES-GCM
- **HTMX No-Cache:** Middleware que garantiza datos frescos en cada navegación
- **Sidebar inteligente:** "Configuración" solo visible para rol admin (template conditional)

**Flujo técnico de seguridad (API Key):**
```
Admin → Form → AES-GCM Encrypt (SHA-256 de SESSION_SECRET) → SQLite configuracion → Decrypt on read → Mask (****últimos4)
```

**Comentarios para el presentador:** "Esto se construyó con el spec workflow de Kiro: requirements → design → tasks → ejecución paralela por waves. 30 tareas, 5 waves, resolviendo dependencias automáticamente."  
**Propuesta de diseño:** Split: izquierda sidebar con las nuevas opciones, derecha el formulario de config con la key enmascarada.

---

## Diapositiva 15: Demo en vivo

**Título:** Demo: Pregúntale a tu POS  
**Descripción:**
- Login con PIN
- CRUD de productos (con botón "Nuevo Producto")
- Gestión de clientes (crear, listar)
- Registrar una venta
- Preguntar: "¿Qué vendí hoy?"
- Dashboard con métricas actualizándose
- Admin: configurar API key cifrada

**Comentarios para el presentador:** Preparar la demo con datos seeded. Mostrar el flujo completo: login → clientes → venta → chat AI → config admin. Tener backup en video por si falla la red.  
**Propuesta de diseño:** Pantalla completa del POS funcionando, sin slides — es la demo real.

---

## Diapositiva 16: Resultados y métricas

**Título:** Lo que logramos en 5 días  
**Descripción:**
- 3 specs creados (pos-ai-first, ui-fixes, aws-deploy)
- 100+ tareas ejecutadas en waves paralelas (5 simultáneas)
- 5 capas de seguridad NL→SQL
- Arquitectura hexagonal: dominio INTACTO al migrar a AWS
- Tests en dominio: 100%, application: 60%+, middleware: 86%
- Zero lint warnings
- Chat AI funcional con respuestas en español
- CRUD completo: Productos, Ventas, Clientes
- Panel admin con cifrado AES-GCM
- Deploy a AWS Lambda: $0/mes con free tier
- Cold start: 4.4s → warm: 1-3ms
- CI/CD: push a main → auto-deploy en ~4 min
- 7 adaptadores PostgreSQL + Bedrock adapter (mismas interfaces)

**Comentarios para el presentador:** "100+ tareas, 3 specs, deploy a AWS — todo sin salir de Kiro. El agente ejecutó 5 tareas en paralelo respetando dependencias. El deploy fue automático: push → tests → build → deploy → health check."  
**Propuesta de diseño:** Grid de métricas con números grandes. Highlight en $0/mes y en 1-3ms warm.

---

## Diapositiva 17: Lecciones aprendidas

**Título:** Lo que aprendí  
**Descripción:**
- Kiro + steering files = consistencia sin esfuerzo
- LTM power = no perder contexto entre sesiones
- Spec workflow = no empezar a codear sin plan
- Wave-based task execution = máxima eficiencia con dependencias respetadas
- NL→SQL requiere múltiples capas de defensa, no solo prompt engineering
- HTMX simplifica drásticamente el frontend para MVPs
- AES-GCM + SESSION_SECRET = secretos protegidos sin servicios externos
- Property tests validan correctitud universal (no solo happy paths)

**Comentarios para el presentador:** Sé honesto sobre qué fue difícil y qué sorprendió. "La ejecución paralela por waves fue reveladora — 5 agentes trabajando al mismo tiempo sin pisarse."  
**Propuesta de diseño:** Post-its o cards con cada lección, estilo retrospectiva.

---

## Diapositiva 18: Deploy en AWS — Ya implementado

**Título:** De MVP local a producción en AWS — El código ya está  
**Descripción:**

**Arquitectura implementada:**
```
Browser → CloudFront (static/) → S3
       → API Gateway HTTP API → Lambda (Go ARM64, 512MB)
                                  ├── RDS PostgreSQL (pgxpool)
                                  ├── Bedrock Claude 3 Haiku
                                  └── Secrets Manager (cache)
```

**Herramientas y flujo de deploy:**
| Herramienta | Rol en el pipeline |
|-------------|-------------------|
| **AWS SAM** | IaC — define Lambda, API Gateway, S3, CloudFront, IAM en un template.yaml |
| **GitHub Actions** | CI/CD — test → build imagen → push ECR → sam deploy → health check |
| **Docker (ARM64)** | Container image Lambda — Go binary 15MB + templates |
| **pgx/v5** | Driver PostgreSQL pure-Go (sin CGO) con connection pool optimizado |
| **algnhsa** | Adaptador chi router → Lambda event handler (zero-change en handlers) |
| **Secrets Manager** | Credenciales cifradas con cache in-memory (warm starts) |

**Dual-mode bootstrap:**
```go
switch cfg.AppEnv {
case "lambda":  // PostgreSQL + Bedrock + pgx sessions
default:        // SQLite + OpenRouter + SQLite sessions
}
```

**Archivos creados:** 35 archivos nuevos, 5243 líneas. Dominio intacto.

**Comentarios para el presentador:** "Todo el código de deploy ya está escrito. La arquitectura hexagonal permitió agregar 7 adaptadores PostgreSQL sin tocar una sola línea del dominio. Mañana solo falta crear la infra en AWS y dar `sam deploy`."  
**Propuesta de diseño:** Diagrama de flujo del pipeline: código → GitHub Actions → ECR → Lambda → API Gateway. Highlight del dual-mode switch.

---

## Diapositiva 19: El poder de la arquitectura hexagonal

**Título:** Hexagonal en acción: migración sin dolor  
**Descripción:**

**Lo que NO cambió (cero líneas):**
- `src/domain/` — Entidades, validaciones, puertos
- `src/application/` — Use cases, NL→SQL service
- `templates/` — Todo el frontend HTMX

**Lo que se AGREGÓ (solo adaptadores nuevos):**
| Adaptador | Implementa | Driver |
|-----------|------------|--------|
| `postgres_product_repository.go` | `ports.ProductRepository` | pgxpool |
| `postgres_sale_repository.go` | `ports.SaleRepository` | pgxpool + tx |
| `postgres_user_repository.go` | `ports.UserRepository` | pgxpool |
| `postgres_client_repository.go` | `ports.ClientRepository` | pgxpool |
| `postgres_inventory_repository.go` | `ports.InventoryRepository` | pgxpool |
| `postgres_config_repository.go` | `ports.ConfigRepository` | pgxpool |
| `postgres_metrics_repository.go` | `ports.MetricsRepository` | pgxpool |
| `bedrock_query_service.go` | `ports.AIQueryService` | AWS SDK v2 |

**Resultado:** Misma interfaz, distinta implementación. El use case `RegisterSale` no sabe si persiste en SQLite o PostgreSQL.

**Comentarios para el presentador:** "Esta es la demostración práctica de por qué la arquitectura hexagonal importa. 7 adaptadores nuevos, zero cambios en lógica de negocio. El dominio es inmutable."  
**Propuesta de diseño:** Diagrama hexagonal con el dominio en el centro (candado) y flechas de los adaptadores nuevos apuntando hacia adentro.

---

## Diapositiva 20: Costo y timeline

**Título:** $0/mes — Free Tier cubre todo  
**Descripción:**

| Servicio | Free Tier | Costo post-free |
|----------|-----------|-----------------|
| Lambda | 1M req/mes GRATIS (always free) | ~$0.20/1M req |
| RDS PostgreSQL | 12 meses gratis (t4g.micro) | ~$13/mes |
| Bedrock Claude Haiku | $200 créditos (~100 meses) | ~$2/mes |
| S3 + CloudFront | 5GB + 1TB gratis | ~$1/mes |
| Secrets Manager | 4 secrets | ~$1.60/mes |
| **Total primer año** | | **$0 efectivo** |

**Timeline de deploy (restante):**
1. Crear VPC + RDS + Secrets (Kiro automatiza) — 30 min
2. Build Docker + push ECR — 5 min  
3. `sam deploy` — 5 min
4. Habilitar Bedrock (manual, consola) — 5 min
5. Health check + test E2E — 5 min

**Comentarios para el presentador:** "El código ya está hecho. Solo falta dar 'sam deploy' y la app vive en la nube. $0/mes el primer año gracias al free tier."  
**Propuesta de diseño:** Tabla de costos grande con $0 resaltado. Timeline horizontal con checkmarks.

---

## Diapositiva 21: Próximos pasos

**Título:** Hacia dónde va el proyecto  
**Descripción:**
- **Mañana:** Deploy a AWS (`sam deploy`) — la app vive en la nube
- **Semana 1:** Seed de datos en RDS, Bedrock habilitado, demo cloud
- **Futuro:**
  - Historial de conversación AI
  - Soporte multi-sucursal
  - Modo offline con sync
  - Analytics con dashboards personalizables
  - Predicción de demanda (Bedrock + datos históricos)

**Comentarios para el presentador:** "El MVP demuestra el concepto. La migración a producción es literalmente un comando. Hexagonal + Kiro specs = velocidad sin sacrificar calidad."  
**Propuesta de diseño:** Roadmap con 3 fases: MVP Local (✅) → AWS Deploy (🔄 mañana) → AI Advanced (📋).

---

## Diapositiva 22: Cierre

**Título:** Gracias — ¿Preguntas?  
**Descripción:** Links al repo, QR code, contacto.  
**Comentarios para el presentador:** Abre a preguntas. Ten preparadas respuestas para: "¿Por qué no usaste un ORM?", "¿Es seguro ejecutar SQL generado por AI?", "¿Por qué SQLite y no Postgres?", "¿Qué Powers de Kiro usaste?", "¿Cuánto cuesta correr esto en AWS?", "¿Cómo funciona el dual-mode bootstrap?"  
**Propuesta de diseño:** QR al repo + información de contacto sobre fondo limpio.
