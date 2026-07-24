# Design Document: AWS Lambda Deploy

## Overview

Deploy the POS AI-First Go application to AWS Lambda using a container image runtime (ARM64). The chi router is wrapped by `algnhsa` to handle API Gateway HTTP events. Domain and application layers remain untouched — only infrastructure adapters change, per the hexagonal architecture.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Clients (Browser)                         │
└───────────────┬────────────────────────────┬────────────────────┘
                │                            │
    ┌───────────▼───────────┐    ┌───────────▼───────────────┐
    │  CloudFront (CDN)     │    │  API Gateway HTTP API (v2) │
    │  + S3 (static/)       │    │  catch-all → Lambda        │
    └───────────────────────┘    └───────────┬───────────────┘
                                             │
                              ┌──────────────▼──────────────┐
                              │   Lambda (Container Image)   │
                              │   ARM64, 512MB, 30s timeout  │
                              │                              │
                              │  cmd/lambda/main.go          │
                              │  ├─ algnhsa.ListenAndServe() │
                              │  ├─ chi router (all routes)  │
                              │  └─ templates/ embedded      │
                              └──┬──────────┬──────────┬─────┘
                                 │          │          │
                   ┌─────────────▼┐  ┌──────▼──────┐  ┌▼───────────────┐
                   │ RDS PostgreSQL│  │   Bedrock   │  │Secrets Manager │
                   │ (pgxpool)     │  │  Claude 3   │  │ DB/Session/AI  │
                   └──────────────┘  └─────────────┘  └────────────────┘
```

## Components and Interfaces

### 1. Lambda Entry Point — `cmd/lambda/main.go`

Separate binary from `cmd/server/main.go`. Shares the same router-building logic via a shared `bootstrap` package.

### 2. Dual-Mode Bootstrap — `internal/bootstrap/`

A shared package that builds the chi router with all dependencies. Both `cmd/server/main.go` and `cmd/lambda/main.go` call into it, differing only in how they start the HTTP listener.

### 3. PostgreSQL Adapters — `src/infrastructure/adapters/`

Seven new adapter files implementing domain ports using `pgxpool.Pool`.

### 4. Bedrock Adapter — `src/infrastructure/adapters/bedrock_query_service.go`

Implements `ports.AIQueryService` using AWS SDK Go v2 `bedrockruntime`.

### 5. Secrets Loader — `src/infrastructure/config/secrets.go`

Retrieves secrets from AWS Secrets Manager on Lambda cold start.

### 6. PostgreSQL Migrations — `migrations/postgres/`

Equivalent DDL to SQLite migrations with PostgreSQL-native types.

### 7. SAM Template — `template.yaml`

Defines Lambda, HttpApi, S3, CloudFront, IAM, and VPC configuration.

### 8. CI/CD — `.github/workflows/deploy.yml`

GitHub Actions workflow: test → build image → push ECR → sam deploy → health check.

---

## New Files and Interfaces

### `cmd/lambda/main.go`

```go
package main

import (
    "log"
    "os"

    "github.com/QuantumEdu/bootcamp-kiro-CF/internal/bootstrap"
    "github.com/akrylysov/algnhsa"
)

func main() {
    router, cleanup, err := bootstrap.BuildRouter(bootstrap.Config{
        AppEnv: "lambda",
    })
    if err != nil {
        log.Fatalf("bootstrap failed: %v", err)
    }
    defer cleanup()

    algnhsa.ListenAndServe(router, nil)
}
```

### `internal/bootstrap/bootstrap.go`

```go
package bootstrap

