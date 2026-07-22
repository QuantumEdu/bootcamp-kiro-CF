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
	baseURL    string
}

// NLSQLResponse represents the structured response from the LLM.
type NLSQLResponse struct {
	SQL         *string `json:"sql"`
	Explanation string  `json:"explanation"`
	Error       *string `json:"error"`
}

// NewOpenRouterClient creates a new OpenRouter API client.
func NewOpenRouterClient(apiKey, model string) *OpenRouterClient {
	return &OpenRouterClient{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://openrouter.ai/api/v1/chat/completions",
	}
}

// chatMessage represents a single message in the conversation.
type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatRequest represents the OpenRouter API request body.
type chatRequest struct {
	Model          string        `json:"model"`
	Messages       []chatMessage `json:"messages"`
	MaxTokens      int           `json:"max_tokens"`
	Temperature    float64       `json:"temperature"`
	ResponseFormat *respFormat   `json:"response_format,omitempty"`
}

type respFormat struct {
	Type string `json:"type"`
}

// chatResponse represents the OpenRouter API response.
type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// GenerateSQL sends a natural language query to OpenRouter and returns the SQL.
func (c *OpenRouterClient) GenerateSQL(ctx context.Context, userQuery, systemPrompt string) (*NLSQLResponse, error) {
	reqBody := chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userQuery},
		},
		MaxTokens:   300,
		Temperature: 0.1,
		ResponseFormat: &respFormat{
			Type: "json_object",
		},
	}

	jsonBody, err := json.Marshal(reqBody)
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
		return nil, fmt.Errorf("calling OpenRouter API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenRouter API returned status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	if chatResp.Error != nil {
		return nil, fmt.Errorf("OpenRouter API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in OpenRouter response")
	}

	content := chatResp.Choices[0].Message.Content

	var nlResp NLSQLResponse
	if err := json.Unmarshal([]byte(content), &nlResp); err != nil {
		return nil, fmt.Errorf("parsing LLM JSON response: %w (raw: %s)", err, content)
	}

	return &nlResp, nil
}
