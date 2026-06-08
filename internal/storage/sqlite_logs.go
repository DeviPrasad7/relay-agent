package storage

import (
	"context"
	"database/sql"
	"relay-agent/internal/ingestion"
)

type SQLiteLogRepository struct {
	db *sql.DB
}

func NewSQLiteLogRepository(db *sql.DB) *SQLiteLogRepository {
	return &SQLiteLogRepository{db: db}
}

func (r *SQLiteLogRepository) Store(ctx context.Context, log *ingestion.LogEntry) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO logs(service, timestamp, severity, message, created_at) VALUES(?, ?, ?, ?, ?)`,
		log.Service, log.Timestamp, log.Severity, log.Message, log.Timestamp, // using event timestamp as created_at for simplicity
	)
	return err
}

func (r *SQLiteLogRepository) Query(ctx context.Context, filter LogFilter) ([]*ingestion.LogEntry, error) {
	query := `SELECT service, timestamp, severity, message FROM logs WHERE 1=1`
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
	if filter.Severity != "" {
		query += " AND severity = ?"
		args = append(args, filter.Severity)
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []*ingestion.LogEntry
	for rows.Next() {
		var l ingestion.LogEntry
		if err := rows.Scan(&l.Service, &l.Timestamp, &l.Severity, &l.Message); err != nil {
			return nil, err
		}
		logs = append(logs, &l)
	}
	return logs, nil
}

func (r *SQLiteLogRepository) DeleteOlderThan(ctx context.Context, beforeUnix int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM logs WHERE created_at < ?`, beforeUnix)
	return err
}
