// Package ingestion validator provides input validation functions.
//
// Responsibilities:
//   - Ensure required fields are present
//   - Check severity/event type enums
//   - Validate timestamp is not in future
//
// Last updated: 2026-06-09
package ingestion

import (
	"fmt"
	"time"
)

func ValidateLog(log *LogEntry) error {
	if log.Service == "" {
		return fmt.Errorf("service is required")
	}
	if log.Timestamp <= 0 || log.Timestamp > time.Now().Unix()+60 {
		return fmt.Errorf("invalid timestamp")
	}
	switch log.Severity {
	case SeverityInfo, SeverityWarn, SeverityError:
		// ok
	default:
		return fmt.Errorf("invalid severity: %s", log.Severity)
	}
	if log.Message == "" {
		return fmt.Errorf("message is required")
	}
	return nil
}

func ValidateTrace(trace *TraceEntry) error {
	if trace.Service == "" {
		return fmt.Errorf("service is required")
	}
	if trace.Timestamp <= 0 || trace.Timestamp > time.Now().Unix()+60 {
		return fmt.Errorf("invalid timestamp")
	}
	if trace.DurationMs < 0 {
		return fmt.Errorf("duration_ms cannot be negative")
	}
	if trace.TraceID == "" {
		return fmt.Errorf("trace_id is required")
	}
	if trace.SpanName == "" {
		return fmt.Errorf("span_name is required")
	}
	if trace.Status != "ok" && trace.Status != "error" {
		return fmt.Errorf("status must be 'ok' or 'error'")
	}
	return nil
}

func ValidateEvent(event *EventEntry) error {
	if event.Service == "" {
		return fmt.Errorf("service is required")
	}
	if event.Timestamp <= 0 || event.Timestamp > time.Now().Unix()+60 {
		return fmt.Errorf("invalid timestamp")
	}
	switch event.Type {
	case EventTypeDeploy, EventTypeRollback, EventTypeHeartbeat:
		// ok
	default:
		return fmt.Errorf("invalid event type: %s", event.Type)
	}
	if (event.Type == EventTypeDeploy || event.Type == EventTypeRollback) && event.Version == "" {
		return fmt.Errorf("version required for deploy/rollback events")
	}
	return nil
}
