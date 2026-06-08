package detection

import (
	"context"
	"testing"
)

func TestErrorSpikeDetector_Configure(t *testing.T) {
	d := &ErrorSpikeDetector{}
	cfg := map[string]interface{}{
		"window_minutes":        10,
		"baseline_window_minutes": 60,
		"std_dev_multiplier":    3.0,
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
	if d.stdDevMultiplier != 3.0 {
		t.Errorf("expected std_dev_multiplier=3.0, got %f", d.stdDevMultiplier)
	}
}

func TestErrorSpikeDetector_Detect(t *testing.T) {
	ctx := context.Background()
	d := &ErrorSpikeDetector{
		windowMinutes:      5,
		baselineWindowMinutes: 30,
		stdDevMultiplier:  2.0,
	}
	tests := []struct {
		name      string
		detCtx    *DetectionContext
		wantAnomaly bool
	}{
		{
			name: "spike above threshold",
			detCtx: &DetectionContext{
				Service: "auth",
				WindowStart: 1000,
				WindowEnd:   1300,
				Metrics: map[string]interface{}{"error_rate": 0.15},
				Baseline: map[string]float64{"mean": 0.02, "std_dev": 0.01},
			},
			wantAnomaly: true,
		},
		{
			name: "no spike below threshold",
			detCtx: &DetectionContext{
				Service: "auth",
				WindowStart: 1000,
				WindowEnd:   1300,
				Metrics: map[string]interface{}{"error_rate": 0.03},
				Baseline: map[string]float64{"mean": 0.02, "std_dev": 0.01},
			},
			wantAnomaly: false,
		},
		{
			name: "missing error_rate",
			detCtx: &DetectionContext{
				Service: "auth",
				Metrics: map[string]interface{}{},
				Baseline: map[string]float64{"mean": 0.02, "std_dev": 0.01},
			},
			wantAnomaly: false, // error returned, but we test for error in separate case
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anomalies, err := d.Detect(ctx, tt.detCtx)
			if tt.name == "missing error_rate" {
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

func TestNewDetector(t *testing.T) {
	cfg := map[string]interface{}{
		"window_minutes": 5,
	}
	d, err := NewDetector("error_spike", cfg)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := d.(*ErrorSpikeDetector); !ok {
		t.Errorf("expected *ErrorSpikeDetector, got %T", d)
	}
	_, err = NewDetector("unknown", cfg)
	if err == nil {
		t.Error("expected error for unknown detector")
	}
}
