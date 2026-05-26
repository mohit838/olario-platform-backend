# Olario — Multi-Tenant Grocery Commerce Platform

This repository is currently in planning and structure phase.
The goal is to establish a clear, scalable foundation before full implementation.

## Start Here
1. Read `notes/VISION.md`
2. Read `notes/main_project.md`
3. Read `notes/FOLDER_STRUCTURE.md`
4. Read `notes/ARCHITECTURE.md`
5. Read `notes/SERVICE_RESPONSIBILITIES.md`
6. Read `notes/DB_DESIGN.md`
7. Read `notes/STACK.md`
8. Read `notes/LOCAL_DEPENDENCIES.md`

## Baseline Stack
- PostgreSQL 17
- Kong latest stable
- Java 21 LTS
- Go 1.26
- Drogon latest stable
- Angular latest LTS
- Node.js v24.15.0

## Planned Later
- Loki
- Prometheus
- cAdvisor
- Grafana

## Current State
- Folder structure scaffold created for scale.
- Service folders created (empty) for later initialization.
- Documentation organized in `notes/` for roadmap-driven development.
- DB/Kafka/MinIO are expected from your local machine environment.

## Security
- `env.example` uses placeholders only.
- Never commit `.env` or real credentials.
