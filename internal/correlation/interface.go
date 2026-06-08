// Package correlation defines interfaces for grouping anomalies and linking deployments.
//
// Responsibilities:
//   - Define Correlator interface
//   - Define IncidentGroup structure
//
// Last updated: 2026-06-09
package correlation

import (
	"context"
	"relay-agent/internal/detection"
	"relay-agent/internal/ingestion"
)

// CorrelationInput holds data for correlation.
type CorrelationInput struct {
	Anomalies          []detection.Anomaly
	DeploymentEvents   []ingestion.EventEntry
	TimeWindowSec      int // anomalies within this window (seconds) are grouped
	DeploymentWindowSec int // deployment within this many seconds before anomaly is linked
}

// IncidentGroup represents a group of correlated anomalies, optionally linked to a deployment.
type IncidentGroup struct {
	StartTime        int64
	EndTime          int64
	Anomalies        []detection.Anomaly
	LinkedDeployment *ingestion.EventEntry
}

// Correlator is implemented by all temporal and deployment correlators.
type Correlator interface {
	Correlate(ctx context.Context, input *CorrelationInput) ([]IncidentGroup, error)
}
