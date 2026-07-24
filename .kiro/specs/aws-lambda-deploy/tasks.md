# Implementation Plan: AWS Lambda Deploy

## Overview

Deploy the POS AI-First Go application to AWS Lambda using container images (ARM64). The implementation proceeds in waves: PostgreSQL migrations and helpers first, then adapters, AWS integrations, the shared bootstrap package, Lambda entry point, IaC, CI/CD, and finally health/monitoring. Domain and application layers remain untouched — only infrastructure changes.

## Tasks

- [x] 1. PostgreSQL migrations and shared helpers
  - [x] 1.1 Create PostgreSQL migration file `migrations/postgres/001_init.sql`
    - Translate all SQLite DDL to PostgreSQL-native types (SERIAL, NUMERIC(12,2), TIMESTAMPTZ, BOOLEAN)
    - Include all tables: usuarios, categorias, productos, clientes, ventas, venta_items, inventario_movimientos, configuracion, sessions
    - Create all indexes equivalent to SQLite migrations
    - Create `fn_inventario_actualiza_stock` PL/pgSQL trigger function and trigger
    - Create sessions table compatible with scs pgxstore format
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6_

  - [x] 1.2 Create error wrapping helper `src/infrastructure/adapters/postgres_helpers.go`
    - Implement `wrapErr(err error, operation string, entityID int64) error` function
    - Return nil for nil errors, include entity ID context when > 0
    - _Requirements: 2.8_

- [x] 2. PostgreSQL repository adapters
  - [x] 2.1 Implement `src/infrastructure/adapters/postgres_product_repository.go`
    - Implement `PostgresProductRepository` struct with `*pgxpool.Pool` dependency
    - Implement Create, Update, FindByID, List, Deactivate, FindLowStock methods
    - Use parameterized queries ($1, $2, ...) for all operations
    - _Requirements: 2.1, 2.9, 2.10_

  - [x] 2.2 Implement `src/infrastructure/adapters/postgres_sale_repository.go`
    - Implement `PostgresSaleRepository` struct with `*pgxpool.Pool` dependency
    - Use `pool.Begin(ctx)` transaction for Create (insert header + items atomically)
    - _Requirements: 2.2, 2.9, 2.10_

  - [x] 2.3 Implement `src/infrastructure/adapters/postgres_user_repository.go`
    - Implement FindByID, FindByPINHash, FindAll, IncrementFailedAttempts, Lock, ResetAttempts
    - _Requirements: 2.3, 2.9, 2.10_

  - [x] 2.4 Implement `src/infrastructure/adapters/postgres_client_repository.go`
    - Implement Create and List operations
    - _Requirements: 2.4, 2.9, 2.10_

  - [x] 2.5 Implement `src/infrastructure/adapters/postgres_inventory_repository.go`
    - Implement Create and FindByProduct operations
    - _Requirements: 2.5, 2.9, 2.10_

  - [x] 2.6 Implement `src/infrastructure/adapters/postgres_config_repository.go`
    - Implement Get (return empty string for missing keys) and Set (upsert with ON CONFLICT)
    - _Requirements: 2.6, 2.9, 2.10_

  - [x] 2.7 Implement `src/infrastructure/adapters/postgres_metrics_repository.go`
    - Implement VentasHoy, VentasSemana, VentasMes using PostgreSQL date_trunc functions
    - Implement TopProductos, StockBajo, ClientesFrecuentes with aggregation queries
    - _Requirements: 2.7, 2.9, 2.10_

  - [x] 2.8 Write property tests for PostgreSQL round-trip persistence
    - **Property 1: PostgreSQL Repository Round-Trip Persistence**
    - Test with random valid entities (Product, Sale, User, Client, InventoryMovement)
    - Use testcontainers-go for PostgreSQL instance
    - **Validates: Requirements 2.1, 2.2, 2.3, 2.4, 2.5**

  - [x] 2.9 Write property test for Config repository upsert semantics
    - **Property 2: Config Repository Upsert Semantics**
    - Test that Set(key, v2) overwrites Set(key, v1) and Get returns v2
    - Test that Get on never-Set keys returns empty string
    - **Validates: Requirements 2.6**

  - [x] 2.10 Write property test for metrics aggregation
    - **Property 3: Metrics Aggregation Correctness**
    - Insert random sales with controlled timestamps, verify VentasHoy count and total
    - **Validates: Requirements 2.7**

  - [x]* 2.11 Write property test for inventory trigger stock update
    - **Property 4: Inventory Trigger Stock Update**
    - Insert random movements, verify product stock_actual equals movement stock_resultante
    - **Validates: Requirements 3.4**

