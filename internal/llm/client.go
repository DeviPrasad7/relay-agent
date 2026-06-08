// Package llm DeepSeek API client.
//
// Last updated: 2026-06-09
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

// DeepSeekClient implements LLMAnalyzer.
type DeepSeekClient struct {
	apiKey      string
	baseURL     string
	model       string
	maxTokens   int
	httpClient  *http.Client
	tokenUsage  TokenUsage
}

// Config holds client configuration.
type Config struct {
	APIKey          string
	BaseURL         string // default https://api.deepseek.com/v1
	Model           string // default deepseek-chat
	MaxTokens       int    // default 1000
	TimeoutSeconds  int
}

// NewDeepSeekClient creates a new client.
func NewDeepSeekClient(cfg Config) *DeepSeekClient {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.deepseek.com/v1"
	}
	if cfg.Model == "" {
		cfg.Model = "deepseek-chat"
	}
	if cfg.MaxTokens == 0 {
		cfg.MaxTokens = 1000
	}
	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &DeepSeekClient{
		apiKey:     cfg.APIKey,
		baseURL:    cfg.BaseURL,
		model:      cfg.Model,
		maxTokens:  cfg.MaxTokens,
		httpClient: &http.Client{Timeout: timeout},
		tokenUsage: TokenUsage{},
	}
}

// chatRequest represents DeepSeek API request body.
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

// chatResponse represents DeepSeek API response.
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

// Analyze sends a prompt to DeepSeek and returns the response.
func (c *DeepSeekClient) Analyze(ctx context.Context, anomalyContext string) (string, error) {
	reqBody := chatRequest{
		Model: c.model,
		Messages: []message{
			{Role: "user", Content: anomalyContext},
		},
		MaxTokens:   c.maxTokens,
		Temperature: 0.3, // deterministic for debugging
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewReader(jsonBody))
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
		return "", fmt.Errorf("DeepSeek error: %s", chatResp.Error.Message)
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

// GetTokenUsage returns current token usage.
func (c *DeepSeekClient) GetTokenUsage() TokenUsage {
	return c.tokenUsage
}

// GetRemainingTokens returns remaining tokens before limit.
func (c *DeepSeekClient) GetRemainingTokens(limit int) int {
	remaining := limit - c.tokenUsage.TotalTokens
	if remaining < 0 {
		return 0
	}
	return remaining
}
