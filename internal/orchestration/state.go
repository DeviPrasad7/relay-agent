// Package orchestration manages processing state and idempotency.
//
// Responsibilities:
//   - Store last processed timestamp in SQLite (key-value table)
//   - Provide atomic update
//
// Last updated: 2026-06-09
package orchestration

import (
	"context"
	"database/sql"
	"errors"
)

// StateStore holds the last processed timestamp.
type StateStore struct {
	db *sql.DB
}

// NewStateStore creates a store and ensures the key-value table exists.
func NewStateStore(db *sql.DB) (*StateStore, error) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS state (
			key TEXT PRIMARY KEY,
			value INTEGER NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}
	return &StateStore{db: db}, nil
}

// GetLastProcessed returns the last processed timestamp (Unix seconds). Returns 0 if never set.
func (s *StateStore) GetLastProcessed(ctx context.Context) (int64, error) {
	var value int64
	err := s.db.QueryRowContext(ctx, `SELECT value FROM state WHERE key = 'last_processed_timestamp'`).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return value, err
}

// SetLastProcessed updates the last processed timestamp.
func (s *StateStore) SetLastProcessed(ctx context.Context, ts int64) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO state(key, value) VALUES('last_processed_timestamp', ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, ts)
	return err
}
