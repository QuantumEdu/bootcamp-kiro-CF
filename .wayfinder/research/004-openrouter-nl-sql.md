# NL→SQL con OpenRouter para POS AI-First MVP

## 1. Modelo recomendado de OpenRouter

| Modelo | ID en OpenRouter | Costo input/M tokens | Costo output/M tokens | Structured Outputs |
|---|---|---|---|---|
| **GPT-4o (recomendado primario)** | `openai/gpt-4o` | $2.50 | $10.00 | Sí |
| GPT-4o-mini (fallback barato) | `openai/gpt-4o-mini` | $0.15 | $0.60 | Sí |
| Mistral Large | `mistralai/mistral-large` | $2.00 | $6.00 | Sí |
| Llama 3.1 70B | `meta-llama/llama-3.1-70b-instruct` | $0.40 | $0.40 | Sí |
| Mistral Nemo (12B) | `mistralai/mistral-nemo` | $0.02 | $0.03 | Sí |

### Recomendación primaria: **GPT-4o** (`openai/gpt-4o`)

**Pros:**
- Structured Outputs nativo (JSON Schema enforcement estricto) — ideal para extraer SQL limpio
- Excelente entendimiento de lenguaje natural en castellano
- Fuerte en tareas de razonamiento y transformación (NL→SQL es esencialmente una tarea de razonamiento)
- Prompt caching via `input_cache_read` ($1.25/M tokens cacheados) si repetís schema en cada request
- Amplia documentación y comunidad

**Contra:**
- Más caro que alternativas open-source
- Dependencia de API externa (latred de red)

### Fallback barato: **GPT-4o-mini** (`openai/gpt-4o-mini`)

Para queries simples o para desarrollo local donde no necesitás precisión máxima. Structured Outputs también disponible.

**Por qué NO Claude 3.5 Sonnet:** Claude no soporta `structured_outputs` nativo en OpenRouter. No podés forzar un JSON Schema estricto y terminás parseando con regex o esperando que el modelo coopere — frágil.

**Por qué NO DeepSeek:** DeepSeek no soporta `structured_outputs` via OpenRouter. Para NL→SQL el enforcement de esquema es crítico.

---

## 2. System Prompt completo

```
Eres un asistente que convierte lenguaje natural a SQL para una base de datos SQLite
de un sistema POS (Point of Sale).

REGLAS ESTRICTAS:
1. Solo genera sentencias SELECT. NUNCA generes INSERT, UPDATE, DELETE, DROP, ALTER, CREATE, o cualquier DDL/DML.
2. Usa solo la sintaxis de SQLite.
3. Si no puedes generar una query segura con la información dada, responde con:
   {"sql": null, "error": "Descripción clara de por qué no se puede generar la query"}

ESQUEMA DE BASE DE DATOS:
{schema_dump}

CONVENCIONES DE NOMBRES:
- Las tablas y columnas están en inglés (snake_case)
- El usuario habla en castellano
- Mapeá conceptos del usuario a nombres reales de tablas/columnas
  Ejemplo: "producto" → products, "venta" → sales, "cliente" → customers

CONTEXTO ADICIONAL:
{tables_description}

Devuelve SIEMPRE un JSON válido con esta estructura exacta:
{"sql": "SELECT ...", "explanation": "Breve explicación de qué hace la query"}
```

El `{schema_dump}` se genera dinámicamente con `.schema` de SQLite (o un dump similar). Ejemplo:

```
CREATE TABLE products (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  price REAL NOT NULL,
  category_id INTEGER REFERENCES categories(id),
  created_at TEXT DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE sales (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  total REAL NOT NULL,
  payment_method TEXT NOT NULL,
  created_at TEXT DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE sale_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  sale_id INTEGER NOT NULL REFERENCES sales(id),
  product_id INTEGER NOT NULL REFERENCES products(id),
  quantity INTEGER NOT NULL,
  unit_price REAL NOT NULL
);
```

El `{tables_description}` humaniza los nombres:

```json
{
  "products": "Productos del catálogo. Cada producto tiene nombre, precio, categoría.",
  "sales": "Ventas realizadas. Cada venta tiene un total, método de pago y fecha.",
  "sale_items": "Items individuales de cada venta. Relaciona productos con ventas.",
  "categories": "Categorías de productos."
}
```

---

## 3. Estrategia de few-shot

### Estructura

Incluir 3-5 ejemplos en el system prompt cubriendo los casos de uso más comunes del POS:

