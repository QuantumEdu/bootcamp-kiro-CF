# Research: Dashboard de Métricas Fijas

**Ticket:** 005  
**Stack:** Go + HTMX + SQLite + sqlc  
**Layout:** Dashboard partido — métricas fijas arriba, chat NL siempre visible abajo  
**Usuario:** Dueño de negocio consultando su POS

---

## Schema de referencia

Asumiendo el modelo del ticket 001:

```sql
CREATE TABLE productos (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre      TEXT NOT NULL,
    precio      REAL NOT NULL,
    categoria_id INTEGER REFERENCES categorias(id),
    stock_actual INTEGER NOT NULL DEFAULT 0,
    sku         TEXT UNIQUE,
    created_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE categorias (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL UNIQUE
);

CREATE TABLE ventas (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    fecha      TEXT NOT NULL DEFAULT (datetime('now')),
    total      REAL NOT NULL,
    usuario_id INTEGER REFERENCES usuarios(id),
    cliente_id INTEGER REFERENCES clientes(id)
);

CREATE TABLE ventas_items (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    venta_id      INTEGER NOT NULL REFERENCES ventas(id),
    producto_id   INTEGER NOT NULL REFERENCES productos(id),
    cantidad      INTEGER NOT NULL,
    precio_unitario REAL NOT NULL
);

CREATE TABLE clientes (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre        TEXT NOT NULL,
    telefono      TEXT,
    direccion     TEXT,
    fecha_registro TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE usuarios (
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre   TEXT NOT NULL,
    pin_hash TEXT NOT NULL,
    rol      TEXT NOT NULL DEFAULT 'cajero'
);

CREATE TABLE movimientos_inventario (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    producto_id INTEGER NOT NULL REFERENCES productos(id),
    tipo        TEXT NOT NULL CHECK(tipo IN ('entrada','salida','ajuste')),
    cantidad    INTEGER NOT NULL,
    motivo      TEXT,
    fecha       TEXT NOT NULL DEFAULT (datetime('now')),
    usuario_id  INTEGER REFERENCES usuarios(id)
);
```

---

## 1. Ventas del día / semana / mes

### Nombre (UI)
**"Ventas hoy"** / **"Ventas esta semana"** / **"Ventas este mes"**

### Queries SQL

```sql
-- Ventas del día de hoy
SELECT COUNT(*)         AS cantidad,
       COALESCE(SUM(total), 0) AS total
FROM ventas
WHERE DATE(fecha) = DATE('now');

-- Ventas de la semana (últimos 7 días, incluido hoy)
SELECT DATE(fecha) AS dia,
       COUNT(*)    AS cantidad,
       SUM(total)  AS total
FROM ventas
WHERE fecha >= datetime('now', '-7 days')
GROUP BY DATE(fecha)
ORDER BY dia DESC;

-- Ventas del mes actual
SELECT COUNT(*)         AS cantidad,
       COALESCE(SUM(total), 0) AS total
FROM ventas
WHERE strftime('%Y-%m', fecha) = strftime('%Y-%m', 'now');
```

### Renderizado
- **Ventas hoy**: número grande + moneda (ej. `$4,230`) con badge de cantidad de transacciones al lado. Color primario.
- **Ventas semana**: mini sparkline (barras verticales delgadas, 7 columnas) + total acumulado al final. HTMX podría renderizar un SVG inline.
- **Ventas mes**: número grande secundario, menos jerarquía que "hoy".

### Frecuencia
- **Ventas hoy**: cada carga de página + polling cada 30 segundos (HTMX `hx-trigger="every 30s"`).
- **Ventas semana/mes**: cada carga de página + polling cada 60 segundos.

### Jerarquía
**Primera fila del dashboard** — es lo primero que el dueño quiere ver.

---

## 2. Top 5 productos más vendidos

### Nombre (UI)
**"Productos más vendidos"**

### Query SQL

```sql
SELECT p.nombre,
       SUM(vi.cantidad)     AS unidades,
       SUM(vi.cantidad * vi.precio_unitario) AS total_venta,
       COUNT(DISTINCT v.id) AS veces_vendido
FROM ventas_items vi
JOIN productos p ON p.id = vi.producto_id
JOIN ventas v ON v.id = vi.venta_id
WHERE v.fecha >= datetime('now', '-30 days')
GROUP BY vi.producto_id
ORDER BY unidades DESC
LIMIT 5;
```

