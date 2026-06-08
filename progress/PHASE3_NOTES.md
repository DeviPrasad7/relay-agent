# Phase 3 Notes – Error Spike Detection

## Decisions & Implementation Details

1. **Common Detection Interfaces**:
   - Created `Detector` interface to represent the anomaly detection engine.
   - Introduced `DetectionContext` to encapsulate the current window and baseline metrics.
   - Introduced `Anomaly` representing the output format of a detected incident trigger.

2. **Error Spike Algorithm**:
   - Calculates if the current error rate exceeds the baseline `mean + multiplier * std_dev` threshold.
   - Designed to work dynamically depending on baseline metrics supplied at runtime.

3. **Factory & Configuration**:
   - Config loads raw YAML into a typed structure.
   - The `NewDetector` factory creates and configures the requested detector dynamically.
