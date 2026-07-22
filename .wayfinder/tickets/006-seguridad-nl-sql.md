---
title: "Seguridad del NL→SQL"
labels: wayfinder:grilling
blocking: []
---

## Question

Definir cómo evitar que el NL→SQL genere queries destructivas o peligrosas:

- ¿Read-only wrapper (solo SELECT)?
- ¿Whitelist de tablas y columnas permitidas?
- ¿Validación post-generación: AST parser que rechace DDL/DML?
- ¿Usuario con permisos restringidos en SQLite?
- ¿Rate limiting por usuario?
- ¿Qué pasa si el prompt del usuario intenta jailbreak ("ignora las instrucciones anteriores")?
- ¿Logging de todas las consultas generadas para auditoría?

Resolver con grilling.
