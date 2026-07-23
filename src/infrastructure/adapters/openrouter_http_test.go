package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newTestClient creates an OpenRouterClient pointed at the given mock server URL.
func newTestClient(serverURL string) *OpenRouterClient {
	return &OpenRouterClient{
		apiKey:        "test-api-key",
		model:         "test/primary-model",
		fallbackModel: "test/fallback-model",
		httpClient:    &http.Client{Timeout: 5 * time.Second},
		baseURL:       serverURL,
	}
}

// mockOpenRouterResponse builds a valid OpenRouter chat completion JSON response.
func mockOpenRouterResponse(content string) []byte {
	resp := chatResp{
		Choices: []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{
			{Message: struct {
				Content string `json:"content"`
			}{Content: content}},
		},
	}
	data, _ := json.Marshal(resp)
	return data
}

func TestGenerateSQL_Success_MockServer(t *testing.T) {
	// Mock returns a valid NLSQLResponse JSON inside the content field.
	sqlContent := `{"sql": "SELECT * FROM productos WHERE stock < 10", "explanation": "Productos con stock bajo"}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("expected Authorization header 'Bearer test-api-key', got: %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got: %s", r.Header.Get("Content-Type"))
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got: %s", r.Method)
		}

		// Decode request body to verify model and messages
		var reqBody chatReq
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if reqBody.Model != "test/primary-model" {
			t.Errorf("expected model 'test/primary-model', got: %s", reqBody.Model)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockOpenRouterResponse(sqlContent))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	resp, err := client.GenerateSQL(context.Background(), "¿qué productos tienen stock bajo?", "You are a SQL assistant")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.SQL == nil {
		t.Fatal("expected SQL field to be non-nil")
	}
	expectedSQL := "SELECT * FROM productos WHERE stock < 10"
	if *resp.SQL != expectedSQL {
		t.Errorf("expected SQL %q, got: %q", expectedSQL, *resp.SQL)
	}
	if resp.Explanation != "Productos con stock bajo" {
		t.Errorf("expected explanation 'Productos con stock bajo', got: %q", resp.Explanation)
	}
}

func TestGenerateSQL_Timeout_Fallback_MockServer(t *testing.T) {
	// First call times out (simulated by hanging), second call with fallback model succeeds.
	callCount := 0
	sqlContent := `{"sql": "SELECT COUNT(*) FROM ventas", "explanation": "Total de ventas"}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody chatReq
		json.NewDecoder(r.Body).Decode(&reqBody)
		callCount++

		if callCount == 1 {
			// First call: verify primary model was used, then simulate timeout
			if reqBody.Model != "test/primary-model" {
				t.Errorf("first call expected model 'test/primary-model', got: %s", reqBody.Model)
			}
			// Block until context is cancelled (simulates timeout)
			<-r.Context().Done()
			return
		}

		// Second call: verify fallback model was used
		if reqBody.Model != "test/fallback-model" {
			t.Errorf("second call expected model 'test/fallback-model', got: %s", reqBody.Model)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockOpenRouterResponse(sqlContent))
	}))
	defer server.Close()

	// Use a very short timeout to trigger timeout quickly
	client := &OpenRouterClient{
		apiKey:        "test-api-key",
		model:         "test/primary-model",
		fallbackModel: "test/fallback-model",
		httpClient:    &http.Client{Timeout: 100 * time.Millisecond},
		baseURL:       server.URL,
	}

	resp, err := client.GenerateSQL(context.Background(), "¿cuántas ventas hay?", "You are a SQL assistant")

	if err != nil {
		t.Fatalf("expected no error (fallback should succeed), got: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response from fallback")
	}
	if resp.SQL == nil {
		t.Fatal("expected SQL field to be non-nil")
	}
	if *resp.SQL != "SELECT COUNT(*) FROM ventas" {
		t.Errorf("expected SQL 'SELECT COUNT(*) FROM ventas', got: %q", *resp.SQL)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls (primary + fallback), got: %d", callCount)
	}
}

func TestGenerateSQL_RateLimit_MockServer(t *testing.T) {
	// Mock returns HTTP 429 (Too Many Requests).
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error": {"message": "rate limited"}}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	resp, err := client.GenerateSQL(context.Background(), "test query", "system prompt")

	if resp != nil {
		t.Errorf("expected nil response, got: %+v", resp)
	}
	if !errors.Is(err, ErrAIRateLimit) {
		t.Errorf("expected ErrAIRateLimit, got: %v", err)
	}
}

func TestGenerateSQL_ServerError_MockServer(t *testing.T) {
	// Mock returns HTTP 500 (Internal Server Error).
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": {"message": "internal error"}}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	resp, err := client.GenerateSQL(context.Background(), "test query", "system prompt")

	if resp != nil {
		t.Errorf("expected nil response, got: %+v", resp)
	}
	if !errors.Is(err, ErrAIUnavailable) {
		t.Errorf("expected ErrAIUnavailable, got: %v", err)
	}
}

func TestGenerateSQL_MalformedJSON_MockServer(t *testing.T) {
	// Mock returns HTTP 200 but content field has invalid JSON (not a valid NLSQLResponse).
	invalidContent := "this is not valid json at all"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockOpenRouterResponse(invalidContent))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	resp, err := client.GenerateSQL(context.Background(), "test query", "system prompt")

	if resp != nil {
		t.Errorf("expected nil response, got: %+v", resp)
	}
	if !errors.Is(err, ErrAIMalformedResponse) {
		t.Errorf("expected ErrAIMalformedResponse, got: %v", err)
	}
}
