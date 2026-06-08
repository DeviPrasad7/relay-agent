package orchestration

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"relay-agent/internal/alerting"
	"relay-agent/internal/correlation"
	"relay-agent/internal/detection"
	"relay-agent/internal/ingestion"
	"relay-agent/internal/llm"
	"relay-agent/internal/storage"
)

// Pipeline orchestrates the entire detection flow.
type Pipeline struct {
	logRepo        storage.LogRepository
	traceRepo      storage.TraceRepository
	eventRepo      storage.EventRepository
	incidentRepo   storage.IncidentRepository
	detectors      []detection.Detector
	heartbeatDetector *detection.HeartbeatFailureDetector
	correlator     correlation.Correlator
	llmAnalyzer    llm.LLMAnalyzer
	jsonReporter   alerting.Reporter
	markdownReporter alerting.Reporter
	cfg            *PipelineConfig
}

// PipelineConfig holds configuration for pipeline.
type PipelineConfig struct {
	TemporalWindowSec   int
	DeploymentWindowSec int
}

// NewPipeline creates a pipeline with all dependencies.
func NewPipeline(
	logRepo storage.LogRepository,
	traceRepo storage.TraceRepository,
	eventRepo storage.EventRepository,
	incidentRepo storage.IncidentRepository,
	detectors []detection.Detector,
	heartbeatDetector *detection.HeartbeatFailureDetector,
	llmAnalyzer llm.LLMAnalyzer,
	jsonReporter alerting.Reporter,
	markdownReporter alerting.Reporter,
	cfg *PipelineConfig,
) *Pipeline {
	// Create a combined correlator that does both temporal grouping and deployment linking
	var correlator correlation.Correlator
	// For MVP, we can use the DeploymentCorrelator which internally uses TemporalCorrelator
	// if we implement its Correlate method properly. Simpler: we create a composite.
	// To avoid overcomplicating, we instantiate a DeploymentCorrelator with windows.
	deployCorr := correlation.NewDeploymentCorrelator(cfg.DeploymentWindowSec)
	// Note: DeploymentCorrelator.Correlate expects anomaly groups already grouped? We'll assume it's the full correlator.
	// Actually from earlier implementation, DeploymentCorrelator.Correlate expects []detection.Anomaly and does grouping itself.
	// We'll use that.
	correlator = deployCorr
	
	return &Pipeline{
		logRepo:        logRepo,
		traceRepo:      traceRepo,
		eventRepo:      eventRepo,
		incidentRepo:   incidentRepo,
		detectors:      detectors,
		heartbeatDetector: heartbeatDetector,
		correlator:     correlator,
		llmAnalyzer:    llmAnalyzer,
		jsonReporter:   jsonReporter,
		markdownReporter: markdownReporter,
		cfg:            cfg,
	}
}

