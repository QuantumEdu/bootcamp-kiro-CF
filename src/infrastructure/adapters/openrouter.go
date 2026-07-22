package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenRouterClient handles communication with OpenRouter API.
type OpenRouterClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
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
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type chatMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type chatReq struct {
	Model       string    `json:"model"`
	Messages    []chatMsg `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	ResponseFormat *struct{ Type string `json:"type"` } `json:"response_format,omitempty"`
}
type chatResp struct {
	Choices []struct{ Message struct{ Content string `json:"content"` } `json:"message"` } `json:"choices"`
	Error   *struct{ Message string `json:"message"` }                                     `json:"error"`
}

// GenerateSQL sends a natural language query to OpenRouter and returns structured SQL response.
func (c *OpenRouterClient) GenerateSQL(ctx context.Context, userQuery, systemPrompt string) (*NLSQLResponse, error) {
	body := chatReq{
		Model:       c.model,
		Messages:    []chatMsg{{Role: "system", Content: systemPrompt}, {Role: "user", Content: userQuery}},
		MaxTokens:   300,
		Temperature: 0.1,
		ResponseFormat: &struct{ Type string `json:"type"` }{Type: "json_object"},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("HTTP-Referer", "https://pos-ai-first.local")
	req.Header.Set("X-Title", "POS AI-First")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API status %d: %s", resp.StatusCode, string(respBody))
	}

	var cr chatResp
	if err := json.Unmarshal(respBody, &cr); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	if cr.Error != nil {
		return nil, fmt.Errorf("API error: %s", cr.Error.Message)
	}
	if len(cr.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}

	var nlResp NLSQLResponse
	if err := json.Unmarshal([]byte(cr.Choices[0].Message.Content), &nlResp); err != nil {
		return nil, fmt.Errorf("parsing LLM JSON: %w (raw: %s)", err, cr.Choices[0].Message.Content)
	}
	return &nlResp, nil
}
