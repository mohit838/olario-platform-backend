# Database Design

Last Updated: 2026-05-26

## Pattern

Database-per-service on PostgreSQL 17.
PostgreSQL is hosted from the local machine environment.

## Databases

- `olario_cms_db` for `cms-service`
- `olario_catalog_db` for `catalog-service`
- `olario_heavy_db` for `heavy-service`

## Rules

o_heavy_db`for`heavy-service`

## Rules_1

- No shared tables across services
- Migrations owned inside each service boundary
- Inter-service data access via APIs/events only
