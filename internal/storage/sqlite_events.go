package storage

import (
	"context"
	"database/sql"
	"relay-agent/internal/ingestion"
)

type SQLiteEventRepository struct {
	db *sql.DB
}

func NewSQLiteEventRepository(db *sql.DB) *SQLiteEventRepository {
	return &SQLiteEventRepository{db: db}
}

func (r *SQLiteEventRepository) Store(ctx context.Context, event *ingestion.EventEntry) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO events(service, timestamp, type, version, metadata, created_at) VALUES(?,?,?,?,?,?)`,
		event.Service, event.Timestamp, event.Type, event.Version, event.Metadata, event.Timestamp,
	)
	return err
}

func (r *SQLiteEventRepository) Query(ctx context.Context, filter EventFilter) ([]*ingestion.EventEntry, error) {
	query := `SELECT service, timestamp, type, version, metadata FROM events WHERE 1=1`
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
	if filter.Type != "" {
		query += " AND type = ?"
		args = append(args, filter.Type)
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []*ingestion.EventEntry
	for rows.Next() {
		var e ingestion.EventEntry
		if err := rows.Scan(&e.Service, &e.Timestamp, &e.Type, &e.Version, &e.Metadata); err != nil {
			return nil, err
		}
		events = append(events, &e)
	}
	return events, nil
}

func (r *SQLiteEventRepository) DeleteOlderThan(ctx context.Context, beforeUnix int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM events WHERE created_at < ?`, beforeUnix)
	return err
}
