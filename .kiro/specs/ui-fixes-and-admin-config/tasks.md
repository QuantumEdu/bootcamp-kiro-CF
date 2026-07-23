# Implementation Plan: UI Fixes and Admin Config

## Overview

This plan implements three groups of changes: (1) Quick UI fixes — logout button in sidebar, product creation button, HTMX cache headers; (2) Client CRUD — domain entity, repository, use cases, handler, templates, and routing; (3) Admin API key configuration — crypto service, config repository, handler, templates, and admin-only routing. All new backend components follow the existing hexagonal architecture.

## Tasks

- [x] 1. Quick UI fixes (logout, product button, HTMX cache)
  - [x] 1.1 Add logout button and "Clientes" nav link to sidebar in layout.html
    - Modify `templates/layout.html` sidebar footer section to include a POST /logout form with a "Cerrar Sesión" button below the user info
    - Add a "Clientes" nav link in the `<nav>` section pointing to `/clientes`
    - Add a conditional "Configuración" nav link visible only when `UserRole` is "admin" pointing to `/admin/config`
    - _Requirements: 1.1, 1.2, 1.4, 3.3, 6.9_

  - [x] 1.2 Add "Nuevo Producto" button in products page template
    - Modify `templates/products/list.html` (or the content template used by `pageHandler.Products`) to include an anchor/button linking to `/productos/new`
    - Ensure it is visible to all authenticated users
    - _Requirements: 2.1, 2.2_

  - [x] 1.3 Add NoCache middleware for HTMX fragment routes
    - Create `src/infrastructure/http/middleware/nocache.go` with a `NoCache` handler that sets `Cache-Control: no-store`
    - Apply the middleware to the HTMX fragment route group (`/api/metrics/*`, `/api/productos/*`) in `cmd/server/main.go`
    - _Requirements: 5.2, 5.3_

  - [ ]* 1.4 Write property test for NoCache middleware
    - **Property 5: HTMX fragment responses include no-store cache header**
    - **Validates: Requirements 5.3**

- [x] 2. Checkpoint — Verify quick fixes compile
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 3. Client domain and infrastructure layer
  - [x] 3.1 Create Client entity with validation
    - Create `src/domain/entities/client.go` with `Client` struct (ID, Nombre, Telefono, Direccion, CreatedAt)
    - Implement `Validate()` method that rejects empty/whitespace-only `Nombre`
    - Define `ErrClientNameRequired` sentinel error
    - _Requirements: 4.2, 4.4_

  - [ ]* 3.2 Write property test for Client.Validate whitespace rejection
    - **Property 2: Whitespace-only client names are rejected**
    - **Validates: Requirements 4.4**

  - [x] 3.3 Create ClientRepository port interface
    - Create `src/domain/ports/client_repository.go` with `ClientRepository` interface (Create, List methods)
    - _Requirements: 3.1, 4.3_

  - [x] 3.4 Create SQLiteClientRepository adapter
    - Create `src/infrastructure/adapters/sqlite_client_repository.go` implementing `ClientRepository`
    - `Create` inserts into `clientes` table and sets the client ID from `LastInsertId`
    - `List` returns all clients ordered by `nombre ASC`
    - _Requirements: 3.2, 4.3_

  - [x] 3.5 Write unit tests for Client entity validation
    - Table-driven tests: valid name, empty string, whitespace-only, trimmed spaces
    - _Requirements: 4.4_

- [ ] 4. Client application layer and handler
  - [x] 4.1 Create Client use cases
    - Create `src/application/use_cases/client_use_cases.go` with `CreateClient` and `ListClients` use cases
    - `CreateClient` validates the client entity before calling the repository
    - _Requirements: 4.3, 4.4_

  - [x] 4.2 Create ClientHandler with List, CreateForm, Create actions
    - Create `src/infrastructure/http/handlers/clients.go`
    - `List`: renders client list page (full page or HTMX fragment)
    - `CreateForm`: renders the new client form template
    - `Create`: parses form, calls CreateClient UC, redirects on success, returns 422 on validation error
    - _Requirements: 3.1, 3.2, 3.4, 4.1, 4.2, 4.3, 4.4_

  - [x] 4.3 Create client templates (list and form)
    - Create `templates/clients/list.html` with client table (nombre, telefono, direccion), empty state message, and "Nuevo Cliente" button
    - Create `templates/clients/form.html` with fields for nombre (required), telefono, direccion
    - _Requirements: 3.2, 3.4, 4.1, 4.2, 4.5_

  - [x] 4.4 Wire client routes in cmd/server/main.go
    - Instantiate `SQLiteClientRepository`, `CreateClient` UC, `ListClients` UC, `ClientHandler`
    - Register routes: `GET /clientes`, `GET /clientes/new`, `POST /clientes` inside the protected group
    - _Requirements: 3.1, 4.1_

  - [x] 4.5 Write unit tests for ClientHandler
    - Test GET /clientes returns 200 with seeded data
    - Test POST /clientes with valid name redirects to /clientes
    - Test POST /clientes with empty name returns 422
    - _Requirements: 3.2, 4.3, 4.4_

