package bootstrap

import (
	"strings"
	"testing"
)

func TestValidateConfig_LambdaMode_AllPresent(t *testing.T) {
	cfg := &Config{
		AppEnv:         "lambda",
		SessionSecret:  "secret123",
		DatabaseURL:    "postgres://localhost/pos",
		BedrockModelID: "anthropic.claude-3-haiku-20240307-v1:0",
	}
	if err := ValidateConfig(cfg); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateConfig_LocalMode_AllPresent(t *testing.T) {
	cfg := &Config{
		AppEnv:        "local",
		SessionSecret: "secret123",
		DatabasePath:  "./data/pos.db",
	}
	if err := ValidateConfig(cfg); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateConfig_LocalMode_DefaultAppEnv(t *testing.T) {
	cfg := &Config{
		AppEnv:        "",
		SessionSecret: "secret123",
		DatabasePath:  "./data/pos.db",
	}
	if err := ValidateConfig(cfg); err != nil {
		t.Errorf("expected no error for unset AppEnv, got: %v", err)
	}
}

func TestValidateConfig_LambdaMode_MissingFields(t *testing.T) {
	tests := []struct {
		name           string
		cfg            *Config
		wantSubstrings []string
	}{
		{
			name: "missing all lambda fields",
			cfg: &Config{
				AppEnv: "lambda",
			},
			wantSubstrings: []string{"SESSION_SECRET", "DATABASE_URL", "BEDROCK_MODEL_ID"},
		},
		{
			name: "missing DatabaseURL only",
			cfg: &Config{
				AppEnv:         "lambda",
				SessionSecret:  "secret",
				BedrockModelID: "model-id",
			},
			wantSubstrings: []string{"DATABASE_URL"},
		},
		{
			name: "missing BedrockModelID only",
			cfg: &Config{
				AppEnv:        "lambda",
				SessionSecret: "secret",
				DatabaseURL:   "postgres://localhost/pos",
			},
			wantSubstrings: []string{"BEDROCK_MODEL_ID"},
		},
		{
			name: "missing SessionSecret only",
			cfg: &Config{
				AppEnv:         "lambda",
				DatabaseURL:    "postgres://localhost/pos",
				BedrockModelID: "model-id",
			},
			wantSubstrings: []string{"SESSION_SECRET"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.cfg)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			errMsg := err.Error()
			for _, sub := range tt.wantSubstrings {
				if !strings.Contains(errMsg, sub) {
					t.Errorf("error %q should mention %q", errMsg, sub)
				}
			}
		})
	}
}

func TestValidateConfig_LocalMode_MissingFields(t *testing.T) {
	tests := []struct {
		name           string
		cfg            *Config
		wantSubstrings []string
	}{
		{
			name: "missing all local fields",
			cfg: &Config{
				AppEnv: "local",
			},
			wantSubstrings: []string{"SESSION_SECRET", "DATABASE_PATH"},
		},
		{
			name: "missing DatabasePath only",
			cfg: &Config{
				AppEnv:        "local",
				SessionSecret: "secret",
			},
			wantSubstrings: []string{"DATABASE_PATH"},
		},
		{
			name: "missing SessionSecret only",
			cfg: &Config{
				AppEnv:       "local",
				DatabasePath: "./data/pos.db",
			},
			wantSubstrings: []string{"SESSION_SECRET"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.cfg)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			errMsg := err.Error()
			for _, sub := range tt.wantSubstrings {
				if !strings.Contains(errMsg, sub) {
					t.Errorf("error %q should mention %q", errMsg, sub)
				}
			}
		})
	}
}

func TestValidateConfig_ErrorListsAllMissingFields(t *testing.T) {
	cfg := &Config{
		AppEnv: "lambda",
	}
	err := ValidateConfig(cfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	errMsg := err.Error()
	// Verify the error mentions all three missing fields
	expected := []string{"SESSION_SECRET", "DATABASE_URL", "BEDROCK_MODEL_ID"}
	for _, field := range expected {
		if !strings.Contains(errMsg, field) {
			t.Errorf("error should list all missing fields; %q not found in %q", field, errMsg)
		}
	}

	// Verify error has the correct prefix
	if !strings.HasPrefix(errMsg, "missing required config: ") {
		t.Errorf("error should start with 'missing required config: ', got %q", errMsg)
	}
}
