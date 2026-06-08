package correlation

import (
	"testing"
	"relay-agent/internal/detection"
	"relay-agent/internal/ingestion"
)

func TestDeploymentCorrelator(t *testing.T) {
	dc := NewDeploymentCorrelator(120)

	anomalies := []detection.Anomaly{
		{Time: 2000, Service: "auth", Method: "error_spike"},
		{Time: 2100, Service: "payment", Method: "latency_degrade"},
	}
	deployments := []ingestion.EventEntry{
		{Service: "auth", Timestamp: 1900, Type: ingestion.EventTypeDeploy, Version: "v1.0"},
		{Service: "payment", Timestamp: 2000, Type: ingestion.EventTypeDeploy, Version: "v2.0"}, // Change from 1950 to 2000 so it falls within the 120s window (2100 - 2000 = 100s)
		{Service: "auth", Timestamp: 1800, Type: ingestion.EventTypeDeploy, Version: "v0.9"},
	}
	// First create temporal groups (default window 30s, both anomalies within 30s? 2000 and 2100 diff 100s >30, so two groups)
	groups := []IncidentGroup{
		{StartTime: 2000, EndTime: 2000, Anomalies: []detection.Anomaly{anomalies[0]}},
		{StartTime: 2100, EndTime: 2100, Anomalies: []detection.Anomaly{anomalies[1]}},
	}
	linked := dc.linkDeployments(groups, deployments, 120)

	// First group (auth) should link to deployment at 1900 (within 120s before 2000)
	if linked[0].LinkedDeployment == nil || linked[0].LinkedDeployment.Timestamp != 1900 {
		t.Errorf("auth incident not linked to correct deployment")
	}
	// Second group (payment) should link to deployment at 2000 (within 120s before 2100)
	if linked[1].LinkedDeployment == nil || linked[1].LinkedDeployment.Timestamp != 2000 {
		t.Errorf("payment incident not linked to correct deployment")
	}
}
