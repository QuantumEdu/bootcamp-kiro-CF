# POS AI-First — Punto de Venta Inteligente

> "¿Qué vendí hoy?" — Pregunta en español, obtén respuestas de tus datos al instante.

[![Deploy to AWS](https://github.com/QuantumEdu/bootcamp-kiro-CF/actions/workflows/deploy.yml/badge.svg)](https://github.com/QuantumEdu/bootcamp-kiro-CF/actions/workflows/deploy.yml)

---

## 🎯 Resumen

**POS AI-First** es un punto de venta con agente conversacional integrado, desarrollado como MVP funcional en 5 días para el **Bootcamp Kiro × Código Facilito × AWS | Hackathon 2026**.

El sistema permite a dueños de pequeños negocios (taquerías, tiendas de abarrotes, etc.) gestionar productos, ventas e inventario, y además **hacer preguntas en lenguaje natural** sobre sus datos — como enviar un mensaje por WhatsApp.

**Presentado por:** Gabriel Magallón  
**Fecha:** 26 de julio de 2026  
**Ubicación:** Michoacán, México

---

## 🌐 Demo en vivo

| Recurso | URL |
|---------|-----|
| **App (AWS Lambda)** | https://zz637vr6cd.execute-api.us-east-1.amazonaws.com |
| **Health Check** | https://zz637vr6cd.execute-api.us-east-1.amazonaws.com/health |

**Credenciales de acceso:**
- Admin: PIN `1234`
- Cajero: PIN `1235`

> ⚠️ La app puede estar desactivada fuera del horario de demo. Para reactivar: ver sección de Deploy.

---

## 🚀 ¿Qué problema resuelve?

Los dueños de pequeños negocios:
- No tienen tiempo de revisar reportes
- Dependen de Excel desordenados para registrar ventas
- Las preguntas son simples: "¿Qué vendí hoy?", "¿Qué se está agotando?"
- La respuesta debería ser tan fácil como preguntar en WhatsApp

**POS AI-First** convierte preguntas en español → SQL seguro → respuestas formateadas, todo en menos de 5 segundos.

---

## ✨ Features

### Core POS
- 🛒 CRUD de productos con categorías y SKU
- 💰 Registro de ventas con múltiples métodos de pago
- 👥 Gestión de clientes
- 📊 Dashboard con métricas en tiempo real (HTMX auto-refresh)
- 🔐 Autenticación por PIN (bcrypt, lockout por intentos)

### AI Chat (NL→SQL)
- 💬 Preguntas en español sobre ventas, productos, inventario
- 🧠 Generación de SQL via OpenRouter (DeepSeek V4 Flash)
- 🛡️ 5 capas de seguridad (prompt, validación Go, read-only, timeout, auditoría)
- 📋 Respuestas formateadas con explicación

### Admin
- ⚙️ Panel de configuración para API keys (AES-GCM cifrado)
- 🔒 Rutas admin protegidas por rol (RequireRole middleware)

---

## 🏗️ Arquitectura

```
┌─────────────────────────────────────────────────────┐
│  Browser (HTMX + Alpine.js + Tailwind CSS)          │
├─────────────────────────────────────────────────────┤
│  API Gateway HTTP API → Lambda (Go ARM64, 512MB)    │
├─────────────────────────────────────────────────────┤
│  Application Layer (Use Cases, NL→SQL Service)      │
├─────────────────────────────────────────────────────┤
│  Domain (Entities, Ports, Value Objects)            │
├─────────────────────────────────────────────────────┤
│  Infrastructure Adapters                            │
│  ├── PostgreSQL (pgxpool) — AWS RDS                 │
│  ├── SQLite (modernc.org) — Local dev              │
│  ├── OpenRouter (DeepSeek V4 Flash) — NL→SQL       │
│  └── Bedrock (Claude 3 Haiku) — AWS production     │
└─────────────────────────────────────────────────────┘
```

**Hexagonal Architecture:** El dominio no importa frameworks. Los adaptadores son intercambiables. La migración de SQLite→PostgreSQL fue **7 archivos nuevos, zero cambios en dominio**.

### Dual-Mode Bootstrap

```go
switch cfg.AppEnv {
case "lambda":  // PostgreSQL + Bedrock + pgx sessions
default:        // SQLite + OpenRouter + SQLite sessions
}
```

---

## 🔒 Seguridad NL→SQL (5 capas)

| Capa | Defensa |
|------|---------|
| 1. Prompt | Instrucción al LLM: "solo genera SELECT" |
| 2. Validación Go | Whitelist SELECT/WITH, reject DDL/DML keywords |
| 3. Conexión | SQLite/PostgreSQL read-only separada |
| 4. Ejecución | Timeout 5s, LIMIT 500 registros |
| 5. Auditoría | Log de toda query generada antes de ejecutar |

---

## 🛠️ Stack Técnico

| Componente | Tecnología |
|------------|-----------|
| Backend | Go 1.26 (chi/v5 router) |
| Frontend | HTMX + Alpine.js + Tailwind CSS (CDN) |
| DB Local | SQLite (modernc.org/sqlite, pure Go) |
| DB Producción | PostgreSQL 16 (Amazon RDS, pgx/v5) |
| AI | OpenRouter → DeepSeek V4 Flash ($0.09/1M tokens) |
| Compute | AWS Lambda (ARM64, container image) |
| IaC | AWS SAM (template.yaml) |
| CI/CD | GitHub Actions (test → build → deploy → health check) |
| Session | alexedwards/scs (pgxstore para AWS) |
| Crypto | AES-GCM (API keys cifradas en reposo) |

---

## 📦 Estructura del proyecto

```
pos-ai-first/
├── cmd/
│   ├── server/main.go       # Entry point local (net/http)
│   ├── lambda/main.go       # Entry point AWS (algnhsa)
│   ├── seed/main.go         # Seed SQLite local
│   └── seedpg/main.go       # Seed PostgreSQL RDS
├── internal/bootstrap/       # Dual-mode router builder
├── src/
│   ├── domain/              # Entities, Ports, Value Objects (0 deps)
│   ├── application/         # Use Cases, NL-SQL Service
│   └── infrastructure/      # Adapters, HTTP, Config
├── templates/               # HTML templates (HTMX)
├── static/                  # JS (Alpine components)
├── migrations/
│   ├── 001_init.sql         # SQLite schema
│   └── postgres/001_init.sql # PostgreSQL schema
├── Dockerfile               # Lambda container (ARM64)
├── template.yaml            # AWS SAM IaC
└── .github/workflows/       # CI/CD pipeline
```

---

## 🚀 Cómo ejecutar

### Local (desarrollo)

```bash
# Clonar y configurar
git clone https://github.com/QuantumEdu/bootcamp-kiro-CF.git
cd bootcamp-kiro-CF
cp .env.example .env
# Editar .env con tu OPENROUTER_API_KEY y SESSION_SECRET

# Ejecutar
go run cmd/server/main.go
# → http://localhost:8080
```

### AWS (producción)

```bash
# Prerequisitos: AWS CLI configurado, SAM CLI instalado
sam deploy --guided

# O via CI/CD: push a main dispara deploy automático
git push origin main
```

### Reactivar la app (si está desactivada)

```bash
aws apigatewayv2 create-stage --api-id zz637vr6cd --stage-name '$default' --auto-deploy --region us-east-1
```

---

## 💰 Costos

| Servicio | Free Tier | Post-free |
|----------|-----------|-----------|
| Lambda | 1M req/mes GRATIS (always free) | ~$0.20/1M |
| RDS PostgreSQL | 12 meses gratis (t4g.micro) | ~$13/mes |
| OpenRouter (DeepSeek) | Pay-per-use ($0.09/1M tokens) | ~$2/mes |
| S3 + CloudFront | 5GB + 1TB gratis | ~$1/mes |
| **Total año 1** | | **~$0-4/mes** |

---

## 🏆 Desarrollado con Kiro

Este proyecto demuestra el poder del **desarrollo asistido por agentes**:

- **Specs workflow:** Requirements → Design → Tasks con ejecución paralela por waves
- **Steering files:** 6 archivos de reglas persistentes (arquitectura, testing, seguridad, quality, convenciones, patrones)
- **Powers:** Long-Term Memory para persistencia entre sesiones, Context7 para docs
- **Hooks:** Auto-documentación de prompts, lint on save
- **Task orchestration:** 5 subagentes ejecutando tareas en paralelo respetando DAG de dependencias

**Métricas:**
- 3 specs creados (100+ tareas)
- Deploy a AWS en ~4 minutos (push → live)
- Cold start: 4.4s → warm: 1-3ms
- 10,000+ líneas de Go
- Zero lint warnings
- Dominio: 100% coverage

---

## ☁️ Infraestructura AWS

### Recursos aprovisionados

Todos los recursos fueron aprovisionados mediante **AWS SAM** (`template.yaml`) con el comando `sam deploy --guided`, excepto ECR que se provisiona por separado como parte del pipeline CI/CD.

| Recurso | Servicio AWS | Rol en el sistema | Cómo se aprovisionó |
|---------|-------------|-------------------|---------------------|
| **Lambda Function** `pos-ai-first-PosFunction` | AWS Lambda | Compute principal. Ejecuta el binario Go (ARM64, 512MB, 30s timeout). Maneja todas las requests HTTP del POS. | SAM `template.yaml` → `AWS::Lambda::Function` |
| **HTTP API** | API Gateway v2 | Punto de entrada público. Rutea requests al Lambda via proxy integration. | SAM `template.yaml` → `AWS::ApiGatewayV2::Api` |
| **API Stage** `$default` | API Gateway v2 | Stage con auto-deploy habilitado. Expone la URL pública del servicio. | SAM `template.yaml` → `AWS::ApiGatewayV2::Stage` |
| **S3 Bucket** `pos-static-production` | Amazon S3 | Almacenamiento de assets estáticos (JS Alpine.js components). Servido via CloudFront. | SAM `template.yaml` → `AWS::S3::Bucket` |
| **Bucket Policy** | Amazon S3 | Permite a CloudFront leer los objetos del bucket via OAC. Bloquea acceso público directo. | SAM `template.yaml` → `AWS::S3::BucketPolicy` |
| **CloudFront Distribution** | Amazon CloudFront | CDN. Sirve assets estáticos desde S3 y rutea `/` al API Gateway. Reduce latencia global. | SAM `template.yaml` → `AWS::CloudFront::Distribution` |
| **Origin Access Control (OAC)** | Amazon CloudFront | Permite a CloudFront autenticarse contra S3 sin exponer el bucket públicamente. | SAM `template.yaml` → `AWS::CloudFront::OriginAccessControl` |
| **Lambda Execution Role** | AWS IAM | Role con permisos mínimos: CloudWatch Logs, VPC networking, Secrets Manager read. | SAM `template.yaml` → `AWS::IAM::Role` |
| **RDS PostgreSQL** `pos-ai-first-db` | Amazon RDS | Base de datos de producción. PostgreSQL 16, instancia `db.t4g.micro`, almacenamiento 20GB gp3. | Manual via consola AWS (fuera de SAM) |
| **Secrets Manager** (3 secrets) | AWS Secrets Manager | Almacena credenciales cifradas: `pos/db-connection`, `pos/ai-config`, `pos/session-secret`. El Lambda los lee al iniciar con cache in-memory. | Manual via consola AWS (fuera de SAM) |
| **ECR Repository** `pos-ai-first` | Amazon ECR | Registry de imágenes Docker. Almacena la imagen del Lambda (Go ARM64, ~15MB). | GitHub Actions CI/CD pipeline |
| **CloudWatch Log Group** `/aws/lambda/pos-ai-first-PosFunction` | Amazon CloudWatch | Logs de ejecución del Lambda. Retención configurable. Creado automáticamente por Lambda al ejecutarse. | Auto-creado por Lambda |

### Recursos desaprovisionados

La limpieza se realizó en el siguiente orden para evitar dependencias rotas:

| Orden | Recurso | Método de eliminación | Notas |
|-------|---------|----------------------|-------|
| 1 | S3 Bucket (objetos) | Consola S3 → Empty bucket | CloudFormation no puede borrar buckets con contenido |
| 2 | CloudFormation Stack `pos-ai-first` | Consola CloudFormation → Delete stack | Eliminó en cascada: Lambda, API Gateway, CloudFront, OAC, S3 bucket, IAM Role |
| 3 | ECR Repository `pos-ai-first` | Consola ECR → Delete | Eliminó repositorio e imágenes Docker |
| 4 | Secrets Manager (3 secrets) | Consola Secrets Manager → Schedule deletion (7 días) | Mínimo de recuperación de 7 días impuesto por AWS |
| 5 | RDS PostgreSQL `pos-ai-first-db` | Consola RDS → Delete (sin snapshot final) | Tardó ~3 minutos en completarse |
| 6 | RDS Snapshots (9 automáticos) | Consola RDS → Snapshots → Delete | Quedaron huérfanos al borrar la DB |
| 7 | CloudWatch Log Group | Consola CloudWatch → Log Management → Delete | Huérfano post-eliminación del stack |

> **Resultado:** Cuenta AWS en estado limpio. Sin recursos activos ni costos recurrentes.

---

## 📄 Licencia

Proyecto desarrollado para el Bootcamp Kiro × Código Facilito × AWS — Hackathon 2026.

---

## 👤 Autor

**Gabriel Magallón**  
Michoacán, México | Julio 2026
