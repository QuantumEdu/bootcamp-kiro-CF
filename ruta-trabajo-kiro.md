# Ruta de Trabajo en Equipo con Kiro — Bootcamp Código Facilito

## Contexto
- **Tiempo:** 5 días de desarrollo
- **Herramientas:** Kiro (IDE con agente AI) + AWS
- **Evaluación:** Impacto tecnológico (30%) + Innovación (30%) + Software funcional (30%) + Uso AWS/Kiro (10%)
- **Top 2 propuestas:** POS AI-First y Code Review Agent (MCP)

---

## Fase 0 — Decisión (Día 0, antes de arrancar)

### Elige tu proyecto según tu perfil:

| Si tu equipo es... | Elige | Razón |
|---|---|---|
| Fullstack + cómodo con LLMs | **POS AI-First** | CRUD sólido + Bedrock = demo killer |
| Backend-heavy + familiarizado con GitHub APIs | **Code Review Agent** | MCP + Kiro nativo = innovación máxima |

> **Nota:** Existe una tercera opción avanzada (Agente Telefónico con llamada en vivo) que ofrece máximo WOW factor pero implica mayor riesgo. No se cubre en esta guía por su complejidad adicional. Si tu equipo tiene experiencia con APIs de telefonía (Twilio/Amazon Connect), investiga por tu cuenta.

---

## Fase 0.5 — Preparación del Entorno (Día 0)

Esta fase consolida toda la configuración previa necesaria antes de escribir código.

### 0.5.1 Configuración del repositorio en GitHub

1. Crear el repositorio en GitHub (público o privado según reglas del bootcamp)
2. Configurar branch protection en `main`:
   - Requerir pull request antes de merge
   - Requerir al menos 1 aprobación
   - Bloquear push directo a main
3. Definir estrategia de branching:
   - `main` → rama estable, solo merges vía PR
   - `feature/<nombre>` → ramas de trabajo individual
   - `fix/<nombre>` → correcciones rápidas
   - Convención de nombres: `feature/chat-ui`, `feature/bedrock-integration`, etc.

### 0.5.2 Configuración de AWS (credenciales y permisos)

1. Crear un IAM user o role dedicado para el proyecto:
   - Permisos mínimos: DynamoDB, Lambda, Bedrock, CloudFormation (o SAM)
   - No usar root account
2. Configurar credenciales locales:
   ```bash
   aws configure --profile bootcamp-project
   # Ingresar: Access Key ID, Secret Access Key, región (us-east-1 recomendado)
   ```
3. Verificar acceso a Bedrock **antes del Día 3**:
   ```bash
   aws bedrock list-foundation-models --profile bootcamp-project --region us-east-1
   ```
   - Si no tienes acceso, solicitar habilitación del modelo Claude en la consola de Bedrock
4. Verificar que DynamoDB funciona:
   ```bash
   aws dynamodb list-tables --profile bootcamp-project
   ```

### 0.5.3 Variables de entorno y secretos

Crear archivo `.env.local` en la raíz del proyecto (nunca subirlo a git):

```env
# AWS
AWS_REGION=us-east-1
AWS_PROFILE=bootcamp-project

# Bedrock
BEDROCK_MODEL_ID=anthropic.claude-3-sonnet-20240229-v1:0

# DynamoDB
DYNAMODB_TABLE_PRODUCTS=bootcamp-products
DYNAMODB_TABLE_SALES=bootcamp-sales

# App
NEXT_PUBLIC_APP_URL=http://localhost:3000
```

**Gestión de secretos en producción:**
- Usar AWS Secrets Manager para API keys y credenciales sensibles
- Nunca hardcodear secretos en el código
- Agregar `.env.local` al `.gitignore`
- Para deploy, configurar variables de entorno en Lambda/Amplify directamente

---

## Fase 1 — Setup del Workspace con Kiro (Día 1 - Mañana)

### 1.1 Crear el proyecto y configurar Kiro

```
Prompt a Kiro (sesión Spec):
"Crea un proyecto [Next.js/React] con la siguiente estructura:
- /src/app → páginas
- /src/components → componentes UI
- /src/lib → lógica de negocio y utils
- /src/api → rutas API (serverless-ready)
- /infra → CDK o SAM para despliegue AWS
Incluye: TypeScript, Tailwind CSS, ESLint"
```

