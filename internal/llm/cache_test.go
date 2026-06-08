package llm

import (
	"context"
	"testing"
	_ "github.com/mattn/go-sqlite3"
)

type mockCacheRepo struct {
	data map[string]string
}

func (m *mockCacheRepo) Get(ctx context.Context, hash string) (string, error) {
	if val, ok := m.data[hash]; ok {
		return val, nil
	}
	return "", nil
}
func (m *mockCacheRepo) Set(ctx context.Context, hash, response string, ttlSeconds int) error {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	m.data[hash] = response
	return nil
}
func (m *mockCacheRepo) CleanExpired(ctx context.Context) error { return nil }

type mockInner struct {
	called   bool
	response string
}

func (m *mockInner) Analyze(ctx context.Context, s string) (string, error) {
	m.called = true
	return m.response, nil
}
func (m *mockInner) GetTokenUsage() TokenUsage        { return TokenUsage{} }
func (m *mockInner) GetRemainingTokens(limit int) int { return limit }

func TestCachedAnalyzer(t *testing.T) {
	cache := &mockCacheRepo{}
	inner := &mockInner{response: "cached response"}
	cached := NewCachedAnalyzer(inner, cache, 24)
	ctx := context.Background()
	// First call: miss, calls inner
	resp1, _ := cached.Analyze(ctx, "test context")
	if resp1 != "cached response" {
		t.Errorf("expected cached response, got %s", resp1)
	}
	if !inner.called {
		t.Error("inner not called on first call")
	}
	// Reset inner called flag
	inner.called = false
	// Second call: should hit cache
	resp2, _ := cached.Analyze(ctx, "test context")
	if resp2 != "cached response" {
		t.Errorf("expected cached response, got %s", resp2)
	}
	if inner.called {
		t.Error("inner called again on cache hit")
	}
}
