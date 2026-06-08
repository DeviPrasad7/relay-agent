package detection

import (
	"context"
	"testing"
)

func TestHeartbeatFailureDetector_Configure(t *testing.T) {
	d := &HeartbeatFailureDetector{}
	cfg := map[string]interface{}{
		"timeout_seconds": 120,
		"service_timeouts": map[string]interface{}{
			"auth-service":    30,
			"payment-service": 90,
		},
	}
	err := d.Configure(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if d.timeoutSeconds != 120 {
		t.Errorf("expected timeout_seconds=120, got %d", d.timeoutSeconds)
	}
	if d.serviceTimeouts["auth-service"] != 30 {
		t.Errorf("expected auth-service timeout=30, got %d", d.serviceTimeouts["auth-service"])
	}
}

func TestHeartbeatFailureDetector_UpdateHeartbeat(t *testing.T) {
	d := &HeartbeatFailureDetector{lastHeartbeats: make(map[string]int64)}
	d.UpdateHeartbeat("auth", 1000)
	if d.lastHeartbeats["auth"] != 1000 {
		t.Errorf("expected heartbeat 1000, got %d", d.lastHeartbeats["auth"])
	}
}

func TestHeartbeatFailureDetector_Detect(t *testing.T) {
	ctx := context.Background()
	d := &HeartbeatFailureDetector{
		timeoutSeconds: 60,
		lastHeartbeats: make(map[string]int64),
		serviceTimeouts: make(map[string]int),
	}
	
	// Add heartbeats
	d.UpdateHeartbeat("auth", 1000)
	d.UpdateHeartbeat("payment", 1050)
	
	tests := []struct {
		name        string
		currentTime int64
		wantAnomalies int
		wantService string
	}{
		{
			name:        "auth heartbeat missing (101 seconds)",
			currentTime: 1101,
			wantAnomalies: 1,
			wantService: "auth",
		},
		{
			name:        "both heartbeats present (50 seconds)",
			currentTime: 1050,
			wantAnomalies: 0,
		},
		{
			name:        "custom timeout for payment (30 seconds)",
			currentTime: 1081, // 31 seconds after payment's last heartbeat
			wantAnomalies: 1,
			wantService: "payment",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset for test that needs custom timeout
			if tt.name == "custom timeout for payment (30 seconds)" {
				d.serviceTimeouts["payment"] = 30
				d.UpdateHeartbeat("auth", tt.currentTime - 10) // Keep auth fresh
			}
			detCtx := &DetectionContext{
				Metrics: map[string]interface{}{
					"current_time": tt.currentTime,
				},
			}
			anomalies, err := d.Detect(ctx, detCtx)
			if err != nil {
				t.Fatal(err)
			}
			if len(anomalies) != tt.wantAnomalies {
				t.Errorf("expected %d anomalies, got %d", tt.wantAnomalies, len(anomalies))
			}
			if tt.wantAnomalies > 0 && anomalies[0].Service != tt.wantService {
				t.Errorf("expected service %s, got %s", tt.wantService, anomalies[0].Service)
			}
		})
	}
}

func TestHeartbeatFailureDetector_NoHeartbeatRegistered(t *testing.T) {
	d := &HeartbeatFailureDetector{
		timeoutSeconds: 60,
		lastHeartbeats: make(map[string]int64),
	}
	ctx := context.Background()
	detCtx := &DetectionContext{
		Metrics: map[string]interface{}{
			"current_time": int64(2000),
		},
	}
	anomalies, err := d.Detect(ctx, detCtx)
	if err != nil {
		t.Fatal(err)
	}
	if len(anomalies) != 0 {
		t.Errorf("expected 0 anomalies when no heartbeats registered, got %d", len(anomalies))
	}
}
