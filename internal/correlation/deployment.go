// Package correlation deployment correlator.
//
// Responsibilities:
//   - Link incident groups to recent deployment events
//   - Choose the closest deployment within time window
//
// Last updated: 2026-06-09
package correlation

import (
	"context"
	"relay-agent/internal/ingestion"
)

// DeploymentCorrelator links incidents to deployment events.
type DeploymentCorrelator struct {
	windowSeconds int
}

// NewDeploymentCorrelator creates a correlator with given window (seconds before incident).
func NewDeploymentCorrelator(windowSeconds int) *DeploymentCorrelator {
	return &DeploymentCorrelator{windowSeconds: windowSeconds}
}

// Correlate implements the Correlator interface.
// It groups anomalies temporally and then links deployments.
func (dc *DeploymentCorrelator) Correlate(ctx context.Context, input *CorrelationInput) ([]IncidentGroup, error) {
	// First, temporal grouping
	tc := NewTemporalCorrelator(input.TimeWindowSec)
	groups, err := tc.Correlate(ctx, input)
	if err != nil {
		return nil, err
	}
	// Then link deployments
	return dc.linkDeployments(groups, input.DeploymentEvents, input.DeploymentWindowSec), nil
}

// linkDeployments is the core logic: for each group, find the most recent deployment
// within windowSeconds before the group start time.
func (dc *DeploymentCorrelator) linkDeployments(groups []IncidentGroup, deployments []ingestion.EventEntry, windowSeconds int) []IncidentGroup {
	if windowSeconds <= 0 {
		windowSeconds = dc.windowSeconds
	}
	if windowSeconds <= 0 {
		windowSeconds = 120 // default 2 minutes
	}
	if len(deployments) == 0 {
		return groups
	}
	// Sort deployments by time descending (latest first)
	sorted := make([]ingestion.EventEntry, len(deployments))
	copy(sorted, deployments)
	// sort descending
	for i := 0; i < len(sorted)-1; i++ {
		for j := i+1; j < len(sorted); j++ {
			if sorted[i].Timestamp < sorted[j].Timestamp {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	for i := range groups {
		groupStart := groups[i].StartTime
		// Find deployment with timestamp within [groupStart - windowSeconds, groupStart)
		var bestDeployment ingestion.EventEntry
		var found bool
		for _, dep := range sorted {
			if dep.Type != ingestion.EventTypeDeploy && dep.Type != ingestion.EventTypeRollback {
				continue
			}
			// Check if deployment service matches any anomaly in the group
			serviceMatches := false
			for _, anomaly := range groups[i].Anomalies {
				if anomaly.Service == dep.Service {
					serviceMatches = true
					break
				}
			}
			if !serviceMatches {
				continue
			}

			if dep.Timestamp < groupStart && dep.Timestamp >= groupStart-int64(windowSeconds) {
				// Found deployment within window; choose the closest (largest timestamp)
				if !found || dep.Timestamp > bestDeployment.Timestamp {
					bestDeployment = dep
					found = true
				}
			}
		}
		if found {
			depCopy := bestDeployment
			groups[i].LinkedDeployment = &depCopy
		}
	}
	return groups
}
