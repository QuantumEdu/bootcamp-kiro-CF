# Design Document

## Overview

This design covers UI fixes (logout button, product creation button, HTMX cache headers) and two new features (client CRUD, admin API key configuration) for the POS AI-First application. All new backend components follow the existing hexagonal architecture: domain ports → infrastructure adapters → HTTP handlers wired through chi/v5.

## Architecture

### Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        HTTP Layer (chi/v5)                    │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────────┐  │
│  │ ClientHandler │  │ AdminConfig  │  │ Existing Handlers │  │
│  │              │  │   Handler    │  │ (Auth, Product,   │  │
│  │ - List       │  │              │  │  Sales, Metrics)  │  │
│  │ - CreateForm │  │ - Show       │  │                   │  │
│  │ - Create     │  │ - Update     │  │ + Logout (Auth)   │  │
│  └──────┬───────┘  └──────┬───────┘  └───────────────────┘  │
│         │                  │                                  │
│  RequireAuth         RequireAuth + RequireRole("admin")       │
└─────────┼──────────────────┼─────────────────────────────────┘
          │                  │
┌─────────▼──────────────────▼─────────────────────────────────┐
│                    Application Layer                           │
│  ┌──────────────────┐  ┌──────────────────────────────────┐  │
│  │ CreateClient UC   │  │ CryptoService (encrypt/decrypt)  │  │
│  │ ListClients UC    │  └──────────────────────────────────┘  │
│  └──────────┬────────┘                                        │
└─────────────┼────────────────────────────────────────────────┘
              │
┌─────────────▼────────────────────────────────────────────────┐
│                      Domain Layer                              │
│  ┌────────────────┐  ┌──────────────────────────────────┐    │
│  │ Client Entity   │  │ Ports                            │    │
│  │ - ID           │  │ - ClientRepository (interface)    │    │
│  │ - Nombre       │  │ - ConfigRepository (interface)    │    │
│  │ - Telefono     │  │                                   │    │
│  │ - Direccion    │  └──────────────────────────────────┘    │
│  │ - CreatedAt    │                                           │
│  └────────────────┘                                           │
└──────────────────────────────────────────────────────────────┘
              │
┌─────────────▼────────────────────────────────────────────────┐
│                   Infrastructure Layer                         │
│  ┌────────────────────────┐  ┌─────────────────────────────┐ │
│  │ SQLiteClientRepository │  │ SQLiteConfigRepository       │ │
│  │ (adapters/)            │  │ (adapters/)                  │ │
│  └────────────────────────┘  └─────────────────────────────┘ │
└──────────────────────────────────────────────────────────────┘
```

### Data Flow

1. **Client CRUD**: Browser → chi router → RequireAuth middleware → ClientHandler → CreateClient/ListClients use case → ClientRepository port → SQLiteClientRepository adapter → SQLite `clientes` table.
2. **Admin Config**: Browser → chi router → RequireAuth → RequireRole("admin") → AdminConfigHandler → CryptoService (encrypt/decrypt) → ConfigRepository port → SQLiteConfigRepository adapter → SQLite `configuracion` table.
3. **Logout**: Sidebar button → POST /logout → AuthHandler.Logout → session.Destroy() → redirect /login.

## Components and Interfaces

### 1. Domain Layer Additions

#### Client Entity

```go
// src/domain/entities/client.go
package entities

import "time"

type Client struct {
    ID        int64
    Nombre    string
    Telefono  string
    Direccion string
    CreatedAt time.Time
}

// Validate checks domain invariants for a client.
func (c *Client) Validate() error {
    if strings.TrimSpace(c.Nombre) == "" {
        return ErrClientNameRequired
    }
    return nil
}
```

#### ClientRepository Port

```go
// src/domain/ports/client_repository.go
package ports

import (
    "context"
    "github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
)

type ClientRepository interface {
    Create(ctx context.Context, client *entities.Client) error
    List(ctx context.Context) ([]entities.Client, error)
}
```

#### ConfigRepository Port

```go
// src/domain/ports/config_repository.go
package ports

