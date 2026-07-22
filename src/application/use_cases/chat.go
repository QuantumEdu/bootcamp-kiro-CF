package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/adapters"
)

// ChatService handles natural language to SQL queries.
type ChatService struct {
	openRouter *adapters.OpenRouterClient
	readDB     *sql.DB
	schema     string
	timeout    time.Duration
}

// ChatResult represents the result of a chat query.
type ChatResult struct {
	Query       string
	SQL         string
	Explanation string
	Columns     []string
	Results     [][]string
	Error       string
}

// NewChatService creates a new chat service.
func NewChatService(openRouter *adapters.OpenRouterClient, readDB *sql.DB, schema string, timeoutSecs int) *ChatService {
	return &ChatService{
		openRouter: openRouter,
		readDB:     readDB,
		schema:     schema,
		timeout:    time.Duration(timeoutSecs) * time.Second,
	}
}

// ProcessQuery processes a natural language query end-to-end.
func (s *ChatService) ProcessQuery(ctx context.Context, userQuery string) *ChatResult {
	result := &ChatResult{Query: userQuery}

	// Generate SQL from natural language
	systemPrompt := s.buildSystemPrompt()
	nlResp, err := s.openRouter.GenerateSQL(ctx, userQuery, systemPrompt)
	if err != nil {
		result.Error = fmt.Sprintf("Error al procesar la consulta: %v", err)
		return result
	}

	if nlResp.Error != nil {
		result.Error = *nlResp.Error
		result.Explanation = nlResp.Explanation
		return result
	}

	if nlResp.SQL == nil {
		result.Error = "No se pudo generar una consulta SQL valida"
		result.Explanation = nlResp.Explanation
		return result
	}

	generatedSQL := *nlResp.SQL
	result.SQL = generatedSQL
	result.Explanation = nlResp.Explanation

	// Validate SQL
	if err := ValidateSQL(generatedSQL); err != nil {
		result.Error = fmt.Sprintf("La consulta generada no es segura: %v", err)
		return result
	}

	// Add LIMIT if not present
	if !strings.Contains(strings.ToUpper(generatedSQL), "LIMIT") {
		generatedSQL += " LIMIT 100"
	}

	// Execute SQL with timeout
	queryCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	columns, rows, err := s.executeQuery(queryCtx, generatedSQL)
	if err != nil {
		result.Error = fmt.Sprintf("Error al ejecutar la consulta: %v", err)
		return result
	}

	result.Columns = columns
	result.Results = rows
	return result
}

func (s *ChatService) executeQuery(ctx context.Context, sqlQuery string) ([]string, [][]string, error) {
	rows, err := s.readDB.QueryContext(ctx, sqlQuery)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	var results [][]string
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, nil, err
		}

		row := make([]string, len(columns))
		for i, v := range values {
			if v == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", v)
			}
		}
		results = append(results, row)
	}

	return columns, results, rows.Err()
}

func (s *ChatService) buildSystemPrompt() string {
	return fmt.Sprintf(`Eres un asistente que convierte lenguaje natural a SQL para una base de datos SQLite de un sistema POS (Point of Sale).

REGLAS ESTRICTAS:
1. Solo genera sentencias SELECT. NUNCA generes INSERT, UPDATE, DELETE, DROP, ALTER, CREATE, o cualquier DDL/DML.
2. Usa solo la sintaxis de SQLite.
3. Si no puedes generar una query segura con la informacion dada, responde con:
   {"sql": null, "error": "Descripcion clara de por que no se puede generar la query", "explanation": "..."}

ESQUEMA DE BASE DE DATOS:
%s

GLOSARIO DE TERMINOS (castellano -> SQL):
- producto, articulos -> productos
- venta, ventas, cobro, factura -> ventas
- item, detalle, producto vendido -> venta_items
- categoria, rubro, tipo -> categorias
- precio, valor, importe, monto -> precio_venta / total / precio_unitario
- fecha, dia, hoy, ayer -> created_at / fecha
- efectivo, cash -> metodo_pago = 'efectivo'
- tarjeta, credito, debito -> metodo_pago = 'tarjeta'
- transferencia, transfer -> metodo_pago = 'transferencia'
- stock, inventario, existencia -> stock_actual
- cliente, comprador -> clientes

EJEMPLOS:
Usuario: "cuantos productos vendi esta semana?"
{"sql": "SELECT COUNT(DISTINCT vi.producto_id) FROM venta_items vi JOIN ventas v ON vi.venta_id = v.id WHERE v.created_at >= datetime('now', '-7 days')", "explanation": "Cuenta productos distintos vendidos en los ultimos 7 dias", "error": null}

Usuario: "que producto se vendio mas ayer"
{"sql": "SELECT p.nombre, SUM(vi.cantidad) as total_qty FROM venta_items vi JOIN productos p ON p.id = vi.producto_id JOIN ventas v ON v.id = vi.venta_id WHERE date(v.created_at) = date('now', '-1 day') GROUP BY p.id ORDER BY total_qty DESC LIMIT 1", "explanation": "Agrupa por producto, suma cantidades, ordena descendente y trae el primero", "error": null}

Usuario: "mostrame las ventas de hoy"
{"sql": "SELECT v.id, v.total, v.metodo_pago, v.created_at FROM ventas v WHERE date(v.created_at) = date('now') ORDER BY v.created_at DESC", "explanation": "Filtra ventas del dia actual ordenadas por fecha descendente", "error": null}

Usuario: "cuanto dinero en efectivo se cobro este mes"
{"sql": "SELECT COALESCE(SUM(total), 0) as total_efectivo FROM ventas WHERE metodo_pago = 'efectivo' AND strftime('%%Y-%%m', created_at) = strftime('%%Y-%%m', 'now')", "explanation": "Suma total de ventas donde el metodo de pago es 'efectivo' en el mes actual", "error": null}

Devuelve SIEMPRE un JSON valido con esta estructura exacta:
{"sql": "SELECT ...", "explanation": "Breve explicacion", "error": null}`, s.schema)
}