import (
    "context"
    "fmt"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

// Config holds bootstrap configuration resolved from env or Secrets Manager.
type Config struct {
    AppEnv              string
    Port                string
    DatabaseURL         string // PostgreSQL connection string (lambda)
    DatabasePath        string // SQLite path (local)
    SessionSecret       string
    BedrockModelID      string
    BedrockRegion       string
    OpenRouterAPIKey    string
    OpenRouterModel     string
    QueryTimeoutSeconds int
    PINMaxAttempts      int
    PINLockoutMinutes   int
}

// BuildRouter constructs the chi router with all dependencies wired.
// Returns the router, a cleanup function, and any initialization error.
func BuildRouter(cfg Config) (http.Handler, func(), error) {
    // ... adapter selection based on cfg.AppEnv
}
```

### `src/infrastructure/adapters/postgres_product_repository.go`

```go
package adapters

import (
    "context"
    "fmt"

    "github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
    "github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
    "github.com/jackc/pgx/v5/pgxpool"
)

// PostgresProductRepository implements ports.ProductRepository using pgxpool.
type PostgresProductRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresProductRepository(pool *pgxpool.Pool) *PostgresProductRepository {
    return &PostgresProductRepository{pool: pool}
}

func (r *PostgresProductRepository) Create(ctx context.Context, p *entities.Product) error {
    const q = `INSERT INTO productos (nombre, sku, categoria_id, precio_venta, 
        precio_compra, stock_actual, stock_minimo, unidad, activo)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, created_at, updated_at`
    return r.pool.QueryRow(ctx, q,
        p.Nombre, p.SKU, p.CategoriaID, p.PrecioVenta,
        p.PrecioCompra, p.StockActual, p.StockMinimo, p.Unidad, p.Activo,
    ).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PostgresProductRepository) Update(ctx context.Context, p *entities.Product) error {
    const q = `UPDATE productos SET nombre=$1, sku=$2, categoria_id=$3, 
        precio_venta=$4, precio_compra=$5, stock_actual=$6, stock_minimo=$7, 
        unidad=$8, activo=$9, updated_at=NOW() WHERE id=$10`
    _, err := r.pool.Exec(ctx, q,
        p.Nombre, p.SKU, p.CategoriaID, p.PrecioVenta,
        p.PrecioCompra, p.StockActual, p.StockMinimo, p.Unidad, p.Activo, p.ID,
    )
    return wrapErr(err, "updating product", p.ID)
}

func (r *PostgresProductRepository) FindByID(ctx context.Context, id int64) (*entities.Product, error) {
    // ... parameterized SELECT with $1 placeholder
}

func (r *PostgresProductRepository) List(ctx context.Context, filter ports.ProductFilter) ([]entities.Product, error) {
    // ... dynamic WHERE clause with parameterized values
}

func (r *PostgresProductRepository) Deactivate(ctx context.Context, id int64) error {
    // ... UPDATE activo=false WHERE id=$1
}

func (r *PostgresProductRepository) FindLowStock(ctx context.Context) ([]entities.Product, error) {
    // ... WHERE stock_actual <= stock_minimo AND activo=true
}
```

### `src/infrastructure/adapters/postgres_sale_repository.go`

```go
package adapters

import (
    "context"
    "fmt"

    "github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/entities"
    "github.com/QuantumEdu/bootcamp-kiro-CF/src/domain/ports"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSaleRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresSaleRepository(pool *pgxpool.Pool) *PostgresSaleRepository {
    return &PostgresSaleRepository{pool: pool}
}

func (r *PostgresSaleRepository) Create(ctx context.Context, sale *entities.Sale) error {
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("beginning sale transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // Insert sale header
    err = tx.QueryRow(ctx,
        `INSERT INTO ventas (usuario_id, cliente_id, total, metodo_pago)
         VALUES ($1,$2,$3,$4) RETURNING id, created_at`,
        sale.UsuarioID, sale.ClienteID, sale.Total, sale.MetodoPago,
    ).Scan(&sale.ID, &sale.CreatedAt)
    if err != nil {
        return fmt.Errorf("inserting sale: %w", err)
    }

    // Insert sale items
    for i := range sale.Items {
        item := &sale.Items[i]
        _, err = tx.Exec(ctx,
            `INSERT INTO venta_items (venta_id, producto_id, cantidad, precio_unitario, subtotal)
             VALUES ($1,$2,$3,$4,$5)`,
            sale.ID, item.ProductoID, item.Cantidad, item.PrecioUnitario, item.Subtotal,
        )
        if err != nil {
            return fmt.Errorf("inserting sale item %d: %w", i, err)
        }
    }

    return tx.Commit(ctx)
}
```

### `src/infrastructure/adapters/postgres_user_repository.go`

```go
type PostgresUserRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository { ... }

func (r *PostgresUserRepository) FindByID(ctx context.Context, id int64) (*entities.User, error) { ... }
func (r *PostgresUserRepository) FindByPINHash(ctx context.Context, pinHash string) (*entities.User, error) { ... }
func (r *PostgresUserRepository) FindAll(ctx context.Context) ([]entities.User, error) { ... }
func (r *PostgresUserRepository) IncrementFailedAttempts(ctx context.Context, id int64) error {
    _, err := r.pool.Exec(ctx,
        `UPDATE usuarios SET failed_attempts = failed_attempts + 1 WHERE id = $1`, id)
    return wrapErr(err, "incrementing failed attempts", id)
}
func (r *PostgresUserRepository) Lock(ctx context.Context, id int64, until time.Time) error { ... }
func (r *PostgresUserRepository) ResetAttempts(ctx context.Context, id int64) error { ... }
```

### `src/infrastructure/adapters/postgres_client_repository.go`

```go
type PostgresClientRepository struct { pool *pgxpool.Pool }
func NewPostgresClientRepository(pool *pgxpool.Pool) *PostgresClientRepository { ... }
func (r *PostgresClientRepository) Create(ctx context.Context, c *entities.Client) error { ... }
func (r *PostgresClientRepository) List(ctx context.Context) ([]entities.Client, error) { ... }
```

### `src/infrastructure/adapters/postgres_inventory_repository.go`

```go
type PostgresInventoryRepository struct { pool *pgxpool.Pool }
func NewPostgresInventoryRepository(pool *pgxpool.Pool) *PostgresInventoryRepository { ... }
func (r *PostgresInventoryRepository) Create(ctx context.Context, m *entities.InventoryMovement) error { ... }
func (r *PostgresInventoryRepository) FindByProduct(ctx context.Context, productID int64) ([]entities.InventoryMovement, error) { ... }
```

### `src/infrastructure/adapters/postgres_config_repository.go`

```go
type PostgresConfigRepository struct { pool *pgxpool.Pool }
func NewPostgresConfigRepository(pool *pgxpool.Pool) *PostgresConfigRepository { ... }
func (r *PostgresConfigRepository) Get(ctx context.Context, clave string) (string, error) {
    var valor string
    err := r.pool.QueryRow(ctx,
        `SELECT valor FROM configuracion WHERE clave = $1`, clave).Scan(&valor)
    if err == pgx.ErrNoRows {
        return "", nil
    }
    return valor, wrapErr(err, "getting config", 0)
}
func (r *PostgresConfigRepository) Set(ctx context.Context, clave, valor string) error {
    _, err := r.pool.Exec(ctx,
        `INSERT INTO configuracion (clave, valor) VALUES ($1, $2)
         ON CONFLICT (clave) DO UPDATE SET valor = EXCLUDED.valor`, clave, valor)
    return wrapErr(err, "setting config", 0)
}
```

### `src/infrastructure/adapters/postgres_metrics_repository.go`

```go
type PostgresMetricsRepository struct { pool *pgxpool.Pool }
func NewPostgresMetricsRepository(pool *pgxpool.Pool) *PostgresMetricsRepository { ... }

func (r *PostgresMetricsRepository) VentasHoy(ctx context.Context) (int, float64, error) {
    var count int
    var total float64
    err := r.pool.QueryRow(ctx,
        `SELECT COUNT(*), COALESCE(SUM(total), 0) FROM ventas
         WHERE created_at >= CURRENT_DATE`).Scan(&count, &total)
    return count, total, wrapErr(err, "metrics ventas hoy", 0)
}

func (r *PostgresMetricsRepository) VentasSemana(ctx context.Context) (int, float64, error) {
    // ... WHERE created_at >= date_trunc('week', CURRENT_DATE)
}
func (r *PostgresMetricsRepository) VentasMes(ctx context.Context) (int, float64, error) {
    // ... WHERE created_at >= date_trunc('month', CURRENT_DATE)
}
func (r *PostgresMetricsRepository) TopProductos(ctx context.Context, limit int) ([]ports.TopProduct, error) {
    // ... GROUP BY p.nombre ORDER BY SUM(vi.cantidad) DESC LIMIT $1
}
func (r *PostgresMetricsRepository) StockBajo(ctx context.Context) ([]ports.LowStockProduct, error) {
    // ... WHERE stock_actual <= stock_minimo AND activo = true
}
func (r *PostgresMetricsRepository) ClientesFrecuentes(ctx context.Context, limit int) ([]ports.FrequentClient, error) {
    // ... GROUP BY c.nombre ORDER BY COUNT(*) DESC LIMIT $1
}
```

### Error Wrapping Helper

```go
// wrapErr adds operation context to database errors.
func wrapErr(err error, operation string, entityID int64) error {
    if err == nil {
        return nil
    }
    if entityID > 0 {
        return fmt.Errorf("%s (id=%d): %w", operation, entityID, err)
    }
    return fmt.Errorf("%s: %w", operation, err)
}
```

### `src/infrastructure/adapters/bedrock_query_service.go`

```go
package adapters

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// BedrockConfig holds configuration for the Bedrock adapter.
type BedrockConfig struct {
    ModelID     string  // e.g., "anthropic.claude-3-haiku-20240307-v1:0"
    Region      string  // e.g., "us-east-1"
    MaxTokens   int     // e.g., 300
    Temperature float64 // e.g., 0.1
}

// BedrockQueryService implements ports.AIQueryService using Amazon Bedrock.
type BedrockQueryService struct {
    client *bedrockruntime.Client
    config BedrockConfig
    schema string
}

func NewBedrockQueryService(client *bedrockruntime.Client, cfg BedrockConfig, schema string) *BedrockQueryService {
    return &BedrockQueryService{client: client, config: cfg, schema: schema}
}

// GenerateSQL implements ports.AIQueryService.
func (s *BedrockQueryService) GenerateSQL(ctx context.Context, question string) (string, string, error) {
    // 1. Build Anthropic Messages API payload
    payload := anthropicRequest{
        AnthropicVersion: "bedrock-2023-05-31",
        MaxTokens:        s.config.MaxTokens,
        Temperature:      s.config.Temperature,
        System:           s.buildSystemPrompt(),
        Messages: []anthropicMessage{
            {Role: "user", Content: question},
        },
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return "", "", fmt.Errorf("marshaling bedrock request: %w", err)
    }

    // 2. Call InvokeModel
    resp, err := s.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
        ModelId:     &s.config.ModelID,
        ContentType: strPtr("application/json"),
        Accept:      strPtr("application/json"),
        Body:        body,
    })
    if err != nil {
        return "", "", sanitizeAWSError(err)
    }

    // 3. Parse Anthropic response
    var result anthropicResponse
    if err := json.Unmarshal(resp.Body, &result); err != nil {
        return "", "", ErrAIMalformedResponse
    }

    // 4. Extract and parse NL-SQL JSON from content
    nlResp, err := parseNLSQLContent(result.Content)
    if err != nil {
        return "", "", err
    }

    if nlResp.Error != nil {
        return "", nlResp.Explanation, fmt.Errorf(*nlResp.Error)
    }
    if nlResp.SQL == nil {
        return "", nlResp.Explanation, fmt.Errorf("no SQL generated")
    }

    return *nlResp.SQL, nlResp.Explanation, nil
}

// sanitizeAWSError wraps AWS errors without exposing internal details.
func sanitizeAWSError(err error) error {
    // Map to domain-friendly errors without leaking ARNs, request IDs, etc.
    return ErrAIUnavailable
}
```

#### Anthropic Messages API Types

```go
type anthropicRequest struct {
    AnthropicVersion string             `json:"anthropic_version"`
    MaxTokens        int                `json:"max_tokens"`
    Temperature      float64            `json:"temperature"`
    System           string             `json:"system"`
    Messages         []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type anthropicResponse struct {
    Content []struct {
        Type string `json:"type"`
        Text string `json:"text"`
    } `json:"content"`
    StopReason string `json:"stop_reason"`
}
```

### `src/infrastructure/config/secrets.go`

```go
package config

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"

    "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretsLoader retrieves and caches secrets from AWS Secrets Manager.
type SecretsLoader struct {
    client *secretsmanager.Client
    cache  map[string]string
    mu     sync.RWMutex
}

func NewSecretsLoader(client *secretsmanager.Client) *SecretsLoader {
    return &SecretsLoader{client: client, cache: make(map[string]string)}
}

// GetSecret retrieves a secret value, using cache on subsequent calls.
func (s *SecretsLoader) GetSecret(ctx context.Context, arn string) (string, error) {
    s.mu.RLock()
    if val, ok := s.cache[arn]; ok {
        s.mu.RUnlock()
        return val, nil
    }
    s.mu.RUnlock()

    out, err := s.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
        SecretId: &arn,
    })
    if err != nil {
        return "", fmt.Errorf("retrieving secret %s: %w", arn, err)
    }

    s.mu.Lock()
    s.cache[arn] = *out.SecretString
    s.mu.Unlock()

    return *out.SecretString, nil
}

