# Olario — Multi-Tenant Grocery Commerce Platform

This repository is in minimal planning/bootstrap mode.
We are intentionally starting small and adding complexity in phases.

## Read First
1. `notes/VISION.md`
2. `notes/main_project.md`
3. `notes/FOLDER_STRUCTURE.md`
4. `notes/ARCHITECTURE.md`
5. `notes/SERVICE_RESPONSIBILITIES.md`
6. `notes/DB_DESIGN.md`
7. `notes/LOCAL_DEPENDENCIES.md`
8. `notes/K8S_PLAN.md`

## Current Goal
- Build a clean minimal base for Java + Go + C++ services behind Kong.
- Use local machine dependencies now.
- Keep project scalable for later VPS and Kubernetes deployment.

## Baseline Stack
- PostgreSQL 17
- Kong latest stable
- Java 21 LTS
- Go 1.26
- Drogon latest stable
- Angular latest LTS
- Node.js v24.15.0

## Later Phases
- Observability: Loki, Prometheus, cAdvisor, Grafana
- Orchestration: Kubernetes

## Security
- `env.example` contains placeholders only.
- Never commit `.env` or real credentials.
