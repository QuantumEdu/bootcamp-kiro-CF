-- 001_init.sql — PostgreSQL equivalent of SQLite 001_init.sql + 002_sessions.sql
-- Schema MVP del POS AI-First (PostgreSQL / RDS)

-- ============================================================
-- USUARIOS: cajeros y administradores del sistema
-- ============================================================
CREATE TABLE IF NOT EXISTS usuarios (
    id              SERIAL PRIMARY KEY,
    nombre          TEXT NOT NULL,
    pin_hash        TEXT NOT NULL,
    rol             TEXT NOT NULL CHECK(rol IN ('admin', 'cajero')),
    activo          BOOLEAN NOT NULL DEFAULT true,
    failed_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- CATEGORIAS: agrupación de productos
-- ============================================================
CREATE TABLE IF NOT EXISTS categorias (
    id              SERIAL PRIMARY KEY,
    nombre          TEXT NOT NULL UNIQUE,
    descripcion     TEXT,
    activo          BOOLEAN NOT NULL DEFAULT true
);

-- ============================================================
-- PRODUCTOS: catálogo de items que se venden
-- ============================================================
CREATE TABLE IF NOT EXISTS productos (
    id              SERIAL PRIMARY KEY,
    nombre          TEXT NOT NULL,
    sku             TEXT UNIQUE,
    categoria_id    INTEGER REFERENCES categorias(id),
    precio_venta    NUMERIC(12,2) NOT NULL CHECK(precio_venta > 0),
    precio_compra   NUMERIC(12,2) NOT NULL DEFAULT 0,
    stock_actual    NUMERIC(12,2) NOT NULL DEFAULT 0,
    stock_minimo    NUMERIC(12,2) NOT NULL DEFAULT 0,
    unidad          TEXT NOT NULL DEFAULT 'unidad'
                    CHECK(unidad IN ('unidad', 'kg', 'litro', 'paquete')),
    activo          BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- CLIENTES: personas que compran
-- ============================================================
CREATE TABLE IF NOT EXISTS clientes (
    id              SERIAL PRIMARY KEY,
    nombre          TEXT NOT NULL,
    telefono        TEXT,
    direccion       TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- VENTAS: cabecera de cada transacción
-- ============================================================
CREATE TABLE IF NOT EXISTS ventas (
    id              SERIAL PRIMARY KEY,
    usuario_id      INTEGER NOT NULL REFERENCES usuarios(id),
    cliente_id      INTEGER REFERENCES clientes(id),
    total           NUMERIC(12,2) NOT NULL CHECK(total >= 0),
    metodo_pago     TEXT NOT NULL CHECK(metodo_pago IN ('efectivo', 'tarjeta', 'transferencia', 'mixto')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- VENTA_ITEMS: detalle de cada producto en una venta
-- ============================================================
CREATE TABLE IF NOT EXISTS venta_items (
    id              SERIAL PRIMARY KEY,
    venta_id        INTEGER NOT NULL REFERENCES ventas(id) ON DELETE CASCADE,
    producto_id     INTEGER NOT NULL REFERENCES productos(id),
    cantidad        NUMERIC(12,2) NOT NULL CHECK(cantidad > 0),
    precio_unitario NUMERIC(12,2) NOT NULL,
    subtotal        NUMERIC(12,2) NOT NULL CHECK(subtotal >= 0)
);

-- ============================================================
-- INVENTARIO_MOVIMIENTOS: registro de cambios de stock
-- ============================================================
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

-- ============================================================
-- CONFIGURACION: clave-valor del sistema
-- ============================================================
CREATE TABLE IF NOT EXISTS configuracion (
    id              SERIAL PRIMARY KEY,
    clave           TEXT NOT NULL UNIQUE,
    valor           TEXT NOT NULL
);

-- ============================================================
-- SESSIONS: compatible con alexedwards/scs pgxstore format
-- ============================================================
CREATE TABLE IF NOT EXISTS sessions (
    token   TEXT PRIMARY KEY,
    data    BYTEA NOT NULL,
    expiry  TIMESTAMPTZ NOT NULL
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

CREATE INDEX IF NOT EXISTS idx_sessions_expiry ON sessions(expiry);

-- ============================================================
-- TRIGGER: actualizar stock_actual al insertar movimiento
-- ============================================================
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