### Renderizado
- Tabla con 4 columnas: Producto | Unidades | Total | Veces vendido.
- Primera fila destacada con badge de oro.
- Si hay 0 ventas, mostrar empty state amigable ("Aún no hay ventas registradas").

### Frecuencia
Cada carga de página + cada 60 segundos.

### Jerarquía
**Segunda fila, columna izquierda** — junto al stock bajo.

---

## 3. Productos con stock bajo (alertas)

### Nombre (UI)
**"Stock bajo — revisar"**

### Query SQL

```sql
SELECT p.nombre,
       p.stock_actual,
       p.sku,
       COALESCE(SUM(vi.cantidad), 0) AS unidades_vendidas_30d
FROM productos p
LEFT JOIN ventas_items vi ON vi.producto_id = p.id
LEFT JOIN ventas v ON v.id = vi.venta_id
    AND v.fecha >= datetime('now', '-30 days')
WHERE p.stock_actual <= 5
GROUP BY p.id
ORDER BY p.stock_actual ASC;
```

### Renderizado
- Tabla con columnas: Producto | Stock actual | SKU | Vendidos (30d).
- Filas con stock = 0 en **rojo** con icono de alerta.
- Filas con stock entre 1-5 en **amarillo/naranja**.
- Si no hay productos con bajo stock, mostrar badge verde "✅ Todo en orden".

### Frecuencia
Cada carga de página. No necesita polling agresivo — el stock cambia con ventas.

### Jerarquía
**Segunda fila, columna derecha** — junto al top 5.

---

## 4. Clientes frecuentes / nuevos

### Nombre (UI)
**"Clientes frecuentes"** y **"Clientes nuevos (30 días)"**

### Queries SQL

```sql
-- Top 5 clientes que más compraron (últimos 30 días)
SELECT c.nombre,
       COUNT(v.id)          AS compras,
       COALESCE(SUM(v.total), 0) AS total_gastado
FROM clientes c
JOIN ventas v ON v.cliente_id = c.id
WHERE v.fecha >= datetime('now', '-30 days')
GROUP BY c.id
ORDER BY compras DESC
LIMIT 5;

-- Clientes nuevos en los últimos 30 días
SELECT COUNT(*) AS nuevos
FROM clientes
WHERE fecha_registro >= datetime('now', '-30 days');
```

### Renderizado
- **Frecuentes**: tabla pequeña: Cliente | Compras | Total. Solo top 5.
- **Nuevos**: número con label "nuevos clientes este mes".

### Frecuencia
Cada carga de página + cada 120 segundos.

### Jerarquía
**Tercera fila, columna izquierda**.

---

## 5. Ingresos totales por período

### Nombre (UI)
**"Resumen de ingresos"**

### Query SQL

```sql
-- Ingresos agrupados por día (últimos 14 días)
SELECT DATE(v.fecha) AS dia,
       COUNT(*)      AS ventas,
       SUM(v.total)  AS ingresos
FROM ventas v
WHERE v.fecha >= datetime('now', '-14 days')
GROUP BY DATE(v.fecha)
ORDER BY dia ASC;

-- Totales acumulados
SELECT 'hoy'       AS periodo, COUNT(*) AS ventas, SUM(v.total) AS ingresos FROM ventas v WHERE DATE(v.fecha) = DATE('now')
UNION ALL
SELECT 'semana'    AS periodo, COUNT(*) AS ventas, SUM(v.total) AS ingresos FROM ventas v WHERE v.fecha >= datetime('now', '-7 days')
UNION ALL
SELECT 'mes'       AS periodo, COUNT(*) AS ventas, SUM(v.total) AS ingresos FROM ventas v WHERE strftime('%Y-%m', v.fecha) = strftime('%Y-%m', 'now')
UNION ALL
SELECT 'total gral' AS periodo, COUNT(*) AS ventas, SUM(v.total) AS ingresos FROM ventas v;
```

