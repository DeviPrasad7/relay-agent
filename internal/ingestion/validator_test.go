package ingestion

import (
	"testing"
	"time"
)

func TestValidateLog(t *testing.T) {
	now := time.Now().Unix()
	tests := []struct {
		name    string
		log     *LogEntry
		wantErr bool
	}{
		{"valid", &LogEntry{Service: "auth", Timestamp: now, Severity: "error", Message: "fail"}, false},
		{"missing service", &LogEntry{Timestamp: now, Severity: "error", Message: "fail"}, true},
		{"future timestamp", &LogEntry{Service: "auth", Timestamp: now + 1000, Severity: "error", Message: "fail"}, true},
		{"invalid severity", &LogEntry{Service: "auth", Timestamp: now, Severity: "critical", Message: "fail"}, true},
		{"empty message", &LogEntry{Service: "auth", Timestamp: now, Severity: "info", Message: ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLog(tt.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTrace(t *testing.T) {
	now := time.Now().Unix()
	tests := []struct {
		name    string
		trace   *TraceEntry
		wantErr bool
	}{
		{"valid", &TraceEntry{Service: "pay", Timestamp: now, DurationMs: 100, TraceID: "abc", SpanName: "charge", Status: "ok"}, false},
		{"missing service", &TraceEntry{Timestamp: now, DurationMs: 100, TraceID: "abc", SpanName: "charge", Status: "ok"}, true},
		{"negative duration", &TraceEntry{Service: "pay", Timestamp: now, DurationMs: -1, TraceID: "abc", SpanName: "charge", Status: "ok"}, true},
		{"invalid status", &TraceEntry{Service: "pay", Timestamp: now, DurationMs: 100, TraceID: "abc", SpanName: "charge", Status: "unknown"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTrace(tt.trace)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTrace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEvent(t *testing.T) {
	now := time.Now().Unix()
	tests := []struct {
		name    string
		event   *EventEntry
		wantErr bool
	}{
		{"valid deploy", &EventEntry{Service: "auth", Timestamp: now, Type: "deploy", Version: "v1.0"}, false},
		{"valid heartbeat", &EventEntry{Service: "auth", Timestamp: now, Type: "heartbeat"}, false},
		{"deploy missing version", &EventEntry{Service: "auth", Timestamp: now, Type: "deploy"}, true},
		{"invalid type", &EventEntry{Service: "auth", Timestamp: now, Type: "unknown"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEvent(tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
