-- 001-init.sql
-- Schema MVP del POS AI-First
-- SQLite + sqlc
-- PRAGMA journal_mode=WAL;

-- ============================================================
-- USUARIOS: cajeros y administradores del sistema
-- ============================================================
CREATE TABLE IF NOT EXISTS usuarios (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre          TEXT NOT NULL,
    pin_hash        TEXT NOT NULL,
    rol             TEXT NOT NULL CHECK(rol IN ('admin', 'cajero')),
    activo          INTEGER NOT NULL DEFAULT 1,
    created_at      TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- ============================================================
-- CATEGORIAS: agrupación de productos
-- ============================================================
CREATE TABLE IF NOT EXISTS categorias (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre          TEXT NOT NULL UNIQUE,
    descripcion     TEXT,
    activo          INTEGER NOT NULL DEFAULT 1
);

-- ============================================================
-- PRODUCTOS: catálogo de items que se venden
-- ============================================================
CREATE TABLE IF NOT EXISTS productos (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre          TEXT NOT NULL,
    sku             TEXT UNIQUE,
    categoria_id    INTEGER REFERENCES categorias(id),
    precio_venta    REAL NOT NULL CHECK(precio_venta > 0),
    precio_compra   REAL NOT NULL DEFAULT 0,
    stock_actual    REAL NOT NULL DEFAULT 0,
    stock_minimo    REAL NOT NULL DEFAULT 0,
    unidad          TEXT NOT NULL DEFAULT 'unidad'
                    CHECK(unidad IN ('unidad', 'kg', 'litro', 'paquete')),
    activo          INTEGER NOT NULL DEFAULT 1,
    created_at      TEXT NOT NULL DEFAULT (datetime('now','localtime')),
    updated_at      TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- ============================================================
-- CLIENTES: personas que compran
-- ============================================================
CREATE TABLE IF NOT EXISTS clientes (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre          TEXT NOT NULL,
    telefono        TEXT,
    direccion       TEXT,
    created_at      TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- ============================================================
-- VENTAS: cabecera de cada transacción
-- ============================================================
CREATE TABLE IF NOT EXISTS ventas (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    usuario_id      INTEGER NOT NULL REFERENCES usuarios(id),
    cliente_id      INTEGER REFERENCES clientes(id),
    total           REAL NOT NULL CHECK(total >= 0),
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
    cantidad        REAL NOT NULL CHECK(cantidad > 0),
    precio_unitario REAL NOT NULL,
    subtotal        REAL NOT NULL CHECK(subtotal >= 0)
);

-- ============================================================
-- INVENTARIO_MOVIMIENTOS: registro de cambios de stock
-- ============================================================
CREATE TABLE IF NOT EXISTS inventario_movimientos (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    producto_id     INTEGER NOT NULL REFERENCES productos(id),
    tipo            TEXT NOT NULL CHECK(tipo IN ('entrada', 'salida', 'ajuste')),
    cantidad        REAL NOT NULL CHECK(cantidad > 0),
    stock_resultante REAL NOT NULL,
    referencia_tipo TEXT,
    referencia_id   INTEGER,
    motivo          TEXT,
    usuario_id      INTEGER REFERENCES usuarios(id),
    created_at      TEXT NOT NULL DEFAULT (datetime('now','localtime'))
);

-- ============================================================
-- CONFIGURACION: clave-valor del sistema
-- ============================================================
CREATE TABLE IF NOT EXISTS configuracion (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    clave           TEXT NOT NULL UNIQUE,
    valor           TEXT NOT NULL
);

-- ============================================================
-- ÍNDICES
-- ============================================================
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

-- ============================================================
-- TRIGGERS
-- ============================================================
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
