package alerting

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"relay-agent/internal/storage"
)

func TestJSONReporter(t *testing.T) {
	tempDir := t.TempDir()
	reporter := NewJSONReporter(tempDir)
	incident := &storage.Incident{
		ID:             1,
		IncidentTime:   1700000000,
		DetectionMethod: "error_spike",
		AnomalyDetails:  `{"error_rate":0.15}`,
		RootCauseSummary: "DB pool exhausted",
		ResolutionStatus: "open",
	}
	err := reporter.Generate(context.Background(), incident)
	if err != nil {
		t.Fatal(err)
	}
	expectedFile := filepath.Join(tempDir, "incident_1_1700000000.json")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("expected file %s not created", expectedFile)
	}
}

func TestMarkdownReporter(t *testing.T) {
	tempDir := t.TempDir()
	reporter := NewMarkdownReporter(tempDir)
	deploymentID := 42
	incident := &storage.Incident{
		ID:                2,
		IncidentTime:      1700000000,
		DetectionMethod:   "latency_degrade",
		AnomalyDetails:    `{"p95_current":250,"p95_baseline":100}`,
		DeploymentEventID: &deploymentID,
		RootCauseSummary:  "Cache misconfiguration",
		RootCauseFull:     "The cache TTL was set too low causing repeated DB calls.",
		ResolutionStatus:  "investigating",
	}
	err := reporter.Generate(context.Background(), incident)
	if err != nil {
		t.Fatal(err)
	}
	expectedFile := filepath.Join(tempDir, "incident_2_1700000000.md")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("expected file %s not created", expectedFile)
	}
	// Read content and verify it contains expected strings
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatal(err)
	}
	if !containsString(string(content), "Cache misconfiguration") {
		t.Error("markdown missing root cause summary")
	}
}

func containsString(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > len(sub) && (s[:len(sub)] == sub || s[len(s)-len(sub):] == sub || (len(s) > len(sub) && helper(s, sub))))
}
func helper(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
