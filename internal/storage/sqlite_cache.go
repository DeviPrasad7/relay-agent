package storage

import (
	"context"
	"database/sql"
	"time"
)

type SQLiteCacheRepository struct {
	db *sql.DB
}

func NewSQLiteCacheRepository(db *sql.DB) *SQLiteCacheRepository {
	return &SQLiteCacheRepository{db: db}
}

func (r *SQLiteCacheRepository) Get(ctx context.Context, hash string) (string, error) {
	var response string
	var expiresAt int64
	err := r.db.QueryRowContext(ctx, `SELECT llm_response, expires_at FROM analysis_cache WHERE context_hash = ?`, hash).Scan(&response, &expiresAt)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	// If expired, treat as not found (cleanup will remove later)
	if expiresAt < currentTime() {
		return "", nil
	}
	return response, nil
}

func (r *SQLiteCacheRepository) Set(ctx context.Context, hash, response string, ttlSeconds int) error {
	expiresAt := currentTime() + int64(ttlSeconds)
	_, err := r.db.ExecContext(ctx,
		`INSERT OR REPLACE INTO analysis_cache(context_hash, llm_response, created_at, expires_at) VALUES(?, ?, ?, ?)`,
		hash, response, currentTime(), expiresAt,
	)
	return err
}

func (r *SQLiteCacheRepository) CleanExpired(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM analysis_cache WHERE expires_at < ?`, currentTime())
	return err
}

// Helper
func currentTime() int64 {
	return time.Now().Unix()
}