```json
[
  {
    "user": "¿cuántos productos vendí esta semana?",
    "assistant": {
      "sql": "SELECT COUNT(DISTINCT si.product_id) FROM sale_items si JOIN sales s ON si.sale_id = s.id WHERE s.created_at >= date('now', '-7 days')",
      "explanation": "Cuenta productos distintos vendidos en los últimos 7 días"
    }
  },
  {
    "user": "qué producto se vendió más ayer",
    "assistant": {
      "sql": "SELECT p.name, SUM(si.quantity) as total_qty FROM sale_items si JOIN products p ON si.product_id = p.id JOIN sales s ON si.sale_id = s.id WHERE date(s.created_at) = date('now', '-1 day') GROUP BY p.id ORDER BY total_qty DESC LIMIT 1",
      "explanation": "Agrupa por producto, suma cantidades, ordena descendente y trae el primero"
    }
  },
  {
    "user": "mostrame las ventas de hoy",
    "assistant": {
      "sql": "SELECT s.id, s.total, s.payment_method, s.created_at FROM sales s WHERE date(s.created_at) = date('now') ORDER BY s.created_at DESC",
      "explanation": "Filtra ventas del día actual ordenadas por fecha descendente"
    }
  },
  {
    "user": "cuánto dinero en efectivo se cobró este mes",
    "assistant": {
      "sql": "SELECT SUM(total) as total_efectivo FROM sales WHERE payment_method = 'cash' AND strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now')",
      "explanation": "Suma total de ventas donde el método de pago es 'cash' en el mes actual"
    }
  }
]
```

**Importante:** Los ejemplos deben reflejar EL MISMO SCHEMA que usás en producción para que el modelo entienda el mapeo. Si el schema cambia, actualizá los ejemplos.

---

## 4. Parsing de respuesta

Usar **Structured Outputs** de OpenRouter con `response_format: json_schema`. Es el approach más robusto.

```json
{
  "model": "openai/gpt-4o",
  "response_format": {
    "type": "json_schema",
    "json_schema": {
      "name": "sql_response",
      "strict": true,
      "schema": {
        "type": "object",
        "properties": {
          "sql": {
            "type": ["string", "null"],
            "description": "SQL query generada, o null si no se puede generar"
          },
          "error": {
            "type": ["string", "null"],
            "description": "Mensaje de error si no se pudo generar SQL"
          },
          "explanation": {
            "type": "string",
            "description": "Explicación breve de lo que hace la query"
          }
        },
        "required": ["sql", "explanation"],
        "additionalProperties": false
      }
    }
  }
}
```

**Fallback si el modelo no soporta structured outputs (Mistral, Llama):**
Usar `response_format: { "type": "json_object" }` y parsear con `json.Unmarshal`. Si falla, intentar regex para extraer `{"sql": "..."}`.

**Response Healing plugin:** OpenRouter ofrece un plugin `response-healing` que repara JSON mal formado automáticamente. Útil para models sin structured outputs.

```json
{
  "plugins": [{ "id": "response-healing" }]
}
```

---

## 5. Manejo de errores

### Pipeline de validación

```
[User query]
    → OpenRouter call (con structured outputs)
    → Parse JSON response
    → ¿SQL presente? → Validar con parser SQLite (sqlparser-go o similar)
        → ¿SELECT válido? → Ejecutar contra SQLite (PRAGMA query_only = ON)
            → ¿Error? → Retry con mensaje de feedback
            → OK → Devolver resultados
        → NO SELECT → Devolver error de seguridad
    → Null SQL → Devolver error descriptivo al usuario
```

### Retry con feedback

Si la query generada es inválida, re-enviar al modelo con:

```
La query SQL que generaste es inválida o insegura. Error: {error}.
Por favor generá una nueva query SELECT de SQLite que resuelva:
{consulta_original}

Esquema disponible:
{schema_dump}
```

Máximo 2 retrys para evitar loops de costo.

### Timeout

OpenRouter tiene timeout default de ~30s. Para mantenerse dentro:

- `max_tokens`: 300 (el SQL no necesita ser largo)
- Si timeout, reintentar 1 vez con `model: "openai/gpt-4o-mini"` (más rápido)

---

## 6. Costos estimados

### Cálculo por consulta típica

| Componente | Tokens estimados |
|---|---|
| System prompt (schema + few-shot) | ~800-1200 |
| Query del usuario | ~20-50 |
| Output SQL | ~50-150 |
| **Total por consulta** | **~900-1400 tokens** |

### Costo por consulta

| Modelo | Costo por consulta (estimado) |
|---|---|
| GPT-4o | $0.0025 - $0.004 (~1400 tokens mix) |
| GPT-4o-mini | $0.00015 - $0.00025 |
| Mistral Nemo | $0.00002 - $0.00005 |

Con **prompt caching** (OpenRouter cachea el system prompt si se repite):
- GPT-4o: input cacheado a $1.25/M → ~$0.001 por consulta cacheada

### Costo mensual estimado (1000 consultas/día = 30K/mes)

