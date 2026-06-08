# Phase 7 Notes – Deployment Correlation

## Decisions & Implementation Details

1. **Deployment Linking**:
   - For each temporal incident group, the correlator scans recent deployments (`deploy` or `rollback` event types).
   - If a deployment event occurred within the configurable deployment window (default 120 seconds) before the start of the incident, it is linked.
   - If multiple deployments occur within the window, the closest one (the latest deployment before the incident) is selected.
