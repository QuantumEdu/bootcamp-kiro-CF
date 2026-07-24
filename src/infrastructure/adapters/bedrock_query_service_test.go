package adapters

import (
	"context"
	"errors"
	"testing"
)

func TestParseNLSQLContent_ValidSQL(t *testing.T) {
	content := []anthropicContentBlock{
		{Type: "text", Text: `{"sql": "SELECT * FROM productos", "explanation": "Todos los productos", "error": null}`},
	}

	resp, err := parseNLSQLContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.SQL == nil || *resp.SQL != "SELECT * FROM productos" {
		t.Errorf("got SQL=%v, want 'SELECT * FROM productos'", resp.SQL)
	}
	if resp.Explanation != "Todos los productos" {
		t.Errorf("got Explanation=%q, want 'Todos los productos'", resp.Explanation)
	}
	if resp.Error != nil {
		t.Errorf("got Error=%v, want nil", resp.Error)
	}
}

func TestParseNLSQLContent_ErrorResponse(t *testing.T) {
	content := []anthropicContentBlock{
		{Type: "text", Text: `{"sql": null, "explanation": "No puedo hacer eso", "error": "consulta no soportada"}`},
	}

	resp, err := parseNLSQLContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.SQL != nil {
		t.Errorf("got SQL=%v, want nil", resp.SQL)
	}
	if resp.Error == nil || *resp.Error != "consulta no soportada" {
		t.Errorf("got Error=%v, want 'consulta no soportada'", resp.Error)
	}
}

func TestParseNLSQLContent_EmptyContent(t *testing.T) {
	_, err := parseNLSQLContent([]anthropicContentBlock{})
	if !errors.Is(err, ErrAIMalformedResponse) {
		t.Errorf("got err=%v, want ErrAIMalformedResponse", err)
	}
}

func TestParseNLSQLContent_NonTextBlocks(t *testing.T) {
	content := []anthropicContentBlock{
		{Type: "image", Text: ""},
	}

	_, err := parseNLSQLContent(content)
	if !errors.Is(err, ErrAIMalformedResponse) {
		t.Errorf("got err=%v, want ErrAIMalformedResponse", err)
	}
}

func TestParseNLSQLContent_InvalidJSON(t *testing.T) {
	content := []anthropicContentBlock{
		{Type: "text", Text: "this is not json"},
	}

	_, err := parseNLSQLContent(content)
	if !errors.Is(err, ErrAIMalformedResponse) {
		t.Errorf("got err=%v, want ErrAIMalformedResponse", err)
	}
}

func TestParseNLSQLContent_WhitespaceWrapped(t *testing.T) {
	content := []anthropicContentBlock{
		{Type: "text", Text: `  {"sql": "SELECT 1", "explanation": "test", "error": null}  `},
	}

	resp, err := parseNLSQLContent(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.SQL == nil || *resp.SQL != "SELECT 1" {
		t.Errorf("got SQL=%v, want 'SELECT 1'", resp.SQL)
	}
}

func TestSanitizeAWSError_Throttling(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr error
	}{
		{"throttling exception", errors.New("ThrottlingException: rate exceeded"), ErrAIRateLimit},
		{"too many requests", errors.New("TooManyRequestsException: slow down"), ErrAIRateLimit},
		{"access denied", errors.New("AccessDeniedException: not authorized"), ErrAIUnavailable},
		{"resource not found", errors.New("ResourceNotFoundException: model not found"), ErrAIUnavailable},
		{"validation exception", errors.New("ValidationException: invalid input"), ErrAIUnavailable},
		{"context deadline", context.DeadlineExceeded, ErrAIUnavailable},
		{"context canceled", context.Canceled, ErrAIUnavailable},
		{"generic error", errors.New("something went wrong"), ErrAIUnavailable},
		{"nil error", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeAWSError(tt.err)
			if !errors.Is(got, tt.wantErr) {
				t.Errorf("sanitizeAWSError(%v) = %v, want %v", tt.err, got, tt.wantErr)
			}
		})
	}
}

func TestSanitizeAWSError_NoLeakedDetails(t *testing.T) {
	// Simulate a real AWS error with internal details
	awsErr := errors.New("operation error Bedrock Runtime: InvokeModel, https response error StatusCode: 403, RequestID: abc123-def456, api error AccessDeniedException: User: arn:aws:iam::123456789012:role/my-role is not authorized")

	sanitized := sanitizeAWSError(awsErr)

	// Must not contain ARNs or request IDs
	errMsg := sanitized.Error()
	if contains(errMsg, "arn:aws") {
		t.Errorf("sanitized error leaks ARN: %s", errMsg)
	}
	if contains(errMsg, "abc123") {
		t.Errorf("sanitized error leaks request ID: %s", errMsg)
	}
	if contains(errMsg, "123456789012") {
		t.Errorf("sanitized error leaks account ID: %s", errMsg)
	}

	// Should be a domain-friendly error
	if !errors.Is(sanitized, ErrAIUnavailable) {
		t.Errorf("expected ErrAIUnavailable, got %v", sanitized)
	}
}

func TestContainsAny(t *testing.T) {
	tests := []struct {
		s      string
		subs   []string
		expect bool
	}{
		{"ThrottlingException: too fast", []string{"ThrottlingException"}, true},
		{"some error", []string{"ThrottlingException", "TooManyRequests"}, false},
		{"THROTTLINGEXCEPTION", []string{"throttlingexception"}, true}, // case insensitive
		{"", []string{"anything"}, false},
	}

	for _, tt := range tests {
		got := containsAny(tt.s, tt.subs...)
		if got != tt.expect {
			t.Errorf("containsAny(%q, %v) = %v, want %v", tt.s, tt.subs, got, tt.expect)
		}
	}
}

func TestNewBedrockQueryService_Defaults(t *testing.T) {
	svc := NewBedrockQueryService(nil, BedrockConfig{
		ModelID: "anthropic.claude-3-haiku-20240307-v1:0",
		Region:  "us-east-1",
	}, "test schema")

	if svc.config.MaxTokens != 300 {
		t.Errorf("expected default MaxTokens=300, got %d", svc.config.MaxTokens)
	}
	if svc.config.Temperature != 0.1 {
		t.Errorf("expected default Temperature=0.1, got %f", svc.config.Temperature)
	}
}

func TestNewBedrockQueryService_CustomConfig(t *testing.T) {
	svc := NewBedrockQueryService(nil, BedrockConfig{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		Region:      "us-west-2",
		MaxTokens:   500,
		Temperature: 0.5,
	}, "custom schema")

	if svc.config.MaxTokens != 500 {
		t.Errorf("expected MaxTokens=500, got %d", svc.config.MaxTokens)
	}
	if svc.config.Temperature != 0.5 {
		t.Errorf("expected Temperature=0.5, got %f", svc.config.Temperature)
	}
	if svc.schema != "custom schema" {
		t.Errorf("expected schema='custom schema', got %q", svc.schema)
	}
}

func TestBuildSystemPrompt_ContainsSchema(t *testing.T) {
	svc := NewBedrockQueryService(nil, BedrockConfig{
		ModelID: "test-model",
	}, "CREATE TABLE productos (...)")

	prompt := svc.buildSystemPrompt()

	if !contains(prompt, "CREATE TABLE productos") {
		t.Error("system prompt should contain the schema")
	}
	if !contains(prompt, "PostgreSQL") {
		t.Error("system prompt should reference PostgreSQL")
	}
	if !contains(prompt, "SELECT") {
		t.Error("system prompt should mention SELECT-only rule")
	}
}

// contains is a helper for string-contains checks in tests.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
