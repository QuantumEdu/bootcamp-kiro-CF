# 001 — Modelo de Datos del MVP

**Estado:** Completado
**Fecha:** 2026-07-21
**Stack:** SQLite + sqlc

---

## Resumen

Schema relacional para un POS AI-First. Cinco dominios: productos, ventas, inventario, clientes, usuarios. Diseñado para ser consumido tanto por la app Go como por el generador NL→SQL vía OpenRouter.

Principios:
- **Pragmático para MVP** — nada de tablas de auditoría, soft-delete, ni normalización excesiva
- **Nombres en castellano** — tablas y columnas consultables directamente por el prompt NL
- **Tipos SQLite nativos** — INTEGER, REAL, TEXT. Sin datetime, boolean, ni decimal
- **sqlc-ready** — DDL en un solo archivo, consultas parametrizables

---

## DDL Completo

```sql
-- 001-init.sql — Schema MVP del POS AI-First
-- SQLite + sqlc ready
-- Compatible con PRAGMA journal_mode=WAL

-- ============================================================
-- USUARIOS: cajeros y administradores del sistema
-- ============================================================
CREATE TABLE IF NOT EXISTS usuarios (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre          TEXT NOT NULL,                          -- nombre completo del usuario
    pin_hash        TEXT NOT NULL,                          -- PIN hasheado con bcrypt
    rol             TEXT NOT NULL CHECK(rol IN ('admin', 'cajero')),  -- admin: todo acceso, cajero: ventas + consultas
    activo          INTEGER NOT NULL DEFAULT 1,             -- 1=activo, 0=desactivado (no se puede eliminar)
    created_at      TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- ============================================================
-- CATEGORIAS: agrupación de productos
-- ============================================================
CREATE TABLE IF NOT EXISTS categorias (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre          TEXT NOT NULL UNIQUE,                   -- ej: "Bebidas", "Lácteos", "Limpieza"
    descripcion     TEXT,                                   -- opcional
    activo          INTEGER NOT NULL DEFAULT 1
);

-- ============================================================
-- PRODUCTOS: catálogo de items que se venden
-- ============================================================
CREATE TABLE IF NOT EXISTS productos (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre          TEXT NOT NULL,                          -- nombre del producto
    sku             TEXT UNIQUE,                            -- código interno (opcional, ej: "BEB-001")
    categoria_id    INTEGER REFERENCES categorias(id),      -- FK a categorias (opcional)
    precio_venta    REAL NOT NULL CHECK(precio_venta > 0),  -- precio de venta al público
    precio_compra   REAL NOT NULL DEFAULT 0,                -- costo (para márgenes, opcional en MVP)
    stock_actual    REAL NOT NULL DEFAULT 0,                -- stock disponible actual
    stock_minimo    REAL NOT NULL DEFAULT 0,                -- alerta de stock bajo
    unidad          TEXT NOT NULL DEFAULT 'unidad'          -- unidad de medida: unidad, kg, litro, paquete, etc.
                                                    CHECK(unidad IN ('unidad', 'kg', 'litro', 'paquete')),
    activo          INTEGER NOT NULL DEFAULT 1,
    created_at      TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    updated_at      TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- ============================================================
-- CLIENTES: personas que compran (datos mínimos)
-- ============================================================
CREATE TABLE IF NOT EXISTS clientes (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre          TEXT NOT NULL,
    telefono        TEXT,                                   -- opcional
    direccion       TEXT,                                   -- opcional
    created_at      TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- ============================================================
-- VENTAS: cabecera de cada transacción de venta
-- ============================================================
CREATE TABLE IF NOT EXISTS ventas (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    usuario_id      INTEGER NOT NULL REFERENCES usuarios(id),   -- cajero que registró
    cliente_id      INTEGER REFERENCES clientes(id),            -- opcional
    total           REAL NOT NULL CHECK(total >= 0),            -- suma de todos los items
    metodo_pago     TEXT NOT NULL CHECK(metodo_pago IN ('efectivo', 'tarjeta', 'transferencia', 'mixto')),
    created_at      TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- ============================================================
-- VENTA_ITEMS: detalle de cada producto en una venta
-- ============================================================
CREATE TABLE IF NOT EXISTS venta_items (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    venta_id        INTEGER NOT NULL REFERENCES ventas(id) ON DELETE CASCADE,
    producto_id     INTEGER NOT NULL REFERENCES productos(id),
    cantidad        REAL NOT NULL CHECK(cantidad > 0),         -- cantidad vendida (ej: 2, 0.5 si es kg)
    precio_unitario REAL NOT NULL,                              -- snapshot del precio al momento de vender
    subtotal        REAL NOT NULL CHECK(subtotal >= 0)         -- cantidad * precio_unitario
);

-- ============================================================
-- INVENTARIO_MOVIMIENTOS: registro de entrada/salida/ajuste de stock
-- ============================================================
CREATE TABLE IF NOT EXISTS inventario_movimientos (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    producto_id     INTEGER NOT NULL REFERENCES productos(id),
    tipo            TEXT NOT NULL CHECK(tipo IN ('entrada', 'salida', 'ajuste')),   -- entrada: compra, salida: venta/merma, ajuste: corrección
    cantidad        REAL NOT NULL CHECK(cantidad > 0),         -- siempre positiva (el tipo define la dirección)
    stock_resultante REAL NOT NULL,                            -- stock_actual del producto DESPUÉS del movimiento
    referencia_tipo TEXT,                                       -- qué originó el movimiento: 'venta', 'compra', 'manual'
    referencia_id   INTEGER,                                    -- ID de la venta o compra que originó el movimiento
    motivo          TEXT,                                       -- explicación (obligatorio en ajustes y mermas)
    usuario_id      INTEGER REFERENCES usuarios(id),            -- quién realizó el movimiento
    created_at      TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- ============================================================
-- CONFIGURACION: clave-valor para settings del sistema
-- ============================================================
CREATE TABLE IF NOT EXISTS configuracion (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    clave           TEXT NOT NULL UNIQUE,                       -- ej: 'nombre_negocio', 'pin_admin'
    valor           TEXT NOT NULL
);

-- ============================================================
-- ÍNDICES
-- ============================================================

-- Para buscar productos rápido en la venta
CREATE INDEX idx_productos_nombre ON productos(nombre);
CREATE INDEX idx_productos_sku ON productos(sku);
CREATE INDEX idx_productos_categoria ON productos(categoria_id);
CREATE INDEX idx_productos_stock ON productos(stock_actual);    -- para alertas de stock bajo

-- Para consultar ventas por fecha, cajero o cliente
CREATE INDEX idx_ventas_created_at ON ventas(created_at);
CREATE INDEX idx_ventas_usuario ON ventas(usuario_id);
CREATE INDEX idx_ventas_cliente ON ventas(cliente_id);

-- Para análisis de productos más vendidos
CREATE INDEX idx_venta_items_venta ON venta_items(venta_id);
CREATE INDEX idx_venta_items_producto ON venta_items(producto_id);

-- Para historial de inventario
CREATE INDEX idx_inventario_producto ON inventario_movimientos(producto_id);
CREATE INDEX idx_inventario_fecha ON inventario_movimientos(created_at);
CREATE INDEX idx_inventario_tipo ON inventario_movimientos(tipo);

-- Para búsqueda de clientes
CREATE INDEX idx_clientes_nombre ON clientes(nombre);
CREATE INDEX idx_clientes_telefono ON clientes(telefono);

-- Para login rápido
CREATE INDEX idx_usuarios_pin ON usuarios(pin_hash);

-- ============================================================
-- TRIGGERS
-- ============================================================

-- Actualizar stock_actual en productos cuando se inserta un movimiento de inventario
CREATE TRIGGER IF NOT EXISTS trg_inventario_actualiza_stock
    AFTER INSERT ON inventario_movimientos
BEGIN
    UPDATE productos
    SET stock_actual = (
        SELECT stock_resultante
        FROM inventario_movimientos
        WHERE id = NEW.id
    ),
        updated_at = datetime('now', 'localtime')
    WHERE id = NEW.producto_id;
END;

-- Registrar automáticamente movimiento de salida cuando se inserta un item de venta
-- (Esto se hace desde la app, no con trigger, porque necesitamos el total de la venta)
```

