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
| Quiere máximo WOW factor y acepta riesgo | **Agente Telefónico** | Demo en vivo con llamada real |

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

### 1.3 Configurar Hooks para el equipo

```
Prompt a Kiro:
"Crea un hook que ejecute TypeScript type-check cada vez que edito un archivo .ts o .tsx"
```

Hooks recomendados:
- **Lint on save** → `fileEdited` en `*.ts,*.tsx` → `npm run lint --fix`
- **Type check** → `fileEdited` en `*.ts,*.tsx` → `npx tsc --noEmit`
- **Test runner** → `postTaskExecution` → `npm run test`

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

---

## Fase 3 — Cómo usar Kiro efectivamente en equipo

### Sesiones Spec (para features complejas)

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

```
Prompts de ejemplo:
- "Genera el componente ProductCard con nombre, precio, stock y botón de editar"
- "Conecta este formulario al endpoint POST /api/products"
- "Agrega loading skeleton mientras carga la lista de productos"
- "Refactoriza esta función para manejar el caso cuando Bedrock no responde"
```

### Patrón de trabajo en paralelo

```
Miembro A (rama feature/chat-ui):
  → Abre Kiro → Trabaja en componentes de chat
  → Usa steering para mantener consistencia

Miembro B (rama feature/bedrock-integration):
  → Abre Kiro → Trabaja en Lambda + Bedrock
  → Usa steering para mantener consistencia

Merge → main (al final de cada día)
```

---

## Fase 4 — Preparación del Demo (Día 5)

### Checklist pre-demo

- [ ] Deploy funcional en AWS (no localhost)
- [ ] 3-5 flujos demostración preparados y probados
- [ ] Datos de prueba cargados (mínimo 20 productos, 50 ventas)
- [ ] Fallbacks para cuando el LLM falle (respuestas predefinidas)
- [ ] Video backup por si falla la conectividad
- [ ] Slide con arquitectura AWS (diagrama)
- [ ] Métricas: "En 5 días, con Kiro generamos X líneas, Y componentes, Z endpoints"

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

## Fase 5 — Tips de productividad con Kiro

### Prompts de alta efectividad

| Necesidad | Prompt recomendado |
|---|---|
| Generar CRUD completo | "Genera CRUD completo para la entidad Product con campos: name, price, stock, category. Incluye API routes, DynamoDB operations, y componentes React" |
| Integrar servicio AWS | "Integra Amazon Bedrock usando el SDK v3. El modelo es Claude. El prompt del sistema es: 'Eres un asistente de punto de venta...'" |
| Debug | "Este endpoint retorna 500. Aquí está el error: [pegar error]. Revisa el handler y sugiere fix" |
| Optimizar | "Esta query a DynamoDB es lenta con muchos registros. Sugiere un GSI o patrón de acceso mejor" |
| UI rápida | "Genera una tabla responsive con sorting para mostrar ventas. Usa Tailwind. Incluye empty state y loading" |

### Archivos de contexto útiles

Crea estos archivos para que Kiro tenga contexto permanente:

- `.kiro/steering/api-patterns.md` → Patrones de API (auth, error handling, pagination)
- `.kiro/steering/aws-services.md` → Servicios AWS usados y cómo conectarlos
- `.kiro/steering/demo-requirements.md` → Lo que debe funcionar para el demo day

---

## Resumen Visual

```
Día 1: 🏗️  Setup + Estructura + Primeras entidades
Día 2: ⚙️  CRUD completo + UI funcional
Día 3: 🧠  Integración AI (Bedrock) + Feature principal
Día 4: 🔧  Polish + Error handling + Fallbacks
Día 5: 🚀  Deploy + Demo prep + Ensayo
```

---

## Decisión final recomendada

**POS AI-First** si quieres el balance perfecto entre viabilidad e impacto.  
**Code Review Agent** si tu equipo es técnicamente fuerte y quiere innovación máxima.

Ambos son ganadores potenciales. La diferencia la hace la ejecución, no la idea.
