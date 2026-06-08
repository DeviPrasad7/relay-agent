package correlation

import (
	"context"
	"testing"
	"relay-agent/internal/detection"
)

func TestTemporalCorrelator(t *testing.T) {
	ctx := context.Background()
	tc := NewTemporalCorrelator(30)

	tests := []struct {
		name        string
		anomalies   []detection.Anomaly
		windowSec   int
		expectedGroups int
	}{
		{
			name:        "empty",
			anomalies:   []detection.Anomaly{},
			expectedGroups: 0,
		},
		{
			name: "single anomaly",
			anomalies: []detection.Anomaly{
				{Time: 1000, Service: "auth", Method: "error_spike"},
			},
			expectedGroups: 1,
		},
		{
			name: "two anomalies within window",
			anomalies: []detection.Anomaly{
				{Time: 1000},
				{Time: 1010},
			},
			expectedGroups: 1,
		},
		{
			name: "two anomalies outside window",
			anomalies: []detection.Anomaly{
				{Time: 1000},
				{Time: 1040},
			},
			expectedGroups: 2,
		},
		{
			name: "three anomalies, two windows",
			anomalies: []detection.Anomaly{
				{Time: 1000},
				{Time: 1010},
				{Time: 1050},
			},
			expectedGroups: 2, // first two grouped, last separate
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &CorrelationInput{
				Anomalies:     tt.anomalies,
				TimeWindowSec: tt.windowSec,
			}
			groups, err := tc.Correlate(ctx, input)
			if err != nil {
				t.Fatal(err)
			}
			if len(groups) != tt.expectedGroups {
				t.Errorf("expected %d groups, got %d", tt.expectedGroups, len(groups))
			}
		})
	}
}