---

## Relaciones

```
usuarios 1──N ventas                    — un cajero registra muchas ventas
usuarios 1──N inventario_movimientos    — un usuario registra movimientos
categorias 1──N productos               — una categoría agrupa productos
ventas 1──N venta_items                — una venta tiene N items
productos 1──N venta_items             — un producto aparece en muchas ventas
productos 1──N inventario_movimientos   — un producto tiene muchos movimientos
clientes 1──N ventas                   — un cliente puede comprar muchas veces
```

---

## Consideraciones para NL→SQL

### Tablas consultables (expuestas al prompt de OpenRouter)

| Tabla | Aliases para NL | Descripción |
|-------|----------------|-------------|
| `productos` | productos, mercadería, artículos | Catálogo completo con precio y stock |
| `ventas` | ventas, facturas, tickets | Cada transacción de venta |
| `venta_items` | items de venta, detalle, productos vendidos | Qué productos se vendieron en cada venta |
| `clientes` | clientes, compradores, personas | Datos de clientes |
| `usuarios` | usuarios, cajeros, empleados | Quién registró cada venta |
| `inventario_movimientos` | movimientos, entradas, salidas, ajustes | Historial de cambios de stock |
| `categorias` | categorías, rubros, secciones | Agrupación de productos |

### Reglas para el generador NL→SQL

