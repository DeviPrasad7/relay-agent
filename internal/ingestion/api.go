package ingestion

import (
	"context"
	"encoding/json"
	"net/http"
)

// Minimal interfaces that match the storage repository methods.
// These avoid an import cycle (ingestion does NOT import storage).
type LogStorer interface {
	Store(ctx context.Context, log *LogEntry) error
}
type TraceStorer interface {
	Store(ctx context.Context, trace *TraceEntry) error
}
type EventStorer interface {
	Store(ctx context.Context, event *EventEntry) error
}

type Handler struct {
	logRepo   LogStorer
	traceRepo TraceStorer
	eventRepo EventStorer
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
		if err := h.logRepo.Store(r.Context(), &log); err != nil {
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
		if err := h.traceRepo.Store(r.Context(), &trace); err != nil {
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
		if err := h.eventRepo.Store(r.Context(), &event); err != nil {
			http.Error(w, "Storage error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
