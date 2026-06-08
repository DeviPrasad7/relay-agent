// Package alerting JSON reporter.
//
// Responsibilities:
//   - Write incident as JSON file to configured directory
//   - Filename: incident_{id}_{timestamp}.json
//
// Last updated: 2026-06-09
package alerting

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"relay-agent/internal/storage"
)

// JSONReporter writes incidents as JSON files.
type JSONReporter struct {
	outputDir string
}

// NewJSONReporter creates a reporter that writes to outputDir.
func NewJSONReporter(outputDir string) *JSONReporter {
	return &JSONReporter{outputDir: outputDir}
}

// Generate writes the incident to a JSON file.
func (r *JSONReporter) Generate(ctx context.Context, incident *storage.Incident) error {
	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return err
	}
	filename := filepath.Join(r.outputDir, fmt.Sprintf("incident_%d_%d.json", incident.ID, incident.IncidentTime))
	data, err := json.MarshalIndent(incident, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
