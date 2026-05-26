# Olario Platform — Folder Structure Guide

**Purpose:** This is the canonical structure reference for developers.  
**Last Updated:** 2026-05-26

## 1. Root Layout
```text
olario-platform/
├── README.md                    # Fast onboarding and first commands
├── env.example                  # Safe environment template (no real secrets)
├── .gitignore                   # Secret and build artifact protection
├── notes/                       # All project documentation (roadmap/architecture/etc.)
│
├── docs/                        # Architecture docs, ADRs, API docs, runbooks
├── docker/                      # Docker compose and image build definitions
├── gateway/                     # Kong gateway config (services/routes/plugins)
├── frontend/                    # Angular app (latest LTS, Node v24.15.0)
├── services/                    # Java/Go/C++ microservices
├── database/                    # Per-service migrations/seeds/schema docs
├── storage/                     # MinIO bootstrap, buckets, policies
├── observability/               # Later: Loki/Prometheus/cAdvisor/Grafana
├── scripts/                     # Setup/dev/ci/release helper scripts
└── tests/                       # Integration, e2e, performance, contract tests
```

## 2. Service Folder Ownership
```text
services/
├── cms-service/                 # Java 21 LTS + Spring Boot
│   ├── src/main/java/...        # API, business logic, repositories
│   ├── src/main/resources/      # Config + migration hooks
│   ├── src/test/java/...        # Unit/integration tests
│   ├── Dockerfile
│   └── README.md
│
├── catalog-service/             # Go 1.26
│   ├── cmd/server/main.go       # Entry point
│   ├── internal/                # Domain logic (handler/service/repository)
│   ├── pkg/                     # Optional reusable public package
│   ├── migrations/
│   ├── Dockerfile
│   └── README.md
│
└── heavy-service/               # C++ Drogon (latest stable)
    ├── src/                     # Controllers, services, repository, grpc
    ├── config/config.json
    ├── CMakeLists.txt
    ├── Dockerfile
    └── README.md
```

## 3. Database Structure (Database-Per-Service)
```text
database/
├── cms/
│   ├── migrations/
│   ├── seeds/
│   └── schema/
├── catalog/
│   ├── migrations/
│   ├── seeds/
│   └── schema/
├── heavy/
│   ├── migrations/
│   ├── seeds/
│   └── schema/
└── scripts/
    ├── init-all.sh
    ├── migrate-all.sh
    ├── seed-all.sh
    └── reset-all.sh
```

## 4. Gateway Structure
```text
gateway/kong/
├── kong.yml
├── services/                    # Upstream service definitions
├── routes/                      # Public route mappings
├── plugins/                     # CORS/rate-limit/logging/auth plugins
└── init/                        # Gateway bootstrap scripts
```

## 5. Frontend Structure
```text
frontend/olario-web/
├── src/app/
│   ├── core/                    # App-wide singleton services/guards/interceptors
│   ├── shared/                  # Shared components/directives/pipes
│   ├── layout/                  # Layout shells
│   └── features/                # Feature modules (cms/catalog/storefront/admin)
├── src/environments/
├── src/assets/
├── src/styles/
├── package.json
└── README.md
```

## 6. Later Observability Structure
```text
observability/
├── loki/
├── prometheus/
├── cadvisor/
└── grafana/
```

## 7. Non-Negotiable Rules
- No direct cross-service database access.
- Each service owns its data, migrations, and domain logic.
- Internal service communication: gRPC.
- External client traffic: REST through Kong.
- PostgreSQL, Kafka, MinIO, Redis are provided from local machine containers.
- `env.example` uses placeholders only.
- Never commit `.env` or secret files.
