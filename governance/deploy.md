# Deploy a AWS — Estado y Próximos Pasos

## Estado actual: CÓDIGO COMPLETO ✅

Todo el código de infraestructura está implementado y en la rama `feat/aws-lambda-deploy` (PR pendiente de merge).

### Lo que ya está hecho (código):

| Componente | Archivo | Status |
|------------|---------|--------|
| PostgreSQL Migrations | `migrations/postgres/001_init.sql` | ✅ |
| PostgreSQL Adapters (7) | `src/infrastructure/adapters/postgres_*.go` | ✅ |
| Bedrock Adapter | `src/infrastructure/adapters/bedrock_query_service.go` | ✅ |
| Secrets Manager Loader | `src/infrastructure/config/secrets.go` | ✅ |
| Bootstrap Dual-Mode | `internal/bootstrap/bootstrap.go` | ✅ |
| Config Validation | `internal/bootstrap/config.go` | ✅ |
| Structured Logging | `internal/bootstrap/logging.go` | ✅ |
| Panic Recovery | `internal/bootstrap/middleware.go` | ✅ |
| Health Endpoint | `internal/bootstrap/health.go` | ✅ |
| Lambda Entry Point | `cmd/lambda/main.go` | ✅ |
| Server Refactored | `cmd/server/main.go` (uses bootstrap) | ✅ |
| Dockerfile (ARM64) | `Dockerfile` | ✅ |
| SAM Template | `template.yaml` | ✅ |
| CI/CD Pipeline | `.github/workflows/deploy.yml` | ✅ |
| MetricsRepository Port | `src/domain/ports/metrics_repository.go` | ✅ |

### Lo que falta (infraestructura AWS — manual):

| # | Paso | Quién | Tiempo est. |
|---|------|-------|-------------|
| 1 | Crear VPC con subnets privadas + NAT Gateway | Kiro (SAM/CLI) | 10 min |
| 2 | Crear RDS PostgreSQL (db.t4g.micro, free tier) | Kiro (CLI) | 15 min |
| 3 | Ejecutar migraciones en RDS | Kiro (psql) | 2 min |
| 4 | Crear secretos en Secrets Manager (DB URL, session key, AI config) | Kiro (CLI) | 5 min |
| 5 | Crear ECR repository | Kiro (CLI) | 2 min |
| 6 | Build & push Docker image a ECR | Kiro (CLI) | 5 min |
| 7 | `sam deploy` con los parámetros | Kiro (CLI) | 5 min |
| 8 | Sync static/ a S3 | Kiro (CLI) | 2 min |
| 9 | Habilitar Bedrock Claude 3 Haiku | TÚ (consola) | 5 min |
| 10 | Health check + testing E2E | Kiro | 5 min |

### Decisiones tomadas:

- **Compute:** Lambda + API Gateway (free tier: 1M req/mes gratis)
- **DB:** RDS PostgreSQL db.t4g.micro (free tier: 12 meses gratis)
- **AI:** Amazon Bedrock Claude 3 Haiku (~$2/mes, cubierto por $200 créditos)
- **Static:** S3 + CloudFront (free tier)
- **Secrets:** AWS Secrets Manager
- **IaC:** AWS SAM (template.yaml)
- **CI/CD:** GitHub Actions → ECR → SAM deploy
- **Región:** us-east-1
- **Arquitectura:** ARM64 (20% más barato)

### Costo estimado: $0/mes (free tier activo)

- Lambda: 1M requests/mes GRATIS (always free)
- RDS: db.t4g.micro GRATIS 12 meses
- Bedrock: ~$2/mes cubierto por $200 créditos
- S3 + CloudFront: 5GB + 1TB GRATIS
- Secrets Manager: ~$1.60/mes (4 secrets × $0.40)
- **Total real: ~$2-4/mes** cubierto por créditos → **$0 efectivo**

### Para retomar mañana:

```bash
# 1. Merge el PR actual
gh pr create --title "feat: AWS Lambda deploy" --base main
gh pr merge --merge --delete-branch

# 2. Crear infraestructura AWS (Kiro ejecuta esto)
# VPC → RDS → Secrets → ECR → SAM deploy → S3 sync

# 3. Único paso manual tuyo:
# Ir a https://console.aws.amazon.com/bedrock/home?region=us-east-1#/modelaccess
# Habilitar Claude 3 Haiku
```

### Comando para retomar:

> "Retomemos el deploy a AWS. El código está en main, falta crear la infra (VPC, RDS, secrets, ECR, sam deploy). Ejecuta los pasos del 1 al 10 de deploy.md."

---

## Arquitectura Final

```
Browser → CloudFront (static/) → S3
       → API Gateway → Lambda (Go ARM64, 512MB)
                         ├── PostgreSQL (RDS t4g.micro)
                         ├── Bedrock (Claude 3 Haiku)
                         └── Secrets Manager
```

**Lo que NO cambió:** `src/domain/`, `src/application/`, `templates/` — zero modificaciones.
**Lo que cambió:** Solo adaptadores de infraestructura + bootstrap dual-mode.