### Renderizado
- Gráfico de barras simple (SVG inline o barra visual tipo progress). 14 barras.
- Debajo, tabla resumen con 4 filas (hoy, semana, mes, total).
- Alternativa MVP: solo la tabla resumen sin gráfico (menos código, mismo valor).

### Frecuencia
Cada carga de página + cada 60 segundos.

### Jerarquía
**Tercera fila, columna derecha** — información complementaria.

---

## 6. Productos sin movimiento (nunca vendidos)

### Nombre (UI)
**"Productos sin rotación"**

### Query SQL

```sql
SELECT p.nombre,
       p.sku,
       p.stock_actual,
       p.precio,
       (p.stock_actual * p.precio) AS valor_inmovilizado
FROM productos p
WHERE p.id NOT IN (
    SELECT DISTINCT vi.producto_id
    FROM ventas_items vi
    JOIN ventas v ON v.id = vi.venta_id
)
ORDER BY valor_inmovilizado DESC;
```

### Renderizado
- Tabla: Producto | SKU | Stock | Precio | Valor inmovilizado.
- Destacar en rojo si `valor_inmovilizado > 1000` (ajustable).
- Si todos los productos se han vendido al menos una vez, mostrar badge verde.

### Frecuencia
On-demand (cada carga de página). No cambia sin ventas nuevas o productos nuevos.

### Jerarquía
**Cuarta fila, oculto tras expander "Ver más métricas"** — no es crítica para el día a día, pero útil para decisiones de inventario.

---

## 7. Input de consulta SQL libre

### Nombre (UI)
**"Consultar datos"** (sección plegable)

### Diseño MVP
- **Textarea** monospace con placeholder: `Escribe tu consulta SQL aquí...`
- **Botón "Ejecutar"** al lado
- **Validación** del lado del servidor:
  - Solo permitir `SELECT` (rechazar `INSERT`, `UPDATE`, `DELETE`, `DROP`, `ALTER`, `PRAGMA`, etc.)
  - Timeout de 5 segundos
  - Límite de 100 filas en resultado (`LIMIT 100` forzado si no lo incluye)
  - Modo read-only en la conexión SQLite (abrir con `?mode=ro` o usar dos conexiones: una RW para la app, una RO para queries)
- **Resultado**: tabla HTML generada con HTMX. Si hay error, mostrar el mensaje de error de SQLite.
- **NO** incluir autocompletado ni syntax highlighting en MVP (scope para día 4 si sobra tiempo).

### Seguridad
- Conexión **read-only** separada para queries del usuario.
- Validar que el SQL comience con `SELECT` o `WITH` (CTE) — y que no contenga palabras clave peligrosas.
- Log de todas las queries ejecutadas para auditoría.

---

## 8. Integración con el chat NL

### Modelo de integración

El chat NL y el dashboard NO están aislados — comparten el mismo backend de datos. El flujo es:

1. Usuario escribe en el chat: *"¿qué vendí ayer?"*
2. Backend recibe, envía a OpenRouter, obtiene SQL, ejecuta, obtiene resultados.
3. **Respuesta en el chat**: texto en castellano con los datos formateados.
4. **Opcional (MVP diferido)**: el chat también puede disparar un reemplazo HTMX de una sección del dashboard. Ej: si preguntó "¿qué vendí ayer?", la respuesta aparece en el chat como texto, y además el backend podría devolver un `OOB` (Out of Band) swap que actualice la tarjeta "Ventas día" con el valor de ayer.

### Estrategia MVP (día 1-3)
- Chat solo responde con texto.
- No toca el dashboard.

### Estrategia extendida (si sobra tiempo, día 4)
- HTMX OOB swaps: el endpoint del chat devuelve `HX-Trigger` o múltiples `<template>` OOB para refrescar secciones afectadas.
- Esto hace que el dashboard "reaccione" a las preguntas del chat sin recarga completa.

### Visualización unificada

