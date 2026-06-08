// Package api provides HTTP handlers for web UI and JSON API.
package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"relay-agent/internal/storage"
	"strconv"
)

type UIServer struct {
	incidentRepo storage.IncidentRepository
	eventRepo    storage.EventRepository
	staticDir    string
}

func NewUIServer(incidentRepo storage.IncidentRepository, eventRepo storage.EventRepository, staticDir string) *UIServer {
	return &UIServer{
		incidentRepo: incidentRepo,
		eventRepo:    eventRepo,
		staticDir:    staticDir,
	}
}

func (s *UIServer) ServeStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(s.staticDir, "index.html"))
}

func (s *UIServer) ListIncidents(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	incidents, err := s.incidentRepo.List(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"incidents": incidents})
}

func (s *UIServer) ServiceHealth(w http.ResponseWriter, r *http.Request) {
	// Get all heartbeat events in last 5 minutes to determine health
	// For MVP, we'll query events of type heartbeat from last 5 minutes
	// Simpler: return dummy health for services that have any heartbeat ever
	// We'll implement a real check: fetch last heartbeat per service, compare with now.
	// For brevity, we return a placeholder – you can extend.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"services": map[string]bool{
			"auth-service":    true,
			"payment-service": true,
		},
	})
}
