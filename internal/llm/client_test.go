package llm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAICompatibleClient_Analyze(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"choices":[{"message":{"content":"Root cause: database connection pool exhausted."}}],"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}}`))
	}))
	defer server.Close()
	client := NewOpenAICompatibleClient(Config{
		APIKey:    "test-key",
		BaseURL:   server.URL,
		Model:     "gpt-3.5-turbo",
		MaxTokens: 1000,
	})
	ctx := context.Background()
	response, err := client.Analyze(ctx, "test anomaly context")
	if err != nil {
		t.Fatal(err)
	}
	if response != "Root cause: database connection pool exhausted." {
		t.Errorf("unexpected response: %s", response)
	}
	usage := client.GetTokenUsage()
	if usage.TotalTokens != 15 {
		t.Errorf("expected total tokens 15, got %d", usage.TotalTokens)
	}
}

func TestOpenAICompatibleClient_GetRemainingTokens(t *testing.T) {
	client := &OpenAICompatibleClient{tokenUsage: TokenUsage{TotalTokens: 100}}
	remaining := client.GetRemainingTokens(500)
	if remaining != 400 {
		t.Errorf("expected 400, got %d", remaining)
	}
}
