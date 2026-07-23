package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ErrAIUnavailable indicates the AI service is temporarily unavailable.
var ErrAIUnavailable = errors.New("Servicio de IA temporalmente no disponible")

// ErrAIRateLimit indicates too many requests to the AI service.
var ErrAIRateLimit = errors.New("Demasiadas consultas, intenta en unos segundos")

// ErrAIMalformedResponse indicates the AI returned an unparseable response.
var ErrAIMalformedResponse = errors.New("Respuesta del servicio de IA no válida")

// OpenRouterClient handles communication with OpenRouter API.
type OpenRouterClient struct {
	apiKey        string
	model         string
	fallbackModel string
	httpClient    *http.Client
	baseURL       string
}

// NLSQLResponse represents the structured response from the LLM.
type NLSQLResponse struct {
	SQL         *string `json:"sql"`
	Explanation string  `json:"explanation"`
	Error       *string `json:"error"`
}

// NewOpenRouterClient creates a new OpenRouter client.
func NewOpenRouterClient(apiKey, model string) *OpenRouterClient {
	return &OpenRouterClient{
		apiKey:        apiKey,
		model:         model,
		fallbackModel: "openai/gpt-4o-mini",
		httpClient:    &http.Client{Timeout: 30 * time.Second},
		baseURL:       "https://openrouter.ai/api/v1/chat/completions",
	}
}

type chatMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type chatReq struct {
	Model          string    `json:"model"`
	Messages       []chatMsg `json:"messages"`
	MaxTokens      int       `json:"max_tokens"`
	Temperature    float64   `json:"temperature"`
	ResponseFormat *struct {
		Type string `json:"type"`
	} `json:"response_format,omitempty"`
}
type chatResp struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// GenerateSQL sends a natural language query to OpenRouter and returns structured SQL response.
func (c *OpenRouterClient) GenerateSQL(ctx context.Context, userQuery, systemPrompt string) (*NLSQLResponse, error) {
	resp, err := c.callAPI(ctx, c.model, userQuery, systemPrompt)
	if err != nil {
		// On timeout, retry once with fallback model
		if isTimeoutError(err) {
			resp, err = c.callAPI(ctx, c.fallbackModel, userQuery, systemPrompt)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}
		return nil, err
	}
	return resp, nil
}

// callAPI performs the actual HTTP call to OpenRouter with the given model.
func (c *OpenRouterClient) callAPI(ctx context.Context, model, userQuery, systemPrompt string) (*NLSQLResponse, error) {
	body := chatReq{
		Model:       model,
		Messages:    []chatMsg{{Role: "system", Content: systemPrompt}, {Role: "user", Content: userQuery}},
		MaxTokens:   300,
		Temperature: 0.1,
		ResponseFormat: &struct {
			Type string `json:"type"`
		}{Type: "json_object"},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("HTTP-Referer", "https://pos-ai-first.local")
	req.Header.Set("X-Title", "POS AI-First")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if isTimeoutError(err) {
			return nil, err
		}
		return nil, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	// Handle HTTP error status codes with domain-friendly messages
	if resp.StatusCode != http.StatusOK {
		return nil, classifyHTTPError(resp.StatusCode)
	}

	var cr chatResp
	if err := json.Unmarshal(respBody, &cr); err != nil {
		return nil, ErrAIMalformedResponse
	}
	if cr.Error != nil {
		return nil, fmt.Errorf("API error: %s", cr.Error.Message)
	}
	if len(cr.Choices) == 0 {
		return nil, ErrAIMalformedResponse
	}

	var nlResp NLSQLResponse
	if err := json.Unmarshal([]byte(cr.Choices[0].Message.Content), &nlResp); err != nil {
		return nil, ErrAIMalformedResponse
	}
	return &nlResp, nil
}

// classifyHTTPError maps HTTP status codes to domain-friendly errors.
func classifyHTTPError(statusCode int) error {
	switch {
	case statusCode == http.StatusTooManyRequests: // 429
		return ErrAIRateLimit
	case statusCode == http.StatusRequestTimeout || statusCode == http.StatusGatewayTimeout: // 408, 504
		return context.DeadlineExceeded
	case statusCode >= 500 && statusCode < 600: // 5xx
		return ErrAIUnavailable
	default:
		return fmt.Errorf("API status %d", statusCode)
	}
}

// isTimeoutError checks if an error is a timeout (context deadline exceeded or HTTP timeout).
func isTimeoutError(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	// Check for net/http timeout errors
	var netErr interface{ Timeout() bool }
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}
	return false
}