### 1.1.1 Instalación de dependencias iniciales

Después de que Kiro genere el scaffold, ejecutar:

```bash
npm install
```

Verificar que el proyecto compila correctamente:

```bash
npx tsc --noEmit        # Verificar tipos
npm run build           # Verificar build completo
npm run dev             # Verificar que inicia en localhost
```

Si hay errores de compilación, resolverlos antes de continuar. Es crítico tener una base limpia desde el inicio.

### 1.2 Configurar Steering Files para el equipo

Crear `.kiro/steering/project-standards.md`:

```markdown
---
inclusion: always
---
# Estándares del Proyecto

## Stack
- Frontend: React/Next.js + Tailwind CSS + TypeScript
- Backend: AWS Lambda (Node.js/TypeScript)
- Base de datos: DynamoDB
- AI: Amazon Bedrock (Claude)
- Infraestructura: AWS CDK o SAM

## Convenciones
- Componentes: PascalCase, un archivo por componente
- APIs: kebab-case en rutas, camelCase en handlers
- Commits: conventional commits (feat:, fix:, chore:)
- Errores: siempre tipar con interfaces, nunca `any`

## Criterios de calidad
- Toda función pública debe tener JSDoc
- Manejo de errores explícito en cada API route
- Loading states y empty states en UI
```

Crear `.kiro/steering/api-patterns.md`:

```markdown
---
inclusion: always
---
# Patrones de API

## Autenticación
- Para el bootcamp: API keys simples o sin auth (MVP)
- Producción: Amazon Cognito

## Manejo de errores
- Siempre retornar { error: string, code: number }
- Loguear errores con contexto (requestId, timestamp)
- Nunca exponer stack traces al cliente

## Paginación
- Usar cursor-based pagination con DynamoDB LastEvaluatedKey
- Default: 20 items por página
- Parámetros: ?limit=20&cursor=<token>
```

Crear `.kiro/steering/aws-services.md`:

```markdown
---
inclusion: always
---
# Servicios AWS

## DynamoDB
- Usar single-table design cuando sea posible
- Definir GSIs según patrones de acceso
- Siempre usar SDK v3 (@aws-sdk/client-dynamodb)

## Lambda
- Runtime: Node.js 20.x
- Timeout: 30s para APIs normales, 60s para Bedrock
- Memory: 256MB mínimo para Bedrock

## Bedrock
- Modelo: Claude 3 Sonnet
- Usar InvokeModel API
- Implementar retry con exponential backoff
- Siempre tener respuestas de fallback
```

Crear `.kiro/steering/demo-requirements.md`:

```markdown
---
inclusion: always
---
# Requisitos del Demo

## Funcionalidades obligatorias para el demo day
- App desplegada en AWS (no localhost)
- Mínimo 3 flujos funcionales de punta a punta
- Integración con Bedrock funcionando
- Datos de prueba cargados

## UX del demo
- Máximo 7 minutos de presentación
- Mostrar arquitectura AWS (diagrama)
- Destacar uso de Kiro en el proceso
- Tener video backup por si falla conectividad
```

### 1.3 Configurar Hooks para el equipo

Configura hooks de automatización en Kiro para mantener la calidad del código. Los hooks permiten ejecutar acciones automáticas ante ciertos eventos.

