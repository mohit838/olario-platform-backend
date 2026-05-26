# Olario Platform — Folder Structure Guide

Purpose: minimal now, scalable later.  
Last Updated: 2026-05-26

## Current Minimal Structure (Now)
```text
olario-platform/
├── README.md
├── env.example
├── .gitignore
├── notes/
├── services/
│   ├── cms-service/
│   ├── catalog-service/
│   └── heavy-service/
├── gateway/
│   └── kong/
│       ├── services/
│       ├── routes/
│       ├── plugins/
│       └── init/
├── database/
│   ├── cms/{migrations,seeds,schema}/
│   ├── catalog/{migrations,seeds,schema}/
│   └── heavy/{migrations,seeds,schema}/
├── frontend/
│   └── olario-web/
└── docker/
```

## Add Later (When Needed)
- `infra/k8s/` for Kubernetes manifests
- `observability/` for Loki/Prometheus/cAdvisor/Grafana
- `tests/` when service test suites start
- `storage/` only if MinIO bootstrap/config is managed in-repo
- `docs/` when implementation documentation grows beyond notes

## Non-Negotiable Rules
- Database per service ownership
- No direct cross-service DB access
- External traffic via Kong
- Internal communication via gRPC
- Local dependencies (DB/Kafka/MinIO/Redis) for now