- [x] 5. Checkpoint — Verify client CRUD compiles and tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 6. Admin config — crypto, repository, handler
  - [x] 6.1 Create CryptoService for AES-GCM encryption
    - Create `src/application/services/crypto.go` with `CryptoService` struct
    - `NewCryptoService(secret)` derives a 256-bit key from `SESSION_SECRET` via SHA-256
    - `Encrypt(plaintext)` returns hex-encoded AES-GCM ciphertext
    - `Decrypt(ciphertextHex)` returns the original plaintext
    - _Requirements: 6.6, 6.7_

  - [ ]* 6.2 Write property test for AES-GCM round-trip
    - **Property 1: AES-GCM encryption round-trip**
    - **Validates: Requirements 6.6, 6.7**

  - [x] 6.3 Create MaskAPIKey helper and ConfigRepository port
    - Create `src/domain/ports/config_repository.go` with `ConfigRepository` interface (Get, Set)
    - Add a `MaskAPIKey` helper function (shows last 4 chars, masks the rest with asterisks)
    - _Requirements: 6.3, 6.4_

  - [ ]* 6.4 Write property test for API key masking
    - **Property 3: API key masking reveals only last four characters**
    - **Validates: Requirements 6.3**

  - [x] 6.5 Create SQLiteConfigRepository adapter
    - Create `src/infrastructure/adapters/sqlite_config_repository.go` implementing `ConfigRepository`
    - `Get` returns value for a key, empty string if not found
    - `Set` uses INSERT ... ON CONFLICT DO UPDATE (upsert)
    - _Requirements: 6.6_

  - [x] 6.6 Create AdminConfigHandler with Show and Update actions
    - Create `src/infrastructure/http/handlers/admin_config.go`
    - `Show`: reads encrypted key from repo, decrypts, masks, renders template
    - `Update`: validates non-empty key, encrypts, stores, redirects
    - _Requirements: 6.1, 6.3, 6.4, 6.5, 6.6_

  - [x] 6.7 Create admin config template
    - Create `templates/admin/config.html` with masked key display, empty state prompt, and update form
    - _Requirements: 6.3, 6.4, 6.5_

  - [x] 6.8 Wire admin routes in cmd/server/main.go
    - Add SESSION_SECRET validation at startup (log.Fatalf if empty)
    - Instantiate `CryptoService`, `SQLiteConfigRepository`, `AdminConfigHandler`
    - Register admin-only route group with `RequireRole("admin")`: `GET /admin/config`, `POST /admin/config`
    - _Requirements: 6.1, 6.2, 6.8_

  - [x] 6.9 Write property test for role enforcement
    - **Property 4: Non-admin users cannot access admin routes**
    - **Validates: Requirements 6.2**

- [x] 7. Final checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- The `clientes` and `configuracion` tables already exist in the DB schema (001_init.sql) — no migration needed
- The `POST /logout` route already exists and is wired — only a UI button in the sidebar is needed
- The `GET /productos/new` and `POST /productos` routes already exist — only a "Nuevo Producto" button in the template is needed
- Property tests validate universal correctness properties from the design document
- The `pages.go` handler currently renders the products page — check how it injects `content` template to determine where the button goes

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1", "1.2", "1.3", "3.1", "3.3"] },
    { "id": 1, "tasks": ["1.4", "3.2", "3.4", "3.5", "6.1", "6.3"] },
    { "id": 2, "tasks": ["4.1", "6.2", "6.4", "6.5"] },
    { "id": 3, "tasks": ["4.2", "4.3", "6.6", "6.7"] },
    { "id": 4, "tasks": ["4.4", "4.5", "6.8", "6.9"] }
  ]
}
```
