# OLARIO — Developer Roadmap

**Status:** R&D, Milestone M1 (Platform Skeleton)  
**Last Updated:** 2026-05-26

## 1) What We Are Building
A multi-tenant grocery commerce platform using microservices with strict service boundaries and database-per-service ownership.

## 2) Version Baseline (Required)
- PostgreSQL 17
- Kong latest stable
- Java 21 LTS
- Go 1.26
- Drogon latest stable
- Angular latest LTS
- Node.js v24.15.0

Later add:
- Loki
- Prometheus
- cAdvisor
- Grafana

## 3) Architecture (High-Level)
```text
Angular Frontend
      |
      v
Kong API Gateway (REST)
  |          |           |
  v          v           v
Java CMS    Go Catalog   C++ Heavy (Drogon)
  |          |           |
  v          v           v
cms_db    catalog_db   heavy_db     (all PostgreSQL 17)
```

Communication:
- External: REST via Kong
- Internal: gRPC
- Async events: Kafka

## 4) Service Responsibilities
### CMS Service (Java 21 LTS)
- CMS/admin workflows
- business orchestration
- tenant/user/order domain coordination
- Does not own product catalog computation details

### Catalog Service (Go 1.26)
- products/categories/variants
- read-heavy product APIs
- fast lookup and data serving
- Does not own order or payment orchestration

### Heavy Service (C++ Drogon)
- pricing/discount/scoring/report computations
- compute-heavy asynchronous logic
- Does not own catalog CRUD or CMS workflows

## 5) Where Developers Start (Step-by-Step)
1. Read `README.md` for quick boot.
2. Read `notes/FOLDER_STRUCTURE.md` for folder ownership.
3. Copy `env.example` to `.env` and fill safe local values.
4. Use local machine dependencies (PostgreSQL/Kafka/MinIO/Redis) and start only project app services.
5. Verify health endpoints through Kong.
6. Pick one service folder and work only in that ownership boundary.

## 6) Milestones
### M1: Platform Skeleton
- Run frontend, gateway, and all three services using local machine dependencies
- Health endpoints return `200` through Kong
- Database connectivity working for each service

### M2-M8: Feature Expansion
- inventory, cart, order, auth, multi-tenancy, payment, analytics

### M9: Observability (later)
- Add Loki + Prometheus + cAdvisor + Grafana
- Build dashboards + alert baselines

## 7) Done Criteria for Contributors
- Code is in correct service boundary
- No secret leakage in code/docs
- Migrations included where schema changes
- Health checks pass
- Docs updated if architecture/folder ownership changed
