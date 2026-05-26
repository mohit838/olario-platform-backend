# OLARIO — Developer Roadmap

Status: R&D, minimal bootstrap phase  
Last Updated: 2026-05-26

## Vision
Start minimal, move step-by-step, preserve scalability from day one.

## Version Baseline
- PostgreSQL 17
- Kong latest stable
- Java 21 LTS
- Go 1.26
- Drogon latest stable
- Angular latest LTS
- Node.js v24.15.0

Later:
- Loki
- Prometheus
- cAdvisor
- Grafana
- Kubernetes

## Architecture
Angular -> Kong -> CMS/Go/Heavy services

Rules:
- DB per service
- No cross-service DB queries
- Internal gRPC
- External REST via Kong

## Start Sequence
1. Prepare local dependencies (Postgres/Redis/MinIO/Kafka-compatible broker)
2. Initialize one service at a time
3. Add health endpoints and route through Kong
4. Add service-specific migrations
5. Add frontend integration

## Milestone Path
- M1: services boot + health checks + DB connectivity
- M2+: domain features and integrations
- M9: observability stack
- M10: Kubernetes migration (after stability)
