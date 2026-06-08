package llm

import "context"

type MockLLMClient struct {
	fixedResponse string
	tokenUsage    TokenUsage
}

func NewMockLLMClient(fixedResponse string) *MockLLMClient {
	if fixedResponse == "" {
		fixedResponse = "Simulated root cause: The issue appears to be a database connection timeout. Check connection pool settings."
	}
	return &MockLLMClient{fixedResponse: fixedResponse}
}

func (m *MockLLMClient) Analyze(ctx context.Context, anomalyContext string) (string, error) {
	// Simulate token usage
	m.tokenUsage.PromptTokens += 50
	m.tokenUsage.CompletionTokens += 30
	m.tokenUsage.TotalTokens += 80
	return m.fixedResponse, nil
}

func (m *MockLLMClient) GetTokenUsage() TokenUsage {
	return m.tokenUsage
}

func (m *MockLLMClient) GetRemainingTokens(limit int) int {
	return limit - m.tokenUsage.TotalTokens
}
