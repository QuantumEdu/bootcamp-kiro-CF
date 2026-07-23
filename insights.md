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