1. **Nunca modificar datos** — solo SELECT. El prompt debe forzar read-only.
2. **Fechas en locale** — todas las created_at están en `datetime('now','localtime')`, las consultas deben usar `date(created_at)` y `strftime` correctamente.
3. **Totales económicos** — ventas.total y venta_items.subtotal son REAL, usar `COALESCE(SUM(...), 0)` para evitar NULLs.
4. **Productos inactivos** — por defecto filtrar `productos.activo = 1` a menos que se pida explícitamente ver inactivos.
5. **Nombres en castellano** — las columnas ya están en español, pero los valores también: `'efectivo'`, `'tarjeta'`, `'entrada'`, `'salida'`, `'ajuste'`.

---

## 5 Consultas SQL de Ejemplo (Few-Shot para el Prompt NL→SQL)

### 1. Ventas del día de hoy

```sql
-- Pregunta: "¿Cuánto vendí hoy?"
SELECT COALESCE(SUM(total), 0) as total_vendido,
       COUNT(*) as cantidad_ventas
FROM ventas
WHERE date(created_at) = date('now', 'localtime');
```

### 2. Productos más vendidos de la semana

```sql
-- Pregunta: "¿Cuáles son los 5 productos más vendidos esta semana?"
SELECT p.nombre                            as producto,
       SUM(vt.cantidad)                    as cantidad_vendida,
       SUM(vt.subtotal)                    as total_generado
FROM venta_items vt
JOIN productos p    ON p.id = vt.producto_id
JOIN ventas v       ON v.id = vt.venta_id
WHERE v.created_at >= datetime('now', '-7 days', 'localtime')
  AND p.activo = 1
GROUP BY p.id
ORDER BY cantidad_vendida DESC
LIMIT 5;
```

### 3. Productos con stock bajo

```sql
-- Pregunta: "¿Qué productos están por agotarse?"
SELECT p.nombre          as producto,
       p.stock_actual    as stock_actual,
       p.stock_minimo    as stock_minimo,
       c.nombre          as categoria
FROM productos p
LEFT JOIN categorias c ON c.id = p.categoria_id
WHERE p.stock_actual <= p.stock_minimo
  AND p.activo = 1
ORDER BY p.stock_actual ASC;
```

### 4. Ventas por cajero hoy

```sql
-- Pregunta: "¿Cuánto vendió cada cajero hoy?"
SELECT u.nombre                           as cajero,
       COUNT(v.id)                        as cantidad_ventas,
       COALESCE(SUM(v.total), 0)          as total_vendido
FROM usuarios u
LEFT JOIN ventas v ON v.usuario_id = u.id
    AND date(v.created_at) = date('now', 'localtime')
WHERE u.activo = 1
GROUP BY u.id
ORDER BY total_vendido DESC;
```

### 5. Últimas ventas de un cliente

```sql
-- Pregunta: "Mostrame las compras de Juan Pérez"
SELECT v.id              as venta_id,
       v.total           as total,
       v.metodo_pago     as metodo_pago,
       v.created_at      as fecha
FROM ventas v
JOIN clientes c ON c.id = v.cliente_id
WHERE c.nombre LIKE '%Juan Pérez%'
ORDER BY v.created_at DESC
LIMIT 10;
```

---

## Archivo `.sql` para sqlc

El DDL completo está en el bloque de código de la sección anterior. Para sqlc se usa exactamente ese archivo. La configuración de sqlc apunta a ese DDL y los queries `.sql` aparte.

Ruta sugerida: `db/migrations/001_init.sql`

La configuración de sqlc (`sqlc.json` o `sqlc.yaml`) debería apuntar a `db/queries/` para los queries y `db/migrations/` para el schema.

---

## Notas de Diseño

- **Stock como REAL** — soporta productos fraccionables (kg, litros) sin conversiones
- **Movimientos de inventario con `referencia_tipo`/`referencia_id`** — permite rastrear el origen sin crear FKs polimórficas complejas
- **`stock_resultante` en cada movimiento** — evita tener que calcular saldo histórico en cada consulta de stock
- **Trigger para stock_actual** — asegura que productos.stock_actual siempre refleje el último movimiento
- **Sin soft-delete** — los registros se marcan como `activo = 0`, no se eliminan físicamente
- **Check constraints en lugar de tablas** — `metodo_pago`, `tipo`, `unidad`, `rol` se validan con CHECK, no con tablas separadas. Si en el futuro se necesitan dinámicos, se migra a tablas.
- **FK sin ON DELETE CASCADE** (excepto venta_items) — se prefiere control explícito desde la app para evitar pérdidas accidentales
