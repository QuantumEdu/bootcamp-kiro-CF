# Guión de Video: POS AI-First MVP (5 minutos)

## Bootcamp Kiro × Código Facilito — Hackathon 2026

---

## Estructura temporal

| Sección | Duración | Acumulado |
|---------|----------|-----------|
| Intro + Problema | 0:30 | 0:30 |
| Solución + Demo conceptual | 0:45 | 1:15 |
| Cómo lo construí (Kiro) | 1:00 | 2:15 |
| Demo en vivo | 1:30 | 3:45 |
| Arquitectura + Seguridad | 0:45 | 4:30 |
| Cierre + Próximos pasos | 0:30 | 5:00 |

---

## Sección 1: Intro + Problema (0:00 – 0:30)

**Visual:** Pantalla con título del proyecto, luego cut a una persona abriendo hojas de cálculo frustrada.

**Narración:**
> "Hola, soy [nombre]. Imaginen al dueño de una taquería que quiere saber cuánto vendió hoy. Actualmente abre Excel, filtra por fecha, suma columnas... o le pide a alguien. ¿Y si pudiera simplemente preguntar, como en WhatsApp, '¿qué vendí hoy?' y recibir la respuesta al instante?"

**Notas de producción:** Transición rápida, energética. Máximo 2 tomas.

---

## Sección 2: Solución + Demo conceptual (0:30 – 1:15)

**Visual:** Screencast del chat bar del POS. Se escribe "¿Qué producto se vendió más esta semana?" y aparece la respuesta.

**Narración:**
> "Construí un POS que habla. El usuario escribe una pregunta en español, el sistema genera una consulta SQL segura usando AI, la ejecuta contra sus datos reales, y devuelve la respuesta formateada. Todo en menos de 5 segundos."
>
> "Pero no es solo chat. Es un POS completo: productos, ventas, inventario, dashboard con métricas en tiempo real, y el chat como feature diferenciador."

**Notas de producción:** Mostrar el flujo completo con overlay de las 5 capas de seguridad como badges.

---

## Sección 3: Cómo lo construí — Kiro (1:15 – 2:15)

**Visual:** Pantalla de Kiro IDE mostrando specs, steering, y el task list.

**Narración:**
> "Lo construí en 5 días usando Kiro. Pero no fue 'open IDE y empezar a codear'. Seguí un proceso estructurado:"
>
> "Primero, análisis. Investigué proyectos similares, definí el alcance, identifiqué riesgos."
>
> "Segundo, configuré Kiro con Steering files — reglas de arquitectura, testing, seguridad, calidad y convenciones que se aplican automáticamente en cada sesión."
>
> "Tercero, usé el spec workflow: requirements, design, tasks. 20 tareas con dependencias, convertidas automáticamente en GitHub Issues y vinculadas a un Project Board Kanban."
>
> "Y cuarto, el Power de Long-Term Memory — memoria local que persiste entre sesiones. Cierro Kiro hoy, vuelvo mañana, y recuerda todo: qué hice, qué decisiones tomé, qué queda pendiente."
>
> "Bonus: un hook que documenta automáticamente cada prompt que envío y cómo debería haberlo pedido de forma más profesional. Meta-aprendizaje en tiempo real."

**Notas de producción:** Mostrar brevemente cada cosa mientras se menciona. Speed up en la navegación de archivos. Highlight en el steering y el LTM.

---

## Sección 4: Demo en vivo (2:15 – 3:45)

**Visual:** Screencast del POS funcionando completamente.

**Narración:**
> "Vamos a la demo. Primero, login con PIN — autenticación simple para un contexto POS real."

*[Muestra login con PIN]*

> "Dashboard: ventas de hoy, productos más vendidos, alertas de stock bajo. Se actualiza solo cada 30 segundos."

*[Muestra dashboard con métricas]*

> "Registro de venta: selecciono productos, agrego al carrito, completo. El inventario se actualiza automáticamente."

*[Muestra flujo de venta]*

> "Y ahora, la estrella: '¿Cuántas ventas hubo esta semana?'"

*[Escribe en chat, espera respuesta]*

> "Respuesta en 3 segundos. SQL generado, validado, ejecutado en read-only con timeout. 5 capas de seguridad entre el LLM y mi base de datos."

**Notas de producción:** Pregrabar la demo como backup. Si es en vivo, tener datos seeded y queries probadas. La parte del chat es el clímax del video.

---

## Sección 5: Arquitectura + Seguridad (3:45 – 4:30)

**Visual:** Diagrama de arquitectura hexagonal + diagrama de 5 capas de seguridad.

**Narración:**
> "La arquitectura es hexagonal: el dominio no importa frameworks, los use-cases orquestan, y la infraestructura implementa los adaptadores."
>
> "Stack: Go, chi router, SQLite con WAL mode, HTMX para el frontend server-driven, y OpenRouter API para el AI."
>
> "Seguridad NL→SQL en 5 capas: prompt, validación Go, conexión read-only, timeout con LIMIT, y auditoría. No confiamos en el LLM — cada capa es independiente."

**Notas de producción:** Animación simple del diagrama, mostrando las capas una por una.

---

## Sección 6: Cierre (4:30 – 5:00)

**Visual:** Resumen de métricas + roadmap + pantalla final.

**Narración:**
> "En 5 días: 20 tareas, arquitectura limpia, tests en dominio al 90%, zero lint warnings, un chat AI funcional con respuestas en español, y un Project Board que se mantiene solo."
>
> "Próximos pasos: migrar a AWS, agregar historial de conversación, y soporte multi-sucursal."
>
> "Esto es lo que pasa cuando combinas buenas herramientas con un proceso estructurado. Gracias."

**Notas de producción:** Terminar con el QR al repo y transición a pantalla de cierre.

---

## Checklist de producción del video

- [ ] Grabar screencast de demo con datos seeded
- [ ] Preparar backup de demo en caso de fallo de red
- [ ] Grabar narración por separado (mejor audio)
- [ ] Editar con overlays de texto para puntos clave
- [ ] Verificar que el video dura exactamente ≤ 5:00
- [ ] Exportar en 1080p mínimo
- [ ] Subir copia a Google Drive como backup
- [ ] Probar reproducción antes de la presentación
