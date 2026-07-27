# 🧮 POS AI-First — Tu negocio responde preguntas

> **Bootcamp Kiro × Código Facilito | Hackathon 2026**  
> MVP construido en 5 días por Gabriel Magallón desde Michoacán, México.

[![Deploy to AWS](https://github.com/QuantumEdu/bootcamp-kiro-CF/actions/workflows/deploy.yml/badge.svg)](https://github.com/QuantumEdu/bootcamp-kiro-CF/actions/workflows/deploy.yml)

---

## 🎯 ¿Qué es?

Un **punto de venta inteligente** donde el dueño de un negocio puede preguntarle a sus datos en español y recibir respuestas inmediatas — como hablar con WhatsApp, pero sobre sus ventas.

```
👤 "¿Qué vendí hoy?"
🤖 "Hoy llevas $1,250 en 8 ventas. Lideran: tacos al pastor (12), agua natural (8), coca cola (6)."
```

## 🌐 Demo en Vivo

| Entorno | URL | Credenciales |
|---------|-----|--------------|
| **AWS Lambda** | [pos-ai-first.aws](https://zz637vr6cd.execute-api.us-east-1.amazonaws.com/login) | Admin: `1234` / Cajero: `1235` |
| **Local** | `http://localhost:8080` | Mismas credenciales |

---

## 💡 El Problema

Los dueños de pequeños negocios (taquerías, abarrotes, tienditas) llevan sus registros en Excel, libretas o de memoria. Cuando quieren saber "¿cuánto vendí esta semana?" tienen que:

1. Abrir Excel
2. Filtrar por fecha
3. Sumar columnas
4. ...o pedirle a alguien que lo haga

**¿Y si pudieran simplemente preguntar?**

## ✨ La Solución

Un POS completo con un **chat conversacional** que convierte preguntas en español a consultas SQL seguras:

- CRUD de productos, clientes y ventas
- Dashboard con métricas en tiempo real
- Chat AI: pregunta → SQL generado → validado → ejecutado → respuesta formateada
- 5 capas de seguridad NL→SQL (prompt, validación Go, read-only, timeout, auditoría)

---

## 🏗️ Arquitectura

```
┌─────────────────────────────────────┐
│  HTMX + Alpine.js + Tailwind CSS   │  ← Frontend server-driven
├─────────────────────────────────────┤
│  Go HTTP (chi router + algnhsa)    │  ← Lambda o servidor local
├─────────────────────────────────────┤
│  Application (use-cases)           │  ← Lógica de negocio
├─────────────────────────────────────┤
│  Domain (entities + ports)         │  ← Inmutable: CERO cambios al migrar
├─────────────────────────────────────┤
│  SQLite (local) │ PostgreSQL (AWS) │  ← Dual-mode via APP_ENV
│  OpenRouter     │ Bedrock (futuro) │
└─────────────────────────────────────┘
```

**Hexagonal en acción:** Al migrar de local a AWS, se crearon 7 adaptadores PostgreSQL nuevos sin tocar una sola línea del dominio.

---

## 🛡️ Seguridad NL→SQL (5 capas)

| Capa | Defensa |
|------|---------|
| 1. Prompt | Instrucción al LLM: solo generar SELECT |
| 2. Validación Go | Whitelist SELECT/WITH, reject DDL/DML |
| 3. Conexión | Read-only separada |
| 4. Ejecución | Timeout 5s + LIMIT 500 |
| 5. Auditoría | Log de toda query generada |

No confiamos en el LLM. Cada capa es independiente.

---

## 🚀 Ejecutar Localmente

```bash
# Clonar
git clone https://github.com/QuantumEdu/bootcamp-kiro-CF.git
cd bootcamp-kiro-CF

# Configurar
cp .env.example .env
# Editar .env con tu OPENROUTER_API_KEY y SESSION_SECRET

# Seed (datos de demo)
go run cmd/seed/main.go

# Ejecutar
go run cmd/server/main.go

# Abrir http://localhost:8080
# PIN Admin: 1234 | PIN Cajero: 1235
```

## ☁️ Deploy a AWS

La app se despliega automáticamente a AWS Lambda en cada push a `main`:

```
Push → GitHub Actions → Test → Build Docker (ARM64) → ECR → SAM Deploy → Health Check ✅
```

**Infraestructura (100% free tier):**
- Lambda + API Gateway (1M req/mes gratis)
- RDS PostgreSQL db.t4g.micro (12 meses gratis)
- DeepSeek V4 Flash via OpenRouter ($0.09/1M tokens)
- S3 + CloudFront para assets estáticos

**Costo mensual: $0** (primer año con free tier)

---

## 🛠️ Stack Técnico

| Capa | Tecnología |
|------|-----------|
| **Backend** | Go 1.26, chi/v5, hexagonal architecture |
| **Frontend** | HTMX 1.9, Alpine.js 3, Tailwind CSS (CDN) |
| **DB Local** | SQLite (modernc.org/sqlite, pure Go) |
| **DB Cloud** | PostgreSQL 16 (RDS, pgx/v5) |
| **AI** | OpenRouter → DeepSeek V4 Flash (NL→SQL) |
| **Infra** | AWS Lambda (ARM64), API Gateway, SAM |
| **CI/CD** | GitHub Actions |
| **IDE** | Kiro (specs, steering, hooks, powers) |

---

## 📊 Métricas del Proyecto

| Métrica | Valor |
|---------|-------|
| Tiempo de desarrollo | 5 días |
| Specs creados | 3 (MVP, UI fixes, AWS deploy) |
| Tareas ejecutadas | 100+ (paralelas por waves) |
| Archivos Go | 60+ |
| Tests | Domain 100%, Middleware 86%, Use Cases 60% |
| Lint warnings | 0 |
| Cold start Lambda | ~4.4s |
| Warm response | 1-3ms |
| Costo AWS | $0/mes |

---

## 🧠 Construido con Kiro

Este proyecto demuestra el flujo completo de desarrollo con [Kiro](https://kiro.dev):

- **Specs:** Requirements → Design → Tasks con dependency graph
- **Steering:** 6 archivos de reglas persistentes (arquitectura, testing, seguridad, quality, convenciones, design patterns)
- **Powers:** Long-Term Memory, Context7
- **Hooks:** Auto-documentación de prompts
- **Ejecución paralela:** 5 tareas simultáneas por wave respetando dependencias

---

## 📁 Estructura del Proyecto

```
├── cmd/
│   ├── server/main.go      # Entry point local
│   ├── lambda/main.go      # Entry point AWS Lambda
│   ├── seed/main.go        # Seed SQLite
│   └── seedpg/main.go      # Seed PostgreSQL
├── internal/bootstrap/      # Dual-mode router builder
├── src/
│   ├── domain/             # Entities + Ports (INMUTABLE)
│   ├── application/        # Use cases + Services
│   └── infrastructure/     # Adapters (SQLite, PostgreSQL, OpenRouter, Bedrock)
├── templates/              # HTMX templates
├── static/                 # JS (Alpine components)
├── migrations/             # SQLite + PostgreSQL DDL
├── governance/             # PRD, AWS plan, deploy state
├── .kiro/specs/            # Kiro specifications
├── template.yaml           # AWS SAM (IaC)
├── Dockerfile              # Lambda container (ARM64)
└── .github/workflows/      # CI/CD pipeline
```

---

## 👤 Autor

**Gabriel Magallón**  
Michoacán, México  
Bootcamp Kiro × Código Facilito — Hackathon 2026

---

## 📄 Licencia

MIT
