// Package llm cache wrapper using storage.CacheRepository.
//
// Last updated: 2026-06-09
package llm

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"relay-agent/internal/storage"
)

// CachedAnalyzer wraps an LLMAnalyzer with caching.
type CachedAnalyzer struct {
	inner    LLMAnalyzer
	cache    storage.CacheRepository
	ttlHours int
}

// NewCachedAnalyzer creates a cached analyzer.
func NewCachedAnalyzer(inner LLMAnalyzer, cache storage.CacheRepository, ttlHours int) *CachedAnalyzer {
	return &CachedAnalyzer{
		inner:    inner,
		cache:    cache,
		ttlHours: ttlHours,
	}
}

// hashContext creates a SHA256 hash of the anomaly context for cache key.
func hashContext(context string) string {
	hash := sha256.Sum256([]byte(context))
	return hex.EncodeToString(hash[:])
}

// Analyze checks cache first, then calls inner if miss.
func (c *CachedAnalyzer) Analyze(ctx context.Context, anomalyContext string) (string, error) {
	hash := hashContext(anomalyContext)
	// Try cache
	cached, err := c.cache.Get(ctx, hash)
	if err == nil && cached != "" {
		return cached, nil
	}
	// Cache miss – call inner
	response, err := c.inner.Analyze(ctx, anomalyContext)
	if err != nil {
		return "", err
	}
	// Store in cache
	ttlSeconds := c.ttlHours * 3600
	if err := c.cache.Set(ctx, hash, response, ttlSeconds); err != nil {
		// Log but don't fail
		fmt.Printf("Cache set error: %v\n", err)
	}
	return response, nil
}

// GetTokenUsage delegates to inner.
func (c *CachedAnalyzer) GetTokenUsage() TokenUsage {
	return c.inner.GetTokenUsage()
}

// GetRemainingTokens delegates to inner.
func (c *CachedAnalyzer) GetRemainingTokens(limit int) int {
	return c.inner.GetRemainingTokens(limit)
}
