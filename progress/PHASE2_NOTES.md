# Phase 2 Notes – SQLite Storage Layer

## Decisions & Implementation Details

1. **SQLite Connection Optimization**:
   - Enabled WAL (Write-Ahead Logging) and NORMAL sync modes to optimize write performance while preserving safety.
   - Set max open and idle connections to `1` to avoid database locking issues typical with SQLite under concurrent writes.

2. **Repository Abstractions**:
   - Decoupled SQL queries from business logic using the `LogRepository`, `TraceRepository`, `EventRepository`, `IncidentRepository`, and `CacheRepository` interfaces.

3. **Retention Background Job**:
   - Built a background `RetentionManager` that deletes data older than `retentionDays` at configurable intervals (e.g., hourly).

4. **Integration Testing**:
   - Created in-memory SQLite schema initialization from the `scripts/init_db.sql` file to run tests without polluting local disks.