```
┌─────────────────────────────────────────────────────────┐
│  DASHBOARD (métricas fijas, se auto-actualizan)         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │ Ventas hoy  │  │ Ventas sem  │  │ Ventas mes  │     │
│  │ $4,230      │  │ ▁▃▄▆▅▇▆     │  │ $42,100     │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
│  ┌──────────────────┐  ┌──────────────────┐             │
│  │ Top 5 productos  │  │ Stock bajo       │             │
│  └──────────────────┘  └──────────────────┘             │
├─────────────────────────────────────────────────────────┤
│  💬 ¿qué vendí ayer?                           [Enviar] │
│  ┌─────────────────────────────────────────────────┐    │
│  │ 🤖 Ayer vendiste $2,340 en 12 transacciones.    │    │
│  │ El producto más vendido fue Café Latte (34 uni- │    │
│  │ dades).                                         │    │
│  └─────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

---

## Maqueta de layout (jerarquía visual)

```
Fila 1 (3 cards, mismo peso visual):
┌──────────────┐ ┌──────────────────┐ ┌──────────────────┐
│  VENTAS HOY   │ │   VENTAS SEMANA  │ │   VENTAS MES     │
│  $4,230       │ │   $18,900        │ │   $42,100        │
│  12 transac   │ │   ▁▃▄▆▅▇▆        │ │   89 transac     │
└──────────────┘ └──────────────────┘ └──────────────────┘

Fila 2 (split 50/50):
┌─────────────────────────┐ ┌──────────────────────────┐
│ TOP 5 PRODUCTOS         │ │ STOCK BAJO (alertas)     │
│ #1 Café Latte    34uds  │ │ ⚠️ Leche        0       │
│ #2 Croissant     22uds  │ │ ⚠️ Azúcar       2       │
│ #3 Muffin        18uds  │ │ ⚠️ Yerba        3       │
│ #4 Tostado       15uds  │ │                          │
│ #5 Jugo Naranja  12uds  │ │ ✅ Todo en orden         │
└─────────────────────────┘ └──────────────────────────┘

Fila 3 (split 50/50):
┌─────────────────────────┐ ┌──────────────────────────┐
│ CLIENTES FRECUENTES     │ │ INGRESOS (14 días)       │
│ Juan P.      8 compras  │ │ ██▃▅▇▆▄▅▆██▇▆▅          │
│ María L.     5 compras  │ │                          │
│ Carlos G.    3 compras  │ │ Hoy:   $4,230            │
│ +5 nuevos este mes      │ │ Sem:   $18,900           │
└─────────────────────────┘ └──────────────────────────┘

Fila 4 (plegable, "Ver más métricas"):
┌─────────────────────────┐
│ PRODUCTOS SIN ROTACIÓN  │
│ (tabla)                 │
└─────────────────────────┘

Sección plegable adicional:
┌─────────────────────────────────────────────────────────┐
│ ⚡ CONSULTA SQL LIBRE                                     │
│ [________________________________________________]      │
│ [Ejecutar]                                               │
│ ┌────────────────────────────────────────────────────┐  │
│ │ Resultados en tabla HTML                           │  │
│ └────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

---

## Resumen de frecuencias

| Métrica | Frecuencia | HTMX trigger |
|---------|-----------|--------------|
| Ventas hoy | 30s | `every 30s` |
| Ventas semana | 60s | `every 60s` |
| Ventas mes | 60s | `every 60s` |
| Top 5 productos | 60s | `every 60s` |
| Stock bajo | on load | `load` |
| Clientes frecuentes | 120s | `every 120s` |
| Ingresos 14 días | 60s | `every 60s` |
| Productos sin rotación | on load | `load` |
| Consulta SQL libre | on-demand | click |
| Chat NL | on-demand | submit |

**Nota MVP**: no sobrecargar de polling. Cada request HTMX es liviano (el server devuelve solo el HTML del fragmento). 6 requests por minuto como máximo.

---

## Decisiones para implementación

1. **Todas las queries van a `sqlc`** — nada de raw SQL en handlers. Cada métrica tiene su query named en `queries/metrics.sql`.
2. **Endpoints HTMX por métrica**: `/metrics/ventas-hoy`, `/metrics/top-5`, etc. Cada uno devuelve un fragmento HTML que se reemplaza vía `hx-target`.
3. **Conexión read-only** para consulta SQL libre: abrir SQLite con `?mode=ro` desde Go, o usar `PRAGMA query_only = ON;`.
4. **Chat y dashboard comparten el mismo pool de conexión** — no hay duplicación de datos.
5. **No hay WebSockets en MVP** — HTMX polling es suficiente para las frecuencias definidas.
