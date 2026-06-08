package llm

import (
	"testing"
	"relay-agent/internal/detection"
	"relay-agent/internal/correlation"
	"relay-agent/internal/ingestion"
)

func TestBuildAnomalyContext(t *testing.T) {
	group := &correlation.IncidentGroup{
		StartTime: 1700000000,
		Anomalies: []detection.Anomaly{
			{Service: "auth", Method: "error_spike", Details: map[string]interface{}{"error_rate": 0.15}},
		},
		LinkedDeployment: &ingestion.EventEntry{
			Service: "auth", Version: "v1.2.3", Timestamp: 1700000000 - 60,
		},
	}
	prompt := BuildAnomalyContext(group)
	if prompt == "" {
		t.Error("prompt should not be empty")
	}
	if !contains(prompt, "error_spike") {
		t.Error("prompt missing anomaly details")
	}
	if !contains(prompt, "v1.2.3") {
		t.Error("prompt missing deployment info")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || (len(s) > len(substr) && containsHelper(s, substr))))
}
func containsHelper(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
