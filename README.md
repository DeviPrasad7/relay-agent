# 🔍 Relay — Distributed Incident Analysis System

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg?style=for-the-badge\&logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg?style=for-the-badge\&logo=docker)](https://www.docker.com/)
[![Grafana](https://img.shields.io/badge/Grafana-Dashboard-F46800.svg?style=for-the-badge\&logo=grafana)](https://grafana.com/)

### AI-powered incident detection, correlation, and root-cause analysis for distributed systems

Relay continuously ingests telemetry from distributed services, detects anomalies, correlates related events, and generates actionable root-cause analysis using pluggable Large Language Models.

</div>

---

## Overview

Modern distributed systems generate massive volumes of logs, traces, metrics, and deployment events. Identifying the actual cause of an incident often requires manually correlating signals across multiple services and infrastructure components.

Relay automates this process by:

* Collecting operational telemetry
* Detecting anomalous behavior
* Correlating related events into incidents
* Enriching incidents with deployment context
* Generating structured root-cause analysis using AI

The platform is designed around clean architecture principles, making every component independently replaceable and extensible.

---

## Features

### Telemetry Ingestion

* HTTP API for log ingestion
* Distributed trace ingestion
* Deployment event tracking
* Service heartbeat monitoring
* Input validation and normalization

### Anomaly Detection

Built-in detectors include:

* Error-rate spike detection
* Latency degradation detection
* Heartbeat failure detection

Additional detectors can be added through the detector interface without modifying orchestration logic.

### Incident Correlation

* Temporal grouping of related anomalies
* Deployment-aware incident enrichment
* Cross-service incident consolidation
* Unified incident timeline generation

### AI-Assisted Root Cause Analysis

* Structured RCA generation
* OpenAI-compatible provider abstraction
* Support for OpenAI, Groq, DeepSeek, Ollama, and compatible APIs
* SHA-256 response caching
* Configurable prompt generation

### Reporting & Observability

* Built-in web dashboard
* Incident history view
* Service health overview
* Grafana integration
* JSON incident reports
* Markdown incident reports

### Reliability

* Idempotent processing pipeline
* Persistent orchestration state
* Safe restarts
* Docker-native deployment
* Minimal operational footprint

---

## Architecture

```text
┌─────────────────────────────────────────────┐
│              Telemetry Sources              │
├─────────────────────────────────────────────┤
│ Logs │ Traces │ Deployments │ Heartbeats    │
└─────────────────────┬───────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────┐
│               Relay Ingestion               │
└─────────────────────┬───────────────────────┘
                      ▼
┌─────────────────────────────────────────────┐
│                 SQLite Store                │
└─────────────────────┬───────────────────────┘
                      ▼
┌─────────────────────────────────────────────┐
│            Anomaly Detection Layer          │
└─────────────────────┬───────────────────────┘
                      ▼
┌─────────────────────────────────────────────┐
│            Incident Correlation             │
└─────────────────────┬───────────────────────┘
                      ▼
┌─────────────────────────────────────────────┐
│             LLM RCA Generation              │
└─────────────────────┬───────────────────────┘
                      ▼
┌─────────────────────────────────────────────┐
│ Dashboard │ Reports │ Grafana │ APIs        │
└─────────────────────────────────────────────┘
```

---

## Detection Pipeline

The orchestration engine executes on a configurable interval (default: 60 seconds).

1. Retrieve newly ingested telemetry.
2. Execute anomaly detectors.
3. Group anomalies into incidents.
4. Correlate incidents with recent deployment events.
5. Generate root-cause analysis using the configured LLM provider.
6. Cache repeated analysis results.
7. Persist incidents and generated reports.
8. Expose results through APIs and dashboards.

---

## Example Incident

### Detected Incident

| Field      | Value                   |
| ---------- | ----------------------- |
| Service    | auth-service            |
| Type       | Error Spike             |
| Error Rate | 83%                     |
| Timestamp  | 2026-06-08 21:43:28 UTC |

### Correlated Events

* payment-service latency degradation
* deployment event detected 78 seconds earlier

### Generated Root Cause Analysis

> The auth-service error spike and payment-service latency degradation occurred within the same correlation window, indicating a likely shared dependency failure. The most probable root cause is exhaustion or degradation of a shared datastore, resulting in increased request latency and elevated error rates.

### Recommended Actions

1. Inspect datastore health metrics.
2. Review recent deployments.
3. Verify connection pool utilization.
4. Search logs for timeout or connection-refused errors.
5. Validate network connectivity between affected services.

---

## Technology Stack

| Component          | Technology                                   |
| ------------------ | -------------------------------------------- |
| Language           | Go 1.23                                      |
| Storage            | SQLite (WAL Mode)                            |
| Detection Engine   | Statistical anomaly detection                |
| Correlation Engine | Time-window and deployment-aware correlation |
| AI Layer           | OpenAI-compatible abstraction                |
| Dashboard          | HTML, CSS, JavaScript                        |
| Visualization      | Chart.js                                     |
| Monitoring         | Grafana                                      |
| Containerization   | Docker & Docker Compose                      |

---

## Getting Started

### Prerequisites

* Docker
* Docker Compose

### Clone Repository

```bash
git clone https://github.com/DeviPrasad7/relay-agent.git
cd relay-agent
```

### Configure Environment

```bash
echo "LLM_PROVIDER=mock" > .env
```

### Start Services

```bash
docker compose -f docker-compose.grafana.yml up --build -d
```

### Access Services

| Service       | URL                   |
| ------------- | --------------------- |
| Web Dashboard | http://localhost:8080 |
| Grafana       | http://localhost:3000 |

---

## Sample Telemetry Ingestion

### Log Event

```bash
curl -X POST http://localhost:8080/ingest/log \
-H "Content-Type: application/json" \
-d '{
  "service":"auth",
  "timestamp":1710000000,
  "severity":"error",
  "message":"database timeout"
}'
```

### Deployment Event

```json
{
  "service": "payment",
  "version": "v1.2.4",
  "timestamp": 1710000000
}
```

---

## AI Analysis Configuration

Relay supports any OpenAI-compatible inference endpoint.

### Mock Provider

```env
LLM_PROVIDER=mock
```

### Groq

```env
LLM_PROVIDER=custom
LLM_BASE_URL=https://api.groq.com/openai/v1
LLM_API_KEY=<api-key>
LLM_MODEL=mixtral-8x7b-32768
```

### OpenAI

```env
LLM_PROVIDER=openai
OPENAI_API_KEY=<api-key>
LLM_MODEL=gpt-4o-mini
```

### DeepSeek

```env
LLM_PROVIDER=deepseek
DEEPSEEK_API_KEY=<api-key>
LLM_MODEL=deepseek-chat
```

### Ollama

```env
LLM_PROVIDER=custom
LLM_BASE_URL=http://localhost:11434/v1
LLM_API_KEY=local
LLM_MODEL=llama3
```

After updating configuration:

```bash
docker compose -f docker-compose.grafana.yml up --build -d
```

---

## Project Structure

```text
relay-agent/
├── cmd/
│   └── relay/
│       └── main.go
│
├── internal/
│   ├── ingestion/
│   ├── storage/
│   ├── detection/
│   ├── correlation/
│   ├── llm/
│   ├── alerting/
│   ├── orchestration/
│   ├── api/
│   └── config/
│
├── web/
├── grafana/
├── scripts/
│   └── init_db.sql
│
├── docker-compose.grafana.yml
└── README.md
```

---

## Extending Relay

### Custom Detectors

Implement the detector interface and register the detector during initialization.

### Custom AI Providers

Implement the analyzer interface or connect any OpenAI-compatible endpoint.

### Alternative Storage Backends

The repository abstraction allows migration from SQLite to PostgreSQL or other databases without modifying business logic.

---

## Roadmap

* Kubernetes auto-remediation workflows
* Native Prometheus exporter
* Additional anomaly detectors
* Slack notifications
* PagerDuty integration
* Discord notifications
* Multi-tenant support
* Role-based access control (RBAC)
* Distributed storage backends

---

## Contributing

Contributions are welcome.

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Submit a pull request

For major changes, open an issue first to discuss the proposed design.

---