| Modelo | Costo/mes |
|---|---|
| GPT-4o | ~$75-120 |
| GPT-4o-mini | ~$5-8 |
| Estrategia híbrida (GPT-4o + fallback a GPT-4o-mini) | ~$30-50 |

**Recomendación MVP:** Arrancar con GPT-4o-mini para desarrollo, pasar a GPT-4o en producción para precisión.

---

## 7. Cache de consultas

### ¿Vale la pena?

Sí, especialmente para consultas frecuentes del dashboard ("ventas de hoy", "producto más vendido", "total del día").

### Estrategia simple (sin Redis)

```go
type CacheEntry struct {
  SQL       string
  CreatedAt time.Time
  HitCount  int
}

var queryCache = make(map[string]*CacheEntry) // en memoria
var mu sync.RWMutex
```

- Key normalizada: `lowercase(trim(query))`
- TTL: 5 minutos (los datos POS cambian constantemente)
- Para MVP: map en memoria. Para producción: SQLite como cache (`.wayfinder/cache.db`)

### Cache por fingerprint de schema

Si el schema cambia (ej: se agrega una columna), invalidar todo el cache. Guardar `schemaVersion` (hash del `.schema` dump) como parte de la key.

---

## 8. Seguridad básica

### Restricciones a nivel de prompt (No alcanzan solas)

El system prompt prohíbe DDL/DML explícitamente. Pero un prompt injection podría bypassearlo.

### Validación runtime (OBLIGATORIA)

```go
import (
  "strings"
  "github.com/blastrain/vitess-sqlparser/sqlparser" // o similar
)

func ValidateSelectOnly(sql string) error {
  stmt, err := sqlparser.Parse(sql)
  if err != nil {
    return fmt.Errorf("SQL inválido: %w", err)
  }
  switch stmt.(type) {
  case *sqlparser.Select:
    return nil
  default:
    return fmt.Errorf("solo SELECT está permitido")
  }
}
```

### Otras medidas

- `PRAGMA query_only = ON` antes de ejecutar la consulta (solo disponible en SQLite ≥ 3.38)
- Ejecutar la query en una conexión de solo lectura
- `LIMIT 500` hardcodeado al final de cualquier query generada (evita queries masivas)
- Timeout de ejecución: 5 segundos máximo

---

## 9. Estrategia para nombres en castellano

### El problema

El usuario dice "producto más vendido" y el modelo necesita mapear a `products.name` + `sale_items.quantity` con `SUM` y `GROUP BY`.

### Solución: descripción semántica de tablas

Incluir en el system prompt un bloque de metadatos que describa cada tabla/columna en castellano:

```
TABLA: products (productos)
  - id: Identificador único
  - name: Nombre del producto (ej: "Coca-Cola 500ml")
  - price: Precio unitario en pesos
  - category_id: Categoría a la que pertenece
  - stock: Stock actual

TABLA: sales (ventas)
  - id: Número de venta
  - total: Monto total de la venta
  - payment_method: "cash" | "card" | "transfer" | "qr"
  - created_at: Fecha y hora ISO 8601

TABLA: sale_items (detalle de venta)
  - product_id: Qué producto se vendió (FK → products)
  - quantity: Cantidad vendida
  - unit_price: Precio al que se vendió (puede diferir del actual)

TABLA: categories (categorías de productos)
  - name: Nombre de la categoría (ej: "Bebidas", "Snacks")
```

### Terms lookup (diccionario de sinónimos)

Incluir en el system prompt un mini-glossary:

```
GLOSARIO DE TÉRMINOS (castellano → SQL):
  - producto, artículos → products
  - venta, ventas, cobro, factura → sales
  - ítem, detalle, producto vendido → sale_items
  - categoría, rubro, tipo → categories
  - precio, valor, importe, monto → price / total / unit_price
  - fecha, día, hoy, ayer → created_at
  - efectivo, cash → payment_method = 'cash'
  - tarjeta, crédito, débito → payment_method = 'card'
  - transferencia, transfer → payment_method = 'transfer'
  - QR, qr → payment_method = 'qr'
```

Los few-shot examples refuerzan este mapeo porque muestran traducciones concretas.

---

## Resumen de recomendaciones

1. **Modelo**: GPT-4o (primario), GPT-4o-mini (fallback barato) — Structured Outputs es el factor decisivo
2. **Parsing**: `response_format: json_schema` con strict mode — no regex, no parches
3. **Seguridad**: Validación con parser SQL en Go + conexión read-only + LIMIT 500
4. **Costos**: ~$30-50/mes para 30K consultas con estrategia híbrida
5. **Cache**: Map en memoria con TTL 5min, invalidación por schema version
6. **Idioma**: Describir tablas en castellano + glossary de sinónimos + few-shot examples bilingües
