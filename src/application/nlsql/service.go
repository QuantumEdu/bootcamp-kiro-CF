package nlsql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/QuantumEdu/bootcamp-kiro-CF/src/infrastructure/adapters"
)

// ChatResult represents the result of a chat query.
type ChatResult struct {
	Query         string
	SQL           string
	Explanation   string
	Columns       []string
	Results       [][]string
	FormattedText string
	Error         string
}

// Service handles NL→SQL query processing.
type Service struct {
	openRouter *adapters.OpenRouterClient
	readDB     *sql.DB
	schema     string
	timeout    time.Duration
	logger     *QueryLogger
}

// NewService creates a new NL→SQL service.
func NewService(openRouter *adapters.OpenRouterClient, readDB *sql.DB, schema string, timeoutSecs int) *Service {
	return &Service{
		openRouter: openRouter,
		readDB:     readDB,
		schema:     schema,
		timeout:    time.Duration(timeoutSecs) * time.Second,
	}
}

// SetLogger attaches a QueryLogger to the service. Nil-safe: if not set, logging is skipped.
func (s *Service) SetLogger(logger *QueryLogger) {
	s.logger = logger
}

// ProcessQuery processes a natural language query end-to-end.
func (s *Service) ProcessQuery(ctx context.Context, userQuery string) *ChatResult {
	start := time.Now()
	result := &ChatResult{Query: userQuery}

	if err := ValidateUserInput(userQuery); err != nil {
		result.Error = err.Error()
		s.logQuery(ctx, userQuery, "", false, err.Error(), time.Since(start))
		return result
	}

	nlResp, err := s.openRouter.GenerateSQL(ctx, userQuery, s.buildSystemPrompt())
	if err != nil {
		result.Error = fmt.Sprintf("Error al procesar: %v", err)
		s.logQuery(ctx, userQuery, "", false, result.Error, time.Since(start))
		return result
	}
	if nlResp.Error != nil {
		result.Error = *nlResp.Error
		result.Explanation = nlResp.Explanation
		s.logQuery(ctx, userQuery, "", false, result.Error, time.Since(start))
		return result
	}
	if nlResp.SQL == nil {
		result.Error = "No se pudo generar SQL"
		result.Explanation = nlResp.Explanation
		s.logQuery(ctx, userQuery, "", false, result.Error, time.Since(start))
		return result
	}

	generatedSQL := *nlResp.SQL
	result.SQL = generatedSQL
	result.Explanation = nlResp.Explanation

	if err := ValidateSQL(generatedSQL); err != nil {
		result.Error = fmt.Sprintf("Consulta insegura: %v", err)
		s.logQuery(ctx, userQuery, generatedSQL, false, result.Error, time.Since(start))
		return result
	}

	if !strings.Contains(strings.ToUpper(generatedSQL), "LIMIT") {
		generatedSQL += " LIMIT 100"
	}

	queryCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	columns, rows, err := s.executeQuery(queryCtx, generatedSQL)
	if err != nil {
		result.Error = fmt.Sprintf("Error ejecutando: %v", err)
		s.logQuery(ctx, userQuery, generatedSQL, false, result.Error, time.Since(start))
		return result
	}
	result.Columns = columns
	result.Results = rows
	result.FormattedText = FormatResults(columns, rows)

	s.logQuery(ctx, userQuery, generatedSQL, true, "", time.Since(start))
	return result
}

// logQuery writes an audit entry via the logger. It is nil-safe.
func (s *Service) logQuery(ctx context.Context, question, generatedSQL string, success bool, errMsg string, elapsed time.Duration) {
	if s.logger == nil {
		return
	}

	entry := QueryLogEntry{
		Question:        question,
		GeneratedSQL:    generatedSQL,
		Success:         success,
		ErrorMessage:    errMsg,
		ExecutionTimeMs: elapsed.Milliseconds(),
	}

	// Extract user_id from context if available.
	if userID, ok := ctx.Value(ContextKeyUserID).(int64); ok {
		entry.UserID = &userID
	}

	// Logging is best-effort; don't fail the request on log errors.
	_ = s.logger.Log(ctx, entry)
}

// contextKey is a type for context keys in this package.
type contextKey string

// ContextKeyUserID is the context key for the user ID.
const ContextKeyUserID = contextKey("user_id")

func (s *Service) executeQuery(ctx context.Context, sqlQuery string) ([]string, [][]string, error) {
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
		if len(results) >= 100 {
			break
		}
		values := make([]interface{}, len(columns))
		ptrs := make([]interface{}, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
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

func (s *Service) buildSystemPrompt() string {
	return fmt.Sprintf(`Eres un asistente que convierte lenguaje natural a SQL para SQLite (sistema POS).

REGLAS:
1. Solo genera SELECT. NUNCA INSERT/UPDATE/DELETE/DROP/ALTER/CREATE.
2. Sintaxis SQLite.
3. Si no puedes, responde: {"sql": null, "error": "motivo", "explanation": "..."}

SCHEMA:
%s

GLOSARIO:
- producto/articulo -> productos
- venta/cobro/factura -> ventas
- item/detalle -> venta_items
- categoria/rubro -> categorias
- precio/valor -> precio_venta / total
- hoy/ayer/fecha -> created_at
- efectivo/cash -> metodo_pago = 'efectivo'
- tarjeta -> metodo_pago = 'tarjeta'
- stock/inventario -> stock_actual
- cliente/comprador -> clientes

EJEMPLOS:
User: "cuantos productos vendi esta semana?"
{"sql": "SELECT COUNT(DISTINCT vi.producto_id) FROM venta_items vi JOIN ventas v ON vi.venta_id = v.id WHERE v.created_at >= datetime('now', '-7 days')", "explanation": "Productos distintos vendidos en 7 dias", "error": null}

User: "mostrame las ventas de hoy"
{"sql": "SELECT v.id, v.total, v.metodo_pago, v.created_at FROM ventas v WHERE date(v.created_at) = date('now') ORDER BY v.created_at DESC", "explanation": "Ventas del dia actual", "error": null}

Responde SIEMPRE JSON: {"sql": "SELECT ...", "explanation": "...", "error": null}`, s.schema)
}
