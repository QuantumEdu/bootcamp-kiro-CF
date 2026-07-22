---
title: "OpenRouter + NL→SQL prompt design"
labels: wayfinder:research
blocking: []
status: completed
completed_at: "2026-07-22"
implemented_in: "feat/tasks-3-4-5 (PR #27)"
---

## Question

Investigar el mejor approach para el NL→SQL vía OpenRouter sobre SQLite:

- ¿Qué modelo de OpenRouter recomendado para SQLite? (Claude, GPT, Mistral, DeepSeek)
- System prompt para generar SQL seguro y correcto para SQLite
- Cómo parsear la respuesta del modelo (esperar solo SQL, o explicación + SQL)
- Few-shot examples: incluir las 5 consultas demostrables como ejemplos en el prompt
- Manejo de errores: SQL mal generado, timeout, sintaxis inválida
- Costos estimados por consulta
- ¿Cache de consultas frecuentes?
- Alternativa: ¿usar un modelo más chico y barato si la query es simple?
