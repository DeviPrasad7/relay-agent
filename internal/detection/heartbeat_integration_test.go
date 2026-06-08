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

func setupEventDB(t *testing.T) (*sql.DB, storage.EventRepository) {
	db, err := storage.OpenDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	schema := `CREATE TABLE events (id INTEGER PRIMARY KEY, service TEXT, timestamp INTEGER, type TEXT, version TEXT, metadata TEXT, created_at INTEGER);`
	_, err = db.Exec(schema)
	if err != nil {
		t.Fatal(err)
	}
	eventRepo := storage.NewSQLiteEventRepository(db)
	return db, eventRepo
}

func TestHeartbeatDetectionIntegration(t *testing.T) {
	db, eventRepo := setupEventDB(t)
	defer db.Close()
	ctx := context.Background()
	
	now := time.Now().Unix()
	
	// Store heartbeats for auth and payment
	authHeartbeat := &ingestion.EventEntry{
		Service:   "auth-service",
		Timestamp: now - 30, // 30 seconds ago
		Type:      ingestion.EventTypeHeartbeat,
	}
	paymentHeartbeat := &ingestion.EventEntry{
		Service:   "payment-service",
		Timestamp: now - 120, // 120 seconds ago (expired)
		Type:      ingestion.EventTypeHeartbeat,
	}
	eventRepo.Store(ctx, authHeartbeat)
	eventRepo.Store(ctx, paymentHeartbeat)
	
	// Query latest heartbeat per service
	// For integration test, we'll manually load into detector
	detector := &HeartbeatFailureDetector{
		timeoutSeconds: 60,
		lastHeartbeats: make(map[string]int64),
		serviceTimeouts: make(map[string]int),
	}
	
	// Simulate loading from DB
	detector.UpdateHeartbeat("auth-service", authHeartbeat.Timestamp)
	detector.UpdateHeartbeat("payment-service", paymentHeartbeat.Timestamp)
	
	detCtx := &DetectionContext{
		Metrics: map[string]interface{}{
			"current_time": now,
		},
	}
	anomalies, err := detector.Detect(ctx, detCtx)
	if err != nil {
		t.Fatal(err)
	}
	
	// Expect only payment-service to have anomaly (120s > 60s timeout)
	if len(anomalies) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(anomalies))
	}
	if anomalies[0].Service != "payment-service" {
		t.Errorf("expected payment-service anomaly, got %s", anomalies[0].Service)
	}
	if anomalies[0].Method != "heartbeat_missing" {
		t.Errorf("expected method heartbeat_missing, got %s", anomalies[0].Method)
	}
}
