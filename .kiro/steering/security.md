---
inclusion: always
---
# Security — Reglas de Seguridad

## NL→SQL (riesgo crítico del proyecto)

1. NUNCA ejecutar SQL generado por LLM directamente sin validación.
2. Whitelist de operaciones: solo SELECT. Rechazar INSERT, UPDATE, DELETE, DROP, ALTER, TRUNCATE.
3. Whitelist de tablas: el LLM solo puede consultar tablas explícitamente permitidas.
4. Parametrizar valores siempre que sea posible (no concatenar strings en queries).
5. Timeout en queries generadas: máximo 5 segundos de ejecución.
6. Log toda query generada antes de ejecutar (auditoría).
7. Si la query falla validación → responder con mensaje amigable, NO exponer el error SQL.

## Validación de entrada

- Todo input del usuario se valida ANTES de llegar al dominio.
- Usar schemas de validación en `src/domain/value-objects/` o `src/schemas/`.
- Nunca confiar en datos del cliente: validar tipo, rango, longitud, formato.
- Sanitizar HTML/scripts en campos de texto libre.

## Autenticación

- PIN se almacena hasheado (bcrypt o argon2), nunca en texto plano.
- Sesiones con expiración (configurable, default 8 horas para turno POS).
- Rate limiting en intentos de PIN: máximo 5 intentos, lockout de 5 minutos.
- No exponer si el PIN es incorrecto vs usuario no existe (mensaje genérico).

## Secretos y configuración

- Variables sensibles (API keys, DB paths) en variables de entorno, nunca hardcodeadas.
- `.env` en `.gitignore` siempre.
- En producción: usar AWS Secrets Manager o SSM Parameter Store.

## Headers HTTP

- CORS restrictivo: solo orígenes conocidos.
- Content-Type validation en todos los endpoints.
- Rate limiting global en API (configurable).

## Dependencias

- Usar versiones exactas en `go.mod` (no rangos).
- Ejecutar `go vet` y `staticcheck` como parte del CI.
- No agregar dependencias sin justificación documentada.

## Principio general

> Si dudas entre "confiar en el input" y "validar de más", siempre valida de más.
> El costo de una validación extra es O(1). El costo de un SQL injection es catastrófico.
