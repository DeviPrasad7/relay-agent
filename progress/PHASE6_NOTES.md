# Phase 6 Notes – Temporal Correlation

## Decisions & Implementation Details

1. **Anomaly Window Grouping**:
   - Anomalies are sorted in chronological order.
   - If an anomaly occurs within `TimeWindowSec` (default 30 seconds) of the end of the current group, it is added to that group, extending its duration.
   - If it occurs after the window, a new group is started.
