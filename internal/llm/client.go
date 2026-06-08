package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAICompatibleClient works with any OpenAI‑compatible API (DeepSeek, OpenAI, local).
type OpenAICompatibleClient struct {
	baseURL    string
	apiKey     string
	model      string
	maxTokens  int
	httpClient *http.Client
	tokenUsage TokenUsage
}

type Config struct {
	BaseURL        string
	APIKey         string
	Model          string
	MaxTokens      int
	TimeoutSeconds int
}

func NewOpenAICompatibleClient(cfg Config) *OpenAICompatibleClient {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.openai.com/v1" // default
	}
	if cfg.Model == "" {
		cfg.Model = "gpt-3.5-turbo"
	}
	if cfg.MaxTokens == 0 {
		cfg.MaxTokens = 1000
	}
	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &OpenAICompatibleClient{
		baseURL:    cfg.BaseURL,
		apiKey:     cfg.APIKey,
		model:      cfg.Model,
		maxTokens:  cfg.MaxTokens,
		httpClient: &http.Client{Timeout: timeout},
	}
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *OpenAICompatibleClient) Analyze(ctx context.Context, anomalyContext string) (string, error) {
	reqBody := chatRequest{
		Model: c.model,
		Messages: []message{
			{Role: "user", Content: anomalyContext},
		},
		MaxTokens:   c.maxTokens,
		Temperature: 0.3,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	url := c.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}
	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", err
	}
	if chatResp.Error.Message != "" {
		return "", fmt.Errorf("LLM error: %s", chatResp.Error.Message)
	}
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	// Update token usage
	c.tokenUsage.PromptTokens += chatResp.Usage.PromptTokens
	c.tokenUsage.CompletionTokens += chatResp.Usage.CompletionTokens
	c.tokenUsage.TotalTokens += chatResp.Usage.TotalTokens
	return chatResp.Choices[0].Message.Content, nil
}

func (c *OpenAICompatibleClient) GetTokenUsage() TokenUsage {
	return c.tokenUsage
}

func (c *OpenAICompatibleClient) GetRemainingTokens(limit int) int {
	remaining := limit - c.tokenUsage.TotalTokens
	if remaining < 0 {
		return 0
	}
	return remaining
}
