package bootstrap

import (
	"encoding/json"
	"log"
	"time"
)

// LogEntry represents a structured log entry for CloudWatch Logs.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	RequestID string `json:"request_id,omitempty"`
	Method    string `json:"method,omitempty"`
	Path      string `json:"path,omitempty"`
	Status    int    `json:"status,omitempty"`
	Duration  string `json:"duration,omitempty"`
	Error     string `json:"error,omitempty"`
}

// LogJSON emits a structured JSON log entry to stdout for CloudWatch consumption.
func LogJSON(entry LogEntry) {
	entry.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf(`{"timestamp":"%s","level":"error","error":"failed to marshal log entry: %s"}`,
			time.Now().UTC().Format(time.RFC3339Nano), err.Error())
		return
	}
	log.Println(string(data))
}
