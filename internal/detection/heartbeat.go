// Package detection heartbeat detector.
//
// Responsibilities:
//   - Track last heartbeat timestamp per service
//   - Detect when heartbeat is missing beyond timeout
//   - Support per-service timeout overrides
//
// Last updated: 2026-06-09
package detection

import (
	"context"
	"time"
)

// HeartbeatFailureDetector detects missing heartbeats.
type HeartbeatFailureDetector struct {
	timeoutSeconds      int
	serviceTimeouts     map[string]int // per-service override
	lastHeartbeats      map[string]int64 // service -> last heartbeat timestamp
}

// Configure sets detector parameters from config map.
func (d *HeartbeatFailureDetector) Configure(cfg map[string]interface{}) error {
	if val, ok := cfg["timeout_seconds"].(int); ok {
		d.timeoutSeconds = val
	}
	if val, ok := cfg["service_timeouts"].(map[string]interface{}); ok {
		d.serviceTimeouts = make(map[string]int)
		for svc, to := range val {
			if timeout, ok := to.(int); ok {
				d.serviceTimeouts[svc] = timeout
			}
		}
	}
	// Default timeout
	if d.timeoutSeconds == 0 {
		d.timeoutSeconds = 60
	}
	if d.lastHeartbeats == nil {
		d.lastHeartbeats = make(map[string]int64)
	}
	return nil
}

// UpdateHeartbeat records a heartbeat for a service (called from ingestion or orchestration).
func (d *HeartbeatFailureDetector) UpdateHeartbeat(service string, timestamp int64) {
	d.lastHeartbeats[service] = timestamp
}

// Detect checks for services missing heartbeat.
// Expects DetectionContext.Metrics to contain "current_time".
// Returns anomalies for each service that has missed heartbeat.
func (d *HeartbeatFailureDetector) Detect(ctx context.Context, detCtx *DetectionContext) ([]Anomaly, error) {
	currentTime, ok := detCtx.Metrics["current_time"].(int64)
	if !ok {
		currentTime = time.Now().Unix()
	}
	
	var anomalies []Anomaly
	for service, lastHeartbeat := range d.lastHeartbeats {
		timeout := d.timeoutSeconds
		if svcTimeout, exists := d.serviceTimeouts[service]; exists {
			timeout = svcTimeout
		}
		
		if lastHeartbeat == 0 {
			// No heartbeat ever received for this service, skip detection
			continue
		}
		
		timeSinceLast := currentTime - lastHeartbeat
		if timeSinceLast > int64(timeout) {
			anomalies = append(anomalies, Anomaly{
				Service: service,
				Method:  "heartbeat_missing",
				Time:    currentTime,
				Details: map[string]interface{}{
					"last_heartbeat":    lastHeartbeat,
					"timeout_seconds":   timeout,
					"seconds_since":     timeSinceLast,
					"current_time":      currentTime,
				},
			})
		}
	}
	return anomalies, nil
}

// GetRegisteredServices returns list of services that have sent heartbeats.
func (d *HeartbeatFailureDetector) GetRegisteredServices() []string {
	services := make([]string, 0, len(d.lastHeartbeats))
	for svc := range d.lastHeartbeats {
		services = append(services, svc)
	}
	return services
}
