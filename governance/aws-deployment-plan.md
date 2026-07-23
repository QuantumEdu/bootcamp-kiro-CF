# Plan de Deployment en AWS — POS AI-First

## Resumen

Migrar el MVP local (Go + SQLite + OpenRouter) a una arquitectura cloud-native en AWS manteniendo la misma arquitectura hexagonal. Gracias a la separación de capas, solo cambian los **adaptadores** de infraestructura.

---

## Arquitectura Propuesta

```
                    ┌────────────────────────┐
                    │   CloudFront (CDN)     │
                    │   + S3 (static assets) │
                    └──────────┬─────────────┘
                               │
                    ┌──────────▼─────────────┐
                    │   App Runner / ECS     │
                    │   (Go binary container)│
                    │   Auto-scaling 0→N     │
                    └──┬──────────────────┬──┘
                       │                  │
          ┌────────────▼───┐    ┌─────────▼──────────┐
          │  RDS PostgreSQL │    │  Amazon Bedrock    │
          │  (t4g.micro)    │    │  Claude 3 Haiku   │
          │  Multi-AZ opt.  │    │  (NL→SQL)         │
          └────────────────┘    └────────────────────┘
                       │
          ┌────────────▼───┐
          │  Secrets Manager│
          │  (API keys,    │
          │   DB creds)    │
          └────────────────┘
```

---

## Componentes y Servicios

### 1. Compute — AWS App Runner (recomendado para MVP)

| Aspecto | Detalle |
|---------|---------|
| Servicio | App Runner (o ECS Fargate para más control) |
| Imagen | Go binary en Alpine container (~15MB) |
| Scaling | 0 → 10 instancias, basado en requests |
| Memoria | 512MB (Go es muy eficiente) |
| Health check | GET /health |

**¿Por qué App Runner?** Zero infrastructure management, deploy desde ECR o GitHub, auto-scaling incluido, HTTPS gratis.

**Alternativa:** ECS Fargate si necesitas más control (VPC placement, sidecars, task definitions).

### 2. Base de datos — Amazon RDS PostgreSQL

| Aspecto | Detalle |
|---------|---------|
| Engine | PostgreSQL 16 |
| Instance | db.t4g.micro (2 vCPU, 1GB RAM) — free tier eligible |
| Storage | 20GB gp3 |
| Backups | Automated daily, 7 days retention |
| Multi-AZ | No (MVP), Sí (producción) |

**Migración de SQLite:**
- Las queries son 95% compatibles (modernc.org/sqlite usa SQL estándar)
- Cambiar driver: `github.com/lib/pq` o `github.com/jackc/pgx/v5`
- Crear migración DDL adaptada (INTEGER → SERIAL, TEXT datetime → TIMESTAMPTZ)
- El adaptador SQLite se reemplaza por un adaptador PostgreSQL (misma interfaz)

### 3. AI — Amazon Bedrock (reemplaza OpenRouter)

| Aspecto | Detalle |
|---------|---------|
| Modelo | Claude 3 Haiku (rápido, barato, bueno en SQL) |
| Pricing | ~$0.25/1M input tokens, ~$1.25/1M output tokens |
| Latencia | <2s típico para queries NL→SQL |
| Región | us-east-1 (mayor disponibilidad de modelos) |

**Cambios necesarios:**
- Nuevo adaptador `BedrockQueryService` implementando `ports.AIQueryService`
- Usar AWS SDK Go v2 (`github.com/aws/aws-sdk-go-v2/service/bedrockruntime`)
- Mismo system prompt, misma validación de respuesta
- El adapter de OpenRouter se mantiene como fallback

### 4. Secrets — AWS Secrets Manager

| Secret | Uso |
|--------|-----|
| `pos/db-connection` | PostgreSQL connection string |
| `pos/bedrock-config` | Región, modelo, parámetros |
| `pos/session-secret` | Cookie encryption key |
| `pos/openrouter-key` | Fallback API key (opcional) |

### 5. CDN + Static Assets — CloudFront + S3

| Aspecto | Detalle |
|---------|---------|
| S3 bucket | Static files (JS, CSS, imágenes) |
| CloudFront | Distribución global, HTTPS, cache |
| Origin | App Runner como origin para HTML dinámico |