import "context"

type ConfigRepository interface {
    Get(ctx context.Context, clave string) (string, error)
    Set(ctx context.Context, clave, valor string) error
}
```

### 2. Application Layer Additions

#### CryptoService

Handles AES-GCM encryption/decryption for the API key. Located at `src/application/services/crypto.go`.

```go
package services

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "io"
)

type CryptoService struct {
    key []byte // 32-byte AES-256 key derived from SESSION_SECRET
}

// NewCryptoService derives a 256-bit key from the secret using SHA-256.
func NewCryptoService(secret string) *CryptoService {
    hash := sha256.Sum256([]byte(secret))
    return &CryptoService{key: hash[:]}
}

// Encrypt encrypts plaintext using AES-GCM. Returns hex-encoded ciphertext.
func (s *CryptoService) Encrypt(plaintext string) (string, error) {
    block, err := aes.NewCipher(s.key)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return hex.EncodeToString(sealed), nil
}

// Decrypt decrypts hex-encoded AES-GCM ciphertext.
func (s *CryptoService) Decrypt(ciphertextHex string) (string, error) {
    data, err := hex.DecodeString(ciphertextHex)
    if err != nil {
        return "", err
    }
    block, err := aes.NewCipher(s.key)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return "", errors.New("ciphertext too short")
    }
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }
    return string(plaintext), nil
}
```

#### Client Use Cases

```go
// src/application/use_cases/client_use_cases.go
package use_cases

