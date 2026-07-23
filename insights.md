# Insights: Retrospectiva de Prompts — POS AI-First MVP

## Cómo usar este archivo

Este documento analiza cada prompt enviado durante el desarrollo del proyecto, comparando cómo se pidió vs. cómo se debería haber pedido para obtener mejores resultados más rápido. Útil para mejorar la comunicación con agentes AI y para la presentación del bootcamp.

---

## Prompt 1: "Update the tasks for pos-ai-first-mvp"

**Lo que pedí:** Actualizar las tareas del spec.  
**Lo que pasó:** El agente detectó que el archivo tenía errores de formato Kiro y ofreció reformatearlo.  
**Versión profesional:**
> "Reformatea `tasks.md` del spec `pos-ai-first-mvp` para que cumpla con el formato requerido por Kiro (diagnostics muestran 22 errores). Mantén todo el contenido existente, solo cambia la estructura."

**Lo que me faltó pedir:**
- Ser explícito sobre QUÉ tipo de actualización quería (formato vs. contenido vs. agregar tareas)
- Mencionar que los diagnostics ya mostraban errores

---

## Prompt 2: "si realiza ese reformateado"

**Lo que pedí:** Confirmación para proceder con el reformateo.  
**Lo que pasó:** El agente ejecutó correctamente.  
**Versión profesional:**
> "Procede con el reformateo. Asegúrate de que el archivo pase los diagnostics de Kiro con cero errores al terminar."

**Lo que me faltó pedir:**
- Pude haberlo incluido en el primer prompt como instrucción directa: "Reformatea y verifica que pase diagnostics"

---

## Prompt 3: "tambien haz commit y push y las tareas en las issues correspondientes, esto de las issues dime si es correcto ? antes de subirlas a issues"

**Lo que pedí:** Commit, push, y crear issues — pero primero validar el esquema de issues conmigo.  
**Lo que pasó:** El agente hizo commit/push en rama nueva y presentó la propuesta de issues antes de crearlas.  
**Versión profesional:**
> "Haz commit del cambio en una rama nueva y push a origin. Luego propón un esquema de GitHub Issues (una por task, con labels por día y categoría, milestone 'POS AI-First MVP', y checklists de sub-tareas). Muéstrame la propuesta antes de ejecutar."

**Lo que me faltó pedir:**
- Especificar el nombre de la rama
- Indicar si quería un PR automático
- Definir si quería labels predefinidos o dejar que el agente propusiera

**Flujo diferente:** Idealmente, commit + push + issues + PR se piden en un solo prompt bien estructurado para evitar ida y vuelta.

---

## Prompt 4: "me parece correcto adelante"

**Lo que pedí:** Aprobación para crear las issues.  
**Lo que pasó:** Se crearon 18 issues con milestone, labels, checklists, y dependencias.  
**Versión profesional:**
> "Aprobado. Crea las 18 issues. Usa el formato propuesto."

**Lo que me faltó pedir:**
- Nada — fue una aprobación limpia después de una propuesta clara. Este es el flujo correcto.

---

## Prompt 5: "Si crea el pr para mergear la rama a main. Añade un par de tareas mas..."

**Lo que pedí:** Múltiples cosas en un solo prompt:
1. Crear PR
2. Usar LTM power y explicar cómo funciona
3. Crear presentacion.md con slides estructuradas
4. Crear video.md con guión de 5 minutos
5. Configurar GitHub Projects
6. Decirme cómo pedirlo profesionalmente

**Lo que pasó:** El agente ejecutó todo secuencialmente y manejó el error de permisos de Projects.  
**Versión profesional:**
> "Ejecuta las siguientes tareas en orden:
> 1. Crea un PR de `chore/reformat-tasks-spec` → `main`
> 2. Añade Tasks 19 y 20 al spec:
>    - Task 19: Preparación de presentación y video
>    - Task 20: Configuración de GitHub Projects y documentación
> 3. Crea `presentacion.md` con formato de slides (número, título, descripción, comentarios para presentador, propuesta de diseño). Incluye contexto del bootcamp, uso de Kiro, steering, powers, hooks, y sincronización.
> 4. Crea `video.md` con guión de 5 minutos basado en la presentación.
> 5. Configura GitHub Project Board V2 con vista Kanban (Backlog/In Progress/Review/Done).
> 6. Explícame brevemente cómo funciona el LTM Power.
>
> Si algo requiere permisos adicionales, indícame el comando para habilitarlos."