> **Importante:** Los nombres exactos de los triggers de hooks pueden variar según la versión de Kiro. Consulta siempre la [documentación oficial de Kiro](https://kiro.dev/docs) para verificar los triggers disponibles y su sintaxis actual.

Hooks recomendados a configurar:
- **Lint on save** → Al editar archivos `.ts` o `.tsx` → `npm run lint --fix`
- **Type check** → Al editar archivos `.ts` o `.tsx` → `npx tsc --noEmit`
- **Test runner** → Al completar una tarea → `npm run test`

```
Prompt a Kiro:
"Configura hooks de automatización según la documentación actual de Kiro:
1. Un hook que ejecute lint y autofix al guardar archivos TypeScript
2. Un hook que ejecute type-check al guardar archivos TypeScript
3. Un hook que ejecute los tests al completar una tarea"
```

---

## Fase 2 — Desarrollo por Día (División de trabajo)

### Para POS AI-First:

| Día | Miembro A (Frontend) | Miembro B (Backend + AI) | Kiro ayuda con |
|---|---|---|---|
| **1** | Setup proyecto + Layout base + Navegación | DynamoDB tables + Lambda CRUD productos | Scaffolding completo, CDK setup |
| **2** | UI Productos (lista, crear, editar) + UI Ventas | API ventas + lógica de inventario | Generación de componentes UI, validaciones |
| **3** | Chat UI (interfaz conversacional) | Bedrock integration: NL → query → respuesta | Integración Bedrock, prompt engineering |
| **4** | Dashboard métricas + gráficas | Queries predefinidas de respaldo + error handling | Charts, optimización de prompts |
| **5** | Polish UI + responsive + demo prep | Deploy AWS + smoke testing + fallbacks | Testing E2E, deployment scripts |

### Para Code Review Agent (MCP):

| Día | Miembro A (Agent Core) | Miembro B (Integration + UI) | Kiro ayuda con |
|---|---|---|---|
| **1** | Setup Lambda + GitHub webhook receiver | Dashboard React básico + autenticación | Estructura MCP, GitHub API scaffolding |
| **2** | Parser de diffs + extracción de contexto | UI: lista de PRs + vista de review | Parsing logic, componentes UI |
| **3** | Bedrock: análisis de código + generación de comentarios | GitHub API: postear comentarios en PRs | Prompt engineering, API integration |
| **4** | Refinamiento de prompts + manejo de repos grandes | Configuración por repo + reglas custom | Optimización, edge cases |
| **5** | Deploy + testing con PRs reales | Demo preparation + UI polish | Deployment, testing |

### Adaptación para equipos de 3 o más miembros

Si tu equipo tiene más de 2 personas, adapta la división de trabajo así:

**Equipo de 3 miembros (POS AI-First):**
| Rol | Responsabilidades |
|---|---|
| Miembro A (Frontend) | UI, componentes, responsive, UX |
| Miembro B (Backend) | APIs, DynamoDB, Lambda, infraestructura |
| Miembro C (AI + Integración) | Bedrock, prompts, chat, dashboard de métricas |

**Equipo de 3 miembros (Code Review Agent):**
| Rol | Responsabilidades |
|---|---|
| Miembro A (Agent Core) | Parser de diffs, lógica de análisis, Bedrock |
| Miembro B (GitHub Integration) | Webhooks, GitHub API, posting comments |
| Miembro C (UI + DevOps) | Dashboard, configuración por repo, deploy |

**Principios para dividir trabajo con N miembros:**
- Cada miembro debe tener una rama `feature/` independiente
- Minimizar dependencias entre miembros (interfaces definidas desde Día 1)
- Un miembro puede asumir el rol de "integrador" que se encarga de merges y resolución de conflictos
- Rotar el rol de integrador si hay conflictos frecuentes

### Estrategia de merge y resolución de conflictos

**Flujo recomendado (al final de cada día):**

1. Cada miembro hace commit y push a su rama `feature/`
2. Crear Pull Request hacia `main`
3. Revisar el PR (mínimo 1 aprobación del compañero)
4. Hacer **squash merge** para mantener el historial limpio:
   ```bash
   # En GitHub: seleccionar "Squash and merge" al aprobar el PR
   ```

**Resolución de conflictos:**
```bash
# Antes de crear el PR, actualizar tu rama con main:
git fetch origin
git rebase origin/main

# Si hay conflictos:
# 1. Resolver manualmente los archivos marcados
# 2. git add <archivos-resueltos>
# 3. git rebase --continue
# 4. Push con force (solo en tu feature branch):
git push --force-with-lease origin feature/<tu-rama>
```

**Recomendaciones:**
- Preferir **rebase** sobre merge para mantener historial lineal
- Usar **squash merge** en PRs para agrupar commits de una feature
- Si un conflicto es complejo, resolverlo en pair programming
- Nunca hacer force push a `main`

---

## Fase 3 — Cómo usar Kiro efectivamente en equipo

> **Importante:** En Kiro, los modos Spec y Vibe son modos de sesión. Para cambiar entre Spec mode y Vibe mode, necesitas iniciar una nueva sesión. No se pueden usar de forma intercambiable dentro de la misma sesión. Planifica qué modo usarás antes de comenzar cada sesión de trabajo.

### Sesiones Spec (para features complejas)

Usa Spec mode cuando necesites planificar una feature completa con requirements y diseño técnico. Inicia una nueva sesión en modo Spec.

```
Prompt:
"Quiero implementar la funcionalidad de chat conversacional para el POS.
El usuario escribe una pregunta en lenguaje natural como '¿qué vendí ayer?'
y el sistema responde con datos reales de DynamoDB vía Bedrock."
```

Kiro generará:
1. Requirements detallados
2. Diseño técnico (architecture + data flow)
3. Tasks ordenadas con dependencias

Luego cada miembro toma tasks y las ejecuta en modo Autopilot.

### Sesiones Vibe (para iteración rápida)

Para iteración rápida y tareas puntuales, inicia una nueva sesión en modo Vibe.

```
Prompts de ejemplo:
- "Genera el componente ProductCard con nombre, precio, stock y botón de editar"
- "Conecta este formulario al endpoint POST /api/products"
- "Agrega loading skeleton mientras carga la lista de productos"
- "Refactoriza esta función para manejar el caso cuando Bedrock no responde"
```

### Cuándo usar cada modo

| Situación | Modo | Razón |
|---|---|---|
| Feature nueva completa | **Spec** (nueva sesión) | Necesitas planificación y diseño |
| Componente UI individual | **Vibe** (nueva sesión) | Iteración rápida, sin planificación pesada |
| Integración con servicio AWS | **Spec** (nueva sesión) | Requiere diseño de arquitectura |
| Bug fix o refactor pequeño | **Vibe** (nueva sesión) | Cambio puntual y rápido |
| Tareas ya planificadas en Spec | **Autopilot** (misma sesión Spec) | Ejecutar tasks generadas |

### Patrón de trabajo en paralelo

```
Miembro A (rama feature/chat-ui):
  → Abre Kiro → Inicia sesión Spec o Vibe según la tarea
  → Trabaja en componentes de chat
  → Usa steering para mantener consistencia

Miembro B (rama feature/bedrock-integration):
  → Abre Kiro → Inicia sesión Spec o Vibe según la tarea
  → Trabaja en Lambda + Bedrock
  → Usa steering para mantener consistencia

Merge → main (al final de cada día vía PR con squash merge)
```

---

## Fase 4 — Preparación del Demo (Día 5)

### Checklist pre-demo

| Tarea | Responsable | Estado |
|---|---|---|
| Deploy funcional en AWS (no localhost) | Miembro B | [ ] |
| 3-5 flujos demostración preparados y probados | Miembro A + B | [ ] |
| Datos de prueba cargados (mínimo 20 productos, 50 ventas) | Miembro B | [ ] |
| Fallbacks para cuando el LLM falle (respuestas predefinidas) | Miembro B | [ ] |
| Video backup por si falla la conectividad | Miembro A | [ ] |
| Slide con arquitectura AWS (diagrama) | Miembro A | [ ] |
| Métricas: "En 5 días, con Kiro generamos X líneas, Y componentes, Z endpoints" | Miembro A + B | [ ] |
| Ensayo completo del demo (cronometrado) | Miembro A + B | [ ] |

> **Nota para equipos de 3+:** Distribuir las tareas de demo entre todos los miembros. Un miembro puede encargarse exclusivamente de la presentación y el slide de arquitectura.

### Script del demo (POS AI-First)

1. **Abrir la app** → mostrar dashboard con métricas
2. **Crear producto** → flujo CRUD completo
3. **Registrar venta** → actualización de inventario
4. **Chat:** "¿Cuánto vendí hoy?" → respuesta instantánea
5. **Chat:** "¿Qué producto tiene más rotación?" → respuesta con datos
6. **Chat:** "Muéstrame ventas de la última semana" → gráfica
7. **Cerrar** con slide de arquitectura + roadmap

### Script del demo (Code Review Agent)

1. **Abrir dashboard** → lista de repos conectados
2. **Crear un PR** con código deliberadamente malo
3. **Esperar 30 segundos** → mostrar comentarios automáticos
4. **Mostrar análisis** → bugs, seguridad, mejores prácticas
5. **Cerrar** con arquitectura MCP + potencial de escalamiento

---

## Fase 5 — Testing

### Estrategia de testing

Define qué tipos de tests escribir según el día y la prioridad:

| Tipo | Herramienta recomendada | Cuándo escribirlos | Prioridad |
|---|---|---|---|
| **Unit tests** | Jest + Testing Library | Días 2-4, para lógica de negocio | Alta |
| **Integration tests** | Jest + supertest | Días 3-4, para APIs y Bedrock | Media |
| **E2E tests** | Playwright | Día 4-5, flujos críticos del demo | Media |
| **Smoke tests** | Script bash / Jest | Día 5, post-deploy | Alta |

### Tests mínimos recomendados

```
Prompt a Kiro (sesión Vibe):
"Genera tests unitarios para:
1. La función que procesa la respuesta de Bedrock
2. El handler de la API de productos (CRUD)
3. La lógica de cálculo de inventario
Usa Jest con TypeScript."
```

**Para POS AI-First:**
- [ ] Test CRUD de productos (crear, leer, actualizar, eliminar)
- [ ] Test de lógica de ventas (descuento de inventario)
- [ ] Test de integración con Bedrock (mock del API)
- [ ] Test E2E del flujo de chat (Playwright)
- [ ] Smoke test post-deploy (health check endpoint)

**Para Code Review Agent:**
- [ ] Test del parser de diffs
- [ ] Test de generación de comentarios (mock Bedrock)
- [ ] Test de integración con GitHub API (mock webhook)
- [ ] Test E2E de flujo completo (PR → análisis → comentario)
- [ ] Smoke test post-deploy (webhook endpoint responde)

### Configuración de Jest

```bash
npm install -D jest @types/jest ts-jest @testing-library/react @testing-library/jest-dom
```

```json
// jest.config.js básico
{
  "preset": "ts-jest",
  "testEnvironment": "jsdom",
  "moduleNameMapper": {
    "^@/(.*)$": "<rootDir>/src/$1"
  }
}
```

---

## Fase 6 — Tips de productividad con Kiro

### Prompts de alta efectividad

| Necesidad | Prompt recomendado |
|---|---|
| Generar CRUD completo | "Genera CRUD completo para la entidad Product con campos: name, price, stock, category. Incluye API routes, DynamoDB operations, y componentes React" |
| Integrar servicio AWS | "Integra Amazon Bedrock usando el SDK v3. El modelo es Claude. El prompt del sistema es: 'Eres un asistente de punto de venta...'" |
| Debug | "Este endpoint retorna 500. Aquí está el error: [pegar error]. Revisa el handler y sugiere fix" |
| Optimizar | "Esta query a DynamoDB es lenta con muchos registros. Sugiere un GSI o patrón de acceso mejor" |
| UI rápida | "Genera una tabla responsive con sorting para mostrar ventas. Usa Tailwind. Incluye empty state y loading" |
| Tests | "Genera tests unitarios con Jest para esta función. Incluye casos edge: input vacío, error de red, timeout" |

---

## Resumen Visual

```
Día 0: 📋  Decisión de proyecto + Preparación del entorno (GitHub, AWS, .env)
Día 1: 🏗️  Setup + Estructura + Steering files + Primeras entidades
Día 2: ⚙️  CRUD completo + UI funcional + Primeros tests
Día 3: 🧠  Integración AI (Bedrock) + Feature principal + Tests integración
Día 4: 🔧  Polish + Error handling + Fallbacks + Tests E2E
Día 5: 🚀  Deploy + Testing final + Demo prep + Ensayo
```

---

## Decisión final recomendada

**POS AI-First** si quieres el balance perfecto entre viabilidad e impacto.  
**Code Review Agent** si tu equipo es técnicamente fuerte y quiere innovación máxima.

Ambos son ganadores potenciales. La diferencia la hace la ejecución, no la idea.
