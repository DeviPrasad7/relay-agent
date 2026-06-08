package storage

import (
	"context"
	"database/sql"
	"relay-agent/internal/ingestion"
)

type SQLiteTraceRepository struct {
	db *sql.DB
}

func NewSQLiteTraceRepository(db *sql.DB) *SQLiteTraceRepository {
	return &SQLiteTraceRepository{db: db}
}

func (r *SQLiteTraceRepository) Store(ctx context.Context, trace *ingestion.TraceEntry) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO traces(service, timestamp, duration_ms, trace_id, span_name, status, created_at) VALUES(?,?,?,?,?,?,?)`,
		trace.Service, trace.Timestamp, trace.DurationMs, trace.TraceID, trace.SpanName, trace.Status, trace.Timestamp,
	)
	return err
}

func (r *SQLiteTraceRepository) Query(ctx context.Context, filter TraceFilter) ([]*ingestion.TraceEntry, error) {
	query := `SELECT service, timestamp, duration_ms, trace_id, span_name, status FROM traces WHERE 1=1`
	args := []interface{}{}
	if filter.Service != "" {
		query += " AND service = ?"
		args = append(args, filter.Service)
	}
	if filter.StartTime > 0 {
		query += " AND timestamp >= ?"
		args = append(args, filter.StartTime)
	}
	if filter.EndTime > 0 {
		query += " AND timestamp < ?"
		args = append(args, filter.EndTime)
	}
	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var traces []*ingestion.TraceEntry
	for rows.Next() {
		var t ingestion.TraceEntry
		if err := rows.Scan(&t.Service, &t.Timestamp, &t.DurationMs, &t.TraceID, &t.SpanName, &t.Status); err != nil {
			return nil, err
		}
		traces = append(traces, &t)
	}
	return traces, nil
}

func (r *SQLiteTraceRepository) DeleteOlderThan(ctx context.Context, beforeUnix int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM traces WHERE created_at < ?`, beforeUnix)
	return err
}
