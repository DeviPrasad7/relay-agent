// Package llm prompt builder.
//
// Last updated: 2026-06-09
package llm

import (
	"encoding/json"
	"fmt"
	"relay-agent/internal/correlation"
)

// BuildAnomalyContext creates a prompt for DeepSeek from an incident group.
func BuildAnomalyContext(group *correlation.IncidentGroup) string {
	prompt := "You are a root cause analysis expert. Analyze the following incident and provide a concise root cause summary (max 200 words).\n\n"
	prompt += fmt.Sprintf("Incident started at: %d\n", group.StartTime)
	prompt += "Anomalies detected:\n"
	for _, a := range group.Anomalies {
		detailsJSON, _ := json.Marshal(a.Details)
		prompt += fmt.Sprintf("- Service: %s, Type: %s, Details: %s\n", a.Service, a.Method, string(detailsJSON))
	}
	if group.LinkedDeployment != nil {
		prompt += fmt.Sprintf("\nLinked deployment: service=%s, version=%s, time=%d\n",
			group.LinkedDeployment.Service, group.LinkedDeployment.Version, group.LinkedDeployment.Timestamp)
	}
	prompt += "\nWhat is the most likely root cause? Suggest specific debugging steps."
	return prompt
}