// Run processes data since lastTimestamp (inclusive).
func (p *Pipeline) Run(ctx context.Context, lastTimestamp int64) (int64, error) {
	now := time.Now().Unix()
	logFilter := storage.LogFilter{StartTime: lastTimestamp + 1, EndTime: now}
	logs, err := p.logRepo.Query(ctx, logFilter)
	if err != nil {
		return lastTimestamp, err
	}
	traceFilter := storage.TraceFilter{StartTime: lastTimestamp + 1, EndTime: now}
	traces, err := p.traceRepo.Query(ctx, traceFilter)
	if err != nil {
		return lastTimestamp, err
	}
	eventFilter := storage.EventFilter{StartTime: lastTimestamp + 1, EndTime: now}
	events, err := p.eventRepo.Query(ctx, eventFilter)
	if err != nil {
		return lastTimestamp, err
	}
	if len(logs) == 0 && len(traces) == 0 && len(events) == 0 {
		return now, nil
	}
	// Update heartbeat detector
	for _, ev := range events {
		if ev.Type == ingestion.EventTypeHeartbeat && p.heartbeatDetector != nil {
			p.heartbeatDetector.UpdateHeartbeat(ev.Service, ev.Timestamp)
		}
	}
	// Aggregate metrics (simplified for MVP)
	serviceMetrics := p.aggregateMetrics(logs, traces, now)
	var allAnomalies []detection.Anomaly
	for _, detector := range p.detectors {
		for svc, metrics := range serviceMetrics {
			detCtx := &detection.DetectionContext{
				Service:     svc,
				WindowStart: lastTimestamp + 1,
				WindowEnd:   now,
				Metrics:     metrics.Metrics,
				Baseline:    metrics.Baseline,
			}
			anomalies, err := detector.Detect(ctx, detCtx)
			if err != nil {
				log.Printf("detector error for %s: %v", svc, err)
				continue
			}
			allAnomalies = append(allAnomalies, anomalies...)
		}
	}
	// Heartbeat detector
	if p.heartbeatDetector != nil {
		hbCtx := &detection.DetectionContext{
			Metrics: map[string]interface{}{"current_time": now},
		}
		hbAnomalies, err := p.heartbeatDetector.Detect(ctx, hbCtx)
		if err != nil {
			log.Printf("heartbeat error: %v", err)
		} else {
			allAnomalies = append(allAnomalies, hbAnomalies...)
		}
	}
	if len(allAnomalies) == 0 {
		return now, nil
	}
	// Convert events slice to value slice for correlation input
	depEvents := make([]ingestion.EventEntry, 0, len(events))
	for _, ev := range events {
		if ev != nil {
			depEvents = append(depEvents, *ev)
		}
	}
	corrInput := &correlation.CorrelationInput{
		Anomalies:           allAnomalies,
		DeploymentEvents:    depEvents,
		TimeWindowSec:       p.cfg.TemporalWindowSec,
		DeploymentWindowSec: p.cfg.DeploymentWindowSec,
	}
	groups, err := p.correlator.Correlate(ctx, corrInput)
	if err != nil {
		return now, err
	}
	for _, group := range groups {
		prompt := llm.BuildAnomalyContext(&group)
		rootCauseFull, err := p.llmAnalyzer.Analyze(ctx, prompt)
		if err != nil {
			log.Printf("LLM error: %v", err)
			rootCauseFull = "LLM analysis failed: " + err.Error()
		}
		summary := rootCauseFull
		if len(summary) > 200 {
			summary = summary[:200] + "..."
		}
		anomalyDetailsJSON, _ := json.Marshal(group.Anomalies)
		incident := &storage.Incident{
			IncidentTime:        group.StartTime,
			DetectionMethod:     group.Anomalies[0].Method,
			AnomalyDetails:      string(anomalyDetailsJSON),
			CorrelatedAnomalies: string(anomalyDetailsJSON),
			RootCauseSummary:    summary,
			RootCauseFull:       rootCauseFull,
			ResolutionStatus:    "open",
			CreatedAt:           now,
		}
		if err := p.incidentRepo.Store(ctx, incident); err != nil {
			log.Printf("store incident error: %v", err)
			continue
		}
		if err := p.jsonReporter.Generate(ctx, incident); err != nil {
			log.Printf("JSON report error: %v", err)
		}
		if err := p.markdownReporter.Generate(ctx, incident); err != nil {
			log.Printf("Markdown report error: %v", err)
		}
	}
	return now, nil
}

type metricSnapshot struct {
	Metrics  map[string]interface{}
	Baseline map[string]float64
}

func (p *Pipeline) aggregateMetrics(logs []*ingestion.LogEntry, traces []*ingestion.TraceEntry, now int64) map[string]*metricSnapshot {
	result := make(map[string]*metricSnapshot)
	serviceLogs := make(map[string]int)
	serviceErrors := make(map[string]int)
	for _, l := range logs {
		serviceLogs[l.Service]++
		if l.Severity == ingestion.SeverityError {
			serviceErrors[l.Service]++
		}
	}
	for svc, total := range serviceLogs {
		errorRate := float64(serviceErrors[svc]) / float64(total)
		result[svc] = &metricSnapshot{
			Metrics:  map[string]interface{}{"error_rate": errorRate},
			Baseline: map[string]float64{"mean": 0.02, "std_dev": 0.01},
		}
	}
	serviceTraces := make(map[string][]int)
	for _, t := range traces {
		serviceTraces[t.Service] = append(serviceTraces[t.Service], t.DurationMs)
	}
	for svc, durations := range serviceTraces {
		p95 := detection.ComputeP95FromDurations(durations)
		if _, exists := result[svc]; !exists {
			result[svc] = &metricSnapshot{
				Metrics:  map[string]interface{}{},
				Baseline: map[string]float64{"mean": 100, "std_dev": 50},
			}
		}
		result[svc].Metrics["p95_current"] = p95
		result[svc].Metrics["p95_baseline"] = 100.0 // default
	}
	return result
}
