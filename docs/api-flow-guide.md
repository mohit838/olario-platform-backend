# How To Add One API

This is the step-by-step flow to add a new API in this project.

## 1. Start With The Business Meaning

Write down what the API does in business language.

Example:

```text
Create a product for a tenant.
The product must have a category, manual product code, price, stock threshold,
and optional barcode.
```

## 2. Check The Database

If the table or column does not exist, add a migration first.

Migration files live in:

```text
cmd/migrate/migrations
```

Run:

```bash
make migrate-up
```

## 3. Add Or Update Domain Types

Domain types live in:

```text
internal/domain
```

Use domain types for important business concepts such as tenant ID, product ID,
order status, and money.

## 4. Create The Request And Response DTOs

DTOs belong near HTTP handlers because they are transport shapes.

Recommended location:

```text
internal/http
```

DTOs should describe JSON input/output. They should not contain SQL logic.

## 5. Add An Application Service

Use cases live in:

```text
internal/application
```

The service should express the business flow and call interfaces such as:

```go
type ProductRepository interface {
    Create(ctx context.Context, input CreateProductInput) (Product, error)
}
```

## 6. Add A Repository Adapter

Postgres adapters live in:

```text
internal/platform/postgres
```

This is where SQL belongs. Keep tenant filtering in every query that touches
tenant-owned data.

## 7. Register The Handler In The Router

Routes live in:

```text
internal/http/router.go
```

Handlers should:

- decode request JSON
- call an application service
- convert the result to JSON
- return clear HTTP status codes

## 8. Wire Dependencies

Dependency wiring lives in:

```text
internal/app/app.go
```

Create repository adapters there, pass them into services, then pass services
into HTTP dependencies.

## 9. Verify

Run:

```bash
make test
make vet
```

If the API needs the database:

```bash
make migrate-up
make run-local
```

Then call the endpoint with `curl` or an API client.

## Current Teaching Example

The local-only teaching API is:

```http
POST /api/v1/dev/full-circle
```

It demonstrates:

- router registration
- handler
- application service
- Postgres transaction
- Redis counter/cache write
- audit log insert
- invoice-style JSON response

## Auth Flow

Current auth endpoints:

```http
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
```

Registration is invite-based. A tenant cannot randomly sign up without a valid
invitation token.

The local seeder creates one test invitation:

```text
email: owner@example.test
token: dev-invite-token
```

The local seeder also creates one test tenant, one tenant admin, two tenant
employees, and one platform superadmin. Seeded accounts use password `12345678`
and have 2FA disabled only for local testing.

The login flow is:

```text
HTTP handler
  -> auth application service
  -> Postgres user lookup
  -> bcrypt password check
  -> access token creation
  -> refresh token creation
  -> Redis refresh session storage
```

The refresh flow rotates tokens:

```text
client sends old refresh token
  -> Redis GETDEL removes old token session
  -> new access token is issued
  -> new refresh token is stored
  -> old refresh token can no longer be used
```
