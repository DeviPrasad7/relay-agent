// Package ingestion defines data structures for logs, traces, and events.
//
// Responsibilities:
//   - Provide JSON‑tagged structs for request body unmarshaling
//   - Define severity and event type enums
//
// Dependencies: none
//
// Last updated: 2026-06-09
package ingestion

type Severity string

const (
	SeverityInfo  Severity = "info"
	SeverityWarn  Severity = "warn"
	SeverityError Severity = "error"
)

type EventType string

const (
	EventTypeDeploy    EventType = "deploy"
	EventTypeRollback  EventType = "rollback"
	EventTypeHeartbeat EventType = "heartbeat"
)

type LogEntry struct {
	Service   string   `json:"service"`
	Timestamp int64    `json:"timestamp"` // Unix seconds
	Severity  Severity `json:"severity"`
	Message   string   `json:"message"`
}

type TraceEntry struct {
	Service    string `json:"service"`
	Timestamp  int64  `json:"timestamp"`
	DurationMs int    `json:"duration_ms"`
	TraceID    string `json:"trace_id"`
	SpanName   string `json:"span_name"`
	Status     string `json:"status"` // "ok" or "error"
}

type EventEntry struct {
	Service   string    `json:"service"`
	Timestamp int64     `json:"timestamp"`
	Type      EventType `json:"type"`
	Version   string    `json:"version,omitempty"`
	Metadata  string    `json:"metadata,omitempty"` // JSON string
}