### 6. Monitoring — CloudWatch + X-Ray

- **CloudWatch Logs:** Todos los logs del container (structured JSON)
- **CloudWatch Metrics:** Request count, latency, errors, DB connections
- **X-Ray:** Tracing del pipeline NL→SQL (request → Bedrock → DB → response)
- **Alarms:** Error rate >5%, latency p99 >5s, Bedrock failures

---

## Cambios en el Código

### Lo que NO cambia (gracias a hexagonal):
- `src/domain/` — Entidades, value objects, ports → 0 cambios
- `src/application/` — Use cases, NL→SQL service logic → 0 cambios
- `templates/` — Frontend completo → 0 cambios
- `static/` → Se mueve a S3 pero el contenido no cambia

### Lo que cambia (solo adaptadores):
1. **Nuevo:** `src/infrastructure/adapters/postgres_product_repository.go`
2. **Nuevo:** `src/infrastructure/adapters/postgres_sale_repository.go`
3. **Nuevo:** `src/infrastructure/adapters/postgres_user_repository.go`
4. **Nuevo:** `src/infrastructure/adapters/postgres_inventory_repository.go`
5. **Nuevo:** `src/infrastructure/adapters/bedrock_query_service.go`
6. **Modificar:** `src/infrastructure/database/connection.go` → PostgreSQL driver
7. **Modificar:** `src/infrastructure/config/config.go` → AWS Secrets Manager
8. **Nuevo:** `Dockerfile` para build del container
9. **Nuevo:** `.github/workflows/deploy.yml` para CI/CD

---

## Estimación de Costos (mensual)

| Servicio | Tier | Costo estimado |
|----------|------|----------------|
| App Runner | 0.5 vCPU, 1GB, min 1 instance | ~$7/mes |
| RDS PostgreSQL | db.t4g.micro | ~$13/mes (o free tier 12 meses) |
| Bedrock (Claude Haiku) | ~1000 queries/mes | ~$2/mes |
| Secrets Manager | 4 secrets | ~$2/mes |
| CloudFront + S3 | Bajo tráfico | ~$1/mes |
| CloudWatch | Básico | Incluido |
| **Total** | | **~$25-40/mes** |

Para una PYME con una tienda, es un costo operativo razonable.

---

## CI/CD Pipeline

```yaml
# .github/workflows/deploy.yml
name: Deploy to AWS
on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - run: go test ./... -cover
      - run: golangci-lint run

  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: aws-actions/configure-aws-credentials@v4
      - uses: aws-actions/amazon-ecr-login@v2
      - run: |
          docker build -t pos-ai-first .
          docker tag pos-ai-first:latest $ECR_REGISTRY/pos-ai-first:latest
          docker push $ECR_REGISTRY/pos-ai-first:latest
      # App Runner auto-deploys from ECR
```

---

## Dockerfile

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /pos cmd/server/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /pos /pos
COPY templates/ /app/templates/
COPY static/ /app/static/
WORKDIR /app
EXPOSE 8080
CMD ["/pos"]
```

---

## Timeline de Migración (estimado)

| Día | Tarea |
|-----|-------|
| 1 | Crear Dockerfile + probar local, setup ECR |
| 2 | Crear RDS, migrar schema, adaptar connection.go |
| 3 | Crear adaptadores PostgreSQL (reemplazar SQLite) |
| 4 | Integrar Bedrock, crear BedrockQueryService |
| 5 | Setup App Runner + CI/CD + CloudFront |
| 6 | Testing E2E en AWS, monitoring, alarmas |
| 7 | Documentar, cleanup, tag release |

---

## Decisión: ¿Cuándo migrar?

El MVP funciona perfecto con SQLite local para demo y desarrollo. La migración a AWS tiene sentido cuando:
- Se necesite acceso multi-usuario remoto (no solo local)
- Se requiera persistencia de datos en la nube
- Se quiera eliminar la dependencia de OpenRouter por Bedrock (más rápido, más barato, en tu VPC)
- Se necesite alta disponibilidad

La arquitectura hexagonal garantiza que la migración es de **adaptadores solamente** — la lógica de negocio no se toca.
