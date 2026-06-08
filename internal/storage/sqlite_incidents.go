package storage

import (
	"context"
	"database/sql"
)

type SQLiteIncidentRepository struct {
	db *sql.DB
}

func NewSQLiteIncidentRepository(db *sql.DB) *SQLiteIncidentRepository {
	return &SQLiteIncidentRepository{db: db}
}

func (r *SQLiteIncidentRepository) Store(ctx context.Context, incident *Incident) error {
	var deploymentID interface{}
	if incident.DeploymentEventID != nil {
		deploymentID = *incident.DeploymentEventID
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO incidents(incident_time, detection_method, anomaly_details, correlated_anomalies, deployment_event_id, root_cause_summary, root_cause_full, resolution_status, created_at)
		 VALUES(?,?,?,?,?,?,?,?,?)`,
		incident.IncidentTime, incident.DetectionMethod, incident.AnomalyDetails, incident.CorrelatedAnomalies,
		deploymentID, incident.RootCauseSummary, incident.RootCauseFull, incident.ResolutionStatus, incident.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	incident.ID = int(id)
	return nil
}

func (r *SQLiteIncidentRepository) List(ctx context.Context, limit int) ([]*Incident, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, incident_time, detection_method, anomaly_details, correlated_anomalies, deployment_event_id, root_cause_summary, root_cause_full, resolution_status, created_at
		 FROM incidents ORDER BY incident_time DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var incidents []*Incident
	for rows.Next() {
		var inc Incident
		var deploymentID sql.NullInt64
		err := rows.Scan(&inc.ID, &inc.IncidentTime, &inc.DetectionMethod, &inc.AnomalyDetails, &inc.CorrelatedAnomalies,
			&deploymentID, &inc.RootCauseSummary, &inc.RootCauseFull, &inc.ResolutionStatus, &inc.CreatedAt)
		if err != nil {
			return nil, err
		}
		if deploymentID.Valid {
			id := int(deploymentID.Int64)
			inc.DeploymentEventID = &id
		}
		incidents = append(incidents, &inc)
	}
	return incidents, nil
}

func (r *SQLiteIncidentRepository) GetByID(ctx context.Context, id int) (*Incident, error) {
	var inc Incident
	var deploymentID sql.NullInt64
	err := r.db.QueryRowContext(ctx,
		`SELECT id, incident_time, detection_method, anomaly_details, correlated_anomalies, deployment_event_id, root_cause_summary, root_cause_full, resolution_status, created_at
		 FROM incidents WHERE id = ?`, id).
		Scan(&inc.ID, &inc.IncidentTime, &inc.DetectionMethod, &inc.AnomalyDetails, &inc.CorrelatedAnomalies,
			&deploymentID, &inc.RootCauseSummary, &inc.RootCauseFull, &inc.ResolutionStatus, &inc.CreatedAt)
	if err != nil {
		return nil, err
	}
	if deploymentID.Valid {
		id := int(deploymentID.Int64)
		inc.DeploymentEventID = &id
	}
	return &inc, nil
}
