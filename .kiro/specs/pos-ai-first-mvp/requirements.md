# Requirements: POS AI-First MVP

## Overview
MVP de un sistema POS (Point of Sale) AI-First para el bootcamp Código Facilito + Kiro. Un dueño de negocio puede preguntarle a su POS en lenguaje natural "¿qué vendí esta semana?" y obtener respuesta inmediata. El sistema incluye CRUD de productos/ventas, dashboard de métricas, y chat conversacional con NL→SQL.

## Timeline
5 días de desarrollo (bootcamp)

## Evaluation Criteria
- Impacto tecnológico (30%)
- Innovación (30%)
- Software funcional y entregables (30%)
- Uso de Servicios de AWS y Kiro (10%)

---

## Functional Requirements

### REQ-1: Project Foundation
**Description:** Go project scaffolded with hexagonal architecture, SQLite database, chi router, and all tooling configured.
**Acceptance Criteria:**
- GIVEN a fresh clone of the repo WHEN I run `make run` THEN the server starts on port 8080
- GIVEN the project structure WHEN I inspect `src/` THEN it follows hexagonal layout (domain/application/infrastructure)
- GIVEN the database WHEN the server starts THEN migrations run automatically creating all 8 tables
- GIVEN the project WHEN I run `make test` THEN tests execute successfully
- GIVEN the project WHEN I run `make lint` THEN golangci-lint runs without errors

### REQ-2: PIN Authentication
**Description:** Multi-user authentication via numeric PIN with role-based access (admin/cajero).
**Acceptance Criteria:**
- GIVEN a user with a valid PIN WHEN they enter it THEN they are authenticated and a session is created
- GIVEN an invalid PIN WHEN entered 5 times THEN the account is locked for 5 minutes
- GIVEN a locked account WHEN the lockout expires THEN the user can try again
- GIVEN an authenticated session WHEN 8 hours pass THEN the session expires automatically
- GIVEN a cajero role WHEN they access admin-only features THEN access is denied

### REQ-3: Product CRUD
**Description:** Full product management with categories, pricing, stock tracking, and SKU support.
**Acceptance Criteria:**
- GIVEN valid product data WHEN I create a product THEN it is stored with all fields including category, price, stock, SKU
- GIVEN an existing product WHEN I update its price THEN the price changes and updated_at is refreshed
- GIVEN a product with sales history WHEN I deactivate it THEN it remains in the database but is hidden from POS
- GIVEN the product list WHEN I filter by category THEN only matching products appear
- GIVEN a product WHEN its stock reaches stock_minimo THEN it appears in low-stock alerts

### REQ-4: Sales Registration
**Description:** Register sales with multiple items, automatic stock deduction, and payment methods.
**Acceptance Criteria:**
- GIVEN a cart with items WHEN I complete a sale THEN total is calculated, stock is deducted, and an inventory movement is created
- GIVEN insufficient stock WHEN I try to sell THEN the sale is rejected with a clear message
- GIVEN a valid sale WHEN payment method is selected (efectivo/tarjeta/transferencia/mixto) THEN the sale is recorded with that method
- GIVEN a completed sale WHEN I check inventory THEN stock_actual reflects the deduction
- GIVEN a sale with multiple items WHEN completed THEN each item generates its own venta_item record

### REQ-5: NL→SQL Chat (AI Core Feature)
**Description:** Conversational interface where users ask questions in natural language (Spanish) and get answers from their data via OpenRouter API generating SQL queries.
**Acceptance Criteria:**
- GIVEN a user question "¿qué vendí hoy?" WHEN submitted THEN the system calls OpenRouter, generates a SELECT query, executes it, and returns a formatted response in Spanish
- GIVEN a generated SQL WHEN it contains INSERT/UPDATE/DELETE/DROP THEN it is rejected before execution
- GIVEN a generated SQL WHEN executed THEN it runs with a 5-second timeout and LIMIT 500
- GIVEN a valid query WHEN executed THEN results are formatted as a human-readable response in the chat
- GIVEN an invalid or unsafe query WHEN generated THEN the user sees a friendly error message, not raw SQL errors
- GIVEN OpenRouter is unavailable WHEN a question is asked THEN a fallback message is shown

### REQ-6: Dashboard Metrics
**Description:** Fixed metrics dashboard showing sales today/week/month, top products, low stock alerts, and frequent customers.
**Acceptance Criteria:**
- GIVEN the dashboard page WHEN loaded THEN it shows ventas hoy, ventas semana, ventas mes as KPI cards
- GIVEN the dashboard WHEN loaded THEN top 5 products sold (last 30 days) are displayed
- GIVEN products with stock ≤ stock_minimo WHEN dashboard loads THEN they appear in the low-stock alert section
- GIVEN the dashboard WHEN 30 seconds pass THEN ventas-hoy refreshes automatically via HTMX polling
- GIVEN no sales data WHEN dashboard loads THEN appropriate empty states are shown

### REQ-7: UI Layout
**Description:** Split-view layout with traditional POS app (CRUD + dashboard) on top and always-visible chat bar on bottom.
**Acceptance Criteria:**
- GIVEN any page WHEN rendered THEN the chat bar is visible at the bottom of the screen
- GIVEN the layout WHEN viewed on desktop THEN navigation shows sidebar with sections: Dashboard, Productos, Ventas, Inventario
- GIVEN the chat bar WHEN a message is submitted THEN the response appears in the chat area without page reload
- GIVEN the UI WHEN rendered THEN Tailwind CSS styles are applied correctly via CDN

---

## Non-Functional Requirements

### NFR-1: Security
- All generated SQL must be validated (SELECT-only) before execution
- PINs stored as bcrypt hashes, never plaintext
- SQLite queries from NL→SQL execute on a read-only connection
- All NL→SQL queries are logged for audit

### NFR-2: Performance
- Dashboard metrics load in < 500ms
- NL→SQL response in < 8 seconds (including API call)
- HTMX fragments render in < 100ms server-side

### NFR-3: Reliability
- OpenRouter failures don't crash the app (graceful degradation)
- SQLite WAL mode for concurrent reads during writes
- Server starts and runs without external dependencies (besides OpenRouter API key)

### NFR-4: Maintainability
- Hexagonal architecture: domain has zero external imports
- All critical paths have tests (auth, inventory, NL→SQL validation)
- Code passes golangci-lint with zero warnings
