// Package detection factory creates detectors from configuration.
//
// Responsibilities:
//   - Instantiate detectors based on type string
//   - Apply common configuration
//
// Last updated: 2026-06-09
package detection

import "fmt"

// NewDetector creates a detector by name and applies config.
func NewDetector(detectorType string, config map[string]interface{}) (Detector, error) {
	var d Detector
	switch detectorType {
	case "error_spike":
		d = &ErrorSpikeDetector{}
	case "latency_degrade":
		d = &LatencyDegradationDetector{}
	default:
		return nil, fmt.Errorf("unknown detector type: %s", detectorType)
	}
	if err := d.Configure(config); err != nil {
		return nil, err
	}
	return d, nil
}