- [x] 3. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 4. Bedrock adapter and Secrets Manager
  - [x] 4.1 Implement `src/infrastructure/adapters/bedrock_query_service.go`
    - Implement BedrockQueryService struct with bedrockruntime.Client + BedrockConfig
    - Implement GenerateSQL using Anthropic Messages API format
    - Implement sanitizeAWSError to map AWS errors to domain-friendly errors (ErrAIUnavailable, ErrAIRateLimit, ErrAIMalformedResponse)
    - Implement parseNLSQLContent to extract SQL and explanation from Anthropic response
    - Reuse existing system prompt logic from OpenRouter adapter
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7_

  - [x] 4.2 Implement `src/infrastructure/config/secrets.go`
    - Implement SecretsLoader struct with secretsmanager.Client + in-memory cache
    - Implement GetSecret with read-through caching (RWMutex-protected)
    - Implement LoadConfig to retrieve DB URL, session key, and AI config from secret ARNs
    - Return descriptive error on retrieval failure for startup termination
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.7_

  - [x]* 4.3 Write property test for Bedrock response parsing
    - **Property 5: Bedrock Response Parsing Round-Trip**
    - Generate random valid NLSQLResponse, wrap in Anthropic Messages format, parse back
    - **Validates: Requirements 4.3**

  - [x] 4.4 Write property test for Bedrock error sanitization
    - **Property 6: Bedrock Error Sanitization**
    - Generate random AWS errors with internal details, verify sanitized output has no ARNs/request IDs
    - **Validates: Requirements 4.4**

  - [x] 4.5 Write property test for SQL validation consistency
    - **Property 7: SQL Validation Consistency Across Adapters**
    - Generate random SQL strings, verify Bedrock validation matches nlsql.ValidateSQL results
    - **Validates: Requirements 4.7**

- [x] 5. Shared bootstrap package (dual-mode)
  - [x] 5.1 Create `internal/bootstrap/bootstrap.go`
    - Define Config struct with all required fields (AppEnv, Port, DatabaseURL, DatabasePath, SessionSecret, BedrockModelID, etc.)
    - Implement BuildRouter(cfg Config) returning (http.Handler, func(), error)
    - Switch on cfg.AppEnv: "lambda" → PostgreSQL + Bedrock path; default → SQLite + OpenRouter path
    - Configure pgxpool with MaxConns=5, ConnectTimeout=5s for Lambda mode
    - Create pgx session manager for Lambda, SQLite session manager for local
    - Build chi router with all middleware and routes (same as current cmd/server/main.go)
    - _Requirements: 1.3, 1.4, 8.1, 8.2, 8.3, 8.4, 8.5_

  - [x] 5.2 Implement config validation in `internal/bootstrap/config.go`
    - Implement ValidateConfig checking all required fields based on AppEnv mode
    - Lambda mode: require DatabaseURL, SessionSecret, BedrockModelID
    - Local mode: require DatabasePath, SessionSecret
    - Return error listing all missing fields
    - _Requirements: 8.5_

  - [x] 5.3 Implement structured logging in `internal/bootstrap/logging.go`
    - Implement LogEntry struct with timestamp, level, request_id, method, path, status, duration, error
    - Implement JSON logger that writes structured entries to stdout
    - _Requirements: 10.6_

  - [x] 5.4 Implement panic recovery middleware in `internal/bootstrap/middleware.go`
    - Implement PanicRecovery middleware: recover panics, log with stack trace, return HTTP 500 JSON
    - Extract X-Amzn-Trace-Id for request_id in Lambda context
    - _Requirements: 10.7_

  - [x]* 5.5 Write property test for config validation completeness
    - **Property 8: Startup Config Validation Completeness**
    - Generate Config structs with random missing fields, verify error mentions all missing field names
    - **Validates: Requirements 8.5**

  - [x]* 5.6 Write property test for structured log format
    - **Property 10: Structured Log Format**
    - Generate random request metadata, verify emitted log is valid JSON with required fields
    - **Validates: Requirements 10.6**

  - [x]* 5.7 Write property test for panic recovery
    - **Property 11: Panic Recovery Safety**
    - Generate random panic values, verify middleware returns 500 JSON and doesn't terminate
    - **Validates: Requirements 10.7**

