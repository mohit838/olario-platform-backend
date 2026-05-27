# Olario Platform Backend

Olario is a multi-tenant backend platform for grocery and commerce workflows.
The project starts as a simple Go service and will grow gradually into a
production-ready API with clean boundaries, tenant-aware data access, and
replaceable infrastructure integrations.

This repository is intentionally built step by step. The first goal is not to
add every technology at once, but to create a clear foundation that is easy to
understand, test, and extend.

## Working Principles

- Keep each step small and reviewable.
- Prefer simple, explicit code over early abstraction.
- Keep business logic separate from infrastructure details.
- Use interfaces at boundaries where replacement is likely.
- Add tests or manual verification notes with every implementation step.
- Keep public documentation free from private local values.

## Current Status

The project is currently a minimal Go backend starter.

The next milestone is the HTTP foundation: application bootstrap, routing,
configuration, logging, health checks, graceful shutdown, and basic rate
limiting.

## Requirements

- Go `1.26.2` or newer
- Air, optional for live reload
- PostgreSQL, planned for persistent data
- Redis, planned for cache, sessions, or rate limiting
- MinIO, planned for object storage
- Docker and Docker Compose, planned for local infrastructure
- Kubernetes or Minikube, planned for later deployment practice

## Configuration

Runtime configuration should be explicit and safe for public repositories.
Never commit real usernames, passwords, access keys, tokens, or private service
URLs.

Use command-line flags for non-secret runtime options:

```bash
go run ./cmd --env=local --http-addr=:8080
```

Examples of safe command-line flags:

- `--env`
- `--config`
- `--http-addr`

Use environment variables for secrets and machine-specific values:

```text
DATABASE_URL=postgres://<user>:<password>@<host>:<port>/<database>?sslmode=disable
REDIS_ADDR=<host>:<port>
REDIS_USERNAME=<username>
REDIS_PASSWORD=<password>
MINIO_ENDPOINT=<host>:<port>
MINIO_ACCESS_KEY=<access-key>
MINIO_SECRET_KEY=<secret-key>
MINIO_BUCKET=<bucket-name>
MINIO_USE_SSL=false
SMTP_PASSWORD=<password>
```

For local development, create private environment files outside version control
or use your shell environment. Public documentation should show placeholders
only.

Recommended private local files, when needed:

- `.env.local`
- `.env.development`
- local shell profile exports

These files should not be committed with real values.

## Architecture Direction

Olario will begin as a modular monolith. That keeps the codebase simple while
still preparing it for future service boundaries.

Planned package direction:

```text
.
├── cmd/                  # application entrypoints
├── internal/app/         # application bootstrap and lifecycle
├── internal/config/      # config loading and validation
├── internal/domain/      # business concepts and rules
├── internal/http/        # REST handlers, routing, middleware
└── internal/platform/    # infrastructure adapters
```

Infrastructure integrations should stay replaceable:

- PostgreSQL behind repository interfaces
- Redis behind cache, session, or rate-limit interfaces
- MinIO behind an object storage interface
- Email providers behind a mailer interface
- Logging behind simple application-level usage patterns

This keeps the project loosely coupled and easier to scale later.

## Step-by-Step Roadmap

### Step 1: Project Roadmap

Goal: document the intended direction before adding more code.

Output:

- public-safe README
- clear local development expectations
- architecture and infrastructure direction
- small implementation milestones

### Step 2: HTTP Foundation

Goal: create a stable web service foundation.

Output:

- application bootstrap
- config loader
- structured logger
- Chi router
- `/healthz` endpoint
- graceful shutdown with `context.Context`
- basic request timeout and rate limiting middleware

### Step 3: PostgreSQL Foundation

Goal: add persistent storage safely.

Output:

- database configuration
- connection pool
- connectivity check
- migration directory
- `golang-migrate` setup
- first schema migration

### Step 4: Tenant Foundation

Goal: make the application tenant-aware before adding business data.

Output:

- tenant ID type or value object
- tenant-aware request context
- basic tenant validation
- clear rule for how tenant ID enters the API

### Step 5: Redis Foundation

Goal: add Redis as an optional infrastructure adapter.

Output:

- Redis configuration
- connectivity check
- cache or rate-limit interface
- Redis implementation behind that interface

### Step 6: MinIO Foundation

Goal: prepare object storage without coupling business logic to MinIO directly.

Output:

- object storage configuration
- storage interface
- MinIO adapter
- placeholder bucket strategy

### Step 7: Grocery Domain V1

Goal: design the first real business capability.

Output:

- product model
- category model
- tenant-aware repository contracts
- REST API design for products and categories

### Step 8: Users and Access Control

Goal: support tenant-aware users later, after the domain foundation is clear.

Output:

- user model
- role and permission direction
- authentication strategy
- authorization middleware plan

### Step 9: Email Sender

Goal: add email without locking the app to one provider.

Output:

- mailer interface
- provider configuration
- development-safe test sender
- production provider later

### Step 10: Local Infrastructure

Goal: make the project easy to run locally.

Output:

- Dockerfile
- Docker Compose setup
- PostgreSQL service
- Redis service
- MinIO service
- documented local commands

### Step 11: Observability

Goal: make the service easier to debug and operate.

Output:

- structured request logs
- daily log rotation or archival direction
- optional log storage in database or object storage
- metrics and tracing plan

### Step 12: Production Evolution

Goal: prepare for larger-scale deployment only after the modular monolith is
stable.

Output:

- Kong gateway direction
- gRPC direction for internal services
- Kubernetes or Helm direction
- Minikube practice setup
- service extraction strategy when needed

## Development Commands

Run the app:

```bash
go run ./cmd
```

Run with runtime flags later:

```bash
go run ./cmd --env=local --http-addr=:8080
```

Install Air, if needed:

```bash
go install github.com/air-verse/air@latest
```

Start live reload:

```bash
air
```

The Air config is in `.air.toml` and builds from `./cmd`.

## Definition of Done for Each Step

Before moving from one roadmap step to the next:

- the app should build successfully
- the changed behavior should be manually verified or covered by tests
- configuration should be documented without real secrets
- errors should be clear enough for local debugging
- new code should follow the project boundaries described above
- the README should be updated when commands or setup steps change

## Migration Direction

Database migrations will use `golang-migrate`.

Planned behavior:

- SQL migration files live in a dedicated migrations directory.
- Local development can run migrations manually.
- The application should fail clearly when required database setup is missing.
- Production migrations should be explicit and controlled.

## API Design Direction

The REST API should be designed carefully before gRPC or microservices are
added.

General API direction:

- versioned routes, for example `/api/v1`
- JSON request and response bodies
- consistent error response format
- request ID support
- tenant-aware requests
- rate-limited public endpoints
- health and readiness endpoints
- clear validation errors

## Security Notes

- Do not commit `.env` files with real values.
- Do not place real passwords or access keys in README examples.
- Prefer environment variables or secret managers for sensitive values.
- Use command-line flags only for non-secret options.
- Keep local-only configuration out of public commits.

## Learning Direction

This project should grow slowly so each part is understandable:

- Go project structure
- `context.Context`
- goroutines and graceful shutdown
- HTTP routing and middleware
- configuration and secret handling
- database connections and migrations
- tenant-aware data access
- caching with Redis
- object storage with MinIO
- Docker and deployment basics
- future microservice boundaries