**Lo que me faltó pedir:**
- Separar las tareas por prioridad (PR primero, luego archivos, luego GitHub Projects que requiere permisos)
- Ser más específico sobre el contenido de las primeras diapositivas
- Indicar idioma preferido para los archivos (español vs. inglés)

**Flujo diferente:** Prompts multi-tarea funcionan mejor con una lista numerada y priorizada. Si el agente puede ejecutar en paralelo, se beneficia de saber qué es independiente vs. secuencial.

---

## Prompt 6: "UTILIZANDO el power long term memory debes de Remember this project..."

**Lo que pedí:** Tres cosas:
1. Ejecutar LTM (bootstrap/recall/checkpoint)
2. Configurar GitHub Projects (con comando de permisos)
3. Crear insights.md

**Versión profesional:**
> "Tres acciones:
> 1. **LTM:** Ejecuta el Power Long-Term Memory — haz recall de la sesión actual y guarda un checkpoint con todo el progreso hasta ahora (spec reformateado, 20 issues, PR, presentación, video).
> 2. **GitHub Projects:** Necesito que corras `gh auth refresh -s project,read:project` para actualizar permisos, y luego configures un Project Board V2 Kanban con columnas Backlog/In Progress/Review/Done, campos custom (Day, Priority, Estimation), y vincules todas las issues del milestone.
> 3. **Insights:** Crea `insights.md` con retrospectiva de cada prompt de esta sesión: qué pedí, versión profesional, qué me faltó, y qué flujo hubiera sido mejor."

**Lo que me faltó pedir:**
- El `gh auth refresh` requiere interacción del usuario en el navegador — no se puede ejecutar automáticamente en un script. Debí preguntar primero si quería ejecutarlo interactivamente.

---

## Prompt 7: "realice la instruccion del paso 2"

**Lo que pedí:** Confirmé que ya ejecuté `gh auth refresh -s project,read:project`.  
**Lo que pasó:** El agente verificó los scopes, creó el Project Board V2, configuró columnas (Backlog/In Progress/Review/Done), campos custom (Day, Priority, Estimation), y vinculó las 20 issues.  
**Versión profesional:**
> "Listo, ejecuté `gh auth refresh` y el token ya tiene scope `project`. Procede con la creación del Project Board V2 y la vinculación de todas las issues."

**Lo que me faltó pedir:**
- Nada crítico — fue una confirmación necesaria para desbloquear la tarea. Flujo correcto.

---

## Prompt 8: "1. Tenemos archivo de presentacion, de video y de insights, actualizalo..."

**Lo que pedí:** Dos cosas:
1. Actualizar los tres archivos con el progreso actual
2. Crear un hook para auto-registrar prompts en insights.md

**Versión profesional:**
> "Dos acciones:
> 1. **Actualizar archivos:** Agrega a `presentacion.md` una diapositiva sobre GitHub Projects. Actualiza `video.md` con la sección de project management. Añade los prompts 7 y 8 a `insights.md` y actualiza las métricas.
> 2. **Hook:** Crea un hook `promptSubmit` que después de cada prompt añada automáticamente una entrada a `insights.md` con: el prompt original, versión profesional, y qué faltó. Formato consistente con las entradas existentes."

**Lo que me faltó pedir:**
- Especificar QUÉ actualización quería en presentacion.md y video.md (agregar slides? actualizar contenido existente?)
- Indicar si el hook debería ejecutarse antes o después de completar la tarea

**Flujo diferente:** Separar "crear hook" (setup) de "actualizar archivos" (contenido) en dos prompts sería más limpio, pero unirlos está bien si se numeran.

