// Package detection latency_degrade detector.
//
// Responsibilities:
//   - Compute p95 latency over current window
//   - Compare against baseline p95 from previous window
//   - Trigger anomaly if ratio exceeds threshold
//
// Last updated: 2026-06-09
package detection

import (
	"context"
	"fmt"
	"sort"
)

// LatencyDegradationDetector detects increases in p95 latency.
type LatencyDegradationDetector struct {
	windowMinutes       int
	baselineWindowMinutes int
	thresholdMultiplier float64
}

// Configure sets detector parameters from config map.
func (d *LatencyDegradationDetector) Configure(cfg map[string]interface{}) error {
	if val, ok := cfg["window_minutes"].(int); ok {
		d.windowMinutes = val
	}
	if val, ok := cfg["baseline_window_minutes"].(int); ok {
		d.baselineWindowMinutes = val
	}
	if val, ok := cfg["threshold_multiplier"].(float64); ok {
		d.thresholdMultiplier = val
	}
	// Defaults
	if d.windowMinutes == 0 {
		d.windowMinutes = 5
	}
	if d.baselineWindowMinutes == 0 {
		d.baselineWindowMinutes = 30
	}
	if d.thresholdMultiplier == 0 {
		d.thresholdMultiplier = 2.0
	}
	return nil
}

// Detect checks for latency degradation.
// Expects DetectionContext.Metrics to contain "p95_current" and "p95_baseline".
func (d *LatencyDegradationDetector) Detect(ctx context.Context, detCtx *DetectionContext) ([]Anomaly, error) {
	p95Current, ok := detCtx.Metrics["p95_current"].(float64)
	if !ok {
		return nil, fmt.Errorf("p95_current not found or invalid in Metrics")
	}
	p95Baseline, ok := detCtx.Metrics["p95_baseline"].(float64)
	if !ok {
		return nil, fmt.Errorf("p95_baseline not found or invalid in Metrics")
	}
	if p95Baseline == 0 {
		return nil, nil // avoid division by zero
	}
	ratio := p95Current / p95Baseline
	if ratio > d.thresholdMultiplier {
		return []Anomaly{{
			Service: detCtx.Service,
			Method:  "latency_degrade",
			Time:    detCtx.WindowEnd,
			Details: map[string]interface{}{
				"p95_current":   p95Current,
				"p95_baseline":  p95Baseline,
				"ratio":         ratio,
				"threshold":     d.thresholdMultiplier,
				"window_start":  detCtx.WindowStart,
				"window_end":    detCtx.WindowEnd,
			},
		}}, nil
	}
	return nil, nil
}

// ComputeP95FromDurations is a helper function to compute p95 latency from a slice of durations (ms).
func ComputeP95FromDurations(durations []int) float64 {
	if len(durations) == 0 {
		return 0
	}
	sort.Ints(durations)
	idx := int(float64(len(durations)) * 0.95)
	if idx >= len(durations) {
		idx = len(durations) - 1
	}
	return float64(durations[idx])
}
