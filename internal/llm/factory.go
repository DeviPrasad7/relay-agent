package llm

import (
	"os"
	"strconv"
)

// NewLLMAnalyzerFromEnv creates a client based on environment variables.
// LLM_PROVIDER: openai, deepseek, custom, mock (default mock)
// For deepseek: sets baseURL to https://api.deepseek.com/v1
// For openai:  baseURL https://api.openai.com/v1
// For custom:   you must set LLM_BASE_URL and LLM_API_KEY
// If LLM_MOCK=true, always returns mock client (overrides provider)
func NewLLMAnalyzerFromEnv() LLMAnalyzer {
	// Check for mock override
	if mock := os.Getenv("LLM_MOCK"); mock == "true" || mock == "1" {
		response := os.Getenv("LLM_MOCK_RESPONSE")
		return NewMockLLMClient(response)
	}
	provider := os.Getenv("LLM_PROVIDER")
	if provider == "" {
		provider = "mock"
	}
	var cfg Config
	switch provider {
	case "deepseek":
		cfg.BaseURL = "https://api.deepseek.com/v1"
		cfg.APIKey = os.Getenv("DEEPSEEK_API_KEY")
		if cfg.APIKey == "" {
			cfg.APIKey = os.Getenv("LLM_API_KEY")
		}
		cfg.Model = os.Getenv("LLM_MODEL")
		if cfg.Model == "" {
			cfg.Model = "deepseek-chat"
		}
	case "openai":
		cfg.BaseURL = "https://api.openai.com/v1"
		cfg.APIKey = os.Getenv("OPENAI_API_KEY")
		if cfg.APIKey == "" {
			cfg.APIKey = os.Getenv("LLM_API_KEY")
		}
		cfg.Model = os.Getenv("LLM_MODEL")
		if cfg.Model == "" {
			cfg.Model = "gpt-3.5-turbo"
		}
	case "custom":
		cfg.BaseURL = os.Getenv("LLM_BASE_URL")
		cfg.APIKey = os.Getenv("LLM_API_KEY")
		cfg.Model = os.Getenv("LLM_MODEL")
		if cfg.Model == "" {
			cfg.Model = "default"
		}
	default: // mock
		response := os.Getenv("LLM_MOCK_RESPONSE")
		return NewMockLLMClient(response)
	}
	if maxTokens := os.Getenv("LLM_MAX_TOKENS"); maxTokens != "" {
		if v, err := strconv.Atoi(maxTokens); err == nil {
			cfg.MaxTokens = v
		}
	}
	if timeout := os.Getenv("LLM_TIMEOUT_SECONDS"); timeout != "" {
		if v, err := strconv.Atoi(timeout); err == nil {
			cfg.TimeoutSeconds = v
		}
	}
	return NewOpenAICompatibleClient(cfg)
}
