# Architecture

Olario starts as a modular monolith. The goal is to keep one deployable Go
service while preserving clear boundaries so pieces can be extracted later if
the project grows.

## Current Layers

```text
cmd/
  api/                 process entrypoint for HTTP API
  migrate/             process entrypoint for database migrations
  seed/                process entrypoint for local/test seed data

internal/
  app/                 application bootstrapping and dependency wiring
  application/         use cases and application services
  config/              YAML/env/flag config loading
  domain/              business types and rules
  http/                router, handlers, middleware, DTOs
  platform/            Postgres, Redis, logger, and future adapters

docs/                  architecture and developer guides
cmd/migrate/migrations SQL migration files
```

Docker Compose currently runs only the API container. Local Postgres, Redis, and
MinIO are expected to run outside Compose and be provided through private config
or environment variables.

## Dependency Rule

- `domain` must not import HTTP, Postgres, Redis, or framework packages.
- `application` coordinates use cases and depends on interfaces.
- `platform` implements adapters for databases, cache, storage, and email.
- `http` translates HTTP requests into application calls.
- `app` wires concrete implementations together.

## Current Request Flow

```text
client
  -> internal/http router
  -> handler
  -> internal/application service
  -> repository/cache interfaces
  -> internal/platform/postgres or internal/platform/redis
  -> response DTO
```

## Auth Request Flow

```text
POST /api/v1/auth/login
  -> internal/http AuthHandler
  -> internal/application/auth Service
  -> internal/platform/postgres AuthRepository
  -> bcrypt password verification
  -> internal/security/token access token signer
  -> internal/platform/redis refresh token store
```

## Why This Shape

This shape keeps code loosely coupled:

- handlers do not contain SQL
- services do not know Chi router details
- domain types do not know where data is stored
- infrastructure can be replaced later

## Later Boundaries

Good future extraction candidates:

- auth and identity
- catalog and inventory
- ordering and invoices
- reports
- billing/subscriptions
- notification/email
