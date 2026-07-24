package bootstrap

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"
)

func TestLogJSON_EmitsValidJSON(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.LstdFlags)

	entry := LogEntry{
		Level:     "info",
		RequestID: "req-123",
		Method:    "GET",
		Path:      "/health",
		Status:    200,
		Duration:  "12ms",
	}

	LogJSON(entry)

	output := buf.String()
	// Remove trailing newline from log.Println
	output = output[:len(output)-1]

	var parsed LogEntry
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v\nraw output: %s", err, output)
	}

	if parsed.Timestamp == "" {
		t.Error("expected timestamp to be populated")
	}
	if _, err := time.Parse(time.RFC3339Nano, parsed.Timestamp); err != nil {
		t.Errorf("timestamp is not RFC3339Nano: %v", err)
	}
	if parsed.Level != "info" {
		t.Errorf("expected level=info, got %s", parsed.Level)
	}
	if parsed.RequestID != "req-123" {
		t.Errorf("expected request_id=req-123, got %s", parsed.RequestID)
	}
	if parsed.Method != "GET" {
		t.Errorf("expected method=GET, got %s", parsed.Method)
	}
	if parsed.Path != "/health" {
		t.Errorf("expected path=/health, got %s", parsed.Path)
	}
	if parsed.Status != 200 {
		t.Errorf("expected status=200, got %d", parsed.Status)
	}
	if parsed.Duration != "12ms" {
		t.Errorf("expected duration=12ms, got %s", parsed.Duration)
	}
}

func TestLogJSON_OmitsEmptyFields(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.LstdFlags)

	entry := LogEntry{
		Level: "error",
		Error: "something failed",
	}

	LogJSON(entry)

	output := buf.String()
	output = output[:len(output)-1]

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(output), &raw); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	// These fields should be omitted (omitempty)
	if _, ok := raw["request_id"]; ok {
		t.Error("expected request_id to be omitted when empty")
	}
	if _, ok := raw["method"]; ok {
		t.Error("expected method to be omitted when empty")
	}
	if _, ok := raw["path"]; ok {
		t.Error("expected path to be omitted when empty")
	}
	if _, ok := raw["duration"]; ok {
		t.Error("expected duration to be omitted when empty")
	}

	// These should be present
	if raw["level"] != "error" {
		t.Errorf("expected level=error, got %v", raw["level"])
	}
	if raw["error"] != "something failed" {
		t.Errorf("expected error='something failed', got %v", raw["error"])
	}
	if _, ok := raw["timestamp"]; !ok {
		t.Error("expected timestamp to always be present")
	}
}

func TestLogJSON_ErrorEntry(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.LstdFlags)

	entry := LogEntry{
		Level:     "error",
		RequestID: "req-456",
		Method:    "POST",
		Path:      "/ventas",
		Status:    500,
		Duration:  "250ms",
		Error:     "database connection timeout",
	}

	LogJSON(entry)

	output := buf.String()
	output = output[:len(output)-1]

	var parsed LogEntry
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if parsed.Error != "database connection timeout" {
		t.Errorf("expected error message, got %s", parsed.Error)
	}
	if parsed.Status != 500 {
		t.Errorf("expected status=500, got %d", parsed.Status)
	}
}

func TestLogJSON_TimestampIsUTC(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.LstdFlags)

	LogJSON(LogEntry{Level: "info"})

	output := buf.String()
	output = output[:len(output)-1]

	var parsed LogEntry
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	ts, err := time.Parse(time.RFC3339Nano, parsed.Timestamp)
	if err != nil {
		t.Fatalf("failed to parse timestamp: %v", err)
	}

	if ts.Location() != time.UTC {
		t.Errorf("expected UTC timestamp, got %v", ts.Location())
	}
}
