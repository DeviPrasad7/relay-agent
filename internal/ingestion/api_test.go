package ingestion

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockLogStore struct {
	storeFunc func(context.Context, *LogEntry) error
}
func (m *mockLogStore) Store(ctx context.Context, log *LogEntry) error {
	if m.storeFunc != nil {
		return m.storeFunc(ctx, log)
	}
	return nil
}

type mockTraceStore struct {
	storeFunc func(context.Context, *TraceEntry) error
}
func (m *mockTraceStore) Store(ctx context.Context, trace *TraceEntry) error {
	if m.storeFunc != nil {
		return m.storeFunc(ctx, trace)
	}
	return nil
}

type mockEventStore struct {
	storeFunc func(context.Context, *EventEntry) error
}
func (m *mockEventStore) Store(ctx context.Context, event *EventEntry) error {
	if m.storeFunc != nil {
		return m.storeFunc(ctx, event)
	}
	return nil
}

func TestIngestLog(t *testing.T) {
	handler := NewHandler(&mockLogStore{}, nil, nil)
	server := httptest.NewServer(http.HandlerFunc(handler.IngestLog))
	defer server.Close()

	logJSON := `{"service":"auth","timestamp":1700000000,"severity":"error","message":"fail"}`
	resp, err := http.Post(server.URL, "application/json", bytes.NewBufferString(logJSON))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// Bad request
	resp2, _ := http.Post(server.URL, "application/json", bytes.NewBufferString(`{"service":""}`))
	if resp2.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp2.StatusCode)
	}
}

func TestIngestTrace(t *testing.T) {
	handler := NewHandler(nil, &mockTraceStore{}, nil)
	server := httptest.NewServer(http.HandlerFunc(handler.IngestTrace))
	defer server.Close()

	traceJSON := `{"service":"pay","timestamp":1700000000,"duration_ms":100,"trace_id":"abc","span_name":"charge","status":"ok"}`
	resp, _ := http.Post(server.URL, "application/json", bytes.NewBufferString(traceJSON))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestIngestEvent(t *testing.T) {
	handler := NewHandler(nil, nil, &mockEventStore{})
	server := httptest.NewServer(http.HandlerFunc(handler.IngestEvent))
	defer server.Close()

	eventJSON := `{"service":"auth","timestamp":1700000000,"type":"deploy","version":"v1.0"}`
	resp, _ := http.Post(server.URL, "application/json", bytes.NewBufferString(eventJSON))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
