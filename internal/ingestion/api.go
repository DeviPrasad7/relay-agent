// Package ingestion provides HTTP handlers for log/trace/event ingestion.
//
// Responsibilities:
//   - Parse JSON requests
//   - Validate using validator
//   - Store via repository interface (to be implemented in Phase 2)
//
// Dependencies: storage.LogRepository, storage.TraceRepository, storage.EventRepository (interfaces)
//
// Last updated: 2026-06-09
package ingestion

import (
	"encoding/json"
	"net/http"
)

// Handler struct holds repository dependencies (interfaces from Phase 2)
// For Phase 1, we use nil placeholders – tests will mock.
type Handler struct {
	logRepo   LogStorer   // will be replaced with storage.LogRepository
	traceRepo TraceStorer
	eventRepo EventStorer
}

// LogStorer is a minimal interface for Phase 1 testing (to avoid importing storage package yet)
type LogStorer interface {
	StoreLog(log *LogEntry) error
}

type TraceStorer interface {
	StoreTrace(trace *TraceEntry) error
}

type EventStorer interface {
	StoreEvent(event *EventEntry) error
}

func NewHandler(logRepo LogStorer, traceRepo TraceStorer, eventRepo EventStorer) *Handler {
	return &Handler{
		logRepo:   logRepo,
		traceRepo: traceRepo,
		eventRepo: eventRepo,
	}
}

func (h *Handler) IngestLog(w http.ResponseWriter, r *http.Request) {
	var log LogEntry
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := ValidateLog(&log); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if h.logRepo != nil {
		if err := h.logRepo.StoreLog(&log); err != nil {
			http.Error(w, "Storage error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) IngestTrace(w http.ResponseWriter, r *http.Request) {
	var trace TraceEntry
	if err := json.NewDecoder(r.Body).Decode(&trace); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := ValidateTrace(&trace); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if h.traceRepo != nil {
		if err := h.traceRepo.StoreTrace(&trace); err != nil {
			http.Error(w, "Storage error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) IngestEvent(w http.ResponseWriter, r *http.Request) {
	var event EventEntry
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := ValidateEvent(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if h.eventRepo != nil {
		if err := h.eventRepo.StoreEvent(&event); err != nil {
			http.Error(w, "Storage error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
