# POS Domain Guide — Business Rules & Invariants

## Domain entities

### Product (`productos`)
- `precio_venta` must be > 0
- `stock_actual` must be >= 0 (never negative)
- `stock_minimo` is the low-stock alert threshold
- Deactivating a product (`activo = false`) does NOT delete it — preserves sales history
- SKU is optional but must be unique if provided

### Sale (`ventas`)
A sale is complete when:
1. All items have sufficient stock (`stock_actual >= cantidad`)
2. Total matches sum of `venta_items.subtotal`
3. Each item's stock is decremented atomically
4. An `inventario_movimientos` record is created per item (type: `venta`)

**Business invariant:** A sale that would cause negative stock MUST be rejected before any DB write.

### User / Authentication (`usuarios`)
- PIN is stored as bcrypt hash — NEVER plain text
- Roles: `admin` (full access) or `cajero` (POS only, no admin config)
- Failed PIN attempts are tracked — lockout after 5 attempts for 5 minutes
- Sessions expire after 8 hours (configurable via `SESSION_LIFETIME_HOURS`)

### Client (`clientes`)
- Name is required (min 2 chars)
- Phone is optional (10 digits if provided)
- Clients are linked to sales via `ventas.cliente_id` (nullable — anonymous sales allowed)

### Inventory movement (`inventario_movimientos`)
Types:
- `entrada` — stock added (purchase, manual adjustment)
- `venta` — stock decremented by a sale
- `ajuste` — manual correction

Every sale MUST generate one `inventario_movimientos` record per line item.

### Configuration (`configuracion`)
Key-value store for runtime settings:
- `openrouter_api_key` — encrypted with AES-GCM using `SESSION_SECRET` as key material
- Access is admin-only via `/admin/config`

## NL→SQL allowed tables

The AI query service may only access these tables:
- `productos`
- `ventas`
- `venta_items`
- `clientes`
- `inventario_movimientos`
- `categorias`

**Never expose:** `usuarios` (PIN hashes), `configuracion` (API keys).

## Payment methods (`metodo_pago`)

Valid values: `efectivo`, `tarjeta`, `transferencia`

## Category names (standard set)

`comida`, `bebida`, `postre`, `botana`, `limpieza`

## Stock update flow

```
RegisterSale use case
  → for each item: check stock_actual >= cantidad
  → if any fails: return ErrInsufficientStock (no partial writes)
  → begin transaction
  → insert venta
  → for each item: insert venta_item, update stock_actual, insert inventario_movimiento
  → commit
```

## Hexagonal architecture constraint

**The domain package (`src/domain/`) MUST NOT import:**
- Any HTTP framework (chi, net/http handlers)
- Any DB driver (sqlite, pgx)
- Any external service (OpenRouter, Bedrock, AWS SDK)
- Any configuration loader (godotenv, os.Getenv)

Domain entities and use cases receive everything they need via **port interfaces** defined in `src/domain/ports/`.
