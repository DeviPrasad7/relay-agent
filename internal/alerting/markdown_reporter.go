// Package alerting Markdown reporter.
//
// Responsibilities:
//   - Write human-readable incident summary as Markdown
//   - Include anomaly details, deployment link, and root cause
//
// Last updated: 2026-06-09
package alerting

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"relay-agent/internal/storage"
)

// MarkdownReporter writes Markdown summaries.
type MarkdownReporter struct {
	outputDir string
}

// NewMarkdownReporter creates a reporter that writes to outputDir.
func NewMarkdownReporter(outputDir string) *MarkdownReporter {
	return &MarkdownReporter{outputDir: outputDir}
}

// Generate writes the incident as a Markdown file.
func (r *MarkdownReporter) Generate(ctx context.Context, incident *storage.Incident) error {
	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return err
	}
	filename := filepath.Join(r.outputDir, fmt.Sprintf("incident_%d_%d.md", incident.ID, incident.IncidentTime))
	content := r.buildMarkdown(incident)
	return os.WriteFile(filename, []byte(content), 0644)
}

func (r *MarkdownReporter) buildMarkdown(incident *storage.Incident) string {
	tm := time.Unix(incident.IncidentTime, 0).Format(time.RFC3339)
	md := fmt.Sprintf("# Incident Report - %s\n\n", tm)
	md += fmt.Sprintf("**Incident ID**: %d\n", incident.ID)
	md += fmt.Sprintf("**Detection Method**: %s\n", incident.DetectionMethod)
	md += fmt.Sprintf("**Resolution Status**: %s\n\n", incident.ResolutionStatus)
	md += "## Anomaly Details\n```json\n" + incident.AnomalyDetails + "\n```\n\n"
	if incident.DeploymentEventID != nil {
		md += "## Linked Deployment\n"
		md += fmt.Sprintf("- Deployment Event ID: %d\n", *incident.DeploymentEventID)
	}
	if incident.RootCauseSummary != "" {
		md += "## Root Cause Summary\n" + incident.RootCauseSummary + "\n\n"
	}
	if incident.RootCauseFull != "" {
		md += "## Full Root Cause Analysis\n" + incident.RootCauseFull + "\n"
	}
	return md
}
