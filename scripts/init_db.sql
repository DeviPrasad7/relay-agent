-- Schema for Relay Agent
-- Last updated: 2026-06-09

CREATE TABLE IF NOT EXISTS logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    severity TEXT NOT NULL,
    message TEXT NOT NULL,
    created_at INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_logs_service_time ON logs(service, timestamp);
CREATE INDEX IF NOT EXISTS idx_logs_created ON logs(created_at);

CREATE TABLE IF NOT EXISTS traces (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    duration_ms INTEGER NOT NULL,
    trace_id TEXT NOT NULL,
    span_name TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_traces_service_time ON traces(service, timestamp);
CREATE INDEX IF NOT EXISTS idx_traces_created ON traces(created_at);

CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    type TEXT NOT NULL,
    version TEXT,
    metadata TEXT,
    created_at INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_events_service_time ON events(service, timestamp);
CREATE INDEX IF NOT EXISTS idx_events_created ON events(created_at);

CREATE TABLE IF NOT EXISTS incidents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    incident_time INTEGER NOT NULL,
    detection_method TEXT NOT NULL,
    anomaly_details TEXT NOT NULL,
    correlated_anomalies TEXT,
    deployment_event_id INTEGER,
    root_cause_summary TEXT,
    root_cause_full TEXT,
    resolution_status TEXT DEFAULT 'open',
    created_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS analysis_cache (
    context_hash TEXT PRIMARY KEY,
    llm_response TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    expires_at INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_cache_expires ON analysis_cache(expires_at);
