# Developer Guide

Last Updated: 2026-05-26

## Start Order (Top Priority)

1. Start local dependencies first (PostgreSQL, Redis, MinIO, Kafka/Redpanda).
2. Start `cms-service` (Java).
3. Start `catalog-service` (Go).
4. Start `heavy-service` (C++).
5. Verify each `/health` endpoint.

## 1) Prerequisites

- Java 21
- Maven (or `./mvnw` wrapper)
- Go 1.26
- CMake + C++ compiler (for heavy service)
- Local dependency containers running

## 2) Start CMS Service (Java Spring Boot)

Path:

- `services/cms-service/olario`

Run:

```bash
cd services/cms-service/olario
./mvnw spring-boot:run
```

Expected:

- Port: `8081`
- App name: `cms-service`

Health check:

```bash
curl http://localhost:8081/actuator/health
```

Note:

- If actuator is not added yet, use a custom health endpoint when available.

## 3) Start Catalog Service (Go)

Path:

- `services/catalog-service`

Run:

```bash
cd services/catalog-service
go run ./cmd/server
```

Expected:

- Port: `8082`

Health check:

```bash
curl http://localhost:8082/health
```

## 4) Start Heavy Service (C++)

Path:

- `services/heavy-service`

Build and run:

```bash
cd services/heavy-service
cmake -S . -B build
cmake --build build
./build/heavy-service
```

Expected:

- Current bootstrap prints startup message (Drogon wiring can be added next)
- Target port convention: `8083`

## 5) Database Setup (CMS First)

Create DB:

```bash
psql -U postgres -h localhost -p 5432 -c "CREATE DATABASE olario_cms_db;"
```

Current CMS DB config is in:

- `services/cms-service/olario/src/main/resources/application.yaml`

## 6) Quick Troubleshooting

- Port in use: change port in service config/env and restart.
- DB auth failed: verify local postgres username/password in `application.yaml`.
- Go not running: check `go version` and run command from `services/catalog-service`.
- C++ build fails: verify `cmake --version` and compiler (`g++ --version`).

## 7) Next Step After Local Start

## 5) Database Setup (CMS First)

Create DB:

```bash
psql -U postgres -h localhost -p 5432 -c "CREATE DATABASE olario_cms_db;"
```

Current CMS DB config is in:

- `services/cms-service/olario/src/main/resources/application.yaml`

## 6) Quick Troubleshooting

- Port in use: change port in service config/env and restart.
- DB auth failed: verify local postgres username/password in `application.yaml`.
- Go not running: check `go version` and run command from `services/catalog-service`.
- C++ build fails: verify `cmake --version` and compiler (`g++ --version`).

## 7) Next Step After Local Start

- Add Kong routes for `cms-service`, `catalog-service`, and `heavy-service`.
- Then validate API flow through gateway.