// LoadConfig loads the full Config from Secrets Manager ARNs.
func (s *SecretsLoader) LoadConfig(ctx context.Context) (*Config, error) {
    // Retrieve DB connection string from SECRET_DB_ARN
    // Retrieve session key from SECRET_SESSION_ARN
    // Retrieve AI config from SECRET_AI_ARN
    // Parse JSON values and build Config struct
}
```

### `internal/bootstrap/bootstrap.go` — Dual-Mode Logic

```go
func BuildRouter(cfg Config) (http.Handler, func(), error) {
    var (
        productRepo   ports.ProductRepository
        saleRepo      ports.SaleRepository
        userRepo      ports.UserRepository
        clientRepo    ports.ClientRepository
        inventoryRepo ports.InventoryRepository
        configRepo    ports.ConfigRepository
        metricsRepo   ports.MetricsRepository
        aiService     ports.AIQueryService
        sessionMgr    *scs.SessionManager
        readDB        *sql.DB  // for NL-SQL query execution
        cleanup       func()
    )

    switch cfg.AppEnv {
    case "lambda":
        // PostgreSQL path
        pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
        // ... configure pool (MaxConns=5, ConnectTimeout=5s)

        productRepo = adapters.NewPostgresProductRepository(pool)
        saleRepo = adapters.NewPostgresSaleRepository(pool)
        userRepo = adapters.NewPostgresUserRepository(pool)
        clientRepo = adapters.NewPostgresClientRepository(pool)
        inventoryRepo = adapters.NewPostgresInventoryRepository(pool)
        configRepo = adapters.NewPostgresConfigRepository(pool)
        metricsRepo = adapters.NewPostgresMetricsRepository(pool)

        // Bedrock AI service
        bedrockClient := bedrockruntime.NewFromConfig(awsCfg)
        aiService = adapters.NewBedrockQueryService(bedrockClient, bedrockCfg, schema)

        // pgx session store
        sessionMgr = newPgxSessionManager(pool)

        // stdlib *sql.DB for NL-SQL read-only queries
        readDB = stdlib.OpenDBFromPool(pool)

        cleanup = func() { pool.Close() }

    default: // "local" or unset
        // SQLite path (existing logic)
        db, err := database.New(cfg.DatabasePath)
        // ... existing adapter initialization
        cleanup = func() { db.Close() }
    }

    // Build router (shared between both modes)
    r := chi.NewRouter()
    // ... middleware, routes (identical to current main.go)

    return r, cleanup, nil
}
```

---

## PostgreSQL Migration — `migrations/postgres/001_init.sql`

```sql
-- 001_init.sql — PostgreSQL equivalent of SQLite 001_init.sql

