// Package detection defines anomaly detection interfaces and types.
//
// Responsibilities:
//   - Define Detector interface for all anomaly detectors
//   - Define Anomaly struct for detection results
//
// Last updated: 2026-06-09
package detection

import "context"

// DetectionContext provides time window and metrics for detection.
type DetectionContext struct {
	Service     string
	WindowStart int64 // Unix seconds
	WindowEnd   int64 // Unix seconds
	Metrics     map[string]interface{}
	Baseline    map[string]float64 // e.g., "mean", "std_dev"
}

// Anomaly represents a detected anomaly.
type Anomaly struct {
	Service   string                 `json:"service"`
	Method    string                 `json:"method"` // error_spike, latency_degrade, heartbeat_missing
	Time      int64                  `json:"time"`
	Details   map[string]interface{} `json:"details"`
}

// Detector is implemented by all anomaly detectors.
type Detector interface {
	Detect(ctx context.Context, detCtx *DetectionContext) ([]Anomaly, error)
	Configure(cfg map[string]interface{}) error
}