- [x] 6. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 7. Lambda entry point and health endpoint
  - [x] 7.1 Create `cmd/lambda/main.go`
    - Import bootstrap package and algnhsa
    - Load config from Secrets Manager when APP_ENV=lambda
    - Call bootstrap.BuildRouter with loaded config
    - Start algnhsa.ListenAndServe(router, nil)
    - Terminate with log.Fatalf on bootstrap failure
    - _Requirements: 1.1, 1.2, 1.3, 1.6, 5.5_

  - [x] 7.2 Implement health endpoint handler in `internal/bootstrap/health.go`
    - Implement healthHandler accepting pool and appEnv
    - Use 2-second timeout for dependency checks (pool.Ping)
    - Return JSON with status, database, app_env fields
    - Return 503 when database ping fails
    - Include memory_mb and cold_start metadata in Lambda mode
    - Register on GET /health route in BuildRouter
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_

  - [x] 7.3 Write property test for health endpoint response structure
    - **Property 9: Health Endpoint Response Structure**
    - Test with mocked pool (ping ok/error) and various APP_ENV values
    - Verify response always contains "status", "database", "app_env" keys
    - Verify 503 status when database fails
    - **Validates: Requirements 10.1, 10.2, 10.3**

- [x] 8. Dockerfile and static asset routing
  - [x] 8.1 Create `Dockerfile`
    - Multi-stage build: golang:1.22-alpine builder → public.ecr.aws/lambda/provided:al2-arm64
    - CGO_ENABLED=0 GOARCH=arm64 GOOS=linux build of cmd/lambda/main.go
    - Copy binary as /var/runtime/bootstrap
    - Copy templates/ directory into image
    - Do NOT include static/ (served via CloudFront)
    - _Requirements: 1.6_

  - [x] 8.2 Update static file routing for dual-mode in bootstrap
    - In Lambda mode: exclude /static/* route (CloudFront handles it)
    - In local mode: keep existing http.FileServer for /static/
    - _Requirements: 9.4, 9.5_

- [x] 9. SAM template (Infrastructure as Code)
  - [x] 9.1 Create `template.yaml` SAM template
    - Define AWS::Serverless::Function with container image package type, ARM64, 512MB, 30s timeout
    - Define AWS::Serverless::HttpApi with catch-all route (/{proxy+} and /)
    - Define Parameters: Environment, SecretDbArn, SecretSessionArn, SecretAiArn, VpcSubnetIds, VpcSecurityGroupIds
    - Configure VpcConfig with private subnets for RDS connectivity
    - Grant IAM permissions: secretsmanager:GetSecretValue (specific ARNs) + bedrock:InvokeModel
    - Define S3 bucket for static assets with CloudFront OAC policy
    - Define CloudFront distribution with S3 origin and CachingOptimized policy
    - Define Outputs: ApiEndpoint, FunctionArn, StaticBucketName, CloudFrontDomain
    - Include Globals section with APP_ENV=lambda and arm64 architecture
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.7, 6.8, 9.1, 9.2, 9.6_

- [x] 10. CI/CD pipeline
  - [x] 10.1 Create `.github/workflows/deploy.yml`
    - Trigger on push to main branch
    - Test job: go test ./... -cover + golangci-lint run (abort on failure)
    - Deploy job (needs test): OIDC auth via aws-actions/configure-aws-credentials
    - Build + push container image to ECR (tagged with SHA + latest)
    - Run sam deploy with parameter overrides from secrets
    - Sync static/ to S3 bucket
    - Post-deploy health check: curl /health and verify HTTP 200
    - Report failure if health check returns non-200
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7, 7.8, 9.3_

- [x] 11. Refactor cmd/server/main.go to use bootstrap package
  - [x] 11.1 Refactor `cmd/server/main.go` to call bootstrap.BuildRouter
    - Replace inline DI and router setup with bootstrap.BuildRouter(cfg) call
    - Load config from env vars / .env file for local mode
    - Keep net/http.ListenAndServe for local startup
    - Ensure existing behavior is preserved (all routes, middleware, session management)
    - _Requirements: 1.4, 8.1, 8.3_

- [ ] 12. Final checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties from the design document
- Unit tests validate specific examples and edge cases
- The implementation language is Go (already used in the project)
- Domain and application layers remain completely untouched — only infrastructure changes
- Property tests requiring PostgreSQL should use testcontainers-go; pure logic tests use httptest/in-memory mocks
- The `internal/bootstrap/` package is the key architectural piece enabling dual-mode operation

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1", "1.2"] },
    { "id": 1, "tasks": ["2.1", "2.2", "2.3", "2.4", "2.5", "2.6", "2.7"] },
    { "id": 2, "tasks": ["2.8", "2.9", "2.10", "2.11", "4.1", "4.2"] },
    { "id": 3, "tasks": ["4.3", "4.4", "4.5", "5.1", "5.2", "5.3", "5.4"] },
    { "id": 4, "tasks": ["5.5", "5.6", "5.7", "7.1", "7.2"] },
    { "id": 5, "tasks": ["7.3", "8.1", "8.2", "11.1"] },
    { "id": 6, "tasks": ["9.1", "10.1"] }
  ]
}
```
