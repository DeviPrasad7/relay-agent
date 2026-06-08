// Package storage defines repository interfaces for data persistence.
//
// Responsibilities:
//   - Abstract database operations
//   - Enable dependency inversion
//
// Last updated: 2026-06-09
package storage

import (
	"context"
	"relay-agent/internal/ingestion"
)

// LogRepository handles log entries.
type LogRepository interface {
	Store(ctx context.Context, log *ingestion.LogEntry) error
	Query(ctx context.Context, filter LogFilter) ([]*ingestion.LogEntry, error)
	DeleteOlderThan(ctx context.Context, beforeUnix int64) error
}

type LogFilter struct {
	Service   string
	StartTime int64 // inclusive
	EndTime   int64 // exclusive
	Severity  string // optional
}

// TraceRepository handles trace entries.
type TraceRepository interface {
	Store(ctx context.Context, trace *ingestion.TraceEntry) error
	Query(ctx context.Context, filter TraceFilter) ([]*ingestion.TraceEntry, error)
	DeleteOlderThan(ctx context.Context, beforeUnix int64) error
}

type TraceFilter struct {
	Service   string
	StartTime int64
	EndTime   int64
	Status    string // "ok" or "error"
}

// EventRepository handles deployment/heartbeat events.
type EventRepository interface {
	Store(ctx context.Context, event *ingestion.EventEntry) error
	Query(ctx context.Context, filter EventFilter) ([]*ingestion.EventEntry, error)
	DeleteOlderThan(ctx context.Context, beforeUnix int64) error
}

type EventFilter struct {
	Service   string
	StartTime int64
	EndTime   int64
	Type      ingestion.EventType
}

// IncidentRepository stores detected incidents.
type IncidentRepository interface {
	Store(ctx context.Context, incident *Incident) error
	List(ctx context.Context, limit int) ([]*Incident, error)
	GetByID(ctx context.Context, id int) (*Incident, error)
}

// CacheRepository stores LLM response cache.
type CacheRepository interface {
	Get(ctx context.Context, hash string) (string, error) // returns empty string if not found
	Set(ctx context.Context, hash, response string, ttlSeconds int) error
	CleanExpired(ctx context.Context) error
}

// Incident is the storage model (matches schema in plan.md).
type Incident struct {
	ID                   int    `json:"id"`
	IncidentTime         int64  `json:"incident_time"`
	DetectionMethod      string `json:"detection_method"`
	AnomalyDetails       string `json:"anomaly_details"` // JSON
	CorrelatedAnomalies  string `json:"correlated_anomalies"` // JSON list
	DeploymentEventID    *int   `json:"deployment_event_id,omitempty"`
	RootCauseSummary     string `json:"root_cause_summary"`
	RootCauseFull        string `json:"root_cause_full"`
	ResolutionStatus     string `json:"resolution_status"`
	CreatedAt            int64  `json:"created_at"`
}
