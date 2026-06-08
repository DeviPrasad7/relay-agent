package detection

import (
	"context"
	"testing"
)

func TestLatencyDegradationDetector_Configure(t *testing.T) {
	d := &LatencyDegradationDetector{}
	cfg := map[string]interface{}{
		"window_minutes":        10,
		"baseline_window_minutes": 60,
		"threshold_multiplier":   2.5,
	}
	err := d.Configure(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if d.windowMinutes != 10 {
		t.Errorf("expected window_minutes=10, got %d", d.windowMinutes)
	}
	if d.baselineWindowMinutes != 60 {
		t.Errorf("expected baseline_window_minutes=60, got %d", d.baselineWindowMinutes)
	}
	if d.thresholdMultiplier != 2.5 {
		t.Errorf("expected threshold_multiplier=2.5, got %f", d.thresholdMultiplier)
	}
}

func TestLatencyDegradationDetector_Detect(t *testing.T) {
	ctx := context.Background()
	d := &LatencyDegradationDetector{
		windowMinutes:       5,
		baselineWindowMinutes: 30,
		thresholdMultiplier: 2.0,
	}
	tests := []struct {
		name        string
		detCtx      *DetectionContext
		wantAnomaly bool
	}{
		{
			name: "degradation above threshold",
			detCtx: &DetectionContext{
				Service: "payment",
				WindowEnd: 1300,
				Metrics: map[string]interface{}{
					"p95_current":  250.0,
					"p95_baseline": 100.0,
				},
			},
			wantAnomaly: true,
		},
		{
			name: "no degradation below threshold",
			detCtx: &DetectionContext{
				Service: "payment",
				Metrics: map[string]interface{}{
					"p95_current":  150.0,
					"p95_baseline": 100.0,
				},
			},
			wantAnomaly: false,
		},
		{
			name: "missing p95_current",
			detCtx: &DetectionContext{
				Service: "payment",
				Metrics: map[string]interface{}{
					"p95_baseline": 100.0,
				},
			},
			wantAnomaly: false, // error
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anomalies, err := d.Detect(ctx, tt.detCtx)
			if tt.name == "missing p95_current" {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if tt.wantAnomaly && len(anomalies) == 0 {
				t.Error("expected anomaly, got none")
			}
			if !tt.wantAnomaly && len(anomalies) > 0 {
				t.Errorf("expected no anomaly, got %v", anomalies)
			}
		})
	}
}

func TestComputeP95FromDurations(t *testing.T) {
	tests := []struct {
		name      string
		durations []int
		expected  float64
	}{
		{"empty", []int{}, 0},
		{"single", []int{100}, 100},
		{"multiple", []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}, 100},
		{"odd count", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, 11},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeP95FromDurations(tt.durations)
			if got != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, got)
			}
		})
	}
}
