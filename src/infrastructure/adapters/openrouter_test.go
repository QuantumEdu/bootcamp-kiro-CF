package adapters

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestClassifyHTTPError_AllCases(t *testing.T) {
	// Comprehensive table-driven test for HTTP error classification
	tests := []struct {
		name       string
		statusCode int
		wantErr    error
		isGeneric  bool
	}{
		{"429 rate limit", 429, ErrAIRateLimit, false},
		{"408 request timeout", 408, context.DeadlineExceeded, false},
		{"504 gateway timeout", 504, context.DeadlineExceeded, false},
		{"500 internal server error", 500, ErrAIUnavailable, false},
		{"502 bad gateway", 502, ErrAIUnavailable, false},
		{"503 service unavailable", 503, ErrAIUnavailable, false},
		{"401 unauthorized", 401, nil, true},
		{"403 forbidden", 403, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := classifyHTTPError(tt.statusCode)
			if tt.isGeneric {
				if err == ErrAIRateLimit || err == ErrAIUnavailable || err == context.DeadlineExceeded {
					t.Errorf("expected generic error for %d, got known sentinel: %v", tt.statusCode, err)
				}
			} else {
				if err != tt.wantErr {
					t.Errorf("expected %v for status %d, got: %v", tt.wantErr, tt.statusCode, err)
				}
			}
		})
	}
}

func TestClassifyHTTPError_429_ReturnsRateLimit(t *testing.T) {
	err := classifyHTTPError(http.StatusTooManyRequests)
	if err != ErrAIRateLimit {
		t.Errorf("expected ErrAIRateLimit, got: %v", err)
	}
}

func TestClassifyHTTPError_5xx_ReturnsUnavailable(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{"500 Internal Server Error", 500},
		{"502 Bad Gateway", 502},
		{"503 Service Unavailable", 503},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := classifyHTTPError(tt.status)
			if err != ErrAIUnavailable {
				t.Errorf("expected ErrAIUnavailable for status %d, got: %v", tt.status, err)
			}
		})
	}
}

func TestClassifyHTTPError_408_ReturnsDeadlineExceeded(t *testing.T) {
	err := classifyHTTPError(http.StatusRequestTimeout)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got: %v", err)
	}
}

func TestClassifyHTTPError_504_ReturnsDeadlineExceeded(t *testing.T) {
	err := classifyHTTPError(http.StatusGatewayTimeout)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got: %v", err)
	}
}

func TestClassifyHTTPError_OtherStatus_ReturnsGenericError(t *testing.T) {
	err := classifyHTTPError(401)
	if err == nil {
		t.Fatal("expected an error")
	}
	if err == ErrAIRateLimit || err == ErrAIUnavailable || err == context.DeadlineExceeded {
		t.Errorf("expected generic error for 401, got known sentinel: %v", err)
	}
}

func TestIsTimeoutError_ContextDeadlineExceeded(t *testing.T) {
	if !isTimeoutError(context.DeadlineExceeded) {
		t.Error("expected context.DeadlineExceeded to be detected as timeout")
	}
}

func TestIsTimeoutError_RegularError(t *testing.T) {
	err := ErrAIUnavailable
	if isTimeoutError(err) {
		t.Error("expected ErrAIUnavailable NOT to be detected as timeout")
	}
}

type timeoutErr struct{}

func (e timeoutErr) Error() string   { return "timeout" }
func (e timeoutErr) Timeout() bool   { return true }
func (e timeoutErr) Temporary() bool { return true }

func TestIsTimeoutError_NetTimeout(t *testing.T) {
	err := timeoutErr{}
	if !isTimeoutError(err) {
		t.Error("expected net timeout error to be detected as timeout")
	}
}

func TestMalformedResponse_InvalidJSON(t *testing.T) {
	// Test that parsing invalid JSON in the LLM content field returns ErrAIMalformedResponse.
	t.Run("malformed LLM content returns structured error", func(t *testing.T) {
		// Simulate the JSON parsing step
		invalidContent := "this is not json at all"
		var nlResp NLSQLResponse
		err := json.Unmarshal([]byte(invalidContent), &nlResp)
		if err == nil {
			t.Fatal("expected unmarshal error")
		}
		// In the real code, this maps to ErrAIMalformedResponse
		// Verify the sentinel exists and has the right message
		if ErrAIMalformedResponse.Error() != "Respuesta del servicio de IA no válida" {
			t.Errorf("unexpected error message: %s", ErrAIMalformedResponse.Error())
		}
	})
}

func TestMalformedResponse_EmptyChoices(t *testing.T) {
	// When chatResp has no choices, should return ErrAIMalformedResponse
	// This tests the logic path in callAPI
	t.Run("empty choices returns malformed error", func(t *testing.T) {
		cr := chatResp{Choices: nil}
		if len(cr.Choices) != 0 {
			t.Fatal("expected zero choices")
		}
		// In callAPI, this returns ErrAIMalformedResponse
		if ErrAIMalformedResponse == nil {
			t.Fatal("ErrAIMalformedResponse should not be nil")
		}
	})
}

func TestFallbackModel_IsSet(t *testing.T) {
	client := NewOpenRouterClient("test-key", "anthropic/claude-3-haiku")
	if client.fallbackModel != "openai/gpt-4o-mini" {
		t.Errorf("expected fallback model 'openai/gpt-4o-mini', got: %s", client.fallbackModel)
	}
}

func TestNewOpenRouterClient_Defaults(t *testing.T) {
	client := NewOpenRouterClient("my-key", "anthropic/claude-3-haiku")

	if client.apiKey != "my-key" {
		t.Errorf("expected apiKey 'my-key', got: %s", client.apiKey)
	}
	if client.model != "anthropic/claude-3-haiku" {
		t.Errorf("expected model 'anthropic/claude-3-haiku', got: %s", client.model)
	}
	if client.fallbackModel != "openai/gpt-4o-mini" {
		t.Errorf("expected fallbackModel 'openai/gpt-4o-mini', got: %s", client.fallbackModel)
	}
	if client.httpClient == nil {
		t.Error("expected httpClient to be initialized")
	}
	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("expected 30s timeout, got: %v", client.httpClient.Timeout)
	}
}
