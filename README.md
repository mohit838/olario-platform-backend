# Grocery Shop Platform

A generic starter repository for building a multi-tenant grocery shop platform.

## Tech Stack

- Backend: Go (REST API)
- Frontend: Angular
- Database: PostgreSQL
- Cache/Queue (optional): Redis

## Project Structure

- `backend/` Go backend service
- `frontend/` Angular frontend application

## Multi-Tenant Approach

The platform is designed to support multi-tenant workflows.

- Tenant context is expected through request metadata (example: `X-Tenant-ID`)
- Tenant resolution and access rules should be enforced in backend middleware
- Data isolation can be implemented via schema-per-tenant, row-level tenancy, or database-per-tenant

## Getting Started

### 1. Backend

```bash
cd backend
go mod init github.com/your-org/grocery-shop/backend
go mod tidy
```

Then create a simple `main.go` and start the API:

```bash
go run .
```

### 2. Frontend

```bash
cd frontend
ng new app --routing --style=scss
cd app
ng serve
```

## Environment Variables

Create environment files for backend and frontend (for example `.env`), and define:

- Database connection URL
- Redis connection URL (if used)
- API base URL for frontend
- Tenant-related feature flags (if needed)

## Development Notes

- Keep backend API contracts versioned (for example `/api/v1`)
- Validate tenant access on every protected request
- Use separate config values for local, staging, and production

## License

Add your preferred license here (MIT, Apache-2.0, or private/internal).
