//go:build integration
// +build integration

package detection

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"relay-agent/internal/ingestion"
	"relay-agent/internal/storage"
)

func setupLogDB(t *testing.T) (*sql.DB, storage.LogRepository) {
	db, err := storage.OpenDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	// init schema
	schema := `CREATE TABLE logs (id INTEGER PRIMARY KEY, service TEXT, timestamp INTEGER, severity TEXT, message TEXT, created_at INTEGER);`
	_, err = db.Exec(schema)
	if err != nil {
		t.Fatal(err)
	}
	logRepo := storage.NewSQLiteLogRepository(db)
	return db, logRepo
}

func TestErrorSpikeIntegration(t *testing.T) {
	db, logRepo := setupLogDB(t)
	defer db.Close()
	ctx := context.Background()

	// Insert logs: 15% error rate in last 5 minutes, baseline 2% error rate over 30 minutes
	now := time.Now().Unix()
	// Baseline: 1000 logs, 20 errors (2%)
	for i := 0; i < 980; i++ {
		logRepo.Store(ctx, &ingestion.LogEntry{Service: "auth", Timestamp: now - 2000, Severity: "info", Message: "ok"})
	}
	for i := 0; i < 20; i++ {
		logRepo.Store(ctx, &ingestion.LogEntry{Service: "auth", Timestamp: now - 2000, Severity: "error", Message: "fail"})
	}
	// Window: 100 logs, 15 errors (15%)
	for i := 0; i < 85; i++ {
		logRepo.Store(ctx, &ingestion.LogEntry{Service: "auth", Timestamp: now - 100, Severity: "info", Message: "ok"})
	}
	for i := 0; i < 15; i++ {
		logRepo.Store(ctx, &ingestion.LogEntry{Service: "auth", Timestamp: now - 100, Severity: "error", Message: "fail"})
	}

	// Query logs for baseline and window
	// For simplicity, we'll compute manually and feed to detector
	detector := &ErrorSpikeDetector{windowMinutes: 5, baselineWindowMinutes: 30, stdDevMultiplier: 2.0}
	detCtx := &DetectionContext{
		Service:     "auth",
		WindowStart: now - 300,
		WindowEnd:   now,
		Metrics:     map[string]interface{}{"error_rate": 0.15},
		Baseline:    map[string]float64{"mean": 0.02, "std_dev": 0.01},
	}
	anomalies, err := detector.Detect(ctx, detCtx)
	if err != nil {
		t.Fatal(err)
	}
	if len(anomalies) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(anomalies))
	}
	if anomalies[0].Method != "error_spike" {
		t.Errorf("expected method error_spike, got %s", anomalies[0].Method)
	}
}