---

## Patrones observados y recomendaciones

### ✅ Lo que funcionó bien
- Pedir aprobación antes de ejecutar acciones masivas (issues)
- Dar contexto suficiente para que el agente proponga un esquema
- Permitir que el agente use rama nueva en lugar de main

### ⚠️ Áreas de mejora
1. **Combinar prompts relacionados:** En lugar de 3 mensajes para "reformatea → confirma → commit", uno solo: "Reformatea, verifica diagnostics, commit en rama nueva, push, y crea PR."
2. **Ser explícito sobre formato de output:** "Muéstrame una tabla con..." es más preciso que "dime si es correcto."
3. **Separar consultas de acciones:** "Explícame X" en un prompt separado de "Ejecuta Y". Mezclarlas genera respuestas más largas.
4. **Numerar y priorizar:** Cuando hay 4+ cosas que hacer, una lista numerada con prioridad ayuda al agente a manejar dependencias.
5. **Indicar qué es bloqueante:** "Si algo falla, reporta y continúa con lo demás" vs. "Si algo falla, para todo."

### 🎯 Template recomendado para prompts complejos

```
## Contexto
[Una línea sobre dónde estamos y qué acabamos de hacer]

## Acciones (en orden de prioridad)
1. [Acción más importante — blocking para las demás]
2. [Acción independiente]
3. [Acción que depende de 1]

## Restricciones
- [Regla: no push a main directo]
- [Regla: mostrar propuesta antes de ejecutar X]

## Output esperado
- [Qué quiero ver al final: links, archivos, confirmación]
```

---

## Métricas de la sesión

| Métrica | Valor |
|---------|-------|
| Total de prompts del usuario | 8 |
| Prompts de confirmación/aprobación | 3 (37%) |
| Prompts multi-acción | 3 |
| Errores encontrados por permisos | 1 (GitHub Projects scope) |
| Archivos creados | 5 (tasks reformateado, presentacion.md, video.md, insights.md, hook) |
| Issues creadas | 20 |
| PRs creados | 1 |
| Commits realizados | 4 |
| GitHub Project Board | ✅ Configurado con 4 columnas + 3 campos custom |
| LTM checkpoint | ✅ Guardado con 5 decisiones |
| Hooks activos | 1 (prompt-insights-log) |

---

## Prompt 9: "Update the tasks for pos-ai-first-mvp"

**Lo que pedí:** Actualizar las tareas del spec pos-ai-first-mvp.
**Lo que pasó:** El agente preguntó qué tipo de actualización quería (status, restructure, add, regenerate). Seleccioné "Update completion status". Escaneó todo el codebase, confirmó que Tasks 1-2 ya estaban correctamente marcadas como done y 3-20 como pending. Hizo dos correcciones menores de precisión: (1) Task 1.7 — el driver SQLite real es `modernc.org/sqlite` no `go-sqlite3`, y `x/crypto`/`scs` aún no están instalados; (2) Task 2.2 — aclaró que usa el driver pure-Go.
**Versión profesional:**
> "Revisa el codebase actual y actualiza el completion status de `tasks.md` en el spec `pos-ai-first-mvp`. Marca como done solo lo que realmente existe implementado, y corrige cualquier discrepancia entre lo que dicen las tareas y lo que hay en el código (dependencias instaladas, driver usado, etc.)."

**Lo que me faltó pedir:**
- Especificar desde el inicio que quería "actualizar el status de completitud basado en el código real"
- Pude haber mencionado que sabía que solo Tasks 1-2 estaban hechas para acelerar

**Flujo diferente:** Un prompt más directo como "sincroniza tasks.md con el estado actual del repo" hubiera evitado la pregunta de clarificación.

---

## Prompt 10: "antes de continuar en la task 3, pudiera continuar en la app o kiro web?"

**Lo que pedí:** Si puedo cambiar de entorno (desktop → web/mobile) para seguir trabajando.  
**Lo que pasó:** El agente confirmó que sí, todo está en main y Sync Files sincroniza.  
**Versión profesional:**
> "Todo está en main. ¿Puedo hacer Sync Files y continuar la Task 3 desde app.kiro.dev? ¿Hay alguna consideración de compatibilidad entre entornos?"

