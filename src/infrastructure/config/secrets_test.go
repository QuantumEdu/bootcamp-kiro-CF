package config

import (
	"context"
	"encoding/json"
	"testing"
)

func TestAISecretPayload_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    AISecretPayload
		wantErr bool
	}{
		{
			name:  "valid payload",
			input: `{"model_id":"anthropic.claude-3-haiku-20240307-v1:0","region":"us-east-1","max_tokens":300,"temperature":0.1}`,
			want: AISecretPayload{
				ModelID:     "anthropic.claude-3-haiku-20240307-v1:0",
				Region:      "us-east-1",
				MaxTokens:   300,
				Temperature: 0.1,
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `not json at all`,
			want:    AISecretPayload{},
			wantErr: true,
		},
		{
			name:  "empty fields",
			input: `{}`,
			want: AISecretPayload{
				ModelID:     "",
				Region:      "",
				MaxTokens:   0,
				Temperature: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got AISecretPayload
			err := json.Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("got err=%v, wantErr=%v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestSecretsLoader_LoadConfig_MissingARNs(t *testing.T) {
	// Ensure environment is clean for this test.
	t.Setenv("SECRET_DB_ARN", "")
	t.Setenv("SECRET_SESSION_ARN", "")
	t.Setenv("SECRET_AI_ARN", "")

	loader := NewSecretsLoader(nil)
	_, err := loader.LoadConfig(context.Background())
	if err == nil {
		t.Fatal("expected error for missing ARNs, got nil")
	}

	// Verify error mentions all three missing vars.
	errMsg := err.Error()
	for _, varName := range []string{"SECRET_DB_ARN", "SECRET_SESSION_ARN", "SECRET_AI_ARN"} {
		if !contains(errMsg, varName) {
			t.Errorf("error message should mention %s, got: %s", varName, errMsg)
		}
	}
}

func TestSecretsLoader_LoadConfig_PartialMissing(t *testing.T) {
	t.Setenv("SECRET_DB_ARN", "arn:aws:secretsmanager:us-east-1:123:secret:db")
	t.Setenv("SECRET_SESSION_ARN", "")
	t.Setenv("SECRET_AI_ARN", "")

	loader := NewSecretsLoader(nil)
	_, err := loader.LoadConfig(context.Background())
	if err == nil {
		t.Fatal("expected error for partial missing ARNs, got nil")
	}

	errMsg := err.Error()
	if contains(errMsg, "SECRET_DB_ARN") {
		t.Errorf("error should NOT mention SECRET_DB_ARN (it's set), got: %s", errMsg)
	}
	if !contains(errMsg, "SECRET_SESSION_ARN") {
		t.Errorf("error should mention SECRET_SESSION_ARN, got: %s", errMsg)
	}
	if !contains(errMsg, "SECRET_AI_ARN") {
		t.Errorf("error should mention SECRET_AI_ARN, got: %s", errMsg)
	}
}

func TestSecretsLoader_GetSecret_CacheHit(t *testing.T) {
	// Pre-populate cache directly to test read-through behavior without AWS client.
	loader := NewSecretsLoader(nil)
	loader.cache["arn:test"] = "cached-value"

	val, err := loader.GetSecret(context.Background(), "arn:test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "cached-value" {
		t.Errorf("got %q, want %q", val, "cached-value")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
