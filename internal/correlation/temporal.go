// Package correlation temporal correlator.
//
// Responsibilities:
//   - Group anomalies by time proximity
//   - Merge overlapping or adjacent windows
//
// Last updated: 2026-06-09
package correlation

import (
	"context"
	"relay-agent/internal/detection"
	"sort"
)

// TemporalCorrelator groups anomalies by time window.
type TemporalCorrelator struct {
	windowSeconds int
}

// NewTemporalCorrelator creates a correlator with given window (seconds).
func NewTemporalCorrelator(windowSeconds int) *TemporalCorrelator {
	return &TemporalCorrelator{windowSeconds: windowSeconds}
}

// Correlate groups anomalies that occur within windowSeconds of each other.
func (tc *TemporalCorrelator) Correlate(ctx context.Context, input *CorrelationInput) ([]IncidentGroup, error) {
	if len(input.Anomalies) == 0 {
		return nil, nil
	}
	window := tc.windowSeconds
	if input.TimeWindowSec > 0 {
		window = input.TimeWindowSec
	}
	if window <= 0 {
		window = 30 // default 30 seconds
	}

	// Sort anomalies by time
	anomalies := make([]detection.Anomaly, len(input.Anomalies))
	copy(anomalies, input.Anomalies)
	sort.Slice(anomalies, func(i, j int) bool {
		return anomalies[i].Time < anomalies[j].Time
	})

	var groups []IncidentGroup
	currentGroup := IncidentGroup{
		StartTime: anomalies[0].Time,
		EndTime:   anomalies[0].Time,
		Anomalies: []detection.Anomaly{anomalies[0]},
	}

	for i := 1; i < len(anomalies); i++ {
		// If anomaly time is within window of current group's end time, extend group
		if anomalies[i].Time <= currentGroup.EndTime+int64(window) {
			currentGroup.EndTime = max(currentGroup.EndTime, anomalies[i].Time)
			currentGroup.Anomalies = append(currentGroup.Anomalies, anomalies[i])
		} else {
			// Finalize current group, start new group
			groups = append(groups, currentGroup)
			currentGroup = IncidentGroup{
				StartTime: anomalies[i].Time,
				EndTime:   anomalies[i].Time,
				Anomalies: []detection.Anomaly{anomalies[i]},
			}
		}
	}
	if len(currentGroup.Anomalies) > 0 {
		groups = append(groups, currentGroup)
	}
	return groups, nil
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
