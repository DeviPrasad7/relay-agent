// Package detection error_spike detector.
//
// Responsibilities:
//   - Calculate error rate over a time window
//   - Compare against baseline (mean + N*std_dev)
//   - Return anomaly if threshold exceeded
//
// Last updated: 2026-06-09
package detection

import (
	"context"
	"fmt"
)

// ErrorSpikeDetector detects increases in error rate.
type ErrorSpikeDetector struct {
	windowMinutes      int
	baselineWindowMinutes int
	stdDevMultiplier   float64
}

// Configure sets detector parameters from config map.
func (d *ErrorSpikeDetector) Configure(cfg map[string]interface{}) error {
	if val, ok := cfg["window_minutes"].(int); ok {
		d.windowMinutes = val
	}
	if val, ok := cfg["baseline_window_minutes"].(int); ok {
		d.baselineWindowMinutes = val
	}
	if val, ok := cfg["std_dev_multiplier"].(float64); ok {
		d.stdDevMultiplier = val
	}
	// Set defaults if not provided
	if d.windowMinutes == 0 {
		d.windowMinutes = 5
	}
	if d.baselineWindowMinutes == 0 {
		d.baselineWindowMinutes = 30
	}
	if d.stdDevMultiplier == 0 {
		d.stdDevMultiplier = 2.0
	}
	return nil
}

// Detect checks for error rate spike.
// It expects DetectionContext.Metrics to contain "error_rate" and
// DetectionContext.Baseline to contain "mean" and "std_dev".
func (d *ErrorSpikeDetector) Detect(ctx context.Context, detCtx *DetectionContext) ([]Anomaly, error) {
	// Get current error rate
	errorRate, ok := detCtx.Metrics["error_rate"].(float64)
	if !ok {
		return nil, fmt.Errorf("error_rate not found or invalid in Metrics")
	}
	mean, ok := detCtx.Baseline["mean"]
	if !ok {
		return nil, fmt.Errorf("baseline mean not found")
	}
	stdDev, ok := detCtx.Baseline["std_dev"]
	if !ok {
		return nil, fmt.Errorf("baseline std_dev not found")
	}
	threshold := mean + d.stdDevMultiplier*stdDev
	if errorRate > threshold {
		return []Anomaly{{
			Service: detCtx.Service,
			Method:  "error_spike",
			Time:    detCtx.WindowEnd,
			Details: map[string]interface{}{
				"error_rate":   errorRate,
				"threshold":    threshold,
				"window_start": detCtx.WindowStart,
				"window_end":   detCtx.WindowEnd,
				"mean":         mean,
				"std_dev":      stdDev,
			},
		}}, nil
	}
	return nil, nil
}
