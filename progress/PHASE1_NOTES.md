# Phase 1 Notes – Ingestion API

## Decisions & Implementation Details

1. **In-Package Mock Interfaces**:
   - To decouple the `ingestion` handlers from the database layer (which is implemented in Phase 2), we defined local `LogStorer`, `TraceStorer`, and `EventStorer` interfaces in `api.go`.
   - This allows unit tests in `api_test.go` to easily mock storage interactions and verify server response logic.

2. **Validation Rules**:
   - Every incoming request's timestamp is validated. It must be non-zero and not more than 60 seconds into the future.
   - Severity enums for logs must be `info`, `warn`, or `error`.
   - Trace status must be `ok` or `error`.
   - Event types must be `deploy`, `rollback`, or `heartbeat`.
   - Version fields are mandatory for deployment and rollback events.

3. **Dependency Cleanup**:
   - Removed the unused `"time"` import from `models.go` to ensure compilation succeeds without Go unused import errors.