**Lo que me faltó pedir:**
- Preguntar si la sincronización es bidireccional (lo preguntó en el siguiente prompt)

---

## Prompt 11: "si le doy sync files sincroniza en ambos sentidos?"

**Lo que pedí:** Confirmar si Sync Files es bidireccional (web→desktop también).  
**Lo que pasó:** Se explicó el flujo de sincronización.  
**Versión profesional:**
> "¿Kiro Sync Files es bidireccional? Si trabajo en app.kiro.dev y luego vuelvo al desktop, ¿los cambios se descargan automáticamente al hacer Sync?"

**Lo que me faltó pedir:**
- Nada — pregunta directa y clara.

---

## Prompt 12: "Se realizaron modificaciones en el repositorio desde kiro dev..."

**Lo que pedí:** Actualizar el proyecto local (git pull) y continuar con las tareas pendientes.
**Lo que pasó:** Se hizo git pull (recibió 25 archivos nuevos con 1211 líneas), se analizó el código existente vs tasks.md, se sincronizó el estado de tareas marcando como completadas las que ya estaban implementadas (tasks 4, 6, 7, 8, 11, 13, 14, 15 parcialmente), y luego se ejecutaron tasks 3.1-3.7, 4.2-4.3, 5.1-5.6, 6.5-6.6, 7.3-7.5, 8.4-8.5, 9.1 usando subagentes. Progreso: 77/137 tareas (56%).
**Versión profesional:**
> "1. Haz `git pull` para sincronizar con el remoto. 2. Analiza qué código nuevo entró y actualiza `tasks.md` marcando como completadas las tareas ya implementadas. 3. Ejecuta todas las tareas pendientes en orden de dependencias usando el orchestrator. Si algo falla, continúa con las que no dependen de la fallida."

**Lo que me faltó pedir:**
- Indicar explícitamente que quería ejecución continua (tuvo que decir "continua" 3 veces)
- Especificar un límite de tareas o tiempo (para no gastar contexto innecesariamente)

**Flujo diferente:** Un solo prompt "git pull + sync tasks.md status + ejecuta todas las tareas ready hasta completar o agotar contexto" hubiera evitado las interrupciones.

---

| Métrica | Valor |
|---------|-------|
| Total de prompts del usuario | 12 |
| Prompts de confirmación/aprobación | 3 (25%) |
| Prompts multi-acción | 4 |
| Errores encontrados por permisos | 1 (GitHub Projects scope) |
| Archivos creados | 5+ (tasks reformateado, presentacion.md, video.md, insights.md, hook) + 20+ domain/test files |
| Issues creadas | 20 |
| PRs creados | 1 |
| Commits realizados | 4+ |
| GitHub Project Board | ✅ Configurado con 4 columnas + 3 campos custom |
| LTM checkpoint | ✅ Guardado con 5 decisiones |
| Hooks activos | 1 (prompt-insights-log) |
| Tasks completadas (spec) | 77/137 (56%) |
| Tests passing | ✅ All packages green |

---

## Prompt 13: "ok continua con las que si puedes y tambien habia una..."

**Lo que pedí:** Continuar ejecutando tareas automáticas + verificar issues 21/22 en GitHub + revisar presentación + analizar deploy AWS + commit/push/merge.
**Lo que pasó:** Se verificó que las issues #22 y #23 (Task 19 y 20) son correctas (no existen 21/22 como tareas extra). Se actualizó la presentación añadiendo 3 diapositivas nuevas (onboarding 3→1, multiplataforma Kiro, deploy AWS). Se creó `governance/aws-deployment-plan.md` con análisis completo (App Runner, RDS, Bedrock, costos ~$25-40/mes). Se resolvieron conflictos de merge con origin/main, se reescribió main.go con wiring completo de todos los handlers, y se mergeó PR #28.
**Versión profesional:**
> "1. Verifica las issues abiertas en GitHub y confirma si hay tareas 21-22. 2. Actualiza `presentacion.md` añadiendo: onboarding (3 proyectos → elegí 1), uso multiplataforma de Kiro (desktop/web/mobile), y plan de deploy AWS con costos. 3. Crea un documento técnico `governance/aws-deployment-plan.md` con arquitectura AWS propuesta. 4. Haz commit, push a la rama actual, y merge a main."

