package storage

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"
	"relay-agent/internal/ingestion"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := OpenDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	// Execute schema from init_db.sql
	schema, err := os.ReadFile("../../scripts/init_db.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func TestLogRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewSQLiteLogRepository(db)

	ctx := context.Background()
	log := &ingestion.LogEntry{
		Service:   "test",
		Timestamp: time.Now().Unix(),
		Severity:  "info",
		Message:   "test message",
	}
	err := repo.Store(ctx, log)
	if err != nil {
		t.Fatal(err)
	}

	filter := LogFilter{Service: "test"}
	logs, err := repo.Query(ctx, filter)
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 1 {
		t.Errorf("expected 1 log, got %d", len(logs))
	}
}

func TestRetention(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	logRepo := NewSQLiteLogRepository(db)
	// Insert a log with old created_at (by using old timestamp)
	oldTimestamp := time.Now().AddDate(0, 0, -10).Unix()
	_, err := db.Exec(`INSERT INTO logs(service, timestamp, severity, message, created_at) VALUES(?,?,?,?,?)`,
		"test", oldTimestamp, "info", "old", oldTimestamp)
	if err != nil {
		t.Fatal(err)
	}
	cutoff := time.Now().AddDate(0, 0, -7).Unix()
	err = logRepo.DeleteOlderThan(context.Background(), cutoff)
	if err != nil {
		t.Fatal(err)
	}
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM logs`).Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 logs after retention, got %d", count)
	}
}
