# 🧾 POS AI-First

> Sistema de Punto de Venta con inteligencia artificial conversacional — pregunta en lenguaje natural y obtén respuestas de tus datos reales.

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)
![SQLite](https://img.shields.io/badge/DB-SQLite-003B57?logo=sqlite&logoColor=white)
![HTMX](https://img.shields.io/badge/Frontend-HTMX-blue)

---

## ✨ Características

- 🔐 **Autenticación por PIN** — login tipo POS real (sin emails, sin contraseñas)
- 🛒 **Gestión de productos** — CRUD completo con categorías y control de stock
- 💰 **Registro de ventas** — carrito interactivo con descuento automático de inventario
- 📊 **Dashboard en tiempo real** — métricas de ventas, productos top, tendencias
- 🤖 **Chat NL→SQL** — pregunta "¿qué vendí esta semana?" y obtén respuesta instantánea
- 👥 **Multi-usuario con roles** — admin y cajero con permisos diferenciados
- ⚡ **Interfaz reactiva** — HTMX + Alpine.js para UX fluida sin SPA

---

## 🏗️ Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                        FRONTEND                              │
│            HTMX + Alpine.js + Tailwind CSS                   │
└──────────────────────────┬──────────────────────────────────┘
                           │ HTTP
┌──────────────────────────▼──────────────────────────────────┐
│                    INFRASTRUCTURE                             │
│  ┌─────────────┐  ┌──────────────┐  ┌───────────────────┐  │
│  │  HTTP/Chi   │  │   SQLite     │  │  OpenRouter API   │  │
│  │  Handlers   │  │   Adapters   │  │  (NL→SQL)         │  │
│  └──────┬──────┘  └──────┬───────┘  └────────┬──────────┘  │
└─────────┼────────────────┼────────────────────┼─────────────┘
          │                │                    │
┌─────────▼────────────────▼────────────────────▼─────────────┐
│                      APPLICATION                             │
│  ┌──────────────────┐  ┌──────────────────────────────┐     │
│  │    Use Cases      │  │     NL→SQL Service           │     │
│  │  (authenticate,   │  │  (validate → call LLM →     │     │
│  │   register sale)  │  │   format response)           │     │
│  └────────┬──────────┘  └──────────────┬──────────────┘     │
└───────────┼────────────────────────────┼────────────────────┘
            │                            │
┌───────────▼────────────────────────────▼────────────────────┐
│                        DOMAIN                                │
│  ┌──────────┐  ┌───────────────┐  ┌──────────────────────┐  │
│  │ Entities │  │ Value Objects │  │   Ports (interfaces) │  │
│  │ Product  │  │ PIN           │  │   ProductRepository  │  │
│  │ Sale     │  │               │  │   UserRepository     │  │
│  │ User     │  │               │  │   AIQueryService     │  │
│  │ Inventory│  │               │  │   SaleRepository     │  │
│  └──────────┘  └───────────────┘  └──────────────────────┘  │
└─────────────────────────────────────────────────────────────┘

Flujo de dependencias: Infrastructure → Application → Domain
(Las capas externas dependen de las internas, nunca al revés)
```

---

## 🚀 Quick Start

### Prerrequisitos

- Go 1.22+ instalado
- (Opcional) Cuenta en [OpenRouter](https://openrouter.ai/) para la funcionalidad NL→SQL

### Instalación

```bash
# 1. Clonar el repositorio
git clone https://github.com/QuantumEdu/bootcamp-kiro-CF.git
cd bootcamp-kiro-CF

# 2. Configurar variables de entorno
cp .env.example .env
# Editar .env y agregar tu OPENROUTER_API_KEY

# 3. Ejecutar (auto-ejecuta migraciones + seed)
make run

# 4. Abrir en el navegador
# http://localhost:8080
```

### Credenciales de demo

| PIN    | Rol    | Permisos                          |
|--------|--------|-----------------------------------|
| `1234` | Admin  | Acceso completo (productos, ventas, métricas, chat) |
| `123`  | Cajero | Ventas y consultas básicas        |

---

## 🎬 Demo Script — Consultas NL→SQL

Una vez logueado, abre el **Chat** y prueba estas consultas en lenguaje natural:

```
1. "¿Cuánto vendí hoy?"
   → Muestra el total de ventas del día actual

2. "¿Cuáles son los 5 productos más vendidos?"
   → Lista los productos con más unidades vendidas

3. "¿Qué productos tienen menos de 10 en stock?"
   → Alerta de inventario bajo

4. "¿Cuál fue la venta más grande de esta semana?"
   → Identifica la transacción de mayor monto

5. "¿Cuántas ventas hizo cada cajero este mes?"
   → Desglose de productividad por usuario
```

> 💡 El sistema traduce tu pregunta a SQL, valida que sea segura (solo SELECT),
> la ejecuta contra tu base de datos real, y formatea la respuesta.

---

## 🛠️ Desarrollo

### Makefile targets

```bash
make run       # Ejecutar servidor (localhost:8080)
make build     # Compilar binario → bin/pos
make test      # Correr todos los tests con cobertura
make lint      # golangci-lint
make fmt       # Formatear código (gofmt)
make seed      # Cargar datos de demostración
make clean     # Limpiar binarios y DB
```

### Testing

```bash
# Todos los tests
go test ./... -v -cover

# Tests de un paquete específico
go test ./src/domain/entities/ -v

# Tests con race detector
go test -race ./...
```

### Linting

```bash
# Lint completo
golangci-lint run

# Análisis estático
go vet ./...
```

---

## 📦 Tech Stack

| Tecnología | Decisión | Justificación |
|------------|----------|---------------|
| **Go 1.22+** | Backend | Compilación rápida, binario único, concurrencia nativa |
| **chi/v5** | Router HTTP | Ligero, idiomático, middleware composable |
| **SQLite** (modernc.org) | Base de datos | Sin servidor externo, ideal para POS local, pure Go |
| **HTMX** | Interactividad | Server-driven UI sin complejidad de SPA |
| **Alpine.js** | Reactividad local | Componentes ligeros para UI interactiva |
| **Tailwind CSS** | Estilos | Utility-first, desarrollo rápido, CDN |
| **OpenRouter** | AI/NL→SQL | Acceso multi-modelo (Claude, GPT), sin config AWS |
| **bcrypt** | Hash de PINs | Seguro, estándar, parte de golang.org/x/crypto |
| **alexedwards/scs** | Sesiones | Ligero, almacenamiento en SQLite |
| **godotenv** | Config | Carga `.env` sin magia |

### Decisiones arquitectónicas clave

1. **SQLite sobre Postgres** — simplicidad para POS local, sin servidor externo
2. **HTMX sobre SPA** — menor complejidad frontend, rendering server-side
3. **OpenRouter sobre Bedrock** — acceso inmediato sin configurar AWS
4. **PIN sobre OAuth** — contexto POS real (meseros no tienen email corporativo)
5. **Hexagonal sobre MVC** — separación clara, testabilidad, bajo acoplamiento

---

## 📁 Estructura del Proyecto

```
bootcamp/
├── cmd/
│   ├── server/main.go            # Entry point del servidor
│   └── seed/main.go              # Script de datos de demostración
├── src/
│   ├── domain/
│   │   ├── entities/             # Product, Sale, User, Inventory
│   │   ├── value_objects/        # PIN (validación + hashing)
│   │   └── ports/                # Interfaces (contratos del dominio)
│   ├── application/
│   │   ├── use_cases/            # AuthenticateUser, RegisterSale
│   │   ├── nlsql/                # Servicio NL→SQL (validar, llamar LLM, formatear)
│   │   └── dtos/                 # Estructuras de transferencia
│   └── infrastructure/
│       ├── adapters/             # SQLite repos, OpenRouter client
│       ├── database/             # Conexión, migraciones
│       ├── http/
│       │   ├── handlers/         # Auth, Chat, Products, Sales, Metrics
│       │   ├── middleware/       # Auth middleware
│       │   └── session.go        # Gestión de sesiones
│       └── config/               # Variables de entorno
├── migrations/                   # SQL schema (001_init, 002_sessions)
├── templates/                    # HTML templates (HTMX)
│   ├── layout.html               # Layout base
│   ├── login.html                # Pantalla de login
│   ├── chat/                     # UI del chat NL→SQL
│   ├── products/                 # CRUD productos
│   ├── sales/                    # Registro de ventas
│   └── metrics/                  # Dashboard
├── static/                       # Assets estáticos (JS)
├── data/                         # SQLite DB (gitignored en prod)
├── governance/                   # PRD y brief del proyecto
├── .env.example                  # Template de variables de entorno
├── Makefile                      # Comandos de desarrollo
├── go.mod / go.sum               # Dependencias Go
└── .golangci.yml                 # Configuración del linter
```

---

## 🎓 Contexto — Bootcamp Código Facilito + Kiro

Este proyecto fue desarrollado durante el **Bootcamp de Código Facilito** en colaboración con **Kiro** (IDE con AI de Amazon).

### Objetivo del bootcamp
Construir un sistema funcional en 5 días que demuestre el impacto de la inteligencia artificial aplicada a un problema real de negocio.

### ¿Por qué un POS con AI?
Las PYMES en Latinoamérica operan con datos que no pueden consultar fácilmente. Este POS permite que el dueño de un restaurante o tienda pregunte en español "¿qué vendí esta semana?" y obtenga respuesta inmediata — sin saber SQL, sin dashboards complicados.

### Criterios de evaluación
- **Impacto tecnológico (30%)** — NL→SQL para PYMES sin acceso a tech
- **Innovación (30%)** — POS conversacional desde el diseño
- **Software funcional (30%)** — CRUD + chat + dashboard demostrable
- **Uso de AWS y Kiro (10%)** — Desarrollo guiado por specs y steering

---

## 📄 Licencia

MIT