**Lo que me faltó pedir:**
- Especificar que "21 y 22" se refería a los numbers de issues en GitHub, no a tasks del spec
- Indicar si el merge debía ser squash o merge commit

**Flujo diferente:** Separar la verificación de issues (consulta) del trabajo de actualización (acción) para tener respuesta más rápida sobre el estado antes de pedir cambios.

---

| Métrica | Valor |
|---------|-------|
| Total de prompts del usuario | 13 |
| Prompts de confirmación/aprobación | 3 (23%) |
| Prompts multi-acción | 5 |
| Archivos creados esta sesión | 55+ (dominio, tests, adapters, handlers, templates, README, AWS plan) |
| Issues en GitHub | 20 closed + 3 open (Tasks 19, 20, Revisión) |
| PRs creados y mergeados | 2 (#28 este session) |
| Tasks completadas (spec) | 129/137 (94%) |
| Tests passing | ✅ All 12 packages green |
| Lint | ✅ Zero warnings |
| Build | ✅ Clean |

---

## Prompt 14: "respecto a amazon el tier gratuito no opera para esta..."

**Lo que pedí:** Si el free tier de AWS aplica para deployar esta app.
**Lo que pasó:** Se investigó el estado actual del AWS Free Tier (post julio 2025): nuevas cuentas reciben $200 en créditos + 30+ servicios always free. Lambda (1M req/mes gratis), S3 (5GB), CloudFront (1TB), CloudWatch — todos aplican. Con Lambda en vez de App Runner, la app podría correr gratis 6+ meses.
**Versión profesional:**
> "Analiza si los servicios AWS que necesita el POS (compute, DB, AI) están cubiertos por el Free Tier actual. Incluye: (1) qué servicios tienen always-free, (2) qué cubren los $200 en créditos, (3) cuánto duraría sin costo, y (4) qué alternativa de arquitectura minimiza costos post-free-tier."

**Lo que me faltó pedir:**
- Indicar si tengo cuenta AWS nueva o existente (los beneficios cambian)
- Preguntar por Lambda como alternativa a App Runner desde el inicio

**Flujo diferente:** Debí preguntar esto ANTES de crear el deployment plan, para que el documento incluyera la opción free-tier desde el principio.

---

## Prompt 15: "guardalo como pendiente porque tambien faltaria ver que..."

**Lo que pedí:** Guardar la investigación del free tier como pendiente y anotar que también hay que investigar qué Kiro Powers de AWS existen para hacer el deploy.
**Lo que pasó:** Se añadió una sección "PENDIENTE" al final de `aws-deployment-plan.md` con los items a investigar: Lambda vs App Runner, Powers de AWS en Kiro, decisiones pendientes (cuenta nueva vs existente, PostgreSQL vs SQLite+EFS, Bedrock vs OpenRouter).
**Versión profesional:**
> "Agrega una sección 'Pendientes' al plan de deploy con: (1) evaluar Lambda+free tier como alternativa, (2) investigar qué Kiro Powers de AWS están disponibles para automatizar el deploy (CDK, SAM, Bedrock), (3) listar las decisiones que necesito tomar antes de ejecutar."

**Lo que me faltó pedir:**
- Nada — es un registro de pendientes correcto para la próxima sesión

**Flujo diferente:** Esto es buen ejemplo de usar el agente como "libreta" — anotar pendientes para no olvidarlos.

---

## Prompt 16: "una pregunta mas en ese free tier de amazon en sus dife..."

**Lo que pedí:** Si con Kiro y sus Powers se puede automatizar todo el deploy a AWS dentro del free tier.
**Lo que pasó:** Se investigaron los MCP servers de AWS disponibles para Kiro: IaC MCP Server (CDK/CloudFormation), Serverless MCP Server (Lambda), Deployment Agent SOPs, Security Agent (Kiro Power), Lambda Durable Power. Se confirmó que sí es posible automatizar el deploy casi por completo con un prompt, combinando estos tools con el free tier ($0/mes).
**Versión profesional:**
> "¿Es posible automatizar el deploy completo a AWS (Lambda, RDS, Bedrock) usando Kiro Powers y MCP servers, manteniéndose dentro del free tier? Lista qué MCP servers necesitaría instalar y qué pasos seguirían siendo manuales."

**Lo que me faltó pedir:**
- Pedir que se instalen los MCP servers ahora para tenerlos listos
- Preguntar si necesito alguna configuración previa (AWS CLI, credentials)

**Flujo diferente:** Esta es una buena pregunta de "viabilidad" que vale la pena hacer ANTES de planificar el deploy, para saber qué es posible.

---

## Prompt 17: "si instalalos de una vez"

**Lo que pedí:** Instalar los MCP servers/Powers de AWS ahora.
**Lo que pasó:** Se abrió el panel de configuración de Powers. La instalación de Powers requiere interacción manual del usuario en el panel visual de Kiro — el agente no puede hacer click en "Install". Se proporcionó la lista de lo que buscar: AWS Documentation, AWS IaC, AWS Serverless, AWS Security Agent.
**Versión profesional:**
> "Abre el panel de Powers e instala los MCP servers de AWS (IaC, Serverless, Documentation, Security). Si requiere configuración adicional (credentials, CLI), indícame los pasos."

**Lo que me faltó pedir:**
- Nada — la limitación es que los Powers se instalan desde UI, no por comando

**Flujo diferente:** El agente podría haber intentado configurar los MCP servers directamente en `.kiro/settings.json` si existiera esa opción programática.

---

## Prompt 18: "antes de continuar ejecuta la aplicacion para ver su f..."

**Lo que pedí:** Ejecutar la aplicación para verificar que funciona.
**Lo que pasó:** Se encontraron 3 bugs al ejecutar: (1) `CREATE INDEX` sin `IF NOT EXISTS` fallaba en re-ejecución de migrations, (2) PIN hashes en seed eran SHA-256 en vez de bcrypt — bcrypt no matcheaba, (3) `entities.Role` type no registrado en gob para session encoding. Se corrigieron los 3 y se verificó el flujo completo: login (1234) → dashboard → metrics APIs → productos. Todo funcional.
**Versión profesional:**
> "Ejecuta `make run` y verifica el flujo completo: (1) server inicia sin errores, (2) login con PIN funciona, (3) dashboard carga métricas, (4) APIs responden. Si hay errores, corrígelos y re-verifica."

**Lo que me faltó pedir:**
- Debí pedir esto mucho antes (después de implementar, antes de hacer merge)
- Pedir que se corra `make seed` también para tener datos de demo más ricos

**Flujo diferente:** El testing E2E debería ser un paso obligatorio ANTES de merge a main, no después. Incluirlo como pre-merge check en el workflow.

---

## Prompt 19: "Mediante specs, algunos issues a corregir, o features..."

**Lo que pedí:** Crear un spec para 5 issues/features (logout, alta productos/clientes, refresh HTMX, config API key admin) + recomendación de Powers AWS a instalar.
**Lo que pasó:** Se creó el spec completo `ui-fixes-and-admin-config` usando Quick Plan workflow (clarify → requirements → design → tasks). 6 requerimientos, design con arquitectura hexagonal + CryptoService + 5 correctness properties, 25 sub-tareas en 5 waves paralelas. Se recomendaron 4 Powers de AWS a instalar (SAM, CDK/CloudFormation, DevOps Agent, Strands).
**Versión profesional:**
> "Crea un nuevo spec llamado 'ui-fixes-and-admin-config' con Quick Plan para: (1) logout button en sidebar, (2) botón 'Nuevo Producto' + CRUD clientes, (3) fix HTMX cache/refresh, (4) panel admin para API key con encriptación AES-GCM. Además, de los AWS Powers disponibles en el panel, recomiéndame cuáles instalar para el futuro deploy a Lambda."

**Lo que me faltó pedir:**
- Indicar si quería ejecutar las tareas inmediatamente después de crear el spec
- Separar los bugs (logout, refresh) de las features (clientes, config) en specs diferentes

**Flujo diferente:** Un spec por categoría (bugfixes vs features) sería más limpio para tracking en GitHub, pero para un MVP de 5 días un solo spec consolidado es más práctico.

---

## Prompt 20: "actualiza tambien la presentacion con todos estos camb..."

**Lo que pedí:** Dos cosas combinadas: (1) actualizar la presentación con los cambios implementados (client CRUD, admin config, UI fixes), (2) incluir los Powers necesarios para la migración a AWS.
**Lo que pasó:** Se actualizó `presentacion.md` añadiendo: diapositiva 14 (nuevas funcionalidades implementadas), diapositiva 19 (Powers para AWS con tabla de instalados vs. necesarios), diapositiva 20 (próximos pasos con timeline). Se actualizaron también las diapositivas 7 (Powers), 16 (Resultados), y 17 (Lecciones). La presentación pasó de 19 a 21 diapositivas. Esto se hizo EN PARALELO con la ejecución de las 22 tareas del spec ui-fixes-and-admin-config.
**Versión profesional:**
> "Actualiza `presentacion.md`: (1) Agrega slide con las features nuevas (Client CRUD, Admin Config AES-GCM, UI fixes). (2) Agrega slide con los Kiro Powers necesarios para AWS (tabla comparativa: instalados vs. pendientes). (3) Actualiza slides de resultados y lecciones con métricas actualizadas. Ejecuta esto en paralelo con las tareas del spec."

**Lo que me faltó pedir:**
- Especificar el número exacto de diapositiva donde insertarlas
- Indicar si quería reemplazar la diapositiva de AWS existente o complementarla

**Flujo diferente:** Pedir la actualización de presentación como tarea separada del "run all tasks" hubiera sido más claro, pero combinarlos aprovechó el paralelismo del agente.

---

## Prompt 21: "realiza commit, push y merge si es necesario"

**Lo que pedí:** Commit, push a rama nueva, y merge a main.
**Lo que pasó:** Se creó branch `feat/ui-fixes-and-admin-config`, commit con mensaje detallado (41 archivos, 2873 insertions), push, PR #29 creado con descripción completa, y merge inmediato con delete de branch. Fast-forward limpio.
**Versión profesional:**
> "Haz commit de todos los cambios en rama `feat/ui-fixes-and-admin-config`, push, crea PR con descripción de los cambios, y merge a main."

**Lo que me faltó pedir:**
- Nada — prompt corto y efectivo para una acción mecánica bien definida
- El agente infirió correctamente el nombre de la rama basado en el spec

**Flujo diferente:** Este es un buen ejemplo de prompt mínimo para operaciones git rutinarias. No necesita más contexto porque el agente ya sabe qué cambió.

---

| Métrica | Valor |
|---------|-------|
| Total de prompts del usuario | 21 |
| Prompts de confirmación/aprobación | 3 (14%) |
| Prompts multi-acción | 7 |
| Archivos creados esta sesión | 68+ (dominio, tests, adapters, handlers, templates, specs, governance) |
| Issues en GitHub | 20 closed + 3 open |
| PRs creados y mergeados | 3 (#28, #29) |
| Tasks completadas (ui-fixes spec) | 22/22 (100%) |
| Tasks completadas (pos-ai-first spec) | 129/137 (94%) |
| Tests passing | ✅ All packages green |
| Property tests | 3 (RequireRole 150iter, CryptoRoundTrip, MaskAPIKey) |
| Lint | ✅ Zero warnings |
| Build | ✅ Clean |
| Presentación | 21 diapositivas |

---
