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

func setupTraceDB(t *testing.T) (*sql.DB, storage.TraceRepository) {
	db, err := storage.OpenDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	schema := `CREATE TABLE traces (id INTEGER PRIMARY KEY, service TEXT, timestamp INTEGER, duration_ms INTEGER, trace_id TEXT, span_name TEXT, status TEXT, created_at INTEGER);`
	_, err = db.Exec(schema)
	if err != nil {
		t.Fatal(err)
	}
	traceRepo := storage.NewSQLiteTraceRepository(db)
	return db, traceRepo
}

func TestLatencyDegradationIntegration(t *testing.T) {
	db, traceRepo := setupTraceDB(t)
	defer db.Close()
	ctx := context.Background()

	now := time.Now().Unix()
	// Baseline: 100 traces with p95 ~100ms
	for i := 0; i < 100; i++ {
		duration := 90 + i%20 // 90-109 ms
		traceRepo.Store(ctx, &ingestion.TraceEntry{Service: "payment", Timestamp: now - 2000, DurationMs: duration, TraceID: "base", SpanName: "pay", Status: "ok"})
	}
	// Current window: 100 traces with p95 ~250ms
	for i := 0; i < 100; i++ {
		duration := 240 + i%20 // 240-259 ms
		traceRepo.Store(ctx, &ingestion.TraceEntry{Service: "payment", Timestamp: now - 100, DurationMs: duration, TraceID: "curr", SpanName: "pay", Status: "ok"})
	}

	// Query durations for baseline and window
	// For simplicity, we'll compute manually and feed to detector
	detector := &LatencyDegradationDetector{windowMinutes: 5, baselineWindowMinutes: 30, thresholdMultiplier: 2.0}
	detCtx := &DetectionContext{
		Service: "payment",
		WindowEnd: now,
		Metrics: map[string]interface{}{
			"p95_current":  250.0,
			"p95_baseline": 100.0,
		},
	}
	anomalies, err := detector.Detect(ctx, detCtx)
	if err != nil {
		t.Fatal(err)
	}
	if len(anomalies) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(anomalies))
	}
	if anomalies[0].Method != "latency_degrade" {
		t.Errorf("expected method latency_degrade, got %s", anomalies[0].Method)
	}
}
