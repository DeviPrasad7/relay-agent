// Package alerting defines report generation interfaces.
//
// Responsibilities:
//   - Define Reporter interface
//   - Provide shared incident report structure
//
// Last updated: 2026-06-09
package alerting

import (
	"context"
	"relay-agent/internal/storage"
)

// Reporter is implemented by JSON and Markdown reporters.
type Reporter interface {
	Generate(ctx context.Context, incident *storage.Incident) error
}

// IncidentReport is the enriched version used for output (same as storage.Incident but with derived fields).
// For simplicity, we reuse storage.Incident and add methods in reporters.
