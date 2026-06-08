# Relay Agent – Automated Incident Detection

See [PLAN.md](progress/PLAN.md) for full design and phase breakdown.

## Quick Start
1. Copy `configs/config.yaml` and set `DEEPSEEK_API_KEY` environment variable.
2. Run `docker-compose up --build`
3. Ingest data: `curl -X POST http://localhost:8080/ingest/log -d '...'`
4. View incidents: `ls data/incidents/`

## Development
- `make test` – run all tests
- `make build` – build binary
- `make run` – run locally (requires SQLite)
