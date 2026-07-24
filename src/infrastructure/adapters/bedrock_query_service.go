package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// BedrockConfig holds configuration for the Bedrock adapter.
type BedrockConfig struct {
	ModelID     string  // e.g., "anthropic.claude-3-haiku-20240307-v1:0"
	Region      string  // e.g., "us-east-1"
	MaxTokens   int     // e.g., 300
	Temperature float64 // e.g., 0.1
}

// BedrockQueryService implements ports.AIQueryService using Amazon Bedrock.
type BedrockQueryService struct {
	client *bedrockruntime.Client
	config BedrockConfig
	schema string
}

// NewBedrockQueryService creates a new Bedrock query service.
// The schema parameter is the database schema string used to build the system prompt.
func NewBedrockQueryService(client *bedrockruntime.Client, cfg BedrockConfig, schema string) *BedrockQueryService {
	if cfg.MaxTokens == 0 {
		cfg.MaxTokens = 300
	}
	if cfg.Temperature == 0 {
		cfg.Temperature = 0.1
	}
	return &BedrockQueryService{client: client, config: cfg, schema: schema}
}

// GenerateSQL implements ports.AIQueryService.
func (s *BedrockQueryService) GenerateSQL(ctx context.Context, question string) (string, string, error) {
	payload := anthropicRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        s.config.MaxTokens,
		Temperature:      s.config.Temperature,
		System:           s.buildSystemPrompt(),
		Messages: []anthropicMessage{
			{Role: "user", Content: question},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("marshaling bedrock request: %w", err)
	}

	modelID := s.config.ModelID
	contentType := "application/json"
	accept := "application/json"

	resp, err := s.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     &modelID,
		ContentType: &contentType,
		Accept:      &accept,
		Body:        body,
	})
	if err != nil {
		return "", "", sanitizeAWSError(err)
	}

	var result anthropicResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return "", "", ErrAIMalformedResponse
	}

	nlResp, err := parseNLSQLContent(result.Content)
	if err != nil {
		return "", "", err
	}

	if nlResp.Error != nil {
		return "", nlResp.Explanation, fmt.Errorf("%s", *nlResp.Error)
	}
	if nlResp.SQL == nil {
		return "", nlResp.Explanation, fmt.Errorf("no SQL generated")
	}

	return *nlResp.SQL, nlResp.Explanation, nil
}

// buildSystemPrompt constructs the system prompt for NL-to-SQL translation.
// This uses the same prompt structure as the OpenRouter adapter.
func (s *BedrockQueryService) buildSystemPrompt() string {
	return fmt.Sprintf(`Eres un asistente que convierte lenguaje natural a SQL para PostgreSQL (sistema POS).

REGLAS:
1. Solo genera SELECT. NUNCA INSERT/UPDATE/DELETE/DROP/ALTER/CREATE.
2. Sintaxis PostgreSQL.
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
{"sql": "SELECT COUNT(DISTINCT vi.producto_id) FROM venta_items vi JOIN ventas v ON vi.venta_id = v.id WHERE v.created_at >= NOW() - INTERVAL '7 days'", "explanation": "Productos distintos vendidos en 7 dias", "error": null}

User: "mostrame las ventas de hoy"
{"sql": "SELECT v.id, v.total, v.metodo_pago, v.created_at FROM ventas v WHERE v.created_at >= CURRENT_DATE ORDER BY v.created_at DESC", "explanation": "Ventas del dia actual", "error": null}

Responde SIEMPRE JSON: {"sql": "SELECT ...", "explanation": "...", "error": null}`, s.schema)
}

// sanitizeAWSError maps AWS errors to domain-friendly errors without exposing
// internal details like ARNs, request IDs, or account information.
func sanitizeAWSError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Check for throttling/rate limit errors
	if containsAny(errMsg, "ThrottlingException", "TooManyRequestsException", "throttl") {
		return ErrAIRateLimit
	}

	// Check for model not found or access denied (treat as unavailable)
	if containsAny(errMsg, "AccessDeniedException", "ResourceNotFoundException", "ValidationException") {
		return ErrAIUnavailable
	}

	// Check for timeout-related errors
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return ErrAIUnavailable
	}

	// Default: service unavailable (don't leak internal error details)
	return ErrAIUnavailable
}

// parseNLSQLContent extracts the NLSQLResponse from an Anthropic Messages API response.
func parseNLSQLContent(content []anthropicContentBlock) (*NLSQLResponse, error) {
	if len(content) == 0 {
		return nil, ErrAIMalformedResponse
	}

	// Find the first text content block
	for _, block := range content {
		if block.Type != "text" {
			continue
		}

		text := strings.TrimSpace(block.Text)
		if text == "" {
			continue
		}

		var nlResp NLSQLResponse
		if err := json.Unmarshal([]byte(text), &nlResp); err != nil {
			return nil, ErrAIMalformedResponse
		}
		return &nlResp, nil
	}

	return nil, ErrAIMalformedResponse
}

// containsAny checks if s contains any of the given substrings (case-insensitive).
func containsAny(s string, substrs ...string) bool {
	lower := strings.ToLower(s)
	for _, sub := range substrs {
		if strings.Contains(lower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}

// --- Anthropic Messages API types ---

type anthropicRequest struct {
	AnthropicVersion string             `json:"anthropic_version"`
	MaxTokens        int                `json:"max_tokens"`
	Temperature      float64            `json:"temperature"`
	System           string             `json:"system"`
	Messages         []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content    []anthropicContentBlock `json:"content"`
	StopReason string                  `json:"stop_reason"`
}

type anthropicContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