import (
    "context"
    "github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
    "github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

type CreateClient struct {
    repo ports.ClientRepository
}

func NewCreateClient(repo ports.ClientRepository) *CreateClient {
    return &CreateClient{repo: repo}
}

type CreateClientInput struct {
    Nombre    string
    Telefono  string
    Direccion string
}

func (uc *CreateClient) Execute(ctx context.Context, input CreateClientInput) (*entities.Client, error) {
    client := &entities.Client{
        Nombre:    input.Nombre,
        Telefono:  input.Telefono,
        Direccion: input.Direccion,
    }
    if err := client.Validate(); err != nil {
        return nil, err
    }
    if err := uc.repo.Create(ctx, client); err != nil {
        return nil, err
    }
    return client, nil
}

type ListClients struct {
    repo ports.ClientRepository
}

func NewListClients(repo ports.ClientRepository) *ListClients {
    return &ListClients{repo: repo}
}

func (uc *ListClients) Execute(ctx context.Context) ([]entities.Client, error) {
    return uc.repo.List(ctx)
}
```

### 3. Infrastructure Layer Additions

#### SQLiteClientRepository

```go
// src/infrastructure/adapters/sqlite_client_repository.go
package adapters

import (
    "context"
    "database/sql"
    "github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
)

type SQLiteClientRepository struct {
    db *sql.DB
}

func NewSQLiteClientRepository(db *sql.DB) *SQLiteClientRepository {
    return &SQLiteClientRepository{db: db}
}

func (r *SQLiteClientRepository) Create(ctx context.Context, client *entities.Client) error {
    result, err := r.db.ExecContext(ctx,
        `INSERT INTO clientes (nombre, telefono, direccion) VALUES (?, ?, ?)`,
        client.Nombre, client.Telefono, client.Direccion)
    if err != nil {
        return err
    }
    id, _ := result.LastInsertId()
    client.ID = id
    return nil
}

func (r *SQLiteClientRepository) List(ctx context.Context) ([]entities.Client, error) {
    rows, err := r.db.QueryContext(ctx,
        `SELECT id, nombre, telefono, direccion, created_at FROM clientes ORDER BY nombre ASC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var clients []entities.Client
    for rows.Next() {
        var c entities.Client
        if err := rows.Scan(&c.ID, &c.Nombre, &c.Telefono, &c.Direccion, &c.CreatedAt); err != nil {
            return nil, err
        }
        clients = append(clients, c)
    }
    return clients, rows.Err()
}
```

#### SQLiteConfigRepository

```go
// src/infrastructure/adapters/sqlite_config_repository.go
package adapters

import (
    "context"
    "database/sql"
)

type SQLiteConfigRepository struct {
    db *sql.DB
}

func NewSQLiteConfigRepository(db *sql.DB) *SQLiteConfigRepository {
    return &SQLiteConfigRepository{db: db}
}

func (r *SQLiteConfigRepository) Get(ctx context.Context, clave string) (string, error) {
    var valor string
    err := r.db.QueryRowContext(ctx,
        `SELECT valor FROM configuracion WHERE clave = ?`, clave).Scan(&valor)
    if err == sql.ErrNoRows {
        return "", nil
    }
    return valor, err
}

func (r *SQLiteConfigRepository) Set(ctx context.Context, clave, valor string) error {
    _, err := r.db.ExecContext(ctx,
        `INSERT INTO configuracion (clave, valor) VALUES (?, ?)
         ON CONFLICT(clave) DO UPDATE SET valor = excluded.valor`,
        clave, valor)
    return err
}
```

### 4. HTTP Handlers

#### ClientHandler

```go
// src/infrastructure/http/handlers/clients.go
package handlers

import (
    "html/template"
    "net/http"
    "strings"

    "github.com/QuantumEdu/bootcamp-kiro-CF/src/application/use_cases"
)

type ClientHandler struct {
    createUC *use_cases.CreateClient
    listUC   *use_cases.ListClients
    tmpl     *template.Template
}

func NewClientHandler(create *use_cases.CreateClient, list *use_cases.ListClients, tmpl *template.Template) *ClientHandler {
    return &ClientHandler{createUC: create, listUC: list, tmpl: tmpl}
}

func (h *ClientHandler) List(w http.ResponseWriter, r *http.Request) {
    clients, err := h.listUC.Execute(r.Context())
    if err != nil {
        http.Error(w, "Error al cargar clientes", http.StatusInternalServerError)
        return
    }
    data := map[string]interface{}{
        "PageTitle": "Clientes",
        "Clients":   clients,
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    h.tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *ClientHandler) CreateForm(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{"PageTitle": "Nuevo Cliente"}
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    h.tmpl.ExecuteTemplate(w, "clients/form.html", data)
}

func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    nombre := strings.TrimSpace(r.FormValue("nombre"))
    input := use_cases.CreateClientInput{
        Nombre:    nombre,
        Telefono:  strings.TrimSpace(r.FormValue("telefono")),
        Direccion: strings.TrimSpace(r.FormValue("direccion")),
    }
    _, err := h.createUC.Execute(r.Context(), input)
    if err != nil {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.WriteHeader(http.StatusUnprocessableEntity)
        // Render validation error inline
        return
    }
    http.Redirect(w, r, "/clientes", http.StatusSeeOther)
}
```

#### AdminConfigHandler

```go
// src/infrastructure/http/handlers/admin_config.go
package handlers

import (
    "html/template"
    "net/http"
    "strings"

    "github.com/QuantumEdu/bootcamp-kiro-CF/src/application/services"
    "github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
)

type AdminConfigHandler struct {
    configRepo ports.ConfigRepository
    crypto     *services.CryptoService
    tmpl       *template.Template
}

func NewAdminConfigHandler(repo ports.ConfigRepository, crypto *services.CryptoService, tmpl *template.Template) *AdminConfigHandler {
    return &AdminConfigHandler{configRepo: repo, crypto: crypto, tmpl: tmpl}
}

func (h *AdminConfigHandler) Show(w http.ResponseWriter, r *http.Request) {
    encrypted, _ := h.configRepo.Get(r.Context(), "openrouter_api_key")
    masked := ""
    if encrypted != "" {
        decrypted, err := h.crypto.Decrypt(encrypted)
        if err == nil && len(decrypted) >= 4 {
            masked = strings.Repeat("*", len(decrypted)-4) + decrypted[len(decrypted)-4:]
        }
    }
    data := map[string]interface{}{
        "PageTitle":  "Configuración",
        "MaskedKey":  masked,
        "HasKey":     encrypted != "",
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    h.tmpl.ExecuteTemplate(w, "layout.html", data)
}

func (h *AdminConfigHandler) Update(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    apiKey := strings.TrimSpace(r.FormValue("api_key"))
    if apiKey == "" {
        http.Error(w, "La API key no puede estar vacía", http.StatusUnprocessableEntity)
        return
    }
    encrypted, err := h.crypto.Encrypt(apiKey)
    if err != nil {
        http.Error(w, "Error al cifrar la clave", http.StatusInternalServerError)
        return
    }
    if err := h.configRepo.Set(r.Context(), "openrouter_api_key", encrypted); err != nil {
        http.Error(w, "Error al guardar la configuración", http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/admin/config", http.StatusSeeOther)
}
```

### 5. HTMX Cache-Control Middleware

```go
// src/infrastructure/http/middleware/nocache.go
package middleware

import "net/http"

// NoCache sets Cache-Control: no-store on responses to prevent HTMX fragments
// from being served stale after navigation or server restart.
func NoCache(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "no-store")
        next.ServeHTTP(w, r)
    })
}
```

### 6. Router Wiring (changes to cmd/server/main.go)

```go
// Inside protected route group:
r.Get("/clientes", clientHandler.List)
r.Get("/clientes/new", clientHandler.CreateForm)
r.Post("/clientes", clientHandler.Create)

// Admin-only group:
r.Group(func(r chi.Router) {
    r.Use(mw.RequireRole(sessionManager, "admin"))
    r.Get("/admin/config", adminConfigHandler.Show)
    r.Post("/admin/config", adminConfigHandler.Update)
})

// HTMX fragment routes get NoCache middleware:
r.Group(func(r chi.Router) {
    r.Use(mw.NoCache)
    r.Get("/api/metrics/ventas-hoy", metricsHandler.VentasHoy)
    // ... all /api/metrics/* and /api/productos routes
})
```

### 7. Template Changes

#### layout.html — Sidebar Modifications

- Add `<a href="/clientes">Clientes</a>` nav link (visible to all authenticated users).
- Add conditional `{{if eq .UserRole "admin"}}<a href="/admin/config">Configuración</a>{{end}}` nav link.
- Add logout button in footer section below user info:

```html
<div class="p-4 border-t border-indigo-800">
    <div class="flex items-center gap-3">
        <div class="w-8 h-8 bg-indigo-600 rounded-full flex items-center justify-center text-sm font-bold">{{.UserInitial}}</div>
        <div class="flex-1"><p class="text-sm font-medium">{{.UserName}}</p><p class="text-xs text-indigo-300">{{.UserRole}}</p></div>
    </div>
    <form method="POST" action="/logout" class="mt-3">
        <button type="submit" class="w-full flex items-center gap-2 px-3 py-2 text-sm text-indigo-200 hover:bg-indigo-800 rounded-lg transition-colors">
            <!-- logout icon -->
            Cerrar Sesión
        </button>
    </form>
</div>
```

#### New Templates

- `templates/clients/list.html` — Client list page with "Nuevo Cliente" button and table.
- `templates/clients/form.html` — Client creation form with nombre, telefono, direccion fields.
- `templates/admin/config.html` — Admin config page with masked key display and update form.

### 8. Interfaces

#### ClientRepository Port

| Method | Input | Output | Description |
|--------|-------|--------|-------------|
| `Create` | `ctx, *Client` | `error` | Insert new client |
| `List` | `ctx` | `([]Client, error)` | List all clients ordered by name |

#### ConfigRepository Port

| Method | Input | Output | Description |
|--------|-------|--------|-------------|
| `Get` | `ctx, clave string` | `(string, error)` | Get config value by key; returns "" if not found |
| `Set` | `ctx, clave, valor string` | `error` | Upsert config key-value pair |

#### CryptoService

| Method | Input | Output | Description |
|--------|-------|--------|-------------|
| `Encrypt` | `plaintext string` | `(string, error)` | AES-GCM encrypt, returns hex |
| `Decrypt` | `ciphertextHex string` | `(string, error)` | AES-GCM decrypt from hex |

#### MaskAPIKey (helper)

| Input | Output | Description |
|-------|--------|-------------|
| `key string` | `string` | Returns `****` + last 4 chars; empty string if key is empty |

## Data Models

### Client Entity

| Field | Type | Constraints |
|-------|------|-------------|
| ID | int64 | Auto-increment PK |
| Nombre | string | Required, non-empty after trim |
| Telefono | string | Optional |
| Direccion | string | Optional |
| CreatedAt | time.Time | Set on insert |

### Configuracion Table (existing)

| Column | Type | Notes |
|--------|------|-------|
| id | INTEGER | PK |
| clave | TEXT | UNIQUE, e.g. "openrouter_api_key" |
| valor | TEXT | Encrypted hex string for sensitive values |

## Testing Strategy

### Unit Tests (Example-based)

- **Logout handler**: POST /logout with valid session → session destroyed, redirect to /login.
- **Client list handler**: GET /clientes with seeded data → 200 with all client fields rendered.
- **Client creation**: POST /clientes with valid nombre → redirect to /clientes, record in DB.
- **Admin config access**: GET /admin/config as admin → 200; as non-admin → 403.
- **Empty state**: GET /clientes with empty DB → empty state message rendered.
- **Template presence**: Sidebar contains logout button, "Clientes" link, conditional "Configuración" link.

### Property Tests (100+ iterations)

- **AES-GCM round-trip** (Property 1): random strings survive encrypt→decrypt cycle.
- **Whitespace rejection** (Property 2): random whitespace strings are rejected by Client.Validate().
- **API key masking** (Property 3): random strings ≥4 chars mask correctly.
- **Role enforcement** (Property 4): random non-admin roles get 403.
- **Cache headers** (Property 5): responses through NoCache middleware always include no-store.

### Integration Tests

- **Client CRUD flow**: create client → list includes it → fields match.
- **Admin config flow**: save encrypted key → reload page → masked key displayed.
- **Logout flow**: login → logout → subsequent requests redirect to /login.

## Error Handling

| Scenario | Behavior |
|----------|----------|
| Empty client name on submit | Return 422 with validation error message |
| DB error on client create | Return 500 with generic error |
| Non-admin accesses /admin/config | Return 403 via RequireRole middleware |
| Empty SESSION_SECRET at startup | `log.Fatalf` — refuse to start |
| Decryption failure (corrupt data) | Show empty state on config page, log error |
| Empty API key on submit | Return 422 with validation error |
| AES-GCM encryption failure | Return 500 with generic error |

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system — essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: AES-GCM encryption round-trip

*For any* non-empty string used as an API key, encrypting it with `CryptoService.Encrypt` and then decrypting the result with `CryptoService.Decrypt` (using the same secret) SHALL produce the original string.

**Validates: Requirements 6.6, 6.7**

### Property 2: Whitespace-only client names are rejected

*For any* string composed entirely of whitespace characters (spaces, tabs, newlines), calling `Client.Validate()` SHALL return an error, preventing the record from being persisted.

**Validates: Requirements 4.4**

### Property 3: API key masking reveals only last four characters

*For any* string of length ≥ 4, the `MaskAPIKey` function SHALL return a string where the first `len - 4` characters are asterisks and the last 4 characters match the original string's last 4 characters. The total length of the masked string equals the original length.

**Validates: Requirements 6.3**

### Property 4: Non-admin users cannot access admin routes

*For any* authenticated user whose role is not "admin", requests to the admin config route SHALL receive a 403 Forbidden response from the `RequireRole("admin")` middleware.

**Validates: Requirements 6.2**

### Property 5: HTMX fragment responses include no-store cache header

*For any* HTTP response served through the NoCache middleware, the `Cache-Control` header SHALL contain "no-store".

**Validates: Requirements 5.3**
