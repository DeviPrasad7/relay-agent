// Package llm provides DeepSeek API integration for root cause analysis.
//
// Responsibilities:
//   - Define LLMAnalyzer interface
//   - Define TokenUsage struct
//
// Last updated: 2026-06-09
package llm

import "context"

// TokenUsage tracks API token consumption.
type TokenUsage struct {
	TotalTokens      int
	PromptTokens     int
	CompletionTokens int
}

// LLMAnalyzer is implemented by DeepSeek client.
type LLMAnalyzer interface {
	Analyze(ctx context.Context, anomalyContext string) (string, error)
	GetTokenUsage() TokenUsage
	GetRemainingTokens(limit int) int // returns limit - total used
}
