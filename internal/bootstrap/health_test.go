package bootstrap

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler_LocalMode_NilPool(t *testing.T) {
	handler := healthHandler(nil, "local")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", resp["status"])
	}
	if resp["database"] != "ok" {
		t.Errorf("expected database=ok, got %v", resp["database"])
	}
	if resp["app_env"] != "local" {
		t.Errorf("expected app_env=local, got %v", resp["app_env"])
	}
	// Local mode should NOT include Lambda metadata
	if _, exists := resp["memory_mb"]; exists {
		t.Error("local mode should not include memory_mb")
	}
	if _, exists := resp["cold_start"]; exists {
		t.Error("local mode should not include cold_start")
	}
}

func TestHealthHandler_LambdaMode_NilPool_IncludesMetadata(t *testing.T) {
	t.Setenv("AWS_LAMBDA_FUNCTION_MEMORY_SIZE", "512")

	handler := healthHandler(nil, "lambda")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", resp["status"])
	}
	if resp["app_env"] != "lambda" {
		t.Errorf("expected app_env=lambda, got %v", resp["app_env"])
	}
	if resp["memory_mb"] != "512" {
		t.Errorf("expected memory_mb=512, got %v", resp["memory_mb"])
	}
	if resp["cold_start"] != true {
		t.Errorf("expected cold_start=true on first call, got %v", resp["cold_start"])
	}

	// Second call should report cold_start=false
	req2 := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	var resp2 map[string]interface{}
	if err := json.Unmarshal(rec2.Body.Bytes(), &resp2); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if resp2["cold_start"] != false {
		t.Errorf("expected cold_start=false on second call, got %v", resp2["cold_start"])
	}
}

func TestHealthHandler_ContentType(t *testing.T) {
	handler := healthHandler(nil, "local")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type=application/json, got %q", ct)
	}
}

func TestHealthHandler_ResponseAlwaysContainsRequiredFields(t *testing.T) {
	tests := []struct {
		name   string
		appEnv string
	}{
		{"local mode", "local"},
		{"lambda mode", "lambda"},
		{"empty env", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := healthHandler(nil, tt.appEnv)
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			var resp map[string]interface{}
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("invalid JSON response: %v", err)
			}

			requiredKeys := []string{"status", "database", "app_env"}
			for _, key := range requiredKeys {
				if _, exists := resp[key]; !exists {
					t.Errorf("response missing required key %q", key)
				}
			}
		})
	}
}
