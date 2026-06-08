# Phase 5 Notes – Heartbeat Failure Detection

## Decisions & Implementation Details

1. **Service Registration & Override**:
   - The detector dynamically registers services when heartbeats are received.
   - Per-service custom timeout values can be configured (e.g., `auth-service` has 30s override). If none is configured, the default `timeout_seconds` is used.

2. **Detection Logic**:
   - Calculates the delta between the context's `current_time` and the service's `last_heartbeat`.
   - If the duration exceeds the timeout threshold, a `heartbeat_missing` anomaly is triggered.