CREATE TABLE IF NOT EXISTS usuarios (
    id          SERIAL PRIMARY KEY,
    nombre      TEXT NOT NULL,
    pin_hash    TEXT NOT NULL,
    rol         TEXT NOT NULL CHECK(rol IN ('admin', 'cajero')),
    activo      BOOLEAN NOT NULL DEFAULT true,
    failed_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS categorias (
    id          SERIAL PRIMARY KEY,
    nombre      TEXT NOT NULL UNIQUE,
    descripcion TEXT,
    activo      BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS productos (
    id            SERIAL PRIMARY KEY,
    nombre        TEXT NOT NULL,
    sku           TEXT UNIQUE,
    categoria_id  INTEGER REFERENCES categorias(id),
    precio_venta  NUMERIC(12,2) NOT NULL CHECK(precio_venta > 0),
    precio_compra NUMERIC(12,2) NOT NULL DEFAULT 0,
    stock_actual  NUMERIC(12,2) NOT NULL DEFAULT 0,
    stock_minimo  NUMERIC(12,2) NOT NULL DEFAULT 0,
    unidad        TEXT NOT NULL DEFAULT 'unidad'
                  CHECK(unidad IN ('unidad', 'kg', 'litro', 'paquete')),
    activo        BOOLEAN NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS clientes (
    id         SERIAL PRIMARY KEY,
    nombre     TEXT NOT NULL,
    telefono   TEXT,
    direccion  TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ventas (
    id          SERIAL PRIMARY KEY,
    usuario_id  INTEGER NOT NULL REFERENCES usuarios(id),
    cliente_id  INTEGER REFERENCES clientes(id),
    total       NUMERIC(12,2) NOT NULL CHECK(total >= 0),
    metodo_pago TEXT NOT NULL CHECK(metodo_pago IN ('efectivo','tarjeta','transferencia','mixto')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS venta_items (
    id              SERIAL PRIMARY KEY,
    venta_id        INTEGER NOT NULL REFERENCES ventas(id) ON DELETE CASCADE,
    producto_id     INTEGER NOT NULL REFERENCES productos(id),
    cantidad        NUMERIC(12,2) NOT NULL CHECK(cantidad > 0),
    precio_unitario NUMERIC(12,2) NOT NULL,
    subtotal        NUMERIC(12,2) NOT NULL CHECK(subtotal >= 0)
);

CREATE TABLE IF NOT EXISTS inventario_movimientos (
    id               SERIAL PRIMARY KEY,
    producto_id      INTEGER NOT NULL REFERENCES productos(id),
    tipo             TEXT NOT NULL CHECK(tipo IN ('entrada', 'salida', 'ajuste')),
    cantidad         NUMERIC(12,2) NOT NULL CHECK(cantidad > 0),
    stock_resultante NUMERIC(12,2) NOT NULL,
    referencia_tipo  TEXT,
    referencia_id    INTEGER,
    motivo           TEXT,
    usuario_id       INTEGER REFERENCES usuarios(id),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS configuracion (
    id    SERIAL PRIMARY KEY,
    clave TEXT NOT NULL UNIQUE,
    valor TEXT NOT NULL
);

-- Sessions table (pgxstore format)
CREATE TABLE IF NOT EXISTS sessions (
    token  TEXT PRIMARY KEY,
    data   BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

-- Indexes (equivalent to SQLite)
CREATE INDEX IF NOT EXISTS idx_productos_nombre ON productos(nombre);
CREATE INDEX IF NOT EXISTS idx_productos_sku ON productos(sku);
CREATE INDEX IF NOT EXISTS idx_productos_categoria ON productos(categoria_id);
CREATE INDEX IF NOT EXISTS idx_productos_stock ON productos(stock_actual);
CREATE INDEX IF NOT EXISTS idx_ventas_created_at ON ventas(created_at);
CREATE INDEX IF NOT EXISTS idx_ventas_usuario ON ventas(usuario_id);
CREATE INDEX IF NOT EXISTS idx_ventas_cliente ON ventas(cliente_id);
CREATE INDEX IF NOT EXISTS idx_venta_items_venta ON venta_items(venta_id);
CREATE INDEX IF NOT EXISTS idx_venta_items_producto ON venta_items(producto_id);
CREATE INDEX IF NOT EXISTS idx_inventario_producto ON inventario_movimientos(producto_id);
CREATE INDEX IF NOT EXISTS idx_inventario_fecha ON inventario_movimientos(created_at);
CREATE INDEX IF NOT EXISTS idx_inventario_tipo ON inventario_movimientos(tipo);
CREATE INDEX IF NOT EXISTS idx_clientes_nombre ON clientes(nombre);
CREATE INDEX IF NOT EXISTS idx_clientes_telefono ON clientes(telefono);
CREATE INDEX IF NOT EXISTS idx_usuarios_pin ON usuarios(pin_hash);
CREATE INDEX IF NOT EXISTS idx_sessions_expiry ON sessions(expiry);

-- Trigger: update stock on inventory movement
CREATE OR REPLACE FUNCTION fn_inventario_actualiza_stock()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE productos
    SET stock_actual = NEW.stock_resultante,
        updated_at = NOW()
    WHERE id = NEW.producto_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_inventario_actualiza_stock
    AFTER INSERT ON inventario_movimientos
    FOR EACH ROW
    EXECUTE FUNCTION fn_inventario_actualiza_stock();
```

---

## SAM Template Structure — `template.yaml`

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Description: POS AI-First Lambda Deployment

Globals:
  Function:
    Timeout: 30
    MemorySize: 512
    Architectures:
      - arm64
    Environment:
      Variables:
        APP_ENV: lambda

Parameters:
  Environment:
    Type: String
    Default: production
  SecretDbArn:
    Type: String
    Description: ARN of the DB connection secret
  SecretSessionArn:
    Type: String
    Description: ARN of the session secret
  SecretAiArn:
    Type: String
    Description: ARN of the AI config secret
  VpcSubnetIds:
    Type: CommaDelimitedList
    Description: Private subnet IDs for Lambda VPC
  VpcSecurityGroupIds:
    Type: CommaDelimitedList
    Description: Security group IDs for Lambda VPC

Resources:
  PosFunction:
    Type: AWS::Serverless::Function
    Properties:
      PackageType: Image
      ImageUri: !Sub "${AWS::AccountId}.dkr.ecr.${AWS::Region}.amazonaws.com/pos-ai-first:latest"
      Environment:
        Variables:
          SECRET_DB_ARN: !Ref SecretDbArn
          SECRET_SESSION_ARN: !Ref SecretSessionArn
          SECRET_AI_ARN: !Ref SecretAiArn
      VpcConfig:
        SubnetIds: !Ref VpcSubnetIds
        SecurityGroupIds: !Ref VpcSecurityGroupIds
      Policies:
        - Statement:
            - Effect: Allow
              Action:
                - secretsmanager:GetSecretValue
              Resource:
                - !Ref SecretDbArn
                - !Ref SecretSessionArn
                - !Ref SecretAiArn
            - Effect: Allow
              Action:
                - bedrock:InvokeModel
              Resource: "*"
      Events:
        CatchAll:
          Type: HttpApi
          Properties:
            ApiId: !Ref PosHttpApi
            Path: /{proxy+}
            Method: ANY
        Root:
          Type: HttpApi
          Properties:
            ApiId: !Ref PosHttpApi
            Path: /
            Method: ANY

  PosHttpApi:
    Type: AWS::Serverless::HttpApi
    Properties:
      StageName: $default

  StaticAssetsBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: !Sub "pos-static-${Environment}"

  StaticBucketPolicy:
    Type: AWS::S3::BucketPolicy
    Properties:
      Bucket: !Ref StaticAssetsBucket
      PolicyDocument:
        Statement:
          - Effect: Allow
            Principal:
              Service: cloudfront.amazonaws.com
            Action: s3:GetObject
            Resource: !Sub "${StaticAssetsBucket.Arn}/*"
            Condition:
              StringEquals:
                AWS:SourceArn: !Sub "arn:aws:cloudfront::${AWS::AccountId}:distribution/${CloudFrontDistribution}"

  CloudFrontDistribution:
    Type: AWS::CloudFront::Distribution
    Properties:
      DistributionConfig:
        Enabled: true
        Origins:
          - Id: S3Static
            DomainName: !GetAtt StaticAssetsBucket.RegionalDomainName
            S3OriginConfig:
              OriginAccessIdentity: ""
            OriginAccessControlId: !Ref CloudFrontOAC
        DefaultCacheBehavior:
          TargetOriginId: S3Static
          ViewerProtocolPolicy: redirect-to-https
          CachePolicyId: "658327ea-f89d-4fab-a63d-7e88639e58f6"  # CachingOptimized

  CloudFrontOAC:
    Type: AWS::CloudFront::OriginAccessControl
    Properties:
      OriginAccessControlConfig:
        Name: !Sub "pos-static-oac-${Environment}"
        OriginAccessControlOriginType: s3
        SigningBehavior: always
        SigningProtocol: sigv4

Outputs:
  ApiEndpoint:
    Description: API Gateway endpoint URL
    Value: !Sub "https://${PosHttpApi}.execute-api.${AWS::Region}.amazonaws.com"
  FunctionArn:
    Description: Lambda function ARN
    Value: !GetAtt PosFunction.Arn
  StaticBucketName:
    Description: S3 bucket for static assets
    Value: !Ref StaticAssetsBucket
  CloudFrontDomain:
    Description: CloudFront distribution domain
    Value: !GetAtt CloudFrontDistribution.DomainName
```

---

## Dockerfile

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -o /handler cmd/lambda/main.go

FROM public.ecr.aws/lambda/provided:al2-arm64
COPY --from=builder /handler /var/runtime/bootstrap
COPY templates/ /var/task/templates/
WORKDIR /var/task
CMD ["bootstrap"]
```

Key decisions:
- `CGO_ENABLED=0`: Pure Go build (pgx doesn't need CGO; removes modernc.org/sqlite dependency in Lambda build)
- `provided:al2-arm64`: AWS Lambda base image for custom runtime on ARM64
- Templates bundled in image; static files served via CloudFront (not in image)

---

## CI/CD Pipeline — `.github/workflows/deploy.yml`

```yaml
name: Deploy to AWS
on:
  push:
    branches: [main]

permissions:
  id-token: write
  contents: read

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - run: go test ./... -cover
      - run: golangci-lint run

  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: ${{ secrets.AWS_DEPLOY_ROLE_ARN }}
          aws-region: us-east-1

      - uses: aws-actions/amazon-ecr-login@v2
        id: ecr

      - name: Build and push image
        env:
          ECR_REGISTRY: ${{ steps.ecr.outputs.registry }}
          IMAGE_TAG: ${{ github.sha }}
        run: |
          docker build -t $ECR_REGISTRY/pos-ai-first:$IMAGE_TAG .
          docker push $ECR_REGISTRY/pos-ai-first:$IMAGE_TAG
          docker tag $ECR_REGISTRY/pos-ai-first:$IMAGE_TAG $ECR_REGISTRY/pos-ai-first:latest
          docker push $ECR_REGISTRY/pos-ai-first:latest

      - name: SAM Deploy
        run: |
          sam deploy \
            --template-file template.yaml \
            --stack-name pos-ai-first \
            --image-repository ${{ steps.ecr.outputs.registry }}/pos-ai-first \
            --parameter-overrides \
              SecretDbArn=${{ secrets.SECRET_DB_ARN }} \
              SecretSessionArn=${{ secrets.SECRET_SESSION_ARN }} \
              SecretAiArn=${{ secrets.SECRET_AI_ARN }} \
              VpcSubnetIds=${{ secrets.VPC_SUBNET_IDS }} \
              VpcSecurityGroupIds=${{ secrets.VPC_SG_IDS }} \
            --capabilities CAPABILITY_IAM \
            --no-confirm-changeset

      - name: Sync static assets to S3
        run: aws s3 sync static/ s3://pos-static-production/ --delete

      - name: Post-deploy health check
        run: |
          API_URL=$(aws cloudformation describe-stacks \
            --stack-name pos-ai-first \
            --query 'Stacks[0].Outputs[?OutputKey==`ApiEndpoint`].OutputValue' \
            --output text)
          STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/health")
          if [ "$STATUS" != "200" ]; then
            echo "Health check failed with status $STATUS"
            exit 1
          fi
```

---

## Health Endpoint Design

```go
// Enhanced health endpoint for Lambda mode
func healthHandler(pool *pgxpool.Pool, appEnv string) http.HandlerFunc {
    var coldStart = true
    return func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
        defer cancel()

        resp := map[string]interface{}{
            "status":  "ok",
            "app_env": appEnv,
        }

        // Database check
        if err := pool.Ping(ctx); err != nil {
            resp["database"] = "error"
            resp["status"] = "degraded"
            w.WriteHeader(http.StatusServiceUnavailable)
        } else {
            resp["database"] = "ok"
        }

        // Lambda metadata
        if appEnv == "lambda" {
            resp["memory_mb"] = os.Getenv("AWS_LAMBDA_FUNCTION_MEMORY_SIZE")
            resp["cold_start"] = coldStart
            coldStart = false
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    }
}
```

## Structured Logging

```go
package logging

import (
    "encoding/json"
    "log"
    "os"
    "time"
)

type LogEntry struct {
    Timestamp string `json:"timestamp"`
    Level     string `json:"level"`
    RequestID string `json:"request_id,omitempty"`
    Method    string `json:"method,omitempty"`
    Path      string `json:"path,omitempty"`
    Status    int    `json:"status,omitempty"`
    Duration  string `json:"duration,omitempty"`
    Error     string `json:"error,omitempty"`
}

func JSON(entry LogEntry) {
    entry.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
    data, _ := json.Marshal(entry)
    log.Println(string(data))
}
```

## Panic Recovery Middleware

```go
func PanicRecovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                stack := debug.Stack()
                logging.JSON(LogEntry{
                    Level:     "error",
                    RequestID: r.Header.Get("X-Amzn-Trace-Id"),
                    Error:     fmt.Sprintf("panic: %v\n%s", rec, stack),
                })
                http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

---

## Data Models

### PostgreSQL Connection Pool Configuration

```go
poolConfig, _ := pgxpool.ParseConfig(databaseURL)
poolConfig.MaxConns = 5
poolConfig.MinConns = 1
poolConfig.MaxConnLifetime = 30 * time.Minute
poolConfig.MaxConnIdleTime = 5 * time.Minute
poolConfig.ConnConfig.ConnectTimeout = 5 * time.Second
```

### Session Manager (PostgreSQL mode)

```go
func newPgxSessionManager(pool *pgxpool.Pool) *scs.SessionManager {
    sm := scs.New()
    sm.Store = pgxstore.New(pool)
    sm.Lifetime = 8 * time.Hour
    sm.Cookie.Name = "pos_session"
    sm.Cookie.HttpOnly = true
    sm.Cookie.SameSite = http.SameSiteLaxMode
    sm.Cookie.Secure = true // HTTPS via API Gateway
    return sm
}
```

### Config Validation

```go
// ValidateConfig checks that all required fields are present for the given mode.
func ValidateConfig(cfg *Config) error {
    var missing []string
    if cfg.SessionSecret == "" {
        missing = append(missing, "SESSION_SECRET")
    }
    if cfg.AppEnv == "lambda" {
        if cfg.DatabaseURL == "" {
            missing = append(missing, "DATABASE_URL (from SECRET_DB_ARN)")
        }
        if cfg.BedrockModelID == "" {
            missing = append(missing, "BEDROCK_MODEL_ID (from SECRET_AI_ARN)")
        }
    } else {
        if cfg.DatabasePath == "" {
            missing = append(missing, "DATABASE_PATH")
        }
    }
    if len(missing) > 0 {
        return fmt.Errorf("missing required config: %s", strings.Join(missing, ", "))
    }
    return nil
}
```

---

## Error Handling

| Layer | Strategy |
|-------|----------|
| PostgreSQL adapters | Wrap with `wrapErr(err, operation, id)` — context but no internal leak |
| Bedrock adapter | `sanitizeAWSError(err)` — map to domain errors, hide ARNs/request IDs |
| Secrets Manager | Wrap with secret ARN context, log full error, terminate on startup |
| Health endpoint | 2s timeout per check, 503 on any dependency failure |
| Panic recovery | `recover()` + stack trace logging + HTTP 500 |

---

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system — essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: PostgreSQL Repository Round-Trip Persistence

*For any* valid domain entity (Product, User, Client, InventoryMovement, Sale), persisting it via the corresponding PostgreSQL adapter and then retrieving it by ID should return an entity with equivalent field values (excluding server-generated timestamps).

**Validates: Requirements 2.1, 2.2, 2.3, 2.4, 2.5**

### Property 2: Config Repository Upsert Semantics

*For any* key-value pair (key, v1) followed by a second Set with (key, v2), Get(key) should always return v2. Additionally, for any key that has never been Set, Get should return an empty string.

**Validates: Requirements 2.6**

### Property 3: Metrics Aggregation Correctness

*For any* set of sales inserted with `created_at` within today's date range, `VentasHoy()` count should equal the number of those sales and the total should equal the sum of their `total` fields (within floating-point tolerance).

**Validates: Requirements 2.7**

### Property 4: Inventory Trigger Stock Update

*For any* inventory movement inserted into the `inventario_movimientos` table, the corresponding product's `stock_actual` should equal the movement's `stock_resultante` value immediately after insertion.

**Validates: Requirements 3.4**

### Property 5: Bedrock Response Parsing Round-Trip

*For any* valid NLSQLResponse (with non-nil SQL, non-empty explanation, nil error), serializing it as a JSON string wrapped in an Anthropic Messages response format, then parsing it via the Bedrock adapter's response parser, should yield an equivalent NLSQLResponse.

**Validates: Requirements 4.3**

### Property 6: Bedrock Error Sanitization

*For any* AWS error containing internal details (request IDs, ARNs, internal service endpoint URLs), the error returned by `sanitizeAWSError` should not contain those internal strings — it should map to one of the domain-friendly error constants (ErrAIUnavailable, ErrAIRateLimit, ErrAIMalformedResponse).

**Validates: Requirements 4.4**

### Property 7: SQL Validation Consistency Across Adapters

*For any* SQL string, the validation result from the Bedrock adapter should be identical to the result from `nlsql.ValidateSQL` — both should accept or reject the same inputs with equivalent error semantics.

**Validates: Requirements 4.7**

### Property 8: Startup Config Validation Completeness

*For any* Config struct with one or more required fields set to empty string (given the appropriate AppEnv mode), `ValidateConfig` should return a non-nil error whose message contains the name of every missing field.

**Validates: Requirements 8.5**

### Property 9: Health Endpoint Response Structure

*For any* combination of database connectivity state (ok/error) and APP_ENV value, the health endpoint JSON response should always contain the keys: "status", "database", and "app_env". When database connectivity fails, HTTP status should be 503.

**Validates: Requirements 10.1, 10.2, 10.3**

### Property 10: Structured Log Format

*For any* HTTP request processed by the Lambda handler, the emitted log entry should be valid JSON containing at minimum the fields: "timestamp" (RFC3339 format), "level" (non-empty string), and "request_id" (string, possibly empty for local mode).

**Validates: Requirements 10.6**

### Property 11: Panic Recovery Safety

*For any* handler that panics with an arbitrary value, the panic recovery middleware should return HTTP 500 with a JSON error body, and the Lambda process should continue serving subsequent requests without termination.

**Validates: Requirements 10.7**

---

## Testing Strategy

### Property-Based Tests (PBT)

| Property | Generator | Framework |
|----------|-----------|-----------|
| 1: Round-trip persistence | Random valid entities (Product, Sale, User, Client, Movement) | `pgx` + testcontainers-go |
| 2: Config upsert | Random key-value strings | `pgx` + testcontainers-go |
| 3: Metrics aggregation | Random sales with controlled timestamps | `pgx` + testcontainers-go |
| 4: Inventory trigger | Random movements + products | `pgx` + testcontainers-go |
| 5: Bedrock parsing | Random NLSQLResponse structs | Pure Go (no DB) |
| 6: Error sanitization | Random AWS-style error strings | Pure Go (no DB) |
| 7: SQL validation | Random SQL strings (valid and invalid) | Pure Go (no DB) |
| 8: Config validation | Config structs with random missing fields | Pure Go (no DB) |
| 9: Health response | Mocked pool (ping success/failure) | `httptest` |
| 10: Log format | Random request metadata | `httptest` |
| 11: Panic recovery | Random panic values | `httptest` |

### Integration Tests

- Lambda + API Gateway end-to-end (post-deploy health check)
- Secrets Manager retrieval (mocked client)
- Full request flow through algnhsa (use `algnhsa.NewTestRequest`)

### Unit Tests (Example-Based)

- Dual-mode startup with APP_ENV="local" vs "lambda"
- Bedrock request payload format (Anthropic Messages API structure)
- PostgreSQL error wrapping format
- Session manager configuration per mode

---

## Dependencies (New)

| Package | Purpose | Justification |
|---------|---------|---------------|
| `github.com/jackc/pgx/v5` | PostgreSQL driver + pool | Pure Go, no CGO, excellent performance, pgxpool |
| `github.com/akrylysov/algnhsa` | Lambda-chi adapter | Minimal bridge, well-maintained, chi-compatible |
| `github.com/aws/aws-sdk-go-v2` | AWS SDK (Bedrock, SM) | Official AWS Go SDK v2, modular imports |
| `github.com/alexedwards/scs/pgxstore` | PostgreSQL session store | Same scs ecosystem as current SQLite store |
| `github.com/jackc/pgx/v5/stdlib` | pgx → database/sql bridge | Needed for NL-SQL read-only query execution via existing *sql.DB interface |

---

## Key Design Decisions

1. **Separate binaries** (`cmd/server/` vs `cmd/lambda/`): Keeps build artifacts clean. Lambda binary excludes SQLite CGO dependency entirely.
2. **Shared bootstrap package**: Avoids duplicating router/DI logic. Both entry points call `BuildRouter()`.
3. **No changes to domain/application**: The hexagonal boundary holds. Only infrastructure adapters change.
4. **pgx native (not database/sql)**: Better performance, native PostgreSQL types, but we bridge to `*sql.DB` via stdlib for the NL-SQL query execution path (which needs generic row scanning).
5. **Container image Lambda**: Avoids ZIP size limits, includes templates, supports custom binary.
6. **ARM64**: Cost-effective (20% cheaper than x86), excellent Go performance.
7. **VPC placement**: Required for RDS connectivity. Bedrock and Secrets Manager accessed via VPC endpoints or NAT Gateway.
8. **Secrets cached in memory**: Avoid repeated API calls on warm Lambda starts. Cache lives for the execution environment lifetime.
