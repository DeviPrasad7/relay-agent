# Phase 4 Notes – Latency Degradation Detection

## Decisions & Implementation Details

1. **Percentile Calculation**:
   - Implemented a standard p95 latency helper function `ComputeP95FromDurations` which sorts latency durations and extracts the value at the 95th percentile index.

2. **Degradation Detection logic**:
   - Checks the ratio of p95 latency between the current window and baseline window.
   - If `p95Current / p95Baseline > thresholdMultiplier`, it flags a latency degradation anomaly.

3. **Integration & Factory**:
   - Registered `latency_degrade` with `NewDetector` factory.
   - Added configuration mapping in config struct.
